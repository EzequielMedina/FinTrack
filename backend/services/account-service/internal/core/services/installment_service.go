package services

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/fintrack/account-service/internal/core/domain/entities"
	"github.com/fintrack/account-service/internal/core/ports"
	"github.com/fintrack/account-service/internal/infrastructure/entrypoints/handlers/card/dto"
	"github.com/google/uuid"
)

// InstallmentService implements the installment service interface
type InstallmentService struct {
	installmentPlanRepo  ports.InstallmentPlanRepositoryInterface
	installmentRepo      ports.InstallmentRepositoryInterface
	installmentAuditRepo ports.InstallmentPlanAuditRepositoryInterface
	cardRepo             ports.CardRepositoryInterface
}

// NewInstallmentService creates a new installment service
func NewInstallmentService(
	installmentPlanRepo ports.InstallmentPlanRepositoryInterface,
	installmentRepo ports.InstallmentRepositoryInterface,
	installmentAuditRepo ports.InstallmentPlanAuditRepositoryInterface,
	cardRepo ports.CardRepositoryInterface,
) ports.InstallmentServiceInterface {
	return &InstallmentService{
		installmentPlanRepo:  installmentPlanRepo,
		installmentRepo:      installmentRepo,
		installmentAuditRepo: installmentAuditRepo,
		cardRepo:             cardRepo,
	}
}

// CalculateInstallmentPlan calculates an installment plan preview
func (s *InstallmentService) CalculateInstallmentPlan(amount float64, installmentsCount int, startDate time.Time, interestRate float64) (*dto.InstallmentPreviewResponse, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}
	if installmentsCount <= 0 || installmentsCount > 24 {
		return nil, errors.New("installments count must be between 1 and 24")
	}
	if startDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return nil, errors.New("start date cannot be in the past")
	}

	// Calculate interest and fees
	totalInterest := amount * (interestRate / 100)
	adminFee := 0.0 // Could be calculated based on business rules
	totalToPay := amount + totalInterest + adminFee

	// Calculate installment amount
	installmentAmount := totalToPay / float64(installmentsCount)

	// Round to 2 decimal places
	installmentAmount = math.Round(installmentAmount*100) / 100

	// Generate installment preview items
	installments := make([]dto.InstallmentPreviewItem, installmentsCount)
	remainingPrincipal := amount

	for i := 0; i < installmentsCount; i++ {
		dueDate := startDate.AddDate(0, i, 0) // Monthly installments

		// Calculate interest for this installment
		interestForInstallment := totalInterest / float64(installmentsCount)
		principalForInstallment := installmentAmount - interestForInstallment

		// Adjust last installment for rounding differences
		if i == installmentsCount-1 {
			// Ensure the last installment covers any rounding differences
			principalForInstallment = remainingPrincipal
			installmentAmount = principalForInstallment + interestForInstallment
		}

		installments[i] = dto.InstallmentPreviewItem{
			Number:             i + 1,
			Amount:             math.Round(installmentAmount*100) / 100,
			DueDate:            dueDate,
			Principal:          math.Round(principalForInstallment*100) / 100,
			Interest:           math.Round(interestForInstallment*100) / 100,
			RemainingPrincipal: math.Round(remainingPrincipal*100) / 100,
		}

		remainingPrincipal -= principalForInstallment
	}

	return &dto.InstallmentPreviewResponse{
		TotalAmount:       amount,
		InstallmentsCount: installmentsCount,
		InstallmentAmount: installmentAmount,
		StartDate:         startDate,
		InterestRate:      interestRate,
		TotalInterest:     math.Round(totalInterest*100) / 100,
		AdminFee:          adminFee,
		TotalToPay:        math.Round(totalToPay*100) / 100,
		Installments:      installments,
	}, nil
}

