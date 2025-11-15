package dto

import (
	"time"

	domuser "github.com/fintrack/user-service/internal/core/domain/entities/user"
)

// CreateUserRequest represents the request to create a new user
// @Description Request payload for creating a new user
type CreateUserRequest struct {
	Email     string `json:"email" binding:"required,email" example:"user@example.com" description:"User email address"`
	Password  string `json:"password" binding:"required,min=8" example:"secretpassword123" description:"User password (minimum 8 characters)"`
	FirstName string `json:"firstName" binding:"required" example:"John" description:"User first name"`
	LastName  string `json:"lastName" binding:"required" example:"Doe" description:"User last name"`
	Role      string `json:"role" binding:"required" example:"user" description:"User role (admin, user, operator, treasurer)" enums:"user,operator,admin,treasurer"`
}

// UpdateUserRequest represents the request to update user information
// @Description Request payload for updating user information
type UpdateUserRequest struct {
	Email         *string `json:"email,omitempty" binding:"omitempty,email" example:"newemail@example.com" description:"New email address"`
	FirstName     *string `json:"firstName,omitempty" example:"John" description:"Updated first name"`
	LastName      *string `json:"lastName,omitempty" example:"Doe" description:"Updated last name"`
	Role          *string `json:"role,omitempty" example:"admin" description:"Updated user role" enums:"user,operator,admin,treasurer"`
	IsActive      *bool   `json:"isActive,omitempty" example:"true" description:"User active status"`
	EmailVerified *bool   `json:"emailVerified,omitempty" example:"true" description:"Email verification status"`
}

// UpdateProfileRequest represents the request to update user profile
// @Description Request payload for updating user profile information
type UpdateProfileRequest struct {
	Phone       *string             `json:"phone,omitempty" example:"+1234567890" description:"Phone number"`
	DateOfBirth *time.Time          `json:"dateOfBirth,omitempty" example:"1990-01-01T00:00:00Z" description:"Date of birth"`
	Address     *AddressRequest     `json:"address,omitempty" description:"Address information"`
	Preferences *PreferencesRequest `json:"preferences,omitempty" description:"User preferences"`
}

// AddressRequest represents address information in requests
// @Description Address information for user profile
type AddressRequest struct {
	Street     string `json:"street,omitempty" example:"123 Main St" description:"Street address"`
	City       string `json:"city,omitempty" example:"New York" description:"City"`
	State      string `json:"state,omitempty" example:"NY" description:"State or province"`
	PostalCode string `json:"postalCode,omitempty" example:"10001" description:"Postal code"`
	Country    string `json:"country,omitempty" example:"USA" description:"Country"`
}

// PreferencesRequest represents user preferences in requests
// @Description User preferences for notifications and display
type PreferencesRequest struct {
	Language          *string `json:"language,omitempty" example:"en" description:"Preferred language"`
	Timezone          *string `json:"timezone,omitempty" example:"UTC" description:"User timezone"`
	NotificationEmail *bool   `json:"notificationEmail,omitempty" example:"true" description:"Enable email notifications"`
	NotificationSMS   *bool   `json:"notificationSMS,omitempty" example:"false" description:"Enable SMS notifications"`
}

// ChangePasswordRequest represents the request to change password
// @Description Request payload for changing user password
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required" example:"oldpassword123" description:"Current password"`
	NewPassword string `json:"newPassword" binding:"required,min=8" example:"newpassword123" description:"New password (minimum 8 characters)"`
}

// ChangeRoleRequest represents the request to change user role
// @Description Request payload for changing user role
type ChangeRoleRequest struct {
	Role string `json:"role" binding:"required" example:"admin" description:"New user role" enums:"user,operator,admin,treasurer"`
}

// ToggleStatusRequest represents the request to toggle user status
// @Description Request payload for toggling user active status
type ToggleStatusRequest struct {
	IsActive bool `json:"isActive" example:"true" description:"User active status"`
}

