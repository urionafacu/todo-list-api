package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"todo-list-api/internal/database"
	"todo-list-api/internal/models"

	"github.com/go-chi/chi/v5"
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
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}

	todo, err := h.todoService.GetTodoByID(r.Context(), uint(id))
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
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
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

	todo, err := h.todoService.UpdateTodo(r.Context(), uint(id), req)
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
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}

	err = h.todoService.DeleteTodo(r.Context(), uint(id))
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
