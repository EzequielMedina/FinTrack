package service

import (
	"time"

	domaintransaction "github.com/fintrack/transaction-service/internal/core/domain/entities/transaction"
)

// TransactionServiceInterface defines the contract for transaction service operations
// This interface follows the Interface Segregation Principle (ISP) from SOLID
type TransactionServiceInterface interface {
	// Transaction CRUD operations
	CreateTransaction(request CreateTransactionRequest, initiatedBy string) (*domaintransaction.Transaction, error)
	GetTransactionByID(id string, userID string) (*domaintransaction.Transaction, error)
	GetTransactionsByUser(userID string, filters TransactionFilters) ([]*domaintransaction.Transaction, int, error)
	UpdateTransactionStatus(id string, status domaintransaction.TransactionStatus, reason string, updatedBy string) (*domaintransaction.Transaction, error)

	// Transaction processing operations
	ProcessTransaction(id string, processedBy string) error
	CompleteTransaction(id string, completedBy string) error
	FailTransaction(id string, reason string, failedBy string) error
	CancelTransaction(id string, reason string, canceledBy string) error
	ReverseTransaction(id string, reason string, reversedBy string) (*domaintransaction.Transaction, error)

	// Balance and account operations
	ProcessWalletDeposit(userID string, accountID string, amount float64, description string, initiatedBy string) (*domaintransaction.Transaction, error)
	ProcessWalletWithdrawal(userID string, accountID string, amount float64, description string, initiatedBy string) (*domaintransaction.Transaction, error)
	ProcessWalletTransfer(userID string, fromAccountID string, toAccountID string, amount float64, description string, initiatedBy string) (*domaintransaction.Transaction, error)

	// Card operations
	ProcessCreditCardCharge(userID string, cardID string, amount float64, description string, merchantName string, initiatedBy string) (*domaintransaction.Transaction, error)
	ProcessCreditCardPayment(userID string, cardID string, amount float64, paymentMethod domaintransaction.PaymentMethod, initiatedBy string) (*domaintransaction.Transaction, error)
	ProcessDebitCardPurchase(userID string, cardID string, amount float64, description string, merchantName string, initiatedBy string) (*domaintransaction.Transaction, error)

	// Account operations
	ProcessAccountTransfer(userID string, fromAccountID string, toAccountID string, amount float64, description string, initiatedBy string) (*domaintransaction.Transaction, error)
	ProcessAccountDeposit(userID string, accountID string, amount float64, description string, initiatedBy string) (*domaintransaction.Transaction, error)
	ProcessAccountWithdraw(userID string, accountID string, amount float64, description string, initiatedBy string) (*domaintransaction.Transaction, error)
}

// TransactionRuleServiceInterface defines the contract for transaction rules management
// Separated for better adherence to Single Responsibility Principle (SRP)
type TransactionRuleServiceInterface interface {
	GetRulesForTransaction(userID string, accountID *string, cardID *string, transactionType domaintransaction.TransactionType) (*domaintransaction.TransactionRule, error)
	ValidateTransactionAgainstRules(transaction *domaintransaction.Transaction) error
	CreateRule(userID string, rule CreateRuleRequest, createdBy string) error
	UpdateRule(ruleID string, updates map[string]interface{}, updatedBy string) error
	DeleteRule(ruleID string, deletedBy string) error
}

// TransactionAuditServiceInterface defines the contract for audit operations
// Separated for better adherence to Single Responsibility Principle (SRP)
type TransactionAuditServiceInterface interface {
	LogTransactionChange(transactionID string, action string, oldStatus *domaintransaction.TransactionStatus, newStatus *domaintransaction.TransactionStatus, changedBy string, reason string) error
	GetAuditTrail(transactionID string) ([]*AuditEntry, error)
	GetUserAuditTrail(userID string, fromDate time.Time, toDate time.Time) ([]*AuditEntry, error)
}

// ExternalServiceInterface defines the contract for communicating with other microservices
// This abstraction allows for easy testing and different implementations (Direct HTTP, Message Queue, etc.)
type ExternalServiceInterface interface {
	// Account service integration
	GetAccountBalance(accountID string) (float64, error)
	UpdateAccountBalance(accountID string, newBalance float64) error
	ValidateAccount(accountID string, userID string) error

	// Card service integration
	GetCardDetails(cardID string) (*CardInfo, error)
	ValidateCard(cardID string, userID string) error
	UpdateCardBalance(cardID string, newBalance float64) error

	// User service integration
	ValidateUser(userID string) error
	GetUserLimits(userID string) (*UserLimits, error)

	// Notification service integration
	SendTransactionNotification(userID string, transaction *domaintransaction.Transaction) error
}

