# User Service - FinTrack

## 📋 Descripción

Microservicio encargado de la gestión de usuarios y autenticación en la plataforma FinTrack. Maneja registro, login, perfiles de usuario y autenticación JWT.

## 🛠️ Tecnologías

- **Lenguaje**: Go 1.24+
- **Framework**: Gin/Echo (HTTP Router)
- **Base de Datos**: MySQL 8.0
- **Autenticación**: JWT (JSON Web Tokens)
- **Contenedor**: Docker multi-stage
- **Arquitectura**: Clean Architecture

## 🏗️ Arquitectura

### Estructura del Proyecto

```
user-service/
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

### Clean Architecture Layers

1. **Domain**: Entidades y reglas de negocio
2. **Service**: Casos de uso y lógica de aplicación
3. **Repository**: Acceso a datos
4. **Handlers**: Controladores HTTP

## 🚀 Desarrollo Local

### Prerrequisitos

- Go 1.24+
- MySQL 8.0+
- Docker (opcional)

### Configuración

```bash
# Clonar y navegar al servicio
cd backend/services/user-service

# Instalar dependencias
go mod download

# Verificar dependencias
go mod verify
```

### Variables de Entorno

```env
# Base de datos
DB_HOST=localhost
DB_PORT=3306
DB_NAME=fintrack
DB_USER=fintrack_user
DB_PASSWORD=fintrack_password

# JWT
JWT_SECRET=your-jwt-secret-key
JWT_EXPIRY=24h

# Servidor
PORT=8080
GIN_MODE=debug

# Logging
LOG_LEVEL=info
```

### Comandos de Desarrollo

```bash
# Ejecutar en modo desarrollo
go run cmd/main.go

# Build del binario
go build -o bin/user-service cmd/main.go

# Ejecutar binario
./bin/user-service

# Tests unitarios
go test ./...

# Tests con coverage
go test -cover ./...

# Tests con reporte detallado
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Linting
golangci-lint run

# Formateo de código
go fmt ./...
```

## 🐳 Docker

### Build Local

```bash
# Build de la imagen
docker build -t fintrack-user-service .

# Ejecutar contenedor
docker run -p 8081:8080 \
  -e DB_HOST=host.docker.internal \
  -e DB_PASSWORD=your_password \
  -e JWT_SECRET=your_secret \
  fintrack-user-service
```

### Docker Compose

```bash
# Desde el directorio raíz del proyecto
docker-compose up user-service

# Con rebuild
docker-compose up --build user-service

# Solo user-service y dependencias
docker-compose up mysql user-service
```

## 📡 API Endpoints

### Autenticación

```http
POST /api/auth/register
POST /api/auth/login
POST /api/auth/refresh
POST /api/auth/logout
```

### Gestión de Usuarios

```http
GET    /api/users/profile
PUT    /api/users/profile
DELETE /api/users/profile
PUT    /api/users/password
GET    /api/users/preferences
PUT    /api/users/preferences
```

### Health Check

```http
GET /health
```

### Ejemplos de Uso

```bash
# Registro de usuario
curl -X POST http://localhost:8081/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword",
    "firstName": "John",
    "lastName": "Doe"
  }'

# Login
curl -X POST http://localhost:8081/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword"
  }'

# Obtener perfil (requiere JWT token)
curl -X GET http://localhost:8081/api/users/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Health check
curl http://localhost:8081/health
```

## 🔐 Seguridad

### Medidas Implementadas

- **JWT Authentication**: Tokens seguros con expiración
- **Password Hashing**: bcrypt para hash de contraseñas
- **Rate Limiting**: Limitación de requests por IP
- **Input Validation**: Validación estricta de inputs
- **SQL Injection Protection**: Prepared statements
- **CORS**: Configuración de CORS apropiada

### Middleware de Seguridad

```go
// Middleware de autenticación JWT
func AuthMiddleware() gin.HandlerFunc

// Middleware de rate limiting
func RateLimitMiddleware() gin.HandlerFunc

// Middleware de validación
func ValidationMiddleware() gin.HandlerFunc
```

## 🧪 Testing

### Estructura de Tests

```
tests/
├── unit/                    # Tests unitarios
│   ├── handlers/
│   ├── service/
│   └── repository/
├── integration/             # Tests de integración
└── mocks/                   # Mocks para testing
```

### Ejecutar Tests

```bash
# Tests unitarios
go test ./internal/...

# Tests de integración
go test ./tests/integration/...

# Tests con base de datos de prueba
DB_NAME=fintrack_test go test ./...

# Benchmark tests
go test -bench=. ./...
```

## 📊 Monitoreo

### Métricas Disponibles

- **Health Check**: `/health`
- **Metrics**: `/metrics` (Prometheus format)
- **Request Duration**: Tiempo de respuesta por endpoint
- **Error Rate**: Tasa de errores por endpoint
- **Active Connections**: Conexiones activas

### Logs Estructurados

```json
{
  "level": "info",
  "timestamp": "2024-01-15T10:30:00Z",
  "service": "user-service",
  "method": "POST",
  "path": "/api/auth/login",
  "status": 200,
  "duration": "45ms",
  "user_id": "12345"
}
```

## 🔧 Configuración Avanzada

### Database Connection Pool

```go
// Configuración de pool de conexiones
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)
```

### JWT Configuration

```go
// Configuración JWT
type JWTConfig struct {
    Secret     string
    Expiry     time.Duration
    RefreshExp time.Duration
    Issuer     string
}
```

## 🚀 Despliegue

### Build de Producción

```bash
# Build optimizado
CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-w -s" \
  -o user-service cmd/main.go

# Verificar binario
./user-service --version
```

### Variables de Producción

```env
GIN_MODE=release
LOG_LEVEL=warn
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
JWT_EXPIRY=1h
RATE_LIMIT=100
```

## 🔍 Troubleshooting

### Problemas Comunes

```bash
# Verificar conectividad a la base de datos
telnet mysql 3306

# Logs del contenedor
docker-compose logs -f user-service

# Verificar health check
curl http://localhost:8081/health

# Debug de JWT tokens
echo "YOUR_JWT_TOKEN" | base64 -d
```

### Debug Mode

```bash
# Ejecutar con debug
GIN_MODE=debug LOG_LEVEL=debug go run cmd/main.go

# Profiling
go tool pprof http://localhost:8081/debug/pprof/profile
```

## 📚 Dependencias Principales

```go
// go.mod principales
require (
    github.com/gin-gonic/gin v1.9.1
    github.com/golang-jwt/jwt/v5 v5.0.0
    github.com/go-sql-driver/mysql v1.7.1
    golang.org/x/crypto v0.14.0
    github.com/go-playground/validator/v10 v10.15.5
)
```

## 🤝 Contribución

1. Seguir las convenciones de Go
2. Escribir tests para nuevas funcionalidades
3. Mantener cobertura de tests > 80%
4. Documentar funciones públicas
5. Usar linting antes de commit

---

**User Service** - Gestión segura de usuarios 👤🔐