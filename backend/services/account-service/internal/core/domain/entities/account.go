package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AccountType represents the type of account
type AccountType string

const (
	AccountTypeChecking AccountType = "checking"
	AccountTypeSavings  AccountType = "savings"
	AccountTypeCredit   AccountType = "credit"
	AccountTypeDebit    AccountType = "debit"
	AccountTypeWallet   AccountType = "wallet"
	// New integrated types
	AccountTypeBankAccount AccountType = "bank_account" // Can have multiple cards
)

// Currency represents the currency type
type Currency string

const (
	CurrencyUSD Currency = "USD"
	CurrencyARS Currency = "ARS"
	CurrencyEUR Currency = "EUR"
)

// CardType represents the type of card
type CardType string

const (
	CardTypeCredit CardType = "credit"
	CardTypeDebit  CardType = "debit"
)

// CardBrand represents the card brand
type CardBrand string

const (
	CardBrandVisa       CardBrand = "visa"
	CardBrandMastercard CardBrand = "mastercard"
	CardBrandAmex       CardBrand = "amex"
	CardBrandDiscover   CardBrand = "discover"
	CardBrandDiners     CardBrand = "diners"
	CardBrandOther      CardBrand = "other"
)

// CardStatus represents the status of a card
type CardStatus string

const (
	CardStatusActive   CardStatus = "active"
	CardStatusInactive CardStatus = "inactive"
	CardStatusBlocked  CardStatus = "blocked"
	CardStatusExpired  CardStatus = "expired"
)

// InstallmentPlanStatus represents the status of an installment plan
type InstallmentPlanStatus string

const (
	InstallmentPlanStatusActive    InstallmentPlanStatus = "active"
	InstallmentPlanStatusCompleted InstallmentPlanStatus = "completed"
	InstallmentPlanStatusCancelled InstallmentPlanStatus = "cancelled"
	InstallmentPlanStatusSuspended InstallmentPlanStatus = "suspended"
)

// InstallmentStatus represents the status of an individual installment
type InstallmentStatus string

const (
	InstallmentStatusPending   InstallmentStatus = "pending"
	InstallmentStatusPaid      InstallmentStatus = "paid"
	InstallmentStatusOverdue   InstallmentStatus = "overdue"
	InstallmentStatusCancelled InstallmentStatus = "cancelled"
	InstallmentStatusPartial   InstallmentStatus = "partial"
)