// Request DTOs for service operations

// CreateTransactionRequest represents the data needed to create a transaction
type CreateTransactionRequest struct {
	UserID        string                            `json:"userId"`
	Type          domaintransaction.TransactionType `json:"type"`
	Amount        float64                           `json:"amount"`
	Currency      string                            `json:"currency"`
	FromAccountID *string                           `json:"fromAccountId"`
	ToAccountID   *string                           `json:"toAccountId"`
	FromCardID    *string                           `json:"fromCardId"`
	ToCardID      *string                           `json:"toCardId"`
	Description   string                            `json:"description"`
	PaymentMethod domaintransaction.PaymentMethod   `json:"paymentMethod"`
	MerchantName  string                            `json:"merchantName"`
	MerchantID    string                            `json:"merchantId"`
	ReferenceID   string                            `json:"referenceId"`
	ExternalID    string                            `json:"externalId"`
	Metadata      map[string]interface{}            `json:"metadata"`
	Tags          []string                          `json:"tags"`
}

// CreateRuleRequest represents the data needed to create a transaction rule
type CreateRuleRequest struct {
	AccountID        *string                           `json:"accountId"`
	CardID           *string                           `json:"cardId"`
	TransactionType  domaintransaction.TransactionType `json:"transactionType"`
	MaxDailyAmount   *float64                          `json:"maxDailyAmount"`
	MaxSingleAmount  *float64                          `json:"maxSingleAmount"`
	RequiresApproval bool                              `json:"requiresApproval"`
	AllowedHours     string                            `json:"allowedHours"`
	EffectiveFrom    *time.Time                        `json:"effectiveFrom"`
	EffectiveUntil   *time.Time                        `json:"effectiveUntil"`
}

// TransactionFilters represents filters for querying transactions
type TransactionFilters struct {
	Types         []domaintransaction.TransactionType   `json:"types"`
	Statuses      []domaintransaction.TransactionStatus `json:"statuses"`
	FromDate      *time.Time                            `json:"fromDate"`
	ToDate        *time.Time                            `json:"toDate"`
	MinAmount     *float64                              `json:"minAmount"`
	MaxAmount     *float64                              `json:"maxAmount"`
	AccountID     *string                               `json:"accountId"`
	CardID        *string                               `json:"cardId"`
	MerchantName  *string                               `json:"merchantName"`
	PaymentMethod *domaintransaction.PaymentMethod      `json:"paymentMethod"`
	Limit         int                                   `json:"limit"`
	Offset        int                                   `json:"offset"`
	OrderBy       string                                `json:"orderBy"`
	Order         string                                `json:"order"`
}

// Response DTOs for service operations

// AuditEntry represents an audit log entry
type AuditEntry struct {
	ID            string                               `json:"id"`
	TransactionID string                               `json:"transactionId"`
	Action        string                               `json:"action"`
	OldStatus     *domaintransaction.TransactionStatus `json:"oldStatus"`
	NewStatus     *domaintransaction.TransactionStatus `json:"newStatus"`
	ChangedFields map[string]interface{}               `json:"changedFields"`
	ChangedBy     string                               `json:"changedBy"`
	ChangeReason  string                               `json:"changeReason"`
	IPAddress     string                               `json:"ipAddress"`
	UserAgent     string                               `json:"userAgent"`
	CreatedAt     time.Time                            `json:"createdAt"`
}

// CardInfo represents card information from external service
type CardInfo struct {
	ID          string   `json:"id"`
	CardType    string   `json:"cardType"`
	Balance     float64  `json:"balance"`
	CreditLimit *float64 `json:"creditLimit"`
	IsActive    bool     `json:"isActive"`
	AccountID   string   `json:"accountId"`
}

// UserLimits represents user transaction limits
type UserLimits struct {
	DailyTransactionLimit   float64 `json:"dailyTransactionLimit"`
	MonthlyTransactionLimit float64 `json:"monthlyTransactionLimit"`
	SingleTransactionLimit  float64 `json:"singleTransactionLimit"`
	RequiresApprovalAbove   float64 `json:"requiresApprovalAbove"`
}
