package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/fintrack/notification-service/internal/core/domain/entities"
	_ "github.com/go-sql-driver/mysql"
)

// CardRepository implementa las operaciones de repositorio para tarjetas
type CardRepository struct {
	db *sql.DB
}

// NewCardRepository crea un nuevo repositorio de tarjetas
func NewCardRepository(db *sql.DB) *CardRepository {
	return &CardRepository{db: db}
}

// GetCardsDueTomorrow obtiene tarjetas que vencen mañana con información del usuario
func (r *CardRepository) GetCardsDueTomorrow() ([]*entities.Card, error) {
	query := `
		SELECT 
			c.id, c.account_id as user_id, c.nickname as card_name, 
			c.card_brand as bank_name, c.last_four_digits, c.due_date,
			u.email, u.first_name, u.last_name
		FROM cards c
		JOIN accounts a ON c.account_id = a.id
		JOIN users u ON a.user_id = u.id
		WHERE DATE(c.due_date) = DATE(NOW() + INTERVAL 1 DAY)
		  AND c.status = 'active'
		  AND c.card_type = 'credit'
		  AND c.due_date IS NOT NULL
		  AND u.is_active = 1
		ORDER BY c.due_date ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying cards due tomorrow: %w", err)
	}
	defer rows.Close()

	var cards []*entities.Card
	for rows.Next() {
		card := &entities.Card{}
		err := rows.Scan(
			&card.ID,
			&card.UserID,
			&card.CardName,
			&card.BankName,
			&card.LastFour,
			&card.DueDate,
			&card.UserEmail,
			&card.UserFirstName,
			&card.UserLastName,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning card row: %w", err)
		}

		// Valores por defecto
		card.IsActive = true
		card.CardType = "credit"

		cards = append(cards, card)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating card rows: %w", err)
	}

	return cards, nil
}

// UpdateExpiredDueDates actualiza las fechas de vencimiento de tarjetas que vencieron ayer
func (r *CardRepository) UpdateExpiredDueDates() (int, error) {
	// Query para actualizar tarjetas que vencieron ayer al próximo mes
	query := `
		UPDATE cards 
		SET due_date = DATE_ADD(due_date, INTERVAL 1 MONTH)
		WHERE DATE(due_date) = DATE(NOW() - INTERVAL 1 DAY)
		  AND status = 'active'
		  AND card_type = 'credit'
		  AND due_date IS NOT NULL
	`

	result, err := r.db.Exec(query)
	if err != nil {
		return 0, fmt.Errorf("error updating expired due dates: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error getting rows affected: %w", err)
	}

	return int(rowsAffected), nil
}

// InstallmentRepository implementa las operaciones de repositorio para cuotas
type InstallmentRepository struct {
	db *sql.DB
}

// NewInstallmentRepository crea un nuevo repositorio de cuotas
func NewInstallmentRepository(db *sql.DB) *InstallmentRepository {
	return &InstallmentRepository{db: db}
}

// GetPendingInstallmentsByCard obtiene cuotas pendientes para una tarjeta hasta su fecha de vencimiento
func (r *InstallmentRepository) GetPendingInstallmentsByCard(cardID string, maxDueDate time.Time) ([]*entities.Installment, error) {
	query := `
		SELECT 
			i.id, i.plan_id, i.amount, i.due_date, i.status,
			ip.description, ip.merchant_name, i.installment_number,
			ip.installments_count
		FROM installments i
		JOIN installment_plans ip ON i.plan_id = ip.id
		WHERE ip.card_id = ?
		  AND i.status IN ('pending', 'overdue')
		  AND i.due_date <= ?
		ORDER BY i.due_date ASC
	`

	rows, err := r.db.Query(query, cardID, maxDueDate)
	if err != nil {
		return nil, fmt.Errorf("error querying installments: %w", err)
	}
	defer rows.Close()

	var installments []*entities.Installment
	for rows.Next() {
		installment := &entities.Installment{}
		err := rows.Scan(
			&installment.ID,
			&installment.PlanID,
			&installment.Amount,
			&installment.DueDate,
			&installment.Status,
			&installment.Description,
			&installment.MerchantName,
			&installment.InstallmentNum,
			&installment.TotalInstallments,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning installment row: %w", err)
		}

		installments = append(installments, installment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating installment rows: %w", err)
	}

	return installments, nil
}

// NotificationRepository implementa las operaciones de repositorio para notificaciones
type NotificationRepository struct {
	db *sql.DB
}

// NewNotificationRepository crea un nuevo repositorio de notificaciones
func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

// SaveNotificationLog guarda un log de notificación
func (r *NotificationRepository) SaveNotificationLog(log *entities.NotificationLog) error {
	query := `
		INSERT INTO notification_logs 
		(id, job_run_id, card_id, user_id, email, status, error_message, sent_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW())
	`

	_, err := r.db.Exec(query,
		log.ID,
		log.JobRunID,
		log.CardID,
		log.UserID,
		log.Email,
		log.Status,
		log.ErrorMessage,
		log.SentAt,
	)

	if err != nil {
		return fmt.Errorf("error saving notification log: %w", err)
	}

	return nil
}

// SaveJobRun guarda un registro de ejecución del job
func (r *NotificationRepository) SaveJobRun(jobRun *entities.JobRun) error {
	query := `
		INSERT INTO job_runs 
		(id, started_at, status, cards_found, emails_sent, errors, error_message)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		jobRun.ID,
		jobRun.StartedAt,
		jobRun.Status,
		jobRun.CardsFound,
		jobRun.EmailsSent,
		jobRun.Errors,
		jobRun.ErrorMessage,
	)

	if err != nil {
		return fmt.Errorf("error saving job run: %w", err)
	}

	return nil
}

