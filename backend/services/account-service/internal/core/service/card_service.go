package service

import (
	"fmt"
	"time"

	"github.com/fintrack/account-service/internal/core/domain/entities"
	"github.com/fintrack/account-service/internal/core/ports"
	"github.com/fintrack/account-service/internal/infrastructure/clients"
	"github.com/fintrack/account-service/internal/infrastructure/entrypoints/handlers/card/dto"
	"github.com/google/uuid"
)

type CardService struct {
	cardRepo           ports.CardRepositoryInterface
	accountRepo        ports.AccountRepositoryInterface  // To validate account exists
	installmentService ports.InstallmentServiceInterface // To handle installment plans
	transactionClient  *clients.TransactionClient        // To record transactions
}

func NewCardService(cardRepo ports.CardRepositoryInterface, accountRepo ports.AccountRepositoryInterface, installmentService ports.InstallmentServiceInterface) *CardService {
	return &CardService{
		cardRepo:           cardRepo,
		accountRepo:        accountRepo,
		installmentService: installmentService,
		transactionClient:  clients.NewTransactionClient(),
	}
}

func (s *CardService) CreateCard(req *dto.CreateCardRequest) (*entities.Card, error) {
	// Validate that the account exists and can have cards
	account, err := s.accountRepo.GetByID(req.AccountID)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}

	if !account.CanHaveCards() {
		return nil, fmt.Errorf("account type %s cannot have cards", account.AccountType)
	}

	// Set default due date if not provided and it's a credit card
	var dueDate *time.Time
	if req.DueDate != nil {
		dueDate = req.DueDate.ToTimePointer()
		fmt.Printf("🗓️ DEBUG - Due date provided: %s\n", dueDate.Format("2006-01-02"))
	} else if req.CardType == "credit" {
		// Set due date to the 5th of next month by default
		now := time.Now()
		nextMonth := now.AddDate(0, 1, 0)
		defaultDueDate := time.Date(nextMonth.Year(), nextMonth.Month(), 5, 0, 0, 0, 0, nextMonth.Location())
		dueDate = &defaultDueDate
		fmt.Printf("🗓️ DEBUG - No due date provided for credit card, setting default to: %s\n", defaultDueDate.Format("2006-01-02"))
	}

	// Extract dates from CustomDate types
	var closingDate *time.Time
	if req.ClosingDate != nil {
		closingDate = req.ClosingDate.ToTimePointer()
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
		CreditLimit:     req.CreditLimit,
		ClosingDate:     closingDate,
		DueDate:         dueDate,
		EncryptedNumber: req.EncryptedNumber,
		KeyFingerprint:  req.KeyFingerprint,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	fmt.Printf("🃏 DEBUG - Creating card with DueDate: %v\n", dueDate)

	// Validate card data
	if err := card.Validate(); err != nil {
		return nil, fmt.Errorf("invalid card data: %w", err)
	}

	// If this is the first card or marked as default, set as default
	if req.IsDefault {
		if err := s.cardRepo.SetDefaultByAccount(req.AccountID, card.ID); err != nil {
			return nil, fmt.Errorf("failed to set default card: %w", err)
		}
	}

	// Create the card
	createdCard, err := s.cardRepo.Create(card)
	if err != nil {
		return nil, fmt.Errorf("failed to create card: %w", err)
	}

	return createdCard, nil
}

func (s *CardService) GetCardByID(cardID string) (*entities.Card, error) {
	card, err := s.cardRepo.GetByID(cardID)
	if err != nil {
		return nil, fmt.Errorf("card not found: %w", err)
	}
	return card, nil
}

func (s *CardService) GetCardsByAccount(accountID string, page, pageSize int) ([]*entities.Card, int64, error) {
	// Validate account exists
	_, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		return nil, 0, fmt.Errorf("account not found: %w", err)
	}

	offset := (page - 1) * pageSize
	cards, total, err := s.cardRepo.GetByAccount(accountID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get cards by account: %w", err)
	}

	return cards, total, nil
}

func (s *CardService) GetCardsByUser(userID string, page, pageSize int) ([]*entities.Card, int64, error) {
	offset := (page - 1) * pageSize
	cards, total, err := s.cardRepo.GetByUser(userID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get cards by user: %w", err)
	}

	return cards, total, nil
}

