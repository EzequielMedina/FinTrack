package errors

import "fmt"

// Domain errors for the account service
var (
	// Account not found errors
	ErrAccountNotFound = fmt.Errorf("account not found")
	ErrUserNotFound    = fmt.Errorf("user not found")

	// Validation errors
	ErrInvalidAccountData   = fmt.Errorf("invalid account data")
	ErrInvalidAccountType   = fmt.Errorf("invalid account type")
	ErrInvalidCurrency      = fmt.Errorf("invalid currency")
	ErrInvalidAmount        = fmt.Errorf("invalid amount")
	ErrInvalidAccountStatus = fmt.Errorf("invalid account status")

	// Business logic errors
	ErrInsufficientBalance            = fmt.Errorf("insufficient balance")
	ErrAccountNotActive               = fmt.Errorf("account is not active")
	ErrAccountClosed                  = fmt.Errorf("account is closed")
	ErrAccountFrozen                  = fmt.Errorf("account is frozen")
	ErrDailyLimitExceeded             = fmt.Errorf("daily limit exceeded")
	ErrMonthlyLimitExceeded           = fmt.Errorf("monthly limit exceeded")
	ErrMaxAccountsReached             = fmt.Errorf("maximum number of accounts reached")
	ErrMaxAccountsExceeded            = fmt.Errorf("maximum number of accounts exceeded")
	ErrCannotDeleteAccountWithBalance = fmt.Errorf("cannot delete account with balance")
	ErrDuplicateAccountNumber         = fmt.Errorf("account number already exists")
	ErrDuplicateAccountName           = fmt.Errorf("account name already exists for user")
	ErrInvalidInput                   = fmt.Errorf("invalid input")

	// Permission errors
	ErrUnauthorized       = fmt.Errorf("unauthorized access")
	ErrInsufficientRights = fmt.Errorf("insufficient rights")

	// External service errors
	ErrExternalServiceUnavailable = fmt.Errorf("external service unavailable")
	ErrCurrencyConversionFailed   = fmt.Errorf("currency conversion failed")

	// Database errors
	ErrDatabaseConnection = fmt.Errorf("database connection error")
	ErrTransactionFailed  = fmt.Errorf("database transaction failed")
)

// AccountError represents a custom error with additional context
type AccountError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
	Field   string `json:"field,omitempty"`
}

func (e AccountError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Details)
	}
	return e.Message
}

// NewAccountError creates a new AccountError
func NewAccountError(code, message, details string) *AccountError {
	return &AccountError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// NewValidationError creates a validation error for a specific field
func NewValidationError(field, message string) *AccountError {
	return &AccountError{
		Code:    "VALIDATION_ERROR",
		Message: message,
		Field:   field,
	}
}

// Error codes constants
const (
	CodeAccountNotFound        = "ACCOUNT_NOT_FOUND"
	CodeUserNotFound           = "USER_NOT_FOUND"
	CodeInvalidData            = "INVALID_DATA"
	CodeInsufficientBalance    = "INSUFFICIENT_BALANCE"
	CodeAccountNotActive       = "ACCOUNT_NOT_ACTIVE"
	CodeAccountClosed          = "ACCOUNT_CLOSED"
	CodeAccountFrozen          = "ACCOUNT_FROZEN"
	CodeDailyLimitExceeded     = "DAILY_LIMIT_EXCEEDED"
	CodeMonthlyLimitExceeded   = "MONTHLY_LIMIT_EXCEEDED"
	CodeMaxAccountsReached     = "MAX_ACCOUNTS_REACHED"
	CodeDuplicateAccountNumber = "DUPLICATE_ACCOUNT_NUMBER"
	CodeDuplicateAccountName   = "DUPLICATE_ACCOUNT_NAME"
	CodeUnauthorized           = "UNAUTHORIZED"
	CodeInsufficientRights     = "INSUFFICIENT_RIGHTS"
	CodeExternalServiceError   = "EXTERNAL_SERVICE_ERROR"
	CodeDatabaseError          = "DATABASE_ERROR"
	CodeValidationError        = "VALIDATION_ERROR"
)

// IsNotFoundError checks if the error is a not found error
func IsNotFoundError(err error) bool {
	return err == ErrAccountNotFound || err == ErrUserNotFound
}

// IsValidationError checks if the error is a validation error
func IsValidationError(err error) bool {
	_, ok := err.(*AccountError)
	return ok || err == ErrInvalidAccountData || err == ErrInvalidAccountType || err == ErrInvalidCurrency
}

// IsBusinessLogicError checks if the error is a business logic error
func IsBusinessLogicError(err error) bool {
	return err == ErrInsufficientBalance ||
		err == ErrAccountNotActive ||
		err == ErrAccountClosed ||
		err == ErrAccountFrozen ||
		err == ErrDailyLimitExceeded ||
		err == ErrMonthlyLimitExceeded ||
		err == ErrMaxAccountsReached ||
		err == ErrDuplicateAccountNumber ||
		err == ErrDuplicateAccountName
}

// IsPermissionError checks if the error is a permission error
func IsPermissionError(err error) bool {
	return err == ErrUnauthorized || err == ErrInsufficientRights
}
