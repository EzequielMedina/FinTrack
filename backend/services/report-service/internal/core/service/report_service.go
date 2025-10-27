package service

import (
	"context"
	"time"

	"github.com/fintrack/report-service/internal/core/domain/dto"
)

// ReportService interfaz del servicio de reportes
type ReportService interface {
	// Reportes de transacciones
	GetTransactionReport(ctx context.Context, req *dto.TransactionReportRequest) (*dto.TransactionReportResponse, error)

	// Reportes de cuotas
	GetInstallmentReport(ctx context.Context, req *dto.InstallmentReportRequest) (*dto.InstallmentReportResponse, error)

	// Reportes de cuentas
	GetAccountReport(ctx context.Context, req *dto.AccountReportRequest) (*dto.AccountReportResponse, error)

	// Reportes de gastos vs ingresos
	GetExpenseIncomeReport(ctx context.Context, req *dto.ExpenseIncomeReportRequest) (*dto.ExpenseIncomeReportResponse, error)

	// Reportes de notificaciones
	GetNotificationReport(ctx context.Context, req *dto.NotificationReportRequest) (*dto.NotificationReportResponse, error)
}

// reportService implementación del servicio
type reportService struct {
	repo ReportRepository
}

// ReportRepository interfaz del repositorio
type ReportRepository interface {
	GetTransactionReport(ctx context.Context, userID string, startDate, endDate time.Time, txType string) (*dto.TransactionReportResponse, error)
	GetInstallmentReport(ctx context.Context, userID string, status string) (*dto.InstallmentReportResponse, error)
	GetAccountReport(ctx context.Context, userID string) (*dto.AccountReportResponse, error)
	GetExpenseIncomeReport(ctx context.Context, userID string, startDate, endDate time.Time) (*dto.ExpenseIncomeReportResponse, error)
	GetNotificationReport(ctx context.Context, startDate, endDate time.Time) (*dto.NotificationReportResponse, error)
}

// NewReportService crea una nueva instancia del servicio
func NewReportService(repo ReportRepository) ReportService {
	return &reportService{
		repo: repo,
	}
}

// GetTransactionReport obtiene el reporte de transacciones
func (s *reportService) GetTransactionReport(ctx context.Context, req *dto.TransactionReportRequest) (*dto.TransactionReportResponse, error) {
	// Si no se especifican fechas, usar el mes actual
	startDate := req.StartDate
	endDate := req.EndDate

	if startDate.IsZero() {
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	}

	if endDate.IsZero() {
		endDate = time.Now()
	}

	return s.repo.GetTransactionReport(ctx, req.UserID, startDate, endDate, req.Type)
}

// GetInstallmentReport obtiene el reporte de cuotas
func (s *reportService) GetInstallmentReport(ctx context.Context, req *dto.InstallmentReportRequest) (*dto.InstallmentReportResponse, error) {
	return s.repo.GetInstallmentReport(ctx, req.UserID, req.Status)
}

// GetAccountReport obtiene el reporte de cuentas
func (s *reportService) GetAccountReport(ctx context.Context, req *dto.AccountReportRequest) (*dto.AccountReportResponse, error) {
	return s.repo.GetAccountReport(ctx, req.UserID)
}

// GetExpenseIncomeReport obtiene el reporte de gastos vs ingresos
func (s *reportService) GetExpenseIncomeReport(ctx context.Context, req *dto.ExpenseIncomeReportRequest) (*dto.ExpenseIncomeReportResponse, error) {
	return s.repo.GetExpenseIncomeReport(ctx, req.UserID, req.StartDate, req.EndDate)
}

// GetNotificationReport obtiene el reporte de notificaciones
func (s *reportService) GetNotificationReport(ctx context.Context, req *dto.NotificationReportRequest) (*dto.NotificationReportResponse, error) {
	// Si no se especifican fechas, usar los últimos 30 días
	startDate := req.StartDate
	endDate := req.EndDate

	if startDate.IsZero() {
		startDate = time.Now().AddDate(0, 0, -30)
	}

	if endDate.IsZero() {
		endDate = time.Now()
	}

	return s.repo.GetNotificationReport(ctx, startDate, endDate)
}
