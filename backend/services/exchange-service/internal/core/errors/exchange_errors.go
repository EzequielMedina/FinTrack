package domerrors

import "errors"

var (
	// Errores de validación
	ErrInvalidExchangeRate = errors.New("invalid exchange rate data")
	ErrEmptyResponse       = errors.New("empty response from exchange API")

	// Errores de API externa
	ErrAPIUnavailable     = errors.New("exchange API is unavailable")
	ErrAPITimeout         = errors.New("exchange API request timeout")
	ErrAPIRateLimit       = errors.New("exchange API rate limit exceeded")
	ErrInvalidAPIResponse = errors.New("invalid response format from exchange API")

	// Errores de configuración
	ErrMissingAPIConfig = errors.New("missing exchange API configuration")
	ErrInvalidURL       = errors.New("invalid API URL")

	// Errores de red
	ErrNetworkFailure = errors.New("network failure while calling exchange API")
	ErrDNSFailure     = errors.New("DNS resolution failure for exchange API")
)
