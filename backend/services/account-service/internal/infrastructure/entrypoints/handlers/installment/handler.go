package installment

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/fintrack/account-service/internal/core/ports"
	"github.com/fintrack/account-service/internal/infrastructure/entrypoints/handlers/card/dto"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	installmentService ports.InstallmentServiceInterface
	cardService        ports.CardServiceInterface
}

func New(installmentService ports.InstallmentServiceInterface, cardService ports.CardServiceInterface) *Handler {
	return &Handler{
		installmentService: installmentService,
		cardService:        cardService,
	}
}

// PreviewInstallmentPlan calculates and previews an installment plan
// @Summary Preview installment plan
// @Description Calculate and preview installment plan for a given amount and terms
// @Tags Installments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param cardId path string true "Card ID"
// @Param preview body dto.InstallmentPreviewRequest true "Preview request data"
// @Success 200 {object} dto.InstallmentPreviewResponse "Installment plan preview"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 404 {object} map[string]string "Card not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/cards/{cardId}/installments/preview [post]
func (h *Handler) PreviewInstallmentPlan(c *gin.Context) {
	cardID := c.Param("cardId")
	if cardID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "card ID is required"})
		return
	}

	var req dto.InstallmentPreviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate card exists and is accessible
	card, err := h.cardService.GetCardByID(cardID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "card not found"})
		return
	}

	// Verify card supports installments (credit cards only)
	if card.CardType != "credit" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only credit cards support installments"})
		return
	}

	// Calculate preview
	preview, err := h.installmentService.CalculateInstallmentPlan(
		req.Amount,
		req.InstallmentsCount,
		req.StartDate,
		req.InterestRate,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, preview)
}

// ChargeCardWithInstallments creates a purchase with installments
// @Summary Charge card with installments
// @Description Create a purchase on a credit card with installment plan
// @Tags Installments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param cardId path string true "Card ID"
// @Param charge body dto.CreateInstallmentPlanRequest true "Charge with installments data"
// @Success 201 {object} dto.ChargeWithInstallmentsResponse "Purchase created with installments"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 404 {object} map[string]string "Card not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/cards/{cardId}/charge-installments [post]
func (h *Handler) ChargeCardWithInstallments(c *gin.Context) {
	fmt.Printf("⭐⭐⭐ HANDLER - ChargeCardWithInstallments STARTED ⭐⭐⭐\n")

	cardID := c.Param("cardId")
	if cardID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "card ID is required"})
		return
	}

	fmt.Printf("⭐⭐⭐ HANDLER - Processing request for CardID: %s ⭐⭐⭐\n", cardID)

	var req dto.CreateInstallmentPlanRequest

	// DEBUG: Log incoming request body
	bodyBytes, _ := c.GetRawData()
	fmt.Printf("DEBUG - Incoming request body: %s\n", string(bodyBytes))

	// Reset body for binding
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("DEBUG - Binding error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("DEBUG - Parsed request: %+v\n", req)

	// Set card ID from URL parameter
	req.CardID = cardID

	// Validate that CardID is set
	if req.CardID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "card ID is required"})
		return
	}

	// Set user context from middleware (would be set by auth middleware)
	// For now, we'll use a placeholder - in real implementation this comes from JWT
	userID := c.GetString("user_id")
	if userID == "" {
		// Fallback for testing - get from header
		userID = c.GetHeader("X-User-ID")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}
	}
	req.UserID = userID
	req.InitiatedBy = userID

	// Create charge with installments
	fmt.Printf("🎯🎯🎯 HANDLER - About to call cardService.ChargeCardWithInstallments for CardID: %s 🎯🎯🎯\n", req.CardID)
	response, err := h.cardService.ChargeCardWithInstallments(&req)
	if err != nil {
		fmt.Printf("🎯🎯🎯 HANDLER - ERROR from cardService: %v 🎯🎯🎯\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("🎯🎯🎯 HANDLER - SUCCESS from cardService 🎯🎯🎯\n")

	c.JSON(http.StatusCreated, response)
}

// GetInstallmentPlansByCard retrieves active installment plans for a card
// @Summary Get installment plans by card
// @Description Get all installment plans for a specific card with pagination
// @Tags Installments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param cardId path string true "Card ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} dto.PaginatedInstallmentPlanResponse "Installment plans list"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 404 {object} map[string]string "Card not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/cards/{cardId}/installment-plans [get]
func (h *Handler) GetInstallmentPlansByCard(c *gin.Context) {
	cardID := c.Param("cardId")
	if cardID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "card ID is required"})
		return
	}

	// Parse pagination parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Verify card exists
	_, err = h.cardService.GetCardByID(cardID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "card not found"})
		return
	}

	// Get installment plans
	plans, total, err := h.installmentService.GetInstallmentPlansByCard(cardID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve installment plans"})
		return
	}

	// Convert to response format
	response := dto.ToPaginatedInstallmentPlanResponse(plans, total, page, pageSize)
	c.JSON(http.StatusOK, response)
}

