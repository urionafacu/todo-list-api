package repository

import (
	"context"
	"todo-list-api/internal/models"

	"gorm.io/gorm"
)

type PostgresAuthRepository struct {
	db *gorm.DB
}

// NewPostgresAuthRepository creates a new instance of PostgresAuthRepository
func NewPostgresAuthRepository(db *gorm.DB) AuthRepository {
	return &PostgresAuthRepository{
		db: db,
	}
}

func (r *PostgresAuthRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (r *PostgresAuthRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *PostgresAuthRepository) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
