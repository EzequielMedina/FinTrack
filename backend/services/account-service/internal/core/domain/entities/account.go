package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AccountType represents the type of account
type AccountType string

const (
	AccountTypeChecking AccountType = "checking"
	AccountTypeSavings  AccountType = "savings"
	AccountTypeCredit   AccountType = "credit"
	AccountTypeDebit    AccountType = "debit"
	AccountTypeWallet   AccountType = "wallet"
)

// Currency represents the currency type
type Currency string

const (
	CurrencyUSD Currency = "USD"
	CurrencyARS Currency = "ARS"
	CurrencyEUR Currency = "EUR"
)

// Account represents a financial account in the system
type Account struct {
	ID          string      `gorm:"type:char(36);primaryKey" json:"id"`
	UserID      string      `gorm:"type:char(36);not null;index" json:"user_id"`
	AccountType AccountType `gorm:"type:varchar(20);not null;index" json:"account_type"`
	Name        string      `gorm:"type:varchar(100);not null" json:"name"`
	Description string      `gorm:"type:text" json:"description"`
	Currency    Currency    `gorm:"type:varchar(3);not null;index" json:"currency"`
	Balance     float64     `gorm:"type:decimal(15,2);not null;default:0" json:"balance"`

	// Credit card specific fields
	CreditLimit *float64   `gorm:"type:decimal(15,2);null" json:"credit_limit,omitempty"`
	ClosingDate *time.Time `gorm:"type:date;null" json:"closing_date,omitempty"`
	DueDate     *time.Time `gorm:"type:date;null" json:"due_date,omitempty"`

	// Personal identification (for virtual wallets)
	DNI *string `gorm:"type:varchar(20);null" json:"dni,omitempty"`

	IsActive  bool           `gorm:"not null;default:true;index" json:"is_active"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName returns the table name for the Account model
func (Account) TableName() string {
	return "accounts"
}

// BeforeCreate is called before creating a new account
func (a *Account) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}

// IsValidAccountType checks if the account type is valid
func IsValidAccountType(accountType AccountType) bool {
	switch accountType {
	case AccountTypeChecking, AccountTypeSavings, AccountTypeCredit, AccountTypeDebit, AccountTypeWallet:
		return true
	default:
		return false
	}
}

// IsValidCurrency checks if the currency is valid
func IsValidCurrency(currency Currency) bool {
	switch currency {
	case CurrencyUSD, CurrencyARS, CurrencyEUR:
		return true
	default:
		return false
	}
}

// Validate validates the account data
func (a *Account) Validate() error {
	if a.UserID == "" {
		return &ValidationError{Field: "user_id", Message: "user ID is required"}
	}
	if a.Name == "" {
		return &ValidationError{Field: "name", Message: "account name is required"}
	}
	if !IsValidAccountType(a.AccountType) {
		return &ValidationError{Field: "account_type", Message: "invalid account type"}
	}
	if !IsValidCurrency(a.Currency) {
		return &ValidationError{Field: "currency", Message: "invalid currency"}
	}
	return nil
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return e.Message
}
