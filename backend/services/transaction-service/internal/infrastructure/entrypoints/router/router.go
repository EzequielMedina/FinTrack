package router

import (
	"database/sql"
	"net/http"
)

// Router handles all HTTP routing for the transaction service
type Router struct {
	handler     *TransactionHandler
	cardHandler *CardHandler
}

// NewRouter creates a new router instance
func NewRouter(db *sql.DB) *Router {
	// Create handlers
	transactionHandler := NewTransactionHandler(db)
	cardHandler := NewCardHandler(db)

	router := &Router{
		handler:     transactionHandler,
		cardHandler: cardHandler,
	}

	return router
}

// SetupRoutes configures all routes using standard net/http
func (r *Router) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", r.healthCheck)

	// Transaction routes - using pattern matching (Go 1.22+)
	mux.HandleFunc("POST /api/v1/transactions", r.handler.CreateTransactionHTTP)
	mux.HandleFunc("GET /api/v1/transactions", r.handler.ListTransactionsHTTP)
	mux.HandleFunc("GET /api/v1/transactions/{id}", r.handler.GetTransactionHTTP)
	mux.HandleFunc("PUT /api/v1/transactions/{id}/status", r.handler.UpdateTransactionStatusHTTP)
	mux.HandleFunc("POST /api/v1/transactions/{id}/process", r.handler.ProcessTransactionHTTP)
	mux.HandleFunc("POST /api/v1/transactions/{id}/reverse", r.handler.ReverseTransactionHTTP)

	// Card transaction routes
	mux.HandleFunc("POST /api/v1/cards/credit/charge", r.cardHandler.ChargeCreditCardHTTP)
	mux.HandleFunc("POST /api/v1/cards/credit/payment", r.cardHandler.PayCreditCardHTTP)
	mux.HandleFunc("POST /api/v1/cards/debit/transaction", r.cardHandler.ProcessDebitCardTransactionHTTP)

	return mux
}

// healthCheck handles health check requests
func (r *Router) healthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy","service":"transaction-service","version":"1.0.0"}`))
}
