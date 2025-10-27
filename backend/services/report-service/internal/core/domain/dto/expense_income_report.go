package dto

import "time"

// ExpenseIncomeReportRequest request para reporte de gastos vs ingresos
type ExpenseIncomeReportRequest struct {
	UserID    string    `json:"user_id" binding:"required"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
	GroupBy   string    `json:"group_by"` // day, week, month
}

// ExpenseIncomeReportResponse respuesta del reporte de gastos vs ingresos
type ExpenseIncomeReportResponse struct {
	UserID     string                    `json:"user_id"`
	Period     Period                    `json:"period"`
	Summary    ExpenseIncomeSummary      `json:"summary"`
	ByPeriod   []ExpenseIncomeByPeriod   `json:"by_period"`
	ByCategory []ExpenseIncomeByCategory `json:"by_category"`
	Trend      TrendAnalysis             `json:"trend"`
}

// ExpenseIncomeSummary resumen de gastos vs ingresos
type ExpenseIncomeSummary struct {
	TotalIncome     float64 `json:"total_income"`
	TotalExpenses   float64 `json:"total_expenses"`
	NetBalance      float64 `json:"net_balance"`
	SavingsRate     float64 `json:"savings_rate"`
	ExpenseRatio    float64 `json:"expense_ratio"`
	AvgDailyIncome  float64 `json:"avg_daily_income"`
	AvgDailyExpense float64 `json:"avg_daily_expense"`
}

// ExpenseIncomeByPeriod gastos e ingresos por período
type ExpenseIncomeByPeriod struct {
	Period      string    `json:"period"`
	Date        time.Time `json:"date"`
	Income      float64   `json:"income"`
	Expenses    float64   `json:"expenses"`
	Net         float64   `json:"net"`
	SavingsRate float64   `json:"savings_rate"`
}

// ExpenseIncomeByCategory gastos e ingresos por categoría
type ExpenseIncomeByCategory struct {
	Category   string  `json:"category"`
	Type       string  `json:"type"` // income, expense
	Amount     float64 `json:"amount"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

// TrendAnalysis análisis de tendencias
type TrendAnalysis struct {
	IncomesTrend  string        `json:"incomes_trend"`  // increasing, decreasing, stable
	ExpensesTrend string        `json:"expenses_trend"` // increasing, decreasing, stable
	NetTrend      string        `json:"net_trend"`      // improving, declining, stable
	IncomeChange  float64       `json:"income_change"`  // % cambio
	ExpenseChange float64       `json:"expense_change"` // % cambio
	Forecast      *ForecastData `json:"forecast,omitempty"`
}

// ForecastData proyección futura
type ForecastData struct {
	NextMonthIncome   float64 `json:"next_month_income"`
	NextMonthExpenses float64 `json:"next_month_expenses"`
	NextMonthNet      float64 `json:"next_month_net"`
}
