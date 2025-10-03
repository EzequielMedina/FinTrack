package mysql

import (
	"fmt"
	"time"

	"github.com/fintrack/account-service/internal/core/domain/entities"
	"github.com/fintrack/account-service/internal/core/ports"
	"gorm.io/gorm"
)

// InstallmentPlanAuditRepository implements InstallmentPlanAuditRepositoryInterface
type InstallmentPlanAuditRepository struct {
	db *gorm.DB
}

// NewInstallmentPlanAuditRepository creates a new installment plan audit repository
func NewInstallmentPlanAuditRepository(db *gorm.DB) ports.InstallmentPlanAuditRepositoryInterface {
	return &InstallmentPlanAuditRepository{db: db}
}

// Create creates a new audit record
func (r *InstallmentPlanAuditRepository) Create(audit *entities.InstallmentPlanAudit) error {
	if err := r.db.Create(audit).Error; err != nil {
		return fmt.Errorf("failed to create audit record: %w", err)
	}
	return nil
}

// GetByPlan retrieves audit records for a specific installment plan
func (r *InstallmentPlanAuditRepository) GetByPlan(planID string, limit, offset int) ([]*entities.InstallmentPlanAudit, int64, error) {
	var audits []*entities.InstallmentPlanAudit
	var total int64

	query := r.db.Model(&entities.InstallmentPlanAudit{}).Where("plan_id = ?", planID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count audit records: %w", err)
	}

	// Get paginated results
	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&audits).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get audit records by plan: %w", err)
	}

	return audits, total, nil
}

// GetByInstallment retrieves audit records for a specific installment
func (r *InstallmentPlanAuditRepository) GetByInstallment(installmentID string) ([]*entities.InstallmentPlanAudit, error) {
	var audits []*entities.InstallmentPlanAudit
	
	err := r.db.Where("installment_id = ?", installmentID).
		Order("created_at DESC").
		Find(&audits).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get audit records by installment: %w", err)
	}

	return audits, nil
}

// GetByUser retrieves audit records for a user with optional action filter
func (r *InstallmentPlanAuditRepository) GetByUser(userID string, action string, limit, offset int) ([]*entities.InstallmentPlanAudit, int64, error) {
	var audits []*entities.InstallmentPlanAudit
	var total int64

	query := r.db.Model(&entities.InstallmentPlanAudit{}).
		Joins("JOIN installment_plans ON installment_plan_audit.plan_id = installment_plans.id").
		Where("installment_plans.user_id = ?", userID)

	if action != "" {
		query = query.Where("installment_plan_audit.action = ?", action)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count audit records by user: %w", err)
	}

	// Get paginated results
	err := query.Order("installment_plan_audit.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&audits).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get audit records by user: %w", err)
	}

	return audits, total, nil
}

// GetByDateRange retrieves audit records within a date range
func (r *InstallmentPlanAuditRepository) GetByDateRange(startDate, endDate time.Time, limit, offset int) ([]*entities.InstallmentPlanAudit, int64, error) {
	var audits []*entities.InstallmentPlanAudit
	var total int64

	query := r.db.Model(&entities.InstallmentPlanAudit{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count audit records by date range: %w", err)
	}

	// Get paginated results
	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&audits).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get audit records by date range: %w", err)
	}

	return audits, total, nil
}

// CreatePaymentAudit creates an audit record for payment activities
func (r *InstallmentPlanAuditRepository) CreatePaymentAudit(planID, installmentID string, oldStatus, newStatus entities.InstallmentStatus, paymentAmount float64, changedBy, reason string) error {
	audit := &entities.InstallmentPlanAudit{
		PlanID:        planID,
		Action:        "payment_applied",
		InstallmentID: &installmentID,
		PaymentAmount: &paymentAmount,
		ChangedBy:     changedBy,
		ChangeReason:  reason,
		Metadata: map[string]interface{}{
			"old_status": string(oldStatus),
			"new_status": string(newStatus),
		},
	}

	return r.Create(audit)
}

// CreateStatusChangeAudit creates an audit record for status changes
func (r *InstallmentPlanAuditRepository) CreateStatusChangeAudit(planID string, oldStatus, newStatus entities.InstallmentPlanStatus, changedBy, reason string) error {
	oldStatusStr := string(oldStatus)
	newStatusStr := string(newStatus)
	
	audit := &entities.InstallmentPlanAudit{
		PlanID:       planID,
		Action:       "status_changed",
		OldStatus:    &oldStatusStr,
		NewStatus:    &newStatusStr,
		ChangedBy:    changedBy,
		ChangeReason: reason,
	}

	return r.Create(audit)
}