// UserResponse represents user information in responses
// @Description Complete user information including profile data
type UserResponse struct {
	ID            string          `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" description:"User unique identifier"`
	Email         string          `json:"email" example:"user@example.com" description:"User email address"`
	FirstName     string          `json:"firstName" example:"John" description:"User first name"`
	LastName      string          `json:"lastName" example:"Doe" description:"User last name"`
	FullName      string          `json:"fullName" example:"John Doe" description:"User full name"`
	Role          string          `json:"role" example:"user" description:"User role" enums:"user,operator,admin,treasurer"`
	IsActive      bool            `json:"isActive" example:"true" description:"User active status"`
	EmailVerified bool            `json:"emailVerified" example:"true" description:"Email verification status"`
	Profile       ProfileResponse `json:"profile" description:"User profile information"`
	CreatedAt     time.Time       `json:"createdAt" example:"2023-01-01T00:00:00Z" description:"Account creation timestamp"`
	UpdatedAt     time.Time       `json:"updatedAt" example:"2023-01-01T00:00:00Z" description:"Last update timestamp"`
	LastLoginAt   *time.Time      `json:"lastLoginAt,omitempty" example:"2023-01-01T00:00:00Z" description:"Last login timestamp"`
}

// ProfileResponse represents profile information in responses
// @Description User profile information including address and preferences
type ProfileResponse struct {
	Phone          string              `json:"phone,omitempty" example:"+1234567890" description:"Phone number"`
	DateOfBirth    *time.Time          `json:"dateOfBirth,omitempty" example:"1990-01-01T00:00:00Z" description:"Date of birth"`
	Address        AddressResponse     `json:"address" description:"Address information"`
	ProfilePicture string              `json:"profilePicture,omitempty" example:"https://example.com/avatar.jpg" description:"Profile picture URL"`
	Preferences    PreferencesResponse `json:"preferences" description:"User preferences"`
}

// AddressResponse represents address information in responses
// @Description Address information for user profile
type AddressResponse struct {
	Street     string `json:"street,omitempty" example:"123 Main St" description:"Street address"`
	City       string `json:"city,omitempty" example:"New York" description:"City"`
	State      string `json:"state,omitempty" example:"NY" description:"State or province"`
	PostalCode string `json:"postalCode,omitempty" example:"10001" description:"Postal code"`
	Country    string `json:"country,omitempty" example:"USA" description:"Country"`
}

// PreferencesResponse represents user preferences in responses
// @Description User preferences for notifications and display
type PreferencesResponse struct {
	Language          string `json:"language,omitempty" example:"en" description:"Preferred language"`
	Timezone          string `json:"timezone,omitempty" example:"UTC" description:"User timezone"`
	NotificationEmail bool   `json:"notificationEmail" example:"true" description:"Email notifications enabled"`
	NotificationSMS   bool   `json:"notificationSMS" example:"false" description:"SMS notifications enabled"`
}

// UsersListResponse represents a paginated list of users
// @Description Paginated list of users with metadata
type UsersListResponse struct {
	Users      []UserResponse `json:"users" description:"List of users"`
	Total      int            `json:"total" example:"100" description:"Total number of users"`
	Page       int            `json:"page" example:"1" description:"Current page number"`
	PageSize   int            `json:"pageSize" example:"20" description:"Number of items per page"`
	TotalPages int            `json:"totalPages" example:"5" description:"Total number of pages"`
}

// UserStatsResponse represents user statistics
// @Description User statistics and counts by different categories
type UserStatsResponse struct {
	TotalUsers    int `json:"totalUsers" example:"100" description:"Total number of users"`
	ActiveUsers   int `json:"activeUsers" example:"80" description:"Number of active users"`
	InactiveUsers int `json:"inactiveUsers" example:"20" description:"Number of inactive users"`
	AdminUsers    int `json:"adminUsers" example:"5" description:"Number of admin users"`
	RegularUsers  int `json:"regularUsers" example:"95" description:"Number of regular users"`
}

// ToUserResponse converts a domain User to a UserResponse DTO
func ToUserResponse(user *domuser.User) UserResponse {
	return UserResponse{
		ID:            user.ID,
		Email:         user.Email,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		FullName:      user.GetFullName(),
		Role:          string(user.Role),
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
		updates["role"] = *req.Role
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
