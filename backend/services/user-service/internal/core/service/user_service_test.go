package service

import (
	"testing"
	"time"

	domuser "github.com/fintrack/user-service/internal/core/domain/entities/user"
	domerrors "github.com/fintrack/user-service/internal/core/errors"
	"github.com/google/uuid"
)

// MockUserRepository implements the userprovider.Repository interface for testing
type MockUserRepository struct {
	users  map[string]*domuser.User
	emails map[string]*domuser.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:  make(map[string]*domuser.User),
		emails: make(map[string]*domuser.User),
	}
}

func (m *MockUserRepository) GetByEmail(email string) (*domuser.User, error) {
	user, exists := m.emails[email]
	if !exists {
		return nil, domerrors.ErrUserNotFound
	}
	return user, nil
}

func (m *MockUserRepository) GetByID(id string) (*domuser.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, domerrors.ErrUserNotFound
	}
	return user, nil
}

func (m *MockUserRepository) Create(u *domuser.User) error {
	if _, exists := m.emails[u.Email]; exists {
		return domerrors.ErrEmailAlreadyExists
	}
	m.users[u.ID] = u
	m.emails[u.Email] = u
	return nil
}

func (m *MockUserRepository) Update(u *domuser.User) error {
	if _, exists := m.users[u.ID]; !exists {
		return domerrors.ErrUserNotFound
	}

	// Remove old email mapping if email changed
	for email, user := range m.emails {
		if user.ID == u.ID && email != u.Email {
			delete(m.emails, email)
			break
		}
	}

	m.users[u.ID] = u
	m.emails[u.Email] = u
	return nil
}

func (m *MockUserRepository) Delete(id string) error {
	user, exists := m.users[id]
	if !exists {
		return domerrors.ErrUserNotFound
	}
	delete(m.users, id)
	delete(m.emails, user.Email)
	return nil
}

func (m *MockUserRepository) GetAll(limit, offset int) ([]*domuser.User, int, error) {
	var users []*domuser.User
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, len(users), nil
}

func (m *MockUserRepository) GetByRole(role domuser.Role, limit, offset int) ([]*domuser.User, int, error) {
	var users []*domuser.User
	for _, user := range m.users {
		if user.Role == role {
			users = append(users, user)
		}
	}
	return users, len(users), nil
}

func (m *MockUserRepository) GetActiveUsers(limit, offset int) ([]*domuser.User, int, error) {
	var users []*domuser.User
	for _, user := range m.users {
		if user.IsActive {
			users = append(users, user)
		}
	}
	return users, len(users), nil
}

func (m *MockUserRepository) ExistsByEmail(email string) (bool, error) {
	_, exists := m.emails[email]
	return exists, nil
}

func (m *MockUserRepository) ExistsByID(id string) (bool, error) {
	_, exists := m.users[id]
	return exists, nil
}

func (m *MockUserRepository) UpdateLastLogin(id string) error {
	user, exists := m.users[id]
	if !exists {
		return domerrors.ErrUserNotFound
	}
	now := time.Now()
	user.LastLoginAt = &now
	return nil
}

func (m *MockUserRepository) ToggleActiveStatus(id string, isActive bool) error {
	user, exists := m.users[id]
	if !exists {
		return domerrors.ErrUserNotFound
	}
	user.IsActive = isActive
	return nil
}

// Test cases

