package pdf

import (
	"fmt"

	"github.com/fintrack/report-service/internal/core/domain/dto"
)

// GenerateTransactionReportPDF genera el PDF del reporte de transacciones
func GenerateTransactionReportPDF(report *dto.TransactionReportResponse) ([]byte, error) {
	gen := NewGenerator()
	gen.AddFooter()

	// Encabezado
	subtitle := fmt.Sprintf("Del %s al %s",
		FormatDate(report.Period.StartDate),
		FormatDate(report.Period.EndDate))
	gen.AddHeader("Reporte de Transacciones", subtitle)

	// Resumen
	gen.AddSection("Resumen")
	summaryData := map[string]string{
		"Total Transacciones": fmt.Sprintf("%d", report.Summary.TotalTransactions),
		"Ingresos":            FormatCurrency(report.Summary.TotalIncome, "ARS"),
		"Egresos":             FormatCurrency(report.Summary.TotalExpenses, "ARS"),
		"Balance Neto":        FormatCurrency(report.Summary.NetBalance, "ARS"),
		"Promedio":            FormatCurrency(report.Summary.AvgTransaction, "ARS"),
	}
	gen.AddSummaryBox(summaryData)

	// Por Tipo
	if len(report.ByType) > 0 {
		gen.AddSection("Por Tipo de Transacción")

		headers := []string{"Tipo", "Cantidad", "Monto", "Porcentaje"}
		widths := []float64{50, 35, 40, 45}

		var tableData [][]string
		for _, item := range report.ByType {
			row := []string{
				translateTransactionType(item.Type),
				fmt.Sprintf("%d", item.Count),
				FormatCurrency(item.Amount, "ARS"),
				fmt.Sprintf("%.1f%%", item.Percentage),
			}
			tableData = append(tableData, row)
		}

		gen.AddTable(headers, widths, tableData)
	}

	// Principales Gastos
	if len(report.TopExpenses) > 0 {
		gen.AddSection("Principales Gastos")

		headers := []string{"Fecha", "Descripción", "Comercio", "Monto"}
		widths := []float64{25, 60, 50, 35}

		var tableData [][]string
		for _, tx := range report.TopExpenses {
			description := tx.Description
			if len(description) > 50 {
				description = description[:47] + "..."
			}

			merchant := tx.MerchantName
			if merchant == "" {
				merchant = "-"
			} else if len(merchant) > 40 {
				merchant = merchant[:37] + "..."
			}

			row := []string{
				FormatDate(tx.Date),
				description,
				merchant,
				FormatCurrency(tx.Amount, "ARS"),
			}
			tableData = append(tableData, row)
		}

		gen.AddTable(headers, widths, tableData)
	}

	// Por Período
	if len(report.ByPeriod) > 0 {
		gen.AddSection("Por Período")

		headers := []string{"Período", "Ingresos", "Egresos", "Balance", "Transacciones"}
		widths := []float64{30, 35, 35, 35, 35}

		var tableData [][]string
		for _, period := range report.ByPeriod {
			row := []string{
				period.Period,
				FormatCurrency(period.Income, "ARS"),
				FormatCurrency(period.Expenses, "ARS"),
				FormatCurrency(period.Net, "ARS"),
				fmt.Sprintf("%d", period.Count),
			}
			tableData = append(tableData, row)
		}

		gen.AddTable(headers, widths, tableData)
	}

	return gen.Output()
}

func translateTransactionType(txType string) string {
	switch txType {
	// Tipos genéricos
	case "income":
		return "Ingreso"
	case "expense":
		return "Egreso"
	case "transfer":
		return "Transferencia"
	// Tipos específicos de billetera
	case "wallet_deposit":
		return "Depósito en Billetera"
	case "wallet_withdrawal":
		return "Retiro de Billetera"
	case "wallet_transfer":
		return "Transferencia de Billetera"
	// Tipos de tarjeta de crédito
	case "credit_charge":
		return "Cargo en Crédito"
	case "credit_payment":
		return "Pago de Crédito"
	case "credit_refund":
		return "Reembolso de Crédito"
	// Tipos de tarjeta de débito
	case "debit_purchase":
		return "Compra con Débito"
	case "debit_refund":
		return "Reembolso de Débito"
	// Tipos de cuenta
	case "account_deposit":
		return "Depósito en Cuenta"
	case "account_withdraw":
		return "Retiro de Cuenta"
	default:
		// Si no se encuentra traducción, devolver el tipo original
		return txType
	}
}
