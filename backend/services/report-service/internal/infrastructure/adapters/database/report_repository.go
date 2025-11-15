package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/fintrack/report-service/internal/core/domain/dto"
)

// ReportRepository implementación del repositorio de reportes
type ReportRepository struct {
	db *sql.DB
}

// NewReportRepository crea una nueva instancia del repositorio
func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

// GetTransactionReport obtiene el reporte de transacciones
func (r *ReportRepository) GetTransactionReport(ctx context.Context, userID string, startDate, endDate time.Time, txType string) (*dto.TransactionReportResponse, error) {
	response := &dto.TransactionReportResponse{
		UserID: userID,
		Period: dto.Period{
			StartDate: startDate,
			EndDate:   endDate,
			Days:      int(endDate.Sub(startDate).Hours() / 24),
		},
	}

	// Query para resumen general
	summaryQuery := `
		SELECT 
			COUNT(*) as total_transactions,
			COALESCE(SUM(CASE 
				WHEN type IN ('wallet_deposit', 'account_deposit', 'credit_payment', 'debit_refund', 'credit_refund') 
				THEN amount ELSE 0 END), 0) as total_income,
			COALESCE(SUM(CASE 
				WHEN type IN ('wallet_withdrawal', 'credit_charge', 'debit_purchase', 'account_withdraw', 'wallet_transfer') 
				THEN amount ELSE 0 END), 0) as total_expenses,
			COALESCE(AVG(amount), 0) as avg_transaction
		FROM transactions
		WHERE user_id = ? 
			AND created_at BETWEEN ? AND ?
			AND status = 'completed'
	`

	var summary dto.TransactionSummary
	err := r.db.QueryRowContext(ctx, summaryQuery, userID, startDate, endDate).Scan(
		&summary.TotalTransactions,
		&summary.TotalIncome,
		&summary.TotalExpenses,
		&summary.AvgTransaction,
	)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo resumen de transacciones: %w", err)
	}

	summary.NetBalance = summary.TotalIncome - summary.TotalExpenses
	response.Summary = summary

	// Query para transacciones por tipo
	byTypeQuery := `
		SELECT 
			type,
			COUNT(*) as count,
			COALESCE(SUM(amount), 0) as amount
		FROM transactions
		WHERE user_id = ? 
			AND created_at BETWEEN ? AND ?
			AND status = 'completed'
		GROUP BY type
		ORDER BY amount DESC
	`

	rows, err := r.db.QueryContext(ctx, byTypeQuery, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo transacciones por tipo: %w", err)
	}
	defer rows.Close()

	var byType []dto.TransactionByType
	totalAmount := summary.TotalIncome + summary.TotalExpenses

	for rows.Next() {
		var item dto.TransactionByType
		if err := rows.Scan(&item.Type, &item.Count, &item.Amount); err != nil {
			return nil, fmt.Errorf("error escaneando transacción por tipo: %w", err)
		}
		if totalAmount > 0 {
			item.Percentage = (item.Amount / totalAmount) * 100
		}
		byType = append(byType, item)
	}
	response.ByType = byType

	// Query para transacciones por período (agrupadas por día)
	byPeriodQuery := `
		SELECT 
			DATE(created_at) as date,
			COALESCE(SUM(CASE 
				WHEN type IN ('wallet_deposit', 'account_deposit', 'credit_payment', 'debit_refund', 'credit_refund') 
				THEN amount ELSE 0 END), 0) as income,
			COALESCE(SUM(CASE 
				WHEN type IN ('wallet_withdrawal', 'credit_charge', 'debit_purchase', 'account_withdraw', 'wallet_transfer') 
				THEN amount ELSE 0 END), 0) as expenses,
			COUNT(*) as count
		FROM transactions
		WHERE user_id = ? 
			AND created_at BETWEEN ? AND ?
			AND status = 'completed'
		GROUP BY DATE(created_at)
		ORDER BY date ASC
	`

	rows, err = r.db.QueryContext(ctx, byPeriodQuery, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo transacciones por período: %w", err)
	}
	defer rows.Close()

	var byPeriod []dto.TransactionByPeriod
	for rows.Next() {
		var item dto.TransactionByPeriod
		if err := rows.Scan(&item.Date, &item.Income, &item.Expenses, &item.Count); err != nil {
			return nil, fmt.Errorf("error escaneando transacción por período: %w", err)
		}
		item.Period = item.Date.Format("2006-01-02")
		item.Net = item.Income - item.Expenses
		byPeriod = append(byPeriod, item)
	}
	response.ByPeriod = byPeriod

	// Query para top gastos
	topExpensesQuery := `
		SELECT 
			id, description, amount, type, created_at, merchant_name
		FROM transactions
		WHERE user_id = ? 
			AND created_at BETWEEN ? AND ?
			AND status = 'completed'
			AND type IN ('wallet_withdrawal', 'credit_charge', 'debit_purchase', 'account_withdraw')
		ORDER BY amount DESC
		LIMIT 10
	`

	rows, err = r.db.QueryContext(ctx, topExpensesQuery, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo top gastos: %w", err)
	}
	defer rows.Close()

	var topExpenses []dto.TransactionItem
	for rows.Next() {
		var item dto.TransactionItem
		var merchantName sql.NullString
		if err := rows.Scan(&item.ID, &item.Description, &item.Amount, &item.Type, &item.Date, &merchantName); err != nil {
			return nil, fmt.Errorf("error escaneando top gasto: %w", err)
		}
		if merchantName.Valid {
			item.MerchantName = merchantName.String
		}
		topExpenses = append(topExpenses, item)
	}
	response.TopExpenses = topExpenses

	return response, nil
}

