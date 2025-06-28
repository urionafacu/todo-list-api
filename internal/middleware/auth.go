package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"
	"todo-list-api/internal/handlers"
	"todo-list-api/internal/utils"
)

type contextKey string

const (
	UserIDKey    contextKey = "user_id"
	UserEmailKey contextKey = "user_email"
	ClaimsKey    contextKey = "claims"
)

// AuthMiddleware validates JWT tokens and adds user information to context
func AuthMiddleware(next http.Handler) http.Handler {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT_SECRET environment variable must be set")
	}

	jwtUtil := &utils.JWT{Secret: jwtSecret}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			handlers.WriteError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		// Check if header starts with "Bearer "
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			handlers.WriteError(w, http.StatusUnauthorized, "Invalid authorization header format")
			return
		}

		tokenString := parts[1]
		if tokenString == "" {
			handlers.WriteError(w, http.StatusUnauthorized, "Token is required")
			return
		}

		// Validate the access token
		claims, err := jwtUtil.ValidateAccessToken(tokenString)
		if err != nil {
			handlers.WriteError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		// Add user information to request context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserEmailKey, claims.UserEmail)
		ctx = context.WithValue(ctx, ClaimsKey, claims)

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuthMiddleware is similar to AuthMiddleware but doesn't require authentication
// If a valid token is provided, user info is added to context
// If no token or invalid token, request continues without user info
func OptionalAuthMiddleware(next http.Handler) http.Handler {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT_SECRET environment variable must be set")
	}

	jwtUtil := &utils.JWT{Secret: jwtSecret}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" && parts[1] != "" {
				// Try to validate the token
				if claims, err := jwtUtil.ValidateAccessToken(parts[1]); err == nil {
					// Add user information to request context
					ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
					ctx = context.WithValue(ctx, UserEmailKey, claims.UserEmail)
					ctx = context.WithValue(ctx, ClaimsKey, claims)
					r = r.WithContext(ctx)
				}
			}
		}

		// Continue regardless of authentication status
		next.ServeHTTP(w, r)
	})
}

// GetUserIDFromContext extracts user ID from request context
func GetUserIDFromContext(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value(UserIDKey).(uint)
	return userID, ok
}

// GetUserEmailFromContext extracts user email from request context
func GetUserEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(UserEmailKey).(string)
	return email, ok
}

// GetClaimsFromContext extracts JWT claims from request context
func GetClaimsFromContext(ctx context.Context) (*utils.Claims, bool) {
	claims, ok := ctx.Value(ClaimsKey).(*utils.Claims)
	return claims, ok
}

// RequireUserID is a helper middleware that ensures a user ID is present in context
func RequireUserID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := GetUserIDFromContext(r.Context()); !ok {
			handlers.WriteError(w, http.StatusUnauthorized, "User authentication required")
			return
		}
		next.ServeHTTP(w, r)
	})
}
