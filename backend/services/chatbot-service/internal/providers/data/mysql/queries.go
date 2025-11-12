package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/fintrack/chatbot-service/internal/core/ports"
)

type DataProvider struct{ db *sql.DB }

func NewDataProvider(conn *Connection) *DataProvider { return &DataProvider{db: conn.DB} }

func (p *DataProvider) GetTotals(ctx context.Context, userID string, from, to time.Time) (ports.Totals, error) {
	q := `SELECT 
        SUM(CASE WHEN type IN ('debit_purchase','credit_charge','wallet_withdrawal','account_withdraw') THEN amount ELSE 0 END) AS total_gastos,
        SUM(CASE WHEN type IN ('wallet_deposit','account_deposit','credit_payment','installment_payment') THEN amount ELSE 0 END) AS total_ingresos
      FROM transactions
      WHERE user_id = ? AND status = 'completed' AND created_at BETWEEN ? AND ?`
	var exp, inc sql.NullFloat64
	err := p.db.QueryRowContext(ctx, q, userID, from, to).Scan(&exp, &inc)
	if err != nil {
		return ports.Totals{}, err
	}
	return ports.Totals{Expenses: exp.Float64, Incomes: inc.Float64}, nil
}

func (p *DataProvider) GetByType(ctx context.Context, userID string, from, to time.Time) (map[string]float64, error) {
	q := `SELECT type, SUM(amount) AS total FROM transactions WHERE user_id=? AND status='completed' AND created_at BETWEEN ? AND ? GROUP BY type`
	rows, err := p.db.QueryContext(ctx, q, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := map[string]float64{}
	for rows.Next() {
		var t string
		var total float64
		if err := rows.Scan(&t, &total); err != nil {
			return nil, err
		}
		res[t] = total
	}
	return res, rows.Err()
}

func (p *DataProvider) GetTopMerchants(ctx context.Context, userID string, from, to time.Time, limit int) ([]ports.MerchantTotal, error) {
	q := `SELECT merchant_name, SUM(amount) AS total 
          FROM transactions 
          WHERE user_id=? 
            AND status='completed' 
            AND merchant_name IS NOT NULL AND merchant_name<>'' 
            AND type IN ('debit_purchase','credit_charge','wallet_withdrawal','account_withdraw') 
            AND created_at BETWEEN ? AND ? 
          GROUP BY merchant_name 
          ORDER BY total DESC 
          LIMIT ?`
	rows, err := p.db.QueryContext(ctx, q, userID, from, to, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []ports.MerchantTotal
	for rows.Next() {
		var m sql.NullString
		var total float64
		if err := rows.Scan(&m, &total); err != nil {
			return nil, err
		}
		res = append(res, ports.MerchantTotal{Merchant: m.String, Total: total})
	}
	return res, rows.Err()
}

func (p *DataProvider) GetByAccountType(ctx context.Context, userID string, from, to time.Time) (map[string]float64, error) {
	q := `SELECT a.account_type, SUM(t.amount) AS total FROM transactions t JOIN accounts a ON (t.from_account_id=a.id OR t.to_account_id=a.id) WHERE t.user_id=? AND t.status='completed' AND t.created_at BETWEEN ? AND ? GROUP BY a.account_type`
	rows, err := p.db.QueryContext(ctx, q, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := map[string]float64{}
	for rows.Next() {
		var typ sql.NullString
		var total float64
		if err := rows.Scan(&typ, &total); err != nil {
			return nil, err
		}
		res[typ.String] = total
	}
	return res, rows.Err()
}

func (p *DataProvider) GetByCard(ctx context.Context, userID string, from, to time.Time) ([]ports.CardTotal, error) {
	q := `SELECT c.card_brand, c.last_four_digits, SUM(t.amount) AS total FROM transactions t JOIN cards c ON (t.from_card_id=c.id OR t.to_card_id=c.id) WHERE t.user_id=? AND t.status='completed' AND t.created_at BETWEEN ? AND ? GROUP BY c.card_brand, c.last_four_digits ORDER BY total DESC`
	rows, err := p.db.QueryContext(ctx, q, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []ports.CardTotal
	for rows.Next() {
		var brand sql.NullString
		var lastFour sql.NullString
		var total float64
		if err := rows.Scan(&brand, &lastFour, &total); err != nil {
			return nil, err
		}
		res = append(res, ports.CardTotal{Brand: brand.String, LastFour: lastFour.String, Total: total})
	}
	return res, rows.Err()
}

// Utility: debug print (optional)
func debugf(format string, a ...any) { _ = fmt.Sprintf(format, a...) }

// Installments: summary across user's plans and payments
func (p *DataProvider) GetInstallmentsSummary(ctx context.Context, userID string, from, to time.Time) (ports.InstallmentsSummary, error) {
	// Count active/completed and sum remaining amount
	q1 := `SELECT 
        COALESCE(SUM(CASE WHEN status='active' THEN 1 ELSE 0 END),0) AS active_cnt,
        COALESCE(SUM(CASE WHEN status='completed' THEN 1 ELSE 0 END),0) AS completed_cnt,
        COALESCE(SUM(remaining_amount),0) AS remaining_sum
      FROM installment_plans WHERE user_id=?`
	var active, completed int
	var remaining sql.NullFloat64
	if err := p.db.QueryRowContext(ctx, q1, userID).Scan(&active, &completed, &remaining); err != nil {
		return ports.InstallmentsSummary{}, err
	}

	// Overdue installments count
	q2 := `SELECT COALESCE(COUNT(*),0)
      FROM installments i JOIN installment_plans ip ON i.plan_id=ip.id 
      WHERE ip.user_id=? AND i.status='overdue'`
	var overdue int
	if err := p.db.QueryRowContext(ctx, q2, userID).Scan(&overdue); err != nil {
		return ports.InstallmentsSummary{}, err
	}

	// Next due date among pending/overdue
	q3 := `SELECT MIN(i.due_date)
      FROM installments i JOIN installment_plans ip ON i.plan_id=ip.id 
      WHERE ip.user_id=? AND i.status IN ('pending','overdue')`
	var nextDue sql.NullTime
	if err := p.db.QueryRowContext(ctx, q3, userID).Scan(&nextDue); err != nil {
		return ports.InstallmentsSummary{}, err
	}

	// Payments in period (from transactions)
	q4 := `SELECT COALESCE(SUM(amount),0) FROM transactions 
      WHERE user_id=? AND status='completed' AND type='installment_payment' AND created_at BETWEEN ? AND ?`
	var paid sql.NullFloat64
	if err := p.db.QueryRowContext(ctx, q4, userID, from, to).Scan(&paid); err != nil {
		return ports.InstallmentsSummary{}, err
	}

	var nextPtr *time.Time
	if nextDue.Valid {
		t := nextDue.Time
		nextPtr = &t
	}
	return ports.InstallmentsSummary{
		Active:          active,
		Completed:       completed,
		Overdue:         overdue,
		RemainingAmount: remaining.Float64,
		NextDueDate:     nextPtr,
		PaidInPeriod:    paid.Float64,
	}, nil
}

// Installment plans: high-level info per plan from view
func (p *DataProvider) GetInstallmentPlans(ctx context.Context, userID string) ([]ports.InstallmentPlanInfo, error) {
	q := `SELECT 
        p.id, p.card_id, p.merchant_name, p.description,
        p.installments_count, (p.installments_count - p.paid_installments) as remaining_installments,
        p.remaining_amount, p.status, 
        (SELECT MIN(i.due_date) FROM installments i WHERE i.plan_id = p.id AND i.status = 'pending') as next_due_date,
        (SELECT i.amount FROM installments i WHERE i.plan_id = p.id AND i.status = 'pending' ORDER BY i.due_date LIMIT 1) as next_installment_amount,
        ROUND((p.paid_installments / p.installments_count) * 100, 2) as completion_percentage, 
        p.created_at
      FROM installment_plans p
      WHERE p.user_id=? AND p.status = 'active'
      ORDER BY p.status DESC, next_due_date ASC`
	rows, err := p.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []ports.InstallmentPlanInfo{}
	for rows.Next() {
		var (
			id           sql.NullString
			cardID       sql.NullString
			merchantName sql.NullString
			description  sql.NullString
			instCount    sql.NullInt64
			remInst      sql.NullInt64
			remAmt       sql.NullFloat64
			status       sql.NullString
			nextDue      sql.NullTime
			nextAmt      sql.NullFloat64
			compl        sql.NullFloat64
			createdAt    sql.NullTime
		)
		if err := rows.Scan(&id, &cardID, &merchantName, &description, &instCount, &remInst, &remAmt, &status, &nextDue, &nextAmt, &compl, &createdAt); err != nil {
			return nil, err
		}
		var nextPtr *time.Time
		if nextDue.Valid {
			t := nextDue.Time
			nextPtr = &t
		}
		var createdPtr *time.Time
		if createdAt.Valid {
			t := createdAt.Time
			createdPtr = &t
		}
		res = append(res, ports.InstallmentPlanInfo{
			ID:                    id.String,
			CardID:                cardID.String,
			MerchantName:          merchantName.String,
			Description:           description.String,
			InstallmentsCount:     int(instCount.Int64),
			RemainingInstallments: int(remInst.Int64),
			RemainingAmount:       remAmt.Float64,
			Status:                status.String,
			NextDueDate:           nextPtr,
			NextInstallmentAmount: nextAmt.Float64,
			CompletionPercentage:  compl.Float64,
			CreatedAt:             createdPtr,
		})
	}
	return res, rows.Err()
}

// GetInstallmentsByMonth obtiene cuotas pendientes agrupadas por mes
func (p *DataProvider) GetInstallmentsByMonth(ctx context.Context, userID string) (map[string]ports.InstallmentMonthSummary, error) {
	q := `SELECT 
        DATE_FORMAT(i.due_date, '%Y-%m') as ym,
        COUNT(*) as count,
        SUM(i.amount) as total
      FROM installments i
      JOIN installment_plans p ON i.plan_id = p.id
      WHERE p.user_id = ? 
        AND i.status = 'pending'
        AND i.due_date >= CURDATE()
      GROUP BY ym
      ORDER BY ym`

	rows, err := p.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]ports.InstallmentMonthSummary)
	for rows.Next() {
		var yearMonth string
		var count int
		var total float64
		if err := rows.Scan(&yearMonth, &count, &total); err != nil {
			return nil, err
		}
		result[yearMonth] = ports.InstallmentMonthSummary{
			YearMonth: yearMonth,
			Count:     count,
			Total:     total,
		}
	}
	return result, rows.Err()
}

// GetRecentTransactions obtiene transacciones recientes con detalles completos
func (p *DataProvider) GetRecentTransactions(ctx context.Context, userID string, from, to time.Time, limit int) ([]ports.TransactionDetail, error) {
	q := `SELECT 
        id, type, amount, 
        COALESCE(merchant_name, '') as merchant_name,
        COALESCE(description, '') as description,
        status, created_at,
        COALESCE(from_account_id, '') as from_account_id,
        COALESCE(to_account_id, '') as to_account_id,
        COALESCE(from_card_id, '') as from_card_id,
        COALESCE(to_card_id, '') as to_card_id,
        COALESCE(currency, 'ARS') as currency
      FROM transactions 
      WHERE user_id = ? AND created_at BETWEEN ? AND ?
      ORDER BY created_at DESC 
      LIMIT ?`

	rows, err := p.db.QueryContext(ctx, q, userID, from, to, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []ports.TransactionDetail
	for rows.Next() {
		var t ports.TransactionDetail
		if err := rows.Scan(&t.ID, &t.Type, &t.Amount, &t.MerchantName, &t.Description,
			&t.Status, &t.CreatedAt, &t.FromAccountID, &t.ToAccountID, &t.FromCardID, &t.ToCardID, &t.Currency); err != nil {
			return nil, err
		}
		res = append(res, t)
	}
	return res, rows.Err()
}

// GetAccountsInfo obtiene información de todas las cuentas del usuario
func (p *DataProvider) GetAccountsInfo(ctx context.Context, userID string) ([]ports.AccountInfo, error) {
	q := `SELECT 
        id, account_type, balance, 
        currency,
        CASE WHEN is_active = 1 THEN 'active' ELSE 'inactive' END as status,
        created_at
      FROM accounts 
      WHERE user_id = ? AND deleted_at IS NULL
      ORDER BY created_at DESC`

	rows, err := p.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []ports.AccountInfo
	for rows.Next() {
		var a ports.AccountInfo
		if err := rows.Scan(&a.ID, &a.AccountType, &a.Balance, &a.Currency, &a.Status, &a.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, a)
	}
	return res, rows.Err()
}

// GetCardsInfo obtiene información de todas las tarjetas del usuario
func (p *DataProvider) GetCardsInfo(ctx context.Context, userID string) ([]ports.CardInfo, error) {
	q := `SELECT 
        c.id, c.card_brand, c.last_four_digits, c.card_type, c.status,
        COALESCE(c.credit_limit, 0) as credit_limit,
        COALESCE(c.balance, 0) as current_debt,
        c.created_at
      FROM cards c
      INNER JOIN accounts a ON c.account_id = a.id 
      WHERE a.user_id = ? AND c.deleted_at IS NULL
      ORDER BY c.created_at DESC`

	rows, err := p.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []ports.CardInfo
	for rows.Next() {
		var c ports.CardInfo
		if err := rows.Scan(&c.ID, &c.CardBrand, &c.LastFour, &c.CardType, &c.Status,
			&c.CreditLimit, &c.CurrentDebt, &c.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, rows.Err()
}

// GetExchangeRates obtiene tasas de cambio recientes
func (p *DataProvider) GetExchangeRates(ctx context.Context, userID string, from, to time.Time) ([]ports.ExchangeRateInfo, error) {
	q := `SELECT 
        id, from_currency, to_currency, rate, 
        COALESCE(source, 'manual') as source,
        created_at
      FROM exchange_rates 
      WHERE created_at BETWEEN ? AND ?
      ORDER BY created_at DESC 
      LIMIT 10`

	rows, err := p.db.QueryContext(ctx, q, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []ports.ExchangeRateInfo
	for rows.Next() {
		var e ports.ExchangeRateInfo
		if err := rows.Scan(&e.ID, &e.FromCurrency, &e.ToCurrency, &e.Rate, &e.Source, &e.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, e)
	}
	return res, rows.Err()
}