// CreateInstallmentPlan creates a new installment plan
func (s *InstallmentService) CreateInstallmentPlan(req *dto.CreateInstallmentPlanRequest) (*entities.InstallmentPlan, error) {
	// Validate card exists and belongs to user
	card, err := s.cardRepo.GetByIDWithAccount(req.CardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get card: %w", err)
	}

	if card.Account.UserID != req.UserID {
		return nil, errors.New("card does not belong to user")
	}

	// Validate card can create installment plans
	if !card.CanCreateInstallmentPlan(req.TotalAmount) {
		return nil, errors.New("card cannot create installment plan with this amount")
	}

	// Calculate installment plan details
	preview, err := s.CalculateInstallmentPlan(req.TotalAmount, req.InstallmentsCount, req.StartDate, req.InterestRate)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate installment plan: %w", err)
	}

	// Create installment plan entity
	plan := &entities.InstallmentPlan{
		ID:                uuid.New().String(),
		TransactionID:     "", // Will be set when transaction is created
		CardID:            req.CardID,
		UserID:            req.UserID,
		TotalAmount:       req.TotalAmount,
		InstallmentsCount: req.InstallmentsCount,
		InstallmentAmount: preview.InstallmentAmount,
		StartDate:         req.StartDate,
		Status:            entities.InstallmentPlanStatusActive,
		PaidInstallments:  0,
		RemainingAmount:   preview.TotalToPay,
		Description:       req.Description,
		MerchantName:      req.MerchantName,
		MerchantID:        req.MerchantID,
		InterestRate:      req.InterestRate,
		TotalInterest:     preview.TotalInterest,
		AdminFee:          req.AdminFee,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Save installment plan
	savedPlan, err := s.installmentPlanRepo.Create(plan)
	if err != nil {
		return nil, fmt.Errorf("failed to create installment plan: %w", err)
	}

	// Create individual installments
	for _, installmentPreview := range preview.Installments {
		installment := &entities.Installment{
			ID:                uuid.New().String(),
			PlanID:            savedPlan.ID,
			InstallmentNumber: installmentPreview.Number,
			Amount:            installmentPreview.Amount,
			DueDate:           installmentPreview.DueDate,
			Status:            entities.InstallmentStatusPending,
			PaidAmount:        0.0,
			RemainingAmount:   installmentPreview.Amount,
			LateFee:           0.0,
			PenaltyAmount:     0.0,
			GracePeriodDays:   7, // Default grace period
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}

		_, err := s.installmentRepo.Create(installment)
		if err != nil {
			return nil, fmt.Errorf("failed to create installment %d: %w", installmentPreview.Number, err)
		}
	}

	// Create audit record
	audit := &entities.InstallmentPlanAudit{
		ID:           uuid.New().String(),
		PlanID:       savedPlan.ID,
		Action:       "created",
		NewStatus:    (*string)(&savedPlan.Status),
		ChangedBy:    req.InitiatedBy,
		ChangeReason: fmt.Sprintf("Plan created with %d installments for %.2f", req.InstallmentsCount, req.TotalAmount),
		Metadata: map[string]interface{}{
			"card_id":       req.CardID,
			"merchant_name": req.MerchantName,
			"user_id":       req.UserID,
		},
		CreatedAt: time.Now(),
	}

	err = s.installmentAuditRepo.Create(audit)
	if err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: failed to create audit record: %v\n", err)
	}

	return savedPlan, nil
}

// GetInstallmentPlan retrieves an installment plan by ID
func (s *InstallmentService) GetInstallmentPlan(planID string) (*entities.InstallmentPlan, error) {
	return s.installmentPlanRepo.GetByIDWithInstallments(planID)
}

// GetInstallmentPlansByCard retrieves installment plans for a specific card
func (s *InstallmentService) GetInstallmentPlansByCard(cardID string, page, pageSize int) ([]*entities.InstallmentPlan, int64, error) {
	offset := (page - 1) * pageSize
	return s.installmentPlanRepo.GetByCard(cardID, "", pageSize, offset)
}

// GetInstallmentPlansByUser retrieves installment plans for a specific user
func (s *InstallmentService) GetInstallmentPlansByUser(userID string, status string, page, pageSize int) ([]*entities.InstallmentPlan, int64, error) {
	offset := (page - 1) * pageSize
	return s.installmentPlanRepo.GetByUser(userID, status, pageSize, offset)
}

