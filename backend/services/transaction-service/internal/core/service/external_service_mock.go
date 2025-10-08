package service

import (
	"fmt"

	domaintransaction "github.com/fintrack/transaction-service/internal/core/domain/entities/transaction"
)

// MockExternalService is a mock implementation for testing and development
// In production, this would be replaced with actual HTTP clients or message queue handlers
type MockExternalService struct {
	// Mock data storage for development/testing
	accounts map[string]float64 // accountID -> balance
	cards    map[string]*CardInfo
	users    map[string]bool // userID -> exists
}

// NewMockExternalService creates a new mock external service
func NewMockExternalService() ExternalServiceInterface {
	return &MockExternalService{
		accounts: make(map[string]float64),
		cards:    make(map[string]*CardInfo),
		users:    make(map[string]bool),
	}
}

// Account service integration methods

// GetAccountBalance retrieves the current balance of an account
func (s *MockExternalService) GetAccountBalance(accountID string) (float64, error) {
	if balance, exists := s.accounts[accountID]; exists {
		return balance, nil
	}
	// Return default balance if account doesn't exist in mock
	s.accounts[accountID] = 1000.0 // Default mock balance
	return 1000.0, nil
}

// UpdateAccountBalance updates the balance of an account
func (s *MockExternalService) UpdateAccountBalance(accountID string, newBalance float64) error {
	if newBalance < 0 {
		return fmt.Errorf("account balance cannot be negative")
	}
	s.accounts[accountID] = newBalance
	return nil
}

// ValidateAccount checks if an account exists and belongs to the user
func (s *MockExternalService) ValidateAccount(accountID string, userID string) error {
	// Mock validation - in real implementation would call account-service
	if accountID == "" {
		return fmt.Errorf("account ID cannot be empty")
	}
	if userID == "" {
		return fmt.Errorf("user ID cannot be empty")
	}
	// Mock: assume all accounts are valid for simplicity
	return nil
}

// Card service integration methods

// GetCardDetails retrieves card information
func (s *MockExternalService) GetCardDetails(cardID string) (*CardInfo, error) {
	if card, exists := s.cards[cardID]; exists {
		return card, nil
	}

	// Return mock card info if not exists
	mockCard := &CardInfo{
		ID:          cardID,
		CardType:    "credit",
		Balance:     500.0,
		CreditLimit: &[]float64{2000.0}[0], // Pointer to 2000.0
		IsActive:    true,
		AccountID:   "mock_account_" + cardID,
	}
	s.cards[cardID] = mockCard
	return mockCard, nil
}

// ValidateCard checks if a card exists and belongs to the user
func (s *MockExternalService) ValidateCard(cardID string, userID string) error {
	// Mock validation - in real implementation would call account-service or card-service
	if cardID == "" {
		return fmt.Errorf("card ID cannot be empty")
	}
	if userID == "" {
		return fmt.Errorf("user ID cannot be empty")
	}
	// Mock: assume all cards are valid for simplicity
	return nil
}

// UpdateCardBalance updates the balance of a card
func (s *MockExternalService) UpdateCardBalance(cardID string, newBalance float64) error {
	card, exists := s.cards[cardID]
	if !exists {
		// Create mock card if doesn't exist
		card = &CardInfo{
			ID:          cardID,
			CardType:    "credit",
			Balance:     newBalance,
			CreditLimit: &[]float64{2000.0}[0],
			IsActive:    true,
			AccountID:   "mock_account_" + cardID,
		}
		s.cards[cardID] = card
		return nil
	}

	card.Balance = newBalance
	return nil
}

// User service integration methods

// ValidateUser checks if a user exists
func (s *MockExternalService) ValidateUser(userID string) error {
	if userID == "" {
		return fmt.Errorf("user ID cannot be empty")
	}

	// Mock validation - assume all users are valid for simplicity
	s.users[userID] = true
	return nil
}

// GetUserLimits retrieves transaction limits for a user
func (s *MockExternalService) GetUserLimits(userID string) (*UserLimits, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	// Return mock user limits
	return &UserLimits{
		DailyTransactionLimit:   5000.0,
		MonthlyTransactionLimit: 50000.0,
		SingleTransactionLimit:  2000.0,
		RequiresApprovalAbove:   1000.0,
	}, nil
}

// Notification service integration methods

// SendTransactionNotification sends a notification about a transaction
func (s *MockExternalService) SendTransactionNotification(userID string, transaction *domaintransaction.Transaction) error {
	// Mock notification - in real implementation would call notification-service
	fmt.Printf("Mock Notification: Transaction %s for user %s - Amount: %.2f %s\n",
		transaction.ID, userID, transaction.Amount, transaction.Currency)
	return nil
}

// Helper methods for testing/development

// SetMockAccountBalance sets a mock account balance for testing
func (s *MockExternalService) SetMockAccountBalance(accountID string, balance float64) {
	s.accounts[accountID] = balance
}

// SetMockCard sets mock card info for testing
func (s *MockExternalService) SetMockCard(cardID string, cardInfo *CardInfo) {
	s.cards[cardID] = cardInfo
}
