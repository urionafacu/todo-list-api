package database

import (
	"context"
	"errors"
	"todo-list-api/internal/models"

	"gorm.io/gorm"
)

type TodoService struct {
	db *gorm.DB
}

func NewTodoService(db *gorm.DB) *TodoService {
	return &TodoService{db}
}

func (s *TodoService) CreateTodo(ctx context.Context, title, description string) (*models.Todo, error) {
	todo := models.Todo{
		Title:       title,
		Description: description,
		Completed:   false,
	}

	result := s.db.WithContext(ctx).Create(&todo)
	if result.Error != nil {
		return nil, result.Error
	}

	return &todo, nil
}

func (s *TodoService) GetTodos(ctx context.Context) ([]models.Todo, error) {
	var todos []models.Todo

	result := s.db.WithContext(ctx).Order("created_at DESC").Find(&todos)
	if result.Error != nil {
		return nil, result.Error
	}

	return todos, nil
}

func (s *TodoService) GetTodoByID(ctx context.Context, id uint) (*models.Todo, error) {
	var todo models.Todo

	result := s.db.WithContext(ctx).First(&todo, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil for not found, similar to sql.ErrNoRows
		}
		return nil, result.Error
	}

	return &todo, nil
}

func (s *TodoService) UpdateTodo(ctx context.Context, id uint, updates models.UpdateTodoRequest) (*models.Todo, error) {
	var todo models.Todo

	// First check if the todo exists
	result := s.db.WithContext(ctx).First(&todo, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	// Update the todo
	updateData := models.Todo{
		Title:       updates.Title,
		Description: updates.Description,
		Completed:   updates.Completed,
	}

	result = s.db.WithContext(ctx).Model(&todo).Updates(updateData)
	if result.Error != nil {
		return nil, result.Error
	}

	// Return the updated todo
	return &todo, nil
}

func (s *TodoService) DeleteTodo(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).Delete(&models.Todo{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
