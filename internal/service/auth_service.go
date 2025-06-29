package service

import (
	"context"
	"todo-list-api/internal/models"
)

// AuthService defines the interface for authentication business logic operations
type AuthService interface {
	Register(ctx context.Context, req *models.CreateUserRequest) (*models.User, error)
	Login(ctx context.Context, req *models.LoginUserRequest) (*models.Token, error)
	RefreshToken(ctx context.Context, req *models.RefreshTokenRequest) (*models.Token, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id uint) (*models.User, error)
}
