package dto

import (
	"time"

	"github.com/fintrack/account-service/internal/core/domain/entities"
)

// CreateInstallmentPlanRequest represents the request to create an installment plan
type CreateInstallmentPlanRequest struct {
	CardID            string    `json:"cardId,omitempty"` // Set from URL parameter, not from request body
	TotalAmount       float64   `json:"totalAmount" binding:"required,gt=0"`
	InstallmentsCount int       `json:"installmentsCount" binding:"required,min=1,max=24"`
	StartDate         time.Time `json:"startDate" binding:"required"`
	Description       string    `json:"description"`
	MerchantName      string    `json:"merchantName"`
	MerchantID        string    `json:"merchantId"`
	InterestRate      float64   `json:"interestRate,omitempty"`
	AdminFee          float64   `json:"adminFee,omitempty"`
	Reference         string    `json:"reference"`

	// User context (usually from authentication)
	UserID      string `json:"-"` // Set by middleware, not from request body
	InitiatedBy string `json:"-"` // Set by middleware, not from request body
}

// InstallmentPreviewRequest represents the request to preview installment calculations
type InstallmentPreviewRequest struct {
	Amount            float64   `json:"amount" binding:"required,gt=0"`
	InstallmentsCount int       `json:"installmentsCount" binding:"required,min=1,max=24"`
	StartDate         time.Time `json:"startDate" binding:"required"`
	InterestRate      float64   `json:"interestRate,omitempty"`
	AdminFee          float64   `json:"adminFee,omitempty"`
}

// InstallmentPreviewResponse represents the preview of an installment plan
type InstallmentPreviewResponse struct {
	TotalAmount       float64                  `json:"totalAmount"`
	InstallmentsCount int                      `json:"installmentsCount"`
	InstallmentAmount float64                  `json:"installmentAmount"`
	StartDate         time.Time                `json:"startDate"`
	InterestRate      float64                  `json:"interestRate"`
	TotalInterest     float64                  `json:"totalInterest"`
	AdminFee          float64                  `json:"adminFee"`
	TotalToPay        float64                  `json:"totalToPay"`
	Installments      []InstallmentPreviewItem `json:"installments"`
}

// InstallmentPreviewItem represents a single installment in the preview
type InstallmentPreviewItem struct {
	Number             int       `json:"number"`
	Amount             float64   `json:"amount"`
	DueDate            time.Time `json:"dueDate"`
	Principal          float64   `json:"principal"`
	Interest           float64   `json:"interest"`
	RemainingPrincipal float64   `json:"remainingPrincipal"`
}

// PayInstallmentRequest represents the request to pay an installment
type PayInstallmentRequest struct {
	InstallmentID    string  `json:"installment_id" binding:"required"`
	Amount           float64 `json:"amount" binding:"required,gt=0"`
	PaymentMethod    string  `json:"payment_method" binding:"required"`
	PaymentReference string  `json:"payment_reference"`

	// User context (usually from authentication)
	UserID      string `json:"-"` // Set by middleware
	InitiatedBy string `json:"-"` // Set by middleware
}

// InstallmentPlanResponse represents an installment plan in API responses
type InstallmentPlanResponse struct {
	ID                string                         `json:"id"`
	TransactionID     string                         `json:"transaction_id"`
	CardID            string                         `json:"card_id"`
	UserID            string                         `json:"user_id"`
	TotalAmount       float64                        `json:"total_amount"`
	InstallmentsCount int                            `json:"installments_count"`
	InstallmentAmount float64                        `json:"installment_amount"`
	StartDate         time.Time                      `json:"start_date"`
	Status            entities.InstallmentPlanStatus `json:"status"`
	PaidInstallments  int                            `json:"paid_installments"`
	RemainingAmount   float64                        `json:"remaining_amount"`
	Description       string                         `json:"description,omitempty"`
	MerchantName      string                         `json:"merchant_name,omitempty"`
	MerchantID        string                         `json:"merchant_id,omitempty"`
	InterestRate      float64                        `json:"interest_rate"`
	TotalInterest     float64                        `json:"total_interest"`
	AdminFee          float64                        `json:"admin_fee"`
	CreatedAt         time.Time                      `json:"created_at"`
	UpdatedAt         time.Time                      `json:"updated_at"`
	CompletedAt       *time.Time                     `json:"completed_at,omitempty"`
	CancelledAt       *time.Time                     `json:"cancelled_at,omitempty"`

	// Calculated fields
	CompletionPercentage  float64    `json:"completion_percentage"`
	RemainingInstallments int        `json:"remaining_installments"`
	NextDueDate           *time.Time `json:"next_due_date,omitempty"`
	NextInstallmentAmount float64    `json:"next_installment_amount,omitempty"`
	OverdueCount          int        `json:"overdue_count"`
	OverdueAmount         float64    `json:"overdue_amount"`

	// Related data (optional)
	Card         *CardResponse         `json:"card,omitempty"`
	Installments []InstallmentResponse `json:"installments,omitempty"`
}

