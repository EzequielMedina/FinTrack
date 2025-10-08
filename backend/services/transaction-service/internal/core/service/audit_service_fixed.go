package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	domaintransaction "github.com/fintrack/transaction-service/internal/core/domain/entities/transaction"
)

// TransactionAuditService implements TransactionAuditServiceInterface
// Handles audit logging for all transaction changes and activities
type TransactionAuditService struct {
	auditRepo TransactionAuditRepositoryInterface
}

// NewTransactionAuditService creates a new transaction audit service
func NewTransactionAuditService(auditRepo TransactionAuditRepositoryInterface) TransactionAuditServiceInterface {
	return &TransactionAuditService{
		auditRepo: auditRepo,
	}
}

// LogTransactionChange logs a change to a transaction for audit purposes
func (s *TransactionAuditService) LogTransactionChange(
	transactionID string,
	action string,
	oldStatus *domaintransaction.TransactionStatus,
	newStatus *domaintransaction.TransactionStatus,
	changedBy string,
	reason string,
) error {
	// Validate required parameters
	if transactionID == "" {
		return fmt.Errorf("transactionID cannot be empty")
	}
	if action == "" {
		return fmt.Errorf("action cannot be empty")
	}
	if changedBy == "" {
		return fmt.Errorf("changedBy cannot be empty")
	}

	// Create audit entry
	auditEntry := &TransactionAuditEntry{
		ID:            s.generateID(),
		TransactionID: transactionID,
		Action:        action,
		OldStatus:     oldStatus,
		NewStatus:     newStatus,
		ChangedBy:     changedBy,
		ChangeReason:  reason,
		CreatedAt:     time.Now(),
	}

	// Build changed fields map
	changedFields := make(map[string]interface{})
	if oldStatus != nil && newStatus != nil && *oldStatus != *newStatus {
		changedFields["status"] = map[string]interface{}{
			"old": string(*oldStatus),
			"new": string(*newStatus),
		}
	}
	auditEntry.ChangedFields = changedFields

	// Save audit entry
	if s.auditRepo == nil {
		return fmt.Errorf("audit repository is not initialized")
	}
	if err := s.auditRepo.Create(auditEntry); err != nil {
		return fmt.Errorf("failed to create audit entry: %w", err)
	}

	return nil
}

// GetAuditTrail retrieves the audit trail for a specific transaction
func (s *TransactionAuditService) GetAuditTrail(transactionID string) ([]*AuditEntry, error) {
	entries, err := s.auditRepo.GetByTransactionID(transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit trail: %w", err)
	}

	// Convert to response format
	auditEntries := make([]*AuditEntry, len(entries))
	for i, entry := range entries {
		auditEntries[i] = &AuditEntry{
			ID:            entry.ID,
			TransactionID: entry.TransactionID,
			Action:        entry.Action,
			OldStatus:     entry.OldStatus,
			NewStatus:     entry.NewStatus,
			ChangedFields: entry.ChangedFields,
			ChangedBy:     entry.ChangedBy,
			ChangeReason:  entry.ChangeReason,
			IPAddress:     entry.IPAddress,
			UserAgent:     entry.UserAgent,
			CreatedAt:     entry.CreatedAt,
		}
	}

	return auditEntries, nil
}

// GetUserAuditTrail retrieves audit entries for all transactions by a user within a date range
func (s *TransactionAuditService) GetUserAuditTrail(userID string, fromDate time.Time, toDate time.Time) ([]*AuditEntry, error) {
	entries, err := s.auditRepo.GetByUserID(userID, fromDate, toDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get user audit trail: %w", err)
	}

	// Convert to response format
	auditEntries := make([]*AuditEntry, len(entries))
	for i, entry := range entries {
		auditEntries[i] = &AuditEntry{
			ID:            entry.ID,
			TransactionID: entry.TransactionID,
			Action:        entry.Action,
			OldStatus:     entry.OldStatus,
			NewStatus:     entry.NewStatus,
			ChangedFields: entry.ChangedFields,
			ChangedBy:     entry.ChangedBy,
			ChangeReason:  entry.ChangeReason,
			IPAddress:     entry.IPAddress,
			UserAgent:     entry.UserAgent,
			CreatedAt:     entry.CreatedAt,
		}
	}

	return auditEntries, nil
}

// generateID creates a simple UUID-like string for audit entries
func (s *TransactionAuditService) generateID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based ID
		return fmt.Sprintf("audit_%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}
