package card

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/fintrack/account-service/internal/core/domain/entities"
	"github.com/fintrack/account-service/internal/core/errors"
	"github.com/fintrack/account-service/internal/core/ports"
	"github.com/fintrack/account-service/internal/infrastructure/entrypoints/handlers/card/dto"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	cardService ports.CardServiceInterface
}

func New(cardService ports.CardServiceInterface) *Handler {
	return &Handler{
		cardService: cardService,
	}
}

// CreateCard creates a new card for an account
// @Summary Create a new card
// @Description Create a new card associated with an account
// @Tags Cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Account ID"
// @Param card body dto.CreateCardRequest true "Card creation data"
// @Success 201 {object} dto.CardResponse "Card created successfully"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 404 {object} map[string]string "Account not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/accounts/{id}/cards [post]
func (h *Handler) CreateCard(c *gin.Context) {
	accountID := c.Param("id")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account ID is required"})
		return
	}

	var req dto.CreateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the account ID from the URL
	req.AccountID = accountID

	card, err := h.cardService.CreateCard(&req)
	if err != nil {
		status := h.getErrorStatus(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToCardResponse(card)
	c.JSON(http.StatusCreated, response)
}

// GetCardsByAccount gets all cards for an account
// @Summary Get cards by account
// @Description Retrieve all cards associated with an account
// @Tags Cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Account ID"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Success 200 {object} dto.PaginatedCardResponse "Cards retrieved successfully"
// @Failure 400 {object} map[string]string "Invalid request parameters"
// @Failure 404 {object} map[string]string "Account not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/accounts/{id}/cards [get]
func (h *Handler) GetCardsByAccount(c *gin.Context) {
	accountID := c.Param("id")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account ID is required"})
		return
	}

	page, pageSize := h.getPaginationParams(c)

	cards, total, err := h.cardService.GetCardsByAccount(accountID, page, pageSize)
	if err != nil {
		status := h.getErrorStatus(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToPaginatedCardResponse(cards, total, page, pageSize)
	c.JSON(http.StatusOK, response)
}

// GetCardsByUser gets all cards for a user across all accounts
// @Summary Get cards by user
// @Description Retrieve all cards owned by a user across all their accounts
// @Tags Cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path string true "User ID"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Success 200 {object} dto.PaginatedCardResponse "Cards retrieved successfully"
// @Failure 400 {object} map[string]string "Invalid request parameters"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/accounts/user/{userId}/cards [get]
func (h *Handler) GetCardsByUser(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	page, pageSize := h.getPaginationParams(c)

	cards, total, err := h.cardService.GetCardsByUser(userID, page, pageSize)
	if err != nil {
		status := h.getErrorStatus(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToPaginatedCardResponse(cards, total, page, pageSize)
	c.JSON(http.StatusOK, response)
}

// GetCard gets a specific card
// @Summary Get card by ID
// @Description Retrieve a specific card by its ID
// @Tags Cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Account ID"
// @Param cardId path string true "Card ID"
// @Success 200 {object} dto.CardResponse "Card retrieved successfully"
// @Failure 400 {object} map[string]string "Invalid card ID"
// @Failure 404 {object} map[string]string "Card not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/accounts/{id}/cards/{cardId} [get]
func (h *Handler) GetCard(c *gin.Context) {
	cardID := c.Param("cardId")
	if cardID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "card ID is required"})
		return
	}

	card, err := h.cardService.GetCardByID(cardID)
	if err != nil {
		status := h.getErrorStatus(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToCardResponse(card)
	c.JSON(http.StatusOK, response)
}

// UpdateCard updates a card
// @Summary Update card
// @Description Update card information
// @Tags Cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Account ID"
// @Param cardId path string true "Card ID"
// @Param card body dto.UpdateCardRequest true "Card update data"
// @Success 200 {object} dto.CardResponse "Card updated successfully"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 404 {object} map[string]string "Card not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/accounts/{id}/cards/{cardId} [put]
func (h *Handler) UpdateCard(c *gin.Context) {
	cardID := c.Param("cardId")
	if cardID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "card ID is required"})
		return
	}

	var req dto.UpdateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	card, err := h.cardService.UpdateCard(cardID, &req)
	if err != nil {
		status := h.getErrorStatus(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToCardResponse(card)
	c.JSON(http.StatusOK, response)
}

// DeleteCard deletes a card
// @Summary Delete card
// @Description Delete a card
// @Tags Cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Account ID"
// @Param cardId path string true "Card ID"
// @Success 200 {object} map[string]string "Card deleted successfully"
// @Failure 400 {object} map[string]string "Invalid card ID"
// @Failure 404 {object} map[string]string "Card not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/accounts/{id}/cards/{cardId} [delete]
func (h *Handler) DeleteCard(c *gin.Context) {
	cardID := c.Param("cardId")
	if cardID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "card ID is required"})
		return
	}

	err := h.cardService.DeleteCard(cardID)
	if err != nil {
		status := h.getErrorStatus(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "card deleted successfully"})
}

