package service

import (
	"context"
	"errors"
	"testing"
	"time"
	"todo-list-api/internal/models"
	"todo-list-api/internal/service/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TodoServiceTestSuite struct {
	suite.Suite
	mockRepo *mocks.MockTodoRepository
	service  TodoService
	ctx      context.Context
}

func (suite *TodoServiceTestSuite) SetupTest() {
	suite.mockRepo = new(mocks.MockTodoRepository)
	suite.service = NewTodoService(suite.mockRepo)
	suite.ctx = context.Background()
}

func (suite *TodoServiceTestSuite) TestCreateTodo_Success() {
	// Arrange
	req := &models.CreateTodoRequest{
		Title:       "Test Todo",
		Description: "Test Description",
		Priority:    "high",
		Category:    "work",
	}

	expectedTodo := &models.Todo{
		ID:          1,
		Title:       "Test Todo",
		Description: "Test Description",
		Priority:    "high",
		Category:    "work",
		Completed:   false,
	}

	suite.mockRepo.On("Create", suite.ctx, mock.MatchedBy(func(todo *models.Todo) bool {
		return todo.Title == "Test Todo" &&
			todo.Description == "Test Description" &&
			todo.Priority == "high" &&
			todo.Category == "work" &&
			todo.Completed == false
	})).Return(expectedTodo, nil)

	// Act
	result, err := suite.service.CreateTodo(suite.ctx, req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), expectedTodo.Title, result.Title)
	assert.Equal(suite.T(), expectedTodo.Priority, result.Priority)
	assert.Equal(suite.T(), expectedTodo.Category, result.Category)
	assert.False(suite.T(), result.Completed)
}

// TestCreateTodo_RepositoryError tests repository error handling
func (suite *TodoServiceTestSuite) TestCreateTodo_RepositoryError() {
	// Arrange
	req := &models.CreateTodoRequest{
		Title:    "Test Todo",
		Priority: "high",
	}

	suite.mockRepo.On("Create", suite.ctx, mock.AnythingOfType("*models.Todo")).
		Return(nil, errors.New("database error"))

	// Act
	result, err := suite.service.CreateTodo(suite.ctx, req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "database error")
}

// TestGetTodoByID_Success tests successful todo retrieval
func (suite *TodoServiceTestSuite) TestGetTodoByID_Success() {
	// Arrange
	todoID := uint(1)
	expectedTodo := &models.Todo{
		ID:    todoID,
		Title: "Test Todo",
	}

	suite.mockRepo.On("GetByID", suite.ctx, todoID).Return(expectedTodo, nil)

	// Act
	result, err := suite.service.GetTodoByID(suite.ctx, todoID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), expectedTodo.ID, result.ID)
	assert.Equal(suite.T(), expectedTodo.Title, result.Title)
}

// TestGetTodoByID_InvalidID tests invalid ID handling
func (suite *TodoServiceTestSuite) TestGetTodoByID_InvalidID() {
	// Act
	result, err := suite.service.GetTodoByID(suite.ctx, 0)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "invalid todo ID")
}

// TestGetTodoByID_NotFound tests todo not found scenario
func (suite *TodoServiceTestSuite) TestGetTodoByID_NotFound() {
	// Arrange
	todoID := uint(999)
	suite.mockRepo.On("GetByID", suite.ctx, todoID).Return(nil, errors.New("todo not found"))

	// Act
	result, err := suite.service.GetTodoByID(suite.ctx, todoID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

// TestUpdateTodo_Success tests successful todo update
func (suite *TodoServiceTestSuite) TestUpdateTodo_Success() {
	// Arrange
	todoID := uint(1)
	existingTodo := &models.Todo{
		ID:        todoID,
		Title:     "Old Title",
		CreatedAt: time.Now().Add(-time.Hour),
	}

	req := &models.UpdateTodoRequest{
		Title:       "  Updated Title  ",
		Description: "Updated Description",
		Completed:   true,
		Priority:    "medium",
	}

	updatedTodo := &models.Todo{
		ID:          todoID,
		Title:       "Updated Title",
		Description: "Updated Description",
		Completed:   true,
		Priority:    "medium",
	}

	suite.mockRepo.On("GetByID", suite.ctx, todoID).Return(existingTodo, nil)
	suite.mockRepo.On("Update", suite.ctx, todoID, mock.MatchedBy(func(todo *models.Todo) bool {
		return todo.Title == "Updated Title" &&
			todo.Description == "Updated Description" &&
			todo.Completed == true &&
			todo.Priority == "medium" &&
			todo.CreatedAt.Equal(existingTodo.CreatedAt)
	})).Return(updatedTodo, nil)

	// Act
	result, err := suite.service.UpdateTodo(suite.ctx, todoID, req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), updatedTodo.Title, result.Title)
	assert.Equal(suite.T(), updatedTodo.Description, result.Description)
	assert.True(suite.T(), result.Completed)
}

// TestUpdateTodo_InvalidID tests invalid ID handling
func (suite *TodoServiceTestSuite) TestUpdateTodo_InvalidID() {
	// Arrange
	req := &models.UpdateTodoRequest{Title: "Updated Title"}

	// Act
	result, err := suite.service.UpdateTodo(suite.ctx, 0, req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "invalid todo ID")
}

// TestUpdateTodo_TodoNotFound tests todo not found scenario
func (suite *TodoServiceTestSuite) TestUpdateTodo_TodoNotFound() {
	// Arrange
	todoID := uint(999)
	req := &models.UpdateTodoRequest{Title: "Updated Title"}

	suite.mockRepo.On("GetByID", suite.ctx, todoID).Return(nil, nil)

	// Act
	result, err := suite.service.UpdateTodo(suite.ctx, todoID, req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "todo not found")
}

// TestDeleteTodo_Success tests successful todo deletion
func (suite *TodoServiceTestSuite) TestDeleteTodo_Success() {
	// Arrange
	todoID := uint(1)
	existingTodo := &models.Todo{ID: todoID, Title: "Test Todo"}

	suite.mockRepo.On("GetByID", suite.ctx, todoID).Return(existingTodo, nil)
	suite.mockRepo.On("Delete", suite.ctx, todoID).Return(nil)

	// Act
	err := suite.service.DeleteTodo(suite.ctx, todoID)

	// Assert
	assert.NoError(suite.T(), err)
}

// TestDeleteTodo_InvalidID tests invalid ID handling
func (suite *TodoServiceTestSuite) TestDeleteTodo_InvalidID() {
	// Act
	err := suite.service.DeleteTodo(suite.ctx, 0)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "invalid todo ID")
}

// TestGetTodosByUserID_Success tests successful user todos retrieval
func (suite *TodoServiceTestSuite) TestGetTodosByUserID_Success() {
	// Arrange
	userID := uint(1)
	expectedTodos := []models.Todo{
		{ID: 1, Title: "Todo 1"},
		{ID: 2, Title: "Todo 2"},
	}

	suite.mockRepo.On("GetByUserID", suite.ctx, userID).Return(expectedTodos, nil)

	// Act
	result, err := suite.service.GetTodosByUserID(suite.ctx, userID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)
	assert.Equal(suite.T(), expectedTodos[0].Title, result[0].Title)
}

// TestGetTodosByUserID_InvalidID tests invalid user ID handling
func (suite *TodoServiceTestSuite) TestGetTodosByUserID_InvalidID() {
	// Act
	result, err := suite.service.GetTodosByUserID(suite.ctx, 0)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "invalid user ID")
}

// TestTodoServiceSuite runs the test suite
func TestTodoServiceSuite(t *testing.T) {
	suite.Run(t, new(TodoServiceTestSuite))
}
