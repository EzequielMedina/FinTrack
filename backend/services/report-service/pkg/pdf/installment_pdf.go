package pdf

import (
	"fmt"

	"github.com/fintrack/report-service/internal/core/domain/dto"
)

// GenerateInstallmentReportPDF genera el PDF del reporte de cuotas
func GenerateInstallmentReportPDF(report *dto.InstallmentReportResponse) ([]byte, error) {
	gen := NewGenerator()
	gen.AddFooter()

	// Encabezado
	subtitle := fmt.Sprintf("Resumen de Planes de Cuotas y Pagos")
	gen.AddHeader("Reporte de Cuotas", subtitle)

	// Resumen
	gen.AddSection("Resumen General")
	summaryData := map[string]string{
		"Total Planes":    fmt.Sprintf("%d", report.Summary.TotalPlans),
		"Planes Activos":  fmt.Sprintf("%d", report.Summary.ActivePlans),
		"Monto Total":     FormatCurrency(report.Summary.TotalAmount, "ARS"),
		"Total Pagado":    FormatCurrency(report.Summary.PaidAmount, "ARS"),
		"Saldo Pendiente": FormatCurrency(report.Summary.RemainingAmount, "ARS"),
		"Monto Vencido":   FormatCurrency(report.Summary.OverdueAmount, "ARS"),
		"Próximo Pago":    FormatCurrency(report.Summary.NextPaymentAmount, "ARS"),
	}
	gen.AddSummaryBox(summaryData)

	// Planes de cuotas
	if len(report.Plans) > 0 {
		gen.AddSection("Planes de Cuotas")

		headers := []string{"Descripción", "Tarjeta", "Total", "Cuotas", "Pagadas", "Pendiente", "Estado"}
		widths := []float64{50, 25, 25, 15, 15, 25, 15}

		var tableData [][]string
		for _, plan := range report.Plans {
			cardInfo := plan.CardLastFour
			if cardInfo == "" {
				cardInfo = "-"
			} else {
				cardInfo = "*" + cardInfo
			}

			description := plan.Description
			if description == "" {
				description = plan.MerchantName
			}
			if len(description) > 30 {
				description = description[:27] + "..."
			}

			row := []string{
				description,
				cardInfo,
				FormatCurrency(plan.TotalAmount, "ARS"),
				fmt.Sprintf("%d", plan.InstallmentsCount),
				fmt.Sprintf("%d", plan.PaidInstallments),
				FormatCurrency(plan.RemainingAmount, "ARS"),
				translateStatus(plan.Status),
			}
			tableData = append(tableData, row)
		}

		gen.AddTable(headers, widths, tableData)
	}

	// Próximos pagos
	if len(report.Upcoming) > 0 {
		gen.AddSection("Próximos Pagos (30 días)")

		headers := []string{"Fecha", "Descripción", "Tarjeta", "Monto", "Días"}
		widths := []float64{25, 60, 25, 30, 20}

		var tableData [][]string
		for _, payment := range report.Upcoming {
			cardInfo := payment.CardLastFour
			if cardInfo == "" {
				cardInfo = "-"
			} else {
				cardInfo = "*" + cardInfo
			}

			description := payment.Description
			if description == "" {
				description = payment.MerchantName
			}
			if len(description) > 35 {
				description = description[:32] + "..."
			}

			row := []string{
				FormatDate(payment.DueDate),
				description,
				cardInfo,
				FormatCurrency(payment.Amount, "ARS"),
				fmt.Sprintf("%d", payment.DaysUntilDue),
			}
			tableData = append(tableData, row)
		}

		gen.AddTable(headers, widths, tableData)
	}

	// Pagos vencidos
	if len(report.Overdue) > 0 {
		gen.AddSection("Pagos Vencidos")

		headers := []string{"Fecha", "Descripción", "Tarjeta", "Monto", "Mora", "Días"}
		widths := []float64{25, 50, 25, 25, 20, 15}

		var tableData [][]string
		for _, payment := range report.Overdue {
			cardInfo := payment.CardLastFour
			if cardInfo == "" {
				cardInfo = "-"
			} else {
				cardInfo = "*" + cardInfo
			}

			description := payment.Description
			if description == "" {
				description = payment.MerchantName
			}
			if len(description) > 30 {
				description = description[:27] + "..."
			}

			row := []string{
				FormatDate(payment.DueDate),
				description,
				cardInfo,
				FormatCurrency(payment.Amount, "ARS"),
				FormatCurrency(payment.LateFee, "ARS"),
				fmt.Sprintf("%d", payment.DaysOverdue),
			}
			tableData = append(tableData, row)
		}

		gen.AddTable(headers, widths, tableData)
	}

	return gen.Output()
}

func translateStatus(status string) string {
	switch status {
	case "active":
		return "Activo"
	case "completed":
		return "Completado"
	case "pending":
		return "Pendiente"
	case "paid":
		return "Pagado"
	case "partial":
		return "Parcial"
	case "cancelled":
		return "Cancelado"
	case "overdue":
		return "Vencido"
	default:
		return status
	}
}
