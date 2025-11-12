package router

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	domaintransaction "github.com/fintrack/transaction-service/internal/core/domain/entities/transaction"
	"github.com/fintrack/transaction-service/internal/core/service"
	"github.com/fintrack/transaction-service/internal/infrastructure/http/clients"
	"github.com/fintrack/transaction-service/internal/infrastructure/repositories/mysql"
)

// TransactionHandler handles HTTP requests for transaction operations
type TransactionHandler struct {
	transactionService service.TransactionServiceInterface
}

// NewTransactionHandler creates a new transaction handler
func NewTransactionHandler(db *sql.DB) *TransactionHandler {
	// Create repositories
	transactionRepo := mysql.NewTransactionRepository(db)

	// Create mock services for now
	ruleService := service.NewTransactionRuleService(nil)   // TODO: Implement rule repository
	auditService := service.NewTransactionAuditService(nil) // TODO: Implement audit repository
	externalService := service.NewMockExternalService()

	// Create account service client - using environment variable or default localhost
	accountServiceURL := "http://localhost:8082" // Correct account service port
	if url := os.Getenv("ACCOUNT_SERVICE_URL"); url != "" {
		accountServiceURL = url
	}
	accountService := clients.NewAccountClient(accountServiceURL)

	// Create transaction service
	transactionService := service.NewTransactionService(
		transactionRepo,
		ruleService,
		auditService,
		externalService,
		accountService,
	)

	return &TransactionHandler{
		transactionService: transactionService,
	}
}

// DTOs for request/response

// CreateTransactionRequest represents the request to create a transaction
type CreateTransactionRequest struct {
	Type          string                 `json:"type"`
	Amount        float64                `json:"amount"`
	Currency      string                 `json:"currency"`
	FromAccountID *string                `json:"fromAccountId"`
	ToAccountID   *string                `json:"toAccountId"`
	FromCardID    *string                `json:"fromCardId"`
	ToCardID      *string                `json:"toCardId"`
	Description   string                 `json:"description"`
	PaymentMethod string                 `json:"paymentMethod"`
	MerchantName  string                 `json:"merchantName"`
	MerchantID    string                 `json:"merchantId"`
	ReferenceID   string                 `json:"referenceId"`
	ExternalID    string                 `json:"externalId"`
	Metadata      map[string]interface{} `json:"metadata"`
	Tags          []string               `json:"tags"`
}

// TransactionResponse represents the response for transaction operations
type TransactionResponse struct {
	ID              string                 `json:"id"`
	ReferenceID     string                 `json:"referenceId"`
	ExternalID      string                 `json:"externalId"`
	Type            string                 `json:"type"`
	Status          string                 `json:"status"`
	Amount          float64                `json:"amount"`
	Currency        string                 `json:"currency"`
	FromAccountID   *string                `json:"fromAccountId"`
	ToAccountID     *string                `json:"toAccountId"`
	FromCardID      *string                `json:"fromCardId"`
	ToCardID        *string                `json:"toCardId"`
	UserID          string                 `json:"userId"`
	InitiatedBy     string                 `json:"initiatedBy"`
	Description     string                 `json:"description"`
	PaymentMethod   string                 `json:"paymentMethod"`
	MerchantName    string                 `json:"merchantName"`
	MerchantID      string                 `json:"merchantId"`
	PreviousBalance float64                `json:"previousBalance"`
	NewBalance      float64                `json:"newBalance"`
	ProcessedAt     *string                `json:"processedAt"`
	FailedAt        *string                `json:"failedAt"`
	FailureReason   string                 `json:"failureReason"`
	Metadata        map[string]interface{} `json:"metadata"`
	Tags            []string               `json:"tags"`
	CreatedAt       string                 `json:"createdAt"`
	UpdatedAt       string                 `json:"updatedAt"`
}

