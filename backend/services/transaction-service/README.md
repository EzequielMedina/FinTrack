# Transaction Service - FinTrack

## 📋 Descripción

Microservicio encargado del procesamiento y gestión de transacciones financieras en la plataforma FinTrack. Maneja transferencias, pagos, historial de transacciones y validaciones de seguridad.

## 🛠️ Tecnologías

- **Lenguaje**: Go 1.24+
- **Framework**: Gin/Echo (HTTP Router)
- **Base de Datos**: MySQL 8.0
- **Comunicación**: HTTP REST APIs
- **Contenedor**: Docker multi-stage
- **Arquitectura**: Clean Architecture

## 🏗️ Arquitectura

### Estructura del Proyecto

```
transaction-service/
├── cmd/
│   └── main.go              # Punto de entrada
├── internal/
│   ├── config/              # Configuración
│   ├── domain/              # Entidades de dominio
│   ├── handlers/            # HTTP handlers
│   ├── repository/          # Capa de datos
│   ├── service/             # Lógica de negocio
│   └── middleware/          # Middlewares HTTP
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
WALLET_SERVICE_URL=http://localhost:8083

# Servidor
PORT=8080
GIN_MODE=debug

# Logging
LOG_LEVEL=info
```

### Comandos de Desarrollo

```bash
# Navegar al servicio
cd backend/services/transaction-service

# Instalar dependencias
go mod download

# Ejecutar en modo desarrollo
go run cmd/main.go

# Build del binario
go build -o bin/transaction-service cmd/main.go

# Tests
go test ./...

# Tests con coverage
go test -cover ./...
```

## 🐳 Docker

```bash
# Build de la imagen
docker build -t fintrack-transaction-service .

# Docker Compose
docker-compose up transaction-service

# Con dependencias
docker-compose up mysql user-service wallet-service transaction-service
```

## 📡 API Endpoints

### Gestión de Transacciones

```http
POST   /api/transactions              # Crear transacción
GET    /api/transactions              # Listar transacciones
GET    /api/transactions/{id}         # Obtener transacción
PUT    /api/transactions/{id}         # Actualizar transacción
DELETE /api/transactions/{id}         # Cancelar transacción
```

### Transferencias

```http
POST   /api/transactions/transfer     # Transferencia entre cuentas
POST   /api/transactions/payment      # Procesar pago
POST   /api/transactions/deposit      # Depósito
POST   /api/transactions/withdrawal   # Retiro
```

### Reportes

```http
GET    /api/transactions/summary      # Resumen de transacciones
GET    /api/transactions/history      # Historial detallado
GET    /api/transactions/balance      # Balance actual
```

### Health Check

```http
GET /health
```

### Ejemplos de Uso

```bash
# Crear transacción
curl -X POST http://localhost:8082/api/transactions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "type": "transfer",
    "amount": 100.50,
    "currency": "USD",
    "fromAccountId": "acc_123",
    "toAccountId": "acc_456",
    "description": "Payment for services"
  }'

# Obtener historial
curl -X GET "http://localhost:8082/api/transactions?limit=10&offset=0" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Health check
curl http://localhost:8082/health
```

## 🔐 Seguridad

### Medidas Implementadas

- **JWT Authentication**: Validación de tokens
- **Transaction Validation**: Validación de fondos y límites
- **Fraud Detection**: Detección básica de fraude
- **Audit Trail**: Registro completo de transacciones
- **Rate Limiting**: Limitación de transacciones por usuario
- **Encryption**: Encriptación de datos sensibles

## 🧪 Testing

```bash
# Tests unitarios
go test ./internal/...

# Tests de integración
go test ./tests/integration/...

# Tests con base de datos de prueba
DB_NAME=fintrack_test go test ./...
```

## 📊 Monitoreo

### Métricas Específicas

- **Transaction Volume**: Volumen de transacciones
- **Success Rate**: Tasa de éxito de transacciones
- **Average Amount**: Monto promedio por transacción
- **Processing Time**: Tiempo de procesamiento
- **Failed Transactions**: Transacciones fallidas

## 🚀 Despliegue

```bash
# Build de producción
CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-w -s" \
  -o transaction-service cmd/main.go
```

---

**Transaction Service** - Procesamiento seguro de transacciones 💳🔒