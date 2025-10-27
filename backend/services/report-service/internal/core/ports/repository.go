package ports

import (
	"context"
	"time"

	"github.com/fintrack/report-service/internal/core/domain/dto"
)

// ReportRepository interfaz del repositorio de reportes
type ReportRepository interface {
	// Reportes de transacciones
	GetTransactionReport(ctx context.Context, userID string, startDate, endDate time.Time, txType string) (*dto.TransactionReportResponse, error)

	// Reportes de cuotas
	GetInstallmentReport(ctx context.Context, userID string, status string) (*dto.InstallmentReportResponse, error)

	// Reportes de cuentas
	GetAccountReport(ctx context.Context, userID string) (*dto.AccountReportResponse, error)

	// Reportes de gastos vs ingresos
	GetExpenseIncomeReport(ctx context.Context, userID string, startDate, endDate time.Time) (*dto.ExpenseIncomeReportResponse, error)

	// Reportes de notificaciones
	GetNotificationReport(ctx context.Context, startDate, endDate time.Time) (*dto.NotificationReportResponse, error)
}
