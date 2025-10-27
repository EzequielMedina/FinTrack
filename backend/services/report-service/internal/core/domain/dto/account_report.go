package dto

// AccountReportRequest request para reporte de cuentas
type AccountReportRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

// AccountReportResponse respuesta del reporte de cuentas
type AccountReportResponse struct {
	UserID       string                `json:"user_id"`
	Summary      AccountSummary        `json:"summary"`
	Accounts     []AccountDetail       `json:"accounts"`
	Cards        []CardDetail          `json:"cards"`
	Distribution []AccountDistribution `json:"distribution"`
}

// AccountSummary resumen de cuentas
type AccountSummary struct {
	TotalBalance      float64 `json:"total_balance"`
	TotalAccounts     int     `json:"total_accounts"`
	TotalCards        int     `json:"total_cards"`
	TotalCreditLimit  float64 `json:"total_credit_limit"`
	TotalCreditUsed   float64 `json:"total_credit_used"`
	AvailableCredit   float64 `json:"available_credit"`
	CreditUtilization float64 `json:"credit_utilization"`
	NetWorth          float64 `json:"net_worth"`
}

// AccountDetail detalle de cuenta
type AccountDetail struct {
	ID          string  `json:"id"`
	AccountType string  `json:"account_type"`
	Name        string  `json:"name"`
	Currency    string  `json:"currency"`
	Balance     float64 `json:"balance"`
	CreditLimit float64 `json:"credit_limit,omitempty"`
	IsActive    bool    `json:"is_active"`
}

// CardDetail detalle de tarjeta
type CardDetail struct {
	ID              string  `json:"id"`
	AccountID       string  `json:"account_id"`
	CardType        string  `json:"card_type"`
	CardBrand       string  `json:"card_brand"`
	LastFourDigits  string  `json:"last_four_digits"`
	HolderName      string  `json:"holder_name"`
	Status          string  `json:"status"`
	CreditLimit     float64 `json:"credit_limit,omitempty"`
	CurrentBalance  float64 `json:"current_balance,omitempty"`
	AvailableCredit float64 `json:"available_credit,omitempty"`
	Nickname        string  `json:"nickname,omitempty"`
}

// AccountDistribution distribución de cuentas
type AccountDistribution struct {
	AccountType  string  `json:"account_type"`
	Count        int     `json:"count"`
	TotalBalance float64 `json:"total_balance"`
	Percentage   float64 `json:"percentage"`
}
