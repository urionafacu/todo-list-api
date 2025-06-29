package repository

import (
	"context"
	"todo-list-api/internal/models"
)

// AuthRepository defines the interface for auth data access operations
type AuthRepository interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id uint) (*models.User, error)
}
