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
		{
			name:        "admin cannot create other admin users",
			email:       "admin2@example.com",
			password:    "password123",
			firstName:   "Another",
			lastName:    "Admin",
			role:        domuser.RoleAdmin,
			currentUser: adminUser,
			wantErr:     domerrors.ErrCannotCreateAdmin,
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

func TestUserService_GetAllUsers(t *testing.T) {
	repo := NewMockUserRepository()
	service := NewUserService(repo)

	// Create multiple test users including two admins
	admin1 := &domuser.User{
		ID:        uuid.NewString(),
		Email:     "admin1@test.com",
		FirstName: "Admin",
		LastName:  "One",
		Role:      domuser.RoleAdmin,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.Create(admin1)

	admin2 := &domuser.User{
		ID:        uuid.NewString(),
		Email:     "admin2@test.com",
		FirstName: "Admin",
		LastName:  "Two",
		Role:      domuser.RoleAdmin,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.Create(admin2)

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

	operatorUser := &domuser.User{
		ID:        uuid.NewString(),
		Email:     "operator@test.com",
		FirstName: "Operator",
		LastName:  "User",
		Role:      domuser.RoleOperator,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.Create(operatorUser)

	tests := []struct {
		name          string
		currentUser   *domuser.User
		expectedCount int
		shouldInclude []string // User IDs that should be included
		shouldExclude []string // User IDs that should be excluded
		wantErr       error
	}{
		{
			name:          "admin1 sees only themselves and non-admin users",
			currentUser:   admin1,
			expectedCount: 3, // admin1, regularUser, operatorUser (not admin2)
			shouldInclude: []string{admin1.ID, regularUser.ID, operatorUser.ID},
			shouldExclude: []string{admin2.ID},
			wantErr:       nil,
		},
		{
			name:          "admin2 sees only themselves and non-admin users",
			currentUser:   admin2,
			expectedCount: 3, // admin2, regularUser, operatorUser (not admin1)
			shouldInclude: []string{admin2.ID, regularUser.ID, operatorUser.ID},
			shouldExclude: []string{admin1.ID},
			wantErr:       nil,
		},
		{
			name:        "non-admin cannot access GetAllUsers",
			currentUser: regularUser,
			wantErr:     domerrors.ErrAdminRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users, total, err := service.GetAllUsers(20, 0, tt.currentUser)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("GetAllUsers() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("GetAllUsers() unexpected error = %v", err)
				return
			}

			if total != tt.expectedCount {
				t.Errorf("GetAllUsers() total = %v, want %v", total, tt.expectedCount)
			}

			if len(users) != tt.expectedCount {
				t.Errorf("GetAllUsers() users count = %v, want %v", len(users), tt.expectedCount)
			}

			// Check that expected users are included
			foundUsers := make(map[string]bool)
			for _, user := range users {
				foundUsers[user.ID] = true
			}

			for _, expectedID := range tt.shouldInclude {
				if !foundUsers[expectedID] {
					t.Errorf("GetAllUsers() should include user %v but didn't", expectedID)
				}
			}

			for _, excludedID := range tt.shouldExclude {
				if foundUsers[excludedID] {
					t.Errorf("GetAllUsers() should exclude user %v but included it", excludedID)
				}
			}
		})
	}
}

func TestUserService_GetUsersByRole(t *testing.T) {
	repo := NewMockUserRepository()
	service := NewUserService(repo)

	// Create multiple test users including two admins
	admin1 := &domuser.User{
		ID:        uuid.NewString(),
		Email:     "admin1@test.com",
		FirstName: "Admin",
		LastName:  "One",
		Role:      domuser.RoleAdmin,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.Create(admin1)

	admin2 := &domuser.User{
		ID:        uuid.NewString(),
		Email:     "admin2@test.com",
		FirstName: "Admin",
		LastName:  "Two",
		Role:      domuser.RoleAdmin,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.Create(admin2)

	user1 := &domuser.User{
		ID:        uuid.NewString(),
		Email:     "user1@test.com",
		FirstName: "User",
		LastName:  "One",
		Role:      domuser.RoleUser,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.Create(user1)

	user2 := &domuser.User{
		ID:        uuid.NewString(),
		Email:     "user2@test.com",
		FirstName: "User",
		LastName:  "Two",
		Role:      domuser.RoleUser,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.Create(user2)

	tests := []struct {
		name          string
		role          domuser.Role
		currentUser   *domuser.User
		expectedCount int
		shouldInclude []string // User IDs that should be included
		shouldExclude []string // User IDs that should be excluded
		wantErr       error
	}{
		{
			name:          "admin1 searches for admin role - only sees themselves",
			role:          domuser.RoleAdmin,
			currentUser:   admin1,
			expectedCount: 1,
			shouldInclude: []string{admin1.ID},
			shouldExclude: []string{admin2.ID},
			wantErr:       nil,
		},
		{
			name:          "admin2 searches for admin role - only sees themselves",
			role:          domuser.RoleAdmin,
			currentUser:   admin2,
			expectedCount: 1,
			shouldInclude: []string{admin2.ID},
			shouldExclude: []string{admin1.ID},
			wantErr:       nil,
		},
		{
			name:          "admin searches for user role - sees all users",
			role:          domuser.RoleUser,
			currentUser:   admin1,
			expectedCount: 2,
			shouldInclude: []string{user1.ID, user2.ID},
			shouldExclude: []string{admin1.ID, admin2.ID},
			wantErr:       nil,
		},
		{
			name:        "non-admin cannot access GetUsersByRole",
			role:        domuser.RoleUser,
			currentUser: user1,
			wantErr:     domerrors.ErrAdminRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users, total, err := service.GetUsersByRole(tt.role, 20, 0, tt.currentUser)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("GetUsersByRole() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("GetUsersByRole() unexpected error = %v", err)
				return
			}

			if total != tt.expectedCount {
				t.Errorf("GetUsersByRole() total = %v, want %v", total, tt.expectedCount)
			}

			if len(users) != tt.expectedCount {
				t.Errorf("GetUsersByRole() users count = %v, want %v", len(users), tt.expectedCount)
			}

			// Check that expected users are included
			foundUsers := make(map[string]bool)
			for _, user := range users {
				foundUsers[user.ID] = true
			}

			for _, expectedID := range tt.shouldInclude {
				if !foundUsers[expectedID] {
					t.Errorf("GetUsersByRole() should include user %v but didn't", expectedID)
				}
			}

			for _, excludedID := range tt.shouldExclude {
				if foundUsers[excludedID] {
					t.Errorf("GetUsersByRole() should exclude user %v but included it", excludedID)
				}
			}
		})
	}
}
