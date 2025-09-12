# Chatbot Service - FinTrack

## 📋 Descripción

Microservicio encargado del asistente virtual inteligente en la plataforma FinTrack. Proporciona soporte automatizado, consultas financieras y asistencia personalizada usando IA.

## 🛠️ Tecnologías

- **Lenguaje**: Go 1.24+
- **Framework**: Gin/Echo (HTTP Router)
- **Base de Datos**: MySQL 8.0
- **IA**: OpenAI GPT API
- **WebSockets**: Comunicación en tiempo real
- **Contenedor**: Docker multi-stage
- **Arquitectura**: Clean Architecture

## 🏗️ Arquitectura

### Estructura del Proyecto

```
chatbot-service/
├── cmd/
│   └── main.go              # Punto de entrada
├── internal/
│   ├── config/              # Configuración
│   ├── domain/              # Entidades de dominio
│   ├── handlers/            # HTTP handlers
│   ├── repository/          # Capa de datos
│   ├── service/             # Lógica de negocio
│   ├── ai/                  # Integración con IA
│   ├── websocket/           # WebSocket handlers
│   └── middleware/          # Middlewares HTTP
├── prompts/                 # Prompts de IA
├── Dockerfile               # Configuración Docker
├── go.mod                   # Dependencias Go
├── go.sum                   # Checksums de dependencias
└── README.md                # Este archivo
```

## 🚀 Desarrollo Local

### Variables de Entorno

```env
# Base de datos
DB_HOST=localhost
DB_PORT=3306
DB_NAME=fintrack
DB_USER=fintrack_user
DB_PASSWORD=fintrack_password

# Servicios externos
USER_SERVICE_URL=http://localhost:8081

# OpenAI
OPENAI_API_KEY=your-openai-api-key
OPENAI_MODEL=gpt-4
OPENAI_MAX_TOKENS=1000
OPENAI_TEMPERATURE=0.7

# WebSocket
WS_READ_BUFFER_SIZE=1024
WS_WRITE_BUFFER_SIZE=1024
WS_MAX_MESSAGE_SIZE=512

# Servidor
PORT=8080
GIN_MODE=debug

# Logging
LOG_LEVEL=info
```

### Comandos de Desarrollo

```bash
# Navegar al servicio
cd backend/services/chatbot-service

# Instalar dependencias
go mod download

# Ejecutar en modo desarrollo
go run cmd/main.go

# Build del binario
go build -o bin/chatbot-service cmd/main.go

# Tests
go test ./...
```

## 🐳 Docker

```bash
# Build de la imagen
docker build -t fintrack-chatbot-service .

# Docker Compose
docker-compose up chatbot-service

# Con dependencias
docker-compose up mysql user-service chatbot-service
```

## 📡 API Endpoints

### Chat HTTP

```http
POST   /api/chat/message              # Enviar mensaje
GET    /api/chat/conversations        # Listar conversaciones
GET    /api/chat/conversations/{id}   # Obtener conversación
DELETE /api/chat/conversations/{id}   # Eliminar conversación
```

### WebSocket

```http
GET    /api/chat/ws                   # Conexión WebSocket
```

### Configuración del Bot

```http
GET    /api/chat/settings             # Configuración del chatbot
PUT    /api/chat/settings             # Actualizar configuración
GET    /api/chat/prompts              # Listar prompts disponibles
PUT    /api/chat/prompts/{id}         # Actualizar prompt
```

### Análisis

```http
GET    /api/chat/analytics            # Métricas del chatbot
GET    /api/chat/feedback             # Feedback de usuarios
POST   /api/chat/feedback             # Enviar feedback
```

### Health Check

```http
GET /health
```

### Ejemplos de Uso

```bash
# Enviar mensaje al chatbot
curl -X POST http://localhost:8086/api/chat/message \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "message": "¿Cuál es mi balance actual?",
    "conversationId": "conv_123"
  }'

# Obtener conversaciones
curl -X GET http://localhost:8086/api/chat/conversations \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# WebSocket connection (JavaScript)
const ws = new WebSocket("ws://localhost:8086/api/chat/ws?token=YOUR_JWT_TOKEN");
ws.onmessage = function(event) {
    const response = JSON.parse(event.data);
    console.log("Bot response:", response.message);
};
ws.send(JSON.stringify({
    "message": "Hola, necesito ayuda con mis finanzas"
}));

# Health check
curl http://localhost:8086/health
```

## 🤖 Capacidades del Chatbot

### Consultas Financieras

