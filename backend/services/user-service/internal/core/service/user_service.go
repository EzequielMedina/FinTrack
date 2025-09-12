package service

import (
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	domuser "github.com/fintrack/user-service/internal/core/domain/entities/user"
	domerrors "github.com/fintrack/user-service/internal/core/errors"
	userprovider "github.com/fintrack/user-service/internal/core/providers/user"
	"github.com/google/uuid"
)

// UserService handles business logic for user management operations
type UserService struct {
	repo userprovider.Repository
}

// NewUserService creates a new UserService instance
func NewUserService(repo userprovider.Repository) *UserService {
	return &UserService{repo: repo}
}

// CreateUser creates a new user with validation
func (s *UserService) CreateUser(email, password, firstName, lastName string, role domuser.Role, currentUser *domuser.User) (*domuser.User, error) {
	// Authorization check
	if currentUser.Role != domuser.RoleAdmin {
		return nil, domerrors.ErrAdminRequired
	}

	// Validate input
	if err := s.validateUserData(email, firstName, lastName); err != nil {
		return nil, err
	}

	if !domuser.IsValidRole(role) {
		return nil, domerrors.ErrInvalidRole
	}

	// Check if email already exists
	exists, err := s.repo.ExistsByEmail(email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domerrors.ErrEmailAlreadyExists
	}

	// Hash password
	hashedPassword, err := s.hashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &domuser.User{
		ID:            uuid.NewString(),
		Email:         strings.ToLower(strings.TrimSpace(email)),
		PasswordHash:  hashedPassword,
		FirstName:     strings.TrimSpace(firstName),
		LastName:      strings.TrimSpace(lastName),
		Role:          role,
		IsActive:      true,
		EmailVerified: false,
		Profile: domuser.Profile{
			Preferences: domuser.Preferences{
				Language:          "en",
				Timezone:          "UTC",
				NotificationEmail: true,
				NotificationSMS:   false,
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(id string, currentUser *domuser.User) (*domuser.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Authorization check - users can see their own data, admins can see anyone's
	if currentUser.Role != domuser.RoleAdmin && currentUser.ID != user.ID {
		return nil, domerrors.ErrUnauthorized
	}

	return user, nil
}

// GetAllUsers retrieves all users with pagination
func (s *UserService) GetAllUsers(limit, offset int, currentUser *domuser.User) ([]*domuser.User, int, error) {
	// Only admins can list all users
	if currentUser.Role != domuser.RoleAdmin {
		return nil, 0, domerrors.ErrAdminRequired
	}

	// Validate pagination
	if limit <= 0 || limit > 100 {
		limit = 20 // Default limit
	}
	if offset < 0 {
		offset = 0
	}

	users, total, err := s.repo.GetAll(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetUsersByRole retrieves users by role with pagination
func (s *UserService) GetUsersByRole(role domuser.Role, limit, offset int, currentUser *domuser.User) ([]*domuser.User, int, error) {
	// Only admins can filter users by role
	if currentUser.Role != domuser.RoleAdmin {
		return nil, 0, domerrors.ErrAdminRequired
	}

	if !domuser.IsValidRole(role) {
		return nil, 0, domerrors.ErrInvalidRole
	}

	// Validate pagination
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	users, total, err := s.repo.GetByRole(role, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(id string, updates map[string]interface{}, currentUser *domuser.User) (*domuser.User, error) {
	// Get existing user
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Authorization check
	if !currentUser.CanModifyUser(user) {
		return nil, domerrors.ErrUnauthorized
	}

	// Apply updates with validation
	if err := s.applyUserUpdates(user, updates, currentUser); err != nil {
		return nil, err
	}

	// Update in repository
	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUserProfile updates user profile information
func (s *UserService) UpdateUserProfile(id string, profile domuser.Profile, currentUser *domuser.User) (*domuser.User, error) {
	// Get existing user
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Authorization check - users can update their own profile
	if currentUser.ID != user.ID && currentUser.Role != domuser.RoleAdmin {
		return nil, domerrors.ErrUnauthorized
	}

	// Update profile
	user.Profile = profile
	user.UpdatedAt = time.Now()

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// ChangeUserRole changes a user's role
func (s *UserService) ChangeUserRole(id string, newRole domuser.Role, currentUser *domuser.User) (*domuser.User, error) {
	// Only admins can change roles
	if currentUser.Role != domuser.RoleAdmin {
		return nil, domerrors.ErrAdminRequired
	}

	if !domuser.IsValidRole(newRole) {
		return nil, domerrors.ErrInvalidRole
	}

	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	user.Role = newRole
	user.UpdatedAt = time.Now()

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// ToggleUserStatus activates or deactivates a user
func (s *UserService) ToggleUserStatus(id string, isActive bool, currentUser *domuser.User) (*domuser.User, error) {
	// Only admins can change user status
	if currentUser.Role != domuser.RoleAdmin {
		return nil, domerrors.ErrAdminRequired
	}

	// Users cannot deactivate themselves
	if currentUser.ID == id && !isActive {
		return nil, domerrors.ErrCannotDeactivateSelf
	}

	if err := s.repo.ToggleActiveStatus(id, isActive); err != nil {
		return nil, err
	}

	return s.repo.GetByID(id)
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(id string, currentUser *domuser.User) error {
	// Only admins can delete users
	if currentUser.Role != domuser.RoleAdmin {
		return domerrors.ErrAdminRequired
	}

	// Get user to check if it exists and for authorization
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// Check if user can delete the target user
	if !currentUser.CanDeleteUser(user) {
		return domerrors.ErrCannotDeleteSelf
	}

	return s.repo.Delete(id)
}

// ChangePassword changes a user's password
func (s *UserService) ChangePassword(id, oldPassword, newPassword string, currentUser *domuser.User) error {
	// Users can change their own password, admins can change anyone's
	if currentUser.ID != id && currentUser.Role != domuser.RoleAdmin {
		return domerrors.ErrUnauthorized
	}

	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// For non-admin users, verify old password
	if currentUser.Role != domuser.RoleAdmin {
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
			return domerrors.ErrInvalidCredentials
		}
	}

	// Validate new password
	if err := s.validatePassword(newPassword); err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := s.hashPassword(newPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = hashedPassword
	user.UpdatedAt = time.Now()

	return s.repo.Update(user)
}

// Helper methods

func (s *UserService) validateUserData(email, firstName, lastName string) error {
	email = strings.TrimSpace(email)
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)

	if email == "" {
		return domerrors.ErrInvalidEmail
	}

	// Create a temporary user to validate email format
	tempUser := &domuser.User{Email: email}
	if !tempUser.IsValidEmail() {
		return domerrors.ErrInvalidEmail
	}

	if firstName == "" {
		return domerrors.ErrEmptyFirstName
	}

	if lastName == "" {
		return domerrors.ErrEmptyLastName
	}

	return nil
}

func (s *UserService) validatePassword(password string) error {
	if len(password) < 8 {
		return domerrors.ErrPasswordTooWeak
	}
	// Add more password validation rules as needed
	return nil
}

func (s *UserService) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *UserService) applyUserUpdates(user *domuser.User, updates map[string]interface{}, currentUser *domuser.User) error {
	for field, value := range updates {
		switch field {
		case "email":
			if email, ok := value.(string); ok {
				if err := s.validateUserData(email, user.FirstName, user.LastName); err != nil {
					return err
				}
				// Check if new email already exists
				exists, err := s.repo.ExistsByEmail(email)
				if err != nil {
					return err
				}
				if exists && email != user.Email {
					return domerrors.ErrEmailAlreadyExists
				}
				user.Email = strings.ToLower(strings.TrimSpace(email))
				user.EmailVerified = false // Reset verification status
			}
		case "firstName":
			if firstName, ok := value.(string); ok {
				firstName = strings.TrimSpace(firstName)
				if firstName == "" {
					return domerrors.ErrEmptyFirstName
				}
				user.FirstName = firstName
			}
		case "lastName":
			if lastName, ok := value.(string); ok {
				lastName = strings.TrimSpace(lastName)
				if lastName == "" {
					return domerrors.ErrEmptyLastName
				}
				user.LastName = lastName
			}
		case "role":
			// Only admins can change roles
			if currentUser.Role != domuser.RoleAdmin {
				return domerrors.ErrAdminRequired
			}
			if roleStr, ok := value.(string); ok {
				role := domuser.Role(roleStr)
				if !domuser.IsValidRole(role) {
					return domerrors.ErrInvalidRole
				}
				user.Role = role
			}
		case "isActive":
			// Only admins can change active status
			if currentUser.Role != domuser.RoleAdmin {
				return domerrors.ErrAdminRequired
			}
			if isActive, ok := value.(bool); ok {
				// Users cannot deactivate themselves
				if currentUser.ID == user.ID && !isActive {
					return domerrors.ErrCannotDeactivateSelf
				}
				user.IsActive = isActive
			}
		case "emailVerified":
			// Only admins can change email verification status
			if currentUser.Role != domuser.RoleAdmin {
				return domerrors.ErrAdminRequired
			}
			if verified, ok := value.(bool); ok {
				user.EmailVerified = verified
			}
		}
	}

	user.UpdatedAt = time.Now()
	return nil
}
