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
