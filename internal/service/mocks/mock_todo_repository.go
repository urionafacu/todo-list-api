package mocks

import (
	"context"
	"todo-list-api/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockTodoRepository struct {
	mock.Mock
}

func (m *MockTodoRepository) Create(ctx context.Context, todo *models.Todo) (*models.Todo, error) {
	args := m.Called(ctx, todo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Todo), args.Error(1)
}

func (m *MockTodoRepository) GetAll(ctx context.Context) ([]models.Todo, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Todo), args.Error(1)
}

func (m *MockTodoRepository) GetByID(ctx context.Context, id uint) (*models.Todo, error) {
	// To implement
	return nil, nil
}

func (m *MockTodoRepository) Update(ctx context.Context, id uint, todo *models.Todo) (*models.Todo, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Todo), args.Error(1)
}

func (m *MockTodoRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTodoRepository) GetByUserID(ctx context.Context, userID uint) ([]models.Todo, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Todo), args.Error(1)
}
