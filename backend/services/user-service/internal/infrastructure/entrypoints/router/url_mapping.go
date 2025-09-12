package router

import (
	"net/http"

	"github.com/fintrack/user-service/internal/app"
	"github.com/fintrack/user-service/internal/config"
	authmw "github.com/fintrack/user-service/internal/infrastructure/entrypoints/middleware"
	mysqlrepo "github.com/fintrack/user-service/internal/infrastructure/repositories/mysql"
	"github.com/gin-gonic/gin"
)

func MapRoutes(r *gin.Engine, h *Handlers, cfg *config.Config, application *app.Application) {
	// Basic health endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	{
		// Authentication routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", h.Auth.Register)
			auth.POST("/login", h.Auth.Login)
		}

		// Get user repository from application
		userRepo := mysqlrepo.NewUserRepository(application.DB)

		// Protected routes requiring authentication with full user context
		protected := api.Group("")
		protected.Use(authmw.AuthWithUser(cfg.JWTSecret, userRepo))
		{
			// Current user routes
			protected.GET("/me", h.User.GetCurrentUser)

			// User management routes
			users := protected.Group("/users")
			{
				// Basic CRUD operations
				users.POST("", h.User.CreateUser)       // POST /api/users
				users.GET("", h.User.GetAllUsers)       // GET /api/users?page=1&pageSize=20
				users.GET("/:id", h.User.GetUser)       // GET /api/users/:id
				users.PUT("/:id", h.User.UpdateUser)    // PUT /api/users/:id
				users.DELETE("/:id", h.User.DeleteUser) // DELETE /api/users/:id

				// Profile management
				users.PUT("/:id/profile", h.User.UpdateUserProfile) // PUT /api/users/:id/profile

				// Role and status management
				users.PUT("/:id/role", h.User.ChangeUserRole)     // PUT /api/users/:id/role
				users.PUT("/:id/status", h.User.ToggleUserStatus) // PUT /api/users/:id/status
				users.PUT("/:id/password", h.User.ChangePassword) // PUT /api/users/:id/password

				// Query by role
				users.GET("/role/:role", h.User.GetUsersByRole) // GET /api/users/role/:role
			}
		}

		// Legacy protected route using basic JWT middleware
		legacyProtected := api.Group("/legacy")
		legacyProtected.Use(authmw.AuthJWT(cfg.JWTSecret))
		{
			legacyProtected.GET("/me", func(c *gin.Context) {
				// token validated by middleware; in a real handler we'd fetch user info
				c.JSON(http.StatusOK, gin.H{"message": "authenticated"})
			})
		}
	}
}
