package repository

import (
	"context"
	"errors"
	"todo-list-api/internal/models"

	"gorm.io/gorm"
)

type postgresTodosRepository struct {
	db *gorm.DB
}

// NewPostgresTodosRepository creates a new PostgreSQL implementation of TodoRepository
func NewPostgresTodosRepository(db *gorm.DB) TodoRepository {
	return &postgresTodosRepository{
		db: db,
	}
}

func (r *postgresTodosRepository) Create(ctx context.Context, todo *models.Todo) (*models.Todo, error) {
	result := r.db.WithContext(ctx).Create(todo)
	if result.Error != nil {
		return nil, result.Error
	}
	return todo, nil
}

func (r *postgresTodosRepository) GetAll(ctx context.Context) ([]models.Todo, error) {
	var todos []models.Todo
	result := r.db.WithContext(ctx).Order("created_at DESC").Find(&todos)
	if result.Error != nil {
		return nil, result.Error
	}
	return todos, nil
}

func (r *postgresTodosRepository) GetByID(ctx context.Context, id uint) (*models.Todo, error) {
	var todo models.Todo
	result := r.db.WithContext(ctx).First(&todo, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil for not found
		}
		return nil, result.Error
	}
	return &todo, nil
}

func (r *postgresTodosRepository) Update(ctx context.Context, id uint, todo *models.Todo) (*models.Todo, error) {
	result := r.db.WithContext(ctx).Model(&models.Todo{}).Where("id = ?", id).Updates(todo)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil // Todo not found
	}

	// Return the updated todo
	return r.GetByID(ctx, id)
}

func (r *postgresTodosRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Todo{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("todo not found")
	}

	return nil
}

func (r *postgresTodosRepository) GetByUserID(ctx context.Context, userID uint) ([]models.Todo, error) {
	var todos []models.Todo
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&todos)
	if result.Error != nil {
		return nil, result.Error
	}
	return todos, nil
}
