package dto

import (
	"time"

	domuser "github.com/fintrack/user-service/internal/core/domain/entities/user"
)

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	Email     string       `json:"email" binding:"required,email"`
	Password  string       `json:"password" binding:"required,min=8"`
	FirstName string       `json:"firstName" binding:"required"`
	LastName  string       `json:"lastName" binding:"required"`
	Role      domuser.Role `json:"role" binding:"required"`
}

// UpdateUserRequest represents the request to update user information
type UpdateUserRequest struct {
	Email         *string       `json:"email,omitempty" binding:"omitempty,email"`
	FirstName     *string       `json:"firstName,omitempty"`
	LastName      *string       `json:"lastName,omitempty"`
	Role          *domuser.Role `json:"role,omitempty"`
	IsActive      *bool         `json:"isActive,omitempty"`
	EmailVerified *bool         `json:"emailVerified,omitempty"`
}

// UpdateProfileRequest represents the request to update user profile
type UpdateProfileRequest struct {
	Phone       *string             `json:"phone,omitempty"`
	DateOfBirth *time.Time          `json:"dateOfBirth,omitempty"`
	Address     *AddressRequest     `json:"address,omitempty"`
	Preferences *PreferencesRequest `json:"preferences,omitempty"`
}

// AddressRequest represents address information in requests
type AddressRequest struct {
	Street     string `json:"street,omitempty"`
	City       string `json:"city,omitempty"`
	State      string `json:"state,omitempty"`
	PostalCode string `json:"postalCode,omitempty"`
	Country    string `json:"country,omitempty"`
}

// PreferencesRequest represents user preferences in requests
type PreferencesRequest struct {
	Language          *string `json:"language,omitempty"`
	Timezone          *string `json:"timezone,omitempty"`
	NotificationEmail *bool   `json:"notificationEmail,omitempty"`
	NotificationSMS   *bool   `json:"notificationSMS,omitempty"`
}

// ChangePasswordRequest represents the request to change password
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=8"`
}

// ChangeRoleRequest represents the request to change user role
type ChangeRoleRequest struct {
	Role domuser.Role `json:"role" binding:"required"`
}

// ToggleStatusRequest represents the request to toggle user status
type ToggleStatusRequest struct {
	IsActive bool `json:"isActive"`
}

// UserResponse represents user information in responses
type UserResponse struct {
	ID            string          `json:"id"`
	Email         string          `json:"email"`
	FirstName     string          `json:"firstName"`
	LastName      string          `json:"lastName"`
	FullName      string          `json:"fullName"`
	Role          domuser.Role    `json:"role"`
	IsActive      bool            `json:"isActive"`
	EmailVerified bool            `json:"emailVerified"`
	Profile       ProfileResponse `json:"profile"`
	CreatedAt     time.Time       `json:"createdAt"`
	UpdatedAt     time.Time       `json:"updatedAt"`
	LastLoginAt   *time.Time      `json:"lastLoginAt,omitempty"`
}

// ProfileResponse represents profile information in responses
type ProfileResponse struct {
	Phone          string              `json:"phone,omitempty"`
	DateOfBirth    *time.Time          `json:"dateOfBirth,omitempty"`
	Address        AddressResponse     `json:"address"`
	ProfilePicture string              `json:"profilePicture,omitempty"`
	Preferences    PreferencesResponse `json:"preferences"`
}

// AddressResponse represents address information in responses
type AddressResponse struct {
	Street     string `json:"street,omitempty"`
	City       string `json:"city,omitempty"`
	State      string `json:"state,omitempty"`
	PostalCode string `json:"postalCode,omitempty"`
	Country    string `json:"country,omitempty"`
}

// PreferencesResponse represents user preferences in responses
type PreferencesResponse struct {
	Language          string `json:"language,omitempty"`
	Timezone          string `json:"timezone,omitempty"`
	NotificationEmail bool   `json:"notificationEmail"`
	NotificationSMS   bool   `json:"notificationSMS"`
}

// UsersListResponse represents a paginated list of users
type UsersListResponse struct {
	Users      []UserResponse `json:"users"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"pageSize"`
	TotalPages int            `json:"totalPages"`
}

