package database

import (
	"context"
	"database/sql"
	"todo-list-api/internal/models"
)

type TodoService struct {
	db *sql.DB
}

func NewTodoService(db *sql.DB) *TodoService {
	return &TodoService{db}
}

func (s *TodoService) CreateTodo(ctx context.Context, title, description string) (*models.Todo, error) {
	query := `
		INSERT INTO todos (title, description, completed, created_at, updated_at)
		VALUES ($1, $2, false, NOW(), NOW())
		RETURNING id, title, description, completed, created_at, updated_at`

	var todo models.Todo
	err := s.db.QueryRowContext(ctx, query, title, description).Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (s *TodoService) GetTodos(ctx context.Context) ([]models.Todo, error) {
	query := `
		SELECT id, title, description, completed, created_at, updated_at
		FROM todos
		ORDER BY created_at DESC`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

func (s *TodoService) GetTodoByID(ctx context.Context, id uint) (*models.Todo, error) {
	query := `
		SELECT id, title, description, completed, created_at, updated_at
		FROM todos
		WHERE id = $1`

	var todo models.Todo
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &todo, nil
}

func (s *TodoService) UpdateTodo(ctx context.Context, id uint, updates models.UpdateTodoRequest) (*models.Todo, error) {
	// First get the current todo
	todo, err := s.GetTodoByID(ctx, id)
	if err != nil || todo == nil {
		return nil, err
	}

	// Apply updates
	todo.Title = updates.Title
	todo.Description = updates.Description
	todo.Completed = updates.Completed

	query := `
		UPDATE todos
		SET title = $1, description = $2, completed = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING id, title, description, completed, created_at, updated_at`

	err = s.db.QueryRowContext(ctx, query, todo.Title, todo.Description, todo.Completed, id).Scan(
		&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (s *TodoService) DeleteTodo(ctx context.Context, id uint) error {
	query := `DELETE FROM todos WHERE id = $1`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
