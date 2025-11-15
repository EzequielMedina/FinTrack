package router

import (
	"github.com/fintrack/report-service/internal/config"
	"github.com/fintrack/report-service/internal/core/service"
	"github.com/fintrack/report-service/internal/infrastructure/entrypoints/handlers/report"
	"github.com/gin-gonic/gin"
)

// SetupRouter configura las rutas del servicio
func SetupRouter(cfg *config.Config, reportService service.ReportService) *gin.Engine {
	// Configurar modo de Gin
	gin.SetMode(cfg.Server.Environment)

	r := gin.Default()

	// Configurar CORS
	r.Use(corsMiddleware(cfg.CORS.AllowedOrigins))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "report-service",
		})
	})

	// API v1
	v1 := r.Group("/api/v1")
	{
		// Crear handlers
		reportHandler := report.NewReportHandler(reportService)

		// Rutas de reportes
		reports := v1.Group("/reports")
		{
			// Reportes JSON
			reports.GET("/transactions", reportHandler.GetTransactionReport)
			reports.GET("/installments", reportHandler.GetInstallmentReport)
			reports.GET("/accounts", reportHandler.GetAccountReport)
			reports.GET("/expenses-income", reportHandler.GetExpenseIncomeReport)
			reports.GET("/notifications", reportHandler.GetNotificationReport)

			// Reportes PDF
			reports.GET("/transactions/pdf", reportHandler.GetTransactionReportPDF)
			reports.GET("/installments/pdf", reportHandler.GetInstallmentReportPDF)
			reports.GET("/accounts/pdf", reportHandler.GetAccountReportPDF)
			reports.GET("/expenses-income/pdf", reportHandler.GetExpenseIncomeReportPDF)
		}
	}

	return r
}

// corsMiddleware configura CORS
func corsMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Verificar si el origen está permitido
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin || allowedOrigin == "*" {
				allowed = true
				break
			}
		}

		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
