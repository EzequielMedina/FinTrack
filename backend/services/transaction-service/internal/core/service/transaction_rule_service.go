package service

import (
	"fmt"
	"time"

	domaintransaction "github.com/fintrack/transaction-service/internal/core/domain/entities/transaction"
)

// TransactionRuleService implements TransactionRuleServiceInterface
// Handles business rules validation and management for transactions
type TransactionRuleService struct {
	ruleRepo TransactionRuleRepositoryInterface
}

// NewTransactionRuleService creates a new transaction rule service
func NewTransactionRuleService(ruleRepo TransactionRuleRepositoryInterface) TransactionRuleServiceInterface {
	return &TransactionRuleService{
		ruleRepo: ruleRepo,
	}
}

// GetRulesForTransaction retrieves applicable rules for a transaction
func (s *TransactionRuleService) GetRulesForTransaction(userID string, accountID *string, cardID *string, transactionType domaintransaction.TransactionType) (*domaintransaction.TransactionRule, error) {
	rules, err := s.ruleRepo.GetActiveRulesForTransaction(userID, accountID, cardID, transactionType)
	if err != nil {
		return nil, fmt.Errorf("failed to get rules: %w", err)
	}

	// If no specific rules found, return default rules
	if len(rules) == 0 {
		return s.getDefaultRules(), nil
	}

	// Merge multiple rules into one (taking the most restrictive values)
	mergedRule := s.mergeRules(rules)
	return mergedRule, nil
}

// ValidateTransactionAgainstRules validates a transaction against applicable business rules
func (s *TransactionRuleService) ValidateTransactionAgainstRules(transaction *domaintransaction.Transaction) error {
	rules, err := s.GetRulesForTransaction(
		transaction.UserID,
		transaction.FromAccountID,
		transaction.FromCardID,
		transaction.Type,
	)
	if err != nil {
		return fmt.Errorf("failed to get rules for validation: %w", err)
	}

	// Check single transaction amount limit
	if rules.MaxSingleAmount > 0 && transaction.Amount > rules.MaxSingleAmount {
		return fmt.Errorf("transaction amount %.2f exceeds maximum allowed %.2f",
			transaction.Amount, rules.MaxSingleAmount)
	}

	// Check if transaction requires approval
	if rules.RequiresApproval {
		// For now, we'll just mark it as requiring approval
		// In a real system, this would trigger an approval workflow
		transaction.Status = domaintransaction.TransactionStatusPending
	}

	// Validate allowed hours (simplified)
	if rules.AllowedHours != "" {
		currentHour := time.Now().Hour()
		if !s.isHourAllowed(currentHour, rules.AllowedHours) {
			return fmt.Errorf("transaction not allowed at current hour")
		}
	}

	return nil
}

// CreateRule creates a new transaction rule
func (s *TransactionRuleService) CreateRule(userID string, rule CreateRuleRequest, createdBy string) error {
	// Create domain entity
	domainRule := &domaintransaction.TransactionRule{
		MaxDailyAmount:   *rule.MaxDailyAmount,
		MaxSingleAmount:  *rule.MaxSingleAmount,
		RequiresApproval: rule.RequiresApproval,
		AllowedHours:     rule.AllowedHours,
	}

	// Save to repository (simplified - would need full implementation)
	_, err := s.ruleRepo.Create(domainRule)
	return err
}

// UpdateRule updates an existing rule
func (s *TransactionRuleService) UpdateRule(ruleID string, updates map[string]interface{}, updatedBy string) error {
	rule, err := s.ruleRepo.GetByID(ruleID)
	if err != nil {
		return fmt.Errorf("failed to get rule: %w", err)
	}

	// Apply updates (simplified)
	if maxDaily, exists := updates["maxDailyAmount"]; exists {
		if amount, ok := maxDaily.(float64); ok {
			rule.MaxDailyAmount = amount
		}
	}

	if maxSingle, exists := updates["maxSingleAmount"]; exists {
		if amount, ok := maxSingle.(float64); ok {
			rule.MaxSingleAmount = amount
		}
	}

	_, err = s.ruleRepo.Update(rule)
	return err
}

// DeleteRule deletes a rule
func (s *TransactionRuleService) DeleteRule(ruleID string, deletedBy string) error {
	return s.ruleRepo.Delete(ruleID)
}

// Helper methods

// getDefaultRules returns default transaction rules
func (s *TransactionRuleService) getDefaultRules() *domaintransaction.TransactionRule {
	return &domaintransaction.TransactionRule{
		MaxDailyAmount:   10000.00, // Default daily limit
		MaxSingleAmount:  1000.00,  // Default single transaction limit
		RequiresApproval: false,    // Default no approval required
		AllowedHours:     "0-23",   // Default 24/7
	}
}

// mergeRules merges multiple rules into one, taking the most restrictive values
func (s *TransactionRuleService) mergeRules(rules []*domaintransaction.TransactionRule) *domaintransaction.TransactionRule {
	if len(rules) == 0 {
		return s.getDefaultRules()
	}

	merged := &domaintransaction.TransactionRule{
		MaxDailyAmount:   rules[0].MaxDailyAmount,
		MaxSingleAmount:  rules[0].MaxSingleAmount,
		RequiresApproval: rules[0].RequiresApproval,
		AllowedHours:     rules[0].AllowedHours,
	}

	for _, rule := range rules[1:] {
		// Take the smaller amount limits (more restrictive)
		if rule.MaxDailyAmount > 0 && rule.MaxDailyAmount < merged.MaxDailyAmount {
			merged.MaxDailyAmount = rule.MaxDailyAmount
		}
		if rule.MaxSingleAmount > 0 && rule.MaxSingleAmount < merged.MaxSingleAmount {
			merged.MaxSingleAmount = rule.MaxSingleAmount
		}

		// If any rule requires approval, merged rule requires approval
		if rule.RequiresApproval {
			merged.RequiresApproval = true
		}
	}

	return merged
}

// isHourAllowed checks if the current hour is within allowed hours (simplified)
func (s *TransactionRuleService) isHourAllowed(currentHour int, allowedHours string) bool {
	// Simplified implementation - in real system would parse complex hour ranges
	return allowedHours == "0-23" || allowedHours == ""
}