// InstallmentResponse represents an installment in API responses
type InstallmentResponse struct {
	ID                   string                     `json:"id"`
	PlanID               string                     `json:"plan_id"`
	InstallmentNumber    int                        `json:"installment_number"`
	Amount               float64                    `json:"amount"`
	DueDate              time.Time                  `json:"due_date"`
	PaidDate             *time.Time                 `json:"paid_date,omitempty"`
	Status               entities.InstallmentStatus `json:"status"`
	PaidAmount           float64                    `json:"paid_amount"`
	RemainingAmount      float64                    `json:"remaining_amount"`
	PaymentMethod        string                     `json:"payment_method,omitempty"`
	PaymentReference     string                     `json:"payment_reference,omitempty"`
	PaymentTransactionID *string                    `json:"payment_transaction_id,omitempty"`
	LateFee              float64                    `json:"late_fee"`
	PenaltyAmount        float64                    `json:"penalty_amount"`
	GracePeriodDays      int                        `json:"grace_period_days"`
	CreatedAt            time.Time                  `json:"created_at"`
	UpdatedAt            time.Time                  `json:"updated_at"`

	// Calculated fields
	IsOverdue       bool      `json:"is_overdue"`
	DaysOverdue     int       `json:"days_overdue"`
	IsInGracePeriod bool      `json:"is_in_grace_period"`
	GracePeriodEnd  time.Time `json:"grace_period_end"`
}

// InstallmentSummaryResponse represents a summary of installments for a user
type InstallmentSummaryResponse struct {
	UserID                   string                       `json:"user_id"`
	TotalActivePlans         int                          `json:"total_active_plans"`
	TotalCompletedPlans      int                          `json:"total_completed_plans"`
	TotalCancelledPlans      int                          `json:"total_cancelled_plans"`
	TotalOutstandingAmount   float64                      `json:"total_outstanding_amount"`
	TotalPaidAmount          float64                      `json:"total_paid_amount"`
	TotalOverdueAmount       float64                      `json:"total_overdue_amount"`
	NextPaymentDue           *time.Time                   `json:"next_payment_due,omitempty"`
	NextPaymentAmount        float64                      `json:"next_payment_amount"`
	OverdueInstallmentsCount int                          `json:"overdue_installments_count"`
	UpcomingInstallments     []UpcomingInstallmentSummary `json:"upcoming_installments"`
	RecentActivity           []InstallmentActivitySummary `json:"recent_activity"`
}

// UpcomingInstallmentSummary represents upcoming installments
type UpcomingInstallmentSummary struct {
	InstallmentID     string    `json:"installment_id"`
	PlanID            string    `json:"plan_id"`
	CardID            string    `json:"card_id"`
	Amount            float64   `json:"amount"`
	DueDate           time.Time `json:"due_date"`
	Description       string    `json:"description"`
	MerchantName      string    `json:"merchant_name"`
	DaysUntilDue      int       `json:"days_until_due"`
	InstallmentNumber int       `json:"installment_number"`
	TotalInstallments int       `json:"total_installments"`
}

// InstallmentActivitySummary represents recent installment activity
type InstallmentActivitySummary struct {
	Date        time.Time `json:"date"`
	Action      string    `json:"action"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount,omitempty"`
	PlanID      string    `json:"plan_id,omitempty"`
}

