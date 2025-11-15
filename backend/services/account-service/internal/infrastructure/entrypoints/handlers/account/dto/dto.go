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
	IsActive       *bool   `json:"is_active,omitempty"`

	// Credit card specific fields
	CreditLimit *float64   `json:"credit_limit,omitempty" binding:"omitempty,min=0"`
	ClosingDate *time.Time `json:"closing_date,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`

	// Personal identification (for virtual wallets)
	DNI *string `json:"dni,omitempty" binding:"omitempty,min=7,max=20"`
}

// UpdateAccountRequest represents the request to update an account
type UpdateAccountRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	AccountType string `json:"account_type,omitempty" binding:"omitempty,oneof=checking savings credit debit wallet bank_account"`

	// Credit card specific fields
	CreditLimit *float64   `json:"credit_limit,omitempty" binding:"omitempty,min=0"`
	ClosingDate *time.Time `json:"closing_date,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`

	// Personal identification (for virtual wallets)
	DNI *string `json:"dni,omitempty"`
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

// CardResponse represents the response for card operations
type CardResponse struct {
	ID              string     `json:"id"`
	AccountID       string     `json:"account_id"`
	CardType        string     `json:"card_type"`
	CardBrand       string     `json:"card_brand"`
	LastFourDigits  string     `json:"last_four_digits"`
	MaskedNumber    string     `json:"masked_number"`
	HolderName      string     `json:"holder_name"`
	ExpirationMonth int        `json:"expiration_month"`
	ExpirationYear  int        `json:"expiration_year"`
	Status          string     `json:"status"`
	IsDefault       bool       `json:"is_default"`
	Nickname        string     `json:"nickname"`
	Balance         float64    `json:"balance"`
	CreditLimit     *float64   `json:"credit_limit,omitempty"`
	ClosingDate     *time.Time `json:"closing_date,omitempty"`
	DueDate         *time.Time `json:"due_date,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
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

	// Cards relationship (for bank_account type)
	Cards []CardResponse `json:"cards,omitempty"`

	// Credit card specific fields (legacy - for backward compatibility)
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
	// Convert cards if present
	var cards []CardResponse
	if len(account.Cards) > 0 {
		cards = make([]CardResponse, len(account.Cards))
		for i, card := range account.Cards {
			cards[i] = CardResponse{
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
				Balance:         card.Balance,
				CreditLimit:     card.CreditLimit,
				ClosingDate:     card.ClosingDate,
				DueDate:         card.DueDate,
				CreatedAt:       card.CreatedAt,
				UpdatedAt:       card.UpdatedAt,
			}
		}
	}

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
		Cards:       cards,
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
