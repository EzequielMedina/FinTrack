# 🤖 Análisis y Diagnóstico del Chatbot-Service FinTrack

## 📊 Estado Actual del Sistema

### ✅ Servicios Activos
- **Chatbot Service**: Puerto 8090 (Saludable)
- **Ollama**: Puerto 11434 (Activo)
- **MySQL**: Puerto 3306 (Saludable)
- **Account Service**: Puerto 8082 (Saludable)
- **Transaction Service**: Puerto 8083 (Saludable)

---

## 🔍 Problemas Identificados

### 1. **Problema Principal: Configuración de Ollama**

#### 🚨 Issue Crítico: Modelo No Descargado
```bash
# El contenedor Ollama está corriendo pero probablemente no tiene el modelo llama3:8b descargado
# Esto causa que las llamadas a /api/chat fallen o retornen vacío
```

**Solución:**
```bash
# Entrar al contenedor de Ollama y descargar el modelo
docker exec -it fintrack-ollama ollama pull llama3:8b

# O si prefieres un modelo más ligero:
docker exec -it fintrack-ollama ollama pull llama3.2:3b
```

### 2. **Problemas de Configuración en Docker Compose**

#### ❌ Configuración Actual Faltante
Tu `docker-compose.yml` necesita estas mejoras:

```yaml
# Agregar al docker-compose.yml
ollama:
  image: ollama/ollama:latest
  container_name: fintrack-ollama
  ports:
    - "11434:11434"
  environment:
    OLLAMA_NUM_PARALLEL: 1          # Limitar concurrencia
    OLLAMA_MAX_LOADED_MODELS: 1     # Limitar modelos cargados
    OLLAMA_HOST: 0.0.0.0           # Permitir conexiones externas
  deploy:
    resources:
      limits:
        memory: 4G                  # Limitar RAM como solicitaste
      reservations:
        memory: 2G
  volumes:
    - ollama_models:/root/.ollama   # Persistir modelos descargados
  networks:
    - fintrack-network
  healthcheck:
    test: ["CMD", "curl", "-f", "http://localhost:11434/api/tags"]
    interval: 30s
    timeout: 10s
    retries: 3

chatbot-service:
  build:
    context: ./backend/services/chatbot-service
    dockerfile: Dockerfile
  container_name: fintrack-chatbot-service
  environment:
    DB_HOST: mysql
    DB_PORT: 3306
    DB_NAME: fintrack
    DB_USER: fintrack_user
    DB_PASSWORD: fintrack_password
    OLLAMA_HOST: http://ollama:11434    # ⚠️ CRÍTICO: Usar nombre del servicio
    OLLAMA_MODEL: llama3:8b            # O llama3.2:3b para menos RAM
    PORT: 8090
    GIN_MODE: release
  ports:
    - "8090:8090"
  depends_on:
    mysql:
      condition: service_healthy
    ollama:
      condition: service_started        # Esperar que Ollama esté disponible
  networks:
    - fintrack-network

volumes:
  ollama_models:                        # Volumen para persistir modelos
    driver: local
```

### 3. **Problemas en el Código**

#### 🔧 Fix en `ollama_client.go`
```go
// Agregar mejor manejo de errores y logs
func (c *Client) Chat(ctx context.Context, systemPrompt string, userPrompt string) (string, error) {
    // ⚠️ AGREGAR: Log para debugging
    log.Printf("🤖 Ollama request to %s with model %s", c.host, c.model)
    
    payload := map[string]any{
        "model": c.model,
        "messages": []map[string]string{
            {"role": "system", "content": systemPrompt},
            {"role": "user", "content": userPrompt},
        },
        "stream": false,
        "options": map[string]any{
            "temperature": 0.7,
            "top_p": 0.9,
            "num_ctx": 2048,  // Reducir contexto para usar menos RAM
        },
    }
    
    b, _ := json.Marshal(payload)
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/api/chat", c.host), bytes.NewReader(b))
    if err != nil { 
        log.Printf("❌ Error creating request: %v", err)
        return "", err 
    }
    
    req.Header.Set("Content-Type", "application/json")
    resp, err := c.http.Do(req)
    if err != nil { 
        log.Printf("❌ Error calling Ollama: %v", err)
        return "", err 
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK { 
        log.Printf("❌ Ollama returned status %d", resp.StatusCode)
        return "", fmt.Errorf("ollama status %d", resp.StatusCode) 
    }
    
    var out struct {
        Message struct{ Content string `json:"content"` } `json:"message"`
        Response string `json:"response"`
        Error    string `json:"error,omitempty"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&out); err != nil { 
        log.Printf("❌ Error decoding response: %v", err)
        return "", err 
    }
    
    if out.Error != "" {
        log.Printf("❌ Ollama error: %s", out.Error)
        return "", fmt.Errorf("ollama error: %s", out.Error)
    }
    
    result := out.Message.Content
    if result == "" {
        result = out.Response
    }
    
    log.Printf("✅ Ollama response length: %d chars", len(result))
    return result, nil
}
```

#### 🔧 Fix en `chatbot_service_impl.go`
```go
// Mejorar el fallback cuando Ollama falla
func (s *ChatbotServiceImpl) HandleQuery(ctx context.Context, req ports.ChatQueryRequest) (ports.ChatQueryResponse, error) {
    // ... código existente ...
    
    // ⚠️ MEJORAR: Llamada a Ollama con mejor fallback
    r, err := s.llm.Chat(ctx, system, user)
    if err != nil {
        // Log del error pero no fallar completamente
        log.Printf("⚠️ Ollama falló, usando respuesta básica: %v", err)
        r = "" // Forzar uso del fallback
    }
    
    // ⚠️ MEJORAR: Fallback más inteligente
    if r == "" || strings.TrimSpace(r) == "" {
        reply = s.generateBasicReply(totals, byType, topMerchants, req.Message)
    } else { 
        reply = r 
    }
    
    // ... resto del código ...
}

