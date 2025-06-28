package database

import (
	"context"
	"todo-list-api/internal/models"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db}
}

func (s *UserService) CreateUser(ctx context.Context, email, firstName, lastName, password string) (*models.User, error) {
	user := models.User{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Password:  password,
	}
	result := s.db.WithContext(ctx).Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	result := s.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
