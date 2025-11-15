package service

import (
	"fmt"
	"log"
	"time"

	"github.com/fintrack/notification-service/internal/core/domain/entities"
	"github.com/fintrack/notification-service/internal/core/ports"
	"github.com/google/uuid"
)

// NotificationService implementa la lógica de negocio para notificaciones
type NotificationService struct {
	cardRepo         ports.CardRepository
	installmentRepo  ports.InstallmentRepository
	notificationRepo ports.NotificationRepository
	emailService     ports.EmailService
}

// NewNotificationService crea un nuevo servicio de notificaciones
func NewNotificationService(
	cardRepo ports.CardRepository,
	installmentRepo ports.InstallmentRepository,
	notificationRepo ports.NotificationRepository,
	emailService ports.EmailService,
) *NotificationService {
	return &NotificationService{
		cardRepo:         cardRepo,
		installmentRepo:  installmentRepo,
		notificationRepo: notificationRepo,
		emailService:     emailService,
	}
}

// ProcessCardDueNotifications procesa las notificaciones de tarjetas que vencen mañana
func (s *NotificationService) ProcessCardDueNotifications() error {
	jobRunID := uuid.New().String()
	jobRun := &entities.JobRun{
		ID:        jobRunID,
		StartedAt: time.Now(),
		Status:    "running",
	}

	// Guardar el inicio del job
	if err := s.notificationRepo.SaveJobRun(jobRun); err != nil {
		log.Printf("Error saving job run: %v", err)
		return fmt.Errorf("error saving job run: %w", err)
	}

	log.Printf("🚀 Starting card due notifications job: %s", jobRunID)

	// 1. Obtener tarjetas que vencen mañana
	cards, err := s.cardRepo.GetCardsDueTomorrow()
	if err != nil {
		jobRun.Status = "failed"
		jobRun.ErrorMessage = err.Error()
		jobRun.CompletedAt = &[]time.Time{time.Now()}[0]
		s.notificationRepo.UpdateJobRun(jobRun)
		return fmt.Errorf("error getting cards due tomorrow: %w", err)
	}

	jobRun.CardsFound = len(cards)
	log.Printf("📅 Found %d cards due tomorrow", len(cards))

	if len(cards) == 0 {
		jobRun.Status = "completed"
		jobRun.CompletedAt = &[]time.Time{time.Now()}[0]
		s.notificationRepo.UpdateJobRun(jobRun)
		log.Printf("✅ Job completed: No cards due tomorrow")
		return nil
	}

	// 2. Procesar cada tarjeta
	emailsSent := 0
	errors := 0

	for _, card := range cards {
		if err := s.processCardNotification(card, jobRunID); err != nil {
			log.Printf("❌ Error processing card %s (%s): %v", card.CardName, card.ID, err)
			errors++
		} else {
			log.Printf("✅ Notification sent for card %s (%s) to %s", card.CardName, card.ID, card.UserEmail)
			emailsSent++
		}
	}

	// 3. Actualizar el job run con los resultados
	jobRun.EmailsSent = emailsSent
	jobRun.Errors = errors
	jobRun.Status = "completed"
	if errors > 0 && emailsSent == 0 {
		jobRun.Status = "failed"
		jobRun.ErrorMessage = fmt.Sprintf("All %d notifications failed", errors)
	}
	jobRun.CompletedAt = &[]time.Time{time.Now()}[0]

	if err := s.notificationRepo.UpdateJobRun(jobRun); err != nil {
		log.Printf("Error updating job run: %v", err)
	}

	log.Printf("🎉 Job completed: %d emails sent, %d errors", emailsSent, errors)
	return nil
}

