package jobs

import (
	"log"
	"time"

	"github.com/fintrack/notification-service/internal/config"
	"github.com/fintrack/notification-service/internal/core/ports"
	"github.com/robfig/cron/v3"
)

// JobScheduler maneja la programación y ejecución de jobs
type JobScheduler struct {
	notificationService ports.NotificationService
	cron                *cron.Cron
	config              *config.JobConfig
}

// NewJobScheduler crea un nuevo scheduler de jobs
func NewJobScheduler(notificationService ports.NotificationService, cfg *config.JobConfig) *JobScheduler {
	// Configurar timezone
	location, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		log.Printf("⚠️  Warning: Invalid timezone %s, using UTC", cfg.Timezone)
		location = time.UTC
	}

	return &JobScheduler{
		notificationService: notificationService,
		cron:                cron.New(cron.WithLocation(location)),
		config:              cfg,
	}
}

// Start inicia el scheduler de jobs
func (j *JobScheduler) Start() error {
	if !j.config.Enabled {
		log.Println("📴 Job scheduler is disabled")
		return nil
	}

	log.Printf("⏰ Starting job scheduler with cron: %s (timezone: %s)", j.config.Schedule, j.config.Timezone)

	// Programar el job de notificaciones de tarjetas
	_, err := j.cron.AddFunc(j.config.Schedule, func() {
		log.Println("🔔 Starting scheduled card due notifications job")

		startTime := time.Now()
		if err := j.notificationService.ProcessCardDueNotifications(); err != nil {
			log.Printf("❌ Scheduled job failed: %v", err)
		} else {
			duration := time.Since(startTime)
			log.Printf("✅ Scheduled job completed successfully in %v", duration)
		}
	})

	if err != nil {
		return err
	}

	// Iniciar el cron scheduler
	j.cron.Start()
	log.Printf("🚀 Job scheduler started - next run will be according to cron: %s", j.config.Schedule)

	return nil
}

// Stop detiene el scheduler de jobs
func (j *JobScheduler) Stop() {
	if j.cron != nil {
		log.Println("⏹️  Stopping job scheduler...")
		ctx := j.cron.Stop()
		<-ctx.Done()
		log.Println("✅ Job scheduler stopped")
	}
}

// TriggerCardDueJob ejecuta manualmente el job de notificaciones de tarjetas
func (j *JobScheduler) TriggerCardDueJob() error {
	log.Println("🔧 Manual trigger for card due notifications job")
	return j.notificationService.TriggerManualJob()
}

// GetNextScheduledRun obtiene la próxima ejecución programada
func (j *JobScheduler) GetNextScheduledRun() time.Time {
	entries := j.cron.Entries()
	if len(entries) > 0 {
		return entries[0].Next
	}
	return time.Time{}
}

// IsRunning verifica si el scheduler está en ejecución
func (j *JobScheduler) IsRunning() bool {
	return j.cron != nil && len(j.cron.Entries()) > 0
}
