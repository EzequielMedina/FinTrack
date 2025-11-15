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
// @Summary Create a new user
// @Description Create a new user account (admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateUserRequest true "User creation data"
// @Success 201 {object} dto.UserResponse "User created successfully"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Admin access required"
// @Failure 409 {object} map[string]string "Email already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/users [post]
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

	user, err := h.userService.CreateUser(req.Email, req.Password, req.FirstName, req.LastName, domuser.Role(req.Role), currentUser)
	if err != nil {
		status := h.getErrorStatus(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToUserResponse(user)
	c.JSON(http.StatusCreated, response)
}

// GetUser retrieves a user by ID
// @Summary Get user by ID
// @Description Retrieve a specific user by their ID
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} dto.UserResponse "User retrieved successfully"
// @Failure 400 {object} map[string]string "Invalid user ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/users/{id} [get]
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
// @Summary Get current user
// @Description Retrieve the current authenticated user's information
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.UserResponse "Current user information"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/me [get]
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
// @Summary Get all users
// @Description Retrieve a paginated list of all users
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Number of items per page" default(20)
// @Success 200 {object} dto.UsersListResponse "Users retrieved successfully"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Admin access required"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/users [get]
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
// @Summary Get users by role
// @Description Retrieve a paginated list of users filtered by role
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param role path string true "User role (admin, user)"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Number of items per page" default(20)
// @Success 200 {object} dto.UsersListResponse "Users retrieved successfully"
// @Failure 400 {object} map[string]string "Invalid role"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Admin access required"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/users/role/{role} [get]
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
// @Summary Update user
// @Description Update user information
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body dto.UpdateUserRequest true "User update data"
// @Success 200 {object} dto.UserResponse "User updated successfully"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Access denied"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/users/{id} [put]
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
// @Summary Update user profile
// @Description Update user profile information including address and preferences
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body dto.UpdateProfileRequest true "Profile update data"
// @Success 200 {object} dto.UserResponse "Profile updated successfully"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Access denied"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/users/{id}/profile [put]
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
// @Summary Change user role
// @Description Change a user's role (admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body dto.ChangeRoleRequest true "Role change data"
// @Success 200 {object} dto.UserResponse "Role changed successfully"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Admin access required"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/users/{id}/role [put]
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

	user, err := h.userService.ChangeUserRole(id, domuser.Role(req.Role), currentUser)
	if err != nil {
		status := h.getErrorStatus(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToUserResponse(user)
	c.JSON(http.StatusOK, response)
}

// ToggleUserStatus activates or deactivates a user
// @Summary Toggle user status
// @Description Activate or deactivate a user (admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body dto.ToggleStatusRequest true "Status toggle data"
// @Success 200 {object} dto.UserResponse "Status changed successfully"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Admin access required or cannot deactivate self"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/users/{id}/status [put]
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
// @Summary Delete user
// @Description Delete a user account (admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 204 "User deleted successfully"
// @Failure 400 {object} map[string]string "Invalid user ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Admin access required or cannot delete self"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/users/{id} [delete]
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
// @Summary Change password
// @Description Change a user's password
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body dto.ChangePasswordRequest true "Password change data"
// @Success 200 {object} map[string]string "Password changed successfully"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Unauthorized or invalid old password"
// @Failure 403 {object} map[string]string "Access denied"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/users/{id}/password [put]
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
		domerrors.ErrInvalidPagination, domerrors.ErrCannotCreateAdmin:
		return http.StatusBadRequest
	case domerrors.ErrUserInactive:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
