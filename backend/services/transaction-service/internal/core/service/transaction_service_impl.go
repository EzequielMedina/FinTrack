package service

import (
	"errors"
	"fmt"

	domaintransaction "github.com/fintrack/transaction-service/internal/core/domain/entities/transaction"
	"github.com/fintrack/transaction-service/internal/core/interfaces"
)

// TransactionService implements TransactionServiceInterface following SOLID principles
// This is the main business logic service that orchestrates transaction operations
type TransactionService struct {
	transactionRepo TransactionRepositoryInterface
	ruleService     TransactionRuleServiceInterface
	auditService    TransactionAuditServiceInterface
	externalService ExternalServiceInterface
	accountService  interfaces.AccountServiceInterface
}

// NewTransactionService creates a new transaction service instance
// Dependency injection is used to promote testability and loose coupling (Dependency Inversion Principle)
func NewTransactionService(
	transactionRepo TransactionRepositoryInterface,
	ruleService TransactionRuleServiceInterface,
	auditService TransactionAuditServiceInterface,
	externalService ExternalServiceInterface,
	accountService interfaces.AccountServiceInterface,
) TransactionServiceInterface {
	return &TransactionService{
		transactionRepo: transactionRepo,
		ruleService:     ruleService,
		auditService:    auditService,
		externalService: externalService,
		accountService:  accountService,
	}
}

