package exchange

import (
	"time"
)

// ExchangeRate representa la cotización del dólar oficial de la API DolarAPI
type ExchangeRate struct {
	Compra             float64   `json:"compra"`
	Venta              float64   `json:"venta"`
	Casa               string    `json:"casa"`
	Nombre             string    `json:"nombre"`
	Moneda             string    `json:"moneda"`
	FechaActualizacion time.Time `json:"fechaActualizacion"`
}

// DolarAPIResponse estructura que mapea exactamente la respuesta de DolarAPI
type DolarAPIResponse struct {
	Compra             float64 `json:"compra"`
	Venta              float64 `json:"venta"`
	Casa               string  `json:"casa"`
	Nombre             string  `json:"nombre"`
	Moneda             string  `json:"moneda"`
	FechaActualizacion string  `json:"fechaActualizacion"`
}

// ConvertToDomain convierte la respuesta de DolarAPI a entidad de dominio
func (r *DolarAPIResponse) ConvertToDomain() (*ExchangeRate, error) {
	// Parsear la fecha de actualización
	fechaActualizacion, err := time.Parse("2006-01-02T15:04:05.000Z", r.FechaActualizacion)
	if err != nil {
		// Intentar con formato alternativo si el primero falla
		fechaActualizacion, err = time.Parse("2006-01-02T15:04:05Z", r.FechaActualizacion)
		if err != nil {
			return nil, err
		}
	}

	return &ExchangeRate{
		Compra:             r.Compra,
		Venta:              r.Venta,
		Casa:               r.Casa,
		Nombre:             r.Nombre,
		Moneda:             r.Moneda,
		FechaActualizacion: fechaActualizacion,
	}, nil
}

// IsValid valida que la tasa de cambio tenga datos válidos
func (e *ExchangeRate) IsValid() bool {
	return e.Compra > 0 && e.Venta > 0 && e.Casa != "" && e.Nombre != "" && e.Moneda != ""
}

// GetSpread calcula la diferencia entre compra y venta
func (e *ExchangeRate) GetSpread() float64 {
	return e.Venta - e.Compra
}

// GetSpreadPercentage calcula el porcentaje de spread
func (e *ExchangeRate) GetSpreadPercentage() float64 {
	if e.Compra == 0 {
		return 0
	}
	return (e.GetSpread() / e.Compra) * 100
}
