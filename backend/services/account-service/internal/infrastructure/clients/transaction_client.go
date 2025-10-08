package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// TransactionClient handles communication with the transaction service
type TransactionClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewTransactionClient creates a new transaction service client
func NewTransactionClient() *TransactionClient {
	baseURL := "http://transaction-service:8083"
	if url := os.Getenv("TRANSACTION_SERVICE_URL"); url != "" {
		baseURL = url
	}

	return &TransactionClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// CreateTransactionRequest represents the request to create a transaction
type CreateTransactionRequest struct {
	Type          string                 `json:"type"`
	Amount        float64                `json:"amount"`
	Currency      string                 `json:"currency"`
	FromAccountID *string                `json:"fromAccountId,omitempty"`
	ToAccountID   *string                `json:"toAccountId,omitempty"`
	Description   string                 `json:"description"`
	PaymentMethod string                 `json:"paymentMethod,omitempty"`
	MerchantName  string                 `json:"merchantName,omitempty"`
	ReferenceID   string                 `json:"referenceId,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// TransactionResponse represents the response from transaction service
type TransactionResponse struct {
	ID            string                 `json:"id"`
	ReferenceID   string                 `json:"referenceId"`
	Type          string                 `json:"type"`
	Status        string                 `json:"status"`
	Amount        float64                `json:"amount"`
	Currency      string                 `json:"currency"`
	FromAccountID *string                `json:"fromAccountId"`
	ToAccountID   *string                `json:"toAccountId"`
	Description   string                 `json:"description"`
	PaymentMethod string                 `json:"paymentMethod"`
	MerchantName  string                 `json:"merchantName"`
	Metadata      map[string]interface{} `json:"metadata"`
	CreatedAt     time.Time              `json:"createdAt"`
	UpdatedAt     time.Time              `json:"updatedAt"`
}

// CreateTransaction creates a transaction in the transaction service
func (c *TransactionClient) CreateTransaction(userID string, req CreateTransactionRequest) (*TransactionResponse, error) {
	url := fmt.Sprintf("%s/api/v1/transactions", c.baseURL)

	// Convert request to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-User-ID", userID)

	// Make request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("transaction service returned status %d", resp.StatusCode)
	}

	// Parse response
	var response TransactionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// CreateDebitCardTransaction creates a debit card transaction record
func (c *TransactionClient) CreateDebitCardTransaction(userID, accountID, cardID string, amount float64, description, merchantName, reference string) error {
	req := CreateTransactionRequest{
		Type:          "debit_purchase",
		Amount:        amount,
		Currency:      "ARS",
		FromAccountID: &accountID,
		Description:   description,
		PaymentMethod: "debit_card",
		MerchantName:  merchantName,
		ReferenceID:   reference,
		Metadata: map[string]interface{}{
			"cardId":     cardID,
			"category":   "purchase",
			"recordOnly": true, // This tells transaction service to only record, not modify balances
		},
	}

	_, err := c.CreateTransaction(userID, req)
	return err
}

// CreateInstallmentTransaction creates a transaction record for installment plan creation
func (c *TransactionClient) CreateInstallmentTransaction(userID, accountID, cardID string, amount float64, installmentsCount int, planID, description, merchantName, reference string) (*TransactionResponse, error) {
	req := CreateTransactionRequest{
		Type:          "credit_purchase_installments",
		Amount:        amount,
		Currency:      "ARS",
		FromAccountID: &accountID,
		Description:   fmt.Sprintf("Purchase with %d installments: %s", installmentsCount, description),
		PaymentMethod: "credit_card_installments",
		MerchantName:  merchantName,
		ReferenceID:   reference,
		Metadata: map[string]interface{}{
			"cardId":            cardID,
			"installmentPlanId": planID,
			"installmentsCount": installmentsCount,
			"category":          "installment_purchase",
			"recordOnly":        false, // This affects the credit balance
		},
	}

	return c.CreateTransaction(userID, req)
}

// CreateInstallmentPaymentTransaction creates a transaction record for installment payment
func (c *TransactionClient) CreateInstallmentPaymentTransaction(userID, accountID, cardID string, amount float64, installmentID, planID, installmentNumber string, description string) (*TransactionResponse, error) {
	req := CreateTransactionRequest{
		Type:          "installment_payment",
		Amount:        amount,
		Currency:      "ARS",
		ToAccountID:   &accountID, // Payment goes TO the account (reduces debt)
		Description:   fmt.Sprintf("Installment #%s payment: %s", installmentNumber, description),
		PaymentMethod: "installment_payment",
		ReferenceID:   fmt.Sprintf("installment-%s", installmentID),
		Metadata: map[string]interface{}{
			"cardId":            cardID,
			"installmentId":     installmentID,
			"installmentPlanId": planID,
			"installmentNumber": installmentNumber,
			"category":          "installment_payment",
			"recordOnly":        false, // This affects the credit balance
		},
	}

	return c.CreateTransaction(userID, req)
}

// CreateInstallmentCancellationTransaction creates a transaction record when an installment plan is cancelled
func (c *TransactionClient) CreateInstallmentCancellationTransaction(userID, accountID, cardID string, remainingAmount float64, planID, reason string) (*TransactionResponse, error) {
	req := CreateTransactionRequest{
		Type:          "installment_cancellation",
		Amount:        remainingAmount,
		Currency:      "ARS",
		ToAccountID:   &accountID, // Cancellation credits back to account
		Description:   fmt.Sprintf("Installment plan cancelled: %s", reason),
		PaymentMethod: "installment_cancellation",
		ReferenceID:   fmt.Sprintf("cancel-plan-%s", planID),
		Metadata: map[string]interface{}{
			"cardId":             cardID,
			"installmentPlanId":  planID,
			"cancellationReason": reason,
			"category":           "installment_cancellation",
			"recordOnly":         false,
		},
	}

	return c.CreateTransaction(userID, req)
}

// GetTransactionsByInstallmentPlan retrieves all transactions related to an installment plan
func (c *TransactionClient) GetTransactionsByInstallmentPlan(userID, planID string) ([]*TransactionResponse, error) {
	url := fmt.Sprintf("%s/api/v1/transactions?installmentPlanId=%s", c.baseURL, planID)

	// Create HTTP request
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("X-User-ID", userID)

	// Make request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("transaction service returned status %d", resp.StatusCode)
	}

	// Parse response
	var transactions []*TransactionResponse
	if err := json.NewDecoder(resp.Body).Decode(&transactions); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return transactions, nil
}

// HealthCheck verifica la conectividad con el transaction service
func (c *TransactionClient) HealthCheck() error {
	url := fmt.Sprintf("%s/health", c.baseURL)

	// Create HTTP request
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	// Make request with shorter timeout for health check
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to connect to transaction service: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("transaction service health check failed with status %d", resp.StatusCode)
	}

	return nil
}
