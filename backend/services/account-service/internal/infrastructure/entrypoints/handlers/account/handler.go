package accounthandler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/fintrack/account-service/internal/core/domain/entities"
	"github.com/fintrack/account-service/internal/core/service"
	"github.com/fintrack/account-service/internal/infrastructure/entrypoints/handlers/account/dto"
)

// Handler handles HTTP requests for account management operations
type Handler struct {
	accountService service.AccountServiceInterface
}

// New creates a new AccountHandler instance
func New(accountService service.AccountServiceInterface) *Handler {
	return &Handler{accountService: accountService}
}

// CreateAccount creates a new account
// @Summary Create a new account
// @Description Create a new financial account
// @Tags Accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateAccountRequest true "Account creation data"
// @Success 201 {object} dto.AccountResponse "Account created successfully"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/accounts [post]
func (h *Handler) CreateAccount(c *gin.Context) {
	var req dto.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account := &entities.Account{
		UserID:      req.UserID,
		AccountType: entities.AccountType(req.AccountType),
		Name:        req.Name,
		Description: req.Description,
		Currency:    entities.Currency(req.Currency),
		Balance:     req.InitialBalance,
		IsActive:    true,
	}

	createdAccount, err := h.accountService.CreateAccount(account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToAccountResponse(createdAccount)
	c.JSON(http.StatusCreated, response)
}

// GetAccount retrieves an account by ID
// @Summary Get account by ID
// @Description Retrieve account details by account ID
// @Tags Accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Account ID"
// @Success 200 {object} dto.AccountResponse "Account retrieved successfully"
// @Failure 400 {object} map[string]string "Invalid account ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Account not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/accounts/{id} [get]
func (h *Handler) GetAccount(c *gin.Context) {
	accountID := c.Param("id")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account ID is required"})
		return
	}

	account, err := h.accountService.GetAccountByID(accountID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		return
	}

	response := dto.ToAccountResponse(account)
	c.JSON(http.StatusOK, response)
}

// GetAccounts retrieves all accounts with pagination
// @Summary Get all accounts
// @Description Retrieve all accounts with pagination
// @Tags Accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Success 200 {object} dto.PaginatedAccountResponse "Accounts retrieved successfully"
// @Failure 400 {object} map[string]string "Invalid pagination parameters"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/accounts [get]
func (h *Handler) GetAccounts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	accounts, total, err := h.accountService.GetAllAccounts(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToPaginatedAccountResponse(accounts, total, page, pageSize)
	c.JSON(http.StatusOK, response)
}

// UpdateAccount updates an existing account
// @Summary Update account
// @Description Update account details
// @Tags Accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Account ID"
// @Param request body dto.UpdateAccountRequest true "Account update data"
// @Success 200 {object} dto.AccountResponse "Account updated successfully"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Account not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/accounts/{id} [put]
func (h *Handler) UpdateAccount(c *gin.Context) {
	accountID := c.Param("id")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account ID is required"})
		return
	}

	var req dto.UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedAccount, err := h.accountService.UpdateAccount(accountID, req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToAccountResponse(updatedAccount)
	c.JSON(http.StatusOK, response)
}

// DeleteAccount deletes an account
// @Summary Delete account
// @Description Delete an account by ID
// @Tags Accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Account ID"
// @Success 204 "Account deleted successfully"
// @Failure 400 {object} map[string]string "Invalid account ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Account not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/accounts/{id} [delete]
func (h *Handler) DeleteAccount(c *gin.Context) {
	accountID := c.Param("id")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account ID is required"})
		return
	}

	err := h.accountService.DeleteAccount(accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetAccountsByUser retrieves accounts for a specific user
// @Summary Get accounts by user ID
// @Description Retrieve all accounts belonging to a specific user
// @Tags Accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path string true "User ID"
// @Success 200 {array} dto.AccountResponse "Accounts retrieved successfully"
// @Failure 400 {object} map[string]string "Invalid user ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/accounts/user/{userId} [get]
func (h *Handler) GetAccountsByUser(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	accounts, err := h.accountService.GetAccountsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []dto.AccountResponse
	for _, account := range accounts {
		response = append(response, dto.ToAccountResponse(account))
	}

	c.JSON(http.StatusOK, response)
}

// GetBalance retrieves account balance
// @Summary Get account balance
// @Description Get the current balance of an account
// @Tags Accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Account ID"
// @Success 200 {object} dto.BalanceResponse "Balance retrieved successfully"
// @Failure 400 {object} map[string]string "Invalid account ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Account not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/accounts/{id}/balance [get]
func (h *Handler) GetBalance(c *gin.Context) {
	accountID := c.Param("id")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account ID is required"})
		return
	}

	balance, err := h.accountService.GetAccountBalance(accountID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		return
	}

	response := dto.BalanceResponse{
		AccountID: accountID,
		Balance:   balance,
	}
	c.JSON(http.StatusOK, response)
}

// UpdateBalance updates account balance
// @Summary Update account balance
// @Description Update the balance of an account
// @Tags Accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Account ID"
// @Param request body dto.UpdateBalanceRequest true "Balance update data"
// @Success 200 {object} dto.BalanceResponse "Balance updated successfully"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Account not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/accounts/{id}/balance [put]
func (h *Handler) UpdateBalance(c *gin.Context) {
	accountID := c.Param("id")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account ID is required"})
		return
	}

	var req dto.UpdateBalanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newBalance, err := h.accountService.UpdateAccountBalance(accountID, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.BalanceResponse{
		AccountID: accountID,
		Balance:   newBalance,
	}
	c.JSON(http.StatusOK, response)
}

// UpdateStatus updates account status
// @Summary Update account status
// @Description Activate or deactivate an account
// @Tags Accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Account ID"
// @Param request body dto.UpdateStatusRequest true "Status update data"
// @Success 200 {object} dto.AccountResponse "Status updated successfully"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Account not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/accounts/{id}/status [put]
func (h *Handler) UpdateStatus(c *gin.Context) {
	accountID := c.Param("id")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account ID is required"})
		return
	}

	var req dto.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedAccount, err := h.accountService.UpdateAccountStatus(accountID, req.IsActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToAccountResponse(updatedAccount)
	c.JSON(http.StatusOK, response)
}
