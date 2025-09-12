# Exchange Service - FinTrack

## 📋 Descripción

Microservicio encargado de la gestión de tipos de cambio y conversiones de moneda en la plataforma FinTrack. Proporciona tasas de cambio en tiempo real y conversiones entre diferentes divisas.

## 🛠️ Tecnologías

- **Lenguaje**: Go 1.24+
- **Framework**: Gin/Echo (HTTP Router)
- **Base de Datos**: MySQL 8.0
- **APIs Externas**: Exchange Rate APIs
- **Cache**: Redis (opcional)
- **Contenedor**: Docker multi-stage
- **Arquitectura**: Clean Architecture

## 🏗️ Arquitectura

### Estructura del Proyecto

```
exchange-service/
├── cmd/
│   └── main.go              # Punto de entrada
├── internal/
│   ├── config/              # Configuración
│   ├── domain/              # Entidades de dominio
│   ├── handlers/            # HTTP handlers
│   ├── repository/          # Capa de datos
│   ├── service/             # Lógica de negocio
│   ├── providers/           # Proveedores de exchange rates
│   ├── cache/               # Sistema de caché
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

# APIs de Exchange Rate
EXCHANGE_API_KEY=your-exchange-api-key
EXCHANGE_API_URL=https://api.exchangerate-api.com/v4
BACKUP_API_KEY=your-backup-api-key
BACKUP_API_URL=https://api.fixer.io

# Cache (Redis)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
CACHE_TTL=300

# Configuración
UPDATE_INTERVAL=300s
BASE_CURRENCY=USD
SUPPORTED_CURRENCIES=USD,EUR,GBP,JPY,CAD,AUD,CHF,CNY,MXN

# Servidor
PORT=8080
GIN_MODE=debug

# Logging
LOG_LEVEL=info
```

### Comandos de Desarrollo

```bash
# Navegar al servicio
cd backend/services/exchange-service

# Instalar dependencias
go mod download

# Ejecutar en modo desarrollo
go run cmd/main.go

# Build del binario
go build -o bin/exchange-service cmd/main.go

# Tests
go test ./...
```

## 🐳 Docker

```bash
# Build de la imagen
docker build -t fintrack-exchange-service .

# Docker Compose
docker-compose up exchange-service

# Con dependencias
docker-compose up mysql exchange-service
```

## 📡 API Endpoints

### Tipos de Cambio

```http
GET    /api/exchange/rates             # Obtener todas las tasas
GET    /api/exchange/rates/{currency}  # Tasa específica
GET    /api/exchange/historical        # Datos históricos
POST   /api/exchange/refresh           # Actualizar tasas
```

### Conversiones

```http
POST   /api/exchange/convert           # Convertir moneda
GET    /api/exchange/convert           # Conversión con parámetros
POST   /api/exchange/batch-convert     # Conversiones múltiples
```

### Configuración

```http
GET    /api/exchange/currencies        # Monedas soportadas
GET    /api/exchange/config            # Configuración del servicio
PUT    /api/exchange/config            # Actualizar configuración
```

### Análisis

```http
GET    /api/exchange/trends            # Tendencias de mercado
GET    /api/exchange/volatility        # Análisis de volatilidad
GET    /api/exchange/alerts            # Alertas de cambio
```

### Health Check

```http
GET /health
```

### Ejemplos de Uso

```bash
# Obtener tasas de cambio actuales
curl -X GET http://localhost:8087/api/exchange/rates \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Obtener tasa específica
curl -X GET http://localhost:8087/api/exchange/rates/EUR \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Convertir moneda
curl -X POST http://localhost:8087/api/exchange/convert \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "amount": 100,
    "fromCurrency": "USD",
    "toCurrency": "EUR"
  }'

# Conversión con parámetros GET
curl -X GET "http://localhost:8087/api/exchange/convert?amount=100&from=USD&to=EUR" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Datos históricos
curl -X GET "http://localhost:8087/api/exchange/historical?currency=EUR&days=30" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Health check
curl http://localhost:8087/health
```

## 💱 Proveedores de Exchange Rate

### APIs Soportadas

```go
type ExchangeProvider interface {
    GetRates(baseCurrency string) (*ExchangeRates, error)
    GetHistoricalRates(currency string, date time.Time) (*ExchangeRates, error)
    IsAvailable() bool
}

// Proveedores implementados
type Providers struct {
    Primary   ExchangeProvider // ExchangeRate-API
    Secondary ExchangeProvider // Fixer.io
    Tertiary  ExchangeProvider // CurrencyAPI
}
```