func (s *CardService) UpdateCard(cardID string, req *dto.UpdateCardRequest) (*entities.Card, error) {
	// Get existing card
	card, err := s.cardRepo.GetByID(cardID)
	if err != nil {
		return nil, fmt.Errorf("card not found: %w", err)
	}

	fmt.Printf("🔄 DEBUG - UpdateCard called for card ID: %s\n", cardID)
	fmt.Printf("🔄 DEBUG - Current card data: HolderName=%s, ExpirationMonth=%d, ExpirationYear=%d, Nickname=%s, IsDefault=%t\n",
		card.HolderName, card.ExpirationMonth, card.ExpirationYear, card.Nickname, card.IsDefault)
	fmt.Printf("🔄 DEBUG - Update request: HolderName=%s, ExpirationMonth=%d, ExpirationYear=%d, Nickname=%s, IsDefault=%v, CreditLimit=%v\n",
		req.HolderName, req.ExpirationMonth, req.ExpirationYear, req.Nickname, req.IsDefault, req.CreditLimit)

	// Track if any updates were made
	updated := false

	// Update fields if provided (check for non-empty strings and non-zero values)
	if req.HolderName != "" && req.HolderName != card.HolderName {
		fmt.Printf("🔄 DEBUG - Updating HolderName from '%s' to '%s'\n", card.HolderName, req.HolderName)
		card.HolderName = req.HolderName
		updated = true
	}

	// For expiration month, only update if it's a valid month and different
	if req.ExpirationMonth > 0 && req.ExpirationMonth != card.ExpirationMonth {
		if req.ExpirationMonth >= 1 && req.ExpirationMonth <= 12 {
			fmt.Printf("🔄 DEBUG - Updating ExpirationMonth from %d to %d\n", card.ExpirationMonth, req.ExpirationMonth)
			card.ExpirationMonth = req.ExpirationMonth
			updated = true
		} else {
			return nil, fmt.Errorf("invalid expiration month: must be between 1 and 12")
		}
	}

	// For expiration year, only update if it's provided and different
	if req.ExpirationYear > 0 && req.ExpirationYear != card.ExpirationYear {
		if req.ExpirationYear >= 2020 { // Allow reasonable range
			fmt.Printf("🔄 DEBUG - Updating ExpirationYear from %d to %d\n", card.ExpirationYear, req.ExpirationYear)
			card.ExpirationYear = req.ExpirationYear
			updated = true
		} else {
			return nil, fmt.Errorf("invalid expiration year: must be 2020 or later")
		}
	}

	// Update nickname (can be empty string to clear it)
	if req.Nickname != card.Nickname {
		fmt.Printf("🔄 DEBUG - Updating Nickname from '%s' to '%s'\n", card.Nickname, req.Nickname)
		card.Nickname = req.Nickname
		updated = true
	}

	// Handle IsDefault field - use pointer to distinguish between not provided and false
	if req.IsDefault != nil {
		if *req.IsDefault != card.IsDefault {
			fmt.Printf("🔄 DEBUG - Updating IsDefault from %t to %t\n", card.IsDefault, *req.IsDefault)
			card.IsDefault = *req.IsDefault
			updated = true

			// If setting as default, ensure no other cards are default for this account
			if *req.IsDefault {
				if err := s.cardRepo.SetDefaultByAccount(card.AccountID, card.ID); err != nil {
					return nil, fmt.Errorf("failed to set default card: %w", err)
				}
			}
		}
	}

	// Handle CreditLimit field - only for credit cards
	if req.CreditLimit != nil {
		// Verify this is a credit card
		if card.CardType != entities.CardTypeCredit {
			return nil, fmt.Errorf("credit limit can only be updated for credit cards")
		}

		// Validate minimum credit limit
		if *req.CreditLimit < 0 {
			return nil, fmt.Errorf("credit limit cannot be negative")
		}

		// Check if it's actually changing
		currentLimit := float64(0)
		if card.CreditLimit != nil {
			currentLimit = *card.CreditLimit
		}

		if *req.CreditLimit != currentLimit {
			fmt.Printf("🔄 DEBUG - Updating CreditLimit from %.2f to %.2f\n", currentLimit, *req.CreditLimit)

			// Validate against current balance
			if *req.CreditLimit < card.Balance {
				return nil, fmt.Errorf("credit limit (%.2f) cannot be lower than current balance (%.2f)", *req.CreditLimit, card.Balance)
			}

			card.CreditLimit = req.CreditLimit
			updated = true
		}
	}

	// If no updates were made, return the existing card
	if !updated {
		fmt.Printf("🔄 DEBUG - No changes detected, returning existing card\n")
		return card, nil
	}

	// Update timestamp
	card.UpdatedAt = time.Now()

	// Validate updated data using the update-specific validation
	if err := card.ValidateForUpdate(); err != nil {
		return nil, fmt.Errorf("invalid card data: %w", err)
	}

	// Save updates to database
	fmt.Printf("🔄 DEBUG - Saving card updates to database\n")
	updatedCard, err := s.cardRepo.Update(card)
	if err != nil {
		return nil, fmt.Errorf("failed to update card: %w", err)
	}

	fmt.Printf("🔄 DEBUG - Card updated successfully\n")
	return updatedCard, nil
}

