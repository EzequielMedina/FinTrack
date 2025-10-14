# 📧 FinTrack Notification Service

Microservicio para el envío de notificaciones automáticas de vencimiento de tarjetas de crédito.

## 🎯 Funcionalidades

- **Job Automático**: Ejecuta diariamente a las 8:00 AM (configurable)
- **Detección de Vencimientos**: Identifica tarjetas que vencen mañana
- **Cálculo de Cuotas**: Suma todas las cuotas pendientes por tarjeta
- **Envío de Emails**: Utiliza EmailJS para enviar notificaciones personalizadas
- **Auditoría Completa**: Registra logs de todas las ejecuciones y notificaciones

## 🏗️ Arquitectura

El servicio sigue Clean Architecture con las siguientes capas:

```
notification-service/
├── cmd/                    # Entry point
├── internal/
│   ├── app/               # Application setup
│   ├── config/            # Configuration management
│   ├── core/
│   │   ├── domain/        # Business entities
│   │   ├── ports/         # Interfaces
│   │   └── service/       # Business logic
│   └── infrastructure/
│       ├── adapters/      # External adapters
│       ├── entrypoints/   # HTTP handlers
│       └── jobs/          # Job scheduler
└── templates/             # Email templates
```

## 🚀 Inicio Rápido

### Usando Docker Compose

```bash
# Construir y ejecutar el servicio
docker-compose up --build notification-service

# Solo el servicio de notificaciones con MySQL
docker-compose up --build mysql notification-service
```

### Desarrollo Local

```bash
# 1. Clonar y navegar al directorio
cd backend/services/notification-service

# 2. Instalar dependencias
go mod download

# 3. Configurar variables de entorno
cp .env.example .env
# Editar .env con tus credenciales

# 4. Ejecutar
go run cmd/main.go
```

## ⚙️ Configuración

### Variables de Entorno

| Variable | Descripción | Valor por Defecto |
|----------|-------------|-------------------|
| `PORT` | Puerto del servidor | `8088` |
| `DB_HOST` | Host de MySQL | `mysql` |
| `DB_NAME` | Nombre de la base de datos | `fintrack` |
| `EMAILJS_SERVICE_ID` | Service ID de EmailJS | Ver keus.txt |
| `EMAILJS_TEMPLATE_ID` | Template ID de EmailJS | Ver keus.txt |
| `JOB_ENABLED` | Habilitar job automático | `true` |
| `JOB_SCHEDULE` | Cron schedule | `0 8 * * *` (8:00 AM) |

### EmailJS Setup

El servicio utiliza EmailJS para el envío de emails. Las credenciales están en `keus.txt`:

- **Service ID**: `service_ceg7xlp`
- **Template ID**: `template_e43va39`
- **Public Key**: `MSBb87-PQcXWr1gWK`
- **Private Key**: `sXLmpEZ8y2EYtCDtN5gZv`

## 📋 API Endpoints

### Job Management

```bash
# Trigger manual del job
POST /api/notifications/trigger-card-due-job

# Historial de ejecuciones
GET /api/notifications/job-history?limit=20

# Estado del scheduler
GET /api/notifications/scheduler/status
```

### Logs y Auditoría

```bash
# Logs de notificaciones para un job
GET /api/notifications/logs?job_run_id=123&limit=50

# Health check
GET /api/notifications/health
GET /health
```

### Ejemplos de Respuesta

```json
// GET /api/notifications/job-history
{
  "data": [
    {
      "run_id": "uuid-123",
      "started_at": "2024-01-15T08:00:00Z",
      "completed_at": "2024-01-15T08:02:30Z",
      "status": "completed",
      "cards_found": 3,
      "emails_sent": 2,
      "errors": 1,
      "duration": "2m30s"
    }
  ],
  "count": 1,
  "limit": 20
}
```

## 🗄️ Base de Datos

### Tablas Creadas

1. **`job_runs`**: Historial de ejecuciones del job
2. **`notification_logs`**: Logs detallados de cada notificación

### Consultas Principales

```sql
-- Tarjetas que vencen mañana
SELECT c.*, u.email, u.first_name, u.last_name
FROM cards c
JOIN accounts a ON c.account_id = a.id  
JOIN users u ON a.user_id = u.id
WHERE DATE(c.due_date) = DATE(NOW() + INTERVAL 1 DAY)
  AND c.status = 'active' AND c.card_type = 'credit';

-- Cuotas pendientes por tarjeta
SELECT i.*, ip.description, ip.merchant_name
FROM installments i
JOIN installment_plans ip ON i.plan_id = ip.id
WHERE ip.card_id = ? AND i.status IN ('pending', 'overdue')
  AND i.due_date <= ?;
```

