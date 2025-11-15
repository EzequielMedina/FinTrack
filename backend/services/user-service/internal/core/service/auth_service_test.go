package service

import (
	"errors"
	"strings"
	"testing"
	"time"

	domuser "github.com/fintrack/user-service/internal/core/domain/entities/user"
	domerrors "github.com/fintrack/user-service/internal/core/errors"
)

// mock repository implementing userprovider.Repository
type mockRepo struct {
	byEmail   *domuser.User
	byID      *domuser.User
	getErr    error
	createErr error
	created   *domuser.User
}

func (m *mockRepo) GetByEmail(email string) (*domuser.User, error) { return m.byEmail, m.getErr }
func (m *mockRepo) GetByID(id string) (*domuser.User, error)       { return m.byID, m.getErr }
func (m *mockRepo) Create(u *domuser.User) error                   { m.created = u; return m.createErr }

// Implement the new interface methods for compatibility
func (m *mockRepo) Update(u *domuser.User) error                           { return nil }
func (m *mockRepo) Delete(id string) error                                 { return nil }
func (m *mockRepo) GetAll(limit, offset int) ([]*domuser.User, int, error) { return nil, 0, nil }
func (m *mockRepo) GetByRole(role domuser.Role, limit, offset int) ([]*domuser.User, int, error) {
	return nil, 0, nil
}
func (m *mockRepo) GetActiveUsers(limit, offset int) ([]*domuser.User, int, error) {
	return nil, 0, nil
}
func (m *mockRepo) ExistsByEmail(email string) (bool, error)          { return m.byEmail != nil, nil }
func (m *mockRepo) ExistsByID(id string) (bool, error)                { return m.byID != nil, nil }
func (m *mockRepo) UpdateLastLogin(id string) error                   { return nil }
func (m *mockRepo) ToggleActiveStatus(id string, isActive bool) error { return nil }

func TestRegister_Success(t *testing.T) {
	repo := &mockRepo{byEmail: nil}
	svc := NewAuthService(repo, "secret", time.Minute, time.Hour)
	u, at, rt, err := svc.Register("a@b.com", "pass", "Ada", "Lovelace")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u == nil || u.ID == "" {
		t.Fatalf("expected user with ID")
	}
	if repo.created == nil || repo.created.Email != "a@b.com" {
		t.Fatalf("user not created in repo")
	}
	if at == "" || rt == "" {
		t.Fatalf("expected tokens")
	}
}

func TestRegister_EmailExists(t *testing.T) {
	repo := &mockRepo{byEmail: &domuser.User{ID: "1", Email: "a@b.com"}}
	svc := NewAuthService(repo, "secret", time.Minute, time.Hour)
	_, _, _, err := svc.Register("a@b.com", "pass", "Ada", "Lovelace")
	if !errors.Is(err, domerrors.ErrEmailAlreadyExists) {
		t.Fatalf("expected ErrEmailAlreadyExists")
	}
}

func TestLogin_Success(t *testing.T) {
	// Prepare a user by registering first
	repo := &mockRepo{}
	svc := NewAuthService(repo, "secret", time.Minute, time.Hour)
	u, _, _, err := svc.Register("x@y.com", "s3cret", "Alan", "Turing")
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	// make sure repo returns the stored user by email
	repo.byEmail = repo.created
	got, at, rt, err := svc.Login("x@y.com", "s3cret")
	if err != nil {
		t.Fatalf("unexpected login error: %v", err)
	}
	if got == nil || got.Email != u.Email {
		t.Fatalf("expected same user")
	}
	if at == "" || rt == "" {
		t.Fatalf("expected tokens")
	}
}

func TestLogin_InvalidPassword(t *testing.T) {
	repo := &mockRepo{}
	svc := NewAuthService(repo, "secret", time.Minute, time.Hour)
	// create a user
	_, _, _, err := svc.Register("x@y.com", "right-pass", "Alan", "Turing")
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	repo.byEmail = repo.created
	// wrong password
	_, _, _, err = svc.Login("x@y.com", "wrong")
	if !errors.Is(err, domerrors.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials")
	}
}

func TestCreateTokens_DeterministicFields(t *testing.T) {
	repo := &mockRepo{}
	svc := NewAuthService(repo, "secret", time.Minute, time.Hour)
	u := &domuser.User{ID: "uid-1", Email: "u@e.com", Role: domuser.RoleUser}
	at, rt, err := svc.createTokens(u)
	if err != nil {
		t.Fatalf("tokens: %v", err)
	}
	if !strings.Contains(at, ".") || !strings.Contains(rt, ".") {
		t.Fatalf("tokens not in JWT format")
	}
}