// CancelInstallmentPlan cancels an installment plan
func (s *InstallmentService) CancelInstallmentPlan(planID, reason string, cancelledBy string) (*entities.InstallmentPlan, error) {
	plan, err := s.installmentPlanRepo.GetByID(planID)
	if err != nil {
		return nil, fmt.Errorf("failed to get installment plan: %w", err)
	}

	if !plan.CanCancel() {
		return nil, errors.New("installment plan cannot be cancelled in current status")
	}

	// Update plan status
	oldStatus := plan.Status
	cancelledStatus := entities.InstallmentPlanStatusCancelled
	plan.Status = entities.InstallmentPlanStatusCancelled
	plan.CancelledAt = &time.Time{}
	*plan.CancelledAt = time.Now()
	plan.UpdatedAt = time.Now()

	updatedPlan, err := s.installmentPlanRepo.Update(plan)
	if err != nil {
		return nil, fmt.Errorf("failed to update installment plan: %w", err)
	}

	// Cancel pending installments
	installments, err := s.installmentRepo.GetPendingByPlan(planID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending installments: %w", err)
	}

	for _, installment := range installments {
		installment.Status = entities.InstallmentStatusCancelled
		installment.UpdatedAt = time.Now()
		_, err := s.installmentRepo.Update(installment)
		if err != nil {
			fmt.Printf("Warning: failed to cancel installment %s: %v\n", installment.ID, err)
		}
	}

	// Create audit record
	audit := &entities.InstallmentPlanAudit{
		ID:           uuid.New().String(),
		PlanID:       planID,
		Action:       "cancelled",
		OldStatus:    (*string)(&oldStatus),
		NewStatus:    (*string)(&cancelledStatus),
		ChangedBy:    cancelledBy,
		ChangeReason: reason,
		Metadata: map[string]interface{}{
			"user_id": plan.UserID,
		},
		CreatedAt: time.Now(),
	}

	err = s.installmentAuditRepo.Create(audit)
	if err != nil {
		fmt.Printf("Warning: failed to create audit record: %v\n", err)
	}

	return updatedPlan, nil
}

// SuspendInstallmentPlan suspends an installment plan
func (s *InstallmentService) SuspendInstallmentPlan(planID, reason string, suspendedBy string) (*entities.InstallmentPlan, error) {
	plan, err := s.installmentPlanRepo.GetByID(planID)
	if err != nil {
		return nil, fmt.Errorf("failed to get installment plan: %w", err)
	}

	if plan.Status != entities.InstallmentPlanStatusActive {
		return nil, errors.New("only active installment plans can be suspended")
	}

	// Update plan status
	oldStatus := plan.Status
	suspendedStatus := entities.InstallmentPlanStatusSuspended
	plan.Status = entities.InstallmentPlanStatusSuspended
	plan.UpdatedAt = time.Now()

	updatedPlan, err := s.installmentPlanRepo.Update(plan)
	if err != nil {
		return nil, fmt.Errorf("failed to update installment plan: %w", err)
	}

	// Create audit record
	audit := &entities.InstallmentPlanAudit{
		ID:           uuid.New().String(),
		PlanID:       planID,
		Action:       "suspended",
		OldStatus:    (*string)(&oldStatus),
		NewStatus:    (*string)(&suspendedStatus),
		ChangedBy:    suspendedBy,
		ChangeReason: reason,
		Metadata: map[string]interface{}{
			"user_id": plan.UserID,
		},
		CreatedAt: time.Now(),
	}

	err = s.installmentAuditRepo.Create(audit)
	if err != nil {
		fmt.Printf("Warning: failed to create audit record: %v\n", err)
	}

	return updatedPlan, nil
}

// ReactivateInstallmentPlan reactivates a suspended installment plan
func (s *InstallmentService) ReactivateInstallmentPlan(planID, reason string, reactivatedBy string) (*entities.InstallmentPlan, error) {
	plan, err := s.installmentPlanRepo.GetByID(planID)
	if err != nil {
		return nil, fmt.Errorf("failed to get installment plan: %w", err)
	}

	if plan.Status != entities.InstallmentPlanStatusSuspended {
		return nil, errors.New("only suspended installment plans can be reactivated")
	}

	// Update plan status
	oldStatus := plan.Status
	activeStatus := entities.InstallmentPlanStatusActive
	plan.Status = entities.InstallmentPlanStatusActive
	plan.UpdatedAt = time.Now()

	updatedPlan, err := s.installmentPlanRepo.Update(plan)
	if err != nil {
		return nil, fmt.Errorf("failed to update installment plan: %w", err)
	}

	// Create audit record
	audit := &entities.InstallmentPlanAudit{
		ID:           uuid.New().String(),
		PlanID:       planID,
		Action:       "reactivated",
		OldStatus:    (*string)(&oldStatus),
		NewStatus:    (*string)(&activeStatus),
		ChangedBy:    reactivatedBy,
		ChangeReason: reason,
		Metadata: map[string]interface{}{
			"user_id": plan.UserID,
		},
		CreatedAt: time.Now(),
	}

	err = s.installmentAuditRepo.Create(audit)
	if err != nil {
		fmt.Printf("Warning: failed to create audit record: %v\n", err)
	}

	return updatedPlan, nil
}

// GetInstallment retrieves an installment by ID
func (s *InstallmentService) GetInstallment(installmentID string) (*entities.Installment, error) {
	return s.installmentRepo.GetByID(installmentID)
}

// GetInstallmentsByPlan retrieves all installments for a plan
func (s *InstallmentService) GetInstallmentsByPlan(planID string) ([]*entities.Installment, error) {
	return s.installmentRepo.GetByPlan(planID)
}

