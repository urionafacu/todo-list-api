package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"todo-list-api/internal/database"
	"todo-list-api/internal/models"

	"gorm.io/gorm"
)

type TodoHandlers struct {
	todoService *database.TodoService
}

func NewTodoHandlers(db *gorm.DB) *TodoHandlers {
	return &TodoHandlers{
		todoService: database.NewTodoService(db),
	}
}

func (h *TodoHandlers) GetTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.todoService.GetTodos(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to get todos")
		return
	}

	WriteJson(w, http.StatusOK, todos)
}

func (h *TodoHandlers) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Basic validation
	if req.Title == "" {
		WriteError(w, http.StatusBadRequest, "Title is required")
		return
	}

	todo, err := h.todoService.CreateTodo(r.Context(), req.Title, req.Description)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to create todo")
		return
	}

	WriteJson(w, http.StatusCreated, todo)
}

func (h *TodoHandlers) GetTodoByID(w http.ResponseWriter, r *http.Request) {
	id, err := h.extractIDFromPath(r.URL.Path)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}

	todo, err := h.todoService.GetTodoByID(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to get todo")
		return
	}

	if todo == nil {
		WriteError(w, http.StatusNotFound, "Todo not found")
		return
	}

	WriteJson(w, http.StatusOK, todo)
}

func (h *TodoHandlers) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id, err := h.extractIDFromPath(r.URL.Path)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}

	var req models.UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Basic validation
	if req.Title == "" {
		WriteError(w, http.StatusBadRequest, "Title is required")
		return
	}

	todo, err := h.todoService.UpdateTodo(r.Context(), id, req)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to update todo")
		return
	}

	if todo == nil {
		WriteError(w, http.StatusNotFound, "Todo not found")
		return
	}

	WriteJson(w, http.StatusOK, todo)
}

func (h *TodoHandlers) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id, err := h.extractIDFromPath(r.URL.Path)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}

	err = h.todoService.DeleteTodo(r.Context(), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			WriteError(w, http.StatusNotFound, "Todo not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to delete todo")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// extractIDFromPath extracts the ID from URL paths like /api/todos/123
func (h *TodoHandlers) extractIDFromPath(path string) (uint, error) {
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		return 0, gorm.ErrInvalidData
	}

	idStr := parts[3] // /api/todos/{id}
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}
