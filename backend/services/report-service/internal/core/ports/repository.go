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

	// Reportes de cuotas (startDate/endDate opcionales: filtran por due_date de cuotas)
	GetInstallmentReport(ctx context.Context, userID string, status string, startDate, endDate time.Time) (*dto.InstallmentReportResponse, error)

	// Reportes de cuentas (startDate/endDate opcionales: filtran por created_at)
	GetAccountReport(ctx context.Context, userID string, startDate, endDate time.Time) (*dto.AccountReportResponse, error)

	// Reportes de gastos vs ingresos
	GetExpenseIncomeReport(ctx context.Context, userID string, startDate, endDate time.Time) (*dto.ExpenseIncomeReportResponse, error)

	// Reportes de notificaciones
	GetNotificationReport(ctx context.Context, startDate, endDate time.Time) (*dto.NotificationReportResponse, error)
}
