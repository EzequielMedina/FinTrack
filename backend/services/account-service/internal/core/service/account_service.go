package service

import (
	"fmt"

	"github.com/fintrack/account-service/internal/core/domain/entities"
	"github.com/fintrack/account-service/internal/infrastructure/entrypoints/handlers/account/dto"
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
func (s *AccountService) UpdateAccount(accountID string, req *dto.UpdateAccountRequest) (*entities.Account, error) {
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	// Get existing account
	account, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	fmt.Printf("🔄 DEBUG - UpdateAccount called for account ID: %s\n", accountID)
	fmt.Printf("🔄 DEBUG - Current account: Name=%s, Type=%s, Description=%s\n",
		account.Name, account.AccountType, account.Description)

	// Track if any updates were made
	updated := false

	// Update fields if provided
	if req.Name != "" && req.Name != account.Name {
		fmt.Printf("🔄 DEBUG - Updating Name from '%s' to '%s'\n", account.Name, req.Name)
		account.Name = req.Name
		updated = true
	}

	if req.Description != account.Description {
		fmt.Printf("🔄 DEBUG - Updating Description from '%s' to '%s'\n", account.Description, req.Description)
		account.Description = req.Description
		updated = true
	}

	// Handle account type change (with validations)
	if req.AccountType != "" && req.AccountType != string(account.AccountType) {
		fmt.Printf("🔄 DEBUG - Updating AccountType from '%s' to '%s'\n", account.AccountType, req.AccountType)

		// Validate account type change is allowed
		if err := s.validateAccountTypeChange(account, entities.AccountType(req.AccountType)); err != nil {
			return nil, fmt.Errorf("cannot change account type: %w", err)
		}

		account.AccountType = entities.AccountType(req.AccountType)
		updated = true
	}

	// Handle credit limit updates
	if req.CreditLimit != nil {
		currentLimit := float64(0)
		if account.CreditLimit != nil {
			currentLimit = *account.CreditLimit
		}

		if *req.CreditLimit != currentLimit {
			fmt.Printf("🔄 DEBUG - Updating CreditLimit from %.2f to %.2f\n", currentLimit, *req.CreditLimit)
			account.CreditLimit = req.CreditLimit
			updated = true
		}
	}

	// Handle credit dates
	if req.ClosingDate != nil {
		account.ClosingDate = req.ClosingDate
		updated = true
	}

	if req.DueDate != nil {
		account.DueDate = req.DueDate
		updated = true
	}

	// Handle DNI
	if req.DNI != nil {
		currentDNI := ""
		if account.DNI != nil {
			currentDNI = *account.DNI
		}

		if *req.DNI != currentDNI {
			fmt.Printf("🔄 DEBUG - Updating DNI\n")
			account.DNI = req.DNI
			updated = true
		}
	}

	// If no updates were made, return the existing account
	if !updated {
		fmt.Printf("🔄 DEBUG - No changes detected, returning existing account\n")
		return account, nil
	}

	// Validate updated account
	if err := account.Validate(); err != nil {
		return nil, fmt.Errorf("invalid account data: %w", err)
	}

	// Save changes
	fmt.Printf("🔄 DEBUG - Saving account updates to database\n")
	if err := s.accountRepo.Update(account); err != nil {
		return nil, fmt.Errorf("failed to update account: %w", err)
	}

	fmt.Printf("🔄 DEBUG - Account updated successfully\n")
	return account, nil
}

// validateAccountTypeChange validates if an account type change is allowed
func (s *AccountService) validateAccountTypeChange(account *entities.Account, newType entities.AccountType) error {
	// Define allowed transitions
	allowedTransitions := map[entities.AccountType][]entities.AccountType{
		entities.AccountTypeChecking:    {entities.AccountTypeSavings, entities.AccountTypeBankAccount},
		entities.AccountTypeSavings:     {entities.AccountTypeChecking, entities.AccountTypeBankAccount},
		entities.AccountTypeBankAccount: {entities.AccountTypeChecking, entities.AccountTypeSavings},
		entities.AccountTypeCredit:      {}, // Credit accounts usually can't change type
		entities.AccountTypeDebit:       {entities.AccountTypeBankAccount},
		entities.AccountTypeWallet:      {}, // Wallet accounts usually can't change type
	}

	allowedTypes, exists := allowedTransitions[account.AccountType]
	if !exists {
		return fmt.Errorf("account type '%s' cannot be changed", account.AccountType)
	}

	// Check if the new type is in the allowed list
	for _, allowedType := range allowedTypes {
		if newType == allowedType {
			return nil // Change is allowed
		}
	}

	return fmt.Errorf("cannot change account type from '%s' to '%s'", account.AccountType, newType)
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
