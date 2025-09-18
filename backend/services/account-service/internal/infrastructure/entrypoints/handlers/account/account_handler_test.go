package accounthandler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/fintrack/account-service/internal/core/domain/entities"
	"github.com/fintrack/account-service/internal/core/errors"
	"github.com/fintrack/account-service/internal/core/service"
	"github.com/fintrack/account-service/internal/infrastructure/entrypoints/handlers/account/dto"
)

// MockAccountService implements a mock AccountService for testing
type MockAccountService struct {
	accounts map[string]*entities.Account
	byUser   map[string][]*entities.Account
}

func NewMockAccountService() *MockAccountService {
	return &MockAccountService{
		accounts: make(map[string]*entities.Account),
		byUser:   make(map[string][]*entities.Account),
	}
}

func (m *MockAccountService) CreateAccount(account *entities.Account) (*entities.Account, error) {
	if account.UserID == "" {
		return nil, errors.ErrAccountNotFound
	}
	if account.Name == "" {
		return nil, errors.ErrInvalidAccountData
	}

	if account.ID == "" {
		account.ID = uuid.NewString()
	}
	account.CreatedAt = time.Now()
	account.UpdatedAt = time.Now()

	m.accounts[account.ID] = account
	m.byUser[account.UserID] = append(m.byUser[account.UserID], account)
	return account, nil
}

func (m *MockAccountService) GetAccountByID(id string) (*entities.Account, error) {
	account, exists := m.accounts[id]
	if !exists {
		return nil, errors.ErrAccountNotFound
	}
	return account, nil
}

func (m *MockAccountService) GetAccountsByUserID(userID string) ([]*entities.Account, error) {
	if userID == "" {
		return nil, errors.ErrInvalidAccountData
	}
	accounts := m.byUser[userID]
	return accounts, nil
}

func (m *MockAccountService) GetAllAccounts(page, pageSize int) ([]*entities.Account, int64, error) {
	var accounts []*entities.Account
	for _, account := range m.accounts {
		accounts = append(accounts, account)
	}
	return accounts, int64(len(accounts)), nil
}

func (m *MockAccountService) UpdateAccount(accountID, name, description string) (*entities.Account, error) {
	account, exists := m.accounts[accountID]
	if !exists {
		return nil, errors.ErrAccountNotFound
	}

	account.Name = name
	account.Description = description
	account.UpdatedAt = time.Now()
	return account, nil
}

func (m *MockAccountService) UpdateAccountBalance(accountID string, amount float64) (float64, error) {
	account, exists := m.accounts[accountID]
	if !exists {
		return 0, errors.ErrAccountNotFound
	}

	newBalance := account.Balance + amount
	if newBalance < 0 {
		return 0, errors.ErrInsufficientBalance
	}

	account.Balance = newBalance
	account.UpdatedAt = time.Now()
	return account.Balance, nil
}

func (m *MockAccountService) UpdateAccountStatus(accountID string, isActive bool) (*entities.Account, error) {
	account, exists := m.accounts[accountID]
	if !exists {
		return nil, errors.ErrAccountNotFound
	}
	account.IsActive = isActive
	account.UpdatedAt = time.Now()
	return account, nil
}

