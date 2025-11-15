# 📧 PLAN DE IMPLEMENTACIÓN - NOTIFICATION SERVICE
## Job de Notificaciones de Vencimiento de Tarjetas de Crédito

---

## 📋 RESUMEN EJECUTIVO

### **Objetivo Principal**
Implementar un notification-service que ejecute un job diario para:
1. **Identificar** tarjetas de crédito que vencen mañana (`due_date = tomorrow`)
2. **Calcular** la suma de todas las cuotas (`installments`) pendientes por tarjeta
3. **Enviar** notificación por email usando EmailJS con template personalizado
4. **Registrar** el historial de notificaciones enviadas

### **Arquitectura Técnica**
- **Microservicio**: Go + Clean Architecture 
- **Base de Datos**: MySQL (lectura de `cards` e `installments`)
- **Email Provider**: EmailJS API
- **Job Scheduler**: Cron job interno en Go
- **Puerto**: 8088 (siguiendo el patrón de FinTrack)

---

## 🏗️ FASE 1: ESTRUCTURA BASE Y CONFIGURACIÓN

### 1.1 Estructura de Clean Architecture ✅
```
notification-service/
├── cmd/
│   └── main.go                    # Entry point
├── internal/
│   ├── app/
│   │   └── application.go         # App setup
│   ├── config/
│   │   └── config.go             # Configuration
│   ├── core/
│   │   ├── domain/
│   │   │   └── entities/         # Business entities
│   │   ├── ports/
│   │   │   ├── repository.go     # Repository interfaces
│   │   │   └── service.go        # Service interfaces
│   │   └── services/             # Business logic
│   └── infrastructure/
│       ├── adapters/
│       │   ├── database/         # DB repositories
│       │   └── email/            # EmailJS client
│       ├── entrypoints/
│       │   └── http/             # HTTP handlers
│       └── jobs/                 # Cron jobs
├── templates/
│   └── card_due_notification.html  # Email template
├── .env.example
├── Dockerfile
├── go.mod
└── go.sum
```

### 1.2 Variables de Entorno
```env
# Database Configuration
DB_HOST=mysql
DB_PORT=3306
DB_NAME=fintrack
DB_USER=fintrack_user
DB_PASSWORD=fintrack_password

# EmailJS Configuration (from docs/emailjs.txt)
EMAILJS_SERVICE_ID=service_ceg7xlp
EMAILJS_TEMPLATE_ID=template_e43va39
EMAILJS_PUBLIC_KEY=MSBb87-PQcXWr1gWK
EMAILJS_PRIVATE_KEY=MSBb87-PQcXWr1gWK

# Server Configuration
PORT=8088
GIN_MODE=debug
LOG_LEVEL=info

# Job Configuration
JOB_ENABLED=true
JOB_SCHEDULE=0 8 * * *  # 8:00 AM daily
JOB_TIMEZONE=America/Argentina/Buenos_Aires
```

---

## 🎯 FASE 2: ENTIDADES DE DOMINIO Y DTOs

### 2.1 Entidades Core
```go
// Card representa una tarjeta para notificaciones
type Card struct {
    ID           string    `json:"id"`
    UserID       string    `json:"user_id"`
    CardName     string    `json:"card_name"`
    BankName     string    `json:"bank_name"`
    LastFour     string    `json:"last_four_digits"`
    DueDate      time.Time `json:"due_date"`
    UserEmail    string    `json:"user_email"`    // From user join
    UserName     string    `json:"user_name"`     // From user join
}

// CardDueNotification contiene datos para el email
type CardDueNotification struct {
    CardID              string    `json:"card_id"`
    UserID              string    `json:"user_id"`
    UserEmail           string    `json:"user_email"`
    UserName            string    `json:"user_name"`
    CardName            string    `json:"card_name"`
    BankName            string    `json:"bank_name"`
    LastFour            string    `json:"last_four"`
    DueDate             time.Time `json:"due_date"`
    TotalPendingAmount  float64   `json:"total_pending_amount"`
    PendingInstallments int       `json:"pending_installments"`
    InstallmentDetails  []InstallmentSummary `json:"installment_details"`
}

// InstallmentSummary para detalles de cuotas
type InstallmentSummary struct {
    Description string  `json:"description"`
    Amount      float64 `json:"amount"`
    DueDate     time.Time `json:"due_date"`
}

// NotificationLog para auditoría
type NotificationLog struct {
    ID           string    `json:"id"`
    JobRunID     string    `json:"job_run_id"`
    CardID       string    `json:"card_id"`
    UserID       string    `json:"user_id"`
    Email        string    `json:"email"`
    Status       string    `json:"status"` // sent, failed, skipped
    ErrorMessage string    `json:"error_message,omitempty"`
    SentAt       time.Time `json:"sent_at"`
}
```

