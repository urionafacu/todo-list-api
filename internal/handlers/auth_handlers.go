package handlers

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"todo-list-api/internal/database"
	"todo-list-api/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	userService *database.UserService
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{
		userService: &database.UserService{},
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
	// Implement login logic here
}
