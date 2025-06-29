package service

import (
	"context"
	"todo-list-api/internal/models"
)

// TodoService defines the interface for todo business logic operations
type TodoService interface {
	CreateTodo(ctx context.Context, req *models.CreateTodoRequest) (*models.Todo, error)
	GetTodos(ctx context.Context) ([]models.Todo, error)
	GetTodoByID(ctx context.Context, id uint) (*models.Todo, error)
	UpdateTodo(ctx context.Context, id uint, req *models.UpdateTodoRequest) (*models.Todo, error)
	DeleteTodo(ctx context.Context, id uint) error
	GetTodosByUserID(ctx context.Context, userID uint) ([]models.Todo, error)
}
