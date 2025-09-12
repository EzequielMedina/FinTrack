# Wallet Service - FinTrack

## 📋 Descripción

Microservicio encargado de la gestión de billeteras digitales en la plataforma FinTrack. Maneja creación, actualización, balance y operaciones de billeteras virtuales.

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
wallet-service/
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

# Servidor
PORT=8080
GIN_MODE=debug

# Logging
LOG_LEVEL=info
```

### Comandos de Desarrollo

```bash
# Navegar al servicio
cd backend/services/wallet-service

# Instalar dependencias
go mod download

# Ejecutar en modo desarrollo
go run cmd/main.go

# Build del binario
go build -o bin/wallet-service cmd/main.go

# Tests
go test ./...
```

## 🐳 Docker

```bash
# Build de la imagen
docker build -t fintrack-wallet-service .

# Docker Compose
docker-compose up wallet-service

# Con dependencias
docker-compose up mysql user-service wallet-service
```

## 📡 API Endpoints

### Gestión de Billeteras

```http
POST   /api/wallets                   # Crear billetera
GET    /api/wallets                   # Listar billeteras del usuario
GET    /api/wallets/{id}              # Obtener billetera específica
PUT    /api/wallets/{id}              # Actualizar billetera
DELETE /api/wallets/{id}              # Eliminar billetera
```

### Balance y Operaciones

```http
GET    /api/wallets/{id}/balance      # Obtener balance
POST   /api/wallets/{id}/deposit      # Depositar fondos
POST   /api/wallets/{id}/withdraw     # Retirar fondos
GET    /api/wallets/{id}/history      # Historial de movimientos
```

### Configuración

```http
GET    /api/wallets/{id}/settings     # Configuración de billetera
PUT    /api/wallets/{id}/settings     # Actualizar configuración
POST   /api/wallets/{id}/freeze       # Congelar billetera
POST   /api/wallets/{id}/unfreeze     # Descongelar billetera
```

### Health Check

```http
GET /health
```

### Ejemplos de Uso

```bash
# Crear billetera
curl -X POST http://localhost:8083/api/wallets \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Mi Billetera Principal",
    "currency": "USD",
    "type": "personal"
  }'

# Obtener balance
curl -X GET http://localhost:8083/api/wallets/wallet_123/balance \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Depositar fondos
curl -X POST http://localhost:8083/api/wallets/wallet_123/deposit \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "amount": 500.00,
    "description": "Depósito inicial"
  }'

# Health check
curl http://localhost:8083/health
```

## 🔐 Seguridad

### Medidas Implementadas

- **JWT Authentication**: Validación de tokens
- **Wallet Ownership**: Verificación de propiedad
- **Balance Validation**: Validación de fondos suficientes
- **Transaction Limits**: Límites de transacciones
- **Audit Trail**: Registro de todas las operaciones
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

- **Total Wallets**: Número total de billeteras
- **Active Wallets**: Billeteras activas
- **Total Balance**: Balance total del sistema
- **Average Balance**: Balance promedio por billetera
- **Wallet Operations**: Operaciones por billetera

## 🚀 Despliegue

```bash
# Build de producción
CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-w -s" \
  -o wallet-service cmd/main.go
```

---

**Wallet Service** - Gestión segura de billeteras digitales 👛💰