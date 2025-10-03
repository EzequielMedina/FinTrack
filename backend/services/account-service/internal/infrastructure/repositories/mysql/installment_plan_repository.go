package mysql

import (
	"fmt"
	"strings"
	"time"

	"github.com/fintrack/account-service/internal/core/domain/entities"
	"github.com/fintrack/account-service/internal/core/ports"
	"github.com/fintrack/account-service/internal/infrastructure/entrypoints/handlers/card/dto"
	"gorm.io/gorm"
)

// InstallmentPlanRepository implements InstallmentPlanRepositoryInterface
type InstallmentPlanRepository struct {
	db *gorm.DB
}

// NewInstallmentPlanRepository creates a new installment plan repository
func NewInstallmentPlanRepository(db *gorm.DB) ports.InstallmentPlanRepositoryInterface {
	return &InstallmentPlanRepository{db: db}
}

// Create creates a new installment plan
func (r *InstallmentPlanRepository) Create(plan *entities.InstallmentPlan) (*entities.InstallmentPlan, error) {
	if err := r.db.Create(plan).Error; err != nil {
		return nil, fmt.Errorf("failed to create installment plan: %w", err)
	}
	return plan, nil
}

// GetByID retrieves an installment plan by ID
func (r *InstallmentPlanRepository) GetByID(planID string) (*entities.InstallmentPlan, error) {
	var plan entities.InstallmentPlan
	err := r.db.Where("id = ?", planID).First(&plan).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("installment plan not found")
		}
		return nil, fmt.Errorf("failed to get installment plan: %w", err)
	}
	return &plan, nil
}

// GetByIDWithInstallments retrieves an installment plan with its installments
func (r *InstallmentPlanRepository) GetByIDWithInstallments(planID string) (*entities.InstallmentPlan, error) {
	var plan entities.InstallmentPlan
	err := r.db.Preload("Installments", func(db *gorm.DB) *gorm.DB {
		return db.Order("installment_number ASC")
	}).Where("id = ?", planID).First(&plan).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("installment plan not found")
		}
		return nil, fmt.Errorf("failed to get installment plan with installments: %w", err)
	}
	return &plan, nil
}

// GetByCard retrieves installment plans by card ID with optional status filter
func (r *InstallmentPlanRepository) GetByCard(cardID string, status string, limit, offset int) ([]*entities.InstallmentPlan, int64, error) {
	var plans []*entities.InstallmentPlan
	var total int64

	query := r.db.Model(&entities.InstallmentPlan{}).Where("card_id = ?", cardID)
	
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count installment plans: %w", err)
	}

	// Get paginated results
	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&plans).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get installment plans by card: %w", err)
	}

	return plans, total, nil
}

// GetByUser retrieves installment plans by user ID with optional status filter
func (r *InstallmentPlanRepository) GetByUser(userID string, status string, limit, offset int) ([]*entities.InstallmentPlan, int64, error) {
	var plans []*entities.InstallmentPlan
	var total int64

	query := r.db.Model(&entities.InstallmentPlan{}).
		Preload("Card").
		Where("user_id = ?", userID)
	
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count installment plans: %w", err)
	}

	// Get paginated results
	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&plans).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get installment plans by user: %w", err)
	}

	return plans, total, nil
}

// GetByTransaction retrieves an installment plan by transaction ID
func (r *InstallmentPlanRepository) GetByTransaction(transactionID string) (*entities.InstallmentPlan, error) {
	var plan entities.InstallmentPlan
	err := r.db.Where("transaction_id = ?", transactionID).First(&plan).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("installment plan not found for transaction")
		}
		return nil, fmt.Errorf("failed to get installment plan by transaction: %w", err)
	}
	return &plan, nil
}

// Update updates an installment plan
func (r *InstallmentPlanRepository) Update(plan *entities.InstallmentPlan) (*entities.InstallmentPlan, error) {
	err := r.db.Save(plan).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update installment plan: %w", err)
	}
	return plan, nil
}

// Delete soft deletes an installment plan
func (r *InstallmentPlanRepository) Delete(planID string) error {
	err := r.db.Where("id = ?", planID).Delete(&entities.InstallmentPlan{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete installment plan: %w", err)
	}
	return nil
}

// GetActiveByCard retrieves active installment plans for a card
func (r *InstallmentPlanRepository) GetActiveByCard(cardID string) ([]*entities.InstallmentPlan, error) {
	var plans []*entities.InstallmentPlan
	err := r.db.Where("card_id = ? AND status = ?", cardID, entities.InstallmentPlanStatusActive).
		Order("created_at DESC").
		Find(&plans).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get active installment plans: %w", err)
	}

	return plans, nil
}

// GetCompletedByCard retrieves completed installment plans for a card
func (r *InstallmentPlanRepository) GetCompletedByCard(cardID string, limit, offset int) ([]*entities.InstallmentPlan, int64, error) {
	var plans []*entities.InstallmentPlan
	var total int64

	query := r.db.Model(&entities.InstallmentPlan{}).
		Where("card_id = ? AND status = ?", cardID, entities.InstallmentPlanStatusCompleted)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count completed installment plans: %w", err)
	}

	// Get paginated results
	err := query.Order("completed_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&plans).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get completed installment plans: %w", err)
	}

	return plans, total, nil
}