// CreateTransaction creates a new transaction with proper validation and business rules
func (s *TransactionService) CreateTransaction(request CreateTransactionRequest, initiatedBy string) (*domaintransaction.Transaction, error) {
	// Validate user exists
	if err := s.externalService.ValidateUser(request.UserID); err != nil {
		return nil, fmt.Errorf("user validation failed: %w", err)
	}

	// Create transaction entity in PENDING status
	transaction := &domaintransaction.Transaction{
		Type:          request.Type,
		UserID:        request.UserID,
		Amount:        request.Amount,
		Currency:      request.Currency,
		FromAccountID: request.FromAccountID,
		ToAccountID:   request.ToAccountID,
		FromCardID:    request.FromCardID,
		ToCardID:      request.ToCardID,
		InitiatedBy:   initiatedBy,
		Status:        domaintransaction.TransactionStatusPending,
		Description:   request.Description,
		PaymentMethod: request.PaymentMethod,
		MerchantName:  request.MerchantName,
		MerchantID:    request.MerchantID,
		ReferenceID:   request.ReferenceID,
		ExternalID:    request.ExternalID,
		Metadata:      request.Metadata,
		Tags:          request.Tags,
	}

	// Validate the transaction
	if err := transaction.Validate(); err != nil {
		return nil, fmt.Errorf("transaction validation failed: %w", err)
	}

	// Perform pre-transaction validations based on transaction type
	if err := s.performPreTransactionValidations(transaction); err != nil {
		transaction.Status = domaintransaction.TransactionStatusFailed
		transaction.FailureReason = err.Error()
		// Save the failed transaction for audit purposes
		s.transactionRepo.Create(transaction)
		return nil, fmt.Errorf("pre-transaction validation failed: %w", err)
	}

	// Save transaction in PENDING status
	savedTransaction, err := s.transactionRepo.Create(transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Check if this is a record-only transaction (balance already updated by another service)
	recordOnly := false
	if savedTransaction.Metadata != nil {
		if recordOnlyValue, exists := savedTransaction.Metadata["recordOnly"]; exists {
			// Handle both boolean and string values
			switch v := recordOnlyValue.(type) {
			case bool:
				recordOnly = v
			case string:
				recordOnly = v == "true"
			}
		}
	}

	// Execute the transaction (update balances) only if not record-only
	if !recordOnly {
		if err := s.executeTransaction(savedTransaction); err != nil {
			// Mark transaction as failed
			savedTransaction.Status = domaintransaction.TransactionStatusFailed
			savedTransaction.FailureReason = err.Error()
			s.transactionRepo.Update(savedTransaction)

			return nil, fmt.Errorf("transaction execution failed: %w", err)
		}
	}

	// Mark transaction as completed
	savedTransaction.Status = domaintransaction.TransactionStatusCompleted
	updatedTransaction, err := s.transactionRepo.Update(savedTransaction)
	if err != nil {
		return nil, fmt.Errorf("failed to update transaction status: %w", err)
	}

	// Log the transaction for audit
	oldStatus := domaintransaction.TransactionStatusPending
	newStatus := domaintransaction.TransactionStatusCompleted
	if err := s.auditService.LogTransactionChange(
		updatedTransaction.ID,
		"complete_transaction",
		&oldStatus,
		&newStatus,
		savedTransaction.InitiatedBy,
		"Transaction completed successfully",
	); err != nil {
		// Don't fail the transaction for audit logging errors, just log it
		fmt.Printf("Warning: Failed to log transaction %s for audit: %v\n", updatedTransaction.ID, err)
	}

	return updatedTransaction, nil
}

// performPreTransactionValidations validates that the transaction can be executed
func (s *TransactionService) performPreTransactionValidations(transaction *domaintransaction.Transaction) error {
	switch transaction.Type {
	case domaintransaction.TransactionTypeWalletWithdrawal, domaintransaction.TransactionTypeDebitWithdrawal, domaintransaction.TransactionTypeAccountWithdraw:
		return s.validateWithdrawal(transaction)
	case domaintransaction.TransactionTypeWalletTransfer, domaintransaction.TransactionTypeAccountTransfer:
		return s.validateTransfer(transaction)
	case domaintransaction.TransactionTypeDebitPurchase, domaintransaction.TransactionTypeCreditCharge:
		return s.validatePurchaseOrPayment(transaction)
	case domaintransaction.TransactionTypeWalletDeposit, domaintransaction.TransactionTypeAccountDeposit:
		return s.validateDeposit(transaction)
	default:
		return nil // No special validation for other types
	}
}

// Helper functions for safe pointer handling
func stringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func isEmpty(s *string) bool {
	return s == nil || *s == ""
}

// validateWithdrawal validates that a withdrawal can be performed
func (s *TransactionService) validateWithdrawal(transaction *domaintransaction.Transaction) error {
	if isEmpty(transaction.FromAccountID) {
		return errors.New("source account required for withdrawal")
	}

	// Check if account exists and is active
	exists, err := s.accountService.ValidateAccountExists(stringValue(transaction.FromAccountID))
	if err != nil {
		return fmt.Errorf("failed to validate account: %w", err)
	}
	if !exists {
		return errors.New("source account not found or inactive")
	}

	// Get current balance
	balance, err := s.accountService.GetAccountBalance(stringValue(transaction.FromAccountID))
	if err != nil {
		return fmt.Errorf("failed to get account balance: %w", err)
	}

	// Check sufficient funds
	if balance.Balance < transaction.Amount {
		return fmt.Errorf("insufficient funds: available %.2f, requested %.2f", balance.Balance, transaction.Amount)
	}

	return nil
}

// validateTransfer validates that a transfer can be performed
func (s *TransactionService) validateTransfer(transaction *domaintransaction.Transaction) error {
	if isEmpty(transaction.FromAccountID) || isEmpty(transaction.ToAccountID) {
		return errors.New("both source and destination accounts required for transfer")
	}

	if stringValue(transaction.FromAccountID) == stringValue(transaction.ToAccountID) {
		return errors.New("cannot transfer to the same account")
	}

	// Validate source account
	if err := s.validateWithdrawal(transaction); err != nil {
		return fmt.Errorf("source account validation failed: %w", err)
	}

	// Validate destination account exists
	exists, err := s.accountService.ValidateAccountExists(stringValue(transaction.ToAccountID))
	if err != nil {
		return fmt.Errorf("failed to validate destination account: %w", err)
	}
	if !exists {
		return errors.New("destination account not found or inactive")
	}

	return nil
}

// validatePurchaseOrPayment validates credit card transactions
func (s *TransactionService) validatePurchaseOrPayment(transaction *domaintransaction.Transaction) error {
	if isEmpty(transaction.FromAccountID) {
		return errors.New("account required for purchase/payment")
	}

	// Get account info to check if it's a credit account
	accountInfo, err := s.accountService.GetAccountInfo(stringValue(transaction.FromAccountID))
	if err != nil {
		return fmt.Errorf("failed to get account info: %w", err)
	}

	if !accountInfo.IsActive {
		return errors.New("account is inactive")
	}

	// If it's a credit account, check available credit
	if accountInfo.AccountType == "credit" {
		availableCredit, err := s.accountService.GetAvailableCredit(stringValue(transaction.FromAccountID))
		if err != nil {
			return fmt.Errorf("failed to get available credit: %w", err)
		}

		if availableCredit.Balance < transaction.Amount {
			return fmt.Errorf("insufficient credit: available %.2f, requested %.2f", availableCredit.Balance, transaction.Amount)
		}
	} else {
		// For debit accounts, check balance like a withdrawal
		return s.validateWithdrawal(transaction)
	}

	return nil
}

// validateDeposit validates that a deposit can be performed
func (s *TransactionService) validateDeposit(transaction *domaintransaction.Transaction) error {
	if isEmpty(transaction.ToAccountID) {
		return errors.New("destination account required for deposit")
	}

	// Check if account exists and is active
	exists, err := s.accountService.ValidateAccountExists(stringValue(transaction.ToAccountID))
	if err != nil {
		return fmt.Errorf("failed to validate account: %w", err)
	}
	if !exists {
		return errors.New("destination account not found or inactive")
	}

	return nil
}

// executeTransaction performs the actual balance updates
func (s *TransactionService) executeTransaction(transaction *domaintransaction.Transaction) error {
	switch transaction.Type {
	case domaintransaction.TransactionTypeWalletDeposit, domaintransaction.TransactionTypeAccountDeposit:
		return s.executeDeposit(transaction)
	case domaintransaction.TransactionTypeWalletWithdrawal, domaintransaction.TransactionTypeAccountWithdraw:
		return s.executeWithdrawal(transaction)
	case domaintransaction.TransactionTypeWalletTransfer, domaintransaction.TransactionTypeAccountTransfer:
		return s.executeTransfer(transaction)
	case domaintransaction.TransactionTypeCreditCharge, domaintransaction.TransactionTypeDebitPurchase:
		return s.executePurchaseOrPayment(transaction)
	case domaintransaction.TransactionTypeCreditPayment:
		return s.executeCreditPayment(transaction)
	case domaintransaction.TransactionTypeInstallmentPayment:
		return s.executeInstallmentPayment(transaction)
	case domaintransaction.TransactionTypeInstallmentRefund:
		return s.executeInstallmentRefund(transaction)
	default:
		// For other transaction types, no balance update is needed
		return nil
	}
}

// executeDeposit adds funds to the destination account
func (s *TransactionService) executeDeposit(transaction *domaintransaction.Transaction) error {
	_, err := s.accountService.AddFunds(
		stringValue(transaction.ToAccountID),
		transaction.Amount,
		fmt.Sprintf("Deposit - %s", transaction.Description),
		transaction.ID,
	)
	return err
}

// executeWithdrawal removes funds from the source account
func (s *TransactionService) executeWithdrawal(transaction *domaintransaction.Transaction) error {
	_, err := s.accountService.WithdrawFunds(
		stringValue(transaction.FromAccountID),
		transaction.Amount,
		fmt.Sprintf("Withdrawal - %s", transaction.Description),
		transaction.ID,
	)
	return err
}

// executeTransfer moves funds between accounts
func (s *TransactionService) executeTransfer(transaction *domaintransaction.Transaction) error {
	// First, withdraw from source account
	if err := s.executeWithdrawal(transaction); err != nil {
		return fmt.Errorf("failed to withdraw from source account: %w", err)
	}

	// Then, deposit to destination account
	// Create a temporary transaction for the deposit part
	depositTransaction := *transaction
	if err := s.executeDeposit(&depositTransaction); err != nil {
		// Rollback the withdrawal - add funds back to source account
		rollbackErr := s.rollbackWithdrawal(transaction)
		if rollbackErr != nil {
			return fmt.Errorf("transfer failed and rollback failed - original error: %w, rollback error: %v", err, rollbackErr)
		}
		return fmt.Errorf("failed to deposit to destination account: %w", err)
	}

	return nil
}

// executePurchaseOrPayment handles credit/debit transactions
func (s *TransactionService) executePurchaseOrPayment(transaction *domaintransaction.Transaction) error {
	// Get account info to determine if it's credit or debit
	accountInfo, err := s.accountService.GetAccountInfo(stringValue(transaction.FromAccountID))
	if err != nil {
		return fmt.Errorf("failed to get account info: %w", err)
	}

	if accountInfo.AccountType == "credit" {
		// For credit accounts, increase the used credit
		_, err := s.accountService.UpdateCreditUsage(
			stringValue(transaction.FromAccountID),
			transaction.Amount,
			fmt.Sprintf("%s - %s", transaction.Type, transaction.Description),
			transaction.ID,
		)
		return err
	} else {
		// For debit accounts, withdraw funds
		return s.executeWithdrawal(transaction)
	}
}

// executeCreditPayment handles credit card payment transactions
func (s *TransactionService) executeCreditPayment(transaction *domaintransaction.Transaction) error {
	// Credit payments reduce the used credit (negative usage)
	_, err := s.accountService.UpdateCreditUsage(
		stringValue(transaction.FromAccountID),
		-transaction.Amount, // Negative amount to reduce used credit
		fmt.Sprintf("Credit Payment - %s", transaction.Description),
		transaction.ID,
	)
	return err
}

// rollbackWithdrawal adds funds back to an account (used for transfer rollbacks)
func (s *TransactionService) rollbackWithdrawal(transaction *domaintransaction.Transaction) error {
	_, err := s.accountService.AddFunds(
		stringValue(transaction.FromAccountID),
		transaction.Amount,
		fmt.Sprintf("Rollback - %s", transaction.Description),
		fmt.Sprintf("rollback-%s", transaction.ID),
	)
	return err
}

// executeInstallmentPayment handles installment payment transactions
func (s *TransactionService) executeInstallmentPayment(transaction *domaintransaction.Transaction) error {
	// Installment payments are withdrawals from the paying account
	return s.executeWithdrawal(transaction)
}

// executeInstallmentRefund handles installment refund transactions
func (s *TransactionService) executeInstallmentRefund(transaction *domaintransaction.Transaction) error {
	// Installment refunds are deposits to the account
	return s.executeDeposit(transaction)
}

// GetTransactionByID retrieves a transaction by its ID with proper authorization
func (s *TransactionService) GetTransactionByID(id string, userID string) (*domaintransaction.Transaction, error) {
	if id == "" {
		return nil, errors.New("transaction ID is required")
	}

	transaction, err := s.transactionRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	// Check if user has permission to view this transaction
	if transaction.UserID != userID {
		return nil, errors.New("unauthorized: user does not have permission to view this transaction")
	}

	return transaction, nil
}

// GetTransactionsByUser retrieves transactions for a user with filtering
func (s *TransactionService) GetTransactionsByUser(userID string, filters TransactionFilters) ([]*domaintransaction.Transaction, int, error) {
	if userID == "" {
		return nil, 0, errors.New("user ID is required")
	}

	// Validate user exists
	if err := s.externalService.ValidateUser(userID); err != nil {
		return nil, 0, fmt.Errorf("user validation failed: %w", err)
	}

	transactions, total, err := s.transactionRepo.GetByUserID(userID, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get transactions: %w", err)
	}

	return transactions, total, nil
}

// UpdateTransactionStatus updates a transaction's status with proper validation
func (s *TransactionService) UpdateTransactionStatus(id string, status domaintransaction.TransactionStatus, reason string, updatedBy string) (*domaintransaction.Transaction, error) {
	transaction, err := s.transactionRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	oldStatus := transaction.Status

	// Validate status transition
	if !s.canTransitionToStatus(transaction.Status, status) {
		return nil, fmt.Errorf("invalid status transition from %s to %s", oldStatus, status)
	}

	// Update status
	transaction.Status = status

	// Save transaction
	updatedTransaction, err := s.transactionRepo.Update(transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to save transaction: %w", err)
	}

	return updatedTransaction, nil
}

// ProcessTransaction processes a pending transaction
func (s *TransactionService) ProcessTransaction(id string, processedBy string) error {
	transaction, err := s.transactionRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	if transaction.Status != domaintransaction.TransactionStatusPending {
		return fmt.Errorf("transaction is not in pending status, current status: %s", transaction.Status)
	}

	// Mark as completed - simplified for this version
	_, err = s.UpdateTransactionStatus(id, domaintransaction.TransactionStatusCompleted, "Transaction processed successfully", processedBy)
	return err
}

// CompleteTransaction marks a transaction as completed
func (s *TransactionService) CompleteTransaction(id string, completedBy string) error {
	_, err := s.UpdateTransactionStatus(id, domaintransaction.TransactionStatusCompleted, "Transaction completed successfully", completedBy)
	return err
}

// FailTransaction marks a transaction as failed
func (s *TransactionService) FailTransaction(id string, reason string, failedBy string) error {
	_, err := s.UpdateTransactionStatus(id, domaintransaction.TransactionStatusFailed, reason, failedBy)
	return err
}

// CancelTransaction marks a transaction as cancelled
func (s *TransactionService) CancelTransaction(id string, reason string, canceledBy string) error {
	transaction, err := s.transactionRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	if !transaction.CanBeCanceled() {
		return fmt.Errorf("transaction cannot be cancelled in its current state: %s", transaction.Status)
	}

	_, err = s.UpdateTransactionStatus(id, domaintransaction.TransactionStatusCanceled, reason, canceledBy)
	return err
}

// ReverseTransaction creates a reversal transaction
func (s *TransactionService) ReverseTransaction(id string, reason string, reversedBy string) (*domaintransaction.Transaction, error) {
	originalTransaction, err := s.transactionRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get original transaction: %w", err)
	}

	if !originalTransaction.CanBeReversed() {
		return nil, fmt.Errorf("transaction cannot be reversed in its current state: %s", originalTransaction.Status)
	}

	// Create reversal transaction using the entity method
	reversalTransaction, err := originalTransaction.Reverse(reversedBy)
	if err != nil {
		return nil, fmt.Errorf("failed to create reversal: %w", err)
	}

	reversalTransaction.Description = fmt.Sprintf("Reversal of transaction %s - %s", id, reason)

	// Save reversal transaction
	savedReversal, err := s.transactionRepo.Create(reversalTransaction)
	if err != nil {
		return nil, fmt.Errorf("failed to create reversal transaction: %w", err)
	}

	// Update original transaction status
	_, err = s.UpdateTransactionStatus(id, domaintransaction.TransactionStatusReversed, reason, reversedBy)
	if err != nil {
		return nil, fmt.Errorf("failed to update original transaction status: %w", err)
	}

	return savedReversal, nil
}

