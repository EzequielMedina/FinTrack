package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/fintrack/account-service/internal/core/domain/entities"
	"github.com/fintrack/account-service/internal/core/ports"
	"github.com/fintrack/account-service/internal/infrastructure/entrypoints/handlers/card/dto"
	"github.com/google/uuid"
)

// CardService implements the card service interface
type CardService struct {
	cardRepo           ports.CardRepositoryInterface
	accountRepo        ports.AccountRepositoryInterface
	installmentService ports.InstallmentServiceInterface
}

// NewCardService creates a new card service
func NewCardService(
	cardRepo ports.CardRepositoryInterface,
	accountRepo ports.AccountRepositoryInterface,
	installmentService ports.InstallmentServiceInterface,
) ports.CardServiceInterface {
	return &CardService{
		cardRepo:           cardRepo,
		accountRepo:        accountRepo,
		installmentService: installmentService,
	}
}

// CreateCard creates a new card
func (s *CardService) CreateCard(req *dto.CreateCardRequest) (*entities.Card, error) {
	// Validate account exists
	_, err := s.accountRepo.GetByID(req.AccountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	// Set default due date if not provided and it's a credit card
	dueDate := req.DueDate
	if dueDate == nil && req.CardType == "credit" {
		// Set due date to the 5th of next month by default
		now := time.Now()
		nextMonth := now.AddDate(0, 1, 0)
		defaultDueDate := time.Date(nextMonth.Year(), nextMonth.Month(), 5, 0, 0, 0, 0, nextMonth.Location())
		dueDate = &defaultDueDate
	}

	// Create card entity
	card := &entities.Card{
		ID:              uuid.New().String(),
		AccountID:       req.AccountID,
		CardType:        entities.CardType(req.CardType),
		CardBrand:       entities.CardBrand(req.CardBrand),
		LastFourDigits:  req.LastFourDigits,
		MaskedNumber:    req.MaskedNumber,
		HolderName:      req.HolderName,
		ExpirationMonth: req.ExpirationMonth,
		ExpirationYear:  req.ExpirationYear,
		Status:          entities.CardStatusActive,
		IsDefault:       req.IsDefault,
		Nickname:        req.Nickname,
		Balance:         0.0,
		CreditLimit:     req.CreditLimit,
		ClosingDate:     req.ClosingDate,
		DueDate:         dueDate,
		EncryptedNumber: req.EncryptedNumber,
		KeyFingerprint:  req.KeyFingerprint,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Validate card data
	if err := card.Validate(); err != nil {
		return nil, fmt.Errorf("card validation failed: %w", err)
	}

	// Save card
	savedCard, err := s.cardRepo.Create(card)
	if err != nil {
		return nil, fmt.Errorf("failed to create card: %w", err)
	}

	// Set as default if requested
	if req.IsDefault {
		err = s.cardRepo.SetDefaultByAccount(req.AccountID, savedCard.ID)
		if err != nil {
			// Log warning but don't fail the creation
			fmt.Printf("Warning: failed to set card as default: %v\n", err)
		}
	}

	return savedCard, nil
}

// GetCardByID retrieves a card by ID
func (s *CardService) GetCardByID(cardID string) (*entities.Card, error) {
	return s.cardRepo.GetByID(cardID)
}

// GetCardByIDWithAccount retrieves a card by ID with account preloaded
func (s *CardService) GetCardByIDWithAccount(cardID string) (*entities.Card, error) {
	return s.cardRepo.GetByIDWithAccount(cardID)
}

// GetCardsByAccount retrieves cards for a specific account
func (s *CardService) GetCardsByAccount(accountID string, page, pageSize int) ([]*entities.Card, int64, error) {
	offset := (page - 1) * pageSize
	return s.cardRepo.GetByAccount(accountID, pageSize, offset)
}

// GetCardsByUser retrieves cards for a specific user
func (s *CardService) GetCardsByUser(userID string, page, pageSize int) ([]*entities.Card, int64, error) {
	offset := (page - 1) * pageSize
	return s.cardRepo.GetByUser(userID, pageSize, offset)
}

// UpdateCard updates an existing card
func (s *CardService) UpdateCard(cardID string, req *dto.UpdateCardRequest) (*entities.Card, error) {
	card, err := s.cardRepo.GetByID(cardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get card: %w", err)
	}

	// Update fields
	if req.HolderName != "" {
		card.HolderName = req.HolderName
	}
	if req.ExpirationMonth != 0 {
		card.ExpirationMonth = req.ExpirationMonth
	}
	if req.ExpirationYear != 0 {
		card.ExpirationYear = req.ExpirationYear
	}
	if req.Nickname != "" {
		card.Nickname = req.Nickname
	}
	if req.IsDefault != nil {
		card.IsDefault = *req.IsDefault
	}

	card.UpdatedAt = time.Now()

	// Validate updated card
	if err := card.Validate(); err != nil {
		return nil, fmt.Errorf("card validation failed: %w", err)
	}

	return s.cardRepo.Update(card)
}

// DeleteCard deletes a card
func (s *CardService) DeleteCard(cardID string) error {
	return s.cardRepo.Delete(cardID)
}

// BlockCard blocks a card
func (s *CardService) BlockCard(cardID string) (*entities.Card, error) {
	card, err := s.cardRepo.GetByID(cardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get card: %w", err)
	}

	card.Status = entities.CardStatusBlocked
	card.UpdatedAt = time.Now()

	return s.cardRepo.Update(card)
}

// UnblockCard unblocks a card
func (s *CardService) UnblockCard(cardID string) (*entities.Card, error) {
	card, err := s.cardRepo.GetByID(cardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get card: %w", err)
	}

	card.Status = entities.CardStatusActive
	card.UpdatedAt = time.Now()

	return s.cardRepo.Update(card)
}

// SetDefaultCard sets a card as default
func (s *CardService) SetDefaultCard(cardID string) (*entities.Card, error) {
	card, err := s.cardRepo.GetByIDWithAccount(cardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get card: %w", err)
	}

	err = s.cardRepo.SetDefaultByAccount(card.AccountID, cardID)
	if err != nil {
		return nil, fmt.Errorf("failed to set default card: %w", err)
	}

	// Return updated card
	return s.cardRepo.GetByID(cardID)
}

// ChargeCard charges a credit card (without installments)
func (s *CardService) ChargeCard(cardID string, amount float64, description, reference string) (*entities.Card, error) {
	card, err := s.cardRepo.GetByIDWithAccount(cardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get card: %w", err)
	}

	if card.CardType != entities.CardTypeCredit {
		return nil, errors.New("only credit cards can be charged")
	}

	if !card.IsActive() {
		return nil, errors.New("card is not active")
	}

	// Check credit limit
	if !card.CanCharge(amount) {
		return nil, errors.New("charge amount exceeds available credit limit")
	}

	// Process charge
	oldBalance := card.Balance
	card.Balance += amount
	card.UpdatedAt = time.Now()

	updatedCard, err := s.cardRepo.Update(card)
	if err != nil {
		return nil, fmt.Errorf("failed to update card balance: %w", err)
	}

	// Here you would typically create a transaction record
	// For now, we'll just log the operation
	fmt.Printf("Credit card charged: CardID=%s, Amount=%.2f, OldBalance=%.2f, NewBalance=%.2f\n",
		cardID, amount, oldBalance, card.Balance)

	return updatedCard, nil
}

// ChargeCardWithInstallments charges a credit card with installment plan
func (s *CardService) ChargeCardWithInstallments(req *dto.CreateInstallmentPlanRequest) (*dto.ChargeWithInstallmentsResponse, error) {
	// Validate card exists and can be charged
	card, err := s.cardRepo.GetByIDWithAccount(req.CardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get card: %w", err)
	}

	if card.CardType != entities.CardTypeCredit {
		return nil, errors.New("only credit cards support installments")
	}

	if !card.IsActive() {
		return nil, errors.New("card is not active")
	}

	if card.Account.UserID != req.UserID {
		return nil, errors.New("card does not belong to user")
	}

	// Check if card can create installment plan
	if !card.CanCreateInstallmentPlan(req.TotalAmount) {
		return nil, errors.New("card cannot create installment plan with this amount")
	}

	// Create installment plan
	installmentPlan, err := s.installmentService.CreateInstallmentPlan(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create installment plan: %w", err)
	}

	// Charge the card for the total amount immediately
	fmt.Printf("DEBUG - About to charge card %s with total amount %.2f\n", req.CardID, req.TotalAmount)
	_, err = s.ChargeCard(req.CardID, req.TotalAmount,
		fmt.Sprintf("Purchase with %d installments - %s", installmentPlan.InstallmentsCount, req.Description),
		req.Reference)
	if err != nil {
		fmt.Printf("DEBUG - Card charge failed: %v\n", err)
		// Try to cleanup the installment plan if card charge fails
		_, cancelErr := s.installmentService.CancelInstallmentPlan(installmentPlan.ID,
			"Card charge failed for purchase", req.InitiatedBy)
		if cancelErr != nil {
			fmt.Printf("Warning: failed to cancel installment plan after charge failure: %v\n", cancelErr)
		}
		return nil, fmt.Errorf("failed to charge card for purchase: %w", err)
	}
	fmt.Printf("DEBUG - Card charged successfully\n")
	firstInstallmentCharged := true

	// Get updated card
	updatedCard, err := s.cardRepo.GetByID(req.CardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated card: %w", err)
	}

	return &dto.ChargeWithInstallmentsResponse{
		InstallmentPlan:         installmentPlan,
		Card:                    updatedCard,
		FirstInstallmentCharged: firstInstallmentCharged,
		TransactionID:           "", // Would be set when transaction service is integrated
	}, nil
}

// PaymentCard processes a payment on a credit card
func (s *CardService) PaymentCard(cardID string, amount float64, paymentMethod, reference string) (*entities.Card, error) {
	card, err := s.cardRepo.GetByID(cardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get card: %w", err)
	}

	if card.CardType != entities.CardTypeCredit {
		return nil, errors.New("only credit cards accept payments")
	}

	if amount <= 0 {
		return nil, errors.New("payment amount must be greater than 0")
	}

	if amount > card.Balance {
		return nil, errors.New("payment amount cannot exceed current balance")
	}

	// Process payment
	oldBalance := card.Balance
	card.Balance -= amount
	card.UpdatedAt = time.Now()

	updatedCard, err := s.cardRepo.Update(card)
	if err != nil {
		return nil, fmt.Errorf("failed to update card balance: %w", err)
	}

	// Here you would typically create a payment transaction record
	fmt.Printf("Credit card payment processed: CardID=%s, Amount=%.2f, OldBalance=%.2f, NewBalance=%.2f\n",
		cardID, amount, oldBalance, card.Balance)

	return updatedCard, nil
}

// ProcessDebitTransaction processes a debit transaction
func (s *CardService) ProcessDebitTransaction(cardID string, amount float64, description, merchantName, reference string) (*entities.Card, error) {
	card, err := s.cardRepo.GetByIDWithAccount(cardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get card: %w", err)
	}

	if card.CardType != entities.CardTypeDebit {
		return nil, errors.New("only debit cards can process debit transactions")
	}

	if !card.IsActive() {
		return nil, errors.New("card is not active")
	}

	// Check account balance
	if card.Account.Balance < amount {
		return nil, errors.New("insufficient account balance")
	}

	// Process debit transaction by reducing account balance
	account := card.Account
	account.Balance -= amount
	account.UpdatedAt = time.Now()

	err = s.accountRepo.Update(&account)
	if err != nil {
		return nil, fmt.Errorf("failed to update account balance: %w", err)
	}

	// Here you would typically create a transaction record
	fmt.Printf("Debit transaction processed: CardID=%s, Amount=%.2f, Account=%s, NewBalance=%.2f\n",
		cardID, amount, card.AccountID, account.Balance)

	// Return updated card with account info
	return s.cardRepo.GetByIDWithAccount(cardID)
}
