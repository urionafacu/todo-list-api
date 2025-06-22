package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"todo-list-api/internal/database"
	"todo-list-api/internal/models"
)

type TodoHandlers struct {
	todoService *database.TodoService
}

func NewTodoHandlers(db *sql.DB) *TodoHandlers {
	return &TodoHandlers{
		todoService: database.NewTodoService(db),
	}
}

func (h *TodoHandlers) GetTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.todoService.GetTodos(r.Context())
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
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

	todo, err := h.todoService.CreateTodo(r.Context(), req.Title, req.Description)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Failed to create todo")
		return
	}

	WriteJson(w, http.StatusOK, todo)
}
