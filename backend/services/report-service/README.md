# Report Service - FinTrack

## 📋 Descripción

Microservicio especializado en la generación de reportes financieros, análisis de datos y visualizaciones para la plataforma FinTrack. Proporciona informes detallados sobre transacciones, balances, tendencias y métricas financieras.

## 🛠️ Tecnologías

- **Lenguaje**: Go 1.24+
- **Framework**: Gin/Echo (HTTP Router)
- **Base de Datos**: MySQL 8.0
- **Analytics**: ClickHouse (opcional)
- **PDF Generation**: wkhtmltopdf
- **Charts**: Chart.js, D3.js
- **Templates**: Go Templates
- **Queue**: RabbitMQ/Redis
- **Storage**: MinIO/S3
- **Contenedor**: Docker multi-stage

## 🏗️ Arquitectura

### Estructura del Proyecto

```
report-service/
├── cmd/
│   └── main.go              # Punto de entrada
├── internal/
│   ├── config/              # Configuración
│   ├── domain/              # Entidades de dominio
│   ├── handlers/            # HTTP handlers
│   ├── repository/          # Capa de datos
│   ├── service/             # Lógica de negocio
│   ├── generator/           # Generadores de reportes
│   ├── templates/           # Plantillas de reportes
│   ├── analytics/           # Motor de análisis
│   ├── scheduler/           # Programador de reportes
│   └── middleware/          # Middlewares HTTP
├── assets/
│   ├── templates/           # Plantillas HTML/PDF
│   ├── styles/              # Estilos CSS
│   └── images/              # Imágenes y logos
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

# Analytics Database (ClickHouse)
CLICKHOUSE_HOST=localhost
CLICKHOUSE_PORT=9000
CLICKHOUSE_DATABASE=fintrack_analytics
CLICKHOUSE_USER=default
CLICKHOUSE_PASSWORD=

# File Storage
STORAGE_TYPE=local # local, s3, minio
STORAGE_PATH=/app/reports
S3_BUCKET=fintrack-reports
S3_REGION=us-east-1
S3_ACCESS_KEY=your-access-key
S3_SECRET_KEY=your-secret-key

# PDF Generation
WKHTMLTOPDF_PATH=/usr/local/bin/wkhtmltopdf
PDF_TIMEOUT=30s
PDF_DPI=300
PDF_PAGE_SIZE=A4

# Queue System
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
REDIS_URL=redis://localhost:6379/1
QUEUE_TYPE=rabbitmq # rabbitmq, redis

# Email (para envío de reportes)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=reports@fintrack.com
SMTP_PASSWORD=your-email-password

# Servidor
PORT=8080
GIN_MODE=debug

# Logging
LOG_LEVEL=info
```

### Comandos de Desarrollo

```bash
# Navegar al servicio
cd backend/services/report-service

# Instalar dependencias
go mod download

# Ejecutar en modo desarrollo
go run cmd/main.go

# Build del binario
go build -o bin/report-service cmd/main.go

# Tests
go test ./...
```

## 🐳 Docker

```bash
# Build de la imagen
docker build -t fintrack-report-service .

# Docker Compose
docker-compose up report-service

# Con dependencias
docker-compose up mysql clickhouse report-service
```

## 📡 API Endpoints

### Generación de Reportes

```http
POST   /api/reports/generate           # Generar reporte
GET    /api/reports/{id}               # Obtener reporte
GET    /api/reports/{id}/download      # Descargar reporte
DELETE /api/reports/{id}               # Eliminar reporte
```

### Reportes Predefinidos

```http
GET    /api/reports/balance            # Reporte de balance
GET    /api/reports/transactions       # Reporte de transacciones
GET    /api/reports/income-expense     # Ingresos vs gastos
GET    /api/reports/cash-flow          # Flujo de efectivo
GET    /api/reports/portfolio          # Reporte de portafolio
```

### Análisis y Métricas

```http
GET    /api/analytics/summary          # Resumen financiero
GET    /api/analytics/trends           # Análisis de tendencias
GET    /api/analytics/categories       # Análisis por categorías
GET    /api/analytics/predictions      # Predicciones financieras
```

### Programación de Reportes

```http
POST   /api/schedules                  # Crear programación
GET    /api/schedules                  # Listar programaciones
PUT    /api/schedules/{id}             # Actualizar programación
DELETE /api/schedules/{id}             # Eliminar programación
```

### Plantillas

```http
GET    /api/templates                  # Listar plantillas
GET    /api/templates/{id}             # Obtener plantilla
POST   /api/templates                  # Crear plantilla
PUT    /api/templates/{id}             # Actualizar plantilla
```

### Health Check

```http
GET /health
```

### Ejemplos de Uso