### Fallback Strategy

```go
func (s *ExchangeService) GetRatesWithFallback() (*ExchangeRates, error) {
    providers := []ExchangeProvider{
        s.providers.Primary,
        s.providers.Secondary,
        s.providers.Tertiary,
    }
    
    for _, provider := range providers {
        if rates, err := provider.GetRates(s.baseCurrency); err == nil {
            return rates, nil
        }
    }
    
    return nil, errors.New("all exchange providers failed")
}
```

## 💾 Sistema de Caché

### Estrategia de Cache

```go
type CacheStrategy struct {
    TTL           time.Duration // 5 minutos
    RefreshBuffer time.Duration // 1 minuto antes de expirar
    MaxRetries    int          // 3 intentos
}

// Niveles de caché
type CacheLevel struct {
    Memory    *sync.Map        // Cache en memoria
    Redis     *redis.Client    // Cache distribuido
    Database  *sql.DB          // Fallback en BD
}
```

### Cache Keys

```go
const (
    RatesCacheKey     = "exchange:rates:%s"           // exchange:rates:USD
    ConvertCacheKey   = "exchange:convert:%s:%s:%f"    // exchange:convert:USD:EUR:100
    HistoricalKey     = "exchange:historical:%s:%s"   // exchange:historical:EUR:2024-01-15
)
```

## 🔐 Seguridad

### Medidas Implementadas

- **API Key Protection**: Rotación de claves de API
- **Rate Limiting**: Limitación de requests
- **Data Validation**: Validación de monedas soportadas
- **Audit Trail**: Registro de conversiones
- **Encryption**: Encriptación de datos sensibles
- **Monitoring**: Monitoreo de APIs externas

## 🧪 Testing

```bash
# Tests unitarios
go test ./internal/...

# Tests de integración
go test ./tests/integration/...

# Tests de proveedores (con mocks)
go test ./internal/providers/...

# Tests de caché
go test ./internal/cache/...
```

### Mock de APIs Externas

```go
type MockExchangeProvider struct {
    rates map[string]float64
    delay time.Duration
}

func (m *MockExchangeProvider) GetRates(base string) (*ExchangeRates, error) {
    time.Sleep(m.delay) // Simular latencia
    return &ExchangeRates{
        Base:  base,
        Rates: m.rates,
        Date:  time.Now(),
    }, nil
}
```

## 📊 Monitoreo

### Métricas Específicas

- **API Calls**: Llamadas a APIs externas
- **Cache Hit Rate**: Tasa de aciertos de caché
- **Response Time**: Tiempo de respuesta
- **Conversion Volume**: Volumen de conversiones
- **Provider Availability**: Disponibilidad de proveedores
- **Rate Update Frequency**: Frecuencia de actualización
- **Error Rate**: Tasa de errores por proveedor

### Alertas

```go
type Alert struct {
    Type        AlertType
    Currency    string
    Threshold   float64
    CurrentRate float64
    Message     string
}

const (
    VolatilityAlert AlertType = "volatility"
    ProviderDown   AlertType = "provider_down"
    RateStale      AlertType = "rate_stale"
)
```

## 🚀 Despliegue

### Variables de Producción

```env
GIN_MODE=release
LOG_LEVEL=warn
UPDATE_INTERVAL=60s
CACHE_TTL=180
MAX_RETRIES=5
TIMEOUT=10s
```

```bash
# Build de producción
CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-w -s" \
  -o exchange-service cmd/main.go
```

## 🔧 Configuración Avanzada

### Monedas Soportadas

```json
{
  "supported_currencies": {
    "USD": {"name": "US Dollar", "symbol": "$", "decimals": 2},
    "EUR": {"name": "Euro", "symbol": "€", "decimals": 2},
    "GBP": {"name": "British Pound", "symbol": "£", "decimals": 2},
    "JPY": {"name": "Japanese Yen", "symbol": "¥", "decimals": 0},
    "BTC": {"name": "Bitcoin", "symbol": "₿", "decimals": 8}
  }
}
```

### Rate Limiting

```go
type RateLimit struct {
    RequestsPerMinute int
    BurstSize        int
    WindowSize       time.Duration
}

var limits = map[string]RateLimit{
    "free":     {RequestsPerMinute: 100, BurstSize: 10},
    "premium":  {RequestsPerMinute: 1000, BurstSize: 50},
    "enterprise": {RequestsPerMinute: 10000, BurstSize: 100},
}
```

---

**Exchange Service** - Tipos de cambio en tiempo real 💱📈