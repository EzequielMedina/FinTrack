package user

import (
	"regexp"
	"strings"
	"time"
)

// Role represents user roles in the system
// @Description User role enumeration
type Role string

const (
	RoleUser      Role = "user"      // @Description Regular user role
	RoleOperator  Role = "operator"  // @Description Operator role
	RoleAdmin     Role = "admin"     // @Description Administrator role
	RoleTreasurer Role = "treasurer" // @Description Treasurer role
)

// User represents the main user entity with authentication and profile information
type User struct {
	// Core identity fields
	ID           string `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`

	// Basic profile information
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`

	// Authorization and status
	Role          Role `json:"role"`
	IsActive      bool `json:"isActive"`
	EmailVerified bool `json:"emailVerified"`

	// Profile information
	Profile Profile `json:"profile"`

	// Audit fields
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	LastLoginAt *time.Time `json:"lastLoginAt,omitempty"`
}

// Profile represents extended user profile information
type Profile struct {
	Phone          string      `json:"phone,omitempty"`
	DateOfBirth    *time.Time  `json:"dateOfBirth,omitempty"`
	Address        Address     `json:"address"`
	ProfilePicture string      `json:"profilePicture,omitempty"`
	Preferences    Preferences `json:"preferences"`
}

// Address represents user address information
type Address struct {
	Street     string `json:"street,omitempty"`
	City       string `json:"city,omitempty"`
	State      string `json:"state,omitempty"`
	PostalCode string `json:"postalCode,omitempty"`
	Country    string `json:"country,omitempty"`
}

// Preferences represents user application preferences
type Preferences struct {
	Language          string `json:"language,omitempty"`
	Timezone          string `json:"timezone,omitempty"`
	NotificationEmail bool   `json:"notificationEmail"`
	NotificationSMS   bool   `json:"notificationSMS"`
}

// IsValidRole checks if the provided role is valid
func IsValidRole(role Role) bool {
	switch role {
	case RoleUser, RoleOperator, RoleAdmin, RoleTreasurer:
		return true
	default:
		return false
	}
}

// GetFullName returns the user's full name
func (u *User) GetFullName() string {
	return strings.TrimSpace(u.FirstName + " " + u.LastName)
}

// IsValidEmail validates email format
func (u *User) IsValidEmail() bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(u.Email)
}

// CanModifyUser checks if the current user can modify another user
func (u *User) CanModifyUser(targetUser *User) bool {
	// Admins can modify anyone except themselves for certain operations
	if u.Role == RoleAdmin {
		return true
	}

	// Users can only modify themselves
	return u.ID == targetUser.ID
}

// CanDeleteUser checks if the current user can delete another user
func (u *User) CanDeleteUser(targetUser *User) bool {
	// Only admins can delete users
	if u.Role != RoleAdmin {
		return false
	}

	// Admins cannot delete themselves
	return u.ID != targetUser.ID
}

// HasPermission checks if user has specific permission based on role
func (u *User) HasPermission(requiredRole Role) bool {
	roleHierarchy := map[Role]int{
		RoleUser:      1,
		RoleOperator:  2,
		RoleTreasurer: 3,
		RoleAdmin:     4,
	}

	userLevel, exists := roleHierarchy[u.Role]
	if !exists {
		return false
	}

	requiredLevel, exists := roleHierarchy[requiredRole]
	if !exists {
		return false
	}

	return userLevel >= requiredLevel
}