// UpdateJobRun actualiza un registro de ejecución del job
func (r *NotificationRepository) UpdateJobRun(jobRun *entities.JobRun) error {
	query := `
		UPDATE job_runs 
		SET completed_at = ?, status = ?, cards_found = ?, emails_sent = ?, errors = ?, error_message = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query,
		jobRun.CompletedAt,
		jobRun.Status,
		jobRun.CardsFound,
		jobRun.EmailsSent,
		jobRun.Errors,
		jobRun.ErrorMessage,
		jobRun.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating job run: %w", err)
	}

	return nil
}

// GetJobRunHistory obtiene el historial de ejecuciones del job
func (r *NotificationRepository) GetJobRunHistory(limit int) ([]*entities.JobRun, error) {
	query := `
		SELECT id, started_at, completed_at, status, cards_found, emails_sent, errors, error_message
		FROM job_runs
		ORDER BY started_at DESC
		LIMIT ?
	`

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("error querying job run history: %w", err)
	}
	defer rows.Close()

	var jobRuns []*entities.JobRun
	for rows.Next() {
		jobRun := &entities.JobRun{}
		err := rows.Scan(
			&jobRun.ID,
			&jobRun.StartedAt,
			&jobRun.CompletedAt,
			&jobRun.Status,
			&jobRun.CardsFound,
			&jobRun.EmailsSent,
			&jobRun.Errors,
			&jobRun.ErrorMessage,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning job run row: %w", err)
		}

		jobRuns = append(jobRuns, jobRun)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating job run rows: %w", err)
	}

	return jobRuns, nil
}

// GetNotificationLogs obtiene los logs de notificaciones para un job run
func (r *NotificationRepository) GetNotificationLogs(jobRunID string, limit int) ([]*entities.NotificationLog, error) {
	query := `
		SELECT id, job_run_id, card_id, user_id, email, status, error_message, sent_at, created_at
		FROM notification_logs
		WHERE job_run_id = ?
		ORDER BY created_at DESC
		LIMIT ?
	`

	rows, err := r.db.Query(query, jobRunID, limit)
	if err != nil {
		return nil, fmt.Errorf("error querying notification logs: %w", err)
	}
	defer rows.Close()

	var logs []*entities.NotificationLog
	for rows.Next() {
		log := &entities.NotificationLog{}
		err := rows.Scan(
			&log.ID,
			&log.JobRunID,
			&log.CardID,
			&log.UserID,
			&log.Email,
			&log.Status,
			&log.ErrorMessage,
			&log.SentAt,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning notification log row: %w", err)
		}

		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating notification log rows: %w", err)
	}

	return logs, nil
}
