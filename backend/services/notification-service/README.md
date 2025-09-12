# Notification Service - FinTrack

## 📋 Descripción

Microservicio encargado del sistema de notificaciones en la plataforma FinTrack. Maneja envío de emails, notificaciones push, SMS y alertas del sistema.

## 🛠️ Tecnologías

- **Lenguaje**: Go 1.24+
- **Framework**: Gin/Echo (HTTP Router)
- **Base de Datos**: MySQL 8.0
- **Email**: SMTP (Gmail, SendGrid, etc.)
- **Push Notifications**: Firebase Cloud Messaging
- **Contenedor**: Docker multi-stage
- **Arquitectura**: Clean Architecture

## 🏗️ Arquitectura

### Estructura del Proyecto

```
notification-service/
├── cmd/
│   └── main.go              # Punto de entrada
├── internal/
│   ├── config/              # Configuración
│   ├── domain/              # Entidades de dominio
│   ├── handlers/            # HTTP handlers
│   ├── repository/          # Capa de datos
│   ├── service/             # Lógica de negocio
│   ├── providers/           # Proveedores de notificaciones
│   └── middleware/          # Middlewares HTTP
├── templates/               # Plantillas de email
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

# Email SMTP
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=noreply@fintrack.com

# Firebase (Push Notifications)
FIREBASE_PROJECT_ID=your-project-id
FIREBASE_PRIVATE_KEY=your-private-key
FIREBASE_CLIENT_EMAIL=your-client-email

# SMS (Twilio ejemplo)
TWILIO_ACCOUNT_SID=your-account-sid
TWILIO_AUTH_TOKEN=your-auth-token
TWILIO_PHONE_NUMBER=+1234567890

# Servidor
PORT=8080
GIN_MODE=debug

# Logging
LOG_LEVEL=info
```

### Comandos de Desarrollo

```bash
# Navegar al servicio
cd backend/services/notification-service

# Instalar dependencias
go mod download

# Ejecutar en modo desarrollo
go run cmd/main.go

# Build del binario
go build -o bin/notification-service cmd/main.go

# Tests
go test ./...
```

## 🐳 Docker

```bash
# Build de la imagen
docker build -t fintrack-notification-service .

# Docker Compose
docker-compose up notification-service

# Con dependencias
docker-compose up mysql user-service notification-service
```

## 📡 API Endpoints

### Envío de Notificaciones

```http
POST   /api/notifications/email       # Enviar email
POST   /api/notifications/push        # Enviar push notification
POST   /api/notifications/sms         # Enviar SMS
POST   /api/notifications/bulk        # Envío masivo
```

### Gestión de Notificaciones

```http
GET    /api/notifications             # Listar notificaciones
GET    /api/notifications/{id}        # Obtener notificación
PUT    /api/notifications/{id}/read   # Marcar como leída
DELETE /api/notifications/{id}        # Eliminar notificación
```

### Preferencias

```http
GET    /api/notifications/preferences # Obtener preferencias
PUT    /api/notifications/preferences # Actualizar preferencias
POST   /api/notifications/subscribe   # Suscribirse a notificaciones
POST   /api/notifications/unsubscribe # Desuscribirse
```

### Plantillas

```http
GET    /api/notifications/templates   # Listar plantillas
GET    /api/notifications/templates/{id} # Obtener plantilla
POST   /api/notifications/templates   # Crear plantilla
PUT    /api/notifications/templates/{id} # Actualizar plantilla
```

### Health Check

```http
GET /health
```

### Ejemplos de Uso

```bash
# Enviar email
curl -X POST http://localhost:8085/api/notifications/email \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "to": "user@example.com",
    "subject": "Transacción Completada",
    "template": "transaction_success",
    "data": {
      "amount": "$100.00",
      "transactionId": "txn_123"
    }
  }'

# Enviar push notification
curl -X POST http://localhost:8085/api/notifications/push \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "userId": "user_123",
    "title": "Nueva Transacción",
    "body": "Has recibido $50.00",
    "data": {
      "type": "transaction",
      "transactionId": "txn_456"
    }
  }'

# Obtener preferencias
curl -X GET http://localhost:8085/api/notifications/preferences \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Health check
curl http://localhost:8085/health
```

## 📧 Tipos de Notificaciones

### Email Templates

```html
<!-- templates/transaction_success.html -->
<!DOCTYPE html>
<html>
<head>
    <title>Transacción Exitosa</title>
</head>
<body>
    <h1>¡Transacción Completada!</h1>
    <p>Tu transacción por {{.Amount}} ha sido procesada exitosamente.</p>
    <p>ID de Transacción: {{.TransactionId}}</p>
</body>
</html>
```

### Push Notification Types

```go
type NotificationType string

const (
    TransactionSuccess NotificationType = "transaction_success"
    TransactionFailed  NotificationType = "transaction_failed"
    LowBalance        NotificationType = "low_balance"
    SecurityAlert     NotificationType = "security_alert"
    AccountUpdate     NotificationType = "account_update"
)
```

## 🔐 Seguridad

### Medidas Implementadas

- **JWT Authentication**: Validación de tokens
- **Rate Limiting**: Limitación de envíos por usuario
- **Email Validation**: Validación de direcciones de email
- **Template Sanitization**: Sanitización de plantillas
- **Spam Protection**: Protección contra spam
- **Audit Trail**: Registro de todas las notificaciones

## 🧪 Testing

```bash
# Tests unitarios
go test ./internal/...

# Tests de integración
go test ./tests/integration/...

# Tests de plantillas
go test ./internal/templates/...

# Tests con proveedores mock
go test ./internal/providers/...
```

## 📊 Monitoreo

### Métricas Específicas

- **Emails Sent**: Emails enviados
- **Push Notifications Sent**: Push notifications enviadas
- **SMS Sent**: SMS enviados
- **Delivery Rate**: Tasa de entrega
- **Open Rate**: Tasa de apertura (emails)
- **Click Rate**: Tasa de clicks
- **Failed Deliveries**: Entregas fallidas

### Logs Estructurados

```json
{
  "level": "info",
  "timestamp": "2024-01-15T10:30:00Z",
  "service": "notification-service",
  "type": "email",
  "recipient": "user@example.com",
  "template": "transaction_success",
  "status": "sent",
  "provider": "smtp"
}
```

## 🚀 Despliegue

### Variables de Producción

```env
GIN_MODE=release
LOG_LEVEL=warn
SMTP_POOL_SIZE=10
PUSH_BATCH_SIZE=100
RATE_LIMIT_PER_HOUR=1000
```

```bash
# Build de producción
CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-w -s" \
  -o notification-service cmd/main.go
```

---

**Notification Service** - Sistema completo de notificaciones 📧📱