// Account represents a financial account in the system
type Account struct {
	ID          string      `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID      string      `gorm:"type:varchar(36);not null;index" json:"user_id"`
	AccountType AccountType `gorm:"type:varchar(20);not null;index" json:"account_type"`
	Name        string      `gorm:"type:varchar(100);not null" json:"name"`
	Description string      `gorm:"type:text" json:"description"`
	Currency    Currency    `gorm:"type:varchar(3);not null;index" json:"currency"`
	Balance     float64     `gorm:"type:decimal(15,2);not null;default:0" json:"balance"`

	// Cards relationship (optional - only for bank_account type)
	Cards []Card `gorm:"foreignKey:AccountID;constraint:OnDelete:CASCADE" json:"cards,omitempty"`

	// Credit card specific fields (legacy - for backward compatibility)
	CreditLimit *float64   `gorm:"type:decimal(15,2);null" json:"credit_limit,omitempty"`
	ClosingDate *time.Time `gorm:"type:date;null" json:"closing_date,omitempty"`
	DueDate     *time.Time `gorm:"type:date;null" json:"due_date,omitempty"`

	// Personal identification (for virtual wallets)
	DNI *string `gorm:"type:varchar(20);null" json:"dni,omitempty"`

	IsActive  bool           `gorm:"not null;default:true;index" json:"is_active"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// Card represents a payment card associated with an account
type Card struct {
	ID              string     `gorm:"type:varchar(36);primaryKey" json:"id"`
	AccountID       string     `gorm:"type:varchar(36);not null;index" json:"account_id"`
	CardType        CardType   `gorm:"type:varchar(10);not null" json:"card_type"`
	CardBrand       CardBrand  `gorm:"type:varchar(20);not null" json:"card_brand"`
	LastFourDigits  string     `gorm:"type:varchar(4);not null" json:"last_four_digits"`
	MaskedNumber    string     `gorm:"type:varchar(19);not null" json:"masked_number"` // **** **** **** 1234
	HolderName      string     `gorm:"type:varchar(100);not null" json:"holder_name"`
	ExpirationMonth int        `gorm:"type:int;not null" json:"expiration_month"`
	ExpirationYear  int        `gorm:"type:int;not null" json:"expiration_year"`
	Status          CardStatus `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	IsDefault       bool       `gorm:"not null;default:false" json:"is_default"`
	Nickname        string     `gorm:"type:varchar(50)" json:"nickname"`

	// Balance field - usage depends on card type:
	// - Credit cards: debt amount (positive = owed to bank)
	// - Debit cards: should always be 0 (uses account balance)
	Balance float64 `gorm:"type:decimal(15,2);not null;default:0" json:"balance"`

	// Credit card specific fields
	CreditLimit *float64   `gorm:"type:decimal(15,2);null" json:"credit_limit,omitempty"`
	ClosingDate *time.Time `gorm:"type:date;null" json:"closing_date,omitempty"`
	DueDate     *time.Time `gorm:"type:date;null" json:"due_date,omitempty"`

	// Security - encrypted fields (stored separately for security)
	EncryptedNumber string `gorm:"type:text;not null" json:"-"` // Never expose in JSON
	KeyFingerprint  string `gorm:"type:varchar(64);not null" json:"-"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relationship back to account
	Account Account `gorm:"foreignKey:AccountID" json:"-"`
}

// TableName returns the table name for the Account model
func (Account) TableName() string {
	return "accounts"
}

// TableName returns the table name for the Card model
func (Card) TableName() string {
	return "cards"
}

// BeforeCreate is called before creating a new account
func (a *Account) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}

// IsValidAccountType checks if the account type is valid
func IsValidAccountType(accountType AccountType) bool {
	switch accountType {
	case AccountTypeChecking, AccountTypeSavings, AccountTypeCredit, AccountTypeDebit,
		AccountTypeWallet, AccountTypeBankAccount:
		return true
	default:
		return false
	}
}

// IsValidCurrency checks if the currency is valid
func IsValidCurrency(currency Currency) bool {
	switch currency {
	case CurrencyUSD, CurrencyARS, CurrencyEUR:
		return true
	default:
		return false
	}
}

// Validate validates the account data
func (a *Account) Validate() error {
	if a.UserID == "" {
		return &ValidationError{Field: "user_id", Message: "user ID is required"}
	}
	if a.Name == "" {
		return &ValidationError{Field: "name", Message: "account name is required"}
	}
	if !IsValidAccountType(a.AccountType) {
		return &ValidationError{Field: "account_type", Message: "invalid account type"}
	}
	if !IsValidCurrency(a.Currency) {
		return &ValidationError{Field: "currency", Message: "invalid currency"}
	}

	// Validate based on account type
	if a.AccountType == AccountTypeWallet && a.DNI == nil {
		return &ValidationError{Field: "dni", Message: "DNI is required for virtual wallets"}
	}

	return nil
}

// IsValidCardType checks if the card type is valid
func IsValidCardType(cardType CardType) bool {
	switch cardType {
	case CardTypeCredit, CardTypeDebit:
		return true
	default:
		return false
	}
}

// IsValidCardBrand checks if the card brand is valid
func IsValidCardBrand(brand CardBrand) bool {
	switch brand {
	case CardBrandVisa, CardBrandMastercard, CardBrandAmex, CardBrandDiscover, CardBrandDiners, CardBrandOther:
		return true
	default:
		return false
	}
}

