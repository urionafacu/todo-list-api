package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"todo-list-api/internal/models"
	"todo-list-api/internal/service"
	httputils "todo-list-api/internal/utils/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type TodoController struct {
	todoService service.TodoService
	validator   *validator.Validate
}

// NewTodoController creates a new instance of TodoController
func NewTodoController(todoService service.TodoService) *TodoController {
	return &TodoController{
		todoService: todoService,
		validator:   validator.New(),
	}
}

func (c *TodoController) GetTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := c.todoService.GetTodos(r.Context())
	if err != nil {
		httputils.WriteError(w, http.StatusInternalServerError, "Failed to get todos")
		return
	}

	httputils.WriteJson(w, http.StatusOK, todos)
}

func (c *TodoController) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTodoRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if err := c.validator.Struct(&req); err != nil {
		httputils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	todo, err := c.todoService.CreateTodo(r.Context(), &req)
	if err != nil {
		httputils.WriteError(w, http.StatusInternalServerError, "Failed to create todo")
		return
	}

	httputils.WriteJson(w, http.StatusCreated, todo)
}

func (c *TodoController) GetTodoByID(w http.ResponseWriter, r *http.Request) {
	id, err := c.parseIDFromURL(r)
	if err != nil {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}

	todo, err := c.todoService.GetTodoByID(r.Context(), id)
	if err != nil {
		if err.Error() == "invalid todo ID" {
			httputils.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		httputils.WriteError(w, http.StatusInternalServerError, "Failed to get todo")
		return
	}

	if todo == nil {
		httputils.WriteError(w, http.StatusNotFound, "Todo not found")
		return
	}

	httputils.WriteJson(w, http.StatusOK, todo)
}

func (c *TodoController) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id, err := c.parseIDFromURL(r)
	if err != nil {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}

	var req models.UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if err := c.validator.Struct(&req); err != nil {
		httputils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	todo, err := c.todoService.UpdateTodo(r.Context(), id, &req)
	if err != nil {
		// Check if it's a not found error
		if err.Error() == "todo not found" || err.Error() == "invalid todo ID" {
			httputils.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		httputils.WriteError(w, http.StatusInternalServerError, "Failed to update todo")
		return
	}

	httputils.WriteJson(w, http.StatusOK, todo)
}

func (c *TodoController) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id, err := c.parseIDFromURL(r)
	if err != nil {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}

	err = c.todoService.DeleteTodo(r.Context(), id)
	if err != nil {
		if err.Error() == "todo not found" || err.Error() == "invalid todo ID" {
			httputils.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		httputils.WriteError(w, http.StatusInternalServerError, "Failed to delete todo")
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
