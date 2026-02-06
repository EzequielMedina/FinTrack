package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/fintrack/report-service/internal/core/domain/dto"
)

// GetAccountReport obtiene el reporte de cuentas (startDate/endDate opcionales: filtran por created_at)
func (r *ReportRepository) GetAccountReport(ctx context.Context, userID string, startDate, endDate time.Time) (*dto.AccountReportResponse, error) {
	response := &dto.AccountReportResponse{
		UserID: userID,
	}

	hasDateFilter := !startDate.IsZero() || !endDate.IsZero()
	dateFilter := ""
	accountArgs := []interface{}{userID}
	if hasDateFilter {
		if !startDate.IsZero() {
			dateFilter += " AND created_at >= ?"
			accountArgs = append(accountArgs, startDate.Format("2006-01-02"))
		}
		if !endDate.IsZero() {
			dateFilter += " AND DATE(created_at) <= ?"
			accountArgs = append(accountArgs, endDate.Format("2006-01-02"))
		}
	}

	// Query para resumen
	summaryQuery := `
		SELECT 
			COALESCE(SUM(balance), 0) as total_balance,
			COUNT(*) as total_accounts,
			COALESCE(SUM(credit_limit), 0) as total_credit_limit
		FROM accounts
		WHERE BINARY user_id = BINARY ? AND is_active = 1 AND deleted_at IS NULL` + dateFilter + `
	`

	var summary dto.AccountSummary
	var err error
	if hasDateFilter {
		err = r.db.QueryRowContext(ctx, summaryQuery, accountArgs...).Scan(
			&summary.TotalBalance,
			&summary.TotalAccounts,
			&summary.TotalCreditLimit,
		)
	} else {
		err = r.db.QueryRowContext(ctx, summaryQuery, userID).Scan(
			&summary.TotalBalance,
			&summary.TotalAccounts,
			&summary.TotalCreditLimit,
		)
	}
	if err != nil {
		return nil, fmt.Errorf("error obteniendo resumen de cuentas: %w", err)
	}

	// Contar tarjetas (con filtro por created_at si aplica)
	cardDateFilter := ""
	if hasDateFilter {
		if !startDate.IsZero() {
			cardDateFilter += " AND c.created_at >= ?"
		}
		if !endDate.IsZero() {
			cardDateFilter += " AND DATE(c.created_at) <= ?"
		}
	}
	cardsCountQuery := `
		SELECT COUNT(*)
		FROM cards c
		JOIN accounts a ON BINARY c.account_id = BINARY a.id
		WHERE BINARY a.user_id = BINARY ? AND c.status = 'active' AND c.deleted_at IS NULL` + cardDateFilter + `
	`
	if hasDateFilter {
		err = r.db.QueryRowContext(ctx, cardsCountQuery, accountArgs...).Scan(&summary.TotalCards)
	} else {
		err = r.db.QueryRowContext(ctx, cardsCountQuery, userID).Scan(&summary.TotalCards)
	}
	if err != nil {
		return nil, fmt.Errorf("error contando tarjetas: %w", err)
	}

	// Calcular crédito usado (no filtramos por fecha en transacciones para el resumen)
	creditUsedQuery := `
		SELECT COALESCE(SUM(t.amount), 0)
		FROM transactions t
		JOIN cards c ON BINARY t.from_card_id = BINARY c.id
		JOIN accounts a ON BINARY c.account_id = BINARY a.id
		WHERE BINARY a.user_id = BINARY ? 
			AND a.account_type = 'credit'
			AND c.card_type = 'credit'
			AND t.status IN ('pending', 'completed')
			AND t.type = 'credit_charge'
	`
	err = r.db.QueryRowContext(ctx, creditUsedQuery, userID).Scan(&summary.TotalCreditUsed)
	if err != nil {
		return nil, fmt.Errorf("error calculando crédito usado: %w", err)
	}

	summary.AvailableCredit = summary.TotalCreditLimit - summary.TotalCreditUsed
	if summary.TotalCreditLimit > 0 {
		summary.CreditUtilization = (summary.TotalCreditUsed / summary.TotalCreditLimit) * 100
	}
	summary.NetWorth = summary.TotalBalance - summary.TotalCreditUsed

	response.Summary = summary

	// Query para detalle de cuentas (con filtro por created_at si aplica)
	accountsQuery := `
		SELECT 
			id, account_type, name, currency, balance, 
			COALESCE(credit_limit, 0) as credit_limit, is_active
		FROM accounts
		WHERE BINARY user_id = BINARY ? AND deleted_at IS NULL` + dateFilter + `
		ORDER BY created_at DESC
	`
	var rows *sql.Rows
	if hasDateFilter {
		rows, err = r.db.QueryContext(ctx, accountsQuery, accountArgs...)
	} else {
		rows, err = r.db.QueryContext(ctx, accountsQuery, userID)
	}
	if err != nil {
		return nil, fmt.Errorf("error obteniendo cuentas: %w", err)
	}
	defer rows.Close()

	var accounts []dto.AccountDetail
	for rows.Next() {
		var account dto.AccountDetail
		err := rows.Scan(
			&account.ID, &account.AccountType, &account.Name, &account.Currency,
			&account.Balance, &account.CreditLimit, &account.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando cuenta: %w", err)
		}
		accounts = append(accounts, account)
	}
	response.Accounts = accounts

	// Query para detalle de tarjetas (con filtro por created_at si aplica)
	// Ahora incluimos el campo balance directamente en la consulta
	cardsQuery := `
		SELECT 
			c.id, c.account_id, c.card_type, c.card_brand, c.last_four_digits,
			c.holder_name, c.status, COALESCE(c.credit_limit, 0) as credit_limit,
			COALESCE(c.balance, 0) as balance,
			COALESCE(c.nickname, '') as nickname
		FROM cards c
		JOIN accounts a ON BINARY c.account_id = BINARY a.id
		WHERE BINARY a.user_id = BINARY ? AND c.deleted_at IS NULL` + cardDateFilter + `
		ORDER BY c.created_at DESC
	`
	var cardRows *sql.Rows
	if hasDateFilter {
		cardRows, err = r.db.QueryContext(ctx, cardsQuery, accountArgs...)
	} else {
		cardRows, err = r.db.QueryContext(ctx, cardsQuery, userID)
	}
	if err != nil {
		return nil, fmt.Errorf("error obteniendo tarjetas: %w", err)
	}
	defer cardRows.Close()

	var cards []dto.CardDetail
	for cardRows.Next() {
		var card dto.CardDetail
		err := cardRows.Scan(
			&card.ID, &card.AccountID, &card.CardType, &card.CardBrand,
			&card.LastFourDigits, &card.HolderName, &card.Status,
			&card.CreditLimit, &card.CurrentBalance, &card.Nickname,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando tarjeta: %w", err)
		}

		// Calcular crédito disponible para tarjetas de crédito
		if card.CardType == "credit" && card.CreditLimit > 0 {
			card.AvailableCredit = card.CreditLimit - card.CurrentBalance
			// Asegurar que no sea negativo
			if card.AvailableCredit < 0 {
				card.AvailableCredit = 0
			}
		}

		cards = append(cards, card)
	}
	response.Cards = cards

	// Query para distribución de cuentas (con filtro por created_at si aplica)
	distributionQuery := `
		SELECT 
			account_type,
			COUNT(*) as count,
			COALESCE(SUM(balance), 0) as total_balance
		FROM accounts
		WHERE BINARY user_id = BINARY ? AND is_active = 1 AND deleted_at IS NULL` + dateFilter + `
		GROUP BY account_type
	`
	var distRows *sql.Rows
	if hasDateFilter {
		distRows, err = r.db.QueryContext(ctx, distributionQuery, accountArgs...)
	} else {
		distRows, err = r.db.QueryContext(ctx, distributionQuery, userID)
	}
	if err != nil {
		return nil, fmt.Errorf("error obteniendo distribución: %w", err)
	}
	defer distRows.Close()

	var distribution []dto.AccountDistribution
	for distRows.Next() {
		var dist dto.AccountDistribution
		err := distRows.Scan(&dist.AccountType, &dist.Count, &dist.TotalBalance)
		if err != nil {
			return nil, fmt.Errorf("error escaneando distribución: %w", err)
		}
		if summary.TotalBalance > 0 {
			dist.Percentage = (dist.TotalBalance / summary.TotalBalance) * 100
		}
		distribution = append(distribution, dist)
	}
	response.Distribution = distribution

	return response, nil
}

