package interfaces

import "github.com/fintrack/transaction-service/internal/infrastructure/http/clients"

// AccountServiceInterface define los métodos para comunicarse con el account-service
type AccountServiceInterface interface {
	GetAccountBalance(accountID string) (*clients.AccountBalance, error)
	GetAccountInfo(accountID string) (*clients.AccountInfo, error)
	AddFunds(accountID string, amount float64, description string, reference string) (*clients.BalanceUpdateResponse, error)
	WithdrawFunds(accountID string, amount float64, description string, reference string) (*clients.BalanceUpdateResponse, error)
	UpdateCreditUsage(accountID string, amount float64, description string, reference string) (*clients.BalanceUpdateResponse, error)
	GetAvailableCredit(accountID string) (*clients.AccountBalance, error)
	ValidateAccountExists(accountID string) (bool, error)
	HealthCheck() error
}
