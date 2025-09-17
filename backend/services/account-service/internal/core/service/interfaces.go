package service

import (
	"github.com/fintrack/account-service/internal/core/domain/entities"
)

// AccountServiceInterface defines the contract for account service operations
type AccountServiceInterface interface {
	// Basic CRUD operations
	CreateAccount(account *entities.Account) (*entities.Account, error)
	GetAccountByID(accountID string) (*entities.Account, error)
	GetAccountsByUserID(userID string) ([]*entities.Account, error)
	GetAllAccounts(page, pageSize int) ([]*entities.Account, int64, error)
	UpdateAccount(accountID, name, description string) (*entities.Account, error)
	DeleteAccount(accountID string) error

	// Balance operations
	GetAccountBalance(accountID string) (float64, error)
	UpdateAccountBalance(accountID string, amount float64) (float64, error)

	// Status operations
	UpdateAccountStatus(accountID string, isActive bool) (*entities.Account, error)
}