func (s *CardService) DeleteCard(cardID string) error {
	// Check if card exists
	_, err := s.cardRepo.GetByID(cardID)
	if err != nil {
		return fmt.Errorf("card not found: %w", err)
	}

	err = s.cardRepo.Delete(cardID)
	if err != nil {
		return fmt.Errorf("failed to delete card: %w", err)
	}

	return nil
}

func (s *CardService) BlockCard(cardID string) (*entities.Card, error) {
	card, err := s.cardRepo.GetByID(cardID)
	if err != nil {
		return nil, fmt.Errorf("card not found: %w", err)
	}

	card.Status = entities.CardStatusBlocked
	card.UpdatedAt = time.Now()

	updatedCard, err := s.cardRepo.Update(card)
	if err != nil {
		return nil, fmt.Errorf("failed to block card: %w", err)
	}

	return updatedCard, nil
}

func (s *CardService) UnblockCard(cardID string) (*entities.Card, error) {
	card, err := s.cardRepo.GetByID(cardID)
	if err != nil {
		return nil, fmt.Errorf("card not found: %w", err)
	}

	card.Status = entities.CardStatusActive
	card.UpdatedAt = time.Now()

	updatedCard, err := s.cardRepo.Update(card)
	if err != nil {
		return nil, fmt.Errorf("failed to unblock card: %w", err)
	}

	return updatedCard, nil
}

func (s *CardService) SetDefaultCard(cardID string) (*entities.Card, error) {
	card, err := s.cardRepo.GetByID(cardID)
	if err != nil {
		return nil, fmt.Errorf("card not found: %w", err)
	}

	if err := s.cardRepo.SetDefaultByAccount(card.AccountID, cardID); err != nil {
		return nil, fmt.Errorf("failed to set default card: %w", err)
	}

	// Get updated card
	updatedCard, err := s.cardRepo.GetByID(cardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated card: %w", err)
	}

	return updatedCard, nil
}

// GetCardByIDWithAccount gets a card by ID with account preloaded
func (s *CardService) GetCardByIDWithAccount(cardID string) (*entities.Card, error) {
	card, err := s.cardRepo.GetByIDWithAccount(cardID)
	if err != nil {
		return nil, fmt.Errorf("card not found: %w", err)
	}
	return card, nil
}

// CREDIT CARD FINANCIAL OPERATIONS

// ChargeCard processes a charge to a credit card
func (s *CardService) ChargeCard(cardID string, amount float64, description, reference string) (*entities.Card, error) {
	// Get card with account data
	card, err := s.cardRepo.GetByIDWithAccount(cardID)
	if err != nil {
		return nil, fmt.Errorf("card not found: %w", err)
	}

	// Validate that it's a credit card
	if card.CardType != entities.CardTypeCredit {
		return nil, fmt.Errorf("charges can only be made to credit cards")
	}

	// Use the business logic from the entity
	if err := card.Charge(amount); err != nil {
		return nil, fmt.Errorf("failed to charge card: %w", err)
	}

	// Save updated card
	updatedCard, err := s.cardRepo.Update(card)
	if err != nil {
		return nil, fmt.Errorf("failed to save card charge: %w", err)
	}

	return updatedCard, nil
}

// PaymentCard processes a payment to a credit card
func (s *CardService) PaymentCard(cardID string, amount float64, paymentMethod, reference string) (*entities.Card, error) {
	// Get card
	card, err := s.cardRepo.GetByID(cardID)
	if err != nil {
		return nil, fmt.Errorf("card not found: %w", err)
	}

	// Use the business logic from the entity
	if err := card.Payment(amount); err != nil {
		return nil, fmt.Errorf("failed to process payment: %w", err)
	}

	// Save updated card
	updatedCard, err := s.cardRepo.Update(card)
	if err != nil {
		return nil, fmt.Errorf("failed to save card payment: %w", err)
	}

	return updatedCard, nil
}