// GetOverdueInstallments retrieves overdue installments for a user
func (s *InstallmentService) GetOverdueInstallments(userID string, limit, offset int) ([]*entities.Installment, int64, error) {
	return s.installmentRepo.GetOverdue(userID, limit, offset)
}

// GetUpcomingInstallments retrieves upcoming installments for a user
func (s *InstallmentService) GetUpcomingInstallments(userID string, days int, limit, offset int) ([]*entities.Installment, int64, error) {
	return s.installmentRepo.GetUpcoming(userID, days, limit, offset)
}

// PayInstallment processes payment for an installment
func (s *InstallmentService) PayInstallment(req *dto.PayInstallmentRequest) (*entities.Installment, error) {
	installment, err := s.installmentRepo.GetByID(req.InstallmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get installment: %w", err)
	}

	if !installment.CanPay() {
		return nil, errors.New("installment cannot be paid in current status")
	}

	if req.Amount > installment.RemainingAmount {
		return nil, errors.New("payment amount exceeds remaining amount")
	}

	// Get the installment plan to verify user ownership
	plan, err := s.installmentPlanRepo.GetByID(installment.PlanID)
	if err != nil {
		return nil, fmt.Errorf("failed to get installment plan: %w", err)
	}

	if plan.UserID != req.UserID {
		return nil, errors.New("installment does not belong to user")
	}

	// Process payment logic
	installment.PaidAmount += req.Amount
	installment.RemainingAmount -= req.Amount
	installment.PaymentMethod = req.PaymentMethod
	installment.PaymentReference = req.PaymentReference
	installment.UpdatedAt = time.Now()

	if installment.RemainingAmount <= 0.01 { // Allow for small rounding differences
		installment.Status = entities.InstallmentStatusPaid
		installment.PaidDate = &time.Time{}
		*installment.PaidDate = time.Now()
		installment.RemainingAmount = 0
	}

	updatedInstallment, err := s.installmentRepo.Update(installment)
	if err != nil {
		return nil, fmt.Errorf("failed to update installment: %w", err)
	}

	// Update installment plan progress
	err = s.updateInstallmentPlanProgress(plan.ID)
	if err != nil {
		fmt.Printf("Warning: failed to update plan progress: %v\n", err)
	}

	// Create audit record
	previousPaidAmount := installment.PaidAmount - req.Amount
	audit := &entities.InstallmentPlanAudit{
		ID:            uuid.New().String(),
		PlanID:        plan.ID,
		InstallmentID: &installment.ID,
		Action:        "installment_paid",
		PaymentAmount: &req.Amount,
		ChangedBy:     req.InitiatedBy,
		ChangeReason:  fmt.Sprintf("Payment: %.2f via %s", req.Amount, req.PaymentMethod),
		Metadata: map[string]interface{}{
			"user_id":       req.UserID,
			"previous_paid": previousPaidAmount,
			"current_paid":  installment.PaidAmount,
		},
		CreatedAt: time.Now(),
	}

	err = s.installmentAuditRepo.Create(audit)
	if err != nil {
		fmt.Printf("Warning: failed to create audit record: %v\n", err)
	}

	return updatedInstallment, nil
}

// GetInstallmentHistory retrieves audit history for an installment
func (s *InstallmentService) GetInstallmentHistory(installmentID string) ([]*entities.InstallmentPlanAudit, error) {
	return s.installmentAuditRepo.GetByInstallment(installmentID)
}