// GetInstallmentPlan retrieves a specific installment plan with details
// @Summary Get installment plan details
// @Description Get detailed information about a specific installment plan including all installments
// @Tags Installments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param planId path string true "Installment Plan ID"
// @Success 200 {object} dto.InstallmentPlanResponse "Installment plan details"
// @Failure 404 {object} map[string]string "Plan not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/installment-plans/{planId} [get]
func (h *Handler) GetInstallmentPlan(c *gin.Context) {
	planID := c.Param("planId")
	if planID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "plan ID is required"})
		return
	}

	// Get installment plan with installments
	plan, err := h.installmentService.GetInstallmentPlan(planID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "installment plan not found"})
		return
	}

	// Convert to response format
	response := dto.ToInstallmentPlanResponse(plan)
	c.JSON(http.StatusOK, response)
}

// PayInstallment processes payment for a specific installment
// @Summary Pay installment
// @Description Process payment for a specific installment
// @Tags Installments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param installmentId path string true "Installment ID"
// @Param payment body dto.PayInstallmentRequest true "Payment data"
// @Success 200 {object} dto.InstallmentResponse "Updated installment"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 404 {object} map[string]string "Installment not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/installments/{installmentId}/pay [post]
func (h *Handler) PayInstallment(c *gin.Context) {
	installmentID := c.Param("installmentId")
	if installmentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "installment ID is required"})
		return
	}

	var req dto.PayInstallmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set installment ID from URL parameter
	req.InstallmentID = installmentID

	// Set user context from middleware
	userID := c.GetString("user_id")
	if userID == "" {
		userID = c.GetHeader("X-User-ID")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}
	}
	req.UserID = userID
	req.InitiatedBy = userID

	// Process payment
	installment, err := h.installmentService.PayInstallment(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert to response format
	response := dto.ToInstallmentResponse(installment)
	c.JSON(http.StatusOK, response)
}

// GetUserInstallmentPlans retrieves all installment plans for the authenticated user
// @Summary Get user installment plans
// @Description Get all installment plans for the authenticated user with filtering and pagination
// @Tags Installments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter by status (active, completed, cancelled, suspended)"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} dto.PaginatedInstallmentPlanResponse "User installment plans"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/installment-plans [get]
func (h *Handler) GetUserInstallmentPlans(c *gin.Context) {
	// Get user ID from context
	userID := c.GetString("user_id")
	if userID == "" {
		userID = c.GetHeader("X-User-ID")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}
	}

	// Parse query parameters
	status := c.Query("status")
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Get user installment plans
	plans, total, err := h.installmentService.GetInstallmentPlansByUser(userID, status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve installment plans"})
		return
	}

	// Convert to response format
	response := dto.ToPaginatedInstallmentPlanResponse(plans, total, page, pageSize)
	c.JSON(http.StatusOK, response)
}

// GetOverdueInstallments retrieves overdue installments for the authenticated user
// @Summary Get overdue installments
// @Description Get all overdue installments for the authenticated user
// @Tags Installments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} dto.PaginatedInstallmentResponse "Overdue installments"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/installments/overdue [get]
func (h *Handler) GetOverdueInstallments(c *gin.Context) {
	// Get user ID from context
	userID := c.GetString("user_id")
	if userID == "" {
		userID = c.GetHeader("X-User-ID")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}
	}

	// Parse pagination parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Get overdue installments
	installments, total, err := h.installmentService.GetOverdueInstallments(userID, pageSize, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve overdue installments"})
		return
	}

	// Convert to response format
	response := dto.ToPaginatedInstallmentResponse(installments, total, page, pageSize)
	c.JSON(http.StatusOK, response)
}

// GetUpcomingInstallments retrieves upcoming installments for the authenticated user
// @Summary Get upcoming installments
// @Description Get upcoming installments for the authenticated user within specified days
// @Tags Installments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param days query int false "Number of days to look ahead" default(30)
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} dto.PaginatedInstallmentResponse "Upcoming installments"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/installments/upcoming [get]
func (h *Handler) GetUpcomingInstallments(c *gin.Context) {
	// Get user ID from context
	userID := c.GetString("user_id")
	if userID == "" {
		userID = c.GetHeader("X-User-ID")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}
	}

	// Parse query parameters
	days, err := strconv.Atoi(c.DefaultQuery("days", "30"))
	if err != nil || days < 1 {
		days = 30
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Get upcoming installments
	installments, total, err := h.installmentService.GetUpcomingInstallments(userID, days, pageSize, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve upcoming installments"})
		return
	}

	// Convert to response format
	response := dto.ToPaginatedInstallmentResponse(installments, total, page, pageSize)
	c.JSON(http.StatusOK, response)
}