// GetOverdueByUser retrieves installment plans with overdue installments for a user
func (r *InstallmentPlanRepository) GetOverdueByUser(userID string) ([]*entities.InstallmentPlan, error) {
	var plans []*entities.InstallmentPlan
	
	// Join with installments table to find plans with overdue installments
	err := r.db.Distinct("installment_plans.*").
		Joins("JOIN installments ON installment_plans.id = installments.plan_id").
		Where("installment_plans.user_id = ? AND installment_plans.status = ? AND installments.status = ?", 
			userID, entities.InstallmentPlanStatusActive, entities.InstallmentStatusOverdue).
		Preload("Installments", "status = ?", entities.InstallmentStatusOverdue).
		Find(&plans).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get overdue installment plans: %w", err)
	}

	return plans, nil
}

// GetSummaryByUser retrieves summary data for a user's installment plans
func (r *InstallmentPlanRepository) GetSummaryByUser(userID string) (*dto.InstallmentSummaryData, error) {
	var summary dto.InstallmentSummaryData

	// Get plan counts by status
	type StatusCount struct {
		Status string
		Count  int
	}

	var statusCounts []StatusCount
	err := r.db.Model(&entities.InstallmentPlan{}).
		Select("status, count(*) as count").
		Where("user_id = ?", userID).
		Group("status").
		Scan(&statusCounts).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get plan status counts: %w", err)
	}

	// Map status counts
	for _, sc := range statusCounts {
		switch sc.Status {
		case string(entities.InstallmentPlanStatusActive):
			summary.TotalActivePlans = sc.Count
		case string(entities.InstallmentPlanStatusCompleted):
			summary.TotalCompletedPlans = sc.Count
		case string(entities.InstallmentPlanStatusCancelled):
			summary.TotalCancelledPlans = sc.Count
		}
	}

	// Get financial summary for active plans
	type FinancialSummary struct {
		TotalOutstanding float64
		TotalPaid       float64
	}

	var financial FinancialSummary
	err = r.db.Model(&entities.InstallmentPlan{}).
		Select("SUM(remaining_amount) as total_outstanding, SUM(total_amount - remaining_amount) as total_paid").
		Where("user_id = ? AND status = ?", userID, entities.InstallmentPlanStatusActive).
		Scan(&financial).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get financial summary: %w", err)
	}

	summary.TotalOutstandingAmount = financial.TotalOutstanding
	summary.TotalPaidAmount = financial.TotalPaid

	// Get overdue information
	type OverdueSummary struct {
		OverdueAmount float64
		OverdueCount  int
	}

	var overdue OverdueSummary
	err = r.db.Model(&entities.Installment{}).
		Select("SUM(remaining_amount) as overdue_amount, COUNT(*) as overdue_count").
		Joins("JOIN installment_plans ON installments.plan_id = installment_plans.id").
		Where("installment_plans.user_id = ? AND installments.status = ?", userID, entities.InstallmentStatusOverdue).
		Scan(&overdue).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get overdue summary: %w", err)
	}

	summary.TotalOverdueAmount = overdue.OverdueAmount
	summary.OverdueInstallmentsCount = overdue.OverdueCount

	return &summary, nil
}

// GetPlansWithUpcomingPayments retrieves plans with payments due in the next N days
func (r *InstallmentPlanRepository) GetPlansWithUpcomingPayments(userID string, days int) ([]*entities.InstallmentPlan, error) {
	cutoffDate := time.Now().AddDate(0, 0, days)
	
	var plans []*entities.InstallmentPlan
	err := r.db.Distinct("installment_plans.*").
		Joins("JOIN installments ON installment_plans.id = installments.plan_id").
		Where("installment_plans.user_id = ? AND installment_plans.status = ? AND installments.status IN ? AND installments.due_date <= ?", 
			userID, entities.InstallmentPlanStatusActive, 
			[]string{string(entities.InstallmentStatusPending), string(entities.InstallmentStatusOverdue)}, 
			cutoffDate).
		Preload("Installments", "status IN ? AND due_date <= ?", 
			[]string{string(entities.InstallmentStatusPending), string(entities.InstallmentStatusOverdue)}, 
			cutoffDate).
		Order("installment_plans.created_at DESC").
		Find(&plans).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get plans with upcoming payments: %w", err)
	}

	return plans, nil
}