// GetInstallmentSummary retrieves installment summary for a user
func (s *InstallmentService) GetInstallmentSummary(userID string) (*dto.InstallmentSummaryResponse, error) {
	summaryData, err := s.installmentPlanRepo.GetSummaryByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get summary data: %w", err)
	}

	// Get upcoming installments
	upcomingInstallments, _, err := s.installmentRepo.GetUpcoming(userID, 30, 10, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming installments: %w", err)
	}

	upcoming := make([]dto.UpcomingInstallmentSummary, len(upcomingInstallments))
	var nextPaymentDue *time.Time
	var nextPaymentAmount float64

	for i, inst := range upcomingInstallments {
		plan, _ := s.installmentPlanRepo.GetByID(inst.PlanID)
		daysUntilDue := int(time.Until(inst.DueDate).Hours() / 24)

		upcoming[i] = dto.UpcomingInstallmentSummary{
			InstallmentID:     inst.ID,
			PlanID:            inst.PlanID,
			CardID:            plan.CardID,
			Amount:            inst.Amount,
			DueDate:           inst.DueDate,
			Description:       plan.Description,
			MerchantName:      plan.MerchantName,
			DaysUntilDue:      daysUntilDue,
			InstallmentNumber: inst.InstallmentNumber,
			TotalInstallments: plan.InstallmentsCount,
		}

		if i == 0 {
			nextPaymentDue = &inst.DueDate
			nextPaymentAmount = inst.Amount
		}
	}

	return &dto.InstallmentSummaryResponse{
		UserID:                   userID,
		TotalActivePlans:         summaryData.TotalActivePlans,
		TotalCompletedPlans:      summaryData.TotalCompletedPlans,
		TotalCancelledPlans:      summaryData.TotalCancelledPlans,
		TotalOutstandingAmount:   summaryData.TotalOutstandingAmount,
		TotalPaidAmount:          summaryData.TotalPaidAmount,
		TotalOverdueAmount:       summaryData.TotalOverdueAmount,
		NextPaymentDue:           nextPaymentDue,
		NextPaymentAmount:        nextPaymentAmount,
		OverdueInstallmentsCount: summaryData.OverdueInstallmentsCount,
		UpcomingInstallments:     upcoming,
		RecentActivity:           []dto.InstallmentActivitySummary{}, // Could be implemented
	}, nil
}

// GetMonthlyInstallmentLoad retrieves monthly installment load
func (s *InstallmentService) GetMonthlyInstallmentLoad(userID string, year, month int) (*dto.MonthlyInstallmentLoadResponse, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)

	installments, err := s.installmentRepo.GetByDueDateRange(userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get installments for date range: %w", err)
	}

	response := &dto.MonthlyInstallmentLoadResponse{
		UserID: userID,
		Year:   year,
		Month:  month,
	}

	// Calculate totals and organize by day
	dailyMap := make(map[int][]dto.UpcomingInstallmentSummary)

	for _, inst := range installments {
		plan, _ := s.installmentPlanRepo.GetByID(inst.PlanID)

		response.TotalInstallments++
		response.TotalAmount += inst.Amount

		switch inst.Status {
		case entities.InstallmentStatusPaid:
			response.PaidInstallments++
			response.PaidAmount += inst.PaidAmount
		case entities.InstallmentStatusPending:
			response.PendingInstallments++
			response.PendingAmount += inst.RemainingAmount
		case entities.InstallmentStatusOverdue:
			response.OverdueInstallments++
			response.OverdueAmount += inst.RemainingAmount
		}

		day := inst.DueDate.Day()
		summary := dto.UpcomingInstallmentSummary{
			InstallmentID:     inst.ID,
			PlanID:            inst.PlanID,
			CardID:            plan.CardID,
			Amount:            inst.Amount,
			DueDate:           inst.DueDate,
			Description:       plan.Description,
			MerchantName:      plan.MerchantName,
			InstallmentNumber: inst.InstallmentNumber,
			TotalInstallments: plan.InstallmentsCount,
		}

		dailyMap[day] = append(dailyMap[day], summary)
	}

	// Create daily breakdown
	for day := 1; day <= endDate.Day(); day++ {
		dailyInstallments := dailyMap[day]
		if dailyInstallments == nil {
			dailyInstallments = []dto.UpcomingInstallmentSummary{}
		}

		var dayTotal float64
		for _, inst := range dailyInstallments {
			dayTotal += inst.Amount
		}

		response.DailyBreakdown = append(response.DailyBreakdown, dto.DailyInstallmentLoad{
			Day:               day,
			Date:              time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC),
			InstallmentsCount: len(dailyInstallments),
			TotalAmount:       dayTotal,
			Installments:      dailyInstallments,
		})
	}

	return response, nil
}

// updateInstallmentPlanProgress updates the progress of an installment plan
func (s *InstallmentService) updateInstallmentPlanProgress(planID string) error {
	plan, err := s.installmentPlanRepo.GetByIDWithInstallments(planID)
	if err != nil {
		return err
	}

	var paidCount int
	var remainingAmount float64

	for _, installment := range plan.Installments {
		if installment.Status == entities.InstallmentStatusPaid {
			paidCount++
		} else {
			remainingAmount += installment.RemainingAmount
		}
	}

	plan.PaidInstallments = paidCount
	plan.RemainingAmount = remainingAmount

	// Check if plan is completed
	if paidCount == plan.InstallmentsCount {
		plan.Status = entities.InstallmentPlanStatusCompleted
		plan.CompletedAt = &time.Time{}
		*plan.CompletedAt = time.Now()
	}

	plan.UpdatedAt = time.Now()

	_, err = s.installmentPlanRepo.Update(plan)
	return err
}
