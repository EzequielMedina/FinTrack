package providers

import (
	"context"
	"fmt"
	"log"

	domexchange "github.com/fintrack/exchange-service/internal/core/domain/entities/exchange"
	domerrors "github.com/fintrack/exchange-service/internal/core/errors"
	"github.com/fintrack/exchange-service/internal/core/providers/exchange"
)

// DolarAPIProvider implementa ExchangeProvider usando DolarAPI
type DolarAPIProvider struct {
	client exchange.ExchangeRateClient
}

// NewDolarAPIProvider crea una nueva instancia del proveedor de DolarAPI
func NewDolarAPIProvider(client exchange.ExchangeRateClient) exchange.ExchangeProvider {
	return &DolarAPIProvider{
		client: client,
	}
}

// GetDolarOficial obtiene la cotización del dólar oficial
func (p *DolarAPIProvider) GetDolarOficial(ctx context.Context) (*domexchange.ExchangeRate, error) {
	log.Println("DolarAPIProvider: obteniendo cotización del dólar oficial")

	// Obtener datos de la API
	apiResponse, err := p.client.GetDolarOficial(ctx)
	if err != nil {
		log.Printf("DolarAPIProvider: error llamando API: %v", err)
		return nil, fmt.Errorf("error obteniendo datos de DolarAPI: %w", err)
	}

	// Convertir respuesta de API a entidad de dominio
	exchangeRate, err := apiResponse.ConvertToDomain()
	if err != nil {
		log.Printf("DolarAPIProvider: error convirtiendo respuesta: %v", err)
		return nil, fmt.Errorf("error procesando respuesta de DolarAPI: %w", domerrors.ErrInvalidAPIResponse)
	}

	// Validar datos
	if !exchangeRate.IsValid() {
		log.Println("DolarAPIProvider: datos de cotización inválidos")
		return nil, domerrors.ErrInvalidExchangeRate
	}

	log.Printf("DolarAPIProvider: cotización obtenida - %s: Compra=%.2f, Venta=%.2f",
		exchangeRate.Nombre, exchangeRate.Compra, exchangeRate.Venta)

	return exchangeRate, nil
}

// IsHealthy verifica si el proveedor está disponible
func (p *DolarAPIProvider) IsHealthy(ctx context.Context) bool {
	log.Println("DolarAPIProvider: verificando health check")

	err := p.client.Health(ctx)
	if err != nil {
		log.Printf("DolarAPIProvider: health check fallido: %v", err)
		return false
	}

	log.Println("DolarAPIProvider: health check exitoso")
	return true
}
