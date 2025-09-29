package service

import (
	"time"

	domaintransaction "github.com/fintrack/transaction-service/internal/core/domain/entities/transaction"
)

// TransactionRepositoryInterface defines the contract for transaction data access
// This follows the Dependency Inversion Principle (DIP) from SOLID
type TransactionRepositoryInterface interface {
	// Basic CRUD operations
	Create(transaction *domaintransaction.Transaction) (*domaintransaction.Transaction, error)
	GetByID(id string) (*domaintransaction.Transaction, error)
	Update(transaction *domaintransaction.Transaction) (*domaintransaction.Transaction, error)
	Delete(id string) error

	// Query operations
	GetByUserID(userID string, filters TransactionFilters) ([]*domaintransaction.Transaction, int, error)
	GetByAccountID(accountID string, filters TransactionFilters) ([]*domaintransaction.Transaction, int, error)
	GetByCardID(cardID string, filters TransactionFilters) ([]*domaintransaction.Transaction, int, error)
	GetByReferenceID(referenceID string) (*domaintransaction.Transaction, error)
	GetByExternalID(externalID string) (*domaintransaction.Transaction, error)

	// Batch operations
	CreateBatch(transactions []*domaintransaction.Transaction) ([]*domaintransaction.Transaction, error)
	UpdateBatch(transactions []*domaintransaction.Transaction) ([]*domaintransaction.Transaction, error)

	// Aggregation operations
	GetUserTransactionSummary(userID string, fromDate, toDate *time.Time) (*TransactionSummary, error)
	GetAccountTransactionSummary(accountID string, fromDate, toDate *time.Time) (*TransactionSummary, error)
	GetDailyTransactionVolume(userID string, date time.Time) (float64, error)
	GetMonthlyTransactionVolume(userID string, year int, month int) (float64, error)
}

// TransactionRuleRepositoryInterface defines the contract for transaction rules data access
type TransactionRuleRepositoryInterface interface {
	Create(rule *domaintransaction.TransactionRule) (*domaintransaction.TransactionRule, error)
	GetByID(id string) (*domaintransaction.TransactionRule, error)
	Update(rule *domaintransaction.TransactionRule) (*domaintransaction.TransactionRule, error)
	Delete(id string) error

	GetByUserID(userID string) ([]*domaintransaction.TransactionRule, error)
	GetByAccountID(accountID string) ([]*domaintransaction.TransactionRule, error)
	GetByCardID(cardID string) ([]*domaintransaction.TransactionRule, error)
	GetActiveRulesForTransaction(userID string, accountID *string, cardID *string, transactionType domaintransaction.TransactionType) ([]*domaintransaction.TransactionRule, error)
}

// TransactionAuditRepositoryInterface defines the contract for audit data access
type TransactionAuditRepositoryInterface interface {
	Create(audit *TransactionAuditEntry) error
	GetByTransactionID(transactionID string) ([]*TransactionAuditEntry, error)
	GetByUserID(userID string, fromDate, toDate time.Time) ([]*TransactionAuditEntry, error)
	GetByDateRange(fromDate, toDate time.Time, limit int) ([]*TransactionAuditEntry, error)
}

// Supporting types for repository operations

// TransactionSummary represents aggregated transaction data
type TransactionSummary struct {
	TotalAmount      float64                                       `json:"totalAmount"`
	TransactionCount int                                           `json:"transactionCount"`
	ByType           map[domaintransaction.TransactionType]float64 `json:"byType"`
	ByStatus         map[domaintransaction.TransactionStatus]int   `json:"byStatus"`
	AverageAmount    float64                                       `json:"averageAmount"`
	MaxAmount        float64                                       `json:"maxAmount"`
	MinAmount        float64                                       `json:"minAmount"`
}

// TransactionAuditEntry represents an audit log entry in the database
type TransactionAuditEntry struct {
	ID            string                               `json:"id" db:"id"`
	TransactionID string                               `json:"transactionId" db:"transaction_id"`
	Action        string                               `json:"action" db:"action"`
	OldStatus     *domaintransaction.TransactionStatus `json:"oldStatus" db:"old_status"`
	NewStatus     *domaintransaction.TransactionStatus `json:"newStatus" db:"new_status"`
	ChangedFields map[string]interface{}               `json:"changedFields" db:"changed_fields"`
	ChangedBy     string                               `json:"changedBy" db:"changed_by"`
	ChangeReason  string                               `json:"changeReason" db:"change_reason"`
	IPAddress     string                               `json:"ipAddress" db:"ip_address"`
	UserAgent     string                               `json:"userAgent" db:"user_agent"`
	CreatedAt     time.Time                            `json:"createdAt" db:"created_at"`
}
