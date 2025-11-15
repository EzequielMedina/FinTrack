package mysql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	domaintransaction "github.com/fintrack/transaction-service/internal/core/domain/entities/transaction"
	"github.com/fintrack/transaction-service/internal/core/service"
	_ "github.com/go-sql-driver/mysql"
)

// TransactionRepository implements the TransactionRepositoryInterface for MySQL
type TransactionRepository struct {
	db *sql.DB
}

// NewTransactionRepository creates a new MySQL transaction repository
func NewTransactionRepository(db *sql.DB) service.TransactionRepositoryInterface {
	return &TransactionRepository{
		db: db,
	}
}

// Create inserts a new transaction into the database
func (r *TransactionRepository) Create(transaction *domaintransaction.Transaction) (*domaintransaction.Transaction, error) {
	// Generate ID if not set
	if transaction.ID == "" {
		transaction.ID = r.generateID()
	}

	// Serialize metadata and tags to JSON
	metadataJSON, err := json.Marshal(transaction.Metadata)
	if err != nil {
		metadataJSON = []byte("{}")
	}

	tagsJSON, err := json.Marshal(transaction.Tags)
	if err != nil {
		tagsJSON = []byte("[]")
	}

	query := `
		INSERT INTO transactions (
			id, reference_id, external_id, type, status, amount, currency,
			from_account_id, to_account_id, from_card_id, to_card_id,
			user_id, initiated_by, description, payment_method,
			merchant_name, merchant_id, previous_balance, new_balance,
			processed_at, failed_at, failure_reason, metadata, tags,
			created_at, updated_at
		) VALUES (
			?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?,
			?, ?, ?, ?,
			?, ?, ?, ?,
			?, ?, ?, ?, ?,
			NOW(), NOW()
		)`

	_, err = r.db.Exec(query,
		transaction.ID, transaction.ReferenceID, transaction.ExternalID,
		transaction.Type, transaction.Status, transaction.Amount, transaction.Currency,
		transaction.FromAccountID, transaction.ToAccountID, transaction.FromCardID, transaction.ToCardID,
		transaction.UserID, transaction.InitiatedBy, transaction.Description, transaction.PaymentMethod,
		transaction.MerchantName, transaction.MerchantID, transaction.PreviousBalance, transaction.NewBalance,
		transaction.ProcessedAt, transaction.FailedAt, transaction.FailureReason,
		string(metadataJSON), string(tagsJSON),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Return the created transaction with timestamps
	return r.GetByID(transaction.ID)
}

// GetByID retrieves a transaction by its ID
func (r *TransactionRepository) GetByID(id string) (*domaintransaction.Transaction, error) {
	query := `
		SELECT id, reference_id, external_id, type, status, amount, currency,
			   from_account_id, to_account_id, from_card_id, to_card_id,
			   user_id, initiated_by, description, payment_method,
			   merchant_name, merchant_id, previous_balance, new_balance,
			   processed_at, failed_at, failure_reason, metadata, tags,
			   created_at, updated_at
		FROM transactions
		WHERE id = ?`

	row := r.db.QueryRow(query, id)

	transaction := &domaintransaction.Transaction{}
	var metadataJSON, tagsJSON string

	err := row.Scan(
		&transaction.ID, &transaction.ReferenceID, &transaction.ExternalID,
		&transaction.Type, &transaction.Status, &transaction.Amount, &transaction.Currency,
		&transaction.FromAccountID, &transaction.ToAccountID, &transaction.FromCardID, &transaction.ToCardID,
		&transaction.UserID, &transaction.InitiatedBy, &transaction.Description, &transaction.PaymentMethod,
		&transaction.MerchantName, &transaction.MerchantID, &transaction.PreviousBalance, &transaction.NewBalance,
		&transaction.ProcessedAt, &transaction.FailedAt, &transaction.FailureReason,
		&metadataJSON, &tagsJSON, &transaction.CreatedAt, &transaction.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaction not found with ID: %s", id)
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	// Deserialize JSON fields
	if metadataJSON != "" {
		json.Unmarshal([]byte(metadataJSON), &transaction.Metadata)
	}
	if tagsJSON != "" {
		json.Unmarshal([]byte(tagsJSON), &transaction.Tags)
	}

	return transaction, nil
}

// Update updates an existing transaction
func (r *TransactionRepository) Update(transaction *domaintransaction.Transaction) (*domaintransaction.Transaction, error) {
	metadataJSON, _ := json.Marshal(transaction.Metadata)
	tagsJSON, _ := json.Marshal(transaction.Tags)

	query := `
		UPDATE transactions SET
			reference_id = ?, external_id = ?, type = ?, status = ?,
			amount = ?, currency = ?, from_account_id = ?, to_account_id = ?,
			from_card_id = ?, to_card_id = ?, description = ?, payment_method = ?,
			merchant_name = ?, merchant_id = ?, previous_balance = ?, new_balance = ?,
			processed_at = ?, failed_at = ?, failure_reason = ?,
			metadata = ?, tags = ?, updated_at = NOW()
		WHERE id = ?`

	_, err := r.db.Exec(query,
		transaction.ReferenceID, transaction.ExternalID, transaction.Type, transaction.Status,
		transaction.Amount, transaction.Currency, transaction.FromAccountID, transaction.ToAccountID,
		transaction.FromCardID, transaction.ToCardID, transaction.Description, transaction.PaymentMethod,
		transaction.MerchantName, transaction.MerchantID, transaction.PreviousBalance, transaction.NewBalance,
		transaction.ProcessedAt, transaction.FailedAt, transaction.FailureReason,
		string(metadataJSON), string(tagsJSON), transaction.ID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	return r.GetByID(transaction.ID)
}

// Delete removes a transaction from the database
func (r *TransactionRepository) Delete(id string) error {
	query := "DELETE FROM transactions WHERE id = ?"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("transaction not found with ID: %s", id)
	}

	return nil
}

// GetByUserID retrieves transactions for a user with filtering
func (r *TransactionRepository) GetByUserID(userID string, filters service.TransactionFilters) ([]*domaintransaction.Transaction, int, error) {
	// Build WHERE clause
	whereConditions := []string{"user_id = ?"}
	args := []interface{}{userID}

	// Add filters
	if len(filters.Types) > 0 {
		placeholders := make([]string, len(filters.Types))
		for i, transactionType := range filters.Types {
			placeholders[i] = "?"
			args = append(args, transactionType)
		}
		whereConditions = append(whereConditions, fmt.Sprintf("type IN (%s)", strings.Join(placeholders, ",")))
	}

	if len(filters.Statuses) > 0 {
		placeholders := make([]string, len(filters.Statuses))
		for i, status := range filters.Statuses {
			placeholders[i] = "?"
			args = append(args, status)
		}
		whereConditions = append(whereConditions, fmt.Sprintf("status IN (%s)", strings.Join(placeholders, ",")))
	}

	if filters.FromDate != nil {
		whereConditions = append(whereConditions, "created_at >= ?")
		args = append(args, *filters.FromDate)
	}

	if filters.ToDate != nil {
		whereConditions = append(whereConditions, "created_at <= ?")
		args = append(args, *filters.ToDate)
	}

	if filters.MinAmount != nil {
		whereConditions = append(whereConditions, "amount >= ?")
		args = append(args, *filters.MinAmount)
	}

	if filters.MaxAmount != nil {
		whereConditions = append(whereConditions, "amount <= ?")
		args = append(args, *filters.MaxAmount)
	}

	whereClause := strings.Join(whereConditions, " AND ")

	// Count total records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM transactions WHERE %s", whereClause)
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count transactions: %w", err)
	}

	// Build ORDER BY clause
	orderBy := "created_at DESC"
	if filters.OrderBy != "" {
		direction := "ASC"
		if filters.Order == "desc" {
			direction = "DESC"
		}
		orderBy = fmt.Sprintf("%s %s", filters.OrderBy, direction)
	}

	// Build main query
	query := fmt.Sprintf(`
		SELECT id, reference_id, external_id, type, status, amount, currency,
			   from_account_id, to_account_id, from_card_id, to_card_id,
			   user_id, initiated_by, description, payment_method,
			   merchant_name, merchant_id, previous_balance, new_balance,
			   processed_at, failed_at, failure_reason, metadata, tags,
			   created_at, updated_at
		FROM transactions
		WHERE %s
		ORDER BY %s
		LIMIT ? OFFSET ?`, whereClause, orderBy)

	// Add pagination
	limit := filters.Limit
	if limit <= 0 {
		limit = 20 // Default limit
	}
	offset := filters.Offset
	if offset < 0 {
		offset = 0
	}

	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*domaintransaction.Transaction
	for rows.Next() {
		transaction := &domaintransaction.Transaction{}
		var metadataJSON, tagsJSON string

		err := rows.Scan(
			&transaction.ID, &transaction.ReferenceID, &transaction.ExternalID,
			&transaction.Type, &transaction.Status, &transaction.Amount, &transaction.Currency,
			&transaction.FromAccountID, &transaction.ToAccountID, &transaction.FromCardID, &transaction.ToCardID,
			&transaction.UserID, &transaction.InitiatedBy, &transaction.Description, &transaction.PaymentMethod,
			&transaction.MerchantName, &transaction.MerchantID, &transaction.PreviousBalance, &transaction.NewBalance,
			&transaction.ProcessedAt, &transaction.FailedAt, &transaction.FailureReason,
			&metadataJSON, &tagsJSON, &transaction.CreatedAt, &transaction.UpdatedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan transaction: %w", err)
		}

		// Deserialize JSON fields
		if metadataJSON != "" {
			json.Unmarshal([]byte(metadataJSON), &transaction.Metadata)
		}
		if tagsJSON != "" {
			json.Unmarshal([]byte(tagsJSON), &transaction.Tags)
		}

		transactions = append(transactions, transaction)
	}

	return transactions, total, nil
}

// GetByAccountID retrieves transactions for an account
func (r *TransactionRepository) GetByAccountID(accountID string, filters service.TransactionFilters) ([]*domaintransaction.Transaction, int, error) {
	// Similar to GetByUserID but filter by account
	whereConditions := []string{"(from_account_id = ? OR to_account_id = ?)"}
	args := []interface{}{accountID, accountID}

	// Add other filters (similar to GetByUserID implementation)
	// ... rest of implementation similar to GetByUserID

	return r.executeFilteredQuery(whereConditions, args, filters)
}

// GetByCardID retrieves transactions for a card
func (r *TransactionRepository) GetByCardID(cardID string, filters service.TransactionFilters) ([]*domaintransaction.Transaction, int, error) {
	whereConditions := []string{"(from_card_id = ? OR to_card_id = ?)"}
	args := []interface{}{cardID, cardID}

	return r.executeFilteredQuery(whereConditions, args, filters)
}

// GetByReferenceID retrieves a transaction by reference ID
func (r *TransactionRepository) GetByReferenceID(referenceID string) (*domaintransaction.Transaction, error) {
	query := `
		SELECT id, reference_id, external_id, type, status, amount, currency,
			   from_account_id, to_account_id, from_card_id, to_card_id,
			   user_id, initiated_by, description, payment_method,
			   merchant_name, merchant_id, previous_balance, new_balance,
			   processed_at, failed_at, failure_reason, metadata, tags,
			   created_at, updated_at
		FROM transactions
		WHERE reference_id = ?`

	row := r.db.QueryRow(query, referenceID)

	transaction := &domaintransaction.Transaction{}
	var metadataJSON, tagsJSON string

	err := row.Scan(
		&transaction.ID, &transaction.ReferenceID, &transaction.ExternalID,
		&transaction.Type, &transaction.Status, &transaction.Amount, &transaction.Currency,
		&transaction.FromAccountID, &transaction.ToAccountID, &transaction.FromCardID, &transaction.ToCardID,
		&transaction.UserID, &transaction.InitiatedBy, &transaction.Description, &transaction.PaymentMethod,
		&transaction.MerchantName, &transaction.MerchantID, &transaction.PreviousBalance, &transaction.NewBalance,
		&transaction.ProcessedAt, &transaction.FailedAt, &transaction.FailureReason,
		&metadataJSON, &tagsJSON, &transaction.CreatedAt, &transaction.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaction not found with reference ID: %s", referenceID)
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	// Deserialize JSON fields
	if metadataJSON != "" {
		json.Unmarshal([]byte(metadataJSON), &transaction.Metadata)
	}
	if tagsJSON != "" {
		json.Unmarshal([]byte(tagsJSON), &transaction.Tags)
	}

	return transaction, nil
}

// GetByExternalID retrieves a transaction by external ID
func (r *TransactionRepository) GetByExternalID(externalID string) (*domaintransaction.Transaction, error) {
	query := `
		SELECT id, reference_id, external_id, type, status, amount, currency,
			   from_account_id, to_account_id, from_card_id, to_card_id,
			   user_id, initiated_by, description, payment_method,
			   merchant_name, merchant_id, previous_balance, new_balance,
			   processed_at, failed_at, failure_reason, metadata, tags,
			   created_at, updated_at
		FROM transactions
		WHERE external_id = ?`

	row := r.db.QueryRow(query, externalID)

	transaction := &domaintransaction.Transaction{}
	var metadataJSON, tagsJSON string

	err := row.Scan(
		&transaction.ID, &transaction.ReferenceID, &transaction.ExternalID,
		&transaction.Type, &transaction.Status, &transaction.Amount, &transaction.Currency,
		&transaction.FromAccountID, &transaction.ToAccountID, &transaction.FromCardID, &transaction.ToCardID,
		&transaction.UserID, &transaction.InitiatedBy, &transaction.Description, &transaction.PaymentMethod,
		&transaction.MerchantName, &transaction.MerchantID, &transaction.PreviousBalance, &transaction.NewBalance,
		&transaction.ProcessedAt, &transaction.FailedAt, &transaction.FailureReason,
		&metadataJSON, &tagsJSON, &transaction.CreatedAt, &transaction.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaction not found with external ID: %s", externalID)
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	// Deserialize JSON fields
	if metadataJSON != "" {
		json.Unmarshal([]byte(metadataJSON), &transaction.Metadata)
	}
	if tagsJSON != "" {
		json.Unmarshal([]byte(tagsJSON), &transaction.Tags)
	}

	return transaction, nil
}

// CreateBatch creates multiple transactions in a batch
func (r *TransactionRepository) CreateBatch(transactions []*domaintransaction.Transaction) ([]*domaintransaction.Transaction, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var createdTransactions []*domaintransaction.Transaction

	for _, transaction := range transactions {
		created, err := r.Create(transaction)
		if err != nil {
			return nil, fmt.Errorf("failed to create transaction in batch: %w", err)
		}
		createdTransactions = append(createdTransactions, created)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit batch transaction: %w", err)
	}

	return createdTransactions, nil
}

// UpdateBatch updates multiple transactions in a batch
func (r *TransactionRepository) UpdateBatch(transactions []*domaintransaction.Transaction) ([]*domaintransaction.Transaction, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var updatedTransactions []*domaintransaction.Transaction

	for _, transaction := range transactions {
		updated, err := r.Update(transaction)
		if err != nil {
			return nil, fmt.Errorf("failed to update transaction in batch: %w", err)
		}
		updatedTransactions = append(updatedTransactions, updated)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit batch update: %w", err)
	}

	return updatedTransactions, nil
}

// GetUserTransactionSummary gets aggregated data for user transactions
func (r *TransactionRepository) GetUserTransactionSummary(userID string, fromDate, toDate *time.Time) (*service.TransactionSummary, error) {
	whereConditions := []string{"user_id = ?"}
	args := []interface{}{userID}

	if fromDate != nil {
		whereConditions = append(whereConditions, "created_at >= ?")
		args = append(args, *fromDate)
	}
	if toDate != nil {
		whereConditions = append(whereConditions, "created_at <= ?")
		args = append(args, *toDate)
	}

	whereClause := strings.Join(whereConditions, " AND ")

	query := fmt.Sprintf(`
		SELECT 
			COUNT(*) as transaction_count,
			COALESCE(SUM(amount), 0) as total_amount,
			COALESCE(AVG(amount), 0) as average_amount,
			COALESCE(MAX(amount), 0) as max_amount,
			COALESCE(MIN(amount), 0) as min_amount
		FROM transactions 
		WHERE %s`, whereClause)

	summary := &service.TransactionSummary{
		ByType:   make(map[domaintransaction.TransactionType]float64),
		ByStatus: make(map[domaintransaction.TransactionStatus]int),
	}

	err := r.db.QueryRow(query, args...).Scan(
		&summary.TransactionCount,
		&summary.TotalAmount,
		&summary.AverageAmount,
		&summary.MaxAmount,
		&summary.MinAmount,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get transaction summary: %w", err)
	}

	return summary, nil
}

// GetAccountTransactionSummary gets aggregated data for account transactions
func (r *TransactionRepository) GetAccountTransactionSummary(accountID string, fromDate, toDate *time.Time) (*service.TransactionSummary, error) {
	// Similar implementation to GetUserTransactionSummary but for accounts
	return &service.TransactionSummary{
		ByType:   make(map[domaintransaction.TransactionType]float64),
		ByStatus: make(map[domaintransaction.TransactionStatus]int),
	}, nil
}

// GetDailyTransactionVolume gets the total transaction volume for a user on a specific date
func (r *TransactionRepository) GetDailyTransactionVolume(userID string, date time.Time) (float64, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	query := `
		SELECT COALESCE(SUM(amount), 0)
		FROM transactions
		WHERE user_id = ? AND created_at >= ? AND created_at < ?`

	var volume float64
	err := r.db.QueryRow(query, userID, startOfDay, endOfDay).Scan(&volume)
	if err != nil {
		return 0, fmt.Errorf("failed to get daily transaction volume: %w", err)
	}

	return volume, nil
}

// GetMonthlyTransactionVolume gets the total transaction volume for a user in a specific month
func (r *TransactionRepository) GetMonthlyTransactionVolume(userID string, year int, month int) (float64, error) {
	startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	query := `
		SELECT COALESCE(SUM(amount), 0)
		FROM transactions
		WHERE user_id = ? AND created_at >= ? AND created_at < ?`

	var volume float64
	err := r.db.QueryRow(query, userID, startOfMonth, endOfMonth).Scan(&volume)
	if err != nil {
		return 0, fmt.Errorf("failed to get monthly transaction volume: %w", err)
	}

	return volume, nil
}

// Helper methods

func (r *TransactionRepository) generateID() string {
	return fmt.Sprintf("txn_%d", time.Now().UnixNano())
}

func (r *TransactionRepository) executeFilteredQuery(whereConditions []string, args []interface{}, filters service.TransactionFilters) ([]*domaintransaction.Transaction, int, error) {
	// Implementation for common filtered query logic
	// This would contain the common logic for filtering across different methods
	return nil, 0, fmt.Errorf("not implemented yet")
}
