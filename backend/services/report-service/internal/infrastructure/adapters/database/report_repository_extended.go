package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/fintrack/report-service/internal/core/domain/dto"
)

// GetAccountReport obtiene el reporte de cuentas
func (r *ReportRepository) GetAccountReport(ctx context.Context, userID string) (*dto.AccountReportResponse, error) {
	response := &dto.AccountReportResponse{
		UserID: userID,
	}

	// Query para resumen
	summaryQuery := `
		SELECT 
			COALESCE(SUM(balance), 0) as total_balance,
			COUNT(*) as total_accounts,
			COALESCE(SUM(credit_limit), 0) as total_credit_limit
		FROM accounts
		WHERE BINARY user_id = BINARY ? AND is_active = 1 AND deleted_at IS NULL
	`

	var summary dto.AccountSummary
	err := r.db.QueryRowContext(ctx, summaryQuery, userID).Scan(
		&summary.TotalBalance,
		&summary.TotalAccounts,
		&summary.TotalCreditLimit,
	)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo resumen de cuentas: %w", err)
	}

	// Contar tarjetas
	cardsCountQuery := `
		SELECT COUNT(*)
		FROM cards c
		JOIN accounts a ON BINARY c.account_id = BINARY a.id
		WHERE BINARY a.user_id = BINARY ? AND c.status = 'active' AND c.deleted_at IS NULL
	`
	err = r.db.QueryRowContext(ctx, cardsCountQuery, userID).Scan(&summary.TotalCards)
	if err != nil {
		return nil, fmt.Errorf("error contando tarjetas: %w", err)
	}

	// Calcular crédito usado (suma de transacciones pendientes en tarjetas de crédito)
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

	// Query para detalle de cuentas
	accountsQuery := `
		SELECT 
			id, account_type, name, currency, balance, 
			COALESCE(credit_limit, 0) as credit_limit, is_active
		FROM accounts
		WHERE BINARY user_id = BINARY ? AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, accountsQuery, userID)
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

	// Query para detalle de tarjetas
	cardsQuery := `
		SELECT 
			c.id, c.account_id, c.card_type, c.card_brand, c.last_four_digits,
			c.holder_name, c.status, COALESCE(c.credit_limit, 0) as credit_limit,
			COALESCE(c.nickname, '') as nickname
		FROM cards c
		JOIN accounts a ON BINARY c.account_id = BINARY a.id
		WHERE BINARY a.user_id = BINARY ? AND c.deleted_at IS NULL
		ORDER BY c.created_at DESC
	`

	cardRows, err := r.db.QueryContext(ctx, cardsQuery, userID)
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
			&card.CreditLimit, &card.Nickname,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando tarjeta: %w", err)
		}

		// Calcular balance actual de la tarjeta si es de crédito
		if card.CardType == "credit" {
			balanceQuery := `
				SELECT COALESCE(SUM(amount), 0)
				FROM transactions
				WHERE from_card_id = ? 
					AND status IN ('pending', 'completed')
					AND type = 'credit_charge'
			`
			err = r.db.QueryRowContext(ctx, balanceQuery, card.ID).Scan(&card.CurrentBalance)
			if err != nil {
				card.CurrentBalance = 0
			}
			card.AvailableCredit = card.CreditLimit - card.CurrentBalance
		}

		cards = append(cards, card)
	}
	response.Cards = cards

	// Query para distribución de cuentas
	distributionQuery := `
		SELECT 
			account_type,
			COUNT(*) as count,
			COALESCE(SUM(balance), 0) as total_balance
		FROM accounts
		WHERE BINARY user_id = BINARY ? AND is_active = 1 AND deleted_at IS NULL
		GROUP BY account_type
	`

	distRows, err := r.db.QueryContext(ctx, distributionQuery, userID)
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
	for rows.Next() {
		var item dto.ExpenseIncomeByCategory
		err := rows.Scan(&item.Category, &item.Type, &item.Count, &item.Amount)
		if err != nil {
			return nil, fmt.Errorf("error escaneando categoría: %w", err)
		}

		total := summary.TotalIncome + summary.TotalExpenses
		if total > 0 {
			item.Percentage = (item.Amount / total) * 100
		}

		byCategory = append(byCategory, item)
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
