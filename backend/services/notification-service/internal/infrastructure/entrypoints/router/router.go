package router

import (
	"github.com/fintrack/notification-service/internal/infrastructure/entrypoints/handlers/notification"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configura todas las rutas del servicio
func SetupRoutes(notificationHandler *notification.Handler) *gin.Engine {
	router := gin.New()

	// Middlewares globales
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Grupo de rutas API
	api := router.Group("/api/notifications")
	{
		// Job management
		api.POST("/trigger-card-due-job", notificationHandler.TriggerCardDueJob)
		api.POST("/trigger-update-due-dates-job", notificationHandler.TriggerUpdateDueDatesJob)
		api.GET("/job-history", notificationHandler.GetJobHistory)

		// Notification logs
		api.GET("/logs", notificationHandler.GetNotificationLogs)

		// Scheduler status
		api.GET("/scheduler/status", notificationHandler.GetSchedulerStatus)

		// Health check
		api.GET("/health", notificationHandler.HealthCheck)
	}

	// Root health check
	router.GET("/health", notificationHandler.HealthCheck)

	// Root info endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "notification-service",
			"version": "1.0.0",
			"status":  "running",
			"endpoints": []string{
				"GET /health",
				"POST /api/notifications/trigger-card-due-job",
				"POST /api/notifications/trigger-update-due-dates-job",
				"GET /api/notifications/job-history",
				"GET /api/notifications/logs",
				"GET /api/notifications/scheduler/status",
				"GET /api/notifications/health",
			},
		})
	})

	return router
}