// CreatePlanCreationAudit creates an audit record for plan creation
func (r *InstallmentPlanAuditRepository) CreatePlanCreationAudit(planID, createdBy string, metadata map[string]interface{}) error {
	audit := &entities.InstallmentPlanAudit{
		PlanID:       planID,
		Action:       "created",
		ChangedBy:    createdBy,
		ChangeReason: "Installment plan created",
		Metadata:     metadata,
	}

	return r.Create(audit)
}

// CreatePlanCancellationAudit creates an audit record for plan cancellation
func (r *InstallmentPlanAuditRepository) CreatePlanCancellationAudit(planID, cancelledBy, reason string) error {
	oldStatus := string(entities.InstallmentPlanStatusActive)
	newStatus := string(entities.InstallmentPlanStatusCancelled)
	
	audit := &entities.InstallmentPlanAudit{
		PlanID:       planID,
		Action:       "cancelled",
		OldStatus:    &oldStatus,
		NewStatus:    &newStatus,
		ChangedBy:    cancelledBy,
		ChangeReason: reason,
	}

	return r.Create(audit)
}

// CreatePlanCompletionAudit creates an audit record for plan completion
func (r *InstallmentPlanAuditRepository) CreatePlanCompletionAudit(planID string, completedBy string) error {
	oldStatus := string(entities.InstallmentPlanStatusActive)
	newStatus := string(entities.InstallmentPlanStatusCompleted)
	
	audit := &entities.InstallmentPlanAudit{
		PlanID:       planID,
		Action:       "completed",
		OldStatus:    &oldStatus,
		NewStatus:    &newStatus,
		ChangedBy:    completedBy,
		ChangeReason: "All installments paid",
	}

	return r.Create(audit)
}

// GetAuditSummary retrieves audit summary statistics
func (r *InstallmentPlanAuditRepository) GetAuditSummary(userID string, dateFrom, dateTo time.Time) (map[string]interface{}, error) {
	type AuditSummary struct {
		Action string `json:"action"`
		Count  int    `json:"count"`
	}

	var summaries []AuditSummary
	
	query := `
		SELECT 
			ipa.action,
			COUNT(*) as count
		FROM installment_plan_audit ipa
		JOIN installment_plans ip ON ipa.plan_id = ip.id
		WHERE ip.user_id = ? AND ipa.created_at BETWEEN ? AND ?
		GROUP BY ipa.action
		ORDER BY count DESC
	`

	err := r.db.Raw(query, userID, dateFrom, dateTo).Scan(&summaries).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get audit summary: %w", err)
	}

	// Convert to map
	result := make(map[string]interface{})
	for _, summary := range summaries {
		result[summary.Action] = summary.Count
	}

	return result, nil
}

// GetRecentActivity retrieves recent audit activity for a user
func (r *InstallmentPlanAuditRepository) GetRecentActivity(userID string, limit int) ([]*entities.InstallmentPlanAudit, error) {
	var audits []*entities.InstallmentPlanAudit
	
	err := r.db.Joins("JOIN installment_plans ON installment_plan_audit.plan_id = installment_plans.id").
		Where("installment_plans.user_id = ?", userID).
		Order("installment_plan_audit.created_at DESC").
		Limit(limit).
		Find(&audits).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get recent activity: %w", err)
	}

	return audits, nil
}

// GetAuditByAction retrieves audit records by specific action
func (r *InstallmentPlanAuditRepository) GetAuditByAction(action string, limit, offset int) ([]*entities.InstallmentPlanAudit, int64, error) {
	var audits []*entities.InstallmentPlanAudit
	var total int64

	query := r.db.Model(&entities.InstallmentPlanAudit{}).Where("action = ?", action)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count audit records by action: %w", err)
	}

	// Get paginated results
	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&audits).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get audit records by action: %w", err)
	}

	return audits, total, nil
}

// BulkCreateAudit creates multiple audit records in a transaction
func (r *InstallmentPlanAuditRepository) BulkCreateAudit(audits []*entities.InstallmentPlanAudit) error {
	if len(audits) == 0 {
		return nil
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, audit := range audits {
			if err := tx.Create(audit).Error; err != nil {
				return fmt.Errorf("failed to create audit record for plan %s: %w", audit.PlanID, err)
			}
		}
		return nil
	})
}

// DeleteOldAuditRecords deletes audit records older than specified date (for maintenance)
func (r *InstallmentPlanAuditRepository) DeleteOldAuditRecords(cutoffDate time.Time) (int64, error) {
	result := r.db.Where("created_at < ?", cutoffDate).Delete(&entities.InstallmentPlanAudit{})
	
	if result.Error != nil {
		return 0, fmt.Errorf("failed to delete old audit records: %w", result.Error)
	}

	return result.RowsAffected, nil
}

