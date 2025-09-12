package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	domuser "github.com/fintrack/user-service/internal/core/domain/entities/user"
	domerrors "github.com/fintrack/user-service/internal/core/errors"
	userprovider "github.com/fintrack/user-service/internal/core/providers/user"
	"github.com/google/uuid"
)

type AuthService struct {
	repo          userprovider.Repository
	jwtSecret     string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

func NewAuthService(repo userprovider.Repository, secret string, accessExp, refreshExp time.Duration) *AuthService {
	return &AuthService{repo: repo, jwtSecret: secret, accessExpiry: accessExp, refreshExpiry: refreshExp}
}

func (s *AuthService) Register(email, password, firstName, lastName string) (*domuser.User, string, string, error) {
	// check existing
	if existing, _ := s.repo.GetByEmail(email); existing != nil {
		return nil, "", "", domerrors.ErrEmailAlreadyExists
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", "", err
	}
	u := &domuser.User{
		ID:            uuid.NewString(),
		Email:         email,
		PasswordHash:  string(hash),
		FirstName:     firstName,
		LastName:      lastName,
		Role:          domuser.RoleAdmin,
		IsActive:      true,
		EmailVerified: false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := s.repo.Create(u); err != nil {
		return nil, "", "", err
	}
	access, refresh, err := s.createTokens(u)
	if err != nil {
		return nil, "", "", err
	}
	return u, access, refresh, nil
}

func (s *AuthService) Login(email, password string) (*domuser.User, string, string, error) {
	u, err := s.repo.GetByEmail(email)
	if err != nil || u == nil {
		return nil, "", "", domerrors.ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, "", "", domerrors.ErrInvalidCredentials
	}
	access, refresh, err := s.createTokens(u)
	if err != nil {
		return nil, "", "", err
	}
	return u, access, refresh, nil
}

func (s *AuthService) createTokens(u *domuser.User) (string, string, error) {
	now := time.Now()
	accessClaims := jwt.MapClaims{
		"sub":   u.ID,
		"email": u.Email,
		"role":  string(u.Role),
		"iat":   now.Unix(),
		"exp":   now.Add(s.accessExpiry).Unix(),
		"iss":   "fintrack-user-service",
		"aud":   "fintrack-api",
	}
	refreshClaims := jwt.MapClaims{
		"sub":  u.ID,
		"type": "refresh",
		"iat":  now.Unix(),
		"exp":  now.Add(s.refreshExpiry).Unix(),
	}
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	at, err := access.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", "", err
	}
	rt, err := refresh.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", "", err
	}
	return at, rt, nil
}
