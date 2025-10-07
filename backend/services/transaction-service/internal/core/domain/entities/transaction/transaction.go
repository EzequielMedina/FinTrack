package transaction

import (
	"errors"
	"fmt"
	"time"
)

// TransactionType represents the type of transaction
type TransactionType string

const (
	// Wallet transactions
	TransactionTypeWalletDeposit    TransactionType = "wallet_deposit"
	TransactionTypeWalletWithdrawal TransactionType = "wallet_withdrawal"
	TransactionTypeWalletTransfer   TransactionType = "wallet_transfer"

	// Credit card transactions
	TransactionTypeCreditCharge  TransactionType = "credit_charge"
	TransactionTypeCreditPayment TransactionType = "credit_payment"
	TransactionTypeCreditRefund  TransactionType = "credit_refund"

	// Debit card transactions
	TransactionTypeDebitPurchase   TransactionType = "debit_purchase"
	TransactionTypeDebitWithdrawal TransactionType = "debit_withdrawal"
	TransactionTypeDebitRefund     TransactionType = "debit_refund"

	// Account transfers
	TransactionTypeAccountTransfer TransactionType = "account_transfer"
	TransactionTypeAccountDeposit  TransactionType = "account_deposit"
	TransactionTypeAccountWithdraw TransactionType = "account_withdraw"

	// Installment transactions
	TransactionTypeInstallmentPayment    TransactionType = "installment_payment"
	TransactionTypeInstallmentRefund     TransactionType = "installment_refund"
	TransactionTypeInstallmentCompletion TransactionType = "installment_plan_completion"
)

// TransactionStatus represents the current status of a transaction
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
	TransactionStatusCanceled  TransactionStatus = "canceled"
	TransactionStatusReversed  TransactionStatus = "reversed"
)

// PaymentMethod represents the method used for payment
type PaymentMethod string

const (
	PaymentMethodCash                PaymentMethod = "cash"
	PaymentMethodBankTransfer        PaymentMethod = "bank_transfer"
	PaymentMethodCreditCard          PaymentMethod = "credit_card"
	PaymentMethodDebitCard           PaymentMethod = "debit_card"
	PaymentMethodWallet              PaymentMethod = "wallet"
	PaymentMethodInstallmentCompletion PaymentMethod = "installment_completion"
)