```bash
# Generar reporte de balance
curl -X POST http://localhost:8088/api/reports/generate \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "type": "balance",
    "format": "pdf",
    "period": {
      "start": "2024-01-01",
      "end": "2024-01-31"
    },
    "filters": {
      "accounts": ["checking", "savings"],
      "currency": "USD"
    }
  }'

# Obtener reporte generado
curl -X GET http://localhost:8088/api/reports/12345 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Descargar reporte PDF
curl -X GET http://localhost:8088/api/reports/12345/download \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -o balance_report.pdf

# Análisis de tendencias
curl -X GET "http://localhost:8088/api/analytics/trends?period=6m&metric=balance" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Programar reporte mensual
curl -X POST http://localhost:8088/api/schedules \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Monthly Balance Report",
    "reportType": "balance",
    "schedule": "0 0 1 * *",
    "format": "pdf",
    "recipients": ["user@example.com"]
  }'

# Health check
curl http://localhost:8088/health
```

## 📊 Tipos de Reportes

### Reportes Financieros

```go
type ReportType string

const (
    BalanceReport       ReportType = "balance"
    TransactionReport   ReportType = "transactions"
    IncomeExpenseReport ReportType = "income_expense"
    CashFlowReport      ReportType = "cash_flow"
    PortfolioReport     ReportType = "portfolio"
    TaxReport          ReportType = "tax"
    BudgetReport       ReportType = "budget"
    NetWorthReport     ReportType = "net_worth"
)
```

### Formatos de Salida

```go
type OutputFormat string

const (
    PDFFormat   OutputFormat = "pdf"
    ExcelFormat OutputFormat = "xlsx"
    CSVFormat   OutputFormat = "csv"
    JSONFormat  OutputFormat = "json"
    HTMLFormat  OutputFormat = "html"
)
```

### Configuración de Reporte

```go
type ReportConfig struct {
    Type        ReportType    `json:"type"`
    Format      OutputFormat  `json:"format"`
    Period      DateRange     `json:"period"`
    Filters     ReportFilters `json:"filters"`
    Template    string        `json:"template,omitempty"`
    Options     ReportOptions `json:"options"`
}

type ReportFilters struct {
    Accounts   []string `json:"accounts,omitempty"`
    Categories []string `json:"categories,omitempty"`
    Currency   string   `json:"currency,omitempty"`
    MinAmount  float64  `json:"min_amount,omitempty"`
    MaxAmount  float64  `json:"max_amount,omitempty"`
}
```

## 📈 Motor de Análisis

### Métricas Calculadas

```go
type FinancialMetrics struct {
    TotalIncome     float64 `json:"total_income"`
    TotalExpenses   float64 `json:"total_expenses"`
    NetIncome       float64 `json:"net_income"`
    SavingsRate     float64 `json:"savings_rate"`
    ExpenseRatio    float64 `json:"expense_ratio"`
    CashFlow        float64 `json:"cash_flow"`
    NetWorth        float64 `json:"net_worth"`
    DebtToIncome    float64 `json:"debt_to_income"`
}
```

### Análisis de Tendencias

```go
type TrendAnalysis struct {
    Metric      string      `json:"metric"`
    Period      string      `json:"period"`
    Trend       TrendType   `json:"trend"`
    Change      float64     `json:"change"`
    ChangeRate  float64     `json:"change_rate"`
    DataPoints  []DataPoint `json:"data_points"`
    Forecast    []DataPoint `json:"forecast,omitempty"`
}

type TrendType string

const (
    UpwardTrend   TrendType = "upward"
    DownwardTrend TrendType = "downward"
    StableTrend   TrendType = "stable"
    VolatileTrend TrendType = "volatile"
)
```

## 🎨 Sistema de Plantillas

### Plantillas HTML

```html
<!-- templates/balance_report.html -->
<!DOCTYPE html>
<html>
<head>
    <title>Balance Report - {{.Period}}</title>
    <link rel="stylesheet" href="/assets/styles/report.css">
</head>
<body>
    <header>
        <img src="/assets/images/logo.png" alt="FinTrack">
        <h1>Balance Report</h1>
        <p>Period: {{.Period}}</p>
    </header>
    
    <section class="summary">
        <div class="metric">
            <h3>Total Assets</h3>
            <span class="amount positive">{{.TotalAssets | currency}}</span>
        </div>
        <div class="metric">
            <h3>Total Liabilities</h3>
            <span class="amount negative">{{.TotalLiabilities | currency}}</span>
        </div>
        <div class="metric">
            <h3>Net Worth</h3>
            <span class="amount {{.NetWorthClass}}">{{.NetWorth | currency}}</span>
        </div>
    </section>
    
    <section class="charts">
        <canvas id="balanceChart"></canvas>
    </section>
    
    <script src="/assets/js/charts.js"></script>
</body>
</html>
```

### Funciones de Template