// IsValidCardStatus checks if the card status is valid
func IsValidCardStatus(status CardStatus) bool {
	switch status {
	case CardStatusActive, CardStatusInactive, CardStatusBlocked, CardStatusExpired:
		return true
	default:
		return false
	}
}

// Validate validates the card data
func (c *Card) Validate() error {
	if c.AccountID == "" {
		return &ValidationError{Field: "account_id", Message: "account ID is required"}
	}
	if c.HolderName == "" {
		return &ValidationError{Field: "holder_name", Message: "holder name is required"}
	}
	if !IsValidCardType(c.CardType) {
		return &ValidationError{Field: "card_type", Message: "invalid card type"}
	}
	if !IsValidCardBrand(c.CardBrand) {
		return &ValidationError{Field: "card_brand", Message: "invalid card brand"}
	}
	if !IsValidCardStatus(c.Status) {
		return &ValidationError{Field: "status", Message: "invalid card status"}
	}
	if c.ExpirationMonth < 1 || c.ExpirationMonth > 12 {
		return &ValidationError{Field: "expiration_month", Message: "expiration month must be between 1 and 12"}
	}
	if c.ExpirationYear < time.Now().Year() {
		return &ValidationError{Field: "expiration_year", Message: "expiration year cannot be in the past"}
	}
	if len(c.LastFourDigits) != 4 {
		return &ValidationError{Field: "last_four_digits", Message: "last four digits must be exactly 4 characters"}
	}

	// Validate credit card specific fields
	if c.CardType == CardTypeCredit && c.CreditLimit == nil {
		return &ValidationError{Field: "credit_limit", Message: "credit limit is required for credit cards"}
	}

	return nil
}

// IsWallet checks if the account is a virtual wallet
func (a *Account) IsWallet() bool {
	return a.AccountType == AccountTypeWallet
}

// IsBankAccount checks if the account is a bank account (can have cards)
func (a *Account) IsBankAccount() bool {
	return a.AccountType == AccountTypeBankAccount
}

// CanHaveCards checks if the account type supports cards
func (a *Account) CanHaveCards() bool {
	return a.IsBankAccount() || a.AccountType == AccountTypeChecking ||
		a.AccountType == AccountTypeCredit || a.AccountType == AccountTypeDebit
}

// GetDefaultCard returns the default card for the account
func (a *Account) GetDefaultCard() *Card {
	for _, card := range a.Cards {
		if card.IsDefault && card.Status == CardStatusActive {
			return &card
		}
	}
	return nil
}

// HasActiveCards checks if the account has any active cards
func (a *Account) HasActiveCards() bool {
	for _, card := range a.Cards {
		if card.Status == CardStatusActive {
			return true
		}
	}
	return false
}

// IsExpired checks if the card is expired
func (c *Card) IsExpired() bool {
	now := time.Now()
	expiry := time.Date(c.ExpirationYear, time.Month(c.ExpirationMonth), 1, 0, 0, 0, 0, time.UTC)
	return now.After(expiry)
}

// IsActive checks if the card is active and not expired
func (c *Card) IsActive() bool {
	return c.Status == CardStatusActive && !c.IsExpired()
}

// GetAvailableBalance returns the available balance for the card
func (c *Card) GetAvailableBalance() float64 {
	if c.CardType == CardTypeDebit {
		// For debit cards, available balance is the account balance
		return c.Account.Balance
	} else if c.CardType == CardTypeCredit && c.CreditLimit != nil {
		// For credit cards, available balance is credit limit minus debt
		return *c.CreditLimit - c.Balance
	}
	return 0
}

// GetDebt returns the debt amount for credit cards
func (c *Card) GetDebt() float64 {
	if c.CardType == CardTypeCredit {
		return c.Balance // Positive balance = debt
	}
	return 0 // Debit cards don't have debt
}

