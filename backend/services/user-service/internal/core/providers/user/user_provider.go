package userprovider

import (
	domuser "github.com/fintrack/user-service/internal/core/domain/entities/user"
)

// Repository defines the interface for user data persistence operations
type Repository interface {
	// Basic CRUD operations
	GetByEmail(email string) (*domuser.User, error)
	GetByID(id string) (*domuser.User, error)
	Create(u *domuser.User) error
	Update(u *domuser.User) error
	Delete(id string) error

	// Query operations
	GetAll(limit, offset int) ([]*domuser.User, int, error) // returns users, total count, error
	GetByRole(role domuser.Role, limit, offset int) ([]*domuser.User, int, error)
	GetActiveUsers(limit, offset int) ([]*domuser.User, int, error)

	// Utility operations
	ExistsByEmail(email string) (bool, error)
	ExistsByID(id string) (bool, error)
	UpdateLastLogin(id string) error
	ToggleActiveStatus(id string, isActive bool) error
}
