package repository

import (
	"context"
	"todo-list-api/internal/models"
)

// TodoRepository defines the interface for todo data access operations
type TodoRepository interface {
	Create(ctx context.Context, todo *models.Todo) (*models.Todo, error)
	GetAll(ctx context.Context) ([]models.Todo, error)
	GetByID(ctx context.Context, id uint) (*models.Todo, error)
	Update(ctx context.Context, id uint, todo *models.Todo) (*models.Todo, error)
	Delete(ctx context.Context, id uint) error
	GetByUserID(ctx context.Context, userID uint) ([]models.Todo, error)
}
