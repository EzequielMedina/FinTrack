package userhandler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	domuser "github.com/fintrack/user-service/internal/core/domain/entities/user"
	"github.com/fintrack/user-service/internal/infrastructure/entrypoints/handlers/user/dto"
	"github.com/google/uuid"
)

// MockUserService implements a mock UserService for testing
type MockUserService struct {
	users map[string]*domuser.User
}

func NewMockUserService() *MockUserService {
	return &MockUserService{
		users: make(map[string]*domuser.User),
	}
}

func (m *MockUserService) CreateUser(email, password, firstName, lastName string, role domuser.Role, currentUser *domuser.User) (*domuser.User, error) {
	user := &domuser.User{
		ID:        uuid.NewString(),
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m.users[user.ID] = user
	return user, nil
}

func (m *MockUserService) GetUserByID(id string, currentUser *domuser.User) (*domuser.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, nil // Simplified for testing
	}
	return user, nil
}

func (m *MockUserService) GetAllUsers(limit, offset int, currentUser *domuser.User) ([]*domuser.User, int, error) {
	var users []*domuser.User
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, len(users), nil
}

func (m *MockUserService) GetUsersByRole(role domuser.Role, limit, offset int, currentUser *domuser.User) ([]*domuser.User, int, error) {
	var users []*domuser.User
	for _, user := range m.users {
		if user.Role == role {
			users = append(users, user)
		}
	}
	return users, len(users), nil
}

func (m *MockUserService) UpdateUser(id string, updates map[string]interface{}, currentUser *domuser.User) (*domuser.User, error) {
	user := m.users[id]
	if user == nil {
		return nil, nil
	}

	// Apply some basic updates for testing
	if firstName, ok := updates["firstName"]; ok {
		user.FirstName = firstName.(string)
	}
	if lastName, ok := updates["lastName"]; ok {
		user.LastName = lastName.(string)
	}

	user.UpdatedAt = time.Now()
	return user, nil
}

func (m *MockUserService) UpdateUserProfile(id string, profile domuser.Profile, currentUser *domuser.User) (*domuser.User, error) {
	user := m.users[id]
	if user == nil {
		return nil, nil
	}
	user.Profile = profile
	user.UpdatedAt = time.Now()
	return user, nil
}

func (m *MockUserService) ChangeUserRole(id string, newRole domuser.Role, currentUser *domuser.User) (*domuser.User, error) {
	user := m.users[id]
	if user == nil {
		return nil, nil
	}
	user.Role = newRole
	user.UpdatedAt = time.Now()
	return user, nil
}

func (m *MockUserService) ToggleUserStatus(id string, isActive bool, currentUser *domuser.User) (*domuser.User, error) {
	user := m.users[id]
	if user == nil {
		return nil, nil
	}
	user.IsActive = isActive
	user.UpdatedAt = time.Now()
	return user, nil
}

func (m *MockUserService) DeleteUser(id string, currentUser *domuser.User) error {
	delete(m.users, id)
	return nil
}

func (m *MockUserService) ChangePassword(id, oldPassword, newPassword string, currentUser *domuser.User) error {
	return nil // Simplified for testing
}

func TestHandler_CreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := NewMockUserService()
	handler := New(mockService)

	// Create admin user for context
	adminUser := &domuser.User{
		ID:       uuid.NewString(),
		Email:    "admin@test.com",
		Role:     domuser.RoleAdmin,
		IsActive: true,
	}

	tests := []struct {
		name           string
		requestBody    dto.CreateUserRequest
		currentUser    *domuser.User
		expectedStatus int
	}{
		{
			name: "successful user creation",
			requestBody: dto.CreateUserRequest{
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
				Role:      domuser.RoleUser,
			},
			currentUser:    adminUser,
			expectedStatus: http.StatusCreated,
		},
		{
			name: "invalid email format",
			requestBody: dto.CreateUserRequest{
				Email:     "invalid-email",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
				Role:      domuser.RoleUser,
			},
			currentUser:    adminUser,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "unauthorized without user",
			requestBody: dto.CreateUserRequest{
				Email:     "test2@example.com",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
				Role:      domuser.RoleUser,
			},
			currentUser:    nil,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Create Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Set current user in context if provided
			if tt.currentUser != nil {
				c.Set("currentUser", tt.currentUser)
			}

			// Call handler
			handler.CreateUser(c)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("CreateUser() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			// For successful creation, check response structure
			if tt.expectedStatus == http.StatusCreated {
				var response dto.UserResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("CreateUser() failed to unmarshal response: %v", err)
				}

				if response.Email != tt.requestBody.Email {
					t.Errorf("CreateUser() response email = %v, want %v", response.Email, tt.requestBody.Email)
				}
			}
		})
	}
}

func TestHandler_GetCurrentUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := NewMockUserService()
	handler := New(mockService)

	// Create test user
	testUser := &domuser.User{
		ID:        uuid.NewString(),
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Role:      domuser.RoleUser,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockService.users[testUser.ID] = testUser

	tests := []struct {
		name           string
		currentUser    *domuser.User
		expectedStatus int
	}{
		{
			name:           "successful get current user",
			currentUser:    testUser,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "unauthorized without user",
			currentUser:    nil,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req := httptest.NewRequest(http.MethodGet, "/api/me", nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Create Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Set current user in context if provided
			if tt.currentUser != nil {
				c.Set("currentUser", tt.currentUser)
			}

			// Call handler
			handler.GetCurrentUser(c)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("GetCurrentUser() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			// For successful request, check response structure
			if tt.expectedStatus == http.StatusOK {
				var response dto.UserResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("GetCurrentUser() failed to unmarshal response: %v", err)
				}

				if response.ID != testUser.ID {
					t.Errorf("GetCurrentUser() response ID = %v, want %v", response.ID, testUser.ID)
				}
			}
		})
	}
}
