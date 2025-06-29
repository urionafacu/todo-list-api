package controller

import (
	"encoding/json"
	"net/http"
	"todo-list-api/internal/models"
	"todo-list-api/internal/service"
)

type AuthController struct {
	authService service.AuthService
}

// NewAuthController creates a new instance of AuthController
func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.writeError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	user, err := c.authService.Register(r.Context(), &req)
	if err != nil {
		// Check if it's a validation error
		if isAuthValidationError(err) {
			c.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		c.writeError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	c.writeJSON(w, http.StatusCreated, user)
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.writeError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	token, err := c.authService.Login(r.Context(), &req)
	if err != nil {
		// Check if it's an authentication error
		if err.Error() == "invalid email or password" {
			c.writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		// Check if it's a validation error
		if isAuthValidationError(err) {
			c.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		c.writeError(w, http.StatusInternalServerError, "Error during login")
		return
	}

	c.writeJSON(w, http.StatusOK, token)
}

func (c *AuthController) Refresh(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.writeError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	token, err := c.authService.RefreshToken(r.Context(), &req)
	if err != nil {
		// Check if it's a token validation error
		if err.Error() == "invalid or expired refresh token" {
			c.writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		// Check if it's a validation error
		if isAuthValidationError(err) {
			c.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		c.writeError(w, http.StatusInternalServerError, "Error refreshing token")
		return
	}

	c.writeJSON(w, http.StatusOK, token)
}

// Helper methods

func (c *AuthController) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (c *AuthController) writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	errorResponse := map[string]string{"error": message}
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
	}
}

// isAuthValidationError checks if the error is a business logic validation error
func isAuthValidationError(err error) bool {
	validationErrors := []string{
		"request cannot be nil",
		"email is required",
		"invalid email format",
		"password is required",
		"password must be at least 8 characters long",
		"password must be less than 128 characters",
		"email and password are required",
		"refresh token is required",
		"invalid user ID",
	}

	for _, validationError := range validationErrors {
		if err.Error() == validationError {
			return true
		}
	}
	return false
}
