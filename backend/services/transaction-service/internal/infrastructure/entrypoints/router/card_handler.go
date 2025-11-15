package router

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"

	"github.com/fintrack/transaction-service/internal/core/service"
	"github.com/fintrack/transaction-service/internal/infrastructure/http/clients"
	"github.com/fintrack/transaction-service/internal/infrastructure/repositories/mysql"
)

// CardHandler handles HTTP requests for card operations
type CardHandler struct {
	transactionService service.TransactionServiceInterface
	accountClient      *clients.AccountClient
}

// NewCardHandler creates a new card handler
func NewCardHandler(db *sql.DB) *CardHandler {
	// Create repositories
	transactionRepo := mysql.NewTransactionRepository(db)

	// Create mock services for now
	ruleService := service.NewTransactionRuleService(nil)   // TODO: Implement rule repository
	auditService := service.NewTransactionAuditService(nil) // TODO: Implement audit repository
	externalService := service.NewMockExternalService()

	// Create account service client - using environment variable or default localhost
	accountServiceURL := "http://localhost:8081" // Default account service port
	if url := os.Getenv("ACCOUNT_SERVICE_URL"); url != "" {
		accountServiceURL = url
	}
	accountClient := clients.NewAccountClient(accountServiceURL)

	// Create transaction service
	transactionService := service.NewTransactionService(
		transactionRepo,
		ruleService,
		auditService,
		externalService,
		accountClient,
	)

	return &CardHandler{
		transactionService: transactionService,
		accountClient:      accountClient,
	}
}

// CreditCardChargeRequest representa una solicitud de cargo a tarjeta de crédito
type CreditCardChargeRequest struct {
	CardID      string  `json:"card_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	MerchantName string `json:"merchant_name,omitempty"`
}

// CreditCardPaymentRequest representa una solicitud de pago a tarjeta de crédito
type CreditCardPaymentRequest struct {
	CardID        string  `json:"card_id"`
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"payment_method"`
}

// DebitCardTransactionRequest representa una solicitud de transacción con tarjeta de débito
type DebitCardTransactionRequest struct {
	CardID       string  `json:"card_id"`
	Amount       float64 `json:"amount"`
	Description  string  `json:"description"`
	MerchantName string  `json:"merchant_name"`
}

// CardTransactionResponse representa la respuesta de una transacción de tarjeta
type CardTransactionResponse struct {
	Success       bool    `json:"success"`
	TransactionID string  `json:"transaction_id"`
	CardID        string  `json:"card_id"`
	Amount        float64 `json:"amount"`
	NewBalance    float64 `json:"new_balance,omitempty"`
	Message       string  `json:"message"`
	Status        string  `json:"status"`
}

// ChargeCreditCardHTTP handles HTTP requests for credit card charges
func (h *CardHandler) ChargeCreditCardHTTP(w http.ResponseWriter, r *http.Request) {
	var request CreditCardChargeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if request.CardID == "" || request.Amount <= 0 || request.Description == "" {
		http.Error(w, "Missing required fields: card_id, amount, description", http.StatusBadRequest)
		return
	}

	// Create charge request for AccountClient
	chargeReq := clients.CardChargeRequest{
		CardID:      request.CardID,
		Amount:      request.Amount,
		Description: request.Description,
		Reference:   request.MerchantName,
	}

	// Process the charge using AccountClient
	chargeResponse, err := h.accountClient.ChargeCard(chargeReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create transaction record using existing TransactionService method
	createRequest := service.CreateTransactionRequest{
		UserID:       "system", // TODO: Extract from card info
		Type:         "CREDIT_CARD_CHARGE",
		Amount:       request.Amount,
		Description:  request.Description,
		MerchantName: request.MerchantName,
		Metadata: map[string]interface{}{
			"card_id":          request.CardID,
			"account_response": chargeResponse.Transaction,
		},
	}

	transaction, err := h.transactionService.CreateTransaction(createRequest, "api")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Mark transaction as completed since AccountService already processed it
	if err := h.transactionService.CompleteTransaction(transaction.ID, "api"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	response := CardTransactionResponse{
		Success:       true,
		TransactionID: transaction.ID,
		CardID:        request.CardID,
		Amount:        request.Amount,
		Message:       "Credit card charge processed successfully",
		Status:        "completed",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// PayCreditCardHTTP handles HTTP requests for credit card payments
func (h *CardHandler) PayCreditCardHTTP(w http.ResponseWriter, r *http.Request) {
	var request CreditCardPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if request.CardID == "" || request.Amount <= 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Create payment request for AccountClient
	paymentReq := clients.CardPaymentRequest{
		CardID:        request.CardID,
		Amount:        request.Amount,
		PaymentMethod: request.PaymentMethod,
	}

	// Process the payment using AccountClient
	paymentResponse, err := h.accountClient.PaymentCard(paymentReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create transaction record using existing TransactionService method
	createRequest := service.CreateTransactionRequest{
		UserID:      "system", // TODO: Extract from card info
		Type:        "CREDIT_CARD_PAYMENT",
		Amount:      request.Amount,
		Description: "Credit card payment",
		Metadata: map[string]interface{}{
			"card_id":          request.CardID,
			"payment_method":   request.PaymentMethod,
			"account_response": paymentResponse.Transaction,
		},
	}

	transaction, err := h.transactionService.CreateTransaction(createRequest, "api")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Mark transaction as completed since AccountService already processed it
	if err := h.transactionService.CompleteTransaction(transaction.ID, "api"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	response := CardTransactionResponse{
		Success:       true,
		TransactionID: transaction.ID,
		CardID:        request.CardID,
		Amount:        request.Amount,
		Message:       "Credit card payment processed successfully",
		Status:        "completed",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ProcessDebitCardTransactionHTTP handles HTTP requests for debit card transactions
func (h *CardHandler) ProcessDebitCardTransactionHTTP(w http.ResponseWriter, r *http.Request) {
	var request DebitCardTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if request.CardID == "" || request.Amount <= 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Create debit transaction request for AccountClient
	debitReq := clients.DebitTransactionRequest{
		CardID:       request.CardID,
		Amount:       request.Amount,
		Description:  request.Description,
		MerchantName: request.MerchantName,
	}

	// Process the debit transaction using AccountClient
	debitResponse, err := h.accountClient.ProcessDebitTransaction(debitReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create transaction record using existing TransactionService method
	createRequest := service.CreateTransactionRequest{
		UserID:       "system", // TODO: Extract from card info
		Type:         "DEBIT_CARD_TRANSACTION",
		Amount:       request.Amount,
		Description:  request.Description,
		MerchantName: request.MerchantName,
		Metadata: map[string]interface{}{
			"card_id":          request.CardID,
			"account_response": debitResponse.Transaction,
		},
	}

	transaction, err := h.transactionService.CreateTransaction(createRequest, "api")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Mark transaction as completed since AccountService already processed it
	if err := h.transactionService.CompleteTransaction(transaction.ID, "api"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	response := CardTransactionResponse{
		Success:       true,
		TransactionID: transaction.ID,
		CardID:        request.CardID,
		Amount:        request.Amount,
		Message:       "Debit card transaction processed successfully",
		Status:        "completed",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}