func (m *MockAccountService) DeleteAccount(id string) error {
	account, exists := m.accounts[id]
	if !exists {
		return errors.ErrAccountNotFound
	}
	if account.Balance > 0 {
		return errors.ErrCannotDeleteAccountWithBalance
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

func (m *MockAccountService) GetAccountBalance(accountID string) (float64, error) {
	account, exists := m.accounts[accountID]
	if !exists {
		return 0, errors.ErrAccountNotFound
	}
	return account.Balance, nil
}

// Verify interface compliance
var _ service.AccountServiceInterface = (*MockAccountService)(nil)

// Helper function to create a test Gin context
func createTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func TestCreateAccount(t *testing.T) {
	service := NewMockAccountService()
	handler := New(service)

	tests := []struct {
		name           string
		requestBody    dto.CreateAccountRequest
		expectedStatus int
		expectError    bool
	}{
		{
			name: "successful account creation",
			requestBody: dto.CreateAccountRequest{
				UserID:         uuid.NewString(),
				AccountType:    string(entities.AccountTypeSavings),
				Name:           "My Savings",
				Currency:       string(entities.CurrencyUSD),
				InitialBalance: 100.0,
			},
			expectedStatus: http.StatusCreated,
			expectError:    false,
		},
		{
			name: "missing user ID",
			requestBody: dto.CreateAccountRequest{
				AccountType:    string(entities.AccountTypeSavings),
				Name:           "Test Account",
				Currency:       string(entities.CurrencyUSD),
				InitialBalance: 0.0,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "missing account name",
			requestBody: dto.CreateAccountRequest{
				UserID:         uuid.NewString(),
				AccountType:    string(entities.AccountTypeSavings),
				Currency:       string(entities.CurrencyUSD),
				InitialBalance: 0.0,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "missing account type",
			requestBody: dto.CreateAccountRequest{
				UserID:         uuid.NewString(),
				Name:           "Test Account",
				Currency:       string(entities.CurrencyUSD),
				InitialBalance: 0.0,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "missing currency",
			requestBody: dto.CreateAccountRequest{
				UserID:         uuid.NewString(),
				AccountType:    string(entities.AccountTypeSavings),
				Name:           "Test Account",
				InitialBalance: 0.0,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := createTestContext()

			// Create request body
			body, _ := json.Marshal(tt.requestBody)
			c.Request = httptest.NewRequest(http.MethodPost, "/api/accounts", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Call handler
			handler.CreateAccount(c)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("CreateAccount() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			// Parse response
			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			// Check if error is expected
			if tt.expectError {
				if _, hasError := response["error"]; !hasError {
					t.Error("CreateAccount() expected error but got none")
				}
			} else {
				if _, hasData := response["id"]; !hasData {
					t.Error("CreateAccount() expected account data but got none")
				} else {
					if response["user_id"] != tt.requestBody.UserID {
						t.Errorf("CreateAccount() userID = %v, want %v", response["user_id"], tt.requestBody.UserID)
					}
				}
			}
		})
	}
}

func TestGetAccount(t *testing.T) {
	service := NewMockAccountService()
	handler := New(service)

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
		name           string
		accountID      string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "existing account",
			accountID:      createdAccount.ID,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "non-existing account",
			accountID:      uuid.NewString(),
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
		{
			name:           "invalid account ID",
			accountID:      "invalid-id",
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := createTestContext()

			// Set URL parameter
			c.Params = gin.Params{
				{Key: "id", Value: tt.accountID},
			}
			c.Request = httptest.NewRequest(http.MethodGet, "/api/accounts/"+tt.accountID, nil)

			// Call handler
			handler.GetAccount(c)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("GetAccount() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			// Parse response
			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			// Check if error is expected
			if tt.expectError {
				if _, hasError := response["error"]; !hasError {
					t.Error("GetAccount() expected error but got none")
				}
			} else {
				if _, hasData := response["id"]; !hasData {
					t.Error("GetAccount() expected account data but got none")
				} else {
					if response["id"] != tt.accountID {
						t.Errorf("GetAccount() id = %v, want %v", response["id"], tt.accountID)
					}
				}
			}
		})
	}
}

func TestUpdateBalance(t *testing.T) {
	service := NewMockAccountService()
	handler := New(service)

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
		name           string
		accountID      string
		requestBody    dto.UpdateBalanceRequest
		expectedStatus int
		expectError    bool
	}{
		{
			name:      "successful balance update",
			accountID: createdAccount.ID,
			requestBody: dto.UpdateBalanceRequest{
				Amount: 50.0,
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:      "zero amount",
			accountID: createdAccount.ID,
			requestBody: dto.UpdateBalanceRequest{
				Amount: 0.0,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:      "insufficient balance",
			accountID: createdAccount.ID,
			requestBody: dto.UpdateBalanceRequest{
				Amount: -200.0,
			},
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
		{
			name:      "non-existing account",
			accountID: uuid.NewString(),
			requestBody: dto.UpdateBalanceRequest{
				Amount: 100.0,
			},
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := createTestContext()

			// Set URL parameter
			c.Params = gin.Params{
				{Key: "id", Value: tt.accountID},
			}

			// Create request body
			body, _ := json.Marshal(tt.requestBody)
			c.Request = httptest.NewRequest(http.MethodPut, "/api/accounts/"+tt.accountID+"/balance", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Call handler
			handler.UpdateBalance(c)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("UpdateBalance() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			// Parse response
			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			// Check if error is expected
			if tt.expectError {
				if _, hasError := response["error"]; !hasError {
					t.Error("UpdateBalance() expected error but got none")
				}
			} else {
				if _, hasData := response["account_id"]; !hasData {
					t.Error("UpdateBalance() expected balance response but got none")
				} else {
					if response["account_id"] != tt.accountID {
						t.Errorf("UpdateBalance() accountID = %v, want %v", response["account_id"], tt.accountID)
					}
				}
			}
		})
	}
}

func TestUpdateStatus(t *testing.T) {
	service := NewMockAccountService()
	handler := New(service)

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
		name           string
		accountID      string
		requestBody    dto.UpdateStatusRequest
		expectedStatus int
		expectError    bool
	}{
		{
			name:      "deactivate account",
			accountID: createdAccount.ID,
			requestBody: dto.UpdateStatusRequest{
				IsActive: false,
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:      "reactivate account",
			accountID: createdAccount.ID,
			requestBody: dto.UpdateStatusRequest{
				IsActive: true,
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:      "non-existing account",
			accountID: uuid.NewString(),
			requestBody: dto.UpdateStatusRequest{
				IsActive: true,
			},
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := createTestContext()

			// Set URL parameter
			c.Params = gin.Params{
				{Key: "id", Value: tt.accountID},
			}

			// Create request body
			body, _ := json.Marshal(tt.requestBody)
			c.Request = httptest.NewRequest(http.MethodPut, "/api/accounts/"+tt.accountID+"/status", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Call handler
			handler.UpdateStatus(c)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("UpdateStatus() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			// Parse response
			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			// Check if error is expected
			if tt.expectError {
				if _, hasError := response["error"]; !hasError {
					t.Error("UpdateStatus() expected error but got none")
				}
			} else {
				if _, hasData := response["id"]; !hasData {
					t.Error("UpdateStatus() expected account response but got none")
				} else {
					if response["id"] != tt.accountID {
						t.Errorf("UpdateStatus() accountID = %v, want %v", response["id"], tt.accountID)
					}
				}
			}
		})
	}
}

func TestGetAccountsByUser(t *testing.T) {
	service := NewMockAccountService()
	handler := New(service)

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
		name           string
		userID         string
		expectedStatus int
		expectError    bool
		expectedCount  int
	}{
		{
			name:           "user with multiple accounts",
			userID:         userID,
			expectedStatus: http.StatusOK,
			expectError:    false,
			expectedCount:  2,
		},
		{
			name:           "user with single account",
			userID:         otherUserID,
			expectedStatus: http.StatusOK,
			expectError:    false,
			expectedCount:  1,
		},
		{
			name:           "user with no accounts",
			userID:         uuid.NewString(),
			expectedStatus: http.StatusOK,
			expectError:    false,
			expectedCount:  0,
		},
		{
			name:           "invalid user ID",
			userID:         "",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := createTestContext()

			// Set URL parameter
			c.Params = gin.Params{
				{Key: "userId", Value: tt.userID},
			}
			c.Request = httptest.NewRequest(http.MethodGet, "/api/accounts/user/"+tt.userID, nil)

			// Call handler
			handler.GetAccountsByUser(c)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("GetAccountsByUser() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			// Parse response
			var response interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			// Check if error is expected
			if tt.expectError {
				responseMap := response.(map[string]interface{})
				if _, hasError := responseMap["error"]; !hasError {
					t.Error("GetAccountsByUser() expected error but got none")
				}
			} else {
				// Handle empty array case
				if response == nil {
					if tt.expectedCount != 0 {
						t.Errorf("GetAccountsByUser() accounts count = 0, want %v", tt.expectedCount)
					}
				} else {
					responseData := response.([]interface{})
					if len(responseData) != tt.expectedCount {
						t.Errorf("GetAccountsByUser() accounts count = %v, want %v", len(responseData), tt.expectedCount)
					}

					// Verify all accounts belong to the user
					for _, accountInterface := range responseData {
						account := accountInterface.(map[string]interface{})
						if account["user_id"] != tt.userID && tt.userID != "" {
							t.Errorf("GetAccountsByUser() account belongs to user %v, want %v", account["user_id"], tt.userID)
						}
					}
				}
			}
		})
	}
}

func TestDeleteAccount(t *testing.T) {
	service := NewMockAccountService()
	handler := New(service)

	// Create accounts with different balances
	zeroBalanceAccount := &entities.Account{
		UserID:      uuid.NewString(),
		AccountType: entities.AccountTypeSavings,
		Name:        "Zero Balance Account",
		Currency:    entities.CurrencyUSD,
		Balance:     0.0,
		IsActive:    true,
	}
	createdZeroAccount, _ := service.CreateAccount(zeroBalanceAccount)

	activeAccount := &entities.Account{
		UserID:      uuid.NewString(),
		AccountType: entities.AccountTypeChecking,
		Name:        "Active Account",
		Currency:    entities.CurrencyUSD,
		Balance:     100.0,
		IsActive:    true,
	}
	createdActiveAccount, _ := service.CreateAccount(activeAccount)

	tests := []struct {
		name           string
		accountID      string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "delete account with zero balance",
			accountID:      createdZeroAccount.ID,
			expectedStatus: http.StatusNoContent,
			expectError:    false,
		},
		{
			name:           "cannot delete account with balance",
			accountID:      createdActiveAccount.ID,
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
		{
			name:           "non-existing account",
			accountID:      uuid.NewString(),
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := createTestContext()

			// Set URL parameter
			c.Params = gin.Params{
				{Key: "id", Value: tt.accountID},
			}
			c.Request = httptest.NewRequest(http.MethodDelete, "/api/accounts/"+tt.accountID, nil)

			// Call handler
			handler.DeleteAccount(c)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("DeleteAccount() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			// Check for error response if error is expected
			if tt.expectError && w.Body.Len() > 0 {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				if _, hasError := response["error"]; !hasError {
					t.Error("DeleteAccount() expected error but got none")
				}
			}
		})
	}
}

func TestAddFunds(t *testing.T) {
	service := NewMockAccountService()
	handler := New(service)

	// Create test wallet account
	account := &entities.Account{
		UserID:      uuid.NewString(),
		AccountType: entities.AccountTypeWallet,
		Name:        "Test Wallet",
		Currency:    entities.CurrencyUSD,
		Balance:     100.0,
		IsActive:    true,
	}
	createdAccount, _ := service.CreateAccount(account)

	tests := []struct {
		name           string
		accountID      string
		requestBody    dto.AddFundsRequest
		expectedStatus int
		expectError    bool
		expectedBalance float64
	}{
		{
			name:      "successful funds addition",
			accountID: createdAccount.ID,
			requestBody: dto.AddFundsRequest{
				Amount:      50.0,
				Description: "Test deposit",
				Reference:   "REF001",
			},
			expectedStatus:  http.StatusOK,
			expectError:     false,
			expectedBalance: 150.0,
		},
		{
			name:      "zero amount",
			accountID: createdAccount.ID,
			requestBody: dto.AddFundsRequest{
				Amount:      0.0,
				Description: "Test deposit",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:      "negative amount",
			accountID: createdAccount.ID,
			requestBody: dto.AddFundsRequest{
				Amount:      -10.0,
				Description: "Test deposit",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:      "missing description",
			accountID: createdAccount.ID,
			requestBody: dto.AddFundsRequest{
				Amount: 50.0,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:      "non-existing account",
			accountID: uuid.NewString(),
			requestBody: dto.AddFundsRequest{
				Amount:      50.0,
				Description: "Test deposit",
			},
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := createTestContext()

			// Set URL parameter
			c.Params = gin.Params{
				{Key: "id", Value: tt.accountID},
			}

			// Create request body
			body, _ := json.Marshal(tt.requestBody)
			c.Request = httptest.NewRequest(http.MethodPost, "/api/accounts/"+tt.accountID+"/add-funds", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Call handler
			handler.AddFunds(c)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("AddFunds() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			// Parse response
			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			// Check if error is expected
			if tt.expectError {
				if _, hasError := response["error"]; !hasError {
					t.Error("AddFunds() expected error but got none")
				}
			} else {
				if _, hasData := response["account_id"]; !hasData {
					t.Error("AddFunds() expected account_id in response")
				} else {
					if balance, ok := response["balance"].(float64); ok {
						if balance != tt.expectedBalance {
							t.Errorf("AddFunds() balance = %v, want %v", balance, tt.expectedBalance)
						}
					}
				}
			}
		})
	}
}

func TestWithdrawFunds(t *testing.T) {
	service := NewMockAccountService()
	handler := New(service)

	// Create test wallet account
	account := &entities.Account{
		UserID:      uuid.NewString(),
		AccountType: entities.AccountTypeWallet,
		Name:        "Test Wallet",
		Currency:    entities.CurrencyUSD,
		Balance:     100.0,
		IsActive:    true,
	}
	createdAccount, _ := service.CreateAccount(account)

	tests := []struct {
		name           string
		accountID      string
		requestBody    dto.WithdrawFundsRequest
		expectedStatus int
		expectError    bool
		expectedBalance float64
	}{
		{
			name:      "successful funds withdrawal",
			accountID: createdAccount.ID,
			requestBody: dto.WithdrawFundsRequest{
				Amount:      30.0,
				Description: "Test withdrawal",
				Reference:   "REF002",
			},
			expectedStatus:  http.StatusOK,
			expectError:     false,
			expectedBalance: 70.0,
		},
		{
			name:      "insufficient balance",
			accountID: createdAccount.ID,
			requestBody: dto.WithdrawFundsRequest{
				Amount:      200.0,
				Description: "Test withdrawal",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:      "zero amount",
			accountID: createdAccount.ID,
			requestBody: dto.WithdrawFundsRequest{
				Amount:      0.0,
				Description: "Test withdrawal",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:      "negative amount",
			accountID: createdAccount.ID,
			requestBody: dto.WithdrawFundsRequest{
				Amount:      -10.0,
				Description: "Test withdrawal",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:      "missing description",
			accountID: createdAccount.ID,
			requestBody: dto.WithdrawFundsRequest{
				Amount: 30.0,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:      "non-existing account",
			accountID: uuid.NewString(),
			requestBody: dto.WithdrawFundsRequest{
				Amount:      30.0,
				Description: "Test withdrawal",
			},
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := createTestContext()

			// Set URL parameter
			c.Params = gin.Params{
				{Key: "id", Value: tt.accountID},
			}

			// Create request body
			body, _ := json.Marshal(tt.requestBody)
			c.Request = httptest.NewRequest(http.MethodPost, "/api/accounts/"+tt.accountID+"/withdraw-funds", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Call handler
			handler.WithdrawFunds(c)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("WithdrawFunds() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			// Parse response
			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			// Check if error is expected
			if tt.expectError {
				if _, hasError := response["error"]; !hasError {
					t.Error("WithdrawFunds() expected error but got none")
				}
			} else {
				if _, hasData := response["account_id"]; !hasData {
					t.Error("WithdrawFunds() expected account_id in response")
				} else {
					if balance, ok := response["balance"].(float64); ok {
						if balance != tt.expectedBalance {
							t.Errorf("WithdrawFunds() balance = %v, want %v", balance, tt.expectedBalance)
						}
					}
				}
			}
		})
	}
}

func TestUpdateCreditLimit(t *testing.T) {
	service := NewMockAccountService()
	handler := New(service)

	// Create test credit card account
	creditLimit := 5000.0
	account := &entities.Account{
		UserID:      uuid.NewString(),
		AccountType: entities.AccountTypeCredit,
		Name:        "Test Credit Card",
		Currency:    entities.CurrencyUSD,
		Balance:     -500.0, // Used 500 of credit
		CreditLimit: &creditLimit,
		IsActive:    true,
	}
	createdAccount, _ := service.CreateAccount(account)

	tests := []struct {
		name           string
		accountID      string
		requestBody    dto.UpdateCreditLimitRequest
		expectedStatus int
		expectError    bool
		expectedLimit  float64
	}{
		{
			name:      "successful credit limit update",
			accountID: createdAccount.ID,
			requestBody: dto.UpdateCreditLimitRequest{
				CreditLimit: 7500.0,
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
			expectedLimit:  7500.0,
		},
		{
			name:      "zero credit limit",
			accountID: createdAccount.ID,
			requestBody: dto.UpdateCreditLimitRequest{
				CreditLimit: 0.0,
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
			expectedLimit:  0.0,
		},
		{
			name:      "negative credit limit",
			accountID: createdAccount.ID,
			requestBody: dto.UpdateCreditLimitRequest{
				CreditLimit: -1000.0,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:      "non-existing account",
			accountID: uuid.NewString(),
			requestBody: dto.UpdateCreditLimitRequest{
				CreditLimit: 5000.0,
			},
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := createTestContext()

			// Set URL parameter
			c.Params = gin.Params{
				{Key: "id", Value: tt.accountID},
			}

			// Create request body
			body, _ := json.Marshal(tt.requestBody)
			c.Request = httptest.NewRequest(http.MethodPut, "/api/accounts/"+tt.accountID+"/credit-limit", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Call handler
			handler.UpdateCreditLimit(c)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("UpdateCreditLimit() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			// Parse response
			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			// Check if error is expected
			if tt.expectError {
				if _, hasError := response["error"]; !hasError {
					t.Error("UpdateCreditLimit() expected error but got none")
				}
			} else {
				if _, hasData := response["id"]; !hasData {
					t.Error("UpdateCreditLimit() expected account data but got none")
				} else {
					if limit, ok := response["credit_limit"].(float64); ok {
						if limit != tt.expectedLimit {
							t.Errorf("UpdateCreditLimit() credit_limit = %v, want %v", limit, tt.expectedLimit)
						}
					}
				}
			}
		})
	}
}

func TestUpdateCreditDates(t *testing.T) {
	service := NewMockAccountService()
	handler := New(service)

	// Create test credit card account
	creditLimit := 5000.0
	account := &entities.Account{
		UserID:      uuid.NewString(),
		AccountType: entities.AccountTypeCredit,
		Name:        "Test Credit Card",
		Currency:    entities.CurrencyUSD,
		Balance:     -500.0,
		CreditLimit: &creditLimit,
		IsActive:    true,
	}
	createdAccount, _ := service.CreateAccount(account)

	closingDate := time.Now().AddDate(0, 1, 0)   // Next month
	dueDate := time.Now().AddDate(0, 1, 15)      // Next month + 15 days

	tests := []struct {
		name           string
		accountID      string
		requestBody    dto.UpdateCreditDatesRequest
		expectedStatus int
		expectError    bool
	}{
		{
			name:      "successful dates update",
			accountID: createdAccount.ID,
			requestBody: dto.UpdateCreditDatesRequest{
				ClosingDate: &closingDate,
				DueDate:     &dueDate,
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:      "update only closing date",
			accountID: createdAccount.ID,
			requestBody: dto.UpdateCreditDatesRequest{
				ClosingDate: &closingDate,
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:      "update only due date",
			accountID: createdAccount.ID,
			requestBody: dto.UpdateCreditDatesRequest{
				DueDate: &dueDate,
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "empty request body",
			accountID:      createdAccount.ID,
			requestBody:    dto.UpdateCreditDatesRequest{},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:      "non-existing account",
			accountID: uuid.NewString(),
			requestBody: dto.UpdateCreditDatesRequest{
				ClosingDate: &closingDate,
				DueDate:     &dueDate,
			},
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := createTestContext()

			// Set URL parameter
			c.Params = gin.Params{
				{Key: "id", Value: tt.accountID},
			}

			// Create request body
			body, _ := json.Marshal(tt.requestBody)
			c.Request = httptest.NewRequest(http.MethodPut, "/api/accounts/"+tt.accountID+"/credit-dates", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Call handler
			handler.UpdateCreditDates(c)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("UpdateCreditDates() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			// Parse response
			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			// Check if error is expected
			if tt.expectError {
				if _, hasError := response["error"]; !hasError {
					t.Error("UpdateCreditDates() expected error but got none")
				}
			} else {
				if _, hasData := response["id"]; !hasData {
					t.Error("UpdateCreditDates() expected account data but got none")
				}
			}
		})
	}
}

func TestGetAvailableCredit(t *testing.T) {
	service := NewMockAccountService()
	handler := New(service)

	// Create test credit card account with used credit
	creditLimit := 5000.0
	account := &entities.Account{
		UserID:      uuid.NewString(),
		AccountType: entities.AccountTypeCredit,
		Name:        "Test Credit Card",
		Currency:    entities.CurrencyUSD,
		Balance:     -1500.0, // Used 1500 of credit
		CreditLimit: &creditLimit,
		IsActive:    true,
	}
	createdAccount, _ := service.CreateAccount(account)

	// Create account without credit limit
	accountNoCreditLimit := &entities.Account{
		UserID:      uuid.NewString(),
		AccountType: entities.AccountTypeWallet,
		Name:        "Test Wallet",
		Currency:    entities.CurrencyUSD,
		Balance:     100.0,
		IsActive:    true,
	}
	createdWalletAccount, _ := service.CreateAccount(accountNoCreditLimit)

	tests := []struct {
		name             string
		accountID        string
		expectedStatus   int
		expectError      bool
		expectedCredit   float64
		expectedUsed     float64
		expectedAvailable float64
	}{
		{
			name:              "successful credit calculation",
			accountID:         createdAccount.ID,
			expectedStatus:    http.StatusOK,
			expectError:       false,
			expectedCredit:    5000.0,
			expectedUsed:      1500.0,
			expectedAvailable: 3500.0,
		},
		{
			name:           "account without credit limit",
			accountID:      createdWalletAccount.ID,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "non-existing account",
			accountID:      uuid.NewString(),
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := createTestContext()

			// Set URL parameter
			c.Params = gin.Params{
				{Key: "id", Value: tt.accountID},
			}
			c.Request = httptest.NewRequest(http.MethodGet, "/api/accounts/"+tt.accountID+"/available-credit", nil)

			// Call handler
			handler.GetAvailableCredit(c)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("GetAvailableCredit() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			// Parse response
			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			// Check if error is expected
			if tt.expectError {
				if _, hasError := response["error"]; !hasError {
					t.Error("GetAvailableCredit() expected error but got none")
				}
			} else {
				if accountID, ok := response["account_id"].(string); !ok || accountID != tt.accountID {
					t.Error("GetAvailableCredit() expected account_id in response")
				}
				if creditLimit, ok := response["credit_limit"].(float64); !ok || creditLimit != tt.expectedCredit {
					t.Errorf("GetAvailableCredit() credit_limit = %v, want %v", creditLimit, tt.expectedCredit)
				}
				if usedCredit, ok := response["used_credit"].(float64); !ok || usedCredit != tt.expectedUsed {
					t.Errorf("GetAvailableCredit() used_credit = %v, want %v", usedCredit, tt.expectedUsed)
				}
				if availableCredit, ok := response["available_credit"].(float64); !ok || availableCredit != tt.expectedAvailable {
					t.Errorf("GetAvailableCredit() available_credit = %v, want %v", availableCredit, tt.expectedAvailable)
				}
			}
		})
	}
}
