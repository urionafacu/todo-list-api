package database

import (
	"log"
	"todo-list-api/internal/models"

	"gorm.io/gorm"
)

// AutoMigrate executes the auto migration for all models
func AutoMigrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		&models.Todo{},
	)
	if err != nil {
		log.Printf("Failed to run migrations: %v", err)
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}