// TransactionListResponse represents the response for listing transactions
type TransactionListResponse struct {
	Transactions []*TransactionResponse `json:"transactions"`
	Total        int                    `json:"total"`
	Page         int                    `json:"page"`
	PageSize     int                    `json:"pageSize"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// HTTP Handler methods using standard net/http

// CreateTransactionHTTP creates a new transaction
func (h *TransactionHandler) CreateTransactionHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("📥 CreateTransactionHTTP called")

	if r.Method != http.MethodPost {
		log.Println("❌ Method not allowed")
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", "Only POST method is allowed")
		return
	}

	var req CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("❌ Failed to decode request: %v\n", err)
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}
	log.Printf("✅ Request decoded: type=%s, amount=%.2f\n", req.Type, req.Amount)

	// Basic validation
	if req.Type == "" {
		log.Println("❌ Transaction type is required")
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request", "Transaction type is required")
		return
	}

	if req.Amount <= 0 {
		log.Println("❌ Amount must be greater than 0")
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request", "Amount must be greater than 0")
		return
	}

	// Get user ID from header (simplified for now)
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		log.Println("❌ User ID is required")
		h.writeErrorResponse(w, http.StatusUnauthorized, "Unauthorized", "User ID is required")
		return
	}
	log.Printf("✅ UserID from header: %s\n", userID)

	// Convert to service request
	serviceReq := service.CreateTransactionRequest{
		UserID:        userID,
		Type:          domaintransaction.TransactionType(req.Type),
		Amount:        req.Amount,
		Currency:      req.Currency,
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		FromCardID:    req.FromCardID,
		ToCardID:      req.ToCardID,
		Description:   req.Description,
		PaymentMethod: domaintransaction.PaymentMethod(req.PaymentMethod),
		MerchantName:  req.MerchantName,
		MerchantID:    req.MerchantID,
		ReferenceID:   req.ReferenceID,
		ExternalID:    req.ExternalID,
		Metadata:      req.Metadata,
		Tags:          req.Tags,
	}

	// Set default currency
	if serviceReq.Currency == "" {
		serviceReq.Currency = "USD"
	}

	log.Printf("🔄 Calling transactionService.CreateTransaction for type=%s, amount=%.2f\n", serviceReq.Type, serviceReq.Amount)

	// Create transaction
	transaction, err := h.transactionService.CreateTransaction(serviceReq, userID)
	if err != nil {
		log.Printf("❌ CreateTransaction failed: %v\n", err)
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to create transaction", err.Error())
		return
	}

	log.Printf("✅ Transaction created successfully: %s\n", transaction.ID)

	// Convert to response
	response := h.toTransactionResponse(transaction)
	h.writeJSONResponse(w, http.StatusCreated, response)
}

// GetTransactionHTTP retrieves a transaction by ID
func (h *TransactionHandler) GetTransactionHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", "Only GET method is allowed")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/v1/transactions/")
	if id == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request", "Transaction ID is required")
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		h.writeErrorResponse(w, http.StatusUnauthorized, "Unauthorized", "User ID is required")
		return
	}

	transaction, err := h.transactionService.GetTransactionByID(id, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.writeErrorResponse(w, http.StatusNotFound, "Transaction not found", err.Error())
			return
		}

		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to get transaction", err.Error())
		return
	}

	response := h.toTransactionResponse(transaction)
	h.writeJSONResponse(w, http.StatusOK, response)
}

// ListTransactionsHTTP retrieves transactions for a user
func (h *TransactionHandler) ListTransactionsHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", "Only GET method is allowed")
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		h.writeErrorResponse(w, http.StatusUnauthorized, "Unauthorized", "User ID is required")
		return
	}

	// Parse query parameters
	filters, err := h.parseFiltersFromURL(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid query parameters", err.Error())
		return
	}

	transactions, total, err := h.transactionService.GetTransactionsByUser(userID, filters)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to get transactions", err.Error())
		return
	}

	// Convert to response
	transactionResponses := make([]*TransactionResponse, len(transactions))
	for i, transaction := range transactions {
		transactionResponses[i] = h.toTransactionResponse(transaction)
	}

	page := filters.Offset/filters.Limit + 1
	if filters.Limit == 0 {
		page = 1
	}

	response := TransactionListResponse{
		Transactions: transactionResponses,
		Total:        total,
		Page:         page,
		PageSize:     filters.Limit,
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}

// UpdateTransactionStatusHTTP updates the status of a transaction
func (h *TransactionHandler) UpdateTransactionStatusHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", "Only PUT method is allowed")
		return
	}

	id := h.extractTransactionID(r.URL.Path, "/status")
	if id == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request", "Transaction ID is required")
		return
	}

	var req struct {
		Status string `json:"status"`
		Reason string `json:"reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	if req.Status == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request", "Status is required")
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		h.writeErrorResponse(w, http.StatusUnauthorized, "Unauthorized", "User ID is required")
		return
	}

	transaction, err := h.transactionService.UpdateTransactionStatus(
		id,
		domaintransaction.TransactionStatus(req.Status),
		req.Reason,
		userID,
	)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.writeErrorResponse(w, http.StatusNotFound, "Transaction not found", err.Error())
			return
		}

		h.writeErrorResponse(w, http.StatusBadRequest, "Failed to update transaction status", err.Error())
		return
	}

	response := h.toTransactionResponse(transaction)
	h.writeJSONResponse(w, http.StatusOK, response)
}