// CanCharge checks if the card can be charged with the specified amount
func (c *Card) CanCharge(amount float64) bool {
	if !c.IsActive() {
		return false
	}

	if c.CardType == CardTypeDebit {
		// For debit cards, check account balance
		return c.Account.Balance >= amount
	} else if c.CardType == CardTypeCredit {
		// For credit cards, check available credit
		return c.GetAvailableBalance() >= amount
	}

	return false
}

// Charge processes a charge to the card
func (c *Card) Charge(amount float64) error {
	if !c.CanCharge(amount) {
		return &ValidationError{Field: "amount", Message: "insufficient funds or credit"}
	}

	if c.CardType == CardTypeDebit {
		// For debit cards, deduct from account balance
		c.Account.Balance -= amount
	} else if c.CardType == CardTypeCredit {
		// For credit cards, add to debt
		c.Balance += amount
	}

	return nil
}

// Payment processes a payment to a credit card
func (c *Card) Payment(amount float64) error {
	if c.CardType != CardTypeCredit {
		return &ValidationError{Field: "card_type", Message: "payments only allowed for credit cards"}
	}

	if amount <= 0 {
		return &ValidationError{Field: "amount", Message: "payment amount must be positive"}
	}

	// Reduce debt (balance can go negative if overpaid)
	c.Balance -= amount
	return nil
}

// GetMinimumPayment calculates minimum payment for credit cards
func (c *Card) GetMinimumPayment() float64 {
	if c.CardType != CardTypeCredit || c.Balance <= 0 {
		return 0
	}

	// Example: 5% of debt or $500, whichever is greater
	minPercentage := c.Balance * 0.05
	minAmount := 500.0

	if minPercentage > minAmount {
		return minPercentage
	}
	return minAmount
}

// IsOverdue checks if the credit card payment is overdue
func (c *Card) IsOverdue() bool {
	if c.CardType != CardTypeCredit || c.DueDate == nil || c.Balance <= 0 {
		return false
	}

	return time.Now().After(*c.DueDate)
}

// GetNextClosingDate calculates the next closing date
func (c *Card) GetNextClosingDate() *time.Time {
	if c.CardType != CardTypeCredit || c.ClosingDate == nil {
		return nil
	}

	now := time.Now()
	closingDay := c.ClosingDate.Day()

	// Find next closing date
	nextClosing := time.Date(now.Year(), now.Month(), closingDay, 0, 0, 0, 0, time.Local)
	if nextClosing.Before(now) {
		nextClosing = nextClosing.AddDate(0, 1, 0)
	}

	return &nextClosing
}

// GetNextDueDate calculates the next due date
func (c *Card) GetNextDueDate() *time.Time {
	if c.CardType != CardTypeCredit || c.DueDate == nil {
		return nil
	}

	now := time.Now()
	dueDay := c.DueDate.Day()

	// Find next due date
	nextDue := time.Date(now.Year(), now.Month(), dueDay, 0, 0, 0, 0, time.Local)
	if nextDue.Before(now) {
		nextDue = nextDue.AddDate(0, 1, 0)
	}

	return &nextDue
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return e.Message
}

