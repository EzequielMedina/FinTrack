package dto

import (
	"time"

	domexchange "github.com/fintrack/exchange-service/internal/core/domain/entities/exchange"
)

// ExchangeRateResponse DTO para la respuesta del endpoint
type ExchangeRateResponse struct {
	Compra             float64   `json:"compra"`
	Venta              float64   `json:"venta"`
	Casa               string    `json:"casa"`
	Nombre             string    `json:"nombre"`
	Moneda             string    `json:"moneda"`
	FechaActualizacion time.Time `json:"fechaActualizacion"`
	Spread             float64   `json:"spread"`
	SpreadPercentage   float64   `json:"spreadPercentage"`
}

// ErrorResponse DTO para respuestas de error
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// HealthResponse DTO para health check
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
}

// NewExchangeRateResponse convierte entidad de dominio a DTO de respuesta
func NewExchangeRateResponse(exchangeRate *domexchange.ExchangeRate) *ExchangeRateResponse {
	return &ExchangeRateResponse{
		Compra:             exchangeRate.Compra,
		Venta:              exchangeRate.Venta,
		Casa:               exchangeRate.Casa,
		Nombre:             exchangeRate.Nombre,
		Moneda:             exchangeRate.Moneda,
		FechaActualizacion: exchangeRate.FechaActualizacion,
		Spread:             exchangeRate.GetSpread(),
		SpreadPercentage:   exchangeRate.GetSpreadPercentage(),
	}
}

// NewErrorResponse crea un nuevo DTO de error
func NewErrorResponse(err error, message string, code int) *ErrorResponse {
	return &ErrorResponse{
		Error:   err.Error(),
		Message: message,
		Code:    code,
	}
}

// NewHealthResponse crea un nuevo DTO de health check
func NewHealthResponse(status string, service string) *HealthResponse {
	return &HealthResponse{
		Status:    status,
		Timestamp: time.Now(),
		Service:   service,
	}
}
