# Account Service - FinTrack

## 📋 Descripción

Microservicio encargado de la gestión de cuentas virtuales y financieras en la plataforma FinTrack. Implementa operaciones CRUD para cuentas, maneja diferentes tipos de cuenta (corriente, ahorro, tarjetas, etc.) y soporta múltiples monedas (ARS, USD).

## ✅ TASK-014 Implementación Completada

Este microservicio implementa TASK-014 "Implementar microservicio de cuentas virtuales" con las siguientes características:

### Funcionalidades Implementadas

#### Tipos de Cuenta Soportados
- **Cuenta Corriente** (`checking`)
- **Caja de Ahorro** (`savings`) 
- **Cuenta USD** (`usd`)
- **Tarjeta de Crédito** (`credit_card`)
- **Tarjeta de Débito** (`debit_card`)
- **Billetera Virtual** (`wallet`)

#### Monedas Soportadas
- **ARS** (Peso Argentino)
- **USD** (Dólar Estadounidense)

#### Estados de Cuenta
- **Activa** (`active`)
- **Inactiva** (`inactive`)
- **Congelada** (`frozen`)
- **Cerrada** (`closed`)

### Arquitectura Clean Architecture

La implementación sigue los principios SOLID y Clean Architecture:

- **Domain Layer** (`internal/domain/`): Entidades de negocio y reglas de dominio
- **Application Layer** (`internal/application/services/`): Casos de uso y lógica de negocio
- **Infrastructure Layer** (`internal/infrastructure/`): HTTP handlers, middleware, persistencia

### Endpoints API REST

```http
# Gestión de Cuentas
POST   /api/v1/users/{user_id}/accounts     # Crear cuenta
GET    /api/v1/users/{user_id}/accounts     # Listar cuentas de usuario
GET    /api/v1/accounts/{account_id}        # Obtener cuenta
PUT    /api/v1/accounts/{account_id}        # Actualizar cuenta
DELETE /api/v1/accounts/{account_id}        # Eliminar cuenta

# Balance y Operaciones
GET    /api/v1/accounts/{account_id}/balance # Consultar balance

# Configuración
GET    /api/v1/accounts/types              # Tipos y monedas disponibles
```

### Seguridad y Autorización

- **JWT Authentication**: Autenticación requerida para todas las operaciones
- **Role-based Authorization**: Roles jerárquicos (admin > treasurer > operator > user)
- **Request ID Tracking**: Trazabilidad de requests para auditoría
- **CORS Support**: Configuración CORS para frontend
- **Input Validation**: Validación de entrada en todos los endpoints

### Validaciones de Negocio

- Límites de cuentas por usuario
- Validación de tipos de cuenta y monedas
- Verificación de permisos de usuario
- Validación de nombres únicos de cuenta por usuario

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
account-service/
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

# APIs bancarias (ejemplo)
BANK_API_URL=https://api.bank.com
BANK_API_KEY=your-bank-api-key

# Servidor
PORT=8080
GIN_MODE=debug

# Logging
LOG_LEVEL=info
```

### Comandos de Desarrollo

```bash
# Navegar al servicio
cd backend/services/account-service

# Instalar dependencias
go mod download

# Ejecutar en modo desarrollo
go run cmd/main.go

# Build del binario
go build -o bin/account-service cmd/main.go

# Tests
go test ./...
```

## 🐳 Docker

```bash
# Build de la imagen
docker build -t fintrack-account-service .

# Docker Compose
docker-compose up account-service

# Con dependencias
docker-compose up mysql user-service account-service
```

## 📡 API Endpoints

### Gestión de Cuentas

```http
POST   /api/accounts                  # Vincular nueva cuenta
GET    /api/accounts                  # Listar cuentas del usuario
GET    /api/accounts/{id}             # Obtener cuenta específica
PUT    /api/accounts/{id}             # Actualizar cuenta
DELETE /api/accounts/{id}             # Desvincular cuenta
```

### Verificación

```http
POST   /api/accounts/{id}/verify      # Verificar cuenta
GET    /api/accounts/{id}/status      # Estado de verificación
POST   /api/accounts/{id}/revalidate  # Re-validar cuenta
```

### Información Bancaria

```http
GET    /api/accounts/{id}/balance     # Balance de cuenta
GET    /api/accounts/{id}/details     # Detalles de cuenta
GET    /api/accounts/{id}/transactions # Transacciones de cuenta
```

### Health Check

```http
GET /health
```

### Ejemplos de Uso

```bash
# Vincular cuenta bancaria
curl -X POST http://localhost:8084/api/accounts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "bankName": "Banco Nacional",
    "accountNumber": "1234567890",
    "accountType": "checking",
    "routingNumber": "021000021",
    "nickname": "Cuenta Principal"
  }'

# Verificar cuenta
curl -X POST http://localhost:8084/api/accounts/acc_123/verify \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "verificationMethod": "microdeposits"
  }'

# Obtener balance
curl -X GET http://localhost:8084/api/accounts/acc_123/balance \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Health check
curl http://localhost:8084/health
```

## 🔐 Seguridad

### Medidas Implementadas

- **JWT Authentication**: Validación de tokens
- **Account Ownership**: Verificación de propiedad
- **Data Encryption**: Encriptación de datos bancarios
- **PCI Compliance**: Cumplimiento de estándares PCI
- **Audit Trail**: Registro de accesos y modificaciones
- **Rate Limiting**: Limitación de requests

### Datos Sensibles

```go
// Encriptación de números de cuenta
type EncryptedAccount struct {
    ID              string
    UserID          string
    BankName        string
    AccountNumber   string `encrypt:"true"`
    RoutingNumber   string `encrypt:"true"`
    AccountType     string
    Nickname        string
    IsVerified      bool
    CreatedAt       time.Time
}
```

## 🧪 Testing

```bash
# Tests unitarios
go test ./internal/...

# Tests de integración
go test ./tests/integration/...

# Tests con base de datos de prueba
DB_NAME=fintrack_test go test ./...

# Tests de encriptación
go test ./internal/crypto/...
```

## 📊 Monitoreo

### Métricas Específicas

- **Linked Accounts**: Cuentas vinculadas
- **Verified Accounts**: Cuentas verificadas
- **Verification Rate**: Tasa de verificación
- **API Response Time**: Tiempo de respuesta de APIs bancarias
- **Failed Verifications**: Verificaciones fallidas

## 🚀 Despliegue

```bash
# Build de producción
CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-w -s" \
  -o account-service cmd/main.go
```

---

**Account Service** - Gestión segura de cuentas bancarias 🏦🔐