// InstallmentPlan represents a credit card installment payment plan
type InstallmentPlan struct {
	ID            string `gorm:"type:varchar(36);primaryKey" json:"id"`
	TransactionID string `gorm:"type:varchar(36);not null;index" json:"transaction_id"`
	CardID        string `gorm:"type:varchar(36);not null;index" json:"card_id"`
	UserID        string `gorm:"type:varchar(36);not null;index" json:"user_id"`

	// Plan details
	TotalAmount       float64   `gorm:"type:decimal(15,2);not null" json:"total_amount"`
	InstallmentsCount int       `gorm:"type:int;not null" json:"installments_count"`
	InstallmentAmount float64   `gorm:"type:decimal(15,2);not null" json:"installment_amount"`
	StartDate         time.Time `gorm:"type:date;not null" json:"start_date"`

	// Merchant information
	MerchantName string `gorm:"type:varchar(255)" json:"merchant_name,omitempty"`
	MerchantID   string `gorm:"type:varchar(100)" json:"merchant_id,omitempty"`
	Description  string `gorm:"type:text" json:"description,omitempty"`

	// Status tracking
	Status           InstallmentPlanStatus `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	PaidInstallments int                   `gorm:"type:int;not null;default:0" json:"paid_installments"`
	RemainingAmount  float64               `gorm:"type:decimal(15,2);not null" json:"remaining_amount"`

	// Interest and fees (for future enhancement)
	InterestRate  float64 `gorm:"type:decimal(5,2);default:0.00" json:"interest_rate"`
	TotalInterest float64 `gorm:"type:decimal(15,2);default:0.00" json:"total_interest"`
	AdminFee      float64 `gorm:"type:decimal(15,2);default:0.00" json:"admin_fee"`

	// Audit fields
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	CompletedAt *time.Time `gorm:"type:timestamp;null" json:"completed_at,omitempty"`
	CancelledAt *time.Time `gorm:"type:timestamp;null" json:"cancelled_at,omitempty"`

	// Relationships
	Card         Card          `gorm:"foreignKey:CardID" json:"card,omitempty"`
	Installments []Installment `gorm:"foreignKey:PlanID;constraint:OnDelete:CASCADE" json:"installments,omitempty"`
}

// Installment represents an individual installment payment
type Installment struct {
	ID                string `gorm:"type:varchar(36);primaryKey" json:"id"`
	PlanID            string `gorm:"type:varchar(36);not null;index" json:"plan_id"`
	InstallmentNumber int    `gorm:"type:int;not null" json:"installment_number"`

	// Payment details
	Amount   float64           `gorm:"type:decimal(15,2);not null" json:"amount"`
	DueDate  time.Time         `gorm:"type:date;not null" json:"due_date"`
	PaidDate *time.Time        `gorm:"type:timestamp;null" json:"paid_date,omitempty"`
	Status   InstallmentStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`

	// Payment information
	PaidAmount       float64 `gorm:"type:decimal(15,2);default:0.00" json:"paid_amount"`
	RemainingAmount  float64 `gorm:"type:decimal(15,2);not null" json:"remaining_amount"`
	PaymentMethod    string  `gorm:"type:varchar(30)" json:"payment_method,omitempty"`
	PaymentReference string  `gorm:"type:varchar(100)" json:"payment_reference,omitempty"`

	// Transaction references
	PaymentTransactionID *string `gorm:"type:varchar(36);null" json:"payment_transaction_id,omitempty"`

	// Late fees and penalties (for future enhancement)
	LateFee         float64 `gorm:"type:decimal(15,2);default:0.00" json:"late_fee"`
	PenaltyAmount   float64 `gorm:"type:decimal(15,2);default:0.00" json:"penalty_amount"`
	GracePeriodDays int     `gorm:"type:int;default:0" json:"grace_period_days"`

	// Audit fields
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Plan InstallmentPlan `gorm:"foreignKey:PlanID" json:"plan,omitempty"`
}