// GetAuditsByIPAddress retrieves audit records by IP address (for security analysis)
func (r *InstallmentPlanAuditRepository) GetAuditsByIPAddress(ipAddress string, limit, offset int) ([]*entities.InstallmentPlanAudit, int64, error) {
	var audits []*entities.InstallmentPlanAudit
	var total int64

	query := r.db.Model(&entities.InstallmentPlanAudit{}).Where("ip_address = ?", ipAddress)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count audit records by IP: %w", err)
	}

	// Get paginated results
	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&audits).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get audit records by IP: %w", err)
	}

	return audits, total, nil
}

// GetAuditStatistics retrieves comprehensive audit statistics
func (r *InstallmentPlanAuditRepository) GetAuditStatistics(dateFrom, dateTo time.Time) (map[string]interface{}, error) {
	type Stats struct {
		TotalAudits          int64  `json:"total_audits"`
		UniqueUsers          int64  `json:"unique_users"`
		UniquePlans          int64  `json:"unique_plans"`
		MostCommonAction     string `json:"most_common_action"`
		MostActiveUser       string `json:"most_active_user"`
		PaymentAudits        int64  `json:"payment_audits"`
		StatusChangeAudits   int64  `json:"status_change_audits"`
		CreationAudits       int64  `json:"creation_audits"`
		CancellationAudits   int64  `json:"cancellation_audits"`
	}

	var stats Stats
	
	// Get basic counts
	err := r.db.Model(&entities.InstallmentPlanAudit{}).
		Where("created_at BETWEEN ? AND ?", dateFrom, dateTo).
		Count(&stats.TotalAudits).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get total audit count: %w", err)
	}

	// Get unique counts
	r.db.Model(&entities.InstallmentPlanAudit{}).
		Joins("JOIN installment_plans ON installment_plan_audit.plan_id = installment_plans.id").
		Where("installment_plan_audit.created_at BETWEEN ? AND ?", dateFrom, dateTo).
		Distinct("installment_plans.user_id").
		Count(&stats.UniqueUsers)

	r.db.Model(&entities.InstallmentPlanAudit{}).
		Where("created_at BETWEEN ? AND ?", dateFrom, dateTo).
		Distinct("plan_id").
		Count(&stats.UniquePlans)

	// Get action-specific counts
	r.db.Model(&entities.InstallmentPlanAudit{}).
		Where("created_at BETWEEN ? AND ? AND action = ?", dateFrom, dateTo, "payment_applied").
		Count(&stats.PaymentAudits)

	r.db.Model(&entities.InstallmentPlanAudit{}).
		Where("created_at BETWEEN ? AND ? AND action = ?", dateFrom, dateTo, "status_changed").
		Count(&stats.StatusChangeAudits)

	r.db.Model(&entities.InstallmentPlanAudit{}).
		Where("created_at BETWEEN ? AND ? AND action = ?", dateFrom, dateTo, "created").
		Count(&stats.CreationAudits)

	r.db.Model(&entities.InstallmentPlanAudit{}).
		Where("created_at BETWEEN ? AND ? AND action = ?", dateFrom, dateTo, "cancelled").
		Count(&stats.CancellationAudits)

	// Get most common action
	type ActionCount struct {
		Action string
		Count  int
	}
	var mostCommon ActionCount
	r.db.Model(&entities.InstallmentPlanAudit{}).
		Select("action, COUNT(*) as count").
		Where("created_at BETWEEN ? AND ?", dateFrom, dateTo).
		Group("action").
		Order("count DESC").
		Limit(1).
		Scan(&mostCommon)
	stats.MostCommonAction = mostCommon.Action

	// Get most active user
	type UserActivity struct {
		UserID string
		Count  int
	}
	var mostActive UserActivity
	r.db.Model(&entities.InstallmentPlanAudit{}).
		Select("installment_plans.user_id, COUNT(*) as count").
		Joins("JOIN installment_plans ON installment_plan_audit.plan_id = installment_plans.id").
		Where("installment_plan_audit.created_at BETWEEN ? AND ?", dateFrom, dateTo).
		Group("installment_plans.user_id").
		Order("count DESC").
		Limit(1).
		Scan(&mostActive)
	stats.MostActiveUser = mostActive.UserID

	// Convert to map
	result := map[string]interface{}{
		"total_audits":          stats.TotalAudits,
		"unique_users":          stats.UniqueUsers,
		"unique_plans":          stats.UniquePlans,
		"most_common_action":    stats.MostCommonAction,
		"most_active_user":      stats.MostActiveUser,
		"payment_audits":        stats.PaymentAudits,
		"status_change_audits":  stats.StatusChangeAudits,
		"creation_audits":       stats.CreationAudits,
		"cancellation_audits":   stats.CancellationAudits,
	}

	return result, nil
}