// DEBIT CARD OPERATIONS

// ProcessDebitTransaction processes a transaction with a debit card
func (s *CardService) ProcessDebitTransaction(cardID string, amount float64, description, merchantName, reference string) (*entities.Card, error) {
	// Get card with account data
	card, err := s.cardRepo.GetByIDWithAccount(cardID)
	if err != nil {
		return nil, fmt.Errorf("card not found: %w", err)
	}

	// Validate that it's a debit card
	if card.CardType != entities.CardTypeDebit {
		return nil, fmt.Errorf("transactions can only be made with debit cards")
	}

	// Use the business logic from the entity
	if err := card.Charge(amount); err != nil {
		return nil, fmt.Errorf("failed to process transaction: %w", err)
	}

	// Save updated account (debit cards deduct from account balance)
	if err := s.accountRepo.Update(&card.Account); err != nil {
		return nil, fmt.Errorf("failed to update account balance: %w", err)
	}

	// Record transaction in transaction service (async, don't fail if this fails)
	go func() {
		userID := card.Account.UserID // Assuming account has UserID field
		if err := s.transactionClient.CreateDebitCardTransaction(
			userID,
			card.Account.ID,
			cardID,
			amount,
			description,
			merchantName,
			reference,
		); err != nil {
			// Log error but don't fail the transaction
			fmt.Printf("Warning: failed to record transaction in transaction service: %v\n", err)
		}
	}()

	// Get updated card with new account balance
	updatedCard, err := s.cardRepo.GetByIDWithAccount(cardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated card: %w", err)
	}

	// Record transaction in transaction service (async, only for logging - does not modify balance)
	go func() {
		// Use a default user ID - in a real implementation, this should come from the request context
		userID := "system"
		err := s.transactionClient.CreateDebitCardTransaction(
			userID,
			card.Account.ID,
			cardID,
			amount,
			description,
			merchantName,
			reference,
		)
		if err != nil {
			// Log error but don't fail the main transaction
			fmt.Printf("Warning: Failed to record transaction in transaction service: %v\n", err)
		}
	}()

	return updatedCard, nil
}

// ChargeCardWithInstallments processes a credit card charge with installment plan
func (s *CardService) ChargeCardWithInstallments(req *dto.CreateInstallmentPlanRequest) (*dto.ChargeWithInstallmentsResponse, error) {
	fmt.Printf("🚨🚨🚨 DEBUG - ChargeCardWithInstallments called with CardID: %s, TotalAmount: %.2f 🚨🚨🚨\n", req.CardID, req.TotalAmount)

	// Verificar que la tarjeta existe y obtener información con cuenta
	card, err := s.cardRepo.GetByIDWithAccount(req.CardID)
	if err != nil {
		return nil, fmt.Errorf("card not found: %w", err)
	}

	if card.CardType != "credit" {
		return nil, fmt.Errorf("installment plans are only available for credit cards")
	}

	// Crear el plan de cuotas usando InstallmentService
	installmentPlan, err := s.installmentService.CreateInstallmentPlan(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create installment plan: %w", err)
	}

	// Cargar el monto total inmediatamente
	fmt.Printf("DEBUG - About to charge card %s with total amount %.2f\n", req.CardID, req.TotalAmount)
	chargedCard, err := s.ChargeCard(req.CardID, req.TotalAmount,
		fmt.Sprintf("Purchase with %d installments - %s", installmentPlan.InstallmentsCount, req.Description),
		req.Reference)
	if err != nil {
		fmt.Printf("DEBUG - Card charge failed: %v\n", err)
		// Tratar de cancelar el plan de cuotas si falla el cargo de tarjeta
		_, cancelErr := s.installmentService.CancelInstallmentPlan(installmentPlan.ID,
			"Card charge failed for purchase", req.InitiatedBy)
		if cancelErr != nil {
			fmt.Printf("Warning: failed to cancel installment plan after charge failure: %v\n", cancelErr)
		}
		return nil, fmt.Errorf("failed to charge card for purchase: %w", err)
	}
	fmt.Printf("DEBUG - Card charged successfully with new balance: %.2f\n", chargedCard.Balance)
	firstInstallmentCharged := true

	return &dto.ChargeWithInstallmentsResponse{
		InstallmentPlan:         installmentPlan,
		Card:                    chargedCard, // Usar la tarjeta actualizada después del cargo
		FirstInstallmentCharged: firstInstallmentCharged,
		TransactionID:           installmentPlan.TransactionID,
	}, nil
}