### 2.2 DTOs para EmailJS
```go
// EmailJSRequest para la API de EmailJS
type EmailJSRequest struct {
    ServiceID  string            `json:"service_id"`
    TemplateID string            `json:"template_id"`
    UserID     string            `json:"user_id"`
    Template   map[string]string `json:"template_params"`
}

// Template params basados en template.html
type EmailTemplateParams struct {
    FromName        string `json:"from_name"`
    Subject         string `json:"subject"`
    ToEmail         string `json:"to_email"`
    ReplyTo         string `json:"reply_to"`
    HTMLContent     string `json:"html_content"`
}
```

---

## 🔗 FASE 3: INTEGRACIÓN CON EMAILJS

### 3.1 Cliente EmailJS
```go
type EmailJSClient struct {
    serviceID  string
    templateID string
    publicKey  string
    privateKey string
    httpClient *http.Client
}

func (c *EmailJSClient) SendEmail(params EmailTemplateParams) error {
    // POST https://api.emailjs.com/api/v1.0/email/send
    // Headers: Content-Type: application/json
    // Body: EmailJSRequest con template_params
}
```

### 3.2 Template HTML Personalizado
Basado en `docs/template.html` pero adaptado para notificaciones de tarjetas:
```html
<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
  <div style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); padding: 20px; text-align: center;">
    <h1 style="color: white; margin: 0;">FinTrack - Recordatorio de Pago</h1>
  </div>
  
  <div style="padding: 20px; background: #f9f9f9;">
    <h2 style="color: #333;">Hola {{user_name}}, tu tarjeta vence mañana 📅</h2>
    <div style="background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
      <h3>{{card_name}} - {{bank_name}} (****{{last_four}})</h3>
      <p><strong>Fecha de vencimiento:</strong> {{due_date_formatted}}</p>
      <p><strong>Total a pagar:</strong> ${{total_amount}}</p>
      <p><strong>Cuotas pendientes:</strong> {{installments_count}}</p>
      
      <div style="margin-top: 20px;">
        <h4>Detalle de cuotas:</h4>
        {{installments_html}}
      </div>
    </div>
  </div>
  
  <div style="padding: 20px; text-align: center; background: #333; color: white;">
    <p style="margin: 0; font-size: 12px;">
      Este email fue enviado desde FinTrack<br>
      Para consultas, responde a: support@fintrack.com
    </p>
  </div>
</div>
```

---

## 🗄️ FASE 4: CAPA DE DATOS Y REPOSITORIOS

### 4.1 Repository Interface
```go
type CardRepository interface {
    GetCardsDueTomorrow() ([]*entities.Card, error)
    GetUserByID(userID string) (*entities.User, error)
}

type InstallmentRepository interface {
    GetPendingInstallmentsByCard(cardID string, maxDueDate time.Time) ([]*entities.Installment, error)
}

type NotificationRepository interface {
    SaveNotificationLog(log *entities.NotificationLog) error
    GetJobRunHistory(limit int) ([]*entities.JobRun, error)
}
```

### 4.2 SQL Queries Críticas
```sql
-- Obtener tarjetas que vencen mañana con info del usuario
SELECT 
    c.id, c.user_id, c.card_name, c.bank_name, 
    c.last_four_digits, c.due_date,
    u.email, u.first_name, u.last_name
FROM cards c
JOIN users u ON c.user_id = u.id
WHERE DATE(c.due_date) = DATE(NOW() + INTERVAL 1 DAY)
  AND c.is_active = 1
  AND c.card_type = 'credit';

-- Obtener cuotas pendientes para una tarjeta hasta su fecha de vencimiento
SELECT 
    i.id, i.amount, i.due_date, ip.description, ip.merchant_name
FROM installments i
JOIN installment_plans ip ON i.plan_id = ip.id
WHERE ip.card_id = ?
  AND i.status IN ('pending', 'overdue')
  AND i.due_date <= ?
ORDER BY i.due_date ASC;
```

---

## ⚙️ FASE 5: LÓGICA DE NEGOCIO Y SERVICIOS

