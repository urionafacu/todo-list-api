package models

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"varchar(200);not null"`
	Description string         `json:"description" gorm:"type:varchar(1000)"`
	Priority    string         `json:"priority" gorm:"type:varchar(10);default:'low'"`
	DueDate     *time.Time     `json:"dueDate"`
	Category    string         `json:"category" gorm:"type:varchar(100)"`
	Completed   bool           `json:"completed" gorm:"default:false"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreateTodoRequest struct {
	Title       string  `json:"title" validate:"required,min=1,max=200"`
	Description string  `json:"description" validate:"max=1000"`
	Priority    string  `json:"priority" validate:"omitempty,oneof=low medium high"`
	DueDate     *string `json:"dueDate" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	Category    string  `json:"category" validate:"omitempty,max=100"`
}

type UpdateTodoRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Completed   bool    `json:"completed"`
	Priority    string  `json:"priority" validate:"omitempty,oneof=low medium high"`
	DueDate     *string `json:"dueDate" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	Category    string  `json:"category" validate:"omitempty,max=100"`
}
