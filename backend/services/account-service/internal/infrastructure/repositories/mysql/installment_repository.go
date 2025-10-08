package mysql

import (
	"fmt"
	"time"

	"github.com/fintrack/account-service/internal/core/domain/entities"
	"github.com/fintrack/account-service/internal/core/ports"
	"gorm.io/gorm"
)

// InstallmentRepository implements InstallmentRepositoryInterface
type InstallmentRepository struct {
	db *gorm.DB
}

// NewInstallmentRepository creates a new installment repository
func NewInstallmentRepository(db *gorm.DB) ports.InstallmentRepositoryInterface {
	return &InstallmentRepository{db: db}
}

// Create creates a new installment
func (r *InstallmentRepository) Create(installment *entities.Installment) (*entities.Installment, error) {
	if err := r.db.Create(installment).Error; err != nil {
		return nil, fmt.Errorf("failed to create installment: %w", err)
	}
	return installment, nil
}

// GetByID retrieves an installment by ID
func (r *InstallmentRepository) GetByID(installmentID string) (*entities.Installment, error) {
	var installment entities.Installment
	err := r.db.Preload("Plan").
		Where("id = ?", installmentID).
		First(&installment).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("installment not found")
		}
		return nil, fmt.Errorf("failed to get installment: %w", err)
	}
	return &installment, nil
}

// GetByPlan retrieves all installments for a plan
func (r *InstallmentRepository) GetByPlan(planID string) ([]*entities.Installment, error) {
	var installments []*entities.Installment
	err := r.db.Where("plan_id = ?", planID).
		Order("installment_number ASC").
		Find(&installments).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get installments by plan: %w", err)
	}
	return installments, nil
}

// GetByPlanAndNumber retrieves a specific installment by plan and number
func (r *InstallmentRepository) GetByPlanAndNumber(planID string, installmentNumber int) (*entities.Installment, error) {
	var installment entities.Installment
	err := r.db.Where("plan_id = ? AND installment_number = ?", planID, installmentNumber).
		First(&installment).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("installment not found")
		}
		return nil, fmt.Errorf("failed to get installment by plan and number: %w", err)
	}
	return &installment, nil
}

// Update updates an installment
func (r *InstallmentRepository) Update(installment *entities.Installment) (*entities.Installment, error) {
	err := r.db.Save(installment).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update installment: %w", err)
	}
	return installment, nil
}

// Delete soft deletes an installment
func (r *InstallmentRepository) Delete(installmentID string) error {
	err := r.db.Where("id = ?", installmentID).Delete(&entities.Installment{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete installment: %w", err)
	}
	return nil
}

// GetOverdue retrieves overdue installments for a user
func (r *InstallmentRepository) GetOverdue(userID string, limit, offset int) ([]*entities.Installment, int64, error) {
	var installments []*entities.Installment
	var total int64

	query := r.db.Model(&entities.Installment{}).
		Joins("JOIN installment_plans ON installments.plan_id = installment_plans.id").
		Where("installment_plans.user_id = ? AND installments.status = ?", userID, entities.InstallmentStatusOverdue).
		Preload("Plan").
		Preload("Plan.Card")

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count overdue installments: %w", err)
	}

	// Get paginated results
	err := query.Order("installments.due_date ASC").
		Limit(limit).
		Offset(offset).
		Find(&installments).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get overdue installments: %w", err)
	}

	return installments, total, nil
}

// GetUpcoming retrieves upcoming installments for a user (due within specified days)
func (r *InstallmentRepository) GetUpcoming(userID string, days int, limit, offset int) ([]*entities.Installment, int64, error) {
	var installments []*entities.Installment
	var total int64

	cutoffDate := time.Now().AddDate(0, 0, days)

	query := r.db.Model(&entities.Installment{}).
		Joins("JOIN installment_plans ON installments.plan_id = installment_plans.id").
		Where("installment_plans.user_id = ? AND installments.status IN ? AND installments.due_date <= ?",
			userID,
			[]string{string(entities.InstallmentStatusPending), string(entities.InstallmentStatusOverdue)},
			cutoffDate).
		Preload("Plan").
		Preload("Plan.Card")

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count upcoming installments: %w", err)
	}

	// Get paginated results
	err := query.Order("installments.due_date ASC").
		Limit(limit).
		Offset(offset).
		Find(&installments).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get upcoming installments: %w", err)
	}

	return installments, total, nil
}