### 5.1 Notification Service
```go
type NotificationService struct {
    cardRepo         ports.CardRepository
    installmentRepo  ports.InstallmentRepository
    notificationRepo ports.NotificationRepository
    emailClient      *email.EmailJSClient
    logger          *slog.Logger
}

func (s *NotificationService) ProcessCardDueNotifications() error {
    // 1. Get cards due tomorrow
    cards, err := s.cardRepo.GetCardsDueTomorrow()
    
    // 2. For each card, calculate installments
    for _, card := range cards {
        installments, err := s.installmentRepo.GetPendingInstallmentsByCard(
            card.ID, card.DueDate)
        
        // 3. Build notification data
        notification := s.buildNotificationData(card, installments)
        
        // 4. Send email
        err = s.sendNotificationEmail(notification)
        
        // 5. Log result
        s.logNotification(notification, err)
    }
}
```

### 5.2 Template Builder
```go
func (s *NotificationService) buildEmailHTML(notification *CardDueNotification) string {
    // Generar HTML para installments
    installmentsHTML := s.buildInstallmentsHTML(notification.InstallmentDetails)
    
    // Reemplazar placeholders en template
    html := strings.ReplaceAll(baseTemplate, "{{user_name}}", notification.UserName)
    html = strings.ReplaceAll(html, "{{card_name}}", notification.CardName)
    html = strings.ReplaceAll(html, "{{total_amount}}", formatCurrency(notification.TotalPendingAmount))
    // ... más reemplazos
    
    return html
}
```

---

## ⏰ FASE 6: JOB SCHEDULER Y CRON

### 6.1 Cron Job Implementation
```go
type JobScheduler struct {
    notificationService *services.NotificationService
    cron                *cron.Cron
    logger             *slog.Logger
}

func (j *JobScheduler) Start() error {
    // Schedule daily at 8:00 AM
    _, err := j.cron.AddFunc("0 8 * * *", func() {
        jobRunID := uuid.New().String()
        j.logger.Info("Starting card due notifications job", "job_run_id", jobRunID)
        
        if err := j.notificationService.ProcessCardDueNotifications(); err != nil {
            j.logger.Error("Job failed", "error", err, "job_run_id", jobRunID)
        } else {
            j.logger.Info("Job completed successfully", "job_run_id", jobRunID)
        }
    })
    
    j.cron.Start()
    return err
}
```

### 6.2 Manual Trigger Endpoint
```go
// POST /api/notifications/trigger-card-due-job
func (h *NotificationHandler) TriggerCardDueJob(c *gin.Context) {
    go func() {
        if err := h.notificationService.ProcessCardDueNotifications(); err != nil {
            h.logger.Error("Manual job trigger failed", "error", err)
        }
    }()
    
    c.JSON(200, gin.H{"message": "Job triggered successfully"})
}
```

---

## 🌐 FASE 7: API ENDPOINTS Y HANDLERS

### 7.1 HTTP Routes
```go
func SetupRoutes(r *gin.Engine, handler *NotificationHandler) {
    api := r.Group("/api/notifications")
    {
        // Job management
        api.POST("/trigger-card-due-job", handler.TriggerCardDueJob)
        api.GET("/job-history", handler.GetJobHistory)
        
        // Notification logs
        api.GET("/logs", handler.GetNotificationLogs)
        api.GET("/logs/:id", handler.GetNotificationLog)
        
        // Health check
        api.GET("/health", handler.HealthCheck)
    }
}
```

### 7.2 Response DTOs
```go
type JobHistoryResponse struct {
    RunID        string    `json:"run_id"`
    StartedAt    time.Time `json:"started_at"`
    CompletedAt  *time.Time `json:"completed_at"`
    Status       string    `json:"status"`
    CardsFound   int       `json:"cards_found"`
    EmailsSent   int       `json:"emails_sent"`
    Errors       int       `json:"errors"`
}

type NotificationLogResponse struct {
    ID           string    `json:"id"`
    CardName     string    `json:"card_name"`
    UserEmail    string    `json:"user_email"`
    Status       string    `json:"status"`
    ErrorMessage string    `json:"error_message,omitempty"`
    SentAt       time.Time `json:"sent_at"`
}
```

---

## 🐳 FASE 8: DOCKERIZACIÓN E INTEGRACIÓN

### 8.1 Dockerfile
```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o notification-service cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/notification-service .
COPY --from=builder /app/templates ./templates
EXPOSE 8088
CMD ["./notification-service"]
```

