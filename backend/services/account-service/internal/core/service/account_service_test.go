package service

import (
	"testing"

	"github.com/fintrack/account-service/internal/core/domain/entities"
	"github.com/fintrack/account-service/internal/core/errors"
	"github.com/fintrack/account-service/internal/infrastructure/repositories"
	"github.com/google/uuid"
)

// MockAccountRepository implements a mock repository for testing
type MockAccountRepository struct {
	accounts map[string]*entities.Account
	byUser   map[string][]*entities.Account
}

func NewMockAccountRepository() *MockAccountRepository {
	return &MockAccountRepository{
		accounts: make(map[string]*entities.Account),
		byUser:   make(map[string][]*entities.Account),
	}
}

func (m *MockAccountRepository) Create(account *entities.Account) error {
	if account.ID == "" {
		account.ID = uuid.NewString()
	}

	m.accounts[account.ID] = account
	m.byUser[account.UserID] = append(m.byUser[account.UserID], account)
	return nil
}

func (m *MockAccountRepository) GetByID(id string) (*entities.Account, error) {
	account, exists := m.accounts[id]
	if !exists {
		return nil, errors.ErrAccountNotFound
	}
	return account, nil
}

func (m *MockAccountRepository) GetByUserID(userID string) ([]*entities.Account, error) {
	return m.byUser[userID], nil
}

func (m *MockAccountRepository) GetAll(limit, offset int) ([]*entities.Account, int64, error) {
	var accounts []*entities.Account
	for _, account := range m.accounts {
		accounts = append(accounts, account)
	}
	return accounts, int64(len(accounts)), nil
}

func (m *MockAccountRepository) Update(account *entities.Account) error {
	if _, exists := m.accounts[account.ID]; !exists {
		return errors.ErrAccountNotFound
	}
	m.accounts[account.ID] = account
	return nil
}

func (m *MockAccountRepository) Delete(id string) error {
	account, exists := m.accounts[id]
	if !exists {
		return errors.ErrAccountNotFound
	}

	// Remove from user accounts
	userAccounts := m.byUser[account.UserID]
	for i, acc := range userAccounts {
		if acc.ID == id {
			m.byUser[account.UserID] = append(userAccounts[:i], userAccounts[i+1:]...)
			break
		}
	}

	delete(m.accounts, id)
	return nil
}

// Verify interface compliance
var _ repositories.AccountRepository = (*MockAccountRepository)(nil)

