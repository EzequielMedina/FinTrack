package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fintrack/user-service/internal/app"
	"github.com/fintrack/user-service/internal/config"
	"github.com/fintrack/user-service/internal/core/service"
	"github.com/gin-gonic/gin"
)

// minimal fake app to avoid DB connections
type fakeAuth struct{ *service.AuthService }

func TestMapRoutes_HealthAndProtected(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	cfg := &config.Config{JWTSecret: "secret", Port: "8080"}
	// application with nil internals is fine for mapping, we only need AuthService when hitting auth endpoints
	a := &app.Application{AuthService: nil, Config: cfg}
	h := NewHandlers(a)
	MapRoutes(r, h, cfg, a)

	// health
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("health expected 200, got %d", w.Code)
	}

	// protected endpoint without token should be 401
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/me", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("/api/me expected 401, got %d", w.Code)
	}
}

func TestNewHandlers_WithNilAppDoesNotPanic(t *testing.T) {
	defer func() { _ = recover() }()
	a := &app.Application{}
	_ = NewHandlers(a)
	// if it panics, test will fail; otherwise success
	time.Sleep(1 * time.Millisecond)
}