// UserStatsResponse represents user statistics
type UserStatsResponse struct {
	TotalUsers    int `json:"totalUsers"`
	ActiveUsers   int `json:"activeUsers"`
	InactiveUsers int `json:"inactiveUsers"`
	AdminUsers    int `json:"adminUsers"`
	RegularUsers  int `json:"regularUsers"`
}

// ToUserResponse converts a domain User to a UserResponse DTO
func ToUserResponse(user *domuser.User) UserResponse {
	return UserResponse{
		ID:            user.ID,
		Email:         user.Email,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		FullName:      user.GetFullName(),
		Role:          user.Role,
		IsActive:      user.IsActive,
		EmailVerified: user.EmailVerified,
		Profile:       ToProfileResponse(user.Profile),
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		LastLoginAt:   user.LastLoginAt,
	}
}

// ToProfileResponse converts a domain Profile to a ProfileResponse DTO
func ToProfileResponse(profile domuser.Profile) ProfileResponse {
	return ProfileResponse{
		Phone:          profile.Phone,
		DateOfBirth:    profile.DateOfBirth,
		Address:        ToAddressResponse(profile.Address),
		ProfilePicture: profile.ProfilePicture,
		Preferences:    ToPreferencesResponse(profile.Preferences),
	}
}

// ToAddressResponse converts a domain Address to an AddressResponse DTO
func ToAddressResponse(address domuser.Address) AddressResponse {
	return AddressResponse{
		Street:     address.Street,
		City:       address.City,
		State:      address.State,
		PostalCode: address.PostalCode,
		Country:    address.Country,
	}
}

// ToPreferencesResponse converts domain Preferences to a PreferencesResponse DTO
func ToPreferencesResponse(prefs domuser.Preferences) PreferencesResponse {
	return PreferencesResponse{
		Language:          prefs.Language,
		Timezone:          prefs.Timezone,
		NotificationEmail: prefs.NotificationEmail,
		NotificationSMS:   prefs.NotificationSMS,
	}
}

// ToUsersListResponse converts a list of users and pagination info to UsersListResponse
func ToUsersListResponse(users []*domuser.User, total, page, pageSize int) UsersListResponse {
	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = ToUserResponse(user)
	}

	totalPages := (total + pageSize - 1) / pageSize
	if totalPages == 0 {
		totalPages = 1
	}

	return UsersListResponse{
		Users:      userResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

// ToUpdateMap converts UpdateUserRequest to a map for service layer
func (req *UpdateUserRequest) ToUpdateMap() map[string]interface{} {
	updates := make(map[string]interface{})

	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.FirstName != nil {
		updates["firstName"] = *req.FirstName
	}
	if req.LastName != nil {
		updates["lastName"] = *req.LastName
	}
	if req.Role != nil {
		updates["role"] = string(*req.Role)
	}
	if req.IsActive != nil {
		updates["isActive"] = *req.IsActive
	}
	if req.EmailVerified != nil {
		updates["emailVerified"] = *req.EmailVerified
	}

	return updates
}

// ToProfile converts UpdateProfileRequest to a domain Profile
func (req *UpdateProfileRequest) ToProfile(existingProfile domuser.Profile) domuser.Profile {
	profile := existingProfile // Start with existing profile

	if req.Phone != nil {
		profile.Phone = *req.Phone
	}
	if req.DateOfBirth != nil {
		profile.DateOfBirth = req.DateOfBirth
	}
	if req.Address != nil {
		profile.Address = domuser.Address{
			Street:     req.Address.Street,
			City:       req.Address.City,
			State:      req.Address.State,
			PostalCode: req.Address.PostalCode,
			Country:    req.Address.Country,
		}
	}
	if req.Preferences != nil {
		if req.Preferences.Language != nil {
			profile.Preferences.Language = *req.Preferences.Language
		}
		if req.Preferences.Timezone != nil {
			profile.Preferences.Timezone = *req.Preferences.Timezone
		}
		if req.Preferences.NotificationEmail != nil {
			profile.Preferences.NotificationEmail = *req.Preferences.NotificationEmail
		}
		if req.Preferences.NotificationSMS != nil {
			profile.Preferences.NotificationSMS = *req.Preferences.NotificationSMS
		}
	}

	return profile
}
