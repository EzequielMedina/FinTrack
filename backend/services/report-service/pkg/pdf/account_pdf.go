package pdf

import (
	"fmt"

	"github.com/fintrack/report-service/internal/core/domain/dto"
)

// GenerateAccountReportPDF genera el PDF del reporte de cuentas
func GenerateAccountReportPDF(report *dto.AccountReportResponse) ([]byte, error) {
	gen := NewGenerator()
	gen.AddFooter()

	// Encabezado
	subtitle := fmt.Sprintf("Resumen de Cuentas y Tarjetas")
	gen.AddHeader("Reporte de Cuentas", subtitle)

	// Resumen
	gen.AddSection("Resumen General")
	summaryData := map[string]string{
		"Total Cuentas":      fmt.Sprintf("%d", report.Summary.TotalAccounts),
		"Total Tarjetas":     fmt.Sprintf("%d", report.Summary.TotalCards),
		"Saldo Total":        FormatCurrency(report.Summary.TotalBalance, "ARS"),
		"Límite Crédito":     FormatCurrency(report.Summary.TotalCreditLimit, "ARS"),
		"Crédito Usado":      FormatCurrency(report.Summary.TotalCreditUsed, "ARS"),
		"Crédito Disponible": FormatCurrency(report.Summary.AvailableCredit, "ARS"),
		"Utilización":        fmt.Sprintf("%.1f%%", report.Summary.CreditUtilization),
		"Patrimonio Neto":    FormatCurrency(report.Summary.NetWorth, "ARS"),
	}
	gen.AddSummaryBox(summaryData)

	// Cuentas
	if len(report.Accounts) > 0 {
		gen.AddSection("Detalle de Cuentas")

		headers := []string{"Nombre", "Tipo", "Moneda", "Saldo", "Límite", "Estado"}
		widths := []float64{50, 30, 20, 30, 30, 15}

		var tableData [][]string
		for _, account := range report.Accounts {
			accountType := translateAccountType(account.AccountType)

			estado := "Inactiva"
			if account.IsActive {
				estado = "Activa"
			}

			limite := "-"
			if account.CreditLimit > 0 {
				limite = FormatCurrency(account.CreditLimit, account.Currency)
			}

			row := []string{
				account.Name,
				accountType,
				account.Currency,
				FormatCurrency(account.Balance, account.Currency),
				limite,
				estado,
			}
			tableData = append(tableData, row)
		}

		gen.AddTable(headers, widths, tableData)
	}

	// Tarjetas
	if len(report.Cards) > 0 {
		gen.AddSection("Tarjetas")

		headers := []string{"Número", "Tipo", "Marca", "Titular", "Límite", "Estado"}
		widths := []float64{30, 25, 25, 40, 30, 20}

		var tableData [][]string
		for _, card := range report.Cards {
			cardType := translateCardType(card.CardType)

			cardNumber := "*" + card.LastFourDigits
			if card.Nickname != "" {
				cardNumber = card.Nickname + " " + cardNumber
			}

			limite := "-"
			if card.CreditLimit > 0 {
				limite = FormatCurrency(card.CreditLimit, "ARS")
			}

			row := []string{
				cardNumber,
				cardType,
				card.CardBrand,
				card.HolderName,
				limite,
				translateCardStatus(card.Status),
			}
			tableData = append(tableData, row)
		}

		gen.AddTable(headers, widths, tableData)
	}

	// Distribución por tipo
	if len(report.Distribution) > 0 {
		gen.AddSection("Distribución por Tipo de Cuenta")

		headers := []string{"Tipo", "Cantidad", "Saldo Total"}
		widths := []float64{70, 40, 60}

		var tableData [][]string
		for _, dist := range report.Distribution {
			row := []string{
				translateAccountType(dist.AccountType),
				fmt.Sprintf("%d", dist.Count),
				FormatCurrency(dist.TotalBalance, "ARS"),
			}
			tableData = append(tableData, row)
		}

		gen.AddTable(headers, widths, tableData)
	}

	return gen.Output()
}

func translateAccountType(accountType string) string {
	switch accountType {
	case "savings":
		return "Caja de Ahorro"
	case "checking":
		return "Cuenta Corriente"
	case "credit":
		return "Tarjeta de Crédito"
	case "investment":
		return "Inversión"
	default:
		return accountType
	}
}

func translateCardType(cardType string) string {
	switch cardType {
	case "debit":
		return "Débito"
	case "credit":
		return "Crédito"
	case "prepaid":
		return "Prepaga"
	default:
		return cardType
	}
}

func translateCardStatus(status string) string {
	switch status {
	case "active":
		return "Activa"
	case "inactive":
		return "Inactiva"
	case "blocked":
		return "Bloqueada"
	case "expired":
		return "Vencida"
	default:
		return status
	}
}
