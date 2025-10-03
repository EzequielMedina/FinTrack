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
			accounts.POST("", h.Account.CreateAccount) // POST /api/accounts
			accounts.GET("", h.Account.GetAccounts)    // GET /api/accounts?page=1&pageSize=20

			// User-specific routes (debe ir antes de /:id para evitar conflictos)
			accounts.GET("/user/:userId", h.Account.GetAccountsByUser) // GET /api/accounts/user/:userId
			accounts.GET("/user/:userId/cards", h.Card.GetCardsByUser) // GET /api/accounts/user/:userId/cards

			// Account-specific routes con :id
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

			// Card management routes (usar :id ya que no hay conflicto con esta estructura)
			accounts.POST("/:id/cards", h.Card.CreateCard)           // POST /api/accounts/:id/cards
			accounts.GET("/:id/cards", h.Card.GetCardsByAccount)     // GET /api/accounts/:id/cards
			accounts.GET("/:id/cards/:cardId", h.Card.GetCard)       // GET /api/accounts/:id/cards/:cardId
			accounts.PUT("/:id/cards/:cardId", h.Card.UpdateCard)    // PUT /api/accounts/:id/cards/:cardId
			accounts.DELETE("/:id/cards/:cardId", h.Card.DeleteCard) // DELETE /api/accounts/:id/cards/:cardId

			// Card status management
			accounts.PUT("/:id/cards/:cardId/block", h.Card.BlockCard)            // PUT /api/accounts/:id/cards/:cardId/block
			accounts.PUT("/:id/cards/:cardId/unblock", h.Card.UnblockCard)        // PUT /api/accounts/:id/cards/:cardId/unblock
			accounts.PUT("/:id/cards/:cardId/set-default", h.Card.SetDefaultCard) // PUT /api/accounts/:id/cards/:cardId/set-default
		}

		// Direct card operations (financial transactions)
		cards := api.Group("/cards")
		{
			// Card balance queries
			cards.GET("/:cardId/balance", h.Card.GetCardBalance) // GET /api/cards/:cardId/balance

			// Credit card operations
			cards.POST("/:cardId/charge", h.Card.ChargeCard)   // POST /api/cards/:cardId/charge
			cards.POST("/:cardId/payment", h.Card.PaymentCard) // POST /api/cards/:cardId/payment

			// Debit card operations
			cards.POST("/:cardId/transaction", h.Card.ProcessDebitTransaction) // POST /api/cards/:cardId/transaction

			// Installment operations
			cards.POST("/:cardId/installments/preview", h.Installment.PreviewInstallmentPlan)    // POST /api/cards/:cardId/installments/preview
			cards.POST("/:cardId/charge-installments", h.Installment.ChargeCardWithInstallments) // POST /api/cards/:cardId/charge-installments
			cards.GET("/:cardId/installment-plans", h.Installment.GetInstallmentPlansByCard)     // GET /api/cards/:cardId/installment-plans
		}

		// Direct installment operations
		installments := api.Group("/installment-plans")
		{
			installments.GET("", h.Installment.GetUserInstallmentPlans)               // GET /api/installment-plans
			installments.GET("/:planId", h.Installment.GetInstallmentPlan)            // GET /api/installment-plans/:planId
			installments.POST("/:planId/cancel", h.Installment.CancelInstallmentPlan) // POST /api/installment-plans/:planId/cancel
		}

		// Individual installment operations
		installmentItems := api.Group("/installments")
		{
			installmentItems.POST("/:installmentId/pay", h.Installment.PayInstallment)     // POST /api/installments/:installmentId/pay
			installmentItems.GET("/overdue", h.Installment.GetOverdueInstallments)         // GET /api/installments/overdue
			installmentItems.GET("/upcoming", h.Installment.GetUpcomingInstallments)       // GET /api/installments/upcoming
			installmentItems.GET("/summary", h.Installment.GetInstallmentSummary)          // GET /api/installments/summary
			installmentItems.GET("/monthly-load", h.Installment.GetMonthlyInstallmentLoad) // GET /api/installments/monthly-load
		}
	}
}
