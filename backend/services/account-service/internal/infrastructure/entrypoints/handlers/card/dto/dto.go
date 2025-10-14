package dto

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/fintrack/account-service/internal/core/domain/entities"
)

// CustomDate handles date parsing for both "2006-01-02" and full RFC3339 formats
type CustomDate struct {
	*time.Time
}

func (cd *CustomDate) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	// Remove quotes from JSON string
	dateStr := strings.Trim(string(data), `"`)

	// Try parsing date-only format first (YYYY-MM-DD)
	if t, err := time.Parse("2006-01-02", dateStr); err == nil {
		cd.Time = &t
		return nil
	}

	// Try parsing full RFC3339 format
	if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
		cd.Time = &t
		return nil
	}

	// Try parsing without timezone
	if t, err := time.Parse("2006-01-02T15:04:05", dateStr); err == nil {
		cd.Time = &t
		return nil
	}

	return fmt.Errorf("cannot parse date: %s", dateStr)
}

func (cd CustomDate) MarshalJSON() ([]byte, error) {
	if cd.Time == nil {
		return []byte("null"), nil
	}
	return json.Marshal(cd.Time.Format("2006-01-02"))
}

// ToCustomDate converts *time.Time to *CustomDate
func ToCustomDate(t *time.Time) *CustomDate {
	if t == nil {
		return nil
	}
	return &CustomDate{Time: t}
}

// ToTimePointer converts *CustomDate to *time.Time
func (cd *CustomDate) ToTimePointer() *time.Time {
	if cd == nil || cd.Time == nil {
		return nil
	}
	return cd.Time
}

// CreateCardRequest represents the request to create a new card
type CreateCardRequest struct {
	AccountID       string `json:"account_id,omitempty"`         // Set from URL parameter, not required in JSON
	CardType        string `json:"card_type" binding:"required"` // "credit" or "debit"
	CardBrand       string `json:"card_brand" binding:"required"`
	LastFourDigits  string `json:"last_four_digits" binding:"required,len=4"`
	MaskedNumber    string `json:"masked_number" binding:"required"`
	HolderName      string `json:"holder_name" binding:"required,min=2,max=100"`
	ExpirationMonth int    `json:"expiration_month" binding:"required,min=1,max=12"`
	ExpirationYear  int    `json:"expiration_year" binding:"required"`
	Nickname        string `json:"nickname,omitempty" binding:"max=50"`
	IsDefault       bool   `json:"is_default,omitempty"`

	// Credit card specific fields
	CreditLimit *float64    `json:"credit_limit,omitempty" binding:"omitempty,min=0"`
	ClosingDate *CustomDate `json:"closing_date,omitempty"`
	DueDate     *CustomDate `json:"due_date,omitempty"`

	// Security fields (encrypted data)
	EncryptedNumber string `json:"encrypted_number" binding:"required"`
	KeyFingerprint  string `json:"key_fingerprint" binding:"required"`
}

// UpdateCardRequest represents the request to update a card
type UpdateCardRequest struct {
	HolderName      string   `json:"holder_name,omitempty" binding:"omitempty,min=2,max=100"`
	ExpirationMonth int      `json:"expiration_month,omitempty" binding:"omitempty,min=1,max=12"`
	ExpirationYear  int      `json:"expiration_year,omitempty" binding:"omitempty,min=2020"` // Allow reasonable past years for testing
	Nickname        string   `json:"nickname,omitempty" binding:"max=50"`
	IsDefault       *bool    `json:"is_default,omitempty"`
	CreditLimit     *float64 `json:"credit_limit,omitempty" binding:"omitempty,min=0,max=1000000"` // Allow credit limit updates
}