// GetExpenseIncomeReport obtiene el reporte de gastos vs ingresos
func (r *ReportRepository) GetExpenseIncomeReport(ctx context.Context, userID string, startDate, endDate time.Time) (*dto.ExpenseIncomeReportResponse, error) {
	response := &dto.ExpenseIncomeReportResponse{
		UserID: userID,
		Period: dto.Period{
			StartDate: startDate,
			EndDate:   endDate,
			Days:      int(endDate.Sub(startDate).Hours() / 24),
		},
	}

	// Query para resumen
	summaryQuery := `
		SELECT 
			COALESCE(SUM(CASE 
				WHEN type IN ('wallet_deposit', 'account_deposit', 'credit_payment', 'debit_refund', 'credit_refund') 
				THEN amount ELSE 0 END), 0) as total_income,
			COALESCE(SUM(CASE 
				WHEN type IN ('wallet_withdrawal', 'credit_charge', 'debit_purchase', 'account_withdraw', 'wallet_transfer') 
				THEN amount ELSE 0 END), 0) as total_expenses
		FROM transactions
		WHERE user_id = ? 
			AND created_at BETWEEN ? AND ?
			AND status = 'completed'
	`

	var summary dto.ExpenseIncomeSummary
	err := r.db.QueryRowContext(ctx, summaryQuery, userID, startDate, endDate).Scan(
		&summary.TotalIncome,
		&summary.TotalExpenses,
	)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo resumen de gastos/ingresos: %w", err)
	}

	summary.NetBalance = summary.TotalIncome - summary.TotalExpenses
	if summary.TotalIncome > 0 {
		summary.SavingsRate = (summary.NetBalance / summary.TotalIncome) * 100
		summary.ExpenseRatio = (summary.TotalExpenses / summary.TotalIncome) * 100
	}

	days := response.Period.Days
	if days > 0 {
		summary.AvgDailyIncome = summary.TotalIncome / float64(days)
		summary.AvgDailyExpense = summary.TotalExpenses / float64(days)
	}

	response.Summary = summary

	// Query para gastos/ingresos por período (por día)
	byPeriodQuery := `
		SELECT 
			DATE(created_at) as date,
			COALESCE(SUM(CASE 
				WHEN type IN ('wallet_deposit', 'account_deposit', 'credit_payment', 'debit_refund', 'credit_refund') 
				THEN amount ELSE 0 END), 0) as income,
			COALESCE(SUM(CASE 
				WHEN type IN ('wallet_withdrawal', 'credit_charge', 'debit_purchase', 'account_withdraw', 'wallet_transfer') 
				THEN amount ELSE 0 END), 0) as expenses
		FROM transactions
		WHERE user_id = ? 
			AND created_at BETWEEN ? AND ?
			AND status = 'completed'
		GROUP BY DATE(created_at)
		ORDER BY date ASC
	`

	rows, err := r.db.QueryContext(ctx, byPeriodQuery, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo gastos/ingresos por período: %w", err)
	}
	defer rows.Close()

	var byPeriod []dto.ExpenseIncomeByPeriod
	for rows.Next() {
		var item dto.ExpenseIncomeByPeriod
		err := rows.Scan(&item.Date, &item.Income, &item.Expenses)
		if err != nil {
			return nil, fmt.Errorf("error escaneando período: %w", err)
		}
		item.Period = item.Date.Format("2006-01-02")
		item.Net = item.Income - item.Expenses
		if item.Income > 0 {
			item.SavingsRate = (item.Net / item.Income) * 100
		}
		byPeriod = append(byPeriod, item)
	}
	response.ByPeriod = byPeriod

	// Query para gastos/ingresos por categoría (tipo de transacción)
	byCategoryQuery := `
		SELECT 
			type,
			CASE 
				WHEN type IN ('wallet_deposit', 'account_deposit', 'credit_payment', 'debit_refund', 'credit_refund') 
				THEN 'income' ELSE 'expense' 
			END as category_type,
			COUNT(*) as count,
			COALESCE(SUM(amount), 0) as amount
		FROM transactions
		WHERE user_id = ? 
			AND created_at BETWEEN ? AND ?
			AND status = 'completed'
		GROUP BY type
		ORDER BY amount DESC
	`

	rows, err = r.db.QueryContext(ctx, byCategoryQuery, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo gastos/ingresos por categoría: %w", err)
	}
	defer rows.Close()

	var byCategory []dto.ExpenseIncomeByCategory
	var totalCategoryAmount float64
	
	// Primera pasada: recolectar todos los items y calcular el total
	for rows.Next() {
		var item dto.ExpenseIncomeByCategory
		err := rows.Scan(&item.Category, &item.Type, &item.Count, &item.Amount)
		if err != nil {
			return nil, fmt.Errorf("error escaneando categoría: %w", err)
		}
		byCategory = append(byCategory, item)
		totalCategoryAmount += item.Amount
	}
	
	// Segunda pasada: calcular porcentajes basados en el total de todas las categorías
	if totalCategoryAmount > 0 {
		for i := range byCategory {
			byCategory[i].Percentage = (byCategory[i].Amount / totalCategoryAmount) * 100
		}
	}
	
	response.ByCategory = byCategory

	// Análisis de tendencias simple
	trend := dto.TrendAnalysis{
		IncomesTrend:  "stable",
		ExpensesTrend: "stable",
		NetTrend:      "stable",
	}

	// Comparar con período anterior
	prevStartDate := startDate.AddDate(0, 0, -days)
	prevEndDate := startDate.AddDate(0, 0, -1)

	var prevIncome, prevExpenses float64
	prevQuery := `
		SELECT 
			COALESCE(SUM(CASE 
				WHEN type IN ('wallet_deposit', 'account_deposit', 'credit_payment', 'debit_refund', 'credit_refund') 
				THEN amount ELSE 0 END), 0) as prev_income,
			COALESCE(SUM(CASE 
				WHEN type IN ('wallet_withdrawal', 'credit_charge', 'debit_purchase', 'account_withdraw', 'wallet_transfer') 
				THEN amount ELSE 0 END), 0) as prev_expenses
		FROM transactions
		WHERE user_id = ? 
			AND created_at BETWEEN ? AND ?
			AND status = 'completed'
	`

	err = r.db.QueryRowContext(ctx, prevQuery, userID, prevStartDate, prevEndDate).Scan(&prevIncome, &prevExpenses)
	if err == nil && prevIncome > 0 {
		trend.IncomeChange = ((summary.TotalIncome - prevIncome) / prevIncome) * 100
		if trend.IncomeChange > 5 {
			trend.IncomesTrend = "increasing"
		} else if trend.IncomeChange < -5 {
			trend.IncomesTrend = "decreasing"
		}
	}

	if err == nil && prevExpenses > 0 {
		trend.ExpenseChange = ((summary.TotalExpenses - prevExpenses) / prevExpenses) * 100
		if trend.ExpenseChange > 5 {
			trend.ExpensesTrend = "increasing"
		} else if trend.ExpenseChange < -5 {
			trend.ExpensesTrend = "decreasing"
		}
	}

	// Determinar tendencia neta
	if trend.IncomesTrend == "increasing" && trend.ExpensesTrend == "decreasing" {
		trend.NetTrend = "improving"
	} else if trend.IncomesTrend == "decreasing" && trend.ExpensesTrend == "increasing" {
		trend.NetTrend = "declining"
	}

	response.Trend = trend

	return response, nil
}

