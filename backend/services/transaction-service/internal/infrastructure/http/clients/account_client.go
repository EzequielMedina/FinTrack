package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// AccountClient maneja la comunicación con el account-service
type AccountClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewAccountClient crea una nueva instancia del cliente
func NewAccountClient(baseURL string) *AccountClient {
	return &AccountClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// AccountBalance representa el balance de una cuenta
type AccountBalance struct {
	AccountID string  `json:"account_id"`
	Balance   float64 `json:"balance"`
	Currency  string  `json:"currency"`
}

// AccountInfo representa información básica de una cuenta
type AccountInfo struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	AccountType string  `json:"account_type"`
	Balance     float64 `json:"balance"`
	Currency    string  `json:"currency"`
	IsActive    bool    `json:"is_active"`
	CreditLimit float64 `json:"creditLimit,omitempty"`
}

// BalanceUpdateRequest representa una solicitud de actualización de balance
type BalanceUpdateRequest struct {
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Reference   string  `json:"reference,omitempty"`
}

// BalanceUpdateResponse representa la respuesta de actualización de balance
type BalanceUpdateResponse struct {
	Success    bool    `json:"success"`
	NewBalance float64 `json:"balance"`
	Message    string  `json:"message,omitempty"`
}

// GetAccountBalance obtiene el balance actual de una cuenta
func (c *AccountClient) GetAccountBalance(accountID string) (*AccountBalance, error) {
	url := fmt.Sprintf("%s/api/accounts/%s/balance", c.baseURL, accountID)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error calling account service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("account service returned status %d", resp.StatusCode)
	}

	var balance AccountBalance
	if err := json.NewDecoder(resp.Body).Decode(&balance); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &balance, nil
}

// GetAccountInfo obtiene información completa de una cuenta
func (c *AccountClient) GetAccountInfo(accountID string) (*AccountInfo, error) {
	url := fmt.Sprintf("%s/api/accounts/%s", c.baseURL, accountID)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error calling account service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("account service returned status %d", resp.StatusCode)
	}

	var account AccountInfo
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &account, nil
}

// AddFunds agrega fondos a una cuenta (para depósitos)
func (c *AccountClient) AddFunds(accountID string, amount float64, description string, reference string) (*BalanceUpdateResponse, error) {
	url := fmt.Sprintf("%s/api/accounts/%s/add-funds", c.baseURL, accountID)

	request := BalanceUpdateRequest{
		Amount:      amount,
		Description: description,
		Reference:   reference,
	}

	return c.updateBalance(url, request)
}

// WithdrawFunds retira fondos de una cuenta (para retiros)
func (c *AccountClient) WithdrawFunds(accountID string, amount float64, description string, reference string) (*BalanceUpdateResponse, error) {
	url := fmt.Sprintf("%s/api/accounts/%s/withdraw-funds", c.baseURL, accountID)

	request := BalanceUpdateRequest{
		Amount:      amount,
		Description: description,
		Reference:   reference,
	}

	return c.updateBalance(url, request)
}

// UpdateCreditUsage actualiza el uso de crédito para tarjetas de crédito
func (c *AccountClient) UpdateCreditUsage(accountID string, amount float64, description string, reference string) (*BalanceUpdateResponse, error) {
	url := fmt.Sprintf("%s/api/accounts/%s/update-credit", c.baseURL, accountID)

	request := BalanceUpdateRequest{
		Amount:      amount,
		Description: description,
		Reference:   reference,
	}

	return c.updateBalance(url, request)
}

// GetAvailableCredit obtiene el crédito disponible para una tarjeta de crédito
func (c *AccountClient) GetAvailableCredit(accountID string) (*AccountBalance, error) {
	url := fmt.Sprintf("%s/api/accounts/%s/available-credit", c.baseURL, accountID)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error calling account service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("account service returned status %d", resp.StatusCode)
	}

	var creditInfo AccountBalance
	if err := json.NewDecoder(resp.Body).Decode(&creditInfo); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &creditInfo, nil
}

// ValidateAccountExists verifica si una cuenta existe y está activa
func (c *AccountClient) ValidateAccountExists(accountID string) (bool, error) {
	account, err := c.GetAccountInfo(accountID)
	if err != nil {
		return false, err
	}

	return account.IsActive, nil
}

// Helper method para actualizaciones de balance
func (c *AccountClient) updateBalance(url string, request BalanceUpdateRequest) (*BalanceUpdateResponse, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error calling account service: %w", err)
	}
	defer resp.Body.Close()

	var response BalanceUpdateResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return &response, fmt.Errorf("account service returned status %d: %s", resp.StatusCode, response.Message)
	}

	return &response, nil
}

// CARD OPERATIONS

// CardChargeRequest representa una solicitud de cargo a tarjeta de crédito
type CardChargeRequest struct {
	CardID      string  `json:"card_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Reference   string  `json:"reference,omitempty"`
}

// CardPaymentRequest representa una solicitud de pago a tarjeta de crédito
type CardPaymentRequest struct {
	CardID        string  `json:"card_id"`
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"payment_method"`
	Reference     string  `json:"reference,omitempty"`
}

// DebitTransactionRequest representa una solicitud de transacción con tarjeta de débito
type DebitTransactionRequest struct {
	CardID       string  `json:"card_id"`
	Amount       float64 `json:"amount"`
	Description  string  `json:"description"`
	MerchantName string  `json:"merchant_name"`
	Reference    string  `json:"reference,omitempty"`
}

// CardOperationResponse representa la respuesta de operaciones de tarjetas
type CardOperationResponse struct {
	Success     bool    `json:"success"`
	CardID      string  `json:"card_id"`
	NewBalance  float64 `json:"new_balance,omitempty"`
	Message     string  `json:"message,omitempty"`
	Transaction string  `json:"transaction_id,omitempty"`
}

// ChargeCard procesa un cargo en una tarjeta de crédito
func (c *AccountClient) ChargeCard(req CardChargeRequest) (*CardOperationResponse, error) {
	url := fmt.Sprintf("%s/api/cards/%s/charge", c.baseURL, req.CardID)
	return c.processCardOperation(url, req)
}

// PaymentCard procesa un pago a una tarjeta de crédito
func (c *AccountClient) PaymentCard(req CardPaymentRequest) (*CardOperationResponse, error) {
	url := fmt.Sprintf("%s/api/cards/%s/payment", c.baseURL, req.CardID)
	return c.processCardOperation(url, req)
}

// ProcessDebitTransaction procesa una transacción con tarjeta de débito
func (c *AccountClient) ProcessDebitTransaction(req DebitTransactionRequest) (*CardOperationResponse, error) {
	url := fmt.Sprintf("%s/api/cards/%s/transaction", c.baseURL, req.CardID)
	return c.processCardOperation(url, req)
}

// Helper method para operaciones de tarjetas
func (c *AccountClient) processCardOperation(url string, request interface{}) (*CardOperationResponse, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error calling account service: %w", err)
	}
	defer resp.Body.Close()

	var response CardOperationResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return &response, fmt.Errorf("account service returned status %d: %s", resp.StatusCode, response.Message)
	}

	return &response, nil
}

// HealthCheck verifica si el account-service está disponible
func (c *AccountClient) HealthCheck() error {
	url := fmt.Sprintf("%s/health", c.baseURL)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("account service unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("account service health check failed with status %d", resp.StatusCode)
	}

	return nil
}
