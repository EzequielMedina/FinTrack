package exchangehandler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	domexchange "github.com/fintrack/exchange-service/internal/core/domain/entities/exchange"
	domerrors "github.com/fintrack/exchange-service/internal/core/errors"
	"github.com/fintrack/exchange-service/internal/infrastructure/entrypoints/handlers/exchange/dto"
)

// MockExchangeService implementa un mock del ExchangeService para testing
type MockExchangeService struct {
	shouldReturnError bool
	isHealthy         bool
	exchangeRate      *domexchange.ExchangeRate
}

func NewMockExchangeService() *MockExchangeService {
	return &MockExchangeService{
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

func (m *MockExchangeService) GetDolarOficial(ctx context.Context) (*domexchange.ExchangeRate, error) {
	if m.shouldReturnError {
		return nil, domerrors.ErrAPIUnavailable
	}
	return m.exchangeRate, nil
}

func (m *MockExchangeService) HealthCheck(ctx context.Context) error {
	if !m.isHealthy {
		return domerrors.ErrAPIUnavailable
	}
	return nil
}

func (m *MockExchangeService) SetShouldReturnError(value bool) {
	m.shouldReturnError = value
}

func (m *MockExchangeService) SetIsHealthy(value bool) {
	m.isHealthy = value
}

func (m *MockExchangeService) SetExchangeRate(rate *domexchange.ExchangeRate) {
	m.exchangeRate = rate
}

// Helper function to create a test Gin context
func createTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func TestGetDolarOficial_Success(t *testing.T) {
	mockService := NewMockExchangeService()
	handler := New(mockService)

	c, w := createTestContext()
	req, _ := http.NewRequest("GET", "/api/exchange/dolar-oficial", nil)
	c.Request = req

	handler.GetDolarOficial(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response dto.ExchangeRateResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Error unmarshaling response: %v", err)
	}

	if response.Compra != 900.50 {
		t.Errorf("Expected Compra to be 900.50, got %f", response.Compra)
	}

	if response.Venta != 950.75 {
		t.Errorf("Expected Venta to be 950.75, got %f", response.Venta)
	}

	if response.Casa != "Banco Nación" {
		t.Errorf("Expected Casa to be 'Banco Nación', got %s", response.Casa)
	}

	if response.Nombre != "Oficial" {
		t.Errorf("Expected Nombre to be 'Oficial', got %s", response.Nombre)
	}

	// Verificar que se calculó el spread
	expectedSpread := 950.75 - 900.50
	if response.Spread != expectedSpread {
		t.Errorf("Expected Spread to be %f, got %f", expectedSpread, response.Spread)
	}
}

func TestGetDolarOficial_ServiceError(t *testing.T) {
	mockService := NewMockExchangeService()
	mockService.SetShouldReturnError(true)
	handler := New(mockService)

	c, w := createTestContext()
	req, _ := http.NewRequest("GET", "/api/exchange/dolar-oficial", nil)
	c.Request = req

	handler.GetDolarOficial(c)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected status %d, got %d", http.StatusServiceUnavailable, w.Code)
	}

	var response dto.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Error unmarshaling error response: %v", err)
	}

	if response.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected error code %d, got %d", http.StatusServiceUnavailable, response.Code)
	}
}

func TestHealthCheck_Healthy(t *testing.T) {
	mockService := NewMockExchangeService()
	handler := New(mockService)

	c, w := createTestContext()
	req, _ := http.NewRequest("GET", "/health", nil)
	c.Request = req

	handler.HealthCheck(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response dto.HealthResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Error unmarshaling health response: %v", err)
	}

	if response.Status != "healthy" {
		t.Errorf("Expected status 'healthy', got %s", response.Status)
	}

	if response.Service != "exchange-service" {
		t.Errorf("Expected service 'exchange-service', got %s", response.Service)
	}
}

func TestHealthCheck_Unhealthy(t *testing.T) {
	mockService := NewMockExchangeService()
	mockService.SetIsHealthy(false)
	handler := New(mockService)

	c, w := createTestContext()
	req, _ := http.NewRequest("GET", "/health", nil)
	c.Request = req

	handler.HealthCheck(c)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected status %d, got %d", http.StatusServiceUnavailable, w.Code)
	}

	var response dto.HealthResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Error unmarshaling health response: %v", err)
	}

	if response.Status != "unhealthy" {
		t.Errorf("Expected status 'unhealthy', got %s", response.Status)
	}
}

func TestHandleError_DifferentErrorTypes(t *testing.T) {
	mockService := NewMockExchangeService()
	handler := New(mockService)

	testCases := []struct {
		name           string
		error          error
		expectedStatus int
	}{
		{
			name:           "API Unavailable",
			error:          domerrors.ErrAPIUnavailable,
			expectedStatus: http.StatusServiceUnavailable,
		},
		{
			name:           "API Timeout",
			error:          domerrors.ErrAPITimeout,
			expectedStatus: http.StatusGatewayTimeout,
		},
		{
			name:           "Rate Limit",
			error:          domerrors.ErrAPIRateLimit,
			expectedStatus: http.StatusTooManyRequests,
		},
		{
			name:           "Invalid Exchange Rate",
			error:          domerrors.ErrInvalidExchangeRate,
			expectedStatus: http.StatusBadGateway,
		},
		{
			name:           "Network Failure",
			error:          domerrors.ErrNetworkFailure,
			expectedStatus: http.StatusBadGateway,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, w := createTestContext()

			handler.handleError(c, tc.error)

			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, w.Code)
			}

			var response dto.ErrorResponse
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("Error unmarshaling error response: %v", err)
			}

			if response.Code != tc.expectedStatus {
				t.Errorf("Expected error code %d, got %d", tc.expectedStatus, response.Code)
			}
		})
	}
}
