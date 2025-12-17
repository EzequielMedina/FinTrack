package pdf

import (
	"bytes"
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

// Generator maneja la generación de PDFs
type Generator struct {
	pdf *gofpdf.Fpdf
}

// NewGenerator crea un nuevo generador de PDF con soporte UTF-8
func NewGenerator() *Generator {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)
	pdf.SetAutoPageBreak(true, 25)
	
	return &Generator{pdf: pdf}
}

// encodeToLatin1 convierte una cadena UTF-8 a ISO-8859-1 (Latin-1)
// para que gofpdf pueda mostrar correctamente los acentos en español
func encodeToLatin1(s string) string {
	encoder := charmap.ISO8859_1.NewEncoder()
	encodedBytes, _, err := transform.Bytes(encoder, []byte(s))
	if err != nil {
		// Si hay error en la conversión, devolver el string original
		// Esto puede pasar con caracteres que no están en Latin-1
		return s
	}
	return string(encodedBytes)
}

// AddHeader agrega el encabezado del documento
func (g *Generator) AddHeader(title, subtitle string) {
	g.pdf.AddPage()

	// Logo/Título principal
	g.pdf.SetFont("Arial", "B", 24)
	g.pdf.SetTextColor(41, 128, 185) // Azul
	g.pdf.CellFormat(0, 15, encodeToLatin1("FinTrack"), "", 1, "L", false, 0, "")

	// Título del reporte
	g.pdf.SetFont("Arial", "B", 18)
	g.pdf.SetTextColor(52, 73, 94) // Gris oscuro
	g.pdf.CellFormat(0, 10, encodeToLatin1(title), "", 1, "L", false, 0, "")

	// Subtítulo
	if subtitle != "" {
		g.pdf.SetFont("Arial", "", 11)
		g.pdf.SetTextColor(127, 140, 141) // Gris claro
		g.pdf.CellFormat(0, 6, encodeToLatin1(subtitle), "", 1, "L", false, 0, "")
	}

	// Fecha de generación
	g.pdf.SetFont("Arial", "I", 9)
	g.pdf.SetTextColor(149, 165, 166)
	g.pdf.CellFormat(0, 5, encodeToLatin1(fmt.Sprintf("Generado: %s", time.Now().Format("02/01/2006 15:04"))), "", 1, "L", false, 0, "")

	// Línea separadora
	g.pdf.Ln(5)
	g.pdf.SetDrawColor(189, 195, 199)
	g.pdf.SetLineWidth(0.5)
	g.pdf.Line(20, g.pdf.GetY(), 190, g.pdf.GetY())
	g.pdf.Ln(8)
}

// AddSection agrega una sección con título
func (g *Generator) AddSection(title string) {
	g.pdf.SetFont("Arial", "B", 14)
	g.pdf.SetTextColor(44, 62, 80)
	g.pdf.CellFormat(0, 8, encodeToLatin1(title), "", 1, "L", false, 0, "")
	g.pdf.Ln(3)
}

// AddKeyValue agrega un par clave-valor
func (g *Generator) AddKeyValue(key, value string) {
	g.pdf.SetFont("Arial", "B", 10)
	g.pdf.SetTextColor(52, 73, 94)
	g.pdf.CellFormat(60, 6, encodeToLatin1(key+":"), "", 0, "L", false, 0, "")

	g.pdf.SetFont("Arial", "", 10)
	g.pdf.SetTextColor(52, 73, 94)
	g.pdf.CellFormat(0, 6, encodeToLatin1(value), "", 1, "L", false, 0, "")
}

// AddTable agrega una tabla con encabezados y datos
func (g *Generator) AddTable(headers []string, widths []float64, data [][]string) {
	// Encabezados
	g.pdf.SetFont("Arial", "B", 9)
	g.pdf.SetFillColor(52, 152, 219)  // Azul
	g.pdf.SetTextColor(255, 255, 255) // Blanco

	for i, header := range headers {
		g.pdf.CellFormat(widths[i], 8, encodeToLatin1(header), "1", 0, "C", true, 0, "")
	}
	g.pdf.Ln(-1)

	// Datos
	g.pdf.SetFont("Arial", "", 9)
	g.pdf.SetTextColor(52, 73, 94)
	fill := false

	for _, row := range data {
		if fill {
			g.pdf.SetFillColor(236, 240, 241) // Gris muy claro
		} else {
			g.pdf.SetFillColor(255, 255, 255) // Blanco
		}

		for i, cell := range row {
			g.pdf.CellFormat(widths[i], 7, encodeToLatin1(cell), "1", 0, "L", true, 0, "")
		}
		g.pdf.Ln(-1)
		fill = !fill
	}

	g.pdf.Ln(5)
}

// AddSummaryBox agrega un cuadro resumen con métricas
func (g *Generator) AddSummaryBox(items map[string]string) {
	g.pdf.SetFillColor(236, 240, 241)
	g.pdf.Rect(20, g.pdf.GetY(), 170, float64(len(items)*8+8), "F")

	currentY := g.pdf.GetY() + 5

	for key, value := range items {
		g.pdf.SetXY(25, currentY)
		g.pdf.SetFont("Arial", "B", 10)
		g.pdf.SetTextColor(52, 73, 94)
		g.pdf.CellFormat(80, 6, encodeToLatin1(key+":"), "", 0, "L", false, 0, "")

		g.pdf.SetFont("Arial", "", 10)
		g.pdf.SetTextColor(41, 128, 185)
		g.pdf.CellFormat(0, 6, encodeToLatin1(value), "", 1, "L", false, 0, "")

		currentY += 8
	}

	g.pdf.SetY(currentY + 3)
	g.pdf.Ln(5)
}

// AddFooter agrega pie de página a todas las páginas
func (g *Generator) AddFooter() {
	g.pdf.SetFooterFunc(func() {
		g.pdf.SetY(-15)
		g.pdf.SetFont("Arial", "I", 8)
		g.pdf.SetTextColor(127, 140, 141)
		g.pdf.CellFormat(0, 10, encodeToLatin1(fmt.Sprintf("Página %d", g.pdf.PageNo())), "", 0, "C", false, 0, "")
	})
}

// Output genera el PDF y devuelve los bytes
func (g *Generator) Output() ([]byte, error) {
	var buf bytes.Buffer

	err := g.pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("error generando PDF: %w", err)
	}

	return buf.Bytes(), nil
}

// FormatCurrency formatea un número como moneda
func FormatCurrency(amount float64, currency string) string {
	if currency == "USD" {
		return fmt.Sprintf("$%.2f USD", amount)
	}
	return fmt.Sprintf("$%.2f ARS", amount)
}

// FormatDate formatea una fecha
func FormatDate(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.Format("02/01/2006")
}

// FormatDateTime formatea fecha y hora
func FormatDateTime(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.Format("02/01/2006 15:04")
}