// MonthlyInstallmentLoadResponse represents monthly installment load
type MonthlyInstallmentLoadResponse struct {
	UserID              string                 `json:"user_id"`
	Year                int                    `json:"year"`
	Month               int                    `json:"month"`
	TotalInstallments   int                    `json:"total_installments"`
	TotalAmount         float64                `json:"total_amount"`
	PaidInstallments    int                    `json:"paid_installments"`
	PaidAmount          float64                `json:"paid_amount"`
	PendingInstallments int                    `json:"pending_installments"`
	PendingAmount       float64                `json:"pending_amount"`
	OverdueInstallments int                    `json:"overdue_installments"`
	OverdueAmount       float64                `json:"overdue_amount"`
	DailyBreakdown      []DailyInstallmentLoad `json:"daily_breakdown"`
}

// DailyInstallmentLoad represents daily installment load within a month
type DailyInstallmentLoad struct {
	Day               int                          `json:"day"`
	Date              time.Time                    `json:"date"`
	InstallmentsCount int                          `json:"installments_count"`
	TotalAmount       float64                      `json:"total_amount"`
	Installments      []UpcomingInstallmentSummary `json:"installments"`
}

// InstallmentSummaryData represents summary data for repository queries
type InstallmentSummaryData struct {
	TotalActivePlans         int     `json:"total_active_plans"`
	TotalCompletedPlans      int     `json:"total_completed_plans"`
	TotalCancelledPlans      int     `json:"total_cancelled_plans"`
	TotalOutstandingAmount   float64 `json:"total_outstanding_amount"`
	TotalPaidAmount          float64 `json:"total_paid_amount"`
	TotalOverdueAmount       float64 `json:"total_overdue_amount"`
	OverdueInstallmentsCount int     `json:"overdue_installments_count"`
}

// CancelInstallmentPlanRequest represents request to cancel an installment plan
type CancelInstallmentPlanRequest struct {
	PlanID      string `json:"plan_id" binding:"required"`
	Reason      string `json:"reason" binding:"required"`
	CancelledBy string `json:"-"` // Set by middleware
}

// SuspendInstallmentPlanRequest represents request to suspend an installment plan
type SuspendInstallmentPlanRequest struct {
	PlanID      string `json:"plan_id" binding:"required"`
	Reason      string `json:"reason" binding:"required"`
	SuspendedBy string `json:"-"` // Set by middleware
}

// ReactivateInstallmentPlanRequest represents request to reactivate an installment plan
type ReactivateInstallmentPlanRequest struct {
	PlanID        string `json:"plan_id" binding:"required"`
	Reason        string `json:"reason" binding:"required"`
	ReactivatedBy string `json:"-"` // Set by middleware
}

// PaginatedInstallmentPlansResponse represents paginated response for installment plans
type PaginatedInstallmentPlansResponse struct {
	Data       []InstallmentPlanResponse `json:"data"`
	Pagination PaginationMeta            `json:"pagination"`
}

// PaginatedInstallmentsResponse represents paginated response for installments
type PaginatedInstallmentsResponse struct {
	Data       []InstallmentResponse `json:"data"`
	Pagination PaginationMeta        `json:"pagination"`
}

// Helper functions for mapping entities to DTOs

