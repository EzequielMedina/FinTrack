# Plan de Implementación: Chatbot Service con Ollama3

## Objetivo

- Implementar un microservicio `chatbot-service` que permita a los usuarios consultar sobre gastos, transacciones, cuentas, tarjetas y billetera.
- El servicio generará respuestas conversacionales, además de reportes descargables en PDF y visualizaciones en gráficos (Chart) para el frontend.
- Integrar con Ollama3 (corriendo en Docker) limitando su uso de memoria RAM para evitar consumo excesivo.
- Basar la arquitectura en el patrón ya utilizado por `notification-service` (Clean Architecture + Gin) para mantener consistencia.

## Alcance

- API HTTP para:
  - Conversación: `/api/chat/query` (texto → respuesta + acciones sugeridas)
  - Reporte PDF: `/api/chat/report/pdf` (parametrizable)
  - Datos para gráficos: `/api/chat/report/chart` (JSON listo para Chart.js)
- Acceso directo a la base MySQL compartida (mismo esquema) para obtener transacciones, cuentas, tarjetas y billetera, evitando integraciones HTTP entre microservicios.
- Provider de LLM (Ollama3) con memoria limitada vía Docker Compose.
- Sistema de prompts seguro y contextual que usa datos del usuario.
- Generación de PDF en backend (server-side) y datos de gráficos para renderizar en frontend.

## Arquitectura (basada en notification-service)

Estructura de carpetas propuesta:

```
internal/
  app/
    application.go            # Bootstrap de app (Gin, DI, config)
  config/
    config.go                 # Carga de variables de entorno
  core/
    domain/
      entities/
        report.go             # Estructuras de reportes (ingresos/gastos, por categoría, etc.)
        chart.go              # Modelos para gráficos (series, labels)
      valueobjects/
        query.go              # Consulta del usuario, filtros, periodo
    ports/
      chatbot_service.go      # Interfaz del servicio principal
      ollama_provider.go      # Interfaz al proveedor LLM
      report_provider.go      # Interfaz para generar reportes/PDF
      data_provider.go        # Interfaz para obtener datos (transactions/accounts/wallet)
    service/
      chatbot_service_impl.go # Implementación core del flujo conversacional
      report_service_impl.go  # Implementación de armado de reportes y gráficos
    providers/
      ollama/
        ollama_client.go      # Cliente HTTP a Ollama (chat/completions)
      data/
        mysql/
          connection.go       # Conexión MySQL (basado en notification-service)
          queries.go          # Consultas y agregaciones
      pdf/
        pdf_generator.go      # Generación de PDF (gofpdf, wkhtmltopdf, etc.)
  infrastructure/
    entrypoints/
      handlers/
        chat_handler.go       # Endpoints HTTP
    repositories/             # Si se requiere persistencia local
cmd/
  main.go                     # Wire-up servidor HTTP
Dockerfile
README.md
```

Principios:
- Clean Architecture: `core` no depende de `infrastructure`.
- Providers intercambiables (interfaces en `ports`).
- Configuración via ENV y `config.go`.

## Endpoints y Contratos

1) `POST /api/chat/query`
- Body:
```json
{
  "userId": "<uuid>",
  "message": "¿Cómo gasté este mes en supermercados?",
  "period": {"from": "2025-10-01", "to": "2025-10-31"},
  "filters": {"categories": ["supermarket"], "accounts": ["..."], "cards": ["..."]}
}
```
- Respuesta:
```json
{
  "reply": "Gastaste ARS 45,230 en supermercados. El pico fue el 12/10.",
  "suggestedActions": [
    {"type": "generate_pdf", "params": {"period": {"from": "2025-10-01", "to": "2025-10-31"}}},
    {"type": "show_chart", "params": {"chartType": "bar", "groupBy": "week"}}
  ],
  "insights": ["Mayor gasto en cadena X", "Promedio por compra ARS 5,033"],
  "dataRefs": {"transactionsCount": 122, "accountsUsed": 3}
}
```

