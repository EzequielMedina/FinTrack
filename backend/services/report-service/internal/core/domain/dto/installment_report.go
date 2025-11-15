package dto

import "time"

// InstallmentReportRequest request para reporte de cuotas
type InstallmentReportRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Status string `json:"status"` // active, completed, overdue
}

// InstallmentReportResponse respuesta del reporte de cuotas
type InstallmentReportResponse struct {
	UserID   string             `json:"user_id"`
	Summary  InstallmentSummary `json:"summary"`
	Plans    []InstallmentPlan  `json:"plans"`
	Upcoming []UpcomingPayment  `json:"upcoming_payments"`
	Overdue  []OverduePayment   `json:"overdue_payments"`
}

// InstallmentSummary resumen de cuotas
type InstallmentSummary struct {
	TotalPlans           int        `json:"total_plans"`
	ActivePlans          int        `json:"active_plans"`
	TotalAmount          float64    `json:"total_amount"`
	PaidAmount           float64    `json:"paid_amount"`
	RemainingAmount      float64    `json:"remaining_amount"`
	OverdueAmount        float64    `json:"overdue_amount"`
	NextPaymentAmount    float64    `json:"next_payment_amount"`
	NextPaymentDate      *time.Time `json:"next_payment_date,omitempty"`
	CompletionPercentage float64    `json:"completion_percentage"`
}

// InstallmentPlan plan de cuotas
type InstallmentPlan struct {
	ID                   string     `json:"id"`
	CardID               string     `json:"card_id"`
	CardLastFour         string     `json:"card_last_four"`
	TotalAmount          float64    `json:"total_amount"`
	InstallmentsCount    int        `json:"installments_count"`
	InstallmentAmount    float64    `json:"installment_amount"`
	PaidInstallments     int        `json:"paid_installments"`
	RemainingAmount      float64    `json:"remaining_amount"`
	Status               string     `json:"status"`
	Description          string     `json:"description,omitempty"`
	MerchantName         string     `json:"merchant_name,omitempty"`
	StartDate            time.Time  `json:"start_date"`
	NextDueDate          *time.Time `json:"next_due_date,omitempty"`
	CompletionPercentage float64    `json:"completion_percentage"`
}

// UpcomingPayment pago próximo
type UpcomingPayment struct {
	InstallmentID string    `json:"installment_id"`
	PlanID        string    `json:"plan_id"`
	CardLastFour  string    `json:"card_last_four"`
	Amount        float64   `json:"amount"`
	DueDate       time.Time `json:"due_date"`
	DaysUntilDue  int       `json:"days_until_due"`
	Description   string    `json:"description,omitempty"`
	MerchantName  string    `json:"merchant_name,omitempty"`
}

// OverduePayment pago vencido
type OverduePayment struct {
	InstallmentID string    `json:"installment_id"`
	PlanID        string    `json:"plan_id"`
	CardLastFour  string    `json:"card_last_four"`
	Amount        float64   `json:"amount"`
	DueDate       time.Time `json:"due_date"`
	DaysOverdue   int       `json:"days_overdue"`
	LateFee       float64   `json:"late_fee"`
	Description   string    `json:"description,omitempty"`
	MerchantName  string    `json:"merchant_name,omitempty"`
}