// ProcessTransactionHTTP processes a pending transaction
func (h *TransactionHandler) ProcessTransactionHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", "Only POST method is allowed")
		return
	}

	id := h.extractTransactionID(r.URL.Path, "/process")
	if id == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request", "Transaction ID is required")
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		h.writeErrorResponse(w, http.StatusUnauthorized, "Unauthorized", "User ID is required")
		return
	}

	err := h.transactionService.ProcessTransaction(id, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.writeErrorResponse(w, http.StatusNotFound, "Transaction not found", err.Error())
			return
		}

		h.writeErrorResponse(w, http.StatusBadRequest, "Failed to process transaction", err.Error())
		return
	}

	// Get updated transaction
	transaction, err := h.transactionService.GetTransactionByID(id, userID)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to get updated transaction", err.Error())
		return
	}

	response := h.toTransactionResponse(transaction)
	h.writeJSONResponse(w, http.StatusOK, response)
}

// ReverseTransactionHTTP creates a reversal for a transaction
func (h *TransactionHandler) ReverseTransactionHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", "Only POST method is allowed")
		return
	}

	id := h.extractTransactionID(r.URL.Path, "/reverse")
	if id == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request", "Transaction ID is required")
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	if req.Reason == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request", "Reason is required")
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		h.writeErrorResponse(w, http.StatusUnauthorized, "Unauthorized", "User ID is required")
		return
	}

	reversalTransaction, err := h.transactionService.ReverseTransaction(id, req.Reason, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.writeErrorResponse(w, http.StatusNotFound, "Transaction not found", err.Error())
			return
		}

		h.writeErrorResponse(w, http.StatusBadRequest, "Failed to reverse transaction", err.Error())
		return
	}

	response := h.toTransactionResponse(reversalTransaction)
	h.writeJSONResponse(w, http.StatusCreated, response)
}

// Helper methods

// extractTransactionID extracts transaction ID from URL path
func (h *TransactionHandler) extractTransactionID(path, suffix string) string {
	// Remove the suffix and extract ID
	// Path format: /api/v1/transactions/{id}/suffix
	prefix := "/api/v1/transactions/"
	if !strings.HasPrefix(path, prefix) {
		return ""
	}

	remaining := strings.TrimPrefix(path, prefix)
	if suffix != "" {
		remaining = strings.TrimSuffix(remaining, suffix)
	}

	return remaining
}

