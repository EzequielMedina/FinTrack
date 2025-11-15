# Chatbot Service - FinTrack

## Descripción

Microservicio encargado del asistente virtual inteligente en FinTrack. Permite consultas sobre gastos, transacciones, cuentas, tarjetas y billetera; genera reportes PDF y datos para gráficos, usando IA vía Ollama.

## Tecnologías

- Lenguaje: Go 1.24+
- Framework: Gin (HTTP Router)
- Base de Datos: MySQL 8.0 (compartida)
- IA: Ollama (modelo configurable)
- Contenedor: Docker multi-stage
- Arquitectura: Clean Architecture

## Variables de Entorno

```env
# Base de datos
DB_HOST=mysql
DB_PORT=3306
DB_NAME=fintrack
DB_USER=fintrack_user
DB_PASSWORD=fintrack_password

# Ollama
OLLAMA_HOST=http://localhost:11434
OLLAMA_MODEL=llama3:8b

# Servidor
PORT=8090
GIN_MODE=release
TIMEZONE=America/Argentina/Buenos_Aires

# Reportes
REPORT_PDF_ENGINE=gofpdf
```

## API Endpoints

```http
POST /api/chat/query        # Consulta conversacional con contexto financiero
POST /api/chat/report/pdf   # Generación de reporte PDF
POST /api/chat/report/chart # Datos listos para Chart.js
GET  /health                # Health check
```

### Ejemplos

```bash
curl -X POST http://localhost:8090/api/chat/query \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "user-123",
    "message": "¿Cómo gasté este mes en supermercados?",
    "period": {"from": "2025-10-01", "to": "2025-10-31"}
  }'

curl -X POST http://localhost:8090/api/chat/report/pdf \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "user-123",
    "title": "Gastos por categoría",
    "period": {"from": "2025-10-01", "to": "2025-10-31"},
    "groupBy": "category",
    "includeCharts": true
  }' --output reporte.pdf

curl http://localhost:8090/health
```

## Desarrollo Local

```bash
cd backend/services/chatbot-service
go mod download
go run cmd/main.go
```

## Docker Compose

```yaml
ollama:
  image: ollama/ollama:latest
  ports: ["11434:11434"]
  environment:
    OLLAMA_NUM_PARALLEL: 1
    OLLAMA_MAX_LOADED_MODELS: 1
  mem_limit: 4g
  mem_reservation: 2g

chatbot-service:
  build: ./backend/services/chatbot-service
  ports: ["8090:8090"]
  environment:
    DB_HOST: mysql
    DB_PORT: 3306
    DB_NAME: fintrack
    DB_USER: fintrack_user
    DB_PASSWORD: fintrack_password
    OLLAMA_HOST: http://ollama:11434
    OLLAMA_MODEL: llama3:8b
    REPORT_PDF_ENGINE: gofpdf
  depends_on:
    - mysql
    - ollama
```

## Capacidades

- Consultas sobre balances y transacciones por periodo
- Análisis de gastos por tipo, comercio, cuenta y tarjeta
- Sugerencias de acciones (generar PDF, ver gráficos)
- Datos agregados listos para Chart.js

## Prompts base (Ollama)

```text
SYSTEM: Eres un asistente financiero de FinTrack. Usa SOLO el contexto provisto. Devuelve resumen, insights y acciones sugeridas.
USER: {mensaje del usuario}
CONTEXT: {totales, por tipo, top comercios, etc.}
```