// GetByDueDateRange retrieves installments within a date range for a user
func (r *InstallmentRepository) GetByDueDateRange(userID string, startDate, endDate time.Time) ([]*entities.Installment, error) {
	var installments []*entities.Installment

	err := r.db.Joins("JOIN installment_plans ON installments.plan_id = installment_plans.id").
		Where("installment_plans.user_id = ? AND installments.due_date BETWEEN ? AND ?",
			userID, startDate, endDate).
		Preload("Plan").
		Preload("Plan.Card").
		Order("installments.due_date ASC").
		Find(&installments).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get installments by date range: %w", err)
	}

	return installments, nil
}

// GetPendingByPlan retrieves pending installments for a plan
func (r *InstallmentRepository) GetPendingByPlan(planID string) ([]*entities.Installment, error) {
	var installments []*entities.Installment

	err := r.db.Where("plan_id = ? AND status IN ?",
		planID,
		[]string{string(entities.InstallmentStatusPending), string(entities.InstallmentStatusOverdue)}).
		Order("installment_number ASC").
		Find(&installments).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get pending installments: %w", err)
	}

	return installments, nil
}

// GetNextDueByPlan retrieves the next due installment for a plan
func (r *InstallmentRepository) GetNextDueByPlan(planID string) (*entities.Installment, error) {
	var installment entities.Installment

	err := r.db.Where("plan_id = ? AND status IN ?",
		planID,
		[]string{string(entities.InstallmentStatusPending), string(entities.InstallmentStatusOverdue)}).
		Order("installment_number ASC").
		First(&installment).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No pending installments
		}
		return nil, fmt.Errorf("failed to get next due installment: %w", err)
	}

	return &installment, nil
}

// MarkOverdue marks installments as overdue based on cutoff date
func (r *InstallmentRepository) MarkOverdue(cutoffDate time.Time) (int64, error) {
	result := r.db.Model(&entities.Installment{}).
		Where("status = ? AND due_date < ?", entities.InstallmentStatusPending, cutoffDate).
		Update("status", entities.InstallmentStatusOverdue)

	if result.Error != nil {
		return 0, fmt.Errorf("failed to mark installments as overdue: %w", result.Error)
	}

	return result.RowsAffected, nil
}

// GetInstallmentsByStatus retrieves installments by status for a user
func (r *InstallmentRepository) GetInstallmentsByStatus(userID string, status entities.InstallmentStatus, limit, offset int) ([]*entities.Installment, int64, error) {
	var installments []*entities.Installment
	var total int64

	query := r.db.Model(&entities.Installment{}).
		Joins("JOIN installment_plans ON installments.plan_id = installment_plans.id").
		Where("installment_plans.user_id = ? AND installments.status = ?", userID, status).
		Preload("Plan").
		Preload("Plan.Card")

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count installments by status: %w", err)
	}

	// Get paginated results
	err := query.Order("installments.due_date ASC").
		Limit(limit).
		Offset(offset).
		Find(&installments).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get installments by status: %w", err)
	}

	return installments, total, nil
}

// BulkCreate creates multiple installments in a transaction
func (r *InstallmentRepository) BulkCreate(installments []*entities.Installment) error {
	if len(installments) == 0 {
		return nil
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, installment := range installments {
			if err := tx.Create(installment).Error; err != nil {
				return fmt.Errorf("failed to create installment %d: %w", installment.InstallmentNumber, err)
			}
		}
		return nil
	})
}

// BulkUpdateStatus updates status for multiple installments
func (r *InstallmentRepository) BulkUpdateStatus(installmentIDs []string, newStatus entities.InstallmentStatus) error {
	if len(installmentIDs) == 0 {
		return nil
	}

	updates := map[string]interface{}{
		"status":     newStatus,
		"updated_at": time.Now(),
	}

	if newStatus == entities.InstallmentStatusPaid {
		updates["paid_date"] = time.Now()
	}

	err := r.db.Model(&entities.Installment{}).
		Where("id IN ?", installmentIDs).
		Updates(updates).Error

	if err != nil {
		return fmt.Errorf("failed to bulk update installment status: %w", err)
	}

	return nil
}

// GetPaymentHistory retrieves payment history for an installment
func (r *InstallmentRepository) GetPaymentHistory(installmentID string) ([]*entities.InstallmentPlanAudit, error) {
	var audits []*entities.InstallmentPlanAudit

	err := r.db.Where("installment_id = ?", installmentID).
		Order("created_at DESC").
		Find(&audits).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get payment history: %w", err)
	}

	return audits, nil
}

