package authhandler

import (
	"net/http"

	domerrors "github.com/fintrack/user-service/internal/core/errors"
	"github.com/fintrack/user-service/internal/core/service"
	"github.com/fintrack/user-service/internal/infrastructure/entrypoints/handlers/auth/dto"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	auth *service.AuthService
}

func New(auth *service.AuthService) *Handler { return &Handler{auth: auth} }

func (h *Handler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, at, rt, err := h.auth.Register(req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		status := http.StatusInternalServerError
		if err == domerrors.ErrEmailAlreadyExists {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, dto.AuthResponse{AccessToken: at, RefreshToken: rt})
}

func (h *Handler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, at, rt, err := h.auth.Login(req.Email, req.Password)
	if err != nil {
		status := http.StatusUnauthorized
		if err != domerrors.ErrInvalidCredentials {
			status = http.StatusInternalServerError
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.AuthResponse{AccessToken: at, RefreshToken: rt})
}