func TestCreateAccount(t *testing.T) {
	repo := NewMockAccountRepository()
	service := NewAccountService(repo)

	tests := []struct {
		name        string
		account     *entities.Account
		expectError bool
	}{
		{
			name: "valid savings account",
			account: &entities.Account{
				UserID:      uuid.NewString(),
				AccountType: entities.AccountTypeSavings,
				Name:        "My Savings",
				Currency:    entities.CurrencyUSD,
				Balance:     100.0,
				IsActive:    true,
			},
			expectError: false,
		},
		{
			name: "valid checking account",
			account: &entities.Account{
				UserID:      uuid.NewString(),
				AccountType: entities.AccountTypeChecking,
				Name:        "My Checking",
				Currency:    entities.CurrencyEUR,
				Balance:     0.0,
				IsActive:    true,
			},
			expectError: false,
		},
		{
			name: "account with empty user ID",
			account: &entities.Account{
				AccountType: entities.AccountTypeSavings,
				Name:        "Test Account",
				Currency:    entities.CurrencyUSD,
				Balance:     0.0,
				IsActive:    true,
			},
			expectError: true,
		},
		{
			name: "account with empty name",
			account: &entities.Account{
				UserID:      uuid.NewString(),
				AccountType: entities.AccountTypeSavings,
				Currency:    entities.CurrencyUSD,
				Balance:     0.0,
				IsActive:    true,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.CreateAccount(tt.account)

			if tt.expectError {
				if err == nil {
					t.Error("CreateAccount() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("CreateAccount() unexpected error: %v", err)
				}
				if result == nil {
					t.Error("CreateAccount() expected account but got nil")
				}
				if result.ID == "" {
					t.Error("CreateAccount() account ID should not be empty")
				}
				if result.UserID != tt.account.UserID {
					t.Errorf("CreateAccount() userID = %v, want %v", result.UserID, tt.account.UserID)
				}
				if result.AccountType != tt.account.AccountType {
					t.Errorf("CreateAccount() accountType = %v, want %v", result.AccountType, tt.account.AccountType)
				}
				if result.Currency != tt.account.Currency {
					t.Errorf("CreateAccount() currency = %v, want %v", result.Currency, tt.account.Currency)
				}
				if !result.IsActive {
					t.Error("CreateAccount() account should be active")
				}
			}
		})
	}
}

func TestGetAccountByID(t *testing.T) {
	repo := NewMockAccountRepository()
	service := NewAccountService(repo)

	// Create test account
	account := &entities.Account{
		UserID:      uuid.NewString(),
		AccountType: entities.AccountTypeSavings,
		Name:        "Test Account",
		Currency:    entities.CurrencyUSD,
		Balance:     100.0,
		IsActive:    true,
	}
	createdAccount, _ := service.CreateAccount(account)

	tests := []struct {
		name        string
		accountID   string
		expectError bool
	}{
		{
			name:        "existing account",
			accountID:   createdAccount.ID,
			expectError: false,
		},
		{
			name:        "non-existing account",
			accountID:   uuid.NewString(),
			expectError: true,
		},
		{
			name:        "empty account ID",
			accountID:   "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.GetAccountByID(tt.accountID)

			if tt.expectError {
				if err == nil {
					t.Error("GetAccountByID() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("GetAccountByID() unexpected error: %v", err)
				}
				if result == nil {
					t.Error("GetAccountByID() expected account but got nil")
				}
				if result.ID != tt.accountID {
					t.Errorf("GetAccountByID() accountID = %v, want %v", result.ID, tt.accountID)
				}
			}
		})
	}
}

func TestUpdateAccountBalance(t *testing.T) {
	repo := NewMockAccountRepository()
	service := NewAccountService(repo)

	// Create test account
	account := &entities.Account{
		UserID:      uuid.NewString(),
		AccountType: entities.AccountTypeSavings,
		Name:        "Test Account",
		Currency:    entities.CurrencyUSD,
		Balance:     100.0,
		IsActive:    true,
	}
	createdAccount, _ := service.CreateAccount(account)

	tests := []struct {
		name            string
		accountID       string
		amount          float64
		expectError     bool
		expectedBalance float64
	}{
		{
			name:            "valid balance update",
			accountID:       createdAccount.ID,
			amount:          50.0,
			expectError:     false,
			expectedBalance: 150.0,
		},
		{
			name:            "set balance to zero",
			accountID:       createdAccount.ID,
			amount:          -150.0,
			expectError:     false,
			expectedBalance: 0.0,
		},
		{
			name:        "negative final balance",
			accountID:   createdAccount.ID,
			amount:      -200.0,
			expectError: true,
		},
		{
			name:        "non-existing account",
			accountID:   uuid.NewString(),
			amount:      50.0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.UpdateAccountBalance(tt.accountID, tt.amount)

			if tt.expectError {
				if err == nil {
					t.Error("UpdateAccountBalance() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("UpdateAccountBalance() unexpected error: %v", err)
				}
				if result != tt.expectedBalance {
					t.Errorf("UpdateAccountBalance() balance = %v, want %v", result, tt.expectedBalance)
				}
			}
		})
	}
}

func TestUpdateAccountStatus(t *testing.T) {
	repo := NewMockAccountRepository()
	service := NewAccountService(repo)

	// Create test account
	account := &entities.Account{
		UserID:      uuid.NewString(),
		AccountType: entities.AccountTypeSavings,
		Name:        "Test Account",
		Currency:    entities.CurrencyUSD,
		Balance:     100.0,
		IsActive:    true,
	}
	createdAccount, _ := service.CreateAccount(account)

	tests := []struct {
		name        string
		accountID   string
		isActive    bool
		expectError bool
	}{
		{
			name:        "deactivate account",
			accountID:   createdAccount.ID,
			isActive:    false,
			expectError: false,
		},
		{
			name:        "reactivate account",
			accountID:   createdAccount.ID,
			isActive:    true,
			expectError: false,
		},
		{
			name:        "non-existing account",
			accountID:   uuid.NewString(),
			isActive:    true,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.UpdateAccountStatus(tt.accountID, tt.isActive)

			if tt.expectError {
				if err == nil {
					t.Error("UpdateAccountStatus() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("UpdateAccountStatus() unexpected error: %v", err)
				}
				if result == nil {
					t.Error("UpdateAccountStatus() expected account but got nil")
				}
				if result.IsActive != tt.isActive {
					t.Errorf("UpdateAccountStatus() isActive = %v, want %v", result.IsActive, tt.isActive)
				}
			}
		})
	}
}

func TestGetAccountsByUserID(t *testing.T) {
	repo := NewMockAccountRepository()
	service := NewAccountService(repo)

	userID := uuid.NewString()
	otherUserID := uuid.NewString()

	// Create multiple accounts for the user
	account1 := &entities.Account{
		UserID:      userID,
		AccountType: entities.AccountTypeSavings,
		Name:        "Savings Account",
		Currency:    entities.CurrencyUSD,
		Balance:     100.0,
		IsActive:    true,
	}
	account2 := &entities.Account{
		UserID:      userID,
		AccountType: entities.AccountTypeChecking,
		Name:        "Checking Account",
		Currency:    entities.CurrencyEUR,
		Balance:     50.0,
		IsActive:    true,
	}
	account3 := &entities.Account{
		UserID:      otherUserID,
		AccountType: entities.AccountTypeSavings,
		Name:        "Other Account",
		Currency:    entities.CurrencyUSD,
		Balance:     200.0,
		IsActive:    true,
	}

	service.CreateAccount(account1)
	service.CreateAccount(account2)
	service.CreateAccount(account3)

	tests := []struct {
		name          string
		userID        string
		expectError   bool
		expectedCount int
	}{
		{
			name:          "user with multiple accounts",
			userID:        userID,
			expectError:   false,
			expectedCount: 2,
		},
		{
			name:          "user with single account",
			userID:        otherUserID,
			expectError:   false,
			expectedCount: 1,
		},
		{
			name:          "user with no accounts",
			userID:        uuid.NewString(),
			expectError:   false,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.GetAccountsByUserID(tt.userID)

			if tt.expectError {
				if err == nil {
					t.Error("GetAccountsByUserID() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("GetAccountsByUserID() unexpected error: %v", err)
				}
				if len(result) != tt.expectedCount {
					t.Errorf("GetAccountsByUserID() account count = %v, want %v", len(result), tt.expectedCount)
				}
				// Verify all accounts belong to the user
				for _, account := range result {
					if account.UserID != tt.userID {
						t.Errorf("GetAccountsByUserID() account belongs to user %v, want %v", account.UserID, tt.userID)
					}
				}
			}
		})
	}
}

func TestDeleteAccount(t *testing.T) {
	repo := NewMockAccountRepository()
	service := NewAccountService(repo)

	// Create test account
	account := &entities.Account{
		UserID:      uuid.NewString(),
		AccountType: entities.AccountTypeSavings,
		Name:        "Test Account",
		Currency:    entities.CurrencyUSD,
		Balance:     0.0,
		IsActive:    true,
	}
	createdAccount, _ := service.CreateAccount(account)

	// Create another account with balance
	accountWithBalance := &entities.Account{
		UserID:      uuid.NewString(),
		AccountType: entities.AccountTypeChecking,
		Name:        "Account with Balance",
		Currency:    entities.CurrencyUSD,
		Balance:     100.0,
		IsActive:    true,
	}
	createdAccountWithBalance, _ := service.CreateAccount(accountWithBalance)

	tests := []struct {
		name        string
		accountID   string
		expectError bool
	}{
		{
			name:        "delete account with zero balance",
			accountID:   createdAccount.ID,
			expectError: false,
		},
		{
			name:        "cannot delete account with balance",
			accountID:   createdAccountWithBalance.ID,
			expectError: true,
		},
		{
			name:        "non-existing account",
			accountID:   uuid.NewString(),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.DeleteAccount(tt.accountID)

			if tt.expectError {
				if err == nil {
					t.Error("DeleteAccount() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("DeleteAccount() unexpected error: %v", err)
				}
				// Verify account is deleted
				_, getErr := service.GetAccountByID(tt.accountID)
				if getErr == nil {
					t.Error("DeleteAccount() account should be deleted but was found")
				}
			}
		})
	}
}
