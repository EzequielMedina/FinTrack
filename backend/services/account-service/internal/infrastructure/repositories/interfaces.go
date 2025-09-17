package repositories

import (
	"github.com/fintrack/account-service/internal/core/domain/entities"
)

// AccountRepository defines the interface for account data access
type AccountRepository interface {
	// Basic CRUD operations
	Create(account *entities.Account) error
	GetByID(id string) (*entities.Account, error)
	GetByUserID(userID string) ([]*entities.Account, error)
	GetAll(limit, offset int) ([]*entities.Account, int64, error)
	Update(account *entities.Account) error
	Delete(id string) error
}
