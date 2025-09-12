package userhandler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	domuser "github.com/fintrack/user-service/internal/core/domain/entities/user"
	domerrors "github.com/fintrack/user-service/internal/core/errors"
	"github.com/fintrack/user-service/internal/core/service"
	"github.com/fintrack/user-service/internal/infrastructure/entrypoints/handlers/user/dto"
)

// Handler handles HTTP requests for user management operations
type Handler struct {
	userService service.UserServiceInterface
}

// New creates a new UserHandler instance
func New(userService service.UserServiceInterface) *Handler {
	return &Handler{userService: userService}
}

// CreateUser creates a new user
// POST /api/users
func (h *Handler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUser := h.getCurrentUser(c)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.userService.CreateUser(req.Email, req.Password, req.FirstName, req.LastName, req.Role, currentUser)
	if err != nil {
		status := h.getErrorStatus(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToUserResponse(user)
	c.JSON(http.StatusCreated, response)
}

// GetUser retrieves a user by ID
// GET /api/users/:id
func (h *Handler) GetUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	currentUser := h.getCurrentUser(c)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.userService.GetUserByID(id, currentUser)
	if err != nil {
		status := h.getErrorStatus(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToUserResponse(user)
	c.JSON(http.StatusOK, response)
}

// GetCurrentUser retrieves the current authenticated user's information
// GET /api/users/me
func (h *Handler) GetCurrentUser(c *gin.Context) {
	currentUser := h.getCurrentUser(c)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Get fresh user data from service
	user, err := h.userService.GetUserByID(currentUser.ID, currentUser)
	if err != nil {
		status := h.getErrorStatus(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToUserResponse(user)
	c.JSON(http.StatusOK, response)
}

// GetAllUsers retrieves all users with pagination
// GET /api/users?page=1&pageSize=20
func (h *Handler) GetAllUsers(c *gin.Context) {
	currentUser := h.getCurrentUser(c)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Parse pagination parameters
	page, pageSize, limit, offset := h.getPaginationParams(c)

	users, total, err := h.userService.GetAllUsers(limit, offset, currentUser)
	if err != nil {
		status := h.getErrorStatus(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToUsersListResponse(users, total, page, pageSize)
	c.JSON(http.StatusOK, response)
}

// GetUsersByRole retrieves users by role with pagination
// GET /api/users/role/:role?page=1&pageSize=20
func (h *Handler) GetUsersByRole(c *gin.Context) {
	roleParam := c.Param("role")
	if roleParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role is required"})
		return
	}

	role := domuser.Role(roleParam)
	if !domuser.IsValidRole(role) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
		return
	}

	currentUser := h.getCurrentUser(c)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Parse pagination parameters
	page, pageSize, limit, offset := h.getPaginationParams(c)

	users, total, err := h.userService.GetUsersByRole(role, limit, offset, currentUser)
	if err != nil {
		status := h.getErrorStatus(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToUsersListResponse(users, total, page, pageSize)
	c.JSON(http.StatusOK, response)
}

// UpdateUser updates user information
// PUT /api/users/:id
func (h *Handler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUser := h.getCurrentUser(c)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	updates := req.ToUpdateMap()
	user, err := h.userService.UpdateUser(id, updates, currentUser)
	if err != nil {
		status := h.getErrorStatus(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToUserResponse(user)
	c.JSON(http.StatusOK, response)
}

// UpdateUserProfile updates user profile information
// PUT /api/users/:id/profile
func (h *Handler) UpdateUserProfile(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUser := h.getCurrentUser(c)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Get existing user to merge profile data
	existingUser, err := h.userService.GetUserByID(id, currentUser)
	if err != nil {
		status := h.getErrorStatus(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	profile := req.ToProfile(existingUser.Profile)
	user, err := h.userService.UpdateUserProfile(id, profile, currentUser)
	if err != nil {
		status := h.getErrorStatus(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToUserResponse(user)
	c.JSON(http.StatusOK, response)
}

// ChangeUserRole changes a user's role
// PUT /api/users/:id/role
func (h *Handler) ChangeUserRole(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	var req dto.ChangeRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUser := h.getCurrentUser(c)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.userService.ChangeUserRole(id, req.Role, currentUser)
	if err != nil {
		status := h.getErrorStatus(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToUserResponse(user)
	c.JSON(http.StatusOK, response)
}

// ToggleUserStatus activates or deactivates a user
// PUT /api/users/:id/status
func (h *Handler) ToggleUserStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	var req dto.ToggleStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUser := h.getCurrentUser(c)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.userService.ToggleUserStatus(id, req.IsActive, currentUser)
	if err != nil {
		status := h.getErrorStatus(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToUserResponse(user)
	c.JSON(http.StatusOK, response)
}

// DeleteUser deletes a user
// DELETE /api/users/:id
func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	currentUser := h.getCurrentUser(c)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	err := h.userService.DeleteUser(id, currentUser)
	if err != nil {
		status := h.getErrorStatus(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ChangePassword changes a user's password
// PUT /api/users/:id/password
func (h *Handler) ChangePassword(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUser := h.getCurrentUser(c)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	err := h.userService.ChangePassword(id, req.OldPassword, req.NewPassword, currentUser)
	if err != nil {
		status := h.getErrorStatus(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
}

// Helper methods

// getCurrentUser extracts the current user from the context (set by auth middleware)
func (h *Handler) getCurrentUser(c *gin.Context) *domuser.User {
	user, exists := c.Get("currentUser")
	if !exists {
		return nil
	}

	domainUser, ok := user.(*domuser.User)
	if !ok {
		return nil
	}

	return domainUser
}

// getPaginationParams extracts and validates pagination parameters
func (h *Handler) getPaginationParams(c *gin.Context) (page, pageSize, limit, offset int) {
	page = 1
	pageSize = 20

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if sizeStr := c.Query("pageSize"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 && s <= 100 {
			pageSize = s
		}
	}

	limit = pageSize
	offset = (page - 1) * pageSize

	return page, pageSize, limit, offset
}

// getErrorStatus maps domain errors to HTTP status codes
func (h *Handler) getErrorStatus(err error) int {
	switch err {
	case domerrors.ErrUserNotFound:
		return http.StatusNotFound
	case domerrors.ErrEmailAlreadyExists:
		return http.StatusConflict
	case domerrors.ErrInvalidCredentials:
		return http.StatusUnauthorized
	case domerrors.ErrUnauthorized:
		return http.StatusUnauthorized
	case domerrors.ErrAdminRequired, domerrors.ErrForbidden:
		return http.StatusForbidden
	case domerrors.ErrCannotDeleteSelf, domerrors.ErrCannotDeactivateSelf:
		return http.StatusForbidden
	case domerrors.ErrInvalidRole, domerrors.ErrInvalidEmail, domerrors.ErrPasswordTooWeak,
		domerrors.ErrInvalidUserData, domerrors.ErrEmptyFirstName, domerrors.ErrEmptyLastName,
		domerrors.ErrInvalidPagination:
		return http.StatusBadRequest
	case domerrors.ErrUserInactive:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
