package dto

import (
	"time"

	"github.com/fintrack/account-service/internal/core/domain/entities"
)

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
	CreditLimit *float64   `json:"credit_limit,omitempty" binding:"omitempty,min=0"`
	ClosingDate *time.Time `json:"closing_date,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`

	// Security fields (encrypted data)
	EncryptedNumber string `json:"encrypted_number" binding:"required"`
	KeyFingerprint  string `json:"key_fingerprint" binding:"required"`
}

// UpdateCardRequest represents the request to update a card
type UpdateCardRequest struct {
	HolderName      string `json:"holder_name,omitempty" binding:"omitempty,min=2,max=100"`
	ExpirationMonth int    `json:"expiration_month,omitempty" binding:"omitempty,min=1,max=12"`
	ExpirationYear  int    `json:"expiration_year,omitempty" binding:"omitempty,min=2024"`
	Nickname        string `json:"nickname,omitempty" binding:"max=50"`
	IsDefault       *bool  `json:"is_default,omitempty"`
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
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// Credit card specific fields
	CreditLimit *float64   `json:"credit_limit,omitempty"`
	ClosingDate *time.Time `json:"closing_date,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`

	// Security fields (encrypted data - never exposed in response)
	// EncryptedNumber and KeyFingerprint are intentionally omitted for security
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
		CreatedAt:       card.CreatedAt,
		UpdatedAt:       card.UpdatedAt,
		CreditLimit:     card.CreditLimit,
		ClosingDate:     card.ClosingDate,
		DueDate:         card.DueDate,
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
