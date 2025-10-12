package service

import (
	"context"
	"testing"
	"time"

	domexchange "github.com/fintrack/exchange-service/internal/core/domain/entities/exchange"
	domerrors "github.com/fintrack/exchange-service/internal/core/errors"
	"github.com/fintrack/exchange-service/internal/core/providers/exchange"
)

// MockExchangeProvider implementa un mock del ExchangeProvider para testing
type MockExchangeProvider struct {
	shouldReturnError bool
	isHealthy         bool
	exchangeRate      *domexchange.ExchangeRate
}

func NewMockExchangeProvider() *MockExchangeProvider {
	return &MockExchangeProvider{
		shouldReturnError: false,
		isHealthy:         true,
		exchangeRate: &domexchange.ExchangeRate{
			Compra:             900.50,
			Venta:              950.75,
			Casa:               "Banco Nación",
			Nombre:             "Oficial",
			Moneda:             "USD",
			FechaActualizacion: time.Now(),
		},
	}
}

func (m *MockExchangeProvider) GetDolarOficial(ctx context.Context) (*domexchange.ExchangeRate, error) {
	if m.shouldReturnError {
		return nil, domerrors.ErrAPIUnavailable
	}
	return m.exchangeRate, nil
}

func (m *MockExchangeProvider) IsHealthy(ctx context.Context) bool {
	return m.isHealthy
}

func (m *MockExchangeProvider) SetShouldReturnError(value bool) {
	m.shouldReturnError = value
}

func (m *MockExchangeProvider) SetIsHealthy(value bool) {
	m.isHealthy = value
}

func (m *MockExchangeProvider) SetExchangeRate(rate *domexchange.ExchangeRate) {
	m.exchangeRate = rate
}

// Verify interface compliance
var _ exchange.ExchangeProvider = (*MockExchangeProvider)(nil)

func TestNewExchangeService(t *testing.T) {
	mockProvider := NewMockExchangeProvider()
	service := NewExchangeService(mockProvider)

	if service == nil {
		t.Fatal("Expected service to be created, got nil")
	}

	if service.exchangeProvider != mockProvider {
		t.Fatal("Expected service to have the provided exchange provider")
	}
}

func TestGetDolarOficial_Success(t *testing.T) {
	mockProvider := NewMockExchangeProvider()
	service := NewExchangeService(mockProvider)
	ctx := context.Background()

	result, err := service.GetDolarOficial(ctx)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected exchange rate, got nil")
	}

	if result.Compra != 900.50 {
		t.Errorf("Expected Compra to be 900.50, got %f", result.Compra)
	}

	if result.Venta != 950.75 {
		t.Errorf("Expected Venta to be 950.75, got %f", result.Venta)
	}

	if result.Casa != "Banco Nación" {
		t.Errorf("Expected Casa to be 'Banco Nación', got %s", result.Casa)
	}
}

func TestGetDolarOficial_ProviderUnhealthy(t *testing.T) {
	mockProvider := NewMockExchangeProvider()
	mockProvider.SetIsHealthy(false)
	service := NewExchangeService(mockProvider)
	ctx := context.Background()

	result, err := service.GetDolarOficial(ctx)

	if err == nil {
		t.Fatal("Expected error when provider is unhealthy, got nil")
	}

	if result != nil {
		t.Fatal("Expected nil result when provider is unhealthy, got exchange rate")
	}

	if err != domerrors.ErrAPIUnavailable {
		t.Errorf("Expected ErrAPIUnavailable, got %v", err)
	}
}

func TestGetDolarOficial_ProviderError(t *testing.T) {
	mockProvider := NewMockExchangeProvider()
	mockProvider.SetShouldReturnError(true)
	service := NewExchangeService(mockProvider)
	ctx := context.Background()

	result, err := service.GetDolarOficial(ctx)

	if err == nil {
		t.Fatal("Expected error when provider returns error, got nil")
	}

	if result != nil {
		t.Fatal("Expected nil result when provider returns error, got exchange rate")
	}
}

func TestGetDolarOficial_InvalidExchangeRate(t *testing.T) {
	mockProvider := NewMockExchangeProvider()
	// Set invalid exchange rate (zero values)
	invalidRate := &domexchange.ExchangeRate{
		Compra: 0,
		Venta:  0,
		Casa:   "",
		Nombre: "",
		Moneda: "",
	}
	mockProvider.SetExchangeRate(invalidRate)
	service := NewExchangeService(mockProvider)
	ctx := context.Background()

	result, err := service.GetDolarOficial(ctx)

	if err == nil {
		t.Fatal("Expected error for invalid exchange rate, got nil")
	}

	if result != nil {
		t.Fatal("Expected nil result for invalid exchange rate, got exchange rate")
	}

	if err != domerrors.ErrInvalidExchangeRate {
		t.Errorf("Expected ErrInvalidExchangeRate, got %v", err)
	}
}

func TestHealthCheck_Success(t *testing.T) {
	mockProvider := NewMockExchangeProvider()
	service := NewExchangeService(mockProvider)
	ctx := context.Background()

	err := service.HealthCheck(ctx)

	if err != nil {
		t.Fatalf("Expected no error for healthy provider, got %v", err)
	}
}

func TestHealthCheck_ProviderUnhealthy(t *testing.T) {
	mockProvider := NewMockExchangeProvider()
	mockProvider.SetIsHealthy(false)
	service := NewExchangeService(mockProvider)
	ctx := context.Background()

	err := service.HealthCheck(ctx)

	if err == nil {
		t.Fatal("Expected error for unhealthy provider, got nil")
	}

	if err != domerrors.ErrAPIUnavailable {
		t.Errorf("Expected ErrAPIUnavailable, got %v", err)
	}
}
