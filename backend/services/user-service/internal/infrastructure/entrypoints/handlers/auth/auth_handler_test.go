package authhandler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	domuser "github.com/fintrack/user-service/internal/core/domain/entities/user"
	userprovider "github.com/fintrack/user-service/internal/core/providers/user"
	"github.com/fintrack/user-service/internal/core/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// mock repo implementing userprovider.Repository for handler tests
type mockRepo struct {
	byEmail    *domuser.User
	byEmailErr error
	createErr  error
}

// Delete implements userprovider.Repository.
func (m *mockRepo) Delete(id string) error {
	panic("unimplemented")
}

// ExistsByEmail implements userprovider.Repository.
func (m *mockRepo) ExistsByEmail(email string) (bool, error) {
	panic("unimplemented")
}

// ExistsByID implements userprovider.Repository.
func (m *mockRepo) ExistsByID(id string) (bool, error) {
	panic("unimplemented")
}

// GetActiveUsers implements userprovider.Repository.
func (m *mockRepo) GetActiveUsers(limit int, offset int) ([]*domuser.User, int, error) {
	panic("unimplemented")
}

// GetAll implements userprovider.Repository.
func (m *mockRepo) GetAll(limit int, offset int) ([]*domuser.User, int, error) {
	panic("unimplemented")
}

// GetByRole implements userprovider.Repository.
func (m *mockRepo) GetByRole(role domuser.Role, limit int, offset int) ([]*domuser.User, int, error) {
	panic("unimplemented")
}

// ToggleActiveStatus implements userprovider.Repository.
func (m *mockRepo) ToggleActiveStatus(id string, isActive bool) error {
	panic("unimplemented")
}

// Update implements userprovider.Repository.
func (m *mockRepo) Update(u *domuser.User) error {
	panic("unimplemented")
}

// UpdateLastLogin implements userprovider.Repository.
func (m *mockRepo) UpdateLastLogin(id string) error {
	panic("unimplemented")
}

func (m *mockRepo) GetByEmail(email string) (*domuser.User, error) { return m.byEmail, m.byEmailErr }
func (m *mockRepo) GetByID(id string) (*domuser.User, error)       { return nil, nil }
func (m *mockRepo) Create(u *domuser.User) error                   { return m.createErr }

// compile-time check
var _ userprovider.Repository = (*mockRepo)(nil)

func setupRouterWithService(repo userprovider.Repository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := service.NewAuthService(repo, "secret", time.Minute, time.Hour)
	h := New(svc)
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
	return r
}

func doJSON(method, path string, body any, r http.Handler) *httptest.ResponseRecorder {
	var buf *bytes.Buffer
	if body != nil {
		b, _ := json.Marshal(body)
		buf = bytes.NewBuffer(b)
	} else {
		buf = bytes.NewBuffer(nil)
	}
	req := httptest.NewRequest(method, path, buf)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestRegister_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupRouterWithService(&mockRepo{})
	// send invalid JSON
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString("{bad json}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestRegister_Conflict(t *testing.T) {
	r := setupRouterWithService(&mockRepo{byEmail: &domuser.User{ID: "1", Email: "a@b.com"}})
	w := doJSON(http.MethodPost, "/register", map[string]string{
		"email":     "a@b.com",
		"password":  "password1",
		"firstName": "Ada",
		"lastName":  "Lovelace",
	}, r)
	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", w.Code)
	}
}

func TestRegister_InternalErrorOnCreate(t *testing.T) {
	r := setupRouterWithService(&mockRepo{byEmail: nil, createErr: errors.New("db down")})
	w := doJSON(http.MethodPost, "/register", map[string]string{
		"email":     "x@y.com",
		"password":  "password1",
		"firstName": "Grace",
		"lastName":  "Hopper",
	}, r)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestRegister_Success(t *testing.T) {
	r := setupRouterWithService(&mockRepo{})
	w := doJSON(http.MethodPost, "/register", map[string]string{
		"email":     "new@user.com",
		"password":  "password1",
		"firstName": "New",
		"lastName":  "User",
	}, r)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	type userInfo struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}
	type response struct {
		AccessToken  string   `json:"accessToken"`
		RefreshToken string   `json:"refreshToken"`
		User         userInfo `json:"user"`
	}
	var resp response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("json: %v", err)
	}
	if resp.AccessToken == "" || resp.RefreshToken == "" {
		t.Fatalf("expected tokens in response")
	}
	if resp.User.Email != "new@user.com" {
		t.Fatalf("expected user email in response")
	}
}

func TestLogin_BadRequest(t *testing.T) {
	r := setupRouterWithService(&mockRepo{})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("{bad json}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestLogin_Unauthorized(t *testing.T) {
	// repo returns nil user
	r := setupRouterWithService(&mockRepo{byEmail: nil})
	w := doJSON(http.MethodPost, "/login", map[string]string{"email": "no@user.com", "password": "anything"}, r)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestLogin_Success(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password1"), bcrypt.DefaultCost)
	repo := &mockRepo{byEmail: &domuser.User{ID: "1", Email: "ok@user.com", PasswordHash: string(hash)}}
	r := setupRouterWithService(repo)
	w := doJSON(http.MethodPost, "/login", map[string]string{"email": "ok@user.com", "password": "password1"}, r)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	type userInfo struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}
	type response struct {
		AccessToken  string   `json:"accessToken"`
		RefreshToken string   `json:"refreshToken"`
		User         userInfo `json:"user"`
	}
	var resp response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("json: %v", err)
	}
	if resp.AccessToken == "" || resp.RefreshToken == "" {
		t.Fatalf("expected tokens in response")
	}
	if resp.User.Email != "ok@user.com" {
		t.Fatalf("expected user email in response")
	}
}