// Transaction represents a financial transaction in the system
type Transaction struct {
	// Core identity
	ID          string `json:"id" gorm:"primaryKey;type:varchar(36)"`
	ReferenceID string `json:"referenceId" gorm:"type:varchar(100);index"`
	ExternalID  string `json:"externalId" gorm:"type:varchar(100);index"`

	// Transaction details
	Type     TransactionType   `json:"type" gorm:"type:varchar(50);not null;index"`
	Status   TransactionStatus `json:"status" gorm:"type:varchar(20);not null;default:'pending';index"`
	Amount   float64           `json:"amount" gorm:"type:decimal(15,2);not null"`
	Currency string            `json:"currency" gorm:"type:varchar(3);not null;default:'ARS'"`

	// Source and destination
	FromAccountID *string `json:"fromAccountId" gorm:"type:varchar(36);index"`
	ToAccountID   *string `json:"toAccountId" gorm:"type:varchar(36);index"`
	FromCardID    *string `json:"fromCardId" gorm:"type:varchar(36);index"`
	ToCardID      *string `json:"toCardId" gorm:"type:varchar(36);index"`

	// User information
	UserID      string `json:"userId" gorm:"type:varchar(36);not null;index"`
	InitiatedBy string `json:"initiatedBy" gorm:"type:varchar(36);not null"`

	// Transaction metadata
	Description   string        `json:"description" gorm:"type:text"`
	PaymentMethod PaymentMethod `json:"paymentMethod" gorm:"type:varchar(30)"`
	MerchantName  string        `json:"merchantName" gorm:"type:varchar(255)"`
	MerchantID    string        `json:"merchantId" gorm:"type:varchar(100)"`

	// Balance tracking
	PreviousBalance float64 `json:"previousBalance" gorm:"type:decimal(15,2)"`
	NewBalance      float64 `json:"newBalance" gorm:"type:decimal(15,2)"`

	// Processing details
	ProcessedAt   *time.Time `json:"processedAt"`
	FailedAt      *time.Time `json:"failedAt"`
	FailureReason string     `json:"failureReason" gorm:"type:text"`

	// Additional metadata
	Metadata map[string]interface{} `json:"metadata" gorm:"type:json"`
	Tags     []string               `json:"tags" gorm:"type:json"`

	// Audit fields
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

// TransactionRule represents business rules for transaction validation
type TransactionRule struct {
	MaxDailyAmount   float64 `json:"maxDailyAmount"`
	MaxSingleAmount  float64 `json:"maxSingleAmount"`
	RequiresApproval bool    `json:"requiresApproval"`
	AllowedHours     string  `json:"allowedHours"`
}

// Validation methods

// IsValidType checks if the transaction type is valid
func IsValidTransactionType(transType TransactionType) bool {
	switch transType {
	case TransactionTypeWalletDeposit, TransactionTypeWalletWithdrawal, TransactionTypeWalletTransfer,
		TransactionTypeCreditCharge, TransactionTypeCreditPayment, TransactionTypeCreditRefund,
		TransactionTypeDebitPurchase, TransactionTypeDebitWithdrawal, TransactionTypeDebitRefund,
		TransactionTypeAccountTransfer, TransactionTypeAccountDeposit, TransactionTypeAccountWithdraw,
		TransactionTypeInstallmentPayment, TransactionTypeInstallmentRefund, TransactionTypeInstallmentCompletion:
		return true
	default:
		return false
	}
}

// IsValidStatus checks if the transaction status is valid
func IsValidTransactionStatus(status TransactionStatus) bool {
	switch status {
	case TransactionStatusPending, TransactionStatusCompleted, TransactionStatusFailed,
		TransactionStatusCanceled, TransactionStatusReversed:
		return true
	default:
		return false
	}
}

// IsValidPaymentMethod checks if the payment method is valid
func IsValidPaymentMethod(method PaymentMethod) bool {
	switch method {
	case PaymentMethodCash, PaymentMethodBankTransfer, PaymentMethodCreditCard,
		PaymentMethodDebitCard, PaymentMethodWallet, PaymentMethodInstallmentCompletion:
		return true
	default:
		return false
	}
}

// Business logic methods

// CanBeCompleted checks if transaction can be marked as completed
func (t *Transaction) CanBeCompleted() bool {
	return t.Status == TransactionStatusPending
}

// CanBeCanceled checks if transaction can be canceled
func (t *Transaction) CanBeCanceled() bool {
	return t.Status == TransactionStatusPending
}

// CanBeReversed checks if transaction can be reversed
func (t *Transaction) CanBeReversed() bool {
	return t.Status == TransactionStatusCompleted &&
		time.Since(t.CreatedAt) <= 24*time.Hour // 24 hour reversal window
}

// Complete marks the transaction as completed
func (t *Transaction) Complete() error {
	if !t.CanBeCompleted() {
		return errors.New("transaction cannot be completed in current state")
	}

	now := time.Now()
	t.Status = TransactionStatusCompleted
	t.ProcessedAt = &now
	return nil
}

// Cancel marks the transaction as canceled
func (t *Transaction) Cancel(reason string) error {
	if !t.CanBeCanceled() {
		return errors.New("transaction cannot be canceled in current state")
	}

	t.Status = TransactionStatusCanceled
	t.FailureReason = reason
	return nil
}

// Fail marks the transaction as failed
func (t *Transaction) Fail(reason string) error {
	now := time.Now()
	t.Status = TransactionStatusFailed
	t.FailedAt = &now
	t.FailureReason = reason
	return nil
}

// Reverse creates a reversal transaction
func (t *Transaction) Reverse(initiatedBy string) (*Transaction, error) {
	if !t.CanBeReversed() {
		return nil, errors.New("transaction cannot be reversed")
	}

	reversalType := t.getReversalType()
	if reversalType == "" {
		return nil, errors.New("no reversal type defined for this transaction type")
	}

	reversal := &Transaction{
		Type:          reversalType,
		Status:        TransactionStatusPending,
		Amount:        t.Amount,
		Currency:      t.Currency,
		FromAccountID: t.ToAccountID, // Swap source and destination
		ToAccountID:   t.FromAccountID,
		FromCardID:    t.ToCardID,
		ToCardID:      t.FromCardID,
		UserID:        t.UserID,
		InitiatedBy:   initiatedBy,
		Description:   fmt.Sprintf("Reversal of transaction %s", t.ID),
		PaymentMethod: t.PaymentMethod,
		ReferenceID:   t.ID, // Reference to original transaction
	}

	return reversal, nil
}

// getReversalType returns the appropriate reversal transaction type
func (t *Transaction) getReversalType() TransactionType {
	switch t.Type {
	case TransactionTypeCreditCharge:
		return TransactionTypeCreditRefund
	case TransactionTypeDebitPurchase:
		return TransactionTypeDebitRefund
	case TransactionTypeWalletWithdrawal:
		return TransactionTypeWalletDeposit
	case TransactionTypeWalletDeposit:
		return TransactionTypeWalletWithdrawal
	default:
		return "" // No direct reversal type
	}
}

// IsDebitTransaction checks if this is a debit (outgoing) transaction
func (t *Transaction) IsDebitTransaction() bool {
	switch t.Type {
	case TransactionTypeWalletWithdrawal, TransactionTypeCreditCharge,
		TransactionTypeDebitPurchase, TransactionTypeAccountWithdraw:
		return true
	default:
		return false
	}
}

// IsCreditTransaction checks if this is a credit (incoming) transaction
func (t *Transaction) IsCreditTransaction() bool {
	switch t.Type {
	case TransactionTypeWalletDeposit, TransactionTypeCreditPayment,
		TransactionTypeDebitRefund, TransactionTypeAccountDeposit:
		return true
	default:
		return false
	}
}

// RequiresSourceAccount checks if transaction requires a source account/card
func (t *Transaction) RequiresSourceAccount() bool {
	switch t.Type {
	case TransactionTypeWalletWithdrawal, TransactionTypeWalletTransfer,
		TransactionTypeCreditCharge, TransactionTypeDebitPurchase,
		TransactionTypeAccountTransfer, TransactionTypeAccountWithdraw:
		return true
	default:
		return false
	}
}

// RequiresDestinationAccount checks if transaction requires a destination account/card
func (t *Transaction) RequiresDestinationAccount() bool {
	switch t.Type {
	case TransactionTypeWalletDeposit, TransactionTypeWalletTransfer,
		TransactionTypeCreditPayment, TransactionTypeAccountTransfer,
		TransactionTypeAccountDeposit:
		return true
	default:
		return false
	}
}

// GetDisplayName returns a human-readable name for the transaction type
func (t *Transaction) GetDisplayName() string {
	displayNames := map[TransactionType]string{
		TransactionTypeWalletDeposit:    "Depósito en Billetera",
		TransactionTypeWalletWithdrawal: "Retiro de Billetera",
		TransactionTypeWalletTransfer:   "Transferencia de Billetera",
		TransactionTypeCreditCharge:     "Cargo Tarjeta de Crédito",
		TransactionTypeCreditPayment:    "Pago Tarjeta de Crédito",
		TransactionTypeCreditRefund:     "Reembolso Tarjeta de Crédito",
		TransactionTypeDebitPurchase:    "Compra con Débito",
		TransactionTypeDebitWithdrawal:  "Retiro con Débito",
		TransactionTypeDebitRefund:      "Reembolso Débito",
		TransactionTypeAccountTransfer:  "Transferencia entre Cuentas",
		TransactionTypeAccountDeposit:   "Depósito en Cuenta",
		TransactionTypeAccountWithdraw:  "Retiro de Cuenta",
	}

	if name, exists := displayNames[t.Type]; exists {
		return name
	}
	return string(t.Type)
}

// Validate performs comprehensive validation of the transaction
func (t *Transaction) Validate() error {
	// Type validation
	if !IsValidTransactionType(t.Type) {
		return errors.New("invalid transaction type")
	}

	// Status validation
	if !IsValidTransactionStatus(t.Status) {
		return errors.New("invalid transaction status")
	}

	// Amount validation
	if t.Amount <= 0 {
		return errors.New("transaction amount must be positive")
	}

	// Currency validation
	if t.Currency == "" {
		return errors.New("currency is required")
	}

	// User validation
	if t.UserID == "" {
		return errors.New("user ID is required")
	}

	if t.InitiatedBy == "" {
		return errors.New("initiated by is required")
	}

	// Source/destination validation
	if t.RequiresSourceAccount() {
		if t.FromAccountID == nil && t.FromCardID == nil {
			return errors.New("source account or card is required for this transaction type")
		}
	}

	if t.RequiresDestinationAccount() {
		if t.ToAccountID == nil && t.ToCardID == nil {
			return errors.New("destination account or card is required for this transaction type")
		}
	}

	// Payment method validation
	if t.PaymentMethod != "" && !IsValidPaymentMethod(t.PaymentMethod) {
		return errors.New("invalid payment method")
	}

	return nil
}
