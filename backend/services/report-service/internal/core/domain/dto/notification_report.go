package dto

import "time"

// NotificationReportRequest request para reporte de notificaciones
type NotificationReportRequest struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// NotificationReportResponse respuesta del reporte de notificaciones
type NotificationReportResponse struct {
	Period   Period                 `json:"period"`
	Summary  NotificationSummary    `json:"summary"`
	ByDay    []NotificationByDay    `json:"by_day"`
	ByStatus []NotificationByStatus `json:"by_status"`
	JobRuns  []JobRunDetail         `json:"job_runs"`
}

// NotificationSummary resumen de notificaciones
type NotificationSummary struct {
	TotalNotifications int     `json:"total_notifications"`
	TotalJobRuns       int     `json:"total_job_runs"`
	SuccessfulSent     int     `json:"successful_sent"`
	Failed             int     `json:"failed"`
	SuccessRate        float64 `json:"success_rate"`
	FailureRate        float64 `json:"failure_rate"`
	AvgEmailsPerRun    float64 `json:"avg_emails_per_run"`
}

// NotificationByDay notificaciones por día
type NotificationByDay struct {
	Date        time.Time `json:"date"`
	Day         string    `json:"day"`
	Sent        int       `json:"sent"`
	Failed      int       `json:"failed"`
	SuccessRate float64   `json:"success_rate"`
}

// NotificationByStatus notificaciones por estado
type NotificationByStatus struct {
	Status     string  `json:"status"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

// JobRunDetail detalle de ejecución de job
type JobRunDetail struct {
	ID           string     `json:"id"`
	StartedAt    time.Time  `json:"started_at"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	Status       string     `json:"status"`
	CardsFound   int        `json:"cards_found"`
	EmailsSent   int        `json:"emails_sent"`
	Errors       int        `json:"errors"`
	Duration     string     `json:"duration,omitempty"`
	ErrorMessage string     `json:"error_message,omitempty"`
}
