package server

import (
	"encoding/json"
	"log"
	"net/http"
	"todo-list-api/internal/controller"
	"todo-list-api/internal/middleware"
	"todo-list-api/internal/repository"
	"todo-list-api/internal/service"

	"github.com/go-chi/chi/v5"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.CorsMiddleware)
	r.Use(middleware.ApiKeyMiddleware(s.apiKey))

	// Basic routes
	r.Get("/", s.HelloWorldHandler)
	r.Get("/health", s.healthHandler)

	// API routes
	r.Route("/api", func(r chi.Router) {
		s.registerAuthRoutes(r)
		s.registerTodoRoutes(r)
	})

	return r
}

func (s *Server) registerTodoRoutes(r chi.Router) {
	// Initialize layers: Repository -> Service -> Controller
	todoRepo := repository.NewPostgresTodosRepository(s.db.GetDB())
	todoService := service.NewTodoService(todoRepo)
	todoController := controller.NewTodoController(todoService)

	r.Route("/todos", func(r chi.Router) {
		// Apply authentication middleware to all todo routes
		r.Use(middleware.AuthMiddleware(s.jwt))

		// Collection routes: /api/todos
		r.Get("/", todoController.GetTodos)
		r.Post("/", todoController.CreateTodo)

		// Individual item routes: /api/todos/{id}
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", todoController.GetTodoByID)
			r.Put("/", todoController.UpdateTodo)
			r.Delete("/", todoController.DeleteTodo)
		})
	})
}

func (s *Server) registerAuthRoutes(r chi.Router) {
	// Initialize layers: Repository -> Service -> Controller
	authRepo := repository.NewPostgresAuthRepository(s.db.GetDB())
	authService := service.NewAuthService(authRepo, s.jwt)
	authController := controller.NewAuthController(authService)

	r.Route("/auth", func(r chi.Router) {
		// Public auth routes (no authentication required)
		r.Post("/register", authController.Register)
		r.Post("/login", authController.Login)
		r.Post("/refresh", authController.Refresh)

		// Protected auth routes (authentication required)
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(s.jwt))
			// Add any protected auth routes here if needed
			// r.Post("/logout", authController.Logout)
			// r.Get("/profile", authController.GetProfile)
		})
	})
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"message": "Hello World"}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(s.db.Health())
	if err != nil {
		http.Error(w, "Failed to marshal health check response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}
