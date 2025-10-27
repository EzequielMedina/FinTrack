package dto

import "time"

// TransactionReportRequest request para reporte de transacciones
type TransactionReportRequest struct {
	UserID    string    `json:"user_id" binding:"required"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Type      string    `json:"type"`
	GroupBy   string    `json:"group_by"` // day, week, month
}

// TransactionReportResponse respuesta del reporte de transacciones
type TransactionReportResponse struct {
	UserID      string                `json:"user_id"`
	Period      Period                `json:"period"`
	Summary     TransactionSummary    `json:"summary"`
	ByType      []TransactionByType   `json:"by_type"`
	ByPeriod    []TransactionByPeriod `json:"by_period"`
	TopExpenses []TransactionItem     `json:"top_expenses"`
}

// TransactionSummary resumen de transacciones
type TransactionSummary struct {
	TotalTransactions int     `json:"total_transactions"`
	TotalIncome       float64 `json:"total_income"`
	TotalExpenses     float64 `json:"total_expenses"`
	NetBalance        float64 `json:"net_balance"`
	AvgTransaction    float64 `json:"avg_transaction"`
}

// TransactionByType transacciones agrupadas por tipo
type TransactionByType struct {
	Type       string  `json:"type"`
	Count      int     `json:"count"`
	Amount     float64 `json:"amount"`
	Percentage float64 `json:"percentage"`
}

// TransactionByPeriod transacciones agrupadas por período
type TransactionByPeriod struct {
	Period   string    `json:"period"`
	Date     time.Time `json:"date"`
	Income   float64   `json:"income"`
	Expenses float64   `json:"expenses"`
	Net      float64   `json:"net"`
	Count    int       `json:"count"`
}

// TransactionItem item individual de transacción
type TransactionItem struct {
	ID           string    `json:"id"`
	Description  string    `json:"description"`
	Amount       float64   `json:"amount"`
	Type         string    `json:"type"`
	Date         time.Time `json:"date"`
	MerchantName string    `json:"merchant_name,omitempty"`
}

// Period período de tiempo
type Period struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Days      int       `json:"days"`
}
