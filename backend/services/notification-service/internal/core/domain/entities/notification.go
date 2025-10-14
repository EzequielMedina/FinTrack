package entities

import "time"

// Card representa una tarjeta para notificaciones
type Card struct {
	ID       string    `json:"id" db:"id"`
	UserID   string    `json:"user_id" db:"user_id"`
	CardName string    `json:"card_name" db:"card_name"`
	BankName string    `json:"bank_name" db:"bank_name"`
	LastFour string    `json:"last_four_digits" db:"last_four_digits"`
	DueDate  time.Time `json:"due_date" db:"due_date"`
	IsActive bool      `json:"is_active" db:"is_active"`
	CardType string    `json:"card_type" db:"card_type"`
	// Campos del usuario (obtenidos via JOIN)
	UserEmail     string `json:"user_email" db:"user_email"`
	UserFirstName string `json:"user_first_name" db:"user_first_name"`
	UserLastName  string `json:"user_last_name" db:"user_last_name"`
}

// GetUserFullName retorna el nombre completo del usuario
func (c *Card) GetUserFullName() string {
	return c.UserFirstName + " " + c.UserLastName
}

// Installment representa una cuota pendiente
type Installment struct {
	ID                string    `json:"id" db:"id"`
	PlanID            string    `json:"plan_id" db:"plan_id"`
	Amount            float64   `json:"amount" db:"amount"`
	DueDate           time.Time `json:"due_date" db:"due_date"`
	Status            string    `json:"status" db:"status"`
	Description       string    `json:"description" db:"description"`
	MerchantName      string    `json:"merchant_name" db:"merchant_name"`
	InstallmentNum    int       `json:"installment_num" db:"installment_num"`
	TotalInstallments int       `json:"total_installments" db:"total_installments"`
}

// InstallmentSummary para detalles de cuotas en el email
type InstallmentSummary struct {
	Description       string    `json:"description"`
	MerchantName      string    `json:"merchant_name"`
	Amount            float64   `json:"amount"`
	DueDate           time.Time `json:"due_date"`
	InstallmentNum    int       `json:"installment_num"`
	TotalInstallments int       `json:"total_installments"`
}

// GetInstallmentText retorna el texto descriptivo de la cuota
func (i *InstallmentSummary) GetInstallmentText() string {
	if i.TotalInstallments > 1 {
		return i.Description + " (" + i.MerchantName + ")"
	}
	return i.Description + " - " + i.MerchantName
}

// CardDueNotification contiene datos para el email
type CardDueNotification struct {
	CardID              string               `json:"card_id"`
	UserID              string               `json:"user_id"`
	UserEmail           string               `json:"user_email"`
	UserName            string               `json:"user_name"`
	CardName            string               `json:"card_name"`
	BankName            string               `json:"bank_name"`
	LastFour            string               `json:"last_four"`
	DueDate             time.Time            `json:"due_date"`
	TotalPendingAmount  float64              `json:"total_pending_amount"`
	PendingInstallments int                  `json:"pending_installments"`
	InstallmentDetails  []InstallmentSummary `json:"installment_details"`
}

// NotificationLog para auditoría
type NotificationLog struct {
	ID           string    `json:"id" db:"id"`
	JobRunID     string    `json:"job_run_id" db:"job_run_id"`
	CardID       string    `json:"card_id" db:"card_id"`
	UserID       string    `json:"user_id" db:"user_id"`
	Email        string    `json:"email" db:"email"`
	Status       string    `json:"status" db:"status"` // sent, failed, skipped
	ErrorMessage string    `json:"error_message,omitempty" db:"error_message"`
	SentAt       time.Time `json:"sent_at" db:"sent_at"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// JobRun para tracking de ejecuciones del job
type JobRun struct {
	ID           string     `json:"id" db:"id"`
	StartedAt    time.Time  `json:"started_at" db:"started_at"`
	CompletedAt  *time.Time `json:"completed_at" db:"completed_at"`
	Status       string     `json:"status" db:"status"` // running, completed, failed
	CardsFound   int        `json:"cards_found" db:"cards_found"`
	EmailsSent   int        `json:"emails_sent" db:"emails_sent"`
	Errors       int        `json:"errors" db:"errors"`
	ErrorMessage string     `json:"error_message,omitempty" db:"error_message"`
}