## 🕐 Job Scheduler

### Configuración del Cron

- **Schedule**: `0 8 * * *` (8:00 AM todos los días)
- **Timezone**: `America/Argentina/Buenos_Aires`
- **Execution**: Asíncrono con logging completo

### Flujo del Job

1. **Inicio**: Crear registro en `job_runs`
2. **Búsqueda**: Obtener tarjetas que vencen mañana
3. **Procesamiento**: Para cada tarjeta:
   - Calcular cuotas pendientes
   - Generar email personalizado
   - Enviar via EmailJS
   - Registrar resultado
4. **Finalización**: Actualizar estadísticas del job

## 📧 Templates de Email

El servicio genera emails HTML personalizados con:

- **Header**: Branding de FinTrack
- **Información de la tarjeta**: Nombre, banco, últimos 4 dígitos
- **Fecha de vencimiento**: Formateada y destacada
- **Total a pagar**: Suma de todas las cuotas
- **Detalle de cuotas**: Lista con descripción, merchant y montos
- **Footer**: Información de contacto

### Ejemplo de Email

```
🏦 FinTrack - Recordatorio de Pago

Hola Juan Pérez, tu tarjeta vence mañana 📅

Visa - Banco Santander (****1234)
Fecha de vencimiento: 15/01/2024

💰 Total a pagar: $25,750.00
Tienes 3 cuotas pendientes que vencen mañana o antes:

📋 Detalle de cuotas:
• Compra en Mercado Libre - Cuota 2/6: $8,500.00
• Combustible YPF - Pago único: $12,250.00  
• Supermercado Coto - Cuota 1/3: $5,000.00
```

## 🐳 Docker

### Dockerfile Multi-stage

```dockerfile
FROM golang:1.24-alpine AS builder
# ... build stage

FROM alpine:latest  
# ... production stage
EXPOSE 8088
CMD ["./notification-service"]
```

### Docker Compose Integration

```yaml
notification-service:
  build: ./backend/services/notification-service
  ports:
    - "8088:8088"
  environment:
    - DB_HOST=mysql
    - JOB_ENABLED=true
  depends_on:
    - mysql
```

## 🧪 Testing

### Trigger Manual

```bash
# Ejecutar job manualmente
curl -X POST http://localhost:8088/api/notifications/trigger-card-due-job

# Verificar estado
curl http://localhost:8088/api/notifications/scheduler/status
```

### Health Check

```bash
curl http://localhost:8088/health
```

## � Monitoreo

### Métricas Importantes

- **Cards Found**: Tarjetas encontradas por ejecución
- **Emails Sent**: Emails enviados exitosamente  
- **Error Rate**: Porcentaje de fallos en envío
- **Execution Time**: Duración de cada job

### Logs Estructurados

```
2024-01-15T08:00:00Z [INFO] 🚀 Starting card due notifications job: job_123
2024-01-15T08:00:15Z [INFO] 📅 Found 3 cards due tomorrow  
2024-01-15T08:01:30Z [INFO] ✅ Notification sent for card Visa (card_456) to user@email.com
2024-01-15T08:02:30Z [INFO] 🎉 Job completed: 2 emails sent, 1 errors
```

## � Desarrollo

### Estructura del Código

- **Entities**: Modelos de dominio (Card, Notification, JobRun)
- **Repositories**: Acceso a datos con interfaces
- **Services**: Lógica de negocio
- **Adapters**: Integraciones externas (EmailJS, MySQL)
- **Jobs**: Scheduler y ejecución de tareas

### Agregar Nuevas Funcionalidades

1. **Nuevos tipos de notificación**: Extender `NotificationService`
2. **Nuevos providers de email**: Implementar `EmailService` interface
3. **Nuevas reglas de negocio**: Modificar `buildCardDueNotification`

## � Troubleshooting

### Problemas Comunes

```bash
# 1. Error de conexión a MySQL
Error: "connection refused"
Solución: Verificar que MySQL esté ejecutándose

# 2. EmailJS API error
Error: "401 Unauthorized"  
Solución: Verificar credenciales en .env

# 3. No se encuentran tarjetas
Info: "Found 0 cards due tomorrow"
Solución: Verificar datos de prueba en base de datos
```

### Debug Mode

```bash
# Activar logs detallados
export GIN_MODE=debug
export LOG_LEVEL=debug
```

---

**🎉 FinTrack Notification Service - Notificaciones automáticas de vencimiento** 📧