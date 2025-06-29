package service

import (
	"context"
	"errors"
	"net/mail"
	"todo-list-api/internal/models"
	"todo-list-api/internal/repository"
	"todo-list-api/internal/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type authServiceImpl struct {
	authRepo repository.AuthRepository
	jwtUtil  *utils.JWT
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(authRepo repository.AuthRepository, jwtUtil *utils.JWT) AuthService {
	return &authServiceImpl{
		authRepo: authRepo,
		jwtUtil:  jwtUtil,
	}
}

func (s *authServiceImpl) Register(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	// Validate request
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	if req.Email == "" {
		return nil, errors.New("email is required")
	}

	// Validate email format
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return nil, errors.New("invalid email format")
	}

	if req.Password == "" {
		return nil, errors.New("password is required")
	}

	if len([]rune(req.Password)) < 8 {
		return nil, errors.New("password must be at least 8 characters long")
	}

	if len([]rune(req.Password)) > 128 {
		return nil, errors.New("password must be less than 128 characters")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("error processing password")
	}

	// Create user
	user := &models.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  string(hashedPassword),
	}

	createdUser, err := s.authRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, errors.New("error creating user")
	}

	return createdUser, nil
}

func (s *authServiceImpl) Login(ctx context.Context, req *models.LoginUserRequest) (*models.Token, error) {
	// Validate request
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	if req.Email == "" || req.Password == "" {
		return nil, errors.New("email and password are required")
	}

	// Get user by email
	user, err := s.authRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, errors.New("error retrieving user")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate token
	token, err := s.jwtUtil.CreateToken(user)
	if err != nil {
		return nil, errors.New("error generating token")
	}

	return token, nil
}

func (s *authServiceImpl) RefreshToken(ctx context.Context, req *models.RefreshTokenRequest) (*models.Token, error) {
	// Validate request
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	if req.RefreshToken == "" {
		return nil, errors.New("refresh token is required")
	}

	// Refresh the access token using the refresh token
	tokenPair, err := s.jwtUtil.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	// Convert to the existing Token model for backward compatibility
	token := &models.Token{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}

	return token, nil
}

func (s *authServiceImpl) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}

	user, err := s.authRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (s *authServiceImpl) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	if id == 0 {
		return nil, errors.New("invalid user ID")
	}

	user, err := s.authRepo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}
