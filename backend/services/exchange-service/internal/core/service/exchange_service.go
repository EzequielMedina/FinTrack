package service

import (
	"context"
	"fmt"
	"log"

	domexchange "github.com/fintrack/exchange-service/internal/core/domain/entities/exchange"
	domerrors "github.com/fintrack/exchange-service/internal/core/errors"
	"github.com/fintrack/exchange-service/internal/core/providers/exchange"
)

// ExchangeService implementa la lógica de negocio para tasas de cambio
type ExchangeService struct {
	exchangeProvider exchange.ExchangeProvider
}

// NewExchangeService crea una nueva instancia del servicio de exchange
func NewExchangeService(exchangeProvider exchange.ExchangeProvider) *ExchangeService {
	return &ExchangeService{
		exchangeProvider: exchangeProvider,
	}
}

// GetDolarOficial obtiene la cotización actual del dólar oficial
func (s *ExchangeService) GetDolarOficial(ctx context.Context) (*domexchange.ExchangeRate, error) {
	log.Println("ExchangeService: obteniendo cotización del dólar oficial")

	// Verificar que el proveedor esté disponible
	if !s.exchangeProvider.IsHealthy(ctx) {
		log.Println("ExchangeService: proveedor de exchange no disponible")
		return nil, domerrors.ErrAPIUnavailable
	}

	// Obtener la cotización del proveedor
	exchangeRate, err := s.exchangeProvider.GetDolarOficial(ctx)
	if err != nil {
		log.Printf("ExchangeService: error obteniendo cotización: %v", err)
		return nil, fmt.Errorf("error obteniendo cotización del dólar: %w", err)
	}

	// Validar que los datos sean válidos
	if exchangeRate == nil || !exchangeRate.IsValid() {
		log.Println("ExchangeService: datos de cotización inválidos")
		return nil, domerrors.ErrInvalidExchangeRate
	}

	log.Printf("ExchangeService: cotización obtenida exitosamente - Compra: %.2f, Venta: %.2f",
		exchangeRate.Compra, exchangeRate.Venta)

	return exchangeRate, nil
}

// HealthCheck verifica el estado del servicio
func (s *ExchangeService) HealthCheck(ctx context.Context) error {
	log.Println("ExchangeService: verificando health check")

	if !s.exchangeProvider.IsHealthy(ctx) {
		return domerrors.ErrAPIUnavailable
	}

	return nil
}