func TestUserService_CreateUser(t *testing.T) {
	repo := NewMockUserRepository()
	service := NewUserService(repo)

	// Create admin user for authorization
	adminUser := &domuser.User{
		ID:        uuid.NewString(),
		Email:     "admin@test.com",
		Role:      domuser.RoleAdmin,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.Create(adminUser)

	tests := []struct {
		name        string
		email       string
		password    string
		firstName   string
		lastName    string
		role        domuser.Role
		currentUser *domuser.User
		wantErr     error
	}{
		{
			name:        "successful user creation",
			email:       "test@example.com",
			password:    "password123",
			firstName:   "Test",
			lastName:    "User",
			role:        domuser.RoleUser,
			currentUser: adminUser,
			wantErr:     nil,
		},
		{
			name:        "non-admin cannot create user",
			email:       "test2@example.com",
			password:    "password123",
			firstName:   "Test",
			lastName:    "User",
			role:        domuser.RoleUser,
			currentUser: &domuser.User{Role: domuser.RoleUser},
			wantErr:     domerrors.ErrAdminRequired,
		},
		{
			name:        "invalid email",
			email:       "invalid-email",
			password:    "password123",
			firstName:   "Test",
			lastName:    "User",
			role:        domuser.RoleUser,
			currentUser: adminUser,
			wantErr:     domerrors.ErrInvalidEmail,
		},
		{
			name:        "empty first name",
			email:       "test3@example.com",
			password:    "password123",
			firstName:   "",
			lastName:    "User",
			role:        domuser.RoleUser,
			currentUser: adminUser,
			wantErr:     domerrors.ErrEmptyFirstName,
		},
		{
			name:        "invalid role",
			email:       "test4@example.com",
			password:    "password123",
			firstName:   "Test",
			lastName:    "User",
			role:        domuser.Role("invalid"),
			currentUser: adminUser,
			wantErr:     domerrors.ErrInvalidRole,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := service.CreateUser(tt.email, tt.password, tt.firstName, tt.lastName, tt.role, tt.currentUser)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("CreateUser() unexpected error = %v", err)
				return
			}

			if user == nil {
				t.Error("CreateUser() returned nil user")
				return
			}

			if user.Email != tt.email {
				t.Errorf("CreateUser() email = %v, want %v", user.Email, tt.email)
			}

			if user.Role != tt.role {
				t.Errorf("CreateUser() role = %v, want %v", user.Role, tt.role)
			}
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	repo := NewMockUserRepository()
	service := NewUserService(repo)

	// Create test users
	adminUser := &domuser.User{
		ID:        uuid.NewString(),
		Email:     "admin@test.com",
		FirstName: "Admin",
		LastName:  "User",
		Role:      domuser.RoleAdmin,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.Create(adminUser)

	regularUser := &domuser.User{
		ID:        uuid.NewString(),
		Email:     "user@test.com",
		FirstName: "Regular",
		LastName:  "User",
		Role:      domuser.RoleUser,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.Create(regularUser)

	tests := []struct {
		name        string
		userID      string
		updates     map[string]interface{}
		currentUser *domuser.User
		wantErr     error
	}{
		{
			name:   "admin can update any user",
			userID: regularUser.ID,
			updates: map[string]interface{}{
				"firstName": "UpdatedFirst",
				"lastName":  "UpdatedLast",
			},
			currentUser: adminUser,
			wantErr:     nil,
		},
		{
			name:   "user can update themselves",
			userID: regularUser.ID,
			updates: map[string]interface{}{
				"firstName": "SelfUpdated",
			},
			currentUser: regularUser,
			wantErr:     nil,
		},
		{
			name:   "user cannot update others",
			userID: adminUser.ID,
			updates: map[string]interface{}{
				"firstName": "Hacked",
			},
			currentUser: regularUser,
			wantErr:     domerrors.ErrUnauthorized,
		},
		{
			name:   "non-admin cannot change role",
			userID: regularUser.ID,
			updates: map[string]interface{}{
				"role": string(domuser.RoleAdmin),
			},
			currentUser: regularUser,
			wantErr:     domerrors.ErrAdminRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := service.UpdateUser(tt.userID, tt.updates, tt.currentUser)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateUser() unexpected error = %v", err)
				return
			}

			if user == nil {
				t.Error("UpdateUser() returned nil user")
				return
			}

			// Verify updates were applied
			if firstName, ok := tt.updates["firstName"]; ok {
				if user.FirstName != firstName {
					t.Errorf("UpdateUser() firstName = %v, want %v", user.FirstName, firstName)
				}
			}
		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	repo := NewMockUserRepository()
	service := NewUserService(repo)

	// Create test users
	adminUser := &domuser.User{
		ID:        uuid.NewString(),
		Email:     "admin@test.com",
		Role:      domuser.RoleAdmin,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.Create(adminUser)

	regularUser := &domuser.User{
		ID:        uuid.NewString(),
		Email:     "user@test.com",
		Role:      domuser.RoleUser,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.Create(regularUser)

	tests := []struct {
		name        string
		userID      string
		currentUser *domuser.User
		wantErr     error
	}{
		{
			name:        "admin can delete other users",
			userID:      regularUser.ID,
			currentUser: adminUser,
			wantErr:     nil,
		},
		{
			name:        "admin cannot delete themselves",
			userID:      adminUser.ID,
			currentUser: adminUser,
			wantErr:     domerrors.ErrCannotDeleteSelf,
		},
		{
			name:        "non-admin cannot delete users",
			userID:      adminUser.ID,
			currentUser: regularUser,
			wantErr:     domerrors.ErrAdminRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.DeleteUser(tt.userID, tt.currentUser)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("DeleteUser() unexpected error = %v", err)
				return
			}

			// Verify user was deleted
			_, err = repo.GetByID(tt.userID)
			if err != domerrors.ErrUserNotFound {
				t.Error("DeleteUser() user was not deleted")
			}
		})
	}
}
