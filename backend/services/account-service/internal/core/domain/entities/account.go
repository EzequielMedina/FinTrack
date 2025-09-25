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