// InstallmentPlanAudit represents audit trail for installment plan changes
type InstallmentPlanAudit struct {
	ID     string `gorm:"type:varchar(36);primaryKey" json:"id"`
	PlanID string `gorm:"type:varchar(36);not null;index" json:"plan_id"`
	Action string `gorm:"type:varchar(50);not null" json:"action"`

	// Change tracking
	OldStatus           *string  `gorm:"type:varchar(20)" json:"old_status,omitempty"`
	NewStatus           *string  `gorm:"type:varchar(20)" json:"new_status,omitempty"`
	OldPaidInstallments *int     `gorm:"type:int" json:"old_paid_installments,omitempty"`
	NewPaidInstallments *int     `gorm:"type:int" json:"new_paid_installments,omitempty"`
	OldRemainingAmount  *float64 `gorm:"type:decimal(15,2)" json:"old_remaining_amount,omitempty"`
	NewRemainingAmount  *float64 `gorm:"type:decimal(15,2)" json:"new_remaining_amount,omitempty"`

	// Additional context
	InstallmentID *string  `gorm:"type:varchar(36)" json:"installment_id,omitempty"`
	PaymentAmount *float64 `gorm:"type:decimal(15,2)" json:"payment_amount,omitempty"`
	ChangedBy     string   `gorm:"type:varchar(36);not null" json:"changed_by"`
	ChangeReason  string   `gorm:"type:text" json:"change_reason,omitempty"`
	IPAddress     string   `gorm:"type:varchar(45)" json:"ip_address,omitempty"`
	UserAgent     string   `gorm:"type:text" json:"user_agent,omitempty"`

	// Metadata (JSON for flexible tracking)
	Metadata map[string]interface{} `gorm:"type:json" json:"metadata,omitempty"`

	// Timestamp
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName returns the table name for the InstallmentPlan model
func (InstallmentPlan) TableName() string {
	return "installment_plans"
}

// TableName returns the table name for the Installment model
func (Installment) TableName() string {
	return "installments"
}

// TableName returns the table name for the InstallmentPlanAudit model
func (InstallmentPlanAudit) TableName() string {
	return "installment_plan_audit"
}

// BeforeCreate is called before creating a new installment plan
func (ip *InstallmentPlan) BeforeCreate(tx *gorm.DB) error {
	if ip.ID == "" {
		ip.ID = uuid.New().String()
	}
	// Set remaining amount to total amount initially
	if ip.RemainingAmount == 0 {
		ip.RemainingAmount = ip.TotalAmount
	}
	return nil
}

// BeforeCreate is called before creating a new installment
func (i *Installment) BeforeCreate(tx *gorm.DB) error {
	if i.ID == "" {
		i.ID = uuid.New().String()
	}
	// Set remaining amount to full amount initially
	if i.RemainingAmount == 0 {
		i.RemainingAmount = i.Amount
	}
	return nil
}

// BeforeCreate is called before creating a new audit record
func (ipa *InstallmentPlanAudit) BeforeCreate(tx *gorm.DB) error {
	if ipa.ID == "" {
		ipa.ID = uuid.New().String()
	}
	return nil
}

// Validate validates the installment plan data
func (ip *InstallmentPlan) Validate() error {
	if ip.TransactionID == "" {
		return &ValidationError{Field: "transaction_id", Message: "transaction ID is required"}
	}
	if ip.CardID == "" {
		return &ValidationError{Field: "card_id", Message: "card ID is required"}
	}
	if ip.UserID == "" {
		return &ValidationError{Field: "user_id", Message: "user ID is required"}
	}
	if ip.TotalAmount <= 0 {
		return &ValidationError{Field: "total_amount", Message: "total amount must be positive"}
	}
	if ip.InstallmentsCount < 1 || ip.InstallmentsCount > 24 {
		return &ValidationError{Field: "installments_count", Message: "installments count must be between 1 and 24"}
	}
	if ip.InstallmentAmount <= 0 {
		return &ValidationError{Field: "installment_amount", Message: "installment amount must be positive"}
	}
	if ip.StartDate.IsZero() {
		return &ValidationError{Field: "start_date", Message: "start date is required"}
	}
	if ip.StartDate.Before(time.Now().AddDate(0, 0, -1)) {
		return &ValidationError{Field: "start_date", Message: "start date cannot be in the past"}
	}

	return nil
}

// Validate validates the installment data
func (i *Installment) Validate() error {
	if i.PlanID == "" {
		return &ValidationError{Field: "plan_id", Message: "plan ID is required"}
	}
	if i.InstallmentNumber <= 0 {
		return &ValidationError{Field: "installment_number", Message: "installment number must be positive"}
	}
	if i.Amount <= 0 {
		return &ValidationError{Field: "amount", Message: "amount must be positive"}
	}
	if i.DueDate.IsZero() {
		return &ValidationError{Field: "due_date", Message: "due date is required"}
	}
	if i.PaidAmount < 0 || i.PaidAmount > i.Amount {
		return &ValidationError{Field: "paid_amount", Message: "paid amount must be between 0 and total amount"}
	}
	if i.RemainingAmount < 0 || i.RemainingAmount > i.Amount {
		return &ValidationError{Field: "remaining_amount", Message: "remaining amount must be between 0 and total amount"}
	}

	return nil
}

// IsActive checks if the installment plan is active
func (ip *InstallmentPlan) IsActive() bool {
	return ip.Status == InstallmentPlanStatusActive
}

// IsCompleted checks if the installment plan is completed
func (ip *InstallmentPlan) IsCompleted() bool {
	return ip.Status == InstallmentPlanStatusCompleted || ip.PaidInstallments >= ip.InstallmentsCount
}

// CanCancel checks if the installment plan can be cancelled
func (ip *InstallmentPlan) CanCancel() bool {
	return ip.Status == InstallmentPlanStatusActive && ip.PaidInstallments == 0
}

// CanSuspend checks if the installment plan can be suspended
func (ip *InstallmentPlan) CanSuspend() bool {
	return ip.Status == InstallmentPlanStatusActive
}

// GetCompletionPercentage calculates the completion percentage
func (ip *InstallmentPlan) GetCompletionPercentage() float64 {
	if ip.InstallmentsCount == 0 {
		return 0
	}
	return float64(ip.PaidInstallments) / float64(ip.InstallmentsCount) * 100
}

// GetNextDueInstallment returns the next installment due for payment
func (ip *InstallmentPlan) GetNextDueInstallment() *Installment {
	for _, installment := range ip.Installments {
		if installment.Status == InstallmentStatusPending || installment.Status == InstallmentStatusOverdue {
			return &installment
		}
	}
	return nil
}

// GetOverdueInstallments returns all overdue installments
func (ip *InstallmentPlan) GetOverdueInstallments() []Installment {
	var overdue []Installment
	for _, installment := range ip.Installments {
		if installment.Status == InstallmentStatusOverdue {
			overdue = append(overdue, installment)
		}
	}
	return overdue
}

// GetOverdueAmount calculates total overdue amount
func (ip *InstallmentPlan) GetOverdueAmount() float64 {
	var total float64
	for _, installment := range ip.Installments {
		if installment.Status == InstallmentStatusOverdue {
			total += installment.RemainingAmount
		}
	}
	return total
}

// GetRemainingInstallments returns count of remaining installments
func (ip *InstallmentPlan) GetRemainingInstallments() int {
	return ip.InstallmentsCount - ip.PaidInstallments
}

// IsOverdue checks if the installment is overdue
func (i *Installment) IsOverdue() bool {
	return i.Status == InstallmentStatusOverdue ||
		(i.Status == InstallmentStatusPending && i.DueDate.Before(time.Now()))
}

// IsPaid checks if the installment is fully paid
func (i *Installment) IsPaid() bool {
	return i.Status == InstallmentStatusPaid
}

// CanPay checks if the installment can be paid
func (i *Installment) CanPay() bool {
	return i.Status == InstallmentStatusPending ||
		i.Status == InstallmentStatusOverdue ||
		i.Status == InstallmentStatusPartial
}

// GetDaysOverdue calculates days overdue (0 if not overdue)
func (i *Installment) GetDaysOverdue() int {
	if !i.IsOverdue() {
		return 0
	}
	return int(time.Since(i.DueDate).Hours() / 24)
}

// GetGracePeriodEnd calculates when grace period ends
func (i *Installment) GetGracePeriodEnd() time.Time {
	return i.DueDate.AddDate(0, 0, i.GracePeriodDays)
}

// IsInGracePeriod checks if installment is in grace period
func (i *Installment) IsInGracePeriod() bool {
	if i.GracePeriodDays == 0 {
		return false
	}
	return time.Now().Before(i.GetGracePeriodEnd())
}

// ProcessPayment processes a payment for the installment
func (i *Installment) ProcessPayment(amount float64, paymentMethod, reference string, transactionID *string) error {
	if !i.CanPay() {
		return &ValidationError{Field: "status", Message: "installment cannot be paid in current status"}
	}

	if amount <= 0 {
		return &ValidationError{Field: "amount", Message: "payment amount must be positive"}
	}

	if amount > i.RemainingAmount {
		return &ValidationError{Field: "amount", Message: "payment amount exceeds remaining amount"}
	}

	// Update payment information
	i.PaidAmount += amount
	i.RemainingAmount -= amount
	i.PaymentMethod = paymentMethod
	i.PaymentReference = reference
	i.PaymentTransactionID = transactionID

	// Update status based on payment
	if i.RemainingAmount <= 0.01 { // Allow for small rounding differences
		i.Status = InstallmentStatusPaid
		i.RemainingAmount = 0
		now := time.Now()
		i.PaidDate = &now
	} else {
		i.Status = InstallmentStatusPartial
	}

	return nil
}

// Cancel marks the installment as cancelled
func (i *Installment) Cancel() error {
	if i.Status == InstallmentStatusPaid {
		return &ValidationError{Field: "status", Message: "cannot cancel paid installment"}
	}

	i.Status = InstallmentStatusCancelled
	return nil
}

// IsValidInstallmentPlanStatus checks if the status is valid
func IsValidInstallmentPlanStatus(status InstallmentPlanStatus) bool {
	switch status {
	case InstallmentPlanStatusActive, InstallmentPlanStatusCompleted,
		InstallmentPlanStatusCancelled, InstallmentPlanStatusSuspended:
		return true
	default:
		return false
	}
}

// IsValidInstallmentStatus checks if the status is valid
func IsValidInstallmentStatus(status InstallmentStatus) bool {
	switch status {
	case InstallmentStatusPending, InstallmentStatusPaid, InstallmentStatusOverdue,
		InstallmentStatusCancelled, InstallmentStatusPartial:
		return true
	default:
		return false
	}
}

// Add new methods to Card entity for installment functionality

// CanCreateInstallmentPlan checks if the card can create installment plans
func (c *Card) CanCreateInstallmentPlan(amount float64) bool {
	if !c.IsActive() {
		return false
	}

	if c.CardType != CardTypeCredit {
		return false
	}

	// Check available credit
	availableCredit := c.GetAvailableBalance()
	return availableCredit >= amount
}

// GetActiveInstallmentPlans returns active installment plans (would be loaded separately)
// This is a placeholder method - actual implementation would load from repository
func (c *Card) GetActiveInstallmentPlans() []InstallmentPlan {
	// This would be populated by the repository when loading card with installment plans
	return []InstallmentPlan{}
}

// GetTotalInstallmentCommitments calculates total committed amount in active installment plans
func (c *Card) GetTotalInstallmentCommitments() float64 {
	var total float64
	plans := c.GetActiveInstallmentPlans()
	for _, plan := range plans {
		if plan.IsActive() {
			total += plan.RemainingAmount
		}
	}
	return total
}

// GetAvailableCreditWithInstallments calculates available credit considering installment commitments
func (c *Card) GetAvailableCreditWithInstallments() float64 {
	if c.CardType != CardTypeCredit || c.CreditLimit == nil {
		return 0
	}

	baseAvailable := c.GetAvailableBalance()
	commitments := c.GetTotalInstallmentCommitments()

	return baseAvailable - commitments
}