// SearchInstallmentPlans performs a flexible search on installment plans
func (r *InstallmentPlanRepository) SearchInstallmentPlans(userID string, filters map[string]interface{}, limit, offset int) ([]*entities.InstallmentPlan, int64, error) {
	var plans []*entities.InstallmentPlan
	var total int64

	query := r.db.Model(&entities.InstallmentPlan{}).
		Preload("Card").
		Where("user_id = ?", userID)

	// Apply filters
	if status, ok := filters["status"]; ok && status != "" {
		if statusList, ok := status.([]string); ok {
			query = query.Where("status IN ?", statusList)
		} else {
			query = query.Where("status = ?", status)
		}
	}

	if cardID, ok := filters["card_id"]; ok && cardID != "" {
		query = query.Where("card_id = ?", cardID)
	}

	if merchantName, ok := filters["merchant_name"]; ok && merchantName != "" {
		query = query.Where("merchant_name ILIKE ?", "%"+merchantName.(string)+"%")
	}

	if description, ok := filters["description"]; ok && description != "" {
		query = query.Where("description ILIKE ?", "%"+description.(string)+"%")
	}

	if amountMin, ok := filters["amount_min"]; ok {
		query = query.Where("total_amount >= ?", amountMin)
	}

	if amountMax, ok := filters["amount_max"]; ok {
		query = query.Where("total_amount <= ?", amountMax)
	}

	if dateFrom, ok := filters["date_from"]; ok {
		query = query.Where("created_at >= ?", dateFrom)
	}

	if dateTo, ok := filters["date_to"]; ok {
		query = query.Where("created_at <= ?", dateTo)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count filtered installment plans: %w", err)
	}

	// Get paginated results
	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&plans).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to search installment plans: %w", err)
	}

	return plans, total, nil
}

// BulkUpdateStatus updates status for multiple plans
func (r *InstallmentPlanRepository) BulkUpdateStatus(planIDs []string, newStatus entities.InstallmentPlanStatus, updatedBy string) error {
	if len(planIDs) == 0 {
		return nil
	}

	updates := map[string]interface{}{
		"status":     newStatus,
		"updated_at": time.Now(),
	}

	if newStatus == entities.InstallmentPlanStatusCompleted {
		updates["completed_at"] = time.Now()
	} else if newStatus == entities.InstallmentPlanStatusCancelled {
		updates["cancelled_at"] = time.Now()
	}

	err := r.db.Model(&entities.InstallmentPlan{}).
		Where("id IN ?", planIDs).
		Updates(updates).Error

	if err != nil {
		return fmt.Errorf("failed to bulk update installment plan status: %w", err)
	}

	return nil
}

// GetMonthlyBreakdown retrieves installment plans breakdown by month
func (r *InstallmentPlanRepository) GetMonthlyBreakdown(userID string, year int) (map[int]dto.MonthlyInstallmentLoadResponse, error) {
	type MonthlyData struct {
		Month               int     `json:"month"`
		TotalInstallments   int     `json:"total_installments"`
		TotalAmount         float64 `json:"total_amount"`
		PaidInstallments    int     `json:"paid_installments"`
		PaidAmount          float64 `json:"paid_amount"`
		PendingInstallments int     `json:"pending_installments"`
		PendingAmount       float64 `json:"pending_amount"`
		OverdueInstallments int     `json:"overdue_installments"`
		OverdueAmount       float64 `json:"overdue_amount"`
	}

	var monthlyData []MonthlyData
	
	// Complex query to get monthly breakdown
	query := `
		SELECT 
			MONTH(i.due_date) as month,
			COUNT(*) as total_installments,
			SUM(i.amount) as total_amount,
			SUM(CASE WHEN i.status = 'paid' THEN 1 ELSE 0 END) as paid_installments,
			SUM(CASE WHEN i.status = 'paid' THEN i.amount ELSE 0 END) as paid_amount,
			SUM(CASE WHEN i.status = 'pending' THEN 1 ELSE 0 END) as pending_installments,
			SUM(CASE WHEN i.status = 'pending' THEN i.remaining_amount ELSE 0 END) as pending_amount,
			SUM(CASE WHEN i.status = 'overdue' THEN 1 ELSE 0 END) as overdue_installments,
			SUM(CASE WHEN i.status = 'overdue' THEN i.remaining_amount ELSE 0 END) as overdue_amount
		FROM installments i
		JOIN installment_plans ip ON i.plan_id = ip.id
		WHERE ip.user_id = ? AND YEAR(i.due_date) = ?
		GROUP BY MONTH(i.due_date)
		ORDER BY month
	`

	err := r.db.Raw(query, userID, year).Scan(&monthlyData).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly breakdown: %w", err)
	}

	// Convert to map
	result := make(map[int]dto.MonthlyInstallmentLoadResponse)
	for _, data := range monthlyData {
		result[data.Month] = dto.MonthlyInstallmentLoadResponse{
			UserID:              userID,
			Year:                year,
			Month:               data.Month,
			TotalInstallments:   data.TotalInstallments,
			TotalAmount:         data.TotalAmount,
			PaidInstallments:    data.PaidInstallments,
			PaidAmount:          data.PaidAmount,
			PendingInstallments: data.PendingInstallments,
			PendingAmount:       data.PendingAmount,
			OverdueInstallments: data.OverdueInstallments,
			OverdueAmount:       data.OverdueAmount,
		}
	}

	return result, nil
}