package errors

import "errors"

var (
	// Authentication errors
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")

	// User management errors
	ErrUserInactive         = errors.New("user is inactive")
	ErrInvalidRole          = errors.New("invalid role")
	ErrCannotDeleteSelf     = errors.New("cannot delete yourself")
	ErrCannotDeactivateSelf = errors.New("cannot deactivate yourself")
	ErrAdminRequired        = errors.New("admin role required for this operation")
	ErrCannotCreateAdmin    = errors.New("admins cannot create other admin users")

	// Validation errors
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrPasswordTooWeak = errors.New("password does not meet requirements")
	ErrInvalidUserData = errors.New("invalid user data")
	ErrEmptyFirstName  = errors.New("first name is required")
	ErrEmptyLastName   = errors.New("last name is required")

	// General errors
	ErrUnauthorized      = errors.New("unauthorized access")
	ErrForbidden         = errors.New("forbidden operation")
	ErrInvalidPagination = errors.New("invalid pagination parameters")
)