// MapInstallmentPlanToResponse maps an InstallmentPlan entity to InstallmentPlanResponse
func MapInstallmentPlanToResponse(plan *entities.InstallmentPlan, includeInstallments bool) *InstallmentPlanResponse {
	response := &InstallmentPlanResponse{
		ID:                    plan.ID,
		TransactionID:         plan.TransactionID,
		CardID:                plan.CardID,
		UserID:                plan.UserID,
		TotalAmount:           plan.TotalAmount,
		InstallmentsCount:     plan.InstallmentsCount,
		InstallmentAmount:     plan.InstallmentAmount,
		StartDate:             plan.StartDate,
		Status:                plan.Status,
		PaidInstallments:      plan.PaidInstallments,
		RemainingAmount:       plan.RemainingAmount,
		Description:           plan.Description,
		MerchantName:          plan.MerchantName,
		MerchantID:            plan.MerchantID,
		InterestRate:          plan.InterestRate,
		TotalInterest:         plan.TotalInterest,
		AdminFee:              plan.AdminFee,
		CreatedAt:             plan.CreatedAt,
		UpdatedAt:             plan.UpdatedAt,
		CompletedAt:           plan.CompletedAt,
		CancelledAt:           plan.CancelledAt,
		CompletionPercentage:  plan.GetCompletionPercentage(),
		RemainingInstallments: plan.GetRemainingInstallments(),
		OverdueAmount:         plan.GetOverdueAmount(),
	}

	// Add overdue count and next due information
	overdueInstallments := plan.GetOverdueInstallments()
	response.OverdueCount = len(overdueInstallments)

	nextDue := plan.GetNextDueInstallment()
	if nextDue != nil {
		response.NextDueDate = &nextDue.DueDate
		response.NextInstallmentAmount = nextDue.Amount
	}

	// Include installments if requested
	if includeInstallments && len(plan.Installments) > 0 {
		response.Installments = make([]InstallmentResponse, len(plan.Installments))
		for i, installment := range plan.Installments {
			response.Installments[i] = *MapInstallmentToResponse(&installment)
		}
	}

	return response
}

