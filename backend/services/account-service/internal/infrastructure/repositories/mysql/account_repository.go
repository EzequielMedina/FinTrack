package mysql

import (
	"github.com/fintrack/account-service/internal/core/domain/entities"
	"github.com/fintrack/account-service/internal/infrastructure/repositories"
	"gorm.io/gorm"
)

// AccountRepository implements the account repository using GORM
type AccountRepository struct {
	db *gorm.DB
}

// NewAccountRepository creates a new account repository
func NewAccountRepository(db *gorm.DB) repositories.AccountRepository {
	return &AccountRepository{db: db}
}

// Create saves a new account to the database
func (r *AccountRepository) Create(account *entities.Account) error {
	return r.db.Create(account).Error
}

// GetByID retrieves an account by its ID
func (r *AccountRepository) GetByID(id string) (*entities.Account, error) {
	var account entities.Account
	err := r.db.Preload("Cards").Where("id = ? AND deleted_at IS NULL", id).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// GetByUserID retrieves all accounts for a user
func (r *AccountRepository) GetByUserID(userID string) ([]*entities.Account, error) {
	var accounts []*entities.Account
	err := r.db.Preload("Cards").Where("user_id = ? AND deleted_at IS NULL", userID).Find(&accounts).Error
	return accounts, err
}

// GetAll retrieves all accounts with pagination
func (r *AccountRepository) GetAll(limit, offset int) ([]*entities.Account, int64, error) {
	var accounts []*entities.Account
	var total int64

	// Count total records
	if err := r.db.Model(&entities.Account{}).Where("deleted_at IS NULL").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results with cards preloaded
	err := r.db.Preload("Cards").Where("deleted_at IS NULL").
		Limit(limit).
		Offset(offset).
		Find(&accounts).Error

	return accounts, total, err
}

// Update updates an existing account
func (r *AccountRepository) Update(account *entities.Account) error {
	return r.db.Save(account).Error
}

// Delete performs soft delete on an account
func (r *AccountRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&entities.Account{}).Error
}
