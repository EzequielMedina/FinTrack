package ports

import (
	"context"
	"time"
)

// ChatbotService define el servicio principal conversacional
type ChatbotService interface {
	HandleQuery(ctx context.Context, req ChatQueryRequest) (ChatQueryResponse, error)
	GeneratePDF(ctx context.Context, req ReportRequest) ([]byte, error)
	GenerateChartData(ctx context.Context, req ChartRequest) (ChartResponse, error)
}

// DataProvider acceso directo a MySQL compartida
type DataProvider interface {
	GetTotals(ctx context.Context, userID string, from, to time.Time) (Totals, error)
	GetByType(ctx context.Context, userID string, from, to time.Time) (map[string]float64, error)
	GetTopMerchants(ctx context.Context, userID string, from, to time.Time, limit int) ([]MerchantTotal, error)
	GetByAccountType(ctx context.Context, userID string, from, to time.Time) (map[string]float64, error)
	GetByCard(ctx context.Context, userID string, from, to time.Time) ([]CardTotal, error)
	// Installments
	GetInstallmentsSummary(ctx context.Context, userID string, from, to time.Time) (InstallmentsSummary, error)
	GetInstallmentPlans(ctx context.Context, userID string) ([]InstallmentPlanInfo, error)
	GetInstallmentsByMonth(ctx context.Context, userID string) (map[string]InstallmentMonthSummary, error)
	// Información detallada adicional
	GetRecentTransactions(ctx context.Context, userID string, from, to time.Time, limit int) ([]TransactionDetail, error)
	GetAccountsInfo(ctx context.Context, userID string) ([]AccountInfo, error)
	GetCardsInfo(ctx context.Context, userID string) ([]CardInfo, error)
	GetExchangeRates(ctx context.Context, userID string, from, to time.Time) ([]ExchangeRateInfo, error)
}

// LLMProvider interfaz al motor LLM (compatible con Ollama, Groq, etc.)
type LLMProvider interface {
	Chat(ctx context.Context, systemPrompt string, userPrompt string) (string, error)
}

// ReportProvider genera PDFs
type ReportProvider interface {
	Generate(ctx context.Context, data ReportData) ([]byte, error)
}

// Modelos de requests/responses
type Period struct {
	From time.Time
	To   time.Time
}

type ChatQueryRequest struct {
	UserID  string         `json:"userId"`
	Message string         `json:"message"`
	Period  Period         `json:"period"`
	Filters map[string]any `json:"filters"`
}

type SuggestedAction struct {
	Type   string         `json:"type"`
	Params map[string]any `json:"params"`
}

type ChatQueryResponse struct {
	Reply            string            `json:"reply"`
	SuggestedActions []SuggestedAction `json:"suggestedActions"`
	Insights         []string          `json:"insights"`
	DataRefs         map[string]any    `json:"dataRefs"`
}

type ReportRequest struct {
	UserID        string `json:"userId"`
	Title         string `json:"title"`
	Period        Period `json:"period"`
	GroupBy       string `json:"groupBy"`
	IncludeCharts bool   `json:"includeCharts"`
}

type ChartRequest struct {
	UserID   string `json:"userId"`
	Period   Period `json:"period"`
	GroupBy  string `json:"groupBy"`
	Currency string `json:"currency"`
}

type ChartResponse struct {
	Labels   []string       `json:"labels"`
	Datasets []ChartDataset `json:"datasets"`
	Meta     map[string]any `json:"meta"`
}

type ChartDataset struct {
	Label           string    `json:"label"`
	Data            []float64 `json:"data"`
	BackgroundColor []string  `json:"backgroundColor"`
}

// Totales y agregaciones
type Totals struct {
	Expenses float64
	Incomes  float64
}

type MerchantTotal struct {
	Merchant string
	Total    float64
}

type CardTotal struct {
	Brand    string
	LastFour string
	Total    float64
}

// Installments summary for conversational context
type InstallmentsSummary struct {
	Active          int        `json:"active"`
	Completed       int        `json:"completed"`
	Overdue         int        `json:"overdue"`
	RemainingAmount float64    `json:"remainingAmount"`
	NextDueDate     *time.Time `json:"nextDueDate"`
	PaidInPeriod    float64    `json:"paidInPeriod"`
}

// Installments grouped by month (for future planning questions)
type InstallmentMonthSummary struct {
	YearMonth string  `json:"yearMonth"` // Format: "2025-11"
	Count     int     `json:"count"`
	Total     float64 `json:"total"`
}

// High-level info per installment plan
type InstallmentPlanInfo struct {
	ID                    string     `json:"id"`
	CardID                string     `json:"cardId"`
	MerchantName          string     `json:"merchantName"`
	Description           string     `json:"description"`
	InstallmentsCount     int        `json:"installmentsCount"`
	RemainingInstallments int        `json:"remainingInstallments"`
	RemainingAmount       float64    `json:"remainingAmount"`
	Status                string     `json:"status"`
	NextDueDate           *time.Time `json:"nextDueDate"`
	NextInstallmentAmount float64    `json:"nextInstallmentAmount"`
	CompletionPercentage  float64    `json:"completionPercentage"`
	CreatedAt             *time.Time `json:"createdAt"`
}

// Información detallada adicional para mejor contexto
type TransactionDetail struct {
	ID            string    `json:"id"`
	Type          string    `json:"type"`
	Amount        float64   `json:"amount"`
	MerchantName  string    `json:"merchantName"`
	Description   string    `json:"description"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
	FromAccountID string    `json:"fromAccountId"`
	ToAccountID   string    `json:"toAccountId"`
	FromCardID    string    `json:"fromCardId"`
	ToCardID      string    `json:"toCardId"`
	Currency      string    `json:"currency"`
}

type AccountInfo struct {
	ID          string    `json:"id"`
	AccountType string    `json:"accountType"`
	Balance     float64   `json:"balance"`
	Currency    string    `json:"currency"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
}

type CardInfo struct {
	ID          string    `json:"id"`
	CardBrand   string    `json:"cardBrand"`
	LastFour    string    `json:"lastFour"`
	CardType    string    `json:"cardType"`
	Status      string    `json:"status"`
	CreditLimit float64   `json:"creditLimit"`
	CurrentDebt float64   `json:"currentDebt"`
	CreatedAt   time.Time `json:"createdAt"`
}

type ExchangeRateInfo struct {
	ID           string    `json:"id"`
	FromCurrency string    `json:"fromCurrency"`
	ToCurrency   string    `json:"toCurrency"`
	Rate         float64   `json:"rate"`
	Source       string    `json:"source"`
	CreatedAt    time.Time `json:"createdAt"`
}

type ReportData struct {
	Title         string
	Period        Period
	Totals        Totals
	ByType        map[string]float64
	TopMerchants  []MerchantTotal
	ByAccountType map[string]float64
	ByCard        []CardTotal
	Currency      string
	// Installments
	Installments InstallmentsSummary
	Plans        []InstallmentPlanInfo
}
