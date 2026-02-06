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
	
	// Calcular el total de transacciones (suma de montos absolutos)
	var totalTransactionAmount float64

	for rows.Next() {
		var item dto.TransactionByType
		if err := rows.Scan(&item.Type, &item.Count, &item.Amount); err != nil {
			return nil, fmt.Errorf("error escaneando transacción por tipo: %w", err)
		}
		byType = append(byType, item)
		totalTransactionAmount += item.Amount
	}

	// Calcular porcentajes sobre el total de transacciones
	if totalTransactionAmount > 0 {
		for i := range byType {
			byType[i].Percentage = (byType[i].Amount / totalTransactionAmount) * 100
		}
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

// GetInstallmentReport obtiene el reporte de cuotas (startDate/endDate opcionales: filtran por due_date)
func (r *ReportRepository) GetInstallmentReport(ctx context.Context, userID string, status string, startDate, endDate time.Time) (*dto.InstallmentReportResponse, error) {
	response := &dto.InstallmentReportResponse{
		UserID: userID,
	}

	hasDateFilter := !startDate.IsZero() || !endDate.IsZero()
	var planFilter string
	var dueDateCond string
	installArgs := []interface{}{userID}
	if hasDateFilter {
		planFilter = ` AND ip.id IN (SELECT BINARY plan_id FROM installments WHERE 1=1 `
		if !startDate.IsZero() {
			planFilter += ` AND due_date >= ?`
			dueDateCond += ` AND i.due_date >= ?`
			installArgs = append(installArgs, startDate.Format("2006-01-02"))
		}
		if !endDate.IsZero() {
			planFilter += ` AND due_date <= ?`
			dueDateCond += ` AND i.due_date <= ?`
			installArgs = append(installArgs, endDate.Format("2006-01-02"))
		}
		planFilter += `)`
	}

	// Query para resumen (solo planes que tengan cuotas en el rango si hay filtro)
	summaryQuery := `
		SELECT 
			COUNT(*) as total_plans,
			COALESCE(SUM(CASE WHEN ip.status = 'active' THEN 1 ELSE 0 END), 0) as active_plans,
			COALESCE(SUM(ip.total_amount), 0) as total_amount,
			COALESCE(SUM(ip.total_amount - ip.remaining_amount), 0) as paid_amount,
			COALESCE(SUM(ip.remaining_amount), 0) as remaining_amount
		FROM installment_plans ip
		WHERE BINARY ip.user_id = BINARY ?` + planFilter

	var summary dto.InstallmentSummary
	var nextPaymentDate sql.NullTime

	var err error
	if hasDateFilter {
		err = r.db.QueryRowContext(ctx, summaryQuery, installArgs...).Scan(
			&summary.TotalPlans,
			&summary.ActivePlans,
			&summary.TotalAmount,
			&summary.PaidAmount,
			&summary.RemainingAmount,
		)
	} else {
		err = r.db.QueryRowContext(ctx, summaryQuery, userID).Scan(
			&summary.TotalPlans,
			&summary.ActivePlans,
			&summary.TotalAmount,
			&summary.PaidAmount,
			&summary.RemainingAmount,
		)
	}
	if err != nil {
		return nil, fmt.Errorf("error obteniendo resumen de cuotas: %w", err)
	}

	// Calcular monto vencido (con filtro de fechas si aplica)
	overdueQuery := `
		SELECT COALESCE(SUM(i.remaining_amount), 0)
		FROM installments i
		JOIN installment_plans ip ON BINARY i.plan_id = BINARY ip.id
		WHERE BINARY ip.user_id = BINARY ? AND i.status = 'overdue'` + dueDateCond + `
	`
	err = r.db.QueryRowContext(ctx, overdueQuery, installArgs...).Scan(&summary.OverdueAmount)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo monto vencido: %w", err)
	}

	// Próximo pago (con filtro de fechas si aplica)
	nextPaymentQuery := `
		SELECT amount, due_date
		FROM installments i
		JOIN installment_plans ip ON BINARY i.plan_id = BINARY ip.id
		WHERE BINARY ip.user_id = BINARY ? AND i.status = 'pending'` + dueDateCond + `
		ORDER BY i.due_date ASC
		LIMIT 1
	`
	err = r.db.QueryRowContext(ctx, nextPaymentQuery, installArgs...).Scan(&summary.NextPaymentAmount, &nextPaymentDate)
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

	// Query para planes de cuotas (con filtro por rango de due_date si aplica)
	plansQuery := `
		SELECT 
			ip.id, ip.card_id, c.last_four_digits, ip.total_amount, 
			ip.installments_count, ip.installment_amount, ip.paid_installments, 
			ip.remaining_amount, ip.status, ip.description, ip.merchant_name, 
			ip.start_date
		FROM installment_plans ip
		LEFT JOIN cards c ON BINARY ip.card_id = BINARY c.id
		WHERE BINARY ip.user_id = BINARY ?` + planFilter + `
		ORDER BY ip.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, plansQuery, installArgs...)
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

	// Pagos próximos (próximos 30 días, o en rango si hay filtro de fechas)
	var upcomingQuery string
	if hasDateFilter {
		upcomingQuery = `
		SELECT 
			i.id, i.plan_id, c.last_four_digits, i.amount, i.due_date,
			DATEDIFF(i.due_date, CURDATE()) as days_until,
			ip.description, ip.merchant_name
		FROM installments i
		JOIN installment_plans ip ON BINARY i.plan_id = BINARY ip.id
		LEFT JOIN cards c ON BINARY ip.card_id = BINARY c.id
		WHERE BINARY ip.user_id = BINARY ? 
			AND i.status = 'pending'` + dueDateCond + `
		ORDER BY i.due_date ASC
		LIMIT 10
		`
	} else {
		upcomingQuery = `
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
	}
	var upcomingRows *sql.Rows
	if hasDateFilter {
		upcomingRows, err = r.db.QueryContext(ctx, upcomingQuery, installArgs...)
	} else {
		upcomingRows, err = r.db.QueryContext(ctx, upcomingQuery, userID)
	}
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

	// Pagos vencidos (con filtro de fechas si aplica)
	overduePaymentsQuery := `
		SELECT 
			i.id, i.plan_id, c.last_four_digits, i.amount, i.due_date,
			DATEDIFF(CURDATE(), i.due_date) as days_overdue,
			i.late_fee, ip.description, ip.merchant_name
		FROM installments i
		JOIN installment_plans ip ON BINARY i.plan_id = BINARY ip.id
		LEFT JOIN cards c ON BINARY ip.card_id = BINARY c.id
		WHERE BINARY ip.user_id = BINARY ? AND i.status = 'overdue'` + dueDateCond + `
		ORDER BY i.due_date ASC
	`

	var overdueRows *sql.Rows
	if hasDateFilter {
		overdueRows, err = r.db.QueryContext(ctx, overduePaymentsQuery, installArgs...)
	} else {
		overdueRows, err = r.db.QueryContext(ctx, overduePaymentsQuery, userID)
	}
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
