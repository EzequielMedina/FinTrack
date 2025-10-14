package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	domexchange "github.com/fintrack/exchange-service/internal/core/domain/entities/exchange"
	domerrors "github.com/fintrack/exchange-service/internal/core/errors"
)

// DolarAPIClient implementa el cliente HTTP para la API de DolarAPI
type DolarAPIClient struct {
	baseURL    string
	httpClient *http.Client
	timeout    time.Duration
}

// NewDolarAPIClient crea una nueva instancia del cliente de DolarAPI
func NewDolarAPIClient(baseURL string, timeout time.Duration) *DolarAPIClient {
	return &DolarAPIClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		timeout: timeout,
	}
}

// GetDolarOficial obtiene la cotización del dólar oficial desde DolarAPI
func (c *DolarAPIClient) GetDolarOficial(ctx context.Context) (*domexchange.DolarAPIResponse, error) {
	url := fmt.Sprintf("%s/v1/dolares/oficial", c.baseURL)

	// Crear request con contexto
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creando request: %w", err)
	}

	// Agregar headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "FinTrack-ExchangeService/1.0")

	// Realizar la petición
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error realizando petición: %w", domerrors.ErrNetworkFailure)
	}
	defer resp.Body.Close()

	// Verificar status code
	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusTooManyRequests:
			return nil, domerrors.ErrAPIRateLimit
		case http.StatusServiceUnavailable, http.StatusBadGateway:
			return nil, domerrors.ErrAPIUnavailable
		default:
			return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
		}
	}

	// Decodificar respuesta JSON
	var apiResponse domexchange.DolarAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta JSON: %w", domerrors.ErrInvalidAPIResponse)
	}

	// Validar que la respuesta no esté vacía
	if apiResponse.Casa == "" || apiResponse.Nombre == "" {
		return nil, domerrors.ErrEmptyResponse
	}

	return &apiResponse, nil
}

// Health verifica el estado de la API de DolarAPI
func (c *DolarAPIClient) Health(ctx context.Context) error {
	url := fmt.Sprintf("%s/v1/dolares/oficial", c.baseURL)

	// Crear request con contexto y timeout corto para health check
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
	if err != nil {
		return fmt.Errorf("error creando health check request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return domerrors.ErrAPIUnavailable
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return domerrors.ErrAPIUnavailable
	}

	return nil
}
