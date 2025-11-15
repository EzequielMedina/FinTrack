package exchange

import (
	"context"

	domexchange "github.com/fintrack/exchange-service/internal/core/domain/entities/exchange"
)

// ExchangeProvider define la interfaz para obtener tasas de cambio
type ExchangeProvider interface {
	// GetDolarOficial obtiene la cotización del dólar oficial
	GetDolarOficial(ctx context.Context) (*domexchange.ExchangeRate, error)

	// IsHealthy verifica si el proveedor está disponible
	IsHealthy(ctx context.Context) bool
}

// ExchangeRateClient define la interfaz para clients HTTP de APIs externas
type ExchangeRateClient interface {
	// GetDolarOficial obtiene datos del dólar oficial desde la API externa
	GetDolarOficial(ctx context.Context) (*domexchange.DolarAPIResponse, error)

	// Health verifica el estado de la API externa
	Health(ctx context.Context) error
}