// GetNotificationReport obtiene el reporte de notificaciones
func (r *ReportRepository) GetNotificationReport(ctx context.Context, startDate, endDate time.Time) (*dto.NotificationReportResponse, error) {
	response := &dto.NotificationReportResponse{
		Period: dto.Period{
			StartDate: startDate,
			EndDate:   endDate,
			Days:      int(endDate.Sub(startDate).Hours() / 24),
		},
	}

	// Query para resumen
	summaryQuery := `
		SELECT 
			COUNT(*) as total_notifications,
			COALESCE(SUM(CASE WHEN status = 'sent' THEN 1 ELSE 0 END), 0) as successful,
			COALESCE(SUM(CASE WHEN status = 'failed' THEN 1 ELSE 0 END), 0) as failed
		FROM notification_logs
		WHERE created_at BETWEEN ? AND ?
	`

	var summary dto.NotificationSummary
	err := r.db.QueryRowContext(ctx, summaryQuery, startDate, endDate).Scan(
		&summary.TotalNotifications,
		&summary.SuccessfulSent,
		&summary.Failed,
	)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo resumen de notificaciones: %w", err)
	}

	// Contar job runs
	jobCountQuery := `SELECT COUNT(*) FROM job_runs WHERE started_at BETWEEN ? AND ?`
	err = r.db.QueryRowContext(ctx, jobCountQuery, startDate, endDate).Scan(&summary.TotalJobRuns)
	if err != nil {
		return nil, fmt.Errorf("error contando jobs: %w", err)
	}

	if summary.TotalNotifications > 0 {
		summary.SuccessRate = (float64(summary.SuccessfulSent) / float64(summary.TotalNotifications)) * 100
		summary.FailureRate = (float64(summary.Failed) / float64(summary.TotalNotifications)) * 100
	}

	if summary.TotalJobRuns > 0 {
		summary.AvgEmailsPerRun = float64(summary.TotalNotifications) / float64(summary.TotalJobRuns)
	}

	response.Summary = summary

	// Query para notificaciones por día
	byDayQuery := `
		SELECT 
			DATE(created_at) as date,
			COALESCE(SUM(CASE WHEN status = 'sent' THEN 1 ELSE 0 END), 0) as sent,
			COALESCE(SUM(CASE WHEN status = 'failed' THEN 1 ELSE 0 END), 0) as failed
		FROM notification_logs
		WHERE created_at BETWEEN ? AND ?
		GROUP BY DATE(created_at)
		ORDER BY date ASC
	`

	rows, err := r.db.QueryContext(ctx, byDayQuery, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo notificaciones por día: %w", err)
	}
	defer rows.Close()

	var byDay []dto.NotificationByDay
	for rows.Next() {
		var item dto.NotificationByDay
		err := rows.Scan(&item.Date, &item.Sent, &item.Failed)
		if err != nil {
			return nil, fmt.Errorf("error escaneando día: %w", err)
		}
		item.Day = item.Date.Format("2006-01-02")
		total := item.Sent + item.Failed
		if total > 0 {
			item.SuccessRate = (float64(item.Sent) / float64(total)) * 100
		}
		byDay = append(byDay, item)
	}
	response.ByDay = byDay

	// Query para notificaciones por estado
	byStatusQuery := `
		SELECT 
			status,
			COUNT(*) as count
		FROM notification_logs
		WHERE created_at BETWEEN ? AND ?
		GROUP BY status
	`

	rows, err = r.db.QueryContext(ctx, byStatusQuery, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo notificaciones por estado: %w", err)
	}
	defer rows.Close()

	var byStatus []dto.NotificationByStatus
	for rows.Next() {
		var item dto.NotificationByStatus
		err := rows.Scan(&item.Status, &item.Count)
		if err != nil {
			return nil, fmt.Errorf("error escaneando estado: %w", err)
		}
		if summary.TotalNotifications > 0 {
			item.Percentage = (float64(item.Count) / float64(summary.TotalNotifications)) * 100
		}
		byStatus = append(byStatus, item)
	}
	response.ByStatus = byStatus

	// Query para detalles de job runs
	jobRunsQuery := `
		SELECT 
			id, started_at, completed_at, status, cards_found,
			emails_sent, errors, error_message
		FROM job_runs
		WHERE started_at BETWEEN ? AND ?
		ORDER BY started_at DESC
		LIMIT 20
	`

	rows, err = r.db.QueryContext(ctx, jobRunsQuery, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo job runs: %w", err)
	}
	defer rows.Close()

	var jobRuns []dto.JobRunDetail
	for rows.Next() {
		var job dto.JobRunDetail
		var completedAt sql.NullTime
		var errorMsg sql.NullString

		err := rows.Scan(
			&job.ID, &job.StartedAt, &completedAt, &job.Status,
			&job.CardsFound, &job.EmailsSent, &job.Errors, &errorMsg,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando job run: %w", err)
		}

		if completedAt.Valid {
			job.CompletedAt = &completedAt.Time
			duration := completedAt.Time.Sub(job.StartedAt)
			job.Duration = duration.String()
		}

		if errorMsg.Valid {
			job.ErrorMessage = errorMsg.String
		}

		jobRuns = append(jobRuns, job)
	}
	response.JobRuns = jobRuns

	return response, nil
}