// GetInstallmentSummary retrieves installment summary for the authenticated user
// @Summary Get installment summary
// @Description Get comprehensive installment summary for the authenticated user
// @Tags Installments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.InstallmentSummaryResponse "Installment summary"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/installments/summary [get]
func (h *Handler) GetInstallmentSummary(c *gin.Context) {
	// Get user ID from context
	userID := c.GetString("user_id")
	if userID == "" {
		userID = c.GetHeader("X-User-ID")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}
	}

	// Get installment summary
	summary, err := h.installmentService.GetInstallmentSummary(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve installment summary"})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetMonthlyInstallmentLoad retrieves monthly installment load for the authenticated user
// @Summary Get monthly installment load
// @Description Get detailed monthly installment load breakdown for the authenticated user
// @Tags Installments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param year query int false "Year" default(current year)
// @Param month query int false "Month (1-12)" default(current month)
// @Success 200 {object} dto.MonthlyInstallmentLoadResponse "Monthly installment load"
// @Failure 400 {object} map[string]string "Invalid parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/installments/monthly-load [get]
func (h *Handler) GetMonthlyInstallmentLoad(c *gin.Context) {
	// Get user ID from context
	userID := c.GetString("user_id")
	if userID == "" {
		userID = c.GetHeader("X-User-ID")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}
	}

	// Parse query parameters
	now := time.Now()
	year, err := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(now.Year())))
	if err != nil || year < 2020 || year > 2030 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year parameter"})
		return
	}

	month, err := strconv.Atoi(c.DefaultQuery("month", strconv.Itoa(int(now.Month()))))
	if err != nil || month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid month parameter"})
		return
	}

	// Get monthly installment load
	load, err := h.installmentService.GetMonthlyInstallmentLoad(userID, year, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve monthly installment load"})
		return
	}

	c.JSON(http.StatusOK, load)
}

// CancelInstallmentPlan cancels an installment plan
// @Summary Cancel installment plan
// @Description Cancel an active installment plan
// @Tags Installments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param planId path string true "Installment Plan ID"
// @Param cancel body dto.CancelInstallmentPlanRequest true "Cancellation data"
// @Success 200 {object} dto.InstallmentPlanResponse "Cancelled installment plan"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 404 {object} map[string]string "Plan not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/installment-plans/{planId}/cancel [post]
func (h *Handler) CancelInstallmentPlan(c *gin.Context) {
	planID := c.Param("planId")
	if planID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "plan ID is required"})
		return
	}

	var req dto.CancelInstallmentPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	userID := c.GetString("user_id")
	if userID == "" {
		userID = c.GetHeader("X-User-ID")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}
	}

	// Cancel installment plan
	plan, err := h.installmentService.CancelInstallmentPlan(planID, req.Reason, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert to response format
	response := dto.ToInstallmentPlanResponse(plan)
	c.JSON(http.StatusOK, response)
}

// GetInstallmentsByPlan retrieves all installments for a specific plan
// @Summary Get installments by plan
// @Description Get all installments for a specific installment plan
// @Tags Installments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param planId path string true "Installment Plan ID"
// @Success 200 {array} dto.InstallmentResponse "List of installments"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 404 {object} map[string]string "Plan not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/installment-plans/{planId}/installments [get]
func (h *Handler) GetInstallmentsByPlan(c *gin.Context) {
	planID := c.Param("planId")
	if planID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "plan ID is required"})
		return
	}

	// Get user ID from context for authorization
	userID := c.GetString("user_id")
	if userID == "" {
		userID = c.GetHeader("X-User-ID")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}
	}

	// First verify that the plan belongs to the user
	plan, err := h.installmentService.GetInstallmentPlan(planID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "installment plan not found"})
		return
	}

	if plan.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied to this installment plan"})
		return
	}

	// Get installments for the plan
	installments, err := h.installmentService.GetInstallmentsByPlan(planID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve installments"})
		return
	}

	// Convert to response format
	var response []dto.InstallmentResponse
	for _, installment := range installments {
		response = append(response, dto.ToInstallmentResponse(installment))
	}

	c.JSON(http.StatusOK, response)
}