// BlockCard blocks a card
// @Summary Block card
// @Description Block a card to prevent its usage
// @Tags Cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Account ID"
// @Param cardId path string true "Card ID"
// @Success 200 {object} dto.CardResponse "Card blocked successfully"
// @Failure 400 {object} map[string]string "Invalid card ID"
// @Failure 404 {object} map[string]string "Card not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/accounts/{id}/cards/{cardId}/block [put]
func (h *Handler) BlockCard(c *gin.Context) {
	cardID := c.Param("cardId")
	if cardID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "card ID is required"})
		return
	}

	card, err := h.cardService.BlockCard(cardID)
	if err != nil {
		status := h.getErrorStatus(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToCardResponse(card)
	c.JSON(http.StatusOK, response)
}

// UnblockCard unblocks a card
// @Summary Unblock card
// @Description Unblock a card to allow its usage
// @Tags Cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Account ID"
// @Param cardId path string true "Card ID"
// @Success 200 {object} dto.CardResponse "Card unblocked successfully"
// @Failure 400 {object} map[string]string "Invalid card ID"
// @Failure 404 {object} map[string]string "Card not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/accounts/{id}/cards/{cardId}/unblock [put]
func (h *Handler) UnblockCard(c *gin.Context) {
	cardID := c.Param("cardId")
	if cardID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "card ID is required"})
		return
	}

	card, err := h.cardService.UnblockCard(cardID)
	if err != nil {
		status := h.getErrorStatus(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToCardResponse(card)
	c.JSON(http.StatusOK, response)
}

// SetDefaultCard sets a card as the default for its account
// @Summary Set default card
// @Description Set a specific card as the default for its account
// @Tags Cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Account ID"
// @Param cardId path string true "Card ID"
// @Success 200 {object} dto.CardResponse "Card set as default successfully"
// @Failure 400 {object} map[string]string "Invalid card ID"
// @Failure 404 {object} map[string]string "Card not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/accounts/{id}/cards/{cardId}/set-default [put]
func (h *Handler) SetDefaultCard(c *gin.Context) {
	cardID := c.Param("cardId")
	if cardID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "card ID is required"})
		return
	}

	card, err := h.cardService.SetDefaultCard(cardID)
	if err != nil {
		status := h.getErrorStatus(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToCardResponse(card)
	c.JSON(http.StatusOK, response)
}

// CREDIT CARD FINANCIAL OPERATIONS

// ChargeCard processes a charge to a credit card
// @Summary Charge credit card
// @Description Process a purchase/charge to a credit card
// @Tags Cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param cardId path string true "Card ID"
// @Param charge body dto.CreditCardChargeRequest true "Charge data"
// @Success 200 {object} dto.CreditCardBalanceResponse "Charge processed successfully"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 404 {object} map[string]string "Card not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/cards/{cardId}/charge [post]
func (h *Handler) ChargeCard(c *gin.Context) {
	cardID := c.Param("cardId")
	if cardID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "card ID is required"})
		return
	}

	var req dto.CreditCardChargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	card, err := h.cardService.ChargeCard(cardID, req.Amount, req.Description, req.Reference)
	if err != nil {
		status := h.getErrorStatus(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToCreditCardBalanceResponse(card)
	c.JSON(http.StatusOK, response)
}

// PaymentCard processes a payment to a credit card
// @Summary Pay credit card
// @Description Process a payment to reduce credit card debt
// @Tags Cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param cardId path string true "Card ID"
// @Param payment body dto.CreditCardPaymentRequest true "Payment data"
// @Success 200 {object} dto.CreditCardBalanceResponse "Payment processed successfully"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 404 {object} map[string]string "Card not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/cards/{cardId}/payment [post]
func (h *Handler) PaymentCard(c *gin.Context) {
	cardID := c.Param("cardId")
	if cardID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "card ID is required"})
		return
	}

	var req dto.CreditCardPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	card, err := h.cardService.PaymentCard(cardID, req.Amount, req.PaymentMethod, req.Reference)
	if err != nil {
		status := h.getErrorStatus(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToCreditCardBalanceResponse(card)
	c.JSON(http.StatusOK, response)
}

// GetCardBalance gets the balance information for a card
// @Summary Get card balance
// @Description Get balance information for either credit or debit card
// @Tags Cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param cardId path string true "Card ID"
// @Success 200 {object} dto.CreditCardBalanceResponse "Credit card balance"
// @Success 200 {object} dto.DebitCardBalanceResponse "Debit card balance"
// @Failure 400 {object} map[string]string "Invalid card ID"
// @Failure 404 {object} map[string]string "Card not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/cards/{cardId}/balance [get]
func (h *Handler) GetCardBalance(c *gin.Context) {
	cardID := c.Param("cardId")
	if cardID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "card ID is required"})
		return
	}

	card, err := h.cardService.GetCardByIDWithAccount(cardID)
	if err != nil {
		status := h.getErrorStatus(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// Return different response based on card type
	if card.CardType == entities.CardTypeCredit {
		response := dto.ToCreditCardBalanceResponse(card)
		c.JSON(http.StatusOK, response)
	} else {
		response := dto.ToDebitCardBalanceResponse(card)
		c.JSON(http.StatusOK, response)
	}
}

// DEBIT CARD OPERATIONS

// ProcessDebitTransaction processes a transaction with a debit card
// @Summary Process debit transaction
// @Description Process a purchase/withdrawal with a debit card
// @Tags Cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param cardId path string true "Card ID"
// @Param transaction body dto.DebitCardTransactionRequest true "Transaction data"
// @Success 200 {object} dto.DebitCardBalanceResponse "Transaction processed successfully"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 404 {object} map[string]string "Card not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/cards/{cardId}/transaction [post]
func (h *Handler) ProcessDebitTransaction(c *gin.Context) {
	cardID := c.Param("cardId")
	if cardID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "card ID is required"})
		return
	}

	var req dto.DebitCardTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	card, err := h.cardService.ProcessDebitTransaction(cardID, req.Amount, req.Description, req.MerchantName, req.Reference)
	if err != nil {
		status := h.getErrorStatus(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToDebitCardBalanceResponse(card)
	c.JSON(http.StatusOK, response)
}

// Helper methods
func (h *Handler) getPaginationParams(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return page, pageSize
}

func (h *Handler) getErrorStatus(err error) int {
	// Handle specific domain errors
	if errors.IsNotFoundError(err) {
		return http.StatusNotFound
	}

	if errors.IsValidationError(err) {
		return http.StatusBadRequest
	}

	if errors.IsBusinessLogicError(err) {
		return http.StatusBadRequest
	}

	if errors.IsPermissionError(err) {
		return http.StatusForbidden
	}

	// Handle errors by message content
	errMsg := strings.ToLower(err.Error())
	if strings.Contains(errMsg, "not found") {
		return http.StatusNotFound
	}

	if strings.Contains(errMsg, "cannot have cards") {
		return http.StatusBadRequest
	}

	if strings.Contains(errMsg, "invalid") || strings.Contains(errMsg, "validation") {
		return http.StatusBadRequest
	}

	if strings.Contains(errMsg, "unauthorized") || strings.Contains(errMsg, "forbidden") {
		return http.StatusForbidden
	}

	// Default to internal server error for unhandled cases
	return http.StatusInternalServerError
}

// ChargeCardWithInstallments charges a credit card with installment plan
// @Summary Charge card with installments
// @Description Create a purchase on a credit card with installment plan
// @Tags Cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Card ID"
// @Param charge body dto.CreditCardChargeWithInstallmentsRequest true "Charge with installments data"
// @Success 201 {object} dto.ChargeWithInstallmentsResponse "Purchase created with installments"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 404 {object} map[string]string "Card not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/cards/{id}/charge-installments [post]
func (h *Handler) ChargeCardWithInstallments(c *gin.Context) {
	cardID := c.Param("id")
	if cardID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "card ID is required"})
		return
	}

	var req dto.CreditCardChargeWithInstallmentsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert to CreateInstallmentPlanRequest
	installmentReq := &dto.CreateInstallmentPlanRequest{
		CardID:            cardID,
		TotalAmount:       req.Amount,
		InstallmentsCount: req.InstallmentsCount,
		StartDate:         req.StartDate,
		Description:       req.Description,
		MerchantName:      req.MerchantName,
		MerchantID:        req.MerchantID,
		InterestRate:      req.InterestRate,
		AdminFee:          req.AdminFee,
		Reference:         req.Reference,
	}

	// Set user context from middleware
	userID := c.GetString("user_id")
	if userID == "" {
		userID = c.GetHeader("X-User-ID")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}
	}
	installmentReq.UserID = userID
	installmentReq.InitiatedBy = userID

	// Create charge with installments
	response, err := h.cardService.ChargeCardWithInstallments(installmentReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}
