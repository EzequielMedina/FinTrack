package notification

import (
	"net/http"
	"strconv"
	"time"

	"github.com/fintrack/notification-service/internal/core/ports"
	"github.com/fintrack/notification-service/internal/infrastructure/jobs"
	"github.com/gin-gonic/gin"
)

// Handler maneja las solicitudes HTTP para notificaciones
type Handler struct {
	notificationService ports.NotificationService
	jobScheduler        *jobs.JobScheduler
}

// New crea un nuevo handler de notificaciones
func New(notificationService ports.NotificationService, jobScheduler *jobs.JobScheduler) *Handler {
	return &Handler{
		notificationService: notificationService,
		jobScheduler:        jobScheduler,
	}
}

// JobHistoryResponse representa la respuesta del historial de jobs
type JobHistoryResponse struct {
	RunID        string     `json:"run_id"`
	StartedAt    time.Time  `json:"started_at"`
	CompletedAt  *time.Time `json:"completed_at"`
	Status       string     `json:"status"`
	CardsFound   int        `json:"cards_found"`
	EmailsSent   int        `json:"emails_sent"`
	Errors       int        `json:"errors"`
	ErrorMessage string     `json:"error_message,omitempty"`
	Duration     *string    `json:"duration,omitempty"`
}

// NotificationLogResponse representa la respuesta de logs de notificación
type NotificationLogResponse struct {
	ID           string    `json:"id"`
	CardID       string    `json:"card_id"`
	UserID       string    `json:"user_id"`
	Email        string    `json:"email"`
	Status       string    `json:"status"`
	ErrorMessage string    `json:"error_message,omitempty"`
	SentAt       time.Time `json:"sent_at"`
}

// HealthResponse representa la respuesta de health check
type HealthResponse struct {
	Status           string     `json:"status"`
	Timestamp        time.Time  `json:"timestamp"`
	Service          string     `json:"service"`
	Version          string     `json:"version"`
	JobScheduler     bool       `json:"job_scheduler_running"`
	NextScheduledRun *time.Time `json:"next_scheduled_run,omitempty"`
}

// TriggerCardDueJob ejecuta manualmente el job de notificaciones de tarjetas
// POST /api/notifications/trigger-card-due-job
func (h *Handler) TriggerCardDueJob(c *gin.Context) {
	// Ejecutar el job en una goroutine para no bloquear la respuesta
	go func() {
		if err := h.jobScheduler.TriggerCardDueJob(); err != nil {
			// Log error pero no retornar al cliente ya que es asíncrono
			// El error se registrará en los logs del job
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"message":   "Card due notifications job triggered successfully",
		"timestamp": time.Now(),
		"async":     true,
	})
}

// GetJobHistory obtiene el historial de ejecuciones del job
// GET /api/notifications/job-history?limit=10
func (h *Handler) GetJobHistory(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 20
	}

	jobRuns, err := h.notificationService.GetJobHistory(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get job history",
			"details": err.Error(),
		})
		return
	}

	// Convertir a respuesta
	response := make([]JobHistoryResponse, 0, len(jobRuns))
	for _, job := range jobRuns {
		jobResponse := JobHistoryResponse{
			RunID:        job.ID,
			StartedAt:    job.StartedAt,
			CompletedAt:  job.CompletedAt,
			Status:       job.Status,
			CardsFound:   job.CardsFound,
			EmailsSent:   job.EmailsSent,
			Errors:       job.Errors,
			ErrorMessage: job.ErrorMessage,
		}

		// Calcular duración si está completado
		if job.CompletedAt != nil {
			duration := job.CompletedAt.Sub(job.StartedAt)
			durationStr := duration.String()
			jobResponse.Duration = &durationStr
		}

		response = append(response, jobResponse)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  response,
		"count": len(response),
		"limit": limit,
	})
}

// GetNotificationLogs obtiene los logs de notificaciones para un job run
// GET /api/notifications/logs?job_run_id=123&limit=10
func (h *Handler) GetNotificationLogs(c *gin.Context) {
	jobRunID := c.Query("job_run_id")
	if jobRunID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "job_run_id parameter is required",
		})
		return
	}

	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 50
	}

	logs, err := h.notificationService.GetNotificationLogs(jobRunID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get notification logs",
			"details": err.Error(),
		})
		return
	}

	// Convertir a respuesta
	response := make([]NotificationLogResponse, 0, len(logs))
	for _, log := range logs {
		response = append(response, NotificationLogResponse{
			ID:           log.ID,
			CardID:       log.CardID,
			UserID:       log.UserID,
			Email:        log.Email,
			Status:       log.Status,
			ErrorMessage: log.ErrorMessage,
			SentAt:       log.SentAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       response,
		"count":      len(response),
		"job_run_id": jobRunID,
		"limit":      limit,
	})
}

// HealthCheck verifica el estado del servicio
// GET /api/notifications/health
func (h *Handler) HealthCheck(c *gin.Context) {
	response := HealthResponse{
		Status:       "healthy",
		Timestamp:    time.Now(),
		Service:      "notification-service",
		Version:      "1.0.0",
		JobScheduler: h.jobScheduler.IsRunning(),
	}

	// Obtener próxima ejecución programada
	if response.JobScheduler {
		nextRun := h.jobScheduler.GetNextScheduledRun()
		if !nextRun.IsZero() {
			response.NextScheduledRun = &nextRun
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetSchedulerStatus obtiene el estado del job scheduler
// GET /api/notifications/scheduler/status
func (h *Handler) GetSchedulerStatus(c *gin.Context) {
	isRunning := h.jobScheduler.IsRunning()
	nextRun := h.jobScheduler.GetNextScheduledRun()

	response := gin.H{
		"scheduler_running": isRunning,
		"status":            "active",
	}

	if isRunning && !nextRun.IsZero() {
		response["next_scheduled_run"] = nextRun
		response["time_to_next_run"] = time.Until(nextRun).String()
	}

	c.JSON(http.StatusOK, response)
}