2) `POST /api/chat/report/pdf`
- Body:
```json
{
  "userId": "<uuid>",
  "report": {
    "title": "Gastos por categoría",
    "period": {"from": "2025-10-01", "to": "2025-10-31"},
    "groupBy": "category",
    "includeCharts": true
  }
}
```
- Respuesta: `application/pdf` (stream / descarga) o URL firmado para descarga.

3) `POST /api/chat/report/chart`
- Body:
```json
{
  "userId": "<uuid>",
  "chart": {
    "type": "bar|line|pie",
    "period": {"from": "2025-10-01", "to": "2025-10-31"},
    "groupBy": "category|week|merchant",
    "currency": "ARS|USD"
  }
}
```
- Respuesta (para Chart.js):
```json
{
  "labels": ["Supermercado", "Transporte", "Comida"],
  "datasets": [{
    "label": "Gastos",
    "data": [45230, 12300, 21990],
    "backgroundColor": ["#3b82f6", "#ef4444", "#10b981"]
  }],
  "meta": {"currency": "ARS", "total": 79520}
}
```

## Flujo Conversacional y Contexto

1) El handler recibe la consulta del usuario.
2) `data_provider` agrega contexto: transacciones del período, categorías, comercios, cuentas y tarjetas.
3) Se construye un `system prompt` con reglas:
   - Solo responder con base en datos del usuario (evitar alucinaciones).
   - Resumir con métricas (total, promedio, top categorías/merchants).
   - Sugerir si conviene PDF o gráfico.
4) `ollama_provider` envía el prompt a Ollama3 (modelo `llama3:8b` por ejemplo).
5) Se post-procesa la respuesta para extraer `suggestedActions` y `insights` (schema JSON dentro del prompt).
6) Se devuelve la respuesta al frontend.

Prompt base (ejemplo simplificado):

```
SYSTEM: Eres un asistente financiero de FinTrack. Usa EXCLUSIVAMENTE los datos del usuario provistos en CONTEXTO.
Devuelve:
{ "reply": string, "suggestedActions": [ {"type": "generate_pdf"|"show_chart", "params": any } ], "insights": string[] }

CONTEXT:
- Periodo: {{from}} .. {{to}}
- Total gastado: {{total}}
- Top categorías: {{topCategories}}
- Comercios destacados: {{topMerchants}}
- Cuentas/tarjetas usadas: {{accountsCards}}

USER: {{message}}
```

## Integración con Ollama3 (Docker + límite de RAM)

Servicio en `docker-compose.yml` (ejemplo):

```yaml
services:
  ollama:
    image: ollama/ollama:latest
    container_name: ollama
    ports:
      - "11434:11434"
    volumes:
      - ollama:/root/.ollama
    command: ["serve"]
    environment:
      # Control de paralelismo y memoria interna
      OLLAMA_NUM_PARALLEL: 1
      OLLAMA_MAX_LOADED_MODELS: 1
    deploy:
      resources:
        limits:
          memory: 4G
        reservations:
          memory: 2G
    mem_limit: 4g       # Para Docker Compose clásico (no Swarm)
    mem_reservation: 2g
    restart: unless-stopped

  chatbot-service:
    build: ./backend/services/chatbot-service
    ports:
      - "8090:8090"
    environment:
      OLLAMA_HOST: http://ollama:11434
      OLLAMA_MODEL: llama3:8b
      CHATBOT_ENABLED: "true"
      # Config DB (MySQL compartida)
      DB_HOST: mysql
      DB_PORT: 3306
      DB_NAME: fintrack
      DB_USER: fintrack
      DB_PASSWORD: fintrack_pwd
    depends_on:
      - ollama
      - mysql
```

Notas:
- `mem_limit` y `mem_reservation` limitan la RAM del contenedor en Compose clásico.
- `deploy.resources.limits.memory` aplica en Swarm; lo mantenemos por compatibilidad futura.
- Ajustar el modelo (`llama3:8b`) según hardware; modelos más grandes requieren más memoria.

## Configuración (`internal/config/config.go`)

