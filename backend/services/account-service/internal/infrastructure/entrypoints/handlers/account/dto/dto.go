package dto

import (
	"time"

	"github.com/fintrack/account-service/internal/core/domain/entities"
)

// CreateAccountRequest represents the request to create a new account
type CreateAccountRequest struct {
	UserID         string  `json:"user_id" binding:"required"`
	AccountType    string  `json:"account_type" binding:"required"`
	Name           string  `json:"name" binding:"required"`
	Description    string  `json:"description"`
	Currency       string  `json:"currency" binding:"required"`
	InitialBalance float64 `json:"initial_balance" binding:"min=0"`

	// Credit card specific fields
	CreditLimit *float64   `json:"credit_limit,omitempty" binding:"omitempty,min=0"`
	ClosingDate *time.Time `json:"closing_date,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`

	// Personal identification (for virtual wallets)
	DNI *string `json:"dni,omitempty" binding:"omitempty,min=7,max=20"`
}

// UpdateAccountRequest represents the request to update an account
type UpdateAccountRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`

	// Credit card specific fields
	CreditLimit *float64   `json:"credit_limit,omitempty" binding:"omitempty,min=0"`
	ClosingDate *time.Time `json:"closing_date,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`

	// Personal identification (for virtual wallets)
	DNI *string `json:"dni,omitempty" binding:"omitempty,min=7,max=20"`
}

// UpdateBalanceRequest represents the request to update account balance
type UpdateBalanceRequest struct {
	Amount float64 `json:"amount" binding:"required"`
}

// UpdateStatusRequest represents the request to update account status
type UpdateStatusRequest struct {
	IsActive bool `json:"is_active"`
}

// AddFundsRequest represents the request to add funds to a wallet
type AddFundsRequest struct {
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description" binding:"required,min=3,max=255"`
	Reference   string  `json:"reference,omitempty" binding:"max=50"`
}

// WithdrawFundsRequest represents the request to withdraw funds from a wallet
type WithdrawFundsRequest struct {
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description" binding:"required,min=3,max=255"`
	Reference   string  `json:"reference,omitempty" binding:"max=50"`
}

// UpdateCreditLimitRequest represents the request to update credit limit
type UpdateCreditLimitRequest struct {
	CreditLimit float64 `json:"credit_limit" binding:"min=0"`
}

// UpdateCreditDatesRequest represents the request to update credit card dates
type UpdateCreditDatesRequest struct {
	ClosingDate *time.Time `json:"closing_date,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

// AvailableCreditResponse represents the response for available credit operations
type AvailableCreditResponse struct {
	AccountID       string  `json:"account_id"`
	CreditLimit     float64 `json:"credit_limit"`
	UsedCredit      float64 `json:"used_credit"`
	AvailableCredit float64 `json:"available_credit"`
}

// AccountResponse represents the response for account operations
type AccountResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	AccountType string    `json:"account_type"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Currency    string    `json:"currency"`
	Balance     float64   `json:"balance"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Credit card specific fields
	CreditLimit *float64   `json:"credit_limit,omitempty"`
	ClosingDate *time.Time `json:"closing_date,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`

	// Personal identification (for virtual wallets)
	DNI *string `json:"dni,omitempty"`
}

// BalanceResponse represents the response for balance operations
type BalanceResponse struct {
	AccountID string  `json:"account_id"`
	Balance   float64 `json:"balance"`
}

// PaginatedAccountResponse represents paginated account list response
type PaginatedAccountResponse struct {
	Data       []AccountResponse `json:"data"`
	Pagination PaginationMeta    `json:"pagination"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
}

// ToAccountResponse converts an Account entity to AccountResponse DTO
func ToAccountResponse(account *entities.Account) AccountResponse {
	return AccountResponse{
		ID:          account.ID,
		UserID:      account.UserID,
		AccountType: string(account.AccountType),
		Name:        account.Name,
		Description: account.Description,
		Currency:    string(account.Currency),
		Balance:     account.Balance,
		IsActive:    account.IsActive,
		CreatedAt:   account.CreatedAt,
		UpdatedAt:   account.UpdatedAt,
		CreditLimit: account.CreditLimit,
		ClosingDate: account.ClosingDate,
		DueDate:     account.DueDate,
		DNI:         account.DNI,
	}
}

// ToPaginatedAccountResponse converts accounts with pagination info to response
func ToPaginatedAccountResponse(accounts []*entities.Account, total int64, page, pageSize int) PaginatedAccountResponse {
	data := make([]AccountResponse, len(accounts))
	for i, account := range accounts {
		data[i] = ToAccountResponse(account)
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return PaginatedAccountResponse{
		Data: data,
		Pagination: PaginationMeta{
			CurrentPage: page,
			PageSize:    pageSize,
			TotalItems:  total,
			TotalPages:  totalPages,
		},
	}
}