// CardResponse represents the response for card operations
type CardResponse struct {
	ID              string    `json:"id"`
	AccountID       string    `json:"account_id"`
	CardType        string    `json:"card_type"`
	CardBrand       string    `json:"card_brand"`
	LastFourDigits  string    `json:"last_four_digits"`
	MaskedNumber    string    `json:"masked_number"`
	HolderName      string    `json:"holder_name"`
	ExpirationMonth int       `json:"expiration_month"`
	ExpirationYear  int       `json:"expiration_year"`
	Status          string    `json:"status"`
	IsDefault       bool      `json:"is_default"`
	Nickname        string    `json:"nickname,omitempty"`
	Balance         float64   `json:"balance"` // New: Card balance (debt for credit, 0 for debit)
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// Credit card specific fields
	CreditLimit *float64    `json:"credit_limit,omitempty"`
	ClosingDate *CustomDate `json:"closing_date,omitempty"`
	DueDate     *CustomDate `json:"due_date,omitempty"`

	// Installment plans summary (optional, when requested)
	InstallmentPlans *InstallmentPlansSummary `json:"installment_plans,omitempty"`

	// Security fields (encrypted data - never exposed in response)
	// EncryptedNumber and KeyFingerprint are intentionally omitted for security
}

// InstallmentPlansSummary represents a summary of installment plans for a card
type InstallmentPlansSummary struct {
	TotalActivePlans       int        `json:"total_active_plans"`
	TotalOutstandingAmount float64    `json:"total_outstanding_amount"`
	NextPaymentDue         *time.Time `json:"next_payment_due,omitempty"`
	NextPaymentAmount      float64    `json:"next_payment_amount"`
}

// Credit Card Financial Operations DTOs

// CreditCardChargeRequest represents a charge to a credit card
type CreditCardChargeRequest struct {
	Amount      float64 `json:"amount" binding:"required,min=0.01"`
	Description string  `json:"description" binding:"required,min=3,max=255"`
	Reference   string  `json:"reference,omitempty" binding:"max=50"`
}

// CreditCardPaymentRequest represents a payment to a credit card
type CreditCardPaymentRequest struct {
	Amount        float64 `json:"amount" binding:"required,min=0.01"`
	PaymentMethod string  `json:"payment_method" binding:"required,oneof=bank_transfer debit_card cash"`
	Reference     string  `json:"reference,omitempty" binding:"max=50"`
}

// CreditCardChargeWithInstallmentsRequest represents a charge to a credit card with installments
type CreditCardChargeWithInstallmentsRequest struct {
	Amount            float64   `json:"amount" binding:"required,min=0.01"`
	InstallmentsCount int       `json:"installments_count" binding:"required,min=1,max=24"`
	StartDate         time.Time `json:"start_date" binding:"required"`
	Description       string    `json:"description" binding:"required,min=3,max=255"`
	MerchantName      string    `json:"merchant_name,omitempty" binding:"max=100"`
	MerchantID        string    `json:"merchant_id,omitempty" binding:"max=50"`
	InterestRate      float64   `json:"interest_rate,omitempty" binding:"min=0,max=100"`
	AdminFee          float64   `json:"admin_fee,omitempty" binding:"min=0"`
	Reference         string    `json:"reference,omitempty" binding:"max=50"`
}

// CreditCardBalanceResponse represents the balance information for a credit card
type CreditCardBalanceResponse struct {
	CardID          string      `json:"card_id"`
	Balance         float64     `json:"balance"`            // Current debt
	CreditLimit     float64     `json:"credit_limit"`       // Total credit limit
	AvailableCredit float64     `json:"available_credit"`   // Remaining credit
	MinimumPayment  float64     `json:"minimum_payment"`    // Minimum payment due
	DueDate         *CustomDate `json:"due_date,omitempty"` // Next payment due date
}

// Debit Card Operations DTOs

// DebitCardTransactionRequest represents a transaction with a debit card
type DebitCardTransactionRequest struct {
	Amount       float64 `json:"amount" binding:"required,min=0.01"`
	Description  string  `json:"description" binding:"required,min=3,max=255"`
	MerchantName string  `json:"merchant_name,omitempty" binding:"max=100"`
	Reference    string  `json:"reference,omitempty" binding:"max=50"`
}

