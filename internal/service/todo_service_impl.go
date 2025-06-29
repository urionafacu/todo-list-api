package service

import (
	"context"
	"errors"
	"strings"
	"time"
	"todo-list-api/internal/models"
	"todo-list-api/internal/repository"
)

type todoServiceImpl struct {
	todoRepo repository.TodoRepository
}

// NewTodoService creates a new instance of TodoService
func NewTodoService(todoRepo repository.TodoRepository) TodoService {
	return &todoServiceImpl{
		todoRepo: todoRepo,
	}
}

func (s *todoServiceImpl) CreateTodo(ctx context.Context, req *models.CreateTodoRequest) (*models.Todo, error) {
	// Business logic validations
	if err := s.validateCreateTodoRequest(req); err != nil {
		return nil, err
	}

	// Create todo entity
	todo := &models.Todo{
		Title:       strings.TrimSpace(req.Title),
		Description: strings.TrimSpace(req.Description),
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Delegate to repository
	return s.todoRepo.Create(ctx, todo)
}

func (s *todoServiceImpl) GetTodos(ctx context.Context) ([]models.Todo, error) {
	return s.todoRepo.GetAll(ctx)
}

func (s *todoServiceImpl) GetTodoByID(ctx context.Context, id uint) (*models.Todo, error) {
	if id == 0 {
		return nil, errors.New("invalid todo ID")
	}

	return s.todoRepo.GetByID(ctx, id)
}

func (s *todoServiceImpl) UpdateTodo(ctx context.Context, id uint, req *models.UpdateTodoRequest) (*models.Todo, error) {
	if id == 0 {
		return nil, errors.New("invalid todo ID")
	}

	// Business logic validations
	if err := s.validateUpdateTodoRequest(req); err != nil {
		return nil, err
	}

	// Check if todo exists
	existingTodo, err := s.todoRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existingTodo == nil {
		return nil, errors.New("todo not found")
	}

	// Create updated todo entity
	updatedTodo := &models.Todo{
		ID:          id,
		Title:       strings.TrimSpace(req.Title),
		Description: strings.TrimSpace(req.Description),
		Completed:   req.Completed,
		UpdatedAt:   time.Now(),
		// Preserve original creation time
		CreatedAt: existingTodo.CreatedAt,
	}

	return s.todoRepo.Update(ctx, id, updatedTodo)
}

func (s *todoServiceImpl) DeleteTodo(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid todo ID")
	}

	// Check if todo exists before deleting
	existingTodo, err := s.todoRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existingTodo == nil {
		return errors.New("todo not found")
	}

	return s.todoRepo.Delete(ctx, id)
}

func (s *todoServiceImpl) GetTodosByUserID(ctx context.Context, userID uint) ([]models.Todo, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}

	return s.todoRepo.GetByUserID(ctx, userID)
}

// Private validation methods

func (s *todoServiceImpl) validateCreateTodoRequest(req *models.CreateTodoRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return errors.New("title is required")
	}

	if len(title) > 200 {
		return errors.New("title cannot exceed 200 characters")
	}

	if len(req.Description) > 1000 {
		return errors.New("description cannot exceed 1000 characters")
	}

	return nil
}

func (s *todoServiceImpl) validateUpdateTodoRequest(req *models.UpdateTodoRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return errors.New("title is required")
	}

	if len(title) > 200 {
		return errors.New("title cannot exceed 200 characters")
	}

	if len(req.Description) > 1000 {
		return errors.New("description cannot exceed 1000 characters")
	}

	return nil
}
