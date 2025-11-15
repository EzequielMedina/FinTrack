package pdf

import (
    "bytes"
    "context"
    "fmt"

    "github.com/fintrack/chatbot-service/internal/core/ports"
)

type Generator struct{}

func NewGenerator() *Generator { return &Generator{} }

// Generate: stub simple que devuelve contenido texto (placeholder)
func (g *Generator) Generate(ctx context.Context, data ports.ReportData) ([]byte, error) {
    // TODO: Integrar gofpdf o wkhtmltopdf según config
    buf := bytes.NewBuffer(nil)
    fmt.Fprintf(buf, "Reporte: %s\nPeriodo: %s..%s\nGastos: %.2f\nIngresos: %.2f\n",
        data.Title,
        data.Period.From.Format("2006-01-02"),
        data.Period.To.Format("2006-01-02"),
        data.Totals.Expenses,
        data.Totals.Incomes,
    )
    return buf.Bytes(), nil
}

var _ ports.ReportProvider = (*Generator)(nil)