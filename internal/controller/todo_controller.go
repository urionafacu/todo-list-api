package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"todo-list-api/internal/models"
	"todo-list-api/internal/service"

	"github.com/go-chi/chi/v5"
)

type TodoController struct {
	todoService service.TodoService
}

// NewTodoController creates a new instance of TodoController
func NewTodoController(todoService service.TodoService) *TodoController {
	return &TodoController{
		todoService: todoService,
	}
}

func (c *TodoController) GetTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := c.todoService.GetTodos(r.Context())
	if err != nil {
		c.writeError(w, http.StatusInternalServerError, "Failed to get todos")
		return
	}

	c.writeJSON(w, http.StatusOK, todos)
}

func (c *TodoController) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.writeError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	todo, err := c.todoService.CreateTodo(r.Context(), &req)
	if err != nil {
		// Check if it's a validation error
		if isValidationError(err) {
			c.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		c.writeError(w, http.StatusInternalServerError, "Failed to create todo")
		return
	}

	c.writeJSON(w, http.StatusCreated, todo)
}

func (c *TodoController) GetTodoByID(w http.ResponseWriter, r *http.Request) {
	id, err := c.parseIDFromURL(r)
	if err != nil {
		c.writeError(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}

	todo, err := c.todoService.GetTodoByID(r.Context(), id)
	if err != nil {
		if err.Error() == "invalid todo ID" {
			c.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		c.writeError(w, http.StatusInternalServerError, "Failed to get todo")
		return
	}

	if todo == nil {
		c.writeError(w, http.StatusNotFound, "Todo not found")
		return
	}

	c.writeJSON(w, http.StatusOK, todo)
}

func (c *TodoController) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id, err := c.parseIDFromURL(r)
	if err != nil {
		c.writeError(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}

	var req models.UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.writeError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	todo, err := c.todoService.UpdateTodo(r.Context(), id, &req)
	if err != nil {
		// Check if it's a validation error
		if isValidationError(err) {
			c.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		// Check if it's a not found error
		if err.Error() == "todo not found" || err.Error() == "invalid todo ID" {
			c.writeError(w, http.StatusNotFound, err.Error())
			return
		}
		c.writeError(w, http.StatusInternalServerError, "Failed to update todo")
		return
	}

	c.writeJSON(w, http.StatusOK, todo)
}

func (c *TodoController) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id, err := c.parseIDFromURL(r)
	if err != nil {
		c.writeError(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}

	err = c.todoService.DeleteTodo(r.Context(), id)
	if err != nil {
		if err.Error() == "todo not found" || err.Error() == "invalid todo ID" {
			c.writeError(w, http.StatusNotFound, err.Error())
			return
		}
		c.writeError(w, http.StatusInternalServerError, "Failed to delete todo")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Helper methods

func (c *TodoController) parseIDFromURL(r *http.Request) (uint, error) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func (c *TodoController) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (c *TodoController) writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	errorResponse := map[string]string{"error": message}
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
	}
}

// isValidationError checks if the error is a business logic validation error
func isValidationError(err error) bool {
	validationErrors := []string{
		"title is required",
		"title cannot exceed 200 characters",
		"description cannot exceed 1000 characters",
		"request cannot be nil",
	}

	for _, validationError := range validationErrors {
		if err.Error() == validationError {
			return true
		}
	}
	return false
}
