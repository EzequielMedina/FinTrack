package mysql

import (
	"github.com/fintrack/account-service/internal/core/domain/entities"
	"github.com/fintrack/account-service/internal/core/ports"
	"gorm.io/gorm"
)

// CardRepository implements the card repository using GORM
type CardRepository struct {
	db *gorm.DB
}

// NewCardRepository creates a new card repository
func NewCardRepository(db *gorm.DB) ports.CardRepositoryInterface {
	return &CardRepository{db: db}
}

// Create saves a new card to the database
func (r *CardRepository) Create(card *entities.Card) (*entities.Card, error) {
	if err := r.db.Create(card).Error; err != nil {
		return nil, err
	}
	return card, nil
}

// GetByID retrieves a card by its ID
func (r *CardRepository) GetByID(cardID string) (*entities.Card, error) {
	var card entities.Card
	err := r.db.Where("id = ? AND deleted_at IS NULL", cardID).First(&card).Error
	if err != nil {
		return nil, err
	}
	return &card, nil
}

// GetByAccount retrieves all cards for a specific account with pagination
func (r *CardRepository) GetByAccount(accountID string, limit, offset int) ([]*entities.Card, int64, error) {
	var cards []*entities.Card
	var total int64

	// Count total cards for the account
	err := r.db.Model(&entities.Card{}).
		Where("account_id = ? AND deleted_at IS NULL", accountID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get cards with pagination
	err = r.db.Where("account_id = ? AND deleted_at IS NULL", accountID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&cards).Error

	return cards, total, err
}

// GetByUser retrieves all cards for a specific user across all accounts with pagination
func (r *CardRepository) GetByUser(userID string, limit, offset int) ([]*entities.Card, int64, error) {
	var cards []*entities.Card
	var total int64

	// Count total cards for the user
	err := r.db.Model(&entities.Card{}).
		Joins("INNER JOIN accounts ON accounts.id = cards.account_id").
		Where("accounts.user_id = ? AND cards.deleted_at IS NULL AND accounts.deleted_at IS NULL", userID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get cards with pagination
	err = r.db.
		Joins("INNER JOIN accounts ON accounts.id = cards.account_id").
		Where("accounts.user_id = ? AND cards.deleted_at IS NULL AND accounts.deleted_at IS NULL", userID).
		Limit(limit).
		Offset(offset).
		Order("cards.created_at DESC").
		Find(&cards).Error

	return cards, total, err
}

// Update updates an existing card
func (r *CardRepository) Update(card *entities.Card) (*entities.Card, error) {
	err := r.db.Save(card).Error
	if err != nil {
		return nil, err
	}
	return card, nil
}

// Delete performs soft delete on a card
func (r *CardRepository) Delete(cardID string) error {
	return r.db.Where("id = ?", cardID).Delete(&entities.Card{}).Error
}

// GetDefaultByAccount retrieves the default card for a specific account
func (r *CardRepository) GetDefaultByAccount(accountID string) (*entities.Card, error) {
	var card entities.Card
	err := r.db.Where("account_id = ? AND is_default = ? AND deleted_at IS NULL", accountID, true).First(&card).Error
	if err != nil {
		return nil, err
	}
	return &card, nil
}

// SetDefaultByAccount sets a card as default for an account (and unsets others)
func (r *CardRepository) SetDefaultByAccount(accountID, cardID string) error {
	// Start a transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Unset all other cards as default for this account
	err := tx.Model(&entities.Card{}).
		Where("account_id = ? AND id != ?", accountID, cardID).
		Update("is_default", false).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Set the specified card as default
	err = tx.Model(&entities.Card{}).
		Where("id = ? AND account_id = ?", cardID, accountID).
		Update("is_default", true).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// GetByIDWithAccount retrieves a card by its ID with account data preloaded
func (r *CardRepository) GetByIDWithAccount(cardID string) (*entities.Card, error) {
	var card entities.Card
	err := r.db.Preload("Account").Where("id = ? AND deleted_at IS NULL", cardID).First(&card).Error
	if err != nil {
		return nil, err
	}
	return &card, nil
}

// GetWithInstallmentPlans retrieves a card with its installment plans preloaded
func (r *CardRepository) GetWithInstallmentPlans(cardID string) (*entities.Card, error) {
	var card entities.Card
	err := r.db.
		Preload("Account").
		Preload("InstallmentPlans", "status IN (?, ?, ?)", "active", "suspended", "pending").
		Preload("InstallmentPlans.Installments", "status IN (?, ?)", "pending", "overdue").
		Where("id = ? AND deleted_at IS NULL", cardID).
		First(&card).Error
	if err != nil {
		return nil, err
	}
	return &card, nil
}

// GetByAccountWithInstallmentPlans retrieves all cards for an account with installment plans preloaded
func (r *CardRepository) GetByAccountWithInstallmentPlans(accountID string, limit, offset int) ([]*entities.Card, int64, error) {
	var cards []*entities.Card
	var total int64

	// Count total cards for the account
	err := r.db.Model(&entities.Card{}).
		Where("account_id = ? AND deleted_at IS NULL", accountID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get cards with installment plans preloaded
	err = r.db.
		Preload("InstallmentPlans", "status IN (?, ?, ?)", "active", "suspended", "pending").
		Preload("InstallmentPlans.Installments", "status IN (?, ?)", "pending", "overdue").
		Where("account_id = ? AND deleted_at IS NULL", accountID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&cards).Error

	return cards, total, err
}
