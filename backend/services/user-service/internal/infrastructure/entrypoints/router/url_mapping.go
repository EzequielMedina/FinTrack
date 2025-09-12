package router

import (
	"net/http"

	"github.com/fintrack/user-service/internal/config"
	authmw "github.com/fintrack/user-service/internal/infrastructure/entrypoints/middleware"
	"github.com/gin-gonic/gin"
)

func MapRoutes(r *gin.Engine, h *Handlers, cfg *config.Config) {
	// Basic health endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", h.Auth.Register)
			auth.POST("/login", h.Auth.Login)
		}

		// Example of a protected route
		protected := api.Group("/me")
		protected.Use(authmw.AuthJWT(cfg.JWTSecret))
		{
			protected.GET("", func(c *gin.Context) {
				// token validated by middleware; in a real handler we'd fetch user info
				c.JSON(http.StatusOK, gin.H{"message": "authenticated"})
			})
		}
	}
}