```go
var templateFuncs = template.FuncMap{
    "currency": func(amount float64) string {
        return fmt.Sprintf("$%.2f", amount)
    },
    "percentage": func(value float64) string {
        return fmt.Sprintf("%.1f%%", value*100)
    },
    "formatDate": func(date time.Time) string {
        return date.Format("January 2, 2006")
    },
    "colorClass": func(value float64) string {
        if value > 0 {
            return "positive"
        } else if value < 0 {
            return "negative"
        }
        return "neutral"
    },
}
```

## ⏰ Programador de Reportes

### Configuración de Cron

```go
type ScheduleConfig struct {
    ID          string        `json:"id"`
    Name        string        `json:"name"`
    ReportType  ReportType    `json:"report_type"`
    Schedule    string        `json:"schedule"` // Cron expression
    Format      OutputFormat  `json:"format"`
    Recipients  []string      `json:"recipients"`
    Enabled     bool          `json:"enabled"`
    LastRun     *time.Time    `json:"last_run,omitempty"`
    NextRun     *time.Time    `json:"next_run,omitempty"`
}

// Ejemplos de programación
var scheduleExamples = map[string]string{
    "daily":    "0 9 * * *",        // 9:00 AM diario
    "weekly":   "0 9 * * 1",        // 9:00 AM lunes
    "monthly":  "0 9 1 * *",        // 9:00 AM primer día del mes
    "quarterly": "0 9 1 */3 *",     // 9:00 AM primer día del trimestre
    "yearly":   "0 9 1 1 *",        // 9:00 AM 1 de enero
}
```

## 🔐 Seguridad

### Medidas Implementadas

- **Access Control**: Control de acceso por roles
- **Data Encryption**: Encriptación de reportes sensibles
- **Audit Trail**: Registro de generación de reportes
- **Rate Limiting**: Limitación de generación de reportes
- **File Security**: Validación de archivos generados
- **Email Security**: Encriptación de envío por email

## 🧪 Testing

```bash
# Tests unitarios
go test ./internal/...

# Tests de generación de reportes
go test ./internal/generator/...

# Tests de plantillas
go test ./internal/templates/...

# Tests de análisis
go test ./internal/analytics/...
```

### Tests de Generación

```go
func TestBalanceReportGeneration(t *testing.T) {
    generator := NewReportGenerator(mockDB, mockStorage)
    
    config := ReportConfig{
        Type:   BalanceReport,
        Format: PDFFormat,
        Period: DateRange{
            Start: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
            End:   time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
        },
    }
    
    report, err := generator.Generate(config)
    assert.NoError(t, err)
    assert.NotNil(t, report)
    assert.Equal(t, PDFFormat, report.Format)
}
```

## 📊 Monitoreo

### Métricas Específicas

- **Report Generation Time**: Tiempo de generación
- **Report Size**: Tamaño de reportes
- **Generation Success Rate**: Tasa de éxito
- **Template Rendering Time**: Tiempo de renderizado
- **Storage Usage**: Uso de almacenamiento
- **Queue Length**: Longitud de cola
- **Email Delivery Rate**: Tasa de entrega de emails

### Dashboard de Métricas

```go
type ReportMetrics struct {
    TotalReports     int64   `json:"total_reports"`
    SuccessRate      float64 `json:"success_rate"`
    AvgGenerationTime float64 `json:"avg_generation_time"`
    QueueLength      int     `json:"queue_length"`
    StorageUsed      int64   `json:"storage_used"`
    PopularFormats   map[OutputFormat]int `json:"popular_formats"`
}
```

## 🚀 Despliegue

### Variables de Producción

```env
GIN_MODE=release
LOG_LEVEL=warn
PDF_TIMEOUT=60s
MAX_CONCURRENT_REPORTS=10
CLEANUP_INTERVAL=24h
RETENTION_DAYS=30
```

```bash
# Build de producción
CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-w -s" \
  -o report-service cmd/main.go
```

## 🔧 Configuración Avanzada

### Optimización de PDF

```go
type PDFOptions struct {
    DPI           int    `json:"dpi"`
    PageSize      string `json:"page_size"`
    Orientation   string `json:"orientation"`
    MarginTop     string `json:"margin_top"`
    MarginBottom  string `json:"margin_bottom"`
    MarginLeft    string `json:"margin_left"`
    MarginRight   string `json:"margin_right"`
    EnableTOC     bool   `json:"enable_toc"`
    EnableFooter  bool   `json:"enable_footer"`
    EnableHeader  bool   `json:"enable_header"`
}
```

### Configuración de Caché

```go
type CacheConfig struct {
    TTL              time.Duration `json:"ttl"`
    MaxSize          int64         `json:"max_size"`
    CleanupInterval  time.Duration `json:"cleanup_interval"`
    EnableCompression bool         `json:"enable_compression"`
}
```

---

**Report Service** - Reportes financieros inteligentes 📊📈