### 8.2 Docker Compose Integration
```yaml
# En docker-compose.yml
notification-service:
  build: ./backend/services/notification-service
  ports:
    - "8088:8088"
  environment:
    - DB_HOST=mysql
    - DB_NAME=fintrack
    - EMAILJS_SERVICE_ID=service_ceg7xlp
    - EMAILJS_TEMPLATE_ID=template_e43va39
    - JOB_ENABLED=true
  depends_on:
    - mysql
    - user-service
    - account-service
  networks:
    - fintrack-network
```

---

## 🧪 FASE 9: TESTING

### 9.1 Unit Tests
```go
func TestNotificationService_ProcessCardDueNotifications(t *testing.T) {
    // Mock repositories
    mockCardRepo := &mocks.CardRepository{}
    mockInstallmentRepo := &mocks.InstallmentRepository{}
    mockEmailClient := &mocks.EmailJSClient{}
    
    // Test data
    tomorrow := time.Now().AddDate(0, 0, 1)
    cards := []*entities.Card{
        {ID: "card1", DueDate: tomorrow, UserEmail: "test@example.com"},
    }
    
    mockCardRepo.On("GetCardsDueTomorrow").Return(cards, nil)
    mockInstallmentRepo.On("GetPendingInstallmentsByCard", "card1", tomorrow).
        Return([]*entities.Installment{}, nil)
    mockEmailClient.On("SendEmail", mock.Anything).Return(nil)
    
    // Execute and assert
    service := services.NewNotificationService(mockCardRepo, mockInstallmentRepo, mockEmailClient)
    err := service.ProcessCardDueNotifications()
    assert.NoError(t, err)
}
```

### 9.2 Integration Tests
- EmailJS API integration test
- Database query tests
- End-to-end job execution test

---

## 📊 FASE 10: MONITOREO Y OBSERVABILIDAD

### 10.1 Métricas Específicas
- Cards processed per job run
- Emails sent successfully
- Email failures by reason
- Job execution time
- Database query performance

### 10.2 Logs Estructurados
```json
{
  "level": "info",
  "timestamp": "2024-01-15T08:00:00Z",
  "service": "notification-service",
  "job_run_id": "job_123",
  "cards_found": 5,
  "emails_sent": 4,
  "errors": 1,
  "duration_ms": 2340
}
```

---

## 🚀 PLAN DE EJECUCIÓN

### **Cronograma de Desarrollo (5 días)**

| Día | Fase | Tareas |
|-----|------|--------|
| **Día 1** | Setup + Base | Estructura Clean Architecture, configuración, entidades |
| **Día 2** | Data Layer | Repositorios, queries SQL, cliente EmailJS |
| **Día 3** | Business Logic | Notification service, template builder, job scheduler |
| **Día 4** | API + Docker | HTTP handlers, Dockerfile, docker-compose |
| **Día 5** | Testing + Polish | Unit tests, integration tests, documentación |

### **Orden de Implementación**
1. ✅ Estructura base y configuración
2. ✅ Entidades y DTOs  
3. ✅ Cliente EmailJS
4. ✅ Repositorios de base de datos
5. ✅ Servicio de notificaciones
6. ✅ Job scheduler/cron
7. ✅ HTTP handlers y rutas
8. ✅ Docker y variables de entorno
9. ✅ Integración con docker-compose
10. ✅ Testing completo

---

## 🎯 CRITERIOS DE ACEPTACIÓN

### ✅ **Funcionalidad Core**
- [ ] Job diario se ejecuta automáticamente a las 8:00 AM
- [ ] Identifica correctamente tarjetas que vencen mañana
- [ ] Calcula suma exacta de cuotas pendientes por tarjeta
- [ ] Envía emails usando EmailJS con template personalizado
- [ ] Registra logs completos de cada ejecución

### ✅ **Integración**
- [ ] Se conecta correctamente a base de datos de FinTrack
- [ ] Usa credenciales de EmailJS del archivo docs/emailjs.txt
- [ ] Se integra con docker-compose del proyecto
- [ ] Sigue Clean Architecture consistente con otros servicios

### ✅ **Calidad**
- [ ] Tests unitarios con >80% cobertura
- [ ] Manejo robusto de errores
- [ ] Logs estructurados para debugging
- [ ] Documentación API completa

---

**🎉 NOTIFICATION SERVICE - Notificaciones automáticas de vencimiento de tarjetas** 📧💳