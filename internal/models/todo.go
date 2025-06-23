package models

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	ID          uint           `json:"id" gorm:"primaryKey;column:id"`
	Title       string         `json:"title" gorm:"not null;column:title"`
	Description string         `json:"description" gorm:"column:description"`
	Completed   bool           `json:"completed" gorm:"default:false;column:completed"`
	CreatedAt   time.Time      `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt   time.Time      `json:"updatedAt" gorm:"column:updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index;column:deleted_at"`
}

type CreateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}