// Nuevo método para respuesta básica inteligente
func (s *ChatbotServiceImpl) generateBasicReply(totals ports.FinancialTotals, byType map[string]float64, merchants []ports.MerchantTotal, message string) string {
    msg := strings.ToLower(message)
    
    // Detectar intención de la pregunta
    if strings.Contains(msg, "gasto") || strings.Contains(msg, "gastado") {
        return fmt.Sprintf("📊 **Resumen de Gastos**\\n\\nTotal gastado: **$%.2f**\\nTotal ingresado: **$%.2f**\\n\\n🏪 **Principales comercios:**\\n%s\\n\\n💡 *Tip: Puedes generar un reporte PDF para ver más detalles.*", 
            totals.Expenses, totals.Incomes, formatMerchantsMarkdown(merchants))
    }
    
    if strings.Contains(msg, "ingreso") || strings.Contains(msg, "ganancia") {
        return fmt.Sprintf("💰 **Resumen de Ingresos**\\n\\nTotal ingresado: **$%.2f**\\nTotal gastado: **$%.2f**\\nBalance: **$%.2f**", 
            totals.Incomes, totals.Expenses, totals.Incomes-totals.Expenses)
    }
    
    if strings.Contains(msg, "tarjeta") {
        ccCharges := getVal(byType, "credit_charge")
        return fmt.Sprintf("💳 **Resumen de Tarjetas**\\n\\nConsumos con tarjeta de crédito: **$%.2f**\\n\\n💡 *Tip: Puedes ver un gráfico por tarjeta usando 'Mostrar gráfico'.*", ccCharges)
    }
    
    // Respuesta general
    return fmt.Sprintf("📈 **Resumen Financiero**\\n\\n• Total gastado: **$%.2f**\\n• Total ingresado: **$%.2f**\\n• Balance: **$%.2f**\\n\\n🏪 **Top comercios:** %s\\n\\n💡 *¿Te gustaría generar un reporte PDF o ver gráficos de tus movimientos?*",
        totals.Expenses, totals.Incomes, totals.Incomes-totals.Expenses, formatMerchantsSimple(merchants))
}
```

---

## 🧪 Pruebas y Diagnóstico

### 1. **Verificar Estado de Ollama**

```bash
# Test 1: Verificar que Ollama responde
curl http://localhost:11434/api/tags

# Test 2: Verificar modelos instalados
docker exec fintrack-ollama ollama list

# Test 3: Descargar modelo si no existe
docker exec fintrack-ollama ollama pull llama3.2:3b

# Test 4: Probar chat directo
curl http://localhost:11434/api/chat -d '{
  "model": "llama3.2:3b",
  "messages": [{"role": "user", "content": "Hola, responde en español"}],
  "stream": false
}'
```

### 2. **Verificar Chatbot Service**

```bash
# Test 1: Health check
curl http://localhost:8090/health

# Test 2: Consulta básica
curl -X POST http://localhost:8090/api/chat/query \\
  -H "Content-Type: application/json" \\
  -d '{
    "userId": "test-user-123",
    "message": "¿Cuánto gasté este mes?",
    "period": {"from": "2025-10-01", "to": "2025-10-31"}
  }'

# Test 3: Verificar logs del contenedor
docker logs fintrack-chatbot-service -f
```

### 3. **Verificar Datos en Base de Datos**

```sql
-- Conectar a MySQL y verificar que hay datos de prueba
SELECT COUNT(*) FROM transactions WHERE user_id = 'test-user-123';
SELECT COUNT(*) FROM accounts WHERE user_id = 'test-user-123';
```

---

## ⚡ Solución Rápida (Quick Fix)

### 1. **Recrear con Configuración Correcta**

```bash
# 1. Detener servicios
docker-compose down