- Variables:
  - `CHATBOT_ENABLED` (bool)
  - `OLLAMA_HOST` (string, default `http://localhost:11434`)
  - `OLLAMA_MODEL` (string, default `llama3:8b`)
  - `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`
  - `REPORT_PDF_ENGINE` (`gofpdf|wkhtmltopdf|chromedp`)
  - `MAX_QUERY_PERIOD_DAYS` (int, p.ej. 366)
  - `TIMEZONE` (string, default `America/Argentina/Buenos_Aires`)

## Acceso a Datos (MySQL compartida)

Basado en `notification-service`:

```go
// providers/data/mysql/connection.go
dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, name)
db, _ := sql.Open("mysql", dsn)
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(25)
db.SetConnMaxLifetime(5 * time.Minute)
```

Consultas típicas:

1) Totales por periodo (gastos/ingresos):
```sql
SELECT 
  SUM(CASE WHEN type IN ('debit_purchase','credit_charge','wallet_withdrawal','account_withdraw') THEN amount ELSE 0 END) AS total_gastos,
  SUM(CASE WHEN type IN ('wallet_deposit','account_deposit','credit_payment','installment_payment') THEN amount ELSE 0 END) AS total_ingresos
FROM transactions
WHERE user_id = ?
  AND status = 'completed'
  AND created_at BETWEEN ? AND ?;
```

2) Gastos por `type`:
```sql
SELECT type, SUM(amount) AS total
FROM transactions
WHERE user_id = ?
  AND status = 'completed'
  AND created_at BETWEEN ? AND ?
GROUP BY type
ORDER BY total DESC;
```

3) Gastos por `merchant_name`:
```sql
SELECT merchant_name, SUM(amount) AS total
FROM transactions
WHERE user_id = ?
  AND status = 'completed'
  AND merchant_name IS NOT NULL AND merchant_name <> ''
  AND created_at BETWEEN ? AND ?
GROUP BY merchant_name
ORDER BY total DESC
LIMIT 15;
```

4) Por cuenta/tarjeta:
```sql
-- Por cuenta
SELECT a.account_type, SUM(t.amount) AS total
FROM transactions t
JOIN accounts a ON (
  t.from_account_id = a.id OR t.to_account_id = a.id
)
WHERE t.user_id = ?
  AND t.status = 'completed'
  AND t.created_at BETWEEN ? AND ?
GROUP BY a.account_type;

-- Por tarjeta (últimos 4 dígitos)
SELECT c.card_brand, c.last_four_digits, SUM(t.amount) AS total
FROM transactions t
JOIN cards c ON (
  t.from_card_id = c.id OR t.to_card_id = c.id
)
WHERE t.user_id = ?
  AND t.status = 'completed'
  AND t.created_at BETWEEN ? AND ?
GROUP BY c.card_brand, c.last_four_digits
ORDER BY total DESC;
```

Índices recomendados:
- `transactions(user_id, created_at)`
- `transactions(type)`
- `transactions(merchant_name)`
- `transactions(status)`

## Seguridad y Privacidad

## Seguridad y Privacidad

- Autenticación JWT en el gateway / middleware; `userId` validado en cada request.
- Data provider filtra por `userId` en todas las consultas SQL.
- Sanitización de prompts: no enviar datos sensibles (números completos de tarjeta; usar enmascaramiento).
- Rate limiting por IP/usuario para `/api/chat/query`.
- Logging estructurado sin datos sensibles.
- CORS configurado para el frontend.

## Generación de PDF

Opciones:
- `gofpdf`: rápido, sin dependencias externas (bueno para texto + tablas, gráficos embebidos como PNG).
- `wkhtmltopdf` / `chromedp`: renderiza HTML a PDF (ideal para templates con CSS/Chart como imagen).

Estrategia:
- Renderizar un template HTML del reporte (título, periodo, KPIs, tabla de gastos) y generar el PDF.
- Incluir gráficos como imagen: generar PNG en backend (con `go-echarts`) o recibir del frontend y embeberlo.
- Retornar descarga directa (`application/pdf`) o URL temporal firmado.

## Gráficos (Chart)