// parseFiltersFromURL parses query parameters into TransactionFilters
func (h *TransactionHandler) parseFiltersFromURL(r *http.Request) (service.TransactionFilters, error) {
	filters := service.TransactionFilters{}
	query := r.URL.Query()

	// Parse pagination
	if page := query.Get("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			filters.Offset = (p - 1) * filters.Limit
		}
	}

	// Accept both 'limit' and 'pageSize' for compatibility
	if limit := query.Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
			filters.Limit = l
		}
	}

	if pageSize := query.Get("pageSize"); pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil && ps > 0 && ps <= 100 {
			filters.Limit = ps
		}
	}

	if filters.Limit == 0 {
		filters.Limit = 10 // Default page size changed to 10
	}

	// Parse types
	if types := query["type"]; len(types) > 0 {
		for _, typeStr := range types {
			filters.Types = append(filters.Types, domaintransaction.TransactionType(typeStr))
		}
	}

	// Parse statuses
	if statuses := query["status"]; len(statuses) > 0 {
		for _, statusStr := range statuses {
			filters.Statuses = append(filters.Statuses, domaintransaction.TransactionStatus(statusStr))
		}
	}

	// Parse amounts
	if minAmount := query.Get("minAmount"); minAmount != "" {
		if ma, err := strconv.ParseFloat(minAmount, 64); err == nil {
			filters.MinAmount = &ma
		}
	}

	if maxAmount := query.Get("maxAmount"); maxAmount != "" {
		if ma, err := strconv.ParseFloat(maxAmount, 64); err == nil {
			filters.MaxAmount = &ma
		}
	}

	// Parse order
	if orderBy := query.Get("orderBy"); orderBy != "" {
		filters.OrderBy = orderBy
	}

	if order := query.Get("order"); order != "" {
		filters.Order = order
	}

	return filters, nil
}

// toTransactionResponse converts domain transaction to response DTO
func (h *TransactionHandler) toTransactionResponse(transaction *domaintransaction.Transaction) *TransactionResponse {
	response := &TransactionResponse{
		ID:              transaction.ID,
		ReferenceID:     transaction.ReferenceID,
		ExternalID:      transaction.ExternalID,
		Type:            string(transaction.Type),
		Status:          string(transaction.Status),
		Amount:          transaction.Amount,
		Currency:        transaction.Currency,
		FromAccountID:   transaction.FromAccountID,
		ToAccountID:     transaction.ToAccountID,
		FromCardID:      transaction.FromCardID,
		ToCardID:        transaction.ToCardID,
		UserID:          transaction.UserID,
		InitiatedBy:     transaction.InitiatedBy,
		Description:     transaction.Description,
		PaymentMethod:   string(transaction.PaymentMethod),
		MerchantName:    transaction.MerchantName,
		MerchantID:      transaction.MerchantID,
		PreviousBalance: transaction.PreviousBalance,
		NewBalance:      transaction.NewBalance,
		FailureReason:   transaction.FailureReason,
		Metadata:        transaction.Metadata,
		Tags:            transaction.Tags,
		CreatedAt:       transaction.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:       transaction.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	// Handle nullable timestamps
	if transaction.ProcessedAt != nil {
		processedAt := transaction.ProcessedAt.Format("2006-01-02T15:04:05Z")
		response.ProcessedAt = &processedAt
	}

	if transaction.FailedAt != nil {
		failedAt := transaction.FailedAt.Format("2006-01-02T15:04:05Z")
		response.FailedAt = &failedAt
	}

	return response
}

// writeJSONResponse writes a JSON response
func (h *TransactionHandler) writeJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeErrorResponse writes an error response
func (h *TransactionHandler) writeErrorResponse(w http.ResponseWriter, status int, error string, message string) {
	response := ErrorResponse{
		Error:   error,
		Message: message,
		Code:    status,
	}
	h.writeJSONResponse(w, status, response)
}