// canTransitionToStatus checks if a status transition is valid
func (s *TransactionService) canTransitionToStatus(currentStatus, newStatus domaintransaction.TransactionStatus) bool {
	// Define valid transitions
	validTransitions := map[domaintransaction.TransactionStatus][]domaintransaction.TransactionStatus{
		domaintransaction.TransactionStatusPending: {
			domaintransaction.TransactionStatusCompleted,
			domaintransaction.TransactionStatusFailed,
			domaintransaction.TransactionStatusCanceled,
		},
		domaintransaction.TransactionStatusCompleted: {
			domaintransaction.TransactionStatusReversed,
		},
		domaintransaction.TransactionStatusFailed:   {},
		domaintransaction.TransactionStatusCanceled: {},
		domaintransaction.TransactionStatusReversed: {},
	}

	allowedTransitions, exists := validTransitions[currentStatus]
	if !exists {
		return false
	}

	for _, allowed := range allowedTransitions {
		if allowed == newStatus {
			return true
		}
	}

	return false
}

// Convenience methods for specific transaction types

func (s *TransactionService) ProcessWalletDeposit(userID string, accountID string, amount float64, description string, initiatedBy string) (*domaintransaction.Transaction, error) {
	request := CreateTransactionRequest{
		UserID:      userID,
		Type:        domaintransaction.TransactionTypeWalletDeposit,
		Amount:      amount,
		Currency:    "USD", // Default currency
		ToAccountID: &accountID,
		Description: description,
	}

	transaction, err := s.CreateTransaction(request, initiatedBy)
	if err != nil {
		return nil, err
	}

	// Process immediately
	if err := s.ProcessTransaction(transaction.ID, initiatedBy); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *TransactionService) ProcessWalletWithdrawal(userID string, accountID string, amount float64, description string, initiatedBy string) (*domaintransaction.Transaction, error) {
	request := CreateTransactionRequest{
		UserID:        userID,
		Type:          domaintransaction.TransactionTypeWalletWithdrawal,
		Amount:        amount,
		Currency:      "USD",
		FromAccountID: &accountID,
		Description:   description,
	}

	transaction, err := s.CreateTransaction(request, initiatedBy)
	if err != nil {
		return nil, err
	}

	if err := s.ProcessTransaction(transaction.ID, initiatedBy); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *TransactionService) ProcessWalletTransfer(userID string, fromAccountID string, toAccountID string, amount float64, description string, initiatedBy string) (*domaintransaction.Transaction, error) {
	request := CreateTransactionRequest{
		UserID:        userID,
		Type:          domaintransaction.TransactionTypeWalletTransfer,
		Amount:        amount,
		Currency:      "USD",
		FromAccountID: &fromAccountID,
		ToAccountID:   &toAccountID,
		Description:   description,
	}

	transaction, err := s.CreateTransaction(request, initiatedBy)
	if err != nil {
		return nil, err
	}

	if err := s.ProcessTransaction(transaction.ID, initiatedBy); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *TransactionService) ProcessCreditCardCharge(userID string, cardID string, amount float64, description string, merchantName string, initiatedBy string) (*domaintransaction.Transaction, error) {
	request := CreateTransactionRequest{
		UserID:       userID,
		Type:         domaintransaction.TransactionTypeCreditCharge,
		Amount:       amount,
		Currency:     "USD",
		FromCardID:   &cardID,
		Description:  description,
		MerchantName: merchantName,
	}

	transaction, err := s.CreateTransaction(request, initiatedBy)
	if err != nil {
		return nil, err
	}

	if err := s.ProcessTransaction(transaction.ID, initiatedBy); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *TransactionService) ProcessCreditCardPayment(userID string, cardID string, amount float64, paymentMethod domaintransaction.PaymentMethod, initiatedBy string) (*domaintransaction.Transaction, error) {
	request := CreateTransactionRequest{
		UserID:        userID,
		Type:          domaintransaction.TransactionTypeCreditPayment,
		Amount:        amount,
		Currency:      "USD",
		ToCardID:      &cardID,
		PaymentMethod: paymentMethod,
		Description:   "Credit card payment",
	}

	transaction, err := s.CreateTransaction(request, initiatedBy)
	if err != nil {
		return nil, err
	}

	if err := s.ProcessTransaction(transaction.ID, initiatedBy); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *TransactionService) ProcessDebitCardPurchase(userID string, cardID string, amount float64, description string, merchantName string, initiatedBy string) (*domaintransaction.Transaction, error) {
	request := CreateTransactionRequest{
		UserID:       userID,
		Type:         domaintransaction.TransactionTypeDebitPurchase,
		Amount:       amount,
		Currency:     "USD",
		FromCardID:   &cardID,
		Description:  description,
		MerchantName: merchantName,
	}

	transaction, err := s.CreateTransaction(request, initiatedBy)
	if err != nil {
		return nil, err
	}

	if err := s.ProcessTransaction(transaction.ID, initiatedBy); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *TransactionService) ProcessAccountTransfer(userID string, fromAccountID string, toAccountID string, amount float64, description string, initiatedBy string) (*domaintransaction.Transaction, error) {
	request := CreateTransactionRequest{
		UserID:        userID,
		Type:          domaintransaction.TransactionTypeAccountTransfer,
		Amount:        amount,
		Currency:      "USD",
		FromAccountID: &fromAccountID,
		ToAccountID:   &toAccountID,
		Description:   description,
	}

	transaction, err := s.CreateTransaction(request, initiatedBy)
	if err != nil {
		return nil, err
	}

	if err := s.ProcessTransaction(transaction.ID, initiatedBy); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *TransactionService) ProcessAccountDeposit(userID string, accountID string, amount float64, description string, initiatedBy string) (*domaintransaction.Transaction, error) {
	request := CreateTransactionRequest{
		UserID:      userID,
		Type:        domaintransaction.TransactionTypeAccountDeposit,
		Amount:      amount,
		Currency:    "USD",
		ToAccountID: &accountID,
		Description: description,
	}

	transaction, err := s.CreateTransaction(request, initiatedBy)
	if err != nil {
		return nil, err
	}

	if err := s.ProcessTransaction(transaction.ID, initiatedBy); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *TransactionService) ProcessAccountWithdraw(userID string, accountID string, amount float64, description string, initiatedBy string) (*domaintransaction.Transaction, error) {
	request := CreateTransactionRequest{
		UserID:        userID,
		Type:          domaintransaction.TransactionTypeAccountWithdraw,
		Amount:        amount,
		Currency:      "USD",
		FromAccountID: &accountID,
		Description:   description,
	}

	transaction, err := s.CreateTransaction(request, initiatedBy)
	if err != nil {
		return nil, err
	}

	if err := s.ProcessTransaction(transaction.ID, initiatedBy); err != nil {
		return nil, err
	}

	return transaction, nil
}