// GetInstallmentReport obtiene el reporte de cuotas
func (r *ReportRepository) GetInstallmentReport(ctx context.Context, userID string, status string) (*dto.InstallmentReportResponse, error) {
	response := &dto.InstallmentReportResponse{
		UserID: userID,
	}

	// Query para resumen
	summaryQuery := `
		SELECT 
			COUNT(*) as total_plans,
			COALESCE(SUM(CASE WHEN status = 'active' THEN 1 ELSE 0 END), 0) as active_plans,
			COALESCE(SUM(total_amount), 0) as total_amount,
			COALESCE(SUM(total_amount - remaining_amount), 0) as paid_amount,
			COALESCE(SUM(remaining_amount), 0) as remaining_amount
		FROM installment_plans
		WHERE user_id = ?
	`

	var summary dto.InstallmentSummary
	var nextPaymentDate sql.NullTime

	err := r.db.QueryRowContext(ctx, summaryQuery, userID).Scan(
		&summary.TotalPlans,
		&summary.ActivePlans,
		&summary.TotalAmount,
		&summary.PaidAmount,
		&summary.RemainingAmount,
	)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo resumen de cuotas: %w", err)
	}

	// Calcular monto vencido
	overdueQuery := `
		SELECT COALESCE(SUM(i.remaining_amount), 0)
		FROM installments i
		JOIN installment_plans ip ON BINARY i.plan_id = BINARY ip.id
		WHERE BINARY ip.user_id = BINARY ? AND i.status = 'overdue'
	`
	err = r.db.QueryRowContext(ctx, overdueQuery, userID).Scan(&summary.OverdueAmount)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo monto vencido: %w", err)
	}

	// Próximo pago
	nextPaymentQuery := `
		SELECT amount, due_date
		FROM installments i
		JOIN installment_plans ip ON BINARY i.plan_id = BINARY ip.id
		WHERE BINARY ip.user_id = BINARY ? AND i.status = 'pending'
		ORDER BY i.due_date ASC
		LIMIT 1
	`
	err = r.db.QueryRowContext(ctx, nextPaymentQuery, userID).Scan(&summary.NextPaymentAmount, &nextPaymentDate)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error obteniendo próximo pago: %w", err)
	}
	if nextPaymentDate.Valid {
		summary.NextPaymentDate = &nextPaymentDate.Time
	}

	if summary.TotalAmount > 0 {
		summary.CompletionPercentage = (summary.PaidAmount / summary.TotalAmount) * 100
	}

	response.Summary = summary

	// Query para planes de cuotas
	plansQuery := `
		SELECT 
			ip.id, ip.card_id, c.last_four_digits, ip.total_amount, 
			ip.installments_count, ip.installment_amount, ip.paid_installments, 
			ip.remaining_amount, ip.status, ip.description, ip.merchant_name, 
			ip.start_date
		FROM installment_plans ip
		LEFT JOIN cards c ON BINARY ip.card_id = BINARY c.id
		WHERE BINARY ip.user_id = BINARY ?
		ORDER BY ip.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, plansQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo planes de cuotas: %w", err)
	}
	defer rows.Close()

	var plans []dto.InstallmentPlan
	for rows.Next() {
		var plan dto.InstallmentPlan
		var description, merchantName, lastFour sql.NullString

		err := rows.Scan(
			&plan.ID, &plan.CardID, &lastFour, &plan.TotalAmount,
			&plan.InstallmentsCount, &plan.InstallmentAmount, &plan.PaidInstallments,
			&plan.RemainingAmount, &plan.Status, &description, &merchantName,
			&plan.StartDate,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando plan de cuotas: %w", err)
		}

		if lastFour.Valid {
			plan.CardLastFour = lastFour.String
		}
		if description.Valid {
			plan.Description = description.String
		}
		if merchantName.Valid {
			plan.MerchantName = merchantName.String
		}

		if plan.TotalAmount > 0 {
			plan.CompletionPercentage = ((plan.TotalAmount - plan.RemainingAmount) / plan.TotalAmount) * 100
		}

		// Obtener próxima fecha de vencimiento
		var nextDue sql.NullTime
		nextDueQuery := `SELECT due_date FROM installments WHERE BINARY plan_id = BINARY ? AND status = 'pending' ORDER BY due_date ASC LIMIT 1`
		err = r.db.QueryRowContext(ctx, nextDueQuery, plan.ID).Scan(&nextDue)
		if err == nil && nextDue.Valid {
			plan.NextDueDate = &nextDue.Time
		}

		plans = append(plans, plan)
	}
	response.Plans = plans

	// Pagos próximos (próximos 30 días)
	upcomingQuery := `
		SELECT 
			i.id, i.plan_id, c.last_four_digits, i.amount, i.due_date,
			DATEDIFF(i.due_date, CURDATE()) as days_until,
			ip.description, ip.merchant_name
		FROM installments i
		JOIN installment_plans ip ON BINARY i.plan_id = BINARY ip.id
		LEFT JOIN cards c ON BINARY ip.card_id = BINARY c.id
		WHERE BINARY ip.user_id = BINARY ? 
			AND i.status = 'pending'
			AND i.due_date BETWEEN CURDATE() AND DATE_ADD(CURDATE(), INTERVAL 30 DAY)
		ORDER BY i.due_date ASC
		LIMIT 10
	`

	upcomingRows, err := r.db.QueryContext(ctx, upcomingQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo pagos próximos: %w", err)
	}
	defer upcomingRows.Close()

	var upcoming []dto.UpcomingPayment
	for upcomingRows.Next() {
		var payment dto.UpcomingPayment
		var lastFour, description, merchantName sql.NullString

		err := upcomingRows.Scan(
			&payment.InstallmentID, &payment.PlanID, &lastFour, &payment.Amount,
			&payment.DueDate, &payment.DaysUntilDue, &description, &merchantName,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando pago próximo: %w", err)
		}

		if lastFour.Valid {
			payment.CardLastFour = lastFour.String
		}
		if description.Valid {
			payment.Description = description.String
		}
		if merchantName.Valid {
			payment.MerchantName = merchantName.String
		}

		upcoming = append(upcoming, payment)
	}
	response.Upcoming = upcoming

	// Pagos vencidos
	overduePaymentsQuery := `
		SELECT 
			i.id, i.plan_id, c.last_four_digits, i.amount, i.due_date,
			DATEDIFF(CURDATE(), i.due_date) as days_overdue,
			i.late_fee, ip.description, ip.merchant_name
		FROM installments i
		JOIN installment_plans ip ON BINARY i.plan_id = BINARY ip.id
		LEFT JOIN cards c ON BINARY ip.card_id = BINARY c.id
		WHERE BINARY ip.user_id = BINARY ? AND i.status = 'overdue'
		ORDER BY i.due_date ASC
	`

	overdueRows, err := r.db.QueryContext(ctx, overduePaymentsQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo pagos vencidos: %w", err)
	}
	defer overdueRows.Close()

	var overdue []dto.OverduePayment
	for overdueRows.Next() {
		var payment dto.OverduePayment
		var lastFour, description, merchantName sql.NullString

		err := overdueRows.Scan(
			&payment.InstallmentID, &payment.PlanID, &lastFour, &payment.Amount,
			&payment.DueDate, &payment.DaysOverdue, &payment.LateFee,
			&description, &merchantName,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando pago vencido: %w", err)
		}

		if lastFour.Valid {
			payment.CardLastFour = lastFour.String
		}
		if description.Valid {
			payment.Description = description.String
		}
		if merchantName.Valid {
			payment.MerchantName = merchantName.String
		}

		overdue = append(overdue, payment)
	}
	response.Overdue = overdue

	return response, nil
}
