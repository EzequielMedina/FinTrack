package ports

import (
	"github.com/fintrack/account-service/internal/core/domain/entities"
	"github.com/fintrack/account-service/internal/infrastructure/entrypoints/handlers/card/dto"
)

// CardServiceInterface defines the contract for card service operations
type CardServiceInterface interface {
	// Basic CRUD operations
	CreateCard(req *dto.CreateCardRequest) (*entities.Card, error)
	GetCardByID(cardID string) (*entities.Card, error)
	GetCardByIDWithAccount(cardID string) (*entities.Card, error) // New: get card with account preloaded
	GetCardsByAccount(accountID string, page, pageSize int) ([]*entities.Card, int64, error)
	GetCardsByUser(userID string, page, pageSize int) ([]*entities.Card, int64, error)
	UpdateCard(cardID string, req *dto.UpdateCardRequest) (*entities.Card, error)
	DeleteCard(cardID string) error
	BlockCard(cardID string) (*entities.Card, error)
	UnblockCard(cardID string) (*entities.Card, error)
	SetDefaultCard(cardID string) (*entities.Card, error)

	// Credit card financial operations
	ChargeCard(cardID string, amount float64, description, reference string) (*entities.Card, error)
	PaymentCard(cardID string, amount float64, paymentMethod, reference string) (*entities.Card, error)

	// Debit card operations
	ProcessDebitTransaction(cardID string, amount float64, description, merchantName, reference string) (*entities.Card, error)
}

// CardRepositoryInterface defines the contract for card repository operations
type CardRepositoryInterface interface {
	Create(card *entities.Card) (*entities.Card, error)
	GetByID(cardID string) (*entities.Card, error)
	GetByIDWithAccount(cardID string) (*entities.Card, error) // New: get card with account preloaded
	GetByAccount(accountID string, limit, offset int) ([]*entities.Card, int64, error)
	GetByUser(userID string, limit, offset int) ([]*entities.Card, int64, error)
	Update(card *entities.Card) (*entities.Card, error)
	Delete(cardID string) error
	GetDefaultByAccount(accountID string) (*entities.Card, error)
	SetDefaultByAccount(accountID, cardID string) error
}

// AccountRepositoryInterface defines the contract for account repository operations
type AccountRepositoryInterface interface {
	Create(account *entities.Account) error
	GetByID(id string) (*entities.Account, error)
	GetByUserID(userID string) ([]*entities.Account, error)
	GetAll(limit, offset int) ([]*entities.Account, int64, error)
	Update(account *entities.Account) error
	Delete(id string) error
}
