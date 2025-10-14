package ports

import (
	"time"

	"github.com/fintrack/notification-service/internal/core/domain/entities"
)

// CardRepository define las operaciones de repositorio para tarjetas
type CardRepository interface {
	GetCardsDueTomorrow() ([]*entities.Card, error)
}

// InstallmentRepository define las operaciones de repositorio para cuotas
type InstallmentRepository interface {
	GetPendingInstallmentsByCard(cardID string, maxDueDate time.Time) ([]*entities.Installment, error)
}

// NotificationRepository define las operaciones de repositorio para notificaciones
type NotificationRepository interface {
	SaveNotificationLog(log *entities.NotificationLog) error
	SaveJobRun(jobRun *entities.JobRun) error
	UpdateJobRun(jobRun *entities.JobRun) error
	GetJobRunHistory(limit int) ([]*entities.JobRun, error)
	GetNotificationLogs(jobRunID string, limit int) ([]*entities.NotificationLog, error)
}

// EmailService define las operaciones de servicio de email
type EmailService interface {
	SendCardDueNotification(notification *entities.CardDueNotification) error
}

// NotificationService define las operaciones del servicio de notificaciones
type NotificationService interface {
	ProcessCardDueNotifications() error
	TriggerManualJob() error
	GetJobHistory(limit int) ([]*entities.JobRun, error)
	GetNotificationLogs(jobRunID string, limit int) ([]*entities.NotificationLog, error)
}
