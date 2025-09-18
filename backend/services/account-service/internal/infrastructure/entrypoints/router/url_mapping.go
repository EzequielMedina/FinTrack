package router

import (
	"github.com/fintrack/account-service/internal/app"
	"github.com/fintrack/account-service/internal/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func MapRoutes(r *gin.Engine, h *Handlers, cfg *config.Config, application *app.Application) {
	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // En producción, especificar dominios exactos
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		// Account management routes
		accounts := api.Group("/accounts")
		{
			// Basic CRUD operations
			accounts.POST("", h.Account.CreateAccount)       // POST /api/accounts
			accounts.GET("", h.Account.GetAccounts)          // GET /api/accounts?page=1&pageSize=20
			accounts.GET("/:id", h.Account.GetAccount)       // GET /api/accounts/:id
			accounts.PUT("/:id", h.Account.UpdateAccount)    // PUT /api/accounts/:id
			accounts.DELETE("/:id", h.Account.DeleteAccount) // DELETE /api/accounts/:id

			// Balance operations
			accounts.PUT("/:id/balance", h.Account.UpdateBalance) // PUT /api/accounts/:id/balance
			accounts.GET("/:id/balance", h.Account.GetBalance)    // GET /api/accounts/:id/balance

			// Status management
			accounts.PUT("/:id/status", h.Account.UpdateStatus) // PUT /api/accounts/:id/status

			// Wallet operations
			accounts.POST("/:id/add-funds", h.Account.AddFunds)           // POST /api/accounts/:id/add-funds
			accounts.POST("/:id/withdraw-funds", h.Account.WithdrawFunds) // POST /api/accounts/:id/withdraw-funds

			// Credit card operations
			accounts.PUT("/:id/credit-limit", h.Account.UpdateCreditLimit)      // PUT /api/accounts/:id/credit-limit
			accounts.PUT("/:id/credit-dates", h.Account.UpdateCreditDates)      // PUT /api/accounts/:id/credit-dates
			accounts.GET("/:id/available-credit", h.Account.GetAvailableCredit) // GET /api/accounts/:id/available-credit

			// User-specific routes
			accounts.GET("/user/:userId", h.Account.GetAccountsByUser) // GET /api/accounts/user/:userId
		}
	}
}
