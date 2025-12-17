package pdf

import (
	"fmt"

	"github.com/fintrack/report-service/internal/core/domain/dto"
)

// GenerateExpenseIncomeReportPDF genera el PDF del reporte de gastos vs ingresos
func GenerateExpenseIncomeReportPDF(report *dto.ExpenseIncomeReportResponse) ([]byte, error) {
	gen := NewGenerator()
	gen.AddFooter()

	// Encabezado
	subtitle := fmt.Sprintf("Del %s al %s",
		FormatDate(report.Period.StartDate),
		FormatDate(report.Period.EndDate))
	gen.AddHeader("Reporte de Gastos vs Ingresos", subtitle)

	// Resumen
	gen.AddSection("Resumen General")
	summaryData := map[string]string{
		"Total Ingresos":  FormatCurrency(report.Summary.TotalIncome, "ARS"),
		"Total Egresos":   FormatCurrency(report.Summary.TotalExpenses, "ARS"),
		"Balance Neto":    FormatCurrency(report.Summary.NetBalance, "ARS"),
		"Tasa de Ahorro":  fmt.Sprintf("%.1f%%", report.Summary.SavingsRate),
		"Ratio de Gastos": fmt.Sprintf("%.1f%%", report.Summary.ExpenseRatio),
		"Ingreso Diario":  FormatCurrency(report.Summary.AvgDailyIncome, "ARS"),
		"Gasto Diario":    FormatCurrency(report.Summary.AvgDailyExpense, "ARS"),
	}
	gen.AddSummaryBox(summaryData)

	// Por Categoría - Ingresos
	var incomeCategories []dto.ExpenseIncomeByCategory
	var expenseCategories []dto.ExpenseIncomeByCategory
	for _, cat := range report.ByCategory {
		if cat.Type == "income" {
			incomeCategories = append(incomeCategories, cat)
		} else {
			expenseCategories = append(expenseCategories, cat)
		}
	}

	if len(incomeCategories) > 0 {
		gen.AddSection("Ingresos por Categoría")

		headers := []string{"Categoría", "Cantidad", "Monto", "Porcentaje"}
		widths := []float64{60, 30, 40, 40}

		var tableData [][]string
		for _, cat := range incomeCategories {
			row := []string{
				translateCategory(cat.Category),
				fmt.Sprintf("%d", cat.Count),
				FormatCurrency(cat.Amount, "ARS"),
				fmt.Sprintf("%.1f%%", cat.Percentage),
			}
			tableData = append(tableData, row)
		}

		gen.AddTable(headers, widths, tableData)
	}

	// Por Categoría - Egresos
	if len(expenseCategories) > 0 {
		gen.AddSection("Egresos por Categoría")

		headers := []string{"Categoría", "Cantidad", "Monto", "Porcentaje"}
		widths := []float64{60, 30, 40, 40}

		var tableData [][]string
		for _, cat := range expenseCategories {
			row := []string{
				translateCategory(cat.Category),
				fmt.Sprintf("%d", cat.Count),
				FormatCurrency(cat.Amount, "ARS"),
				fmt.Sprintf("%.1f%%", cat.Percentage),
			}
			tableData = append(tableData, row)
		}

		gen.AddTable(headers, widths, tableData)
	}

	// Por Período
	if len(report.ByPeriod) > 0 {
		gen.AddSection("Evolución por Período")

		headers := []string{"Período", "Ingresos", "Egresos", "Balance", "Ahorro %"}
		widths := []float64{30, 35, 35, 35, 35}

		var tableData [][]string
		for _, period := range report.ByPeriod {
			row := []string{
				period.Period,
				FormatCurrency(period.Income, "ARS"),
				FormatCurrency(period.Expenses, "ARS"),
				FormatCurrency(period.Net, "ARS"),
				fmt.Sprintf("%.1f%%", period.SavingsRate),
			}
			tableData = append(tableData, row)
		}

		gen.AddTable(headers, widths, tableData)
	}

	// Análisis de Tendencias
	gen.AddSection("Análisis de Tendencias")
	trendData := map[string]string{
		"Tendencia Ingresos": translateTrend(report.Trend.IncomesTrend) +
			fmt.Sprintf(" (%.1f%%)", report.Trend.IncomeChange),
		"Tendencia Egresos": translateTrend(report.Trend.ExpensesTrend) +
			fmt.Sprintf(" (%.1f%%)", report.Trend.ExpenseChange),
		"Tendencia Neta": translateTrend(report.Trend.NetTrend),
	}

	// Agregar pronóstico si existe
	if report.Trend.Forecast != nil {
		trendData["Pronóstico Ingresos"] = FormatCurrency(report.Trend.Forecast.NextMonthIncome, "ARS")
		trendData["Pronóstico Egresos"] = FormatCurrency(report.Trend.Forecast.NextMonthExpenses, "ARS")
		trendData["Pronóstico Balance"] = FormatCurrency(report.Trend.Forecast.NextMonthNet, "ARS")
	}

	gen.AddSummaryBox(trendData)

	return gen.Output()
}

func translateTrend(trend string) string {
	switch trend {
	case "increasing":
		return "↑ Aumentando"
	case "decreasing":
		return "↓ Disminuyendo"
	case "stable":
		return "→ Estable"
	case "improving":
		return "✓ Mejorando"
	case "declining":
		return "✗ Decayendo"
	default:
		return trend
	}
}

func translateCategory(category string) string {
	// Traducción de categorías comunes
	categoryMap := map[string]string{
		// Ingresos
		"salary":       "Salario",
		"investment":   "Inversión",
		"business":     "Negocio",
		"other_income": "Otros Ingresos",
		// Gastos
		"food":         "Alimentos",
		"transport":    "Transporte",
		"utilities":    "Servicios",
		"entertainment": "Entretenimiento",
		"healthcare":   "Salud",
		"education":    "Educación",
		"shopping":     "Compras",
		"other":        "Otros",
		// Si la categoría ya está en español o no está en el mapa, devolverla tal cual
	}
	
	if translated, exists := categoryMap[category]; exists {
		return translated
	}
	// Si no está en el mapa, devolver la categoría original (puede estar ya en español)
	return category
}