// Helper functions for pointer to string conversion
func stringPtrToString(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

// MapInstallmentToResponse maps an Installment entity to InstallmentResponse
func MapInstallmentToResponse(installment *entities.Installment) *InstallmentResponse {
	return &InstallmentResponse{
		ID:                   installment.ID,
		PlanID:               installment.PlanID,
		InstallmentNumber:    installment.InstallmentNumber,
		Amount:               installment.Amount,
		DueDate:              installment.DueDate,
		PaidDate:             installment.PaidDate,
		Status:               installment.Status,
		PaidAmount:           installment.PaidAmount,
		RemainingAmount:      installment.RemainingAmount,
		PaymentMethod:        stringPtrToString(installment.PaymentMethod),
		PaymentReference:     stringPtrToString(installment.PaymentReference),
		PaymentTransactionID: installment.PaymentTransactionID,
		LateFee:              installment.LateFee,
		PenaltyAmount:        installment.PenaltyAmount,
		GracePeriodDays:      installment.GracePeriodDays,
		CreatedAt:            installment.CreatedAt,
		UpdatedAt:            installment.UpdatedAt,
		IsOverdue:            installment.IsOverdue(),
		DaysOverdue:          installment.GetDaysOverdue(),
		IsInGracePeriod:      installment.IsInGracePeriod(),
		GracePeriodEnd:       installment.GetGracePeriodEnd(),
	}
}

// PaginatedInstallmentPlanResponse represents paginated installment plan list response
type PaginatedInstallmentPlanResponse struct {
	Data       []InstallmentPlanResponse `json:"data"`
	Pagination PaginationMeta            `json:"pagination"`
}

// PaginatedInstallmentResponse represents paginated installment list response
type PaginatedInstallmentResponse struct {
	Data       []InstallmentResponse `json:"data"`
	Pagination PaginationMeta        `json:"pagination"`
}

// ToInstallmentResponse converts an Installment entity to InstallmentResponse DTO
func ToInstallmentResponse(installment *entities.Installment) InstallmentResponse {
	return InstallmentResponse{
		ID:                   installment.ID,
		PlanID:               installment.PlanID,
		InstallmentNumber:    installment.InstallmentNumber,
		Amount:               installment.Amount,
		DueDate:              installment.DueDate,
		PaidDate:             installment.PaidDate,
		Status:               installment.Status,
		PaidAmount:           installment.PaidAmount,
		RemainingAmount:      installment.RemainingAmount,
		PaymentMethod:        stringPtrToString(installment.PaymentMethod),
		PaymentReference:     stringPtrToString(installment.PaymentReference),
		PaymentTransactionID: installment.PaymentTransactionID,
		LateFee:              installment.LateFee,
		PenaltyAmount:        installment.PenaltyAmount,
		GracePeriodDays:      installment.GracePeriodDays,
		CreatedAt:            installment.CreatedAt,
		UpdatedAt:            installment.UpdatedAt,
		IsOverdue:            installment.IsOverdue(),
		DaysOverdue:          installment.GetDaysOverdue(),
		IsInGracePeriod:      installment.IsInGracePeriod(),
		GracePeriodEnd:       installment.GetGracePeriodEnd(),
	}
}

// ToInstallmentPlanResponse converts an InstallmentPlan entity to InstallmentPlanResponse DTO
func ToInstallmentPlanResponse(plan *entities.InstallmentPlan) InstallmentPlanResponse {
	response := InstallmentPlanResponse{
		ID:                plan.ID,
		TransactionID:     plan.TransactionID,
		CardID:            plan.CardID,
		UserID:            plan.UserID,
		TotalAmount:       plan.TotalAmount,
		InstallmentsCount: plan.InstallmentsCount,
		InstallmentAmount: plan.InstallmentAmount,
		StartDate:         plan.StartDate,
		Status:            plan.Status,
		PaidInstallments:  plan.PaidInstallments,
		RemainingAmount:   plan.RemainingAmount,
		Description:       plan.Description,
		MerchantName:      plan.MerchantName,
		MerchantID:        plan.MerchantID,
		InterestRate:      plan.InterestRate,
		TotalInterest:     plan.TotalInterest,
		AdminFee:          plan.AdminFee,
		CreatedAt:         plan.CreatedAt,
		UpdatedAt:         plan.UpdatedAt,
		CompletedAt:       plan.CompletedAt,
		CancelledAt:       plan.CancelledAt,
	}

	// Calculate derived fields
	if plan.InstallmentsCount > 0 {
		response.CompletionPercentage = float64(plan.PaidInstallments) / float64(plan.InstallmentsCount) * 100
		response.RemainingInstallments = plan.InstallmentsCount - plan.PaidInstallments
	}

	// Get next due installment info if available
	if len(plan.Installments) > 0 {
		for _, installment := range plan.Installments {
			if installment.Status == entities.InstallmentStatusPending || installment.Status == entities.InstallmentStatusOverdue {
				response.NextDueDate = &installment.DueDate
				response.NextInstallmentAmount = installment.Amount
				break
			}
		}

		// Count overdue installments
		var overdueCount int
		var overdueAmount float64
		for _, installment := range plan.Installments {
			if installment.Status == entities.InstallmentStatusOverdue {
				overdueCount++
				overdueAmount += installment.RemainingAmount
			}
		}
		response.OverdueCount = overdueCount
		response.OverdueAmount = overdueAmount
	}

	// Convert installments if included
	if len(plan.Installments) > 0 {
		response.Installments = make([]InstallmentResponse, len(plan.Installments))
		for i, installment := range plan.Installments {
			response.Installments[i] = ToInstallmentResponse(&installment)
		}
	}

	return response
}

// ToPaginatedInstallmentPlanResponse converts installment plans with pagination info to response
func ToPaginatedInstallmentPlanResponse(plans []*entities.InstallmentPlan, total int64, page, pageSize int) PaginatedInstallmentPlanResponse {
	data := make([]InstallmentPlanResponse, len(plans))
	for i, plan := range plans {
		data[i] = ToInstallmentPlanResponse(plan)
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return PaginatedInstallmentPlanResponse{
		Data: data,
		Pagination: PaginationMeta{
			CurrentPage: page,
			PageSize:    pageSize,
			TotalItems:  total,
			TotalPages:  totalPages,
		},
	}
}

// ToPaginatedInstallmentResponse converts installments with pagination info to response
func ToPaginatedInstallmentResponse(installments []*entities.Installment, total int64, page, pageSize int) PaginatedInstallmentResponse {
	data := make([]InstallmentResponse, len(installments))
	for i, installment := range installments {
		data[i] = ToInstallmentResponse(installment)
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return PaginatedInstallmentResponse{
		Data: data,
		Pagination: PaginationMeta{
			CurrentPage: page,
			PageSize:    pageSize,
			TotalItems:  total,
			TotalPages:  totalPages,
		},
	}
}
