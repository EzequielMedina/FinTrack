package dto

// RegisterRequest represents the request payload for user registration
// @Description User registration request payload
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email" example:"user@example.com" description:"User email address"`
	Password  string `json:"password" binding:"required,min=8" example:"secretpassword123" description:"User password (minimum 8 characters)"`
	FirstName string `json:"firstName" binding:"required" example:"John" description:"User first name"`
	LastName  string `json:"lastName" binding:"required" example:"Doe" description:"User last name"`
}

// LoginRequest represents the request payload for user authentication
// @Description User login request payload
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com" description:"User email address"`
	Password string `json:"password" binding:"required" example:"secretpassword123" description:"User password"`
}

// AuthResponse represents the response for successful authentication
// @Description Authentication response with tokens and user information
type AuthResponse struct {
	AccessToken  string   `json:"accessToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." description:"JWT access token"`
	RefreshToken string   `json:"refreshToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." description:"JWT refresh token"`
	User         UserInfo `json:"user" description:"User information"`
}

// UserInfo represents basic user information
// @Description Basic user information
type UserInfo struct {
	Email     string `json:"email" example:"user@example.com" description:"User email address"`
	FirstName string `json:"firstName" example:"John" description:"User first name"`
	LastName  string `json:"lastName" example:"Doe" description:"User last name"`
}