# 2. Agregar configuración de Ollama al docker-compose.yml (ver arriba)

# 3. Reiniciar todo
docker-compose up --build -d

# 4. Descargar modelo Ollama
docker exec fintrack-ollama ollama pull llama3.2:3b

# 5. Verificar que funciona
curl http://localhost:8090/health
```

### 2. **Configuración Alternativa (Modelo Más Liviano)**

Si tienes problemas de RAM, usa un modelo más pequeño:

```yaml
# En docker-compose.yml
chatbot-service:
  environment:
    OLLAMA_MODEL: llama3.2:1b  # Modelo más pequeño (1B parámetros)
    # O usar qwen2.5:0.5b para aún menos RAM
```

---

## 🔧 Mejoras Recomendadas

### 1. **Sistema de Fallback Robusto**

```go
// Implementar en chatbot_service_impl.go
type ChatbotServiceImpl struct {
    data   ports.DataProvider
    llm    ports.OllamaProvider
    report ports.ReportProvider
    fallbackEnabled bool  // ⚠️ NUEVO: Flag para activar fallback
}

func (s *ChatbotServiceImpl) HandleQuery(ctx context.Context, req ports.ChatQueryRequest) (ports.ChatQueryResponse, error) {
    // Intentar Ollama primero
    reply, err := s.tryOllamaQuery(ctx, req)
    if err != nil && s.fallbackEnabled {
        // Usar respuesta estructurada como fallback
        reply = s.generateStructuredReply(ctx, req)
    }
    
    return ports.ChatQueryResponse{
        Reply: reply,
        SuggestedActions: s.generateSuggestedActions(req),
        Insights: s.generateInsights(ctx, req),
        DataRefs: s.gatherDataRefs(ctx, req),
    }, nil
}
```

### 2. **Configuración de Environment**

```bash
# Agregar al .env del chatbot-service
CHATBOT_FALLBACK_ENABLED=true
CHATBOT_MAX_CONTEXT_LENGTH=2048
CHATBOT_RESPONSE_TIMEOUT=30s
OLLAMA_TEMPERATURE=0.7
OLLAMA_TOP_P=0.9
```

### 3. **Health Check Mejorado**

```go
// En chat_handler.go
func (h *ChatHandler) Health(c *gin.Context) {
    // Verificar Ollama
    ollamaStatus := "unknown"
    if err := h.svc.TestOllamaConnection(); err == nil {
        ollamaStatus = "healthy"
    } else {
        ollamaStatus = "unhealthy: " + err.Error()
    }
    
    c.JSON(http.StatusOK, gin.H{
        "status": "healthy", 
        "service": "chatbot-service",
        "ollama": ollamaStatus,
        "timestamp": time.Now(),
    })
}
```

---

## 🎯 Pasos Inmediatos para Solucionar

### ✅ Checklist de Solución

1. **[ ] Verificar modelo Ollama descargado**
   ```bash
   docker exec fintrack-ollama ollama list
   ```

2. **[ ] Descargar modelo si falta**
   ```bash
   docker exec fintrack-ollama ollama pull llama3.2:3b
   ```

3. **[ ] Verificar configuración de red**
   - Asegurar que chatbot-service use `http://ollama:11434` (no localhost)

4. **[ ] Agregar logs de debugging**
   - Implementar logs en ollama_client.go como se muestra arriba

5. **[ ] Probar endpoint directamente**
   ```bash
   curl -X POST http://localhost:8090/api/chat/query -H "Content-Type: application/json" -d '{"userId":"test","message":"¿Cuánto gasté?","period":{"from":"2025-10-01","to":"2025-10-31"}}'
   ```

6. **[ ] Verificar logs del contenedor**
   ```bash
   docker logs fintrack-chatbot-service -f
   ```

### 🔥 Si Todo Falla - Modo Debugging

```bash
# 1. Entrar al contenedor del chatbot
docker exec -it fintrack-chatbot-service sh

# 2. Verificar conectividad a Ollama
curl http://ollama:11434/api/tags

# 3. Verificar variables de entorno
env | grep OLLAMA

# 4. Verificar conectividad a MySQL
nc -zv mysql 3306
```

---

## 💡 Próximos Pasos

Una vez que el chatbot funcione básicamente, te recomiendo:

1. **Implementar cache** para respuestas frecuentes
2. **Agregar sistema de contexto** para conversaciones multi-turno
3. **Implementar rate limiting** para evitar sobrecarga
4. **Agregar métricas** de performance y uso
5. **Crear interfaz de chat** en el frontend Angular

¿Te gustaría que implemente alguna de estas soluciones específicas o prefieres que revisemos los logs para diagnosticar el problema actual?