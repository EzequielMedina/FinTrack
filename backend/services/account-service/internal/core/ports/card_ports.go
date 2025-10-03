package ports

import (
	"time"
	"github.com/fintrack/account-service/internal/core/domain/entities"
	"github.com/fintrack/account-service/internal/infrastructure/entrypoints/handlers/card/dto"
)

// CardServiceInterface defines the contract for card service operations
type CardServiceInterface interface {
	// Basic CRUD operations
	CreateCard(req *dto.CreateCardRequest) (*entities.Card, error)
	GetCardByID(cardID string) (*entities.Card, error)
	GetCardByIDWithAccount(cardID string) (*entities.Card, error) // New: get card with account preloaded
	GetCardsByAccount(accountID string, page, pageSize int) ([]*entities.Card, int64, error)
	GetCardsByUser(userID string, page, pageSize int) ([]*entities.Card, int64, error)
	UpdateCard(cardID string, req *dto.UpdateCardRequest) (*entities.Card, error)
	DeleteCard(cardID string) error
	BlockCard(cardID string) (*entities.Card, error)
	UnblockCard(cardID string) (*entities.Card, error)
	SetDefaultCard(cardID string) (*entities.Card, error)

	// Credit card financial operations
	ChargeCard(cardID string, amount float64, description, reference string) (*entities.Card, error)
	PaymentCard(cardID string, amount float64, paymentMethod, reference string) (*entities.Card, error)

	// Debit card operations
	ProcessDebitTransaction(cardID string, amount float64, description, merchantName, reference string) (*entities.Card, error)
}

// InstallmentServiceInterface defines the contract for installment service operations
type InstallmentServiceInterface interface {
	// Installment plan operations
	CalculateInstallmentPlan(amount float64, installmentsCount int, startDate time.Time, interestRate float64) (*dto.InstallmentPreviewResponse, error)
	CreateInstallmentPlan(req *dto.CreateInstallmentPlanRequest) (*entities.InstallmentPlan, error)
	GetInstallmentPlan(planID string) (*entities.InstallmentPlan, error)
	GetInstallmentPlansByCard(cardID string, page, pageSize int) ([]*entities.InstallmentPlan, int64, error)
	GetInstallmentPlansByUser(userID string, status string, page, pageSize int) ([]*entities.InstallmentPlan, int64, error)
	CancelInstallmentPlan(planID, reason string, cancelledBy string) (*entities.InstallmentPlan, error)
	SuspendInstallmentPlan(planID, reason string, suspendedBy string) (*entities.InstallmentPlan, error)
	ReactivateInstallmentPlan(planID, reason string, reactivatedBy string) (*entities.InstallmentPlan, error)

	// Individual installment operations
	GetInstallment(installmentID string) (*entities.Installment, error)
	GetInstallmentsByPlan(planID string) ([]*entities.Installment, error)
	GetOverdueInstallments(userID string, limit, offset int) ([]*entities.Installment, int64, error)
	GetUpcomingInstallments(userID string, days int, limit, offset int) ([]*entities.Installment, int64, error)
	PayInstallment(req *dto.PayInstallmentRequest) (*entities.Installment, error)
	GetInstallmentHistory(installmentID string) ([]*entities.InstallmentPlanAudit, error)

	// Reporting and analytics
	GetInstallmentSummary(userID string) (*dto.InstallmentSummaryResponse, error)
	GetMonthlyInstallmentLoad(userID string, year, month int) (*dto.MonthlyInstallmentLoadResponse, error)
}

// CardRepositoryInterface defines the contract for card repository operations
type CardRepositoryInterface interface {
	Create(card *entities.Card) (*entities.Card, error)
	GetByID(cardID string) (*entities.Card, error)
	GetByIDWithAccount(cardID string) (*entities.Card, error) // New: get card with account preloaded
	GetWithInstallmentPlans(cardID string) (*entities.Card, error) // New: get card with installment plans preloaded
	GetByAccount(accountID string, limit, offset int) ([]*entities.Card, int64, error)
	GetByAccountWithInstallmentPlans(accountID string, limit, offset int) ([]*entities.Card, int64, error) // New: get cards with installment plans
	GetByUser(userID string, limit, offset int) ([]*entities.Card, int64, error)
	Update(card *entities.Card) (*entities.Card, error)
	Delete(cardID string) error
	GetDefaultByAccount(accountID string) (*entities.Card, error)
	SetDefaultByAccount(accountID, cardID string) error
}

// InstallmentPlanRepositoryInterface defines the contract for installment plan repository operations
type InstallmentPlanRepositoryInterface interface {
	Create(plan *entities.InstallmentPlan) (*entities.InstallmentPlan, error)
	GetByID(planID string) (*entities.InstallmentPlan, error)
	GetByIDWithInstallments(planID string) (*entities.InstallmentPlan, error)
	GetByCard(cardID string, status string, limit, offset int) ([]*entities.InstallmentPlan, int64, error)
	GetByUser(userID string, status string, limit, offset int) ([]*entities.InstallmentPlan, int64, error)
	GetByTransaction(transactionID string) (*entities.InstallmentPlan, error)
	Update(plan *entities.InstallmentPlan) (*entities.InstallmentPlan, error)
	Delete(planID string) error
	GetActiveByCard(cardID string) ([]*entities.InstallmentPlan, error)
	GetCompletedByCard(cardID string, limit, offset int) ([]*entities.InstallmentPlan, int64, error)
	GetOverdueByUser(userID string) ([]*entities.InstallmentPlan, error)
	GetSummaryByUser(userID string) (*dto.InstallmentSummaryData, error)
}

// InstallmentRepositoryInterface defines the contract for installment repository operations
type InstallmentRepositoryInterface interface {
	Create(installment *entities.Installment) (*entities.Installment, error)
	GetByID(installmentID string) (*entities.Installment, error)
	GetByPlan(planID string) ([]*entities.Installment, error)
	GetByPlanAndNumber(planID string, installmentNumber int) (*entities.Installment, error)
	Update(installment *entities.Installment) (*entities.Installment, error)
	Delete(installmentID string) error
	GetOverdue(userID string, limit, offset int) ([]*entities.Installment, int64, error)
	GetUpcoming(userID string, days int, limit, offset int) ([]*entities.Installment, int64, error)
	GetByDueDateRange(userID string, startDate, endDate time.Time) ([]*entities.Installment, error)
	GetPendingByPlan(planID string) ([]*entities.Installment, error)
	GetNextDueByPlan(planID string) (*entities.Installment, error)
	MarkOverdue(cutoffDate time.Time) (int64, error)
}

// InstallmentPlanAuditRepositoryInterface defines the contract for audit repository operations
type InstallmentPlanAuditRepositoryInterface interface {
	Create(audit *entities.InstallmentPlanAudit) error
	GetByPlan(planID string, limit, offset int) ([]*entities.InstallmentPlanAudit, int64, error)
	GetByInstallment(installmentID string) ([]*entities.InstallmentPlanAudit, error)
	GetByUser(userID string, action string, limit, offset int) ([]*entities.InstallmentPlanAudit, int64, error)
	GetByDateRange(startDate, endDate time.Time, limit, offset int) ([]*entities.InstallmentPlanAudit, int64, error)
}

// AccountRepositoryInterface defines the contract for account repository operations
type AccountRepositoryInterface interface {
	Create(account *entities.Account) error
	GetByID(id string) (*entities.Account, error)
	GetByUserID(userID string) ([]*entities.Account, error)
	GetAll(limit, offset int) ([]*entities.Account, int64, error)
	Update(account *entities.Account) error
	Delete(id string) error
}