- El backend devuelve datos normalizados para Chart.js: `labels`, `datasets`, `meta`.
- El frontend renderiza el gráfico con Chart.js (bar/line/pie) y permite descarga (imagen) o exportación a PDF.
- Alternativa server-side: `go-echarts` para generar PNG y devolver recurso estático (útil para PDFs).

## Integración con Servicios Existentes (Opcional)

Si se requiere lógica adicional, puede integrarse con servicios vía HTTP. Por defecto, el `chatbot-service` leerá directamente de MySQL compartida y realizará las agregaciones necesarias. Para conversión de moneda se puede consultar `exchange-service` si no hay cache local.

## Ejemplo de Flujo End-to-End

1) Usuario escribe: "Comparame mis gastos de octubre vs septiembre y generá un PDF".
2) Frontend llama `/api/chat/query` con periodo [sep, oct].
3) `chatbot-service` arma contexto y consulta a Ollama3; devuelve respuesta con `suggestedActions=[generate_pdf]`.
4) Frontend invoca `/api/chat/report/pdf` con los parámetros sugeridos.
5) Backend genera PDF y retorna descarga.
6) Frontend opcionalmente llama `/api/chat/report/chart` para visualizar el comparativo.

## Testing

- Unit tests en `core/service` y `providers` (mocks de Ollama y data provider).
- Endpoint tests (Gin) para `/query`, `/report/pdf`, `/report/chart`.
- Pruebas de límites: periodo máximo, filtros, performance con dataset grande.
- Integración: Compose con `ollama` + servicios stub o reales.

## Observabilidad

- Logs estructurados (nivel info/debug/error) por request-id.
- Métricas Prometheus (latencia, éxito/fracaso, memoria usada).
- Health check: `/health` y readiness cuando `OLLAMA_HOST` responde.

## Roadmap de Implementación

1) Bootstrap del proyecto (clonar estructura de notification-service).
2) Configuración y Dockerfile del `chatbot-service`.
3) Cliente Ollama (`providers/ollama`).
4) Data provider (clients a transaction/account/wallet).
5) Servicio conversacional (`core/service/chatbot_service_impl.go`).
6) Servicio de reportes y generación de PDF.
7) Handlers HTTP y contratos.
8) Integración con Compose (servicio `ollama` con límites de RAM).
9) Pruebas unitarias e integración.
10) Documentación y ejemplos de uso.

## Ejemplos de Configuración de Entorno

```
CHATBOT_ENABLED=true
OLLAMA_HOST=http://localhost:11434
OLLAMA_MODEL=llama3:8b
TRANSACTION_SERVICE_URL=http://localhost:8083
ACCOUNT_SERVICE_URL=http://localhost:8082
WALLET_SERVICE_URL=http://localhost:8080
REPORT_PDF_ENGINE=gofpdf
MAX_QUERY_PERIOD_DAYS=366
TIMEZONE=America/Argentina/Buenos_Aires
DB_HOST=mysql
DB_PORT=3306
DB_NAME=fintrack
DB_USER=fintrack
DB_PASSWORD=fintrack_pwd
```

## Consideraciones de Rendimiento

- Limitar paralelismo de Ollama y modelos cargados (`OLLAMA_NUM_PARALLEL=1`, `OLLAMA_MAX_LOADED_MODELS=1`).
- Cache de resultados de agregación por periodo y filtro.
- Paginación y límites en consultas a `transaction-service`.
- Streaming de respuesta (SSE/WebSocket) opcional para UX.

## Seguridad en Prompts

- Nunca incluir datos sensibles (números de tarjeta completos); usar `masked_number` y `last_four_digits`.
- Evitar prompts abiertos; usar schema JSON en la respuesta para parseo confiable.
- Validar salida del modelo (tipos y rangos) antes de exponer al usuario.

## Próximos Pasos

- Crear esqueletos de archivos indicados.
- Agregar el servicio `ollama` a `docker-compose.yml` con límites de memoria.
- Implementar primer flujo `/api/chat/query` con contexto simple (por mes actual).
- Iterar en reportes PDF y gráfico por categoría.