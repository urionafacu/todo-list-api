package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"todo-list-api/internal/database"
	"todo-list-api/internal/utils"
)

type Server struct {
	port   int
	apiKey string
	db     database.Service
	jwt    *utils.JWT
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	apiKey := os.Getenv("API_KEY")

	// Initialize JWT secret once
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable must be set")
	}

	NewServer := &Server{
		port:   port,
		apiKey: apiKey,
		db:     database.New(),
		jwt:    &utils.JWT{Secret: jwtSecret},
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
