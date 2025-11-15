package service

import (
	"context"

	domexchange "github.com/fintrack/exchange-service/internal/core/domain/entities/exchange"
)

// ExchangeServiceInterface define la interfaz para el servicio de exchange
type ExchangeServiceInterface interface {
	GetDolarOficial(ctx context.Context) (*domexchange.ExchangeRate, error)
	HealthCheck(ctx context.Context) error
}