// GetDailyInstallments retrieves installments due on a specific date
func (r *InstallmentRepository) GetDailyInstallments(userID string, date time.Time) ([]*entities.Installment, error) {
	var installments []*entities.Installment

	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour).Add(-1 * time.Nanosecond)

	err := r.db.Joins("JOIN installment_plans ON installments.plan_id = installment_plans.id").
		Where("installment_plans.user_id = ? AND installments.due_date BETWEEN ? AND ?",
			userID, startOfDay, endOfDay).
		Preload("Plan").
		Preload("Plan.Card").
		Order("installments.installment_number ASC").
		Find(&installments).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get daily installments: %w", err)
	}

	return installments, nil
}

// GetInstallmentStatistics retrieves statistics for installments
func (r *InstallmentRepository) GetInstallmentStatistics(userID string) (map[string]interface{}, error) {
	type Stats struct {
		TotalInstallments    int        `json:"total_installments"`
		PaidInstallments     int        `json:"paid_installments"`
		PendingInstallments  int        `json:"pending_installments"`
		OverdueInstallments  int        `json:"overdue_installments"`
		TotalAmount          float64    `json:"total_amount"`
		PaidAmount           float64    `json:"paid_amount"`
		PendingAmount        float64    `json:"pending_amount"`
		OverdueAmount        float64    `json:"overdue_amount"`
		AvgInstallmentAmount float64    `json:"avg_installment_amount"`
		EarliestDueDate      *time.Time `json:"earliest_due_date"`
		LatestDueDate        *time.Time `json:"latest_due_date"`
	}

	var stats Stats

	query := `
		SELECT 
			COUNT(*) as total_installments,
			SUM(CASE WHEN status = 'paid' THEN 1 ELSE 0 END) as paid_installments,
			SUM(CASE WHEN status = 'pending' THEN 1 ELSE 0 END) as pending_installments,
			SUM(CASE WHEN status = 'overdue' THEN 1 ELSE 0 END) as overdue_installments,
			SUM(amount) as total_amount,
			SUM(CASE WHEN status = 'paid' THEN paid_amount ELSE 0 END) as paid_amount,
			SUM(CASE WHEN status = 'pending' THEN remaining_amount ELSE 0 END) as pending_amount,
			SUM(CASE WHEN status = 'overdue' THEN remaining_amount ELSE 0 END) as overdue_amount,
			AVG(amount) as avg_installment_amount,
			MIN(CASE WHEN status IN ('pending', 'overdue') THEN due_date END) as earliest_due_date,
			MAX(due_date) as latest_due_date
		FROM installments i
		JOIN installment_plans ip ON i.plan_id = ip.id
		WHERE ip.user_id = ?
	`

	err := r.db.Raw(query, userID).Scan(&stats).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get installment statistics: %w", err)
	}

	// Convert to map for flexible response
	result := map[string]interface{}{
		"total_installments":     stats.TotalInstallments,
		"paid_installments":      stats.PaidInstallments,
		"pending_installments":   stats.PendingInstallments,
		"overdue_installments":   stats.OverdueInstallments,
		"total_amount":           stats.TotalAmount,
		"paid_amount":            stats.PaidAmount,
		"pending_amount":         stats.PendingAmount,
		"overdue_amount":         stats.OverdueAmount,
		"avg_installment_amount": stats.AvgInstallmentAmount,
		"earliest_due_date":      stats.EarliestDueDate,
		"latest_due_date":        stats.LatestDueDate,
	}

	return result, nil
}

// GetInstallmentsByCard retrieves installments for a specific card
func (r *InstallmentRepository) GetInstallmentsByCard(cardID string, status string, limit, offset int) ([]*entities.Installment, int64, error) {
	var installments []*entities.Installment
	var total int64

	query := r.db.Model(&entities.Installment{}).
		Joins("JOIN installment_plans ON installments.plan_id = installment_plans.id").
		Where("installment_plans.card_id = ?", cardID).
		Preload("Plan")

	if status != "" {
		query = query.Where("installments.status = ?", status)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count installments by card: %w", err)
	}

	// Get paginated results
	err := query.Order("installments.due_date ASC").
		Limit(limit).
		Offset(offset).
		Find(&installments).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get installments by card: %w", err)
	}

	return installments, total, nil
}
