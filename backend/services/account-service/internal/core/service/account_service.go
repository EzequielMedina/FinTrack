package service

import (
	"fmt"

	"github.com/fintrack/account-service/internal/core/domain/entities"
	"github.com/fintrack/account-service/internal/infrastructure/repositories"
)

// AccountService provides business logic for account operations
type AccountService struct {
	accountRepo repositories.AccountRepository
}

// NewAccountService creates a new account service instance
func NewAccountService(accountRepo repositories.AccountRepository) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
	}
}

// CreateAccount creates a new account
func (s *AccountService) CreateAccount(account *entities.Account) (*entities.Account, error) {
	// Validate input
	if account.UserID == "" {
		return nil, fmt.Errorf("user ID is required")
	}
	if account.Name == "" {
		return nil, fmt.Errorf("account name is required")
	}
	if account.AccountType == "" {
		return nil, fmt.Errorf("account type is required")
	}
	if account.Currency == "" {
		return nil, fmt.Errorf("currency is required")
	}

	// Save to database
	if err := s.accountRepo.Create(account); err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return account, nil
}

// GetAccountByID retrieves an account by its ID
func (s *AccountService) GetAccountByID(accountID string) (*entities.Account, error) {
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	account, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return account, nil
}

// GetAccountsByUserID retrieves all accounts for a user
func (s *AccountService) GetAccountsByUserID(userID string) ([]*entities.Account, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	accounts, err := s.accountRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts for user: %w", err)
	}

	return accounts, nil
}

// GetAllAccounts retrieves all accounts with pagination
func (s *AccountService) GetAllAccounts(page, pageSize int) ([]*entities.Account, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	accounts, total, err := s.accountRepo.GetAll(pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get accounts: %w", err)
	}

	return accounts, total, nil
}

// UpdateAccount updates an existing account
func (s *AccountService) UpdateAccount(accountID, name, description string) (*entities.Account, error) {
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	// Get existing account
	account, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	// Update fields
	if name != "" {
		account.Name = name
	}
	account.Description = description

	// Save changes
	if err := s.accountRepo.Update(account); err != nil {
		return nil, fmt.Errorf("failed to update account: %w", err)
	}

	return account, nil
}

// DeleteAccount deletes an account
func (s *AccountService) DeleteAccount(accountID string) error {
	if accountID == "" {
		return fmt.Errorf("account ID is required")
	}

	// Get existing account to check balance
	account, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		return fmt.Errorf("failed to get account: %w", err)
	}

	// Check if account has balance
	if account.Balance > 0 {
		return fmt.Errorf("cannot delete account with balance: current balance %.2f", account.Balance)
	}

	if err := s.accountRepo.Delete(accountID); err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}

	return nil
}

// GetAccountBalance retrieves the balance of an account
func (s *AccountService) GetAccountBalance(accountID string) (float64, error) {
	if accountID == "" {
		return 0, fmt.Errorf("account ID is required")
	}

	account, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		return 0, fmt.Errorf("failed to get account: %w", err)
	}

	return account.Balance, nil
}

// UpdateAccountBalance updates the balance of an account
func (s *AccountService) UpdateAccountBalance(accountID string, amount float64) (float64, error) {
	if accountID == "" {
		return 0, fmt.Errorf("account ID is required")
	}

	// Get existing account
	account, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		return 0, fmt.Errorf("failed to get account: %w", err)
	}

	// Calculate new balance
	newBalance := account.Balance + amount

	// Validate new balance is not negative
	if newBalance < 0 {
		return 0, fmt.Errorf("insufficient balance: current balance %.2f, requested change %.2f", account.Balance, amount)
	}

	// Update balance
	account.Balance = newBalance

	// Save changes
	if err := s.accountRepo.Update(account); err != nil {
		return 0, fmt.Errorf("failed to update account balance: %w", err)
	}

	return account.Balance, nil
}

// UpdateAccountStatus updates the status of an account
func (s *AccountService) UpdateAccountStatus(accountID string, isActive bool) (*entities.Account, error) {
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	// Get existing account
	account, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	// Update status
	account.IsActive = isActive

	// Save changes
	if err := s.accountRepo.Update(account); err != nil {
		return nil, fmt.Errorf("failed to update account status: %w", err)
	}

	return account, nil
}
