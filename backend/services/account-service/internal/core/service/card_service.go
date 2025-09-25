package service

import (
	"fmt"
	"time"

	"github.com/fintrack/account-service/internal/core/domain/entities"
	"github.com/fintrack/account-service/internal/core/ports"
	"github.com/fintrack/account-service/internal/infrastructure/entrypoints/handlers/card/dto"
	"github.com/google/uuid"
)

type CardService struct {
	cardRepo    ports.CardRepositoryInterface
	accountRepo ports.AccountRepositoryInterface // To validate account exists
}

func NewCardService(cardRepo ports.CardRepositoryInterface, accountRepo ports.AccountRepositoryInterface) *CardService {
	return &CardService{
		cardRepo:    cardRepo,
		accountRepo: accountRepo,
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
		ClosingDate:     req.ClosingDate,
		DueDate:         req.DueDate,
		EncryptedNumber: req.EncryptedNumber,
		KeyFingerprint:  req.KeyFingerprint,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

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

	// Update fields if provided
	if req.HolderName != "" {
		card.HolderName = req.HolderName
	}
	if req.ExpirationMonth > 0 {
		card.ExpirationMonth = req.ExpirationMonth
	}
	if req.ExpirationYear > 0 {
		card.ExpirationYear = req.ExpirationYear
	}
	if req.Nickname != "" {
		card.Nickname = req.Nickname
	}
	if req.IsDefault != nil {
		card.IsDefault = *req.IsDefault
		if *req.IsDefault {
			if err := s.cardRepo.SetDefaultByAccount(card.AccountID, card.ID); err != nil {
				return nil, fmt.Errorf("failed to set default card: %w", err)
			}
		}
	}

	card.UpdatedAt = time.Now()

	// Validate updated data
	if err := card.Validate(); err != nil {
		return nil, fmt.Errorf("invalid card data: %w", err)
	}

	updatedCard, err := s.cardRepo.Update(card)
	if err != nil {
		return nil, fmt.Errorf("failed to update card: %w", err)
	}

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
