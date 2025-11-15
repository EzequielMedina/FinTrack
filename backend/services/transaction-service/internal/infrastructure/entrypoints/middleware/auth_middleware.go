package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
	// UserIDKey is the context key for storing the user ID
	UserIDKey contextKey = "userID"
)

// AuthMiddleware extracts user ID from JWT token and adds it to the request context
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get JWT secret from environment
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			jwtSecret = "fintrack-secret-key-2024" // Default secret (should match other services)
		}

		// Check if X-User-ID header is already present (from other service)
		userID := r.Header.Get("X-User-ID")
		
		// If not present, try to extract from JWT token
		if userID == "" {
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				// Extract token from "Bearer <token>" format
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					tokenStr := parts[1]
					
					// Parse and validate JWT token
					token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
						// Verify signing method
						if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
							return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
						}
						return []byte(jwtSecret), nil
					})

					if err == nil && token.Valid {
						// Extract user ID from claims
						if claims, ok := token.Claims.(jwt.MapClaims); ok {
							if sub, ok := claims["sub"].(string); ok && sub != "" {
								userID = sub
							}
						}
					}
				}
			}
		}

		// If still no user ID, check if this is a health check endpoint
		if userID == "" && r.URL.Path != "/health" {
			// Return 401 Unauthorized for requests without valid authentication
			http.Error(w, `{"error":"Unauthorized","message":"Missing or invalid authentication token"}`, http.StatusUnauthorized)
			return
		}

		// Add user ID to request header for downstream handlers
		if userID != "" {
			r.Header.Set("X-User-ID", userID)
			
			// Also add to context for easier access
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			r = r.WithContext(ctx)
		}

		// Continue to the next handler
		next.ServeHTTP(w, r)
	})
}

// GetUserIDFromContext extracts the user ID from the request context
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}