// processCardNotification procesa la notificación para una tarjeta específica
func (s *NotificationService) processCardNotification(card *entities.Card, jobRunID string) error {
	notificationLog := &entities.NotificationLog{
		ID:       uuid.New().String(),
		JobRunID: jobRunID,
		CardID:   card.ID,
		UserID:   card.UserID,
		Email:    card.UserEmail,
		SentAt:   time.Now(),
	}

	// 1. Obtener cuotas pendientes hasta la fecha de vencimiento
	installments, err := s.installmentRepo.GetPendingInstallmentsByCard(card.ID, card.DueDate)
	if err != nil {
		notificationLog.Status = "failed"
		notificationLog.ErrorMessage = fmt.Sprintf("Error getting installments: %v", err)
		s.notificationRepo.SaveNotificationLog(notificationLog)
		return fmt.Errorf("error getting installments for card %s: %w", card.ID, err)
	}

	// 2. Calcular el total y construir los detalles
	notification := s.buildCardDueNotification(card, installments)

	// 3. Mostrar información sobre cuotas pendientes pero siempre enviar email
	if notification.PendingInstallments == 0 {
		log.Printf("📧 Sending notification for card %s: no pending installments (informational email)", card.CardName)
	} else {
		log.Printf("📧 Sending notification for card %s: %d pending installments", card.CardName, notification.PendingInstallments)
	}

	// 4. Enviar email (siempre, independientemente de si hay cuotas pendientes)
	if err := s.emailService.SendCardDueNotification(notification); err != nil {
		notificationLog.Status = "failed"
		notificationLog.ErrorMessage = fmt.Sprintf("Error sending email: %v", err)
		s.notificationRepo.SaveNotificationLog(notificationLog)
		return fmt.Errorf("error sending email for card %s: %w", card.ID, err)
	}

	// 5. Guardar log exitoso
	notificationLog.Status = "sent"
	s.notificationRepo.SaveNotificationLog(notificationLog)

	return nil
}

// buildCardDueNotification construye la notificación a partir de la tarjeta y cuotas
func (s *NotificationService) buildCardDueNotification(card *entities.Card, installments []*entities.Installment) *entities.CardDueNotification {
	notification := &entities.CardDueNotification{
		CardID:              card.ID,
		UserID:              card.UserID,
		UserEmail:           card.UserEmail,
		UserName:            card.GetUserFullName(),
		CardName:            card.CardName,
		BankName:            card.BankName,
		LastFour:            card.LastFour,
		DueDate:             card.DueDate,
		PendingInstallments: len(installments),
		InstallmentDetails:  make([]entities.InstallmentSummary, 0, len(installments)),
	}

	// Calcular total y construir detalles
	totalAmount := 0.0
	for _, installment := range installments {
		totalAmount += installment.Amount

		summary := entities.InstallmentSummary{
			Description:       installment.Description,
			MerchantName:      installment.MerchantName,
			Amount:            installment.Amount,
			DueDate:           installment.DueDate,
			InstallmentNum:    installment.InstallmentNum,
			TotalInstallments: installment.TotalInstallments,
		}

		notification.InstallmentDetails = append(notification.InstallmentDetails, summary)
	}

	notification.TotalPendingAmount = totalAmount

	return notification
}

// TriggerManualJob ejecuta manualmente el job de notificaciones
func (s *NotificationService) TriggerManualJob() error {
	log.Printf("🔧 Manual trigger for card due notifications job")
	return s.ProcessCardDueNotifications()
}

// UpdateExpiredDueDates actualiza las fechas de vencimiento de tarjetas que vencieron ayer
func (s *NotificationService) UpdateExpiredDueDates() error {
	log.Printf("🔧 Starting update expired due dates job")

	startTime := time.Now()
	cardsUpdated, err := s.cardRepo.UpdateExpiredDueDates()
	if err != nil {
		log.Printf("❌ Error updating expired due dates: %v", err)
		return fmt.Errorf("failed to update expired due dates: %w", err)
	}

	duration := time.Since(startTime)
	if cardsUpdated > 0 {
		log.Printf("✅ Updated %d cards with expired due dates in %v", cardsUpdated, duration)
	} else {
		log.Printf("ℹ️  No cards found with expired due dates (yesterday). Duration: %v", duration)
	}

	return nil
}

// GetJobHistory obtiene el historial de ejecuciones del job
func (s *NotificationService) GetJobHistory(limit int) ([]*entities.JobRun, error) {
	return s.notificationRepo.GetJobRunHistory(limit)
}

// GetNotificationLogs obtiene los logs de notificaciones para un job run
func (s *NotificationService) GetNotificationLogs(jobRunID string, limit int) ([]*entities.NotificationLog, error) {
	return s.notificationRepo.GetNotificationLogs(jobRunID, limit)
}

// SendSupportEmail envía un email de soporte desde el FAQ
func (s *NotificationService) SendSupportEmail(name, email, subject, message string) error {
	log.Printf("📧 Sending support email from %s (%s): %s", name, email, subject)

	if err := s.emailService.SendSupportEmail(name, email, subject, message); err != nil {
		log.Printf("❌ Error sending support email: %v", err)
		return fmt.Errorf("failed to send support email: %w", err)
	}

	log.Printf("✅ Support email sent successfully")
	return nil
}
