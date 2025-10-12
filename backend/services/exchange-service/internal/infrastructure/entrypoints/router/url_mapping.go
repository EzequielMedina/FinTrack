package router

import (
	"github.com/gin-gonic/gin"

	exchangehandler "github.com/fintrack/exchange-service/internal/infrastructure/entrypoints/handlers/exchange"
)

// SetupRoutes configura todas las rutas del servicio
func SetupRoutes(exchangeHandler *exchangehandler.ExchangeHandler) *gin.Engine {
	router := gin.Default()

	// Middleware global
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(CORSMiddleware())

	// Health check endpoint
	router.GET("/health", exchangeHandler.HealthCheck)

	// API routes
	apiGroup := router.Group("/api")
	{
		exchangeGroup := apiGroup.Group("/exchange")
		{
			// Endpoint principal para obtener dólar oficial
			exchangeGroup.GET("/dolar-oficial", exchangeHandler.GetDolarOficial)
		}
	}

	return router
}

// CORSMiddleware configura CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
