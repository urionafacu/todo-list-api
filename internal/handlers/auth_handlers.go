package handlers

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"os"
	"todo-list-api/internal/database"
	"todo-list-api/internal/models"
	"todo-list-api/internal/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	userService *database.UserService
	jwtSecret   string
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT_SECRET environment variable must be set")
	}

	return &AuthHandler{
		userService: database.NewUserService(db),
		jwtSecret:   os.Getenv("JWT_SECRET"),
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Basic validation
	if req.Email == "" {
		WriteError(w, http.StatusBadRequest, "Email is required")
		return
	}

	// Validate email format
	if _, err := mail.ParseAddress(req.Email); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid email format")
		return
	}

	if req.Password == "" {
		WriteError(w, http.StatusBadRequest, "Password is required")
		return
	}

	if len([]rune(req.Password)) < 8 {
		WriteError(w, http.StatusBadRequest, "Password must be at least 8 characters long")
		return
	}

	if (len([]rune(req.Password))) > 128 {
		WriteError(w, http.StatusBadRequest, "Password must be less than 128 characters")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Error processing password")
		return
	}

	user, err := h.userService.CreateUser(r.Context(), req.Email, req.FirstName, req.LastName, string(hashedPassword))
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Error creating user")
		return
	}
	WriteJson(w, http.StatusCreated, user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Basic validations
	if req.Email == "" || req.Password == "" {
		WriteError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	user, err := h.userService.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			WriteError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Error retrieving user")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		WriteError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	jwtUtil := &utils.JWT{Secret: h.jwtSecret}

	token, err := jwtUtil.CreateToken(user)

	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	WriteJson(w, http.StatusOK, token)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Basic validation
	if req.RefreshToken == "" {
		WriteError(w, http.StatusBadRequest, "Refresh token is required")
		return
	}

	jwtUtil := &utils.JWT{Secret: h.jwtSecret}

	// Refresh the access token using the refresh token
	tokenPair, err := jwtUtil.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "Invalid or expired refresh token")
		return
	}

	// Convert to the existing Token model for backward compatibility
	token := &models.Token{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}

	WriteJson(w, http.StatusOK, token)
}