// DebitCardBalanceResponse represents the balance information for a debit card
type DebitCardBalanceResponse struct {
	CardID           string  `json:"card_id"`
	AccountBalance   float64 `json:"account_balance"`   // Current account balance
	AvailableBalance float64 `json:"available_balance"` // Available balance (same as account)
}

// PaginatedCardResponse represents paginated card list response
type PaginatedCardResponse struct {
	Data       []CardResponse `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

// PaginationMeta represents pagination metadata (reused from account dto)
type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
}

// ToCardResponse converts a Card entity to CardResponse DTO
func ToCardResponse(card *entities.Card) CardResponse {
	return CardResponse{
		ID:              card.ID,
		AccountID:       card.AccountID,
		CardType:        string(card.CardType),
		CardBrand:       string(card.CardBrand),
		LastFourDigits:  card.LastFourDigits,
		MaskedNumber:    card.MaskedNumber,
		HolderName:      card.HolderName,
		ExpirationMonth: card.ExpirationMonth,
		ExpirationYear:  card.ExpirationYear,
		Status:          string(card.Status),
		IsDefault:       card.IsDefault,
		Nickname:        card.Nickname,
		Balance:         card.Balance, // New: Include balance
		CreatedAt:       card.CreatedAt,
		UpdatedAt:       card.UpdatedAt,
		CreditLimit:     card.CreditLimit,
		ClosingDate:     ToCustomDate(card.ClosingDate),
		DueDate:         ToCustomDate(card.DueDate),
	}
}

// ToPaginatedCardResponse converts cards with pagination info to response
func ToPaginatedCardResponse(cards []*entities.Card, total int64, page, pageSize int) PaginatedCardResponse {
	data := make([]CardResponse, len(cards))
	for i, card := range cards {
		data[i] = ToCardResponse(card)
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return PaginatedCardResponse{
		Data: data,
		Pagination: PaginationMeta{
			CurrentPage: page,
			PageSize:    pageSize,
			TotalItems:  total,
			TotalPages:  totalPages,
		},
	}
}

// ToCreditCardBalanceResponse converts card entity to credit balance response
func ToCreditCardBalanceResponse(card *entities.Card) CreditCardBalanceResponse {
	var creditLimit, availableCredit, minimumPayment float64

	if card.CreditLimit != nil {
		creditLimit = *card.CreditLimit
		availableCredit = creditLimit - card.Balance
	}

	// Calculate minimum payment (5% of balance or minimum $500)
	if card.Balance > 0 {
		minimumPayment = card.Balance * 0.05
		if minimumPayment < 500 {
			minimumPayment = 500
		}
		// Can't be more than the total balance
		if minimumPayment > card.Balance {
			minimumPayment = card.Balance
		}
	}

	return CreditCardBalanceResponse{
		CardID:          card.ID,
		Balance:         card.Balance,
		CreditLimit:     creditLimit,
		AvailableCredit: availableCredit,
		MinimumPayment:  minimumPayment,
		DueDate:         ToCustomDate(card.DueDate),
	}
}

// ToDebitCardBalanceResponse converts card entity to debit balance response
func ToDebitCardBalanceResponse(card *entities.Card) DebitCardBalanceResponse {
	// For debit cards, available balance comes from the associated account
	accountBalance := card.Account.Balance

	return DebitCardBalanceResponse{
		CardID:           card.ID,
		AccountBalance:   accountBalance,
		AvailableBalance: accountBalance, // Same as account balance for debit cards
	}
}

// ChargeWithInstallmentsResponse represents the response after charging a card with installments
type ChargeWithInstallmentsResponse struct {
	InstallmentPlan         *entities.InstallmentPlan `json:"installment_plan"`
	Card                    *entities.Card            `json:"card"`
	FirstInstallmentCharged bool                      `json:"first_installment_charged"`
	TransactionID           string                    `json:"transaction_id,omitempty"`
}