- **Balance de cuentas**: "¿Cuál es mi balance?"
- **Historial de transacciones**: "Muéstrame mis últimas transacciones"
- **Análisis de gastos**: "¿En qué he gastado más este mes?"
- **Presupuestos**: "¿Cómo va mi presupuesto?"
- **Inversiones**: "¿Cómo están mis inversiones?"

### Operaciones

- **Transferencias**: "Transfiere $100 a mi cuenta de ahorros"
- **Pagos**: "Paga la factura de electricidad"
- **Recordatorios**: "Recuérdame pagar el alquiler el día 1"

### Soporte

- **Preguntas frecuentes**: Respuestas automáticas
- **Guías**: Tutoriales paso a paso
- **Escalación**: Transferencia a soporte humano

## 🧠 Configuración de IA

### System Prompts

```go
const (
    SystemPrompt = `Eres un asistente financiero virtual para FinTrack.
    Ayudas a los usuarios con:
    - Consultas sobre balances y transacciones
    - Análisis de gastos e ingresos
    - Consejos financieros básicos
    - Navegación por la plataforma
    
    Siempre sé útil, preciso y mantén un tono profesional pero amigable.
    Si no puedes realizar una acción, explica cómo el usuario puede hacerlo.`
    
    FinancialAnalysisPrompt = `Analiza los datos financieros del usuario y proporciona:
    - Resumen de gastos por categoría
    - Tendencias de ahorro
    - Recomendaciones personalizadas
    - Alertas sobre gastos inusuales`
)
```

### Context Management

```go
type ConversationContext struct {
    UserID          string
    ConversationID  string
    LastMessages    []Message
    UserProfile     UserProfile
    FinancialData   FinancialSummary
    SessionData     map[string]interface{}
}
```

## 🔐 Seguridad

### Medidas Implementadas

- **JWT Authentication**: Validación de tokens
- **Data Privacy**: No almacenamiento de datos sensibles en logs
- **Rate Limiting**: Limitación de mensajes por usuario
- **Content Filtering**: Filtrado de contenido inapropiado
- **Audit Trail**: Registro de conversaciones
- **API Key Protection**: Protección de claves de OpenAI

### Privacidad de Datos

```go
// Sanitización de datos antes de enviar a OpenAI
func sanitizeFinancialData(data FinancialData) FinancialData {
    // Remover números de cuenta, SSN, etc.
    // Mantener solo datos agregados y categorías
    return sanitizedData
}
```

## 🧪 Testing

```bash
# Tests unitarios
go test ./internal/...

# Tests de integración
go test ./tests/integration/...

# Tests de WebSocket
go test ./internal/websocket/...

# Tests de IA (con mocks)
go test ./internal/ai/...
```

### Mock de OpenAI

```go
type MockOpenAIClient struct {
    responses map[string]string
}

func (m *MockOpenAIClient) CreateCompletion(prompt string) (string, error) {
    if response, exists := m.responses[prompt]; exists {
        return response, nil
    }
    return "I'm a mock response", nil
}
```

## 📊 Monitoreo

### Métricas Específicas

- **Messages Processed**: Mensajes procesados
- **Response Time**: Tiempo de respuesta
- **User Satisfaction**: Satisfacción del usuario
- **Conversation Length**: Duración de conversaciones
- **API Usage**: Uso de API de OpenAI
- **Error Rate**: Tasa de errores
- **Active Connections**: Conexiones WebSocket activas

### Analytics Dashboard

```json
{
  "daily_metrics": {
    "messages_sent": 1250,
    "conversations_started": 89,
    "avg_response_time": "1.2s",
    "user_satisfaction": 4.2,
    "api_cost": "$12.45"
  }
}
```

## 🚀 Despliegue

### Variables de Producción

```env
GIN_MODE=release
LOG_LEVEL=warn
OPENAI_MAX_TOKENS=500
OPENAI_TEMPERATURE=0.5
WS_MAX_CONNECTIONS=1000
RATE_LIMIT_PER_MINUTE=30
```

```bash
# Build de producción
CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-w -s" \
  -o chatbot-service cmd/main.go
```

## 🔧 Configuración Avanzada

### Fine-tuning del Modelo

```bash
# Preparar datos de entrenamiento
go run scripts/prepare-training-data.go

# Subir a OpenAI para fine-tuning
curl -X POST "https://api.openai.com/v1/fine-tuning/jobs" \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "training_file": "file-abc123",
    "model": "gpt-3.5-turbo"
  }'
```

---

**Chatbot Service** - Asistente virtual inteligente 🤖💬