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

// Register godoc
// @Summary Register a new user
// @Description Create a new user account with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User registration data"
// @Success 201 {object} dto.AuthResponse "User created successfully"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 409 {object} map[string]string "Email already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, at, rt, err := h.auth.Register(req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		status := http.StatusInternalServerError
		if err == domerrors.ErrEmailAlreadyExists {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, dto.AuthResponse{
		AccessToken:  at,
		RefreshToken: rt,
		User: dto.UserInfo{
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
	})
}

// Login godoc
// @Summary Authenticate user
// @Description Authenticate user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "User login credentials"
// @Success 200 {object} dto.AuthResponse "Authentication successful"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, at, rt, err := h.auth.Login(req.Email, req.Password)
	if err != nil {
		status := http.StatusUnauthorized
		if err != domerrors.ErrInvalidCredentials {
			status = http.StatusInternalServerError
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.AuthResponse{
		AccessToken:  at,
		RefreshToken: rt,
		User: dto.UserInfo{
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
	})
}
