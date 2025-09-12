package service

import (
	domuser "github.com/fintrack/user-service/internal/core/domain/entities/user"
)

// UserServiceInterface defines the contract for user service operations
type UserServiceInterface interface {
	CreateUser(email, password, firstName, lastName string, role domuser.Role, currentUser *domuser.User) (*domuser.User, error)
	GetUserByID(id string, currentUser *domuser.User) (*domuser.User, error)
	GetAllUsers(limit, offset int, currentUser *domuser.User) ([]*domuser.User, int, error)
	GetUsersByRole(role domuser.Role, limit, offset int, currentUser *domuser.User) ([]*domuser.User, int, error)
	UpdateUser(id string, updates map[string]interface{}, currentUser *domuser.User) (*domuser.User, error)
	UpdateUserProfile(id string, profile domuser.Profile, currentUser *domuser.User) (*domuser.User, error)
	ChangeUserRole(id string, newRole domuser.Role, currentUser *domuser.User) (*domuser.User, error)
	ToggleUserStatus(id string, isActive bool, currentUser *domuser.User) (*domuser.User, error)
	DeleteUser(id string, currentUser *domuser.User) error
	ChangePassword(id, oldPassword, newPassword string, currentUser *domuser.User) error
}