package service

import (
	"github.com/fintrack/account-service/internal/core/domain/entities"
	"github.com/fintrack/account-service/internal/infrastructure/entrypoints/handlers/account/dto"
)

// AccountServiceInterface defines the contract for account service operations
type AccountServiceInterface interface {
	// Basic CRUD operations
	CreateAccount(account *entities.Account) (*entities.Account, error)
	GetAccountByID(accountID string) (*entities.Account, error)
	GetAccountsByUserID(userID string) ([]*entities.Account, error)
	GetAllAccounts(page, pageSize int) ([]*entities.Account, int64, error)
	UpdateAccount(accountID string, req *dto.UpdateAccountRequest) (*entities.Account, error)
	DeleteAccount(accountID string) error

	// Balance operations
	GetAccountBalance(accountID string) (float64, error)
	UpdateAccountBalance(accountID string, amount float64) (float64, error)

	// Status operations
	UpdateAccountStatus(accountID string, isActive bool) (*entities.Account, error)
}

// InstallmentServiceInterface defines the contract for installment service operations
type InstallmentServiceInterface interface {
	// Preview and calculation
	CalculateInstallmentPlan(amount float64, installments int, interestRate float64) (*entities.InstallmentPlan, []*entities.Installment, error)

	// Plan management
	CreateInstallmentPlan(cardID string, amount float64, installments int, interestRate float64, description string) (*entities.InstallmentPlan, error)
	GetInstallmentPlansByCard(cardID string, page, pageSize int) ([]*entities.InstallmentPlan, int64, error)
	GetInstallmentPlanByID(planID string) (*entities.InstallmentPlan, error)
	CancelInstallmentPlan(planID string, reason string) error

	// Installment operations
	PayInstallment(installmentID string, amount float64) (*entities.Installment, error)
	GetInstallmentsByStatus(status string, page, pageSize int) ([]*entities.Installment, int64, error)
	GetOverdueInstallments(page, pageSize int) ([]*entities.Installment, int64, error)
	GetUpcomingInstallments(days int, page, pageSize int) ([]*entities.Installment, int64, error)

	// Summary and reporting
	GetInstallmentsSummary(userID string) (map[string]interface{}, error)
	GetMonthlyInstallmentsLoad(userID string, year, month int) ([]map[string]interface{}, error)
}
