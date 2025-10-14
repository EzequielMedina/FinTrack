package exchangehandler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	domerrors "github.com/fintrack/exchange-service/internal/core/errors"
	"github.com/fintrack/exchange-service/internal/core/service"
	"github.com/fintrack/exchange-service/internal/infrastructure/entrypoints/handlers/exchange/dto"
)

// ExchangeHandler maneja las peticiones HTTP relacionadas con tasas de cambio
type ExchangeHandler struct {
	exchangeService service.ExchangeServiceInterface
}

// New crea una nueva instancia del handler de exchange
func New(exchangeService service.ExchangeServiceInterface) *ExchangeHandler {
	return &ExchangeHandler{
		exchangeService: exchangeService,
	}
}

// GetDolarOficial maneja GET /api/exchange/dolar-oficial
func (h *ExchangeHandler) GetDolarOficial(c *gin.Context) {
	ctx := c.Request.Context()

	// Obtener cotización del servicio
	exchangeRate, err := h.exchangeService.GetDolarOficial(ctx)
	if err != nil {
		h.handleError(c, err)
		return
	}

	// Convertir a DTO de respuesta
	response := dto.NewExchangeRateResponse(exchangeRate)

	c.JSON(http.StatusOK, response)
}

// HealthCheck maneja GET /health
func (h *ExchangeHandler) HealthCheck(c *gin.Context) {
	ctx := c.Request.Context()

	err := h.exchangeService.HealthCheck(ctx)
	if err != nil {
		response := dto.NewHealthResponse("unhealthy", "exchange-service")
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	response := dto.NewHealthResponse("healthy", "exchange-service")
	c.JSON(http.StatusOK, response)
}

// handleError maneja los errores y devuelve la respuesta HTTP apropiada
func (h *ExchangeHandler) handleError(c *gin.Context, err error) {
	var statusCode int
	var message string

	switch {
	case errors.Is(err, domerrors.ErrAPIUnavailable):
		statusCode = http.StatusServiceUnavailable
		message = "El servicio de cotizaciones no está disponible temporalmente"
	case errors.Is(err, domerrors.ErrAPITimeout):
		statusCode = http.StatusGatewayTimeout
		message = "Timeout al obtener cotizaciones"
	case errors.Is(err, domerrors.ErrAPIRateLimit):
		statusCode = http.StatusTooManyRequests
		message = "Límite de peticiones excedido"
	case errors.Is(err, domerrors.ErrInvalidExchangeRate):
		statusCode = http.StatusBadGateway
		message = "Datos de cotización inválidos"
	case errors.Is(err, domerrors.ErrEmptyResponse):
		statusCode = http.StatusBadGateway
		message = "Respuesta vacía del proveedor de cotizaciones"
	case errors.Is(err, domerrors.ErrNetworkFailure):
		statusCode = http.StatusBadGateway
		message = "Error de red al obtener cotizaciones"
	default:
		statusCode = http.StatusInternalServerError
		message = "Error interno del servidor"
	}

	errorResponse := dto.NewErrorResponse(err, message, statusCode)
	c.JSON(statusCode, errorResponse)
}
