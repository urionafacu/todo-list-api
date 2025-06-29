package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"todo-list-api/internal/controller"
	"todo-list-api/internal/middleware"
	"todo-list-api/internal/repository"
	"todo-list-api/internal/service"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	// Global middleware (CORS only)
	r.Use(middleware.CorsMiddleware)

	// Static files (no API key required)
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "web/static"))
	r.Handle("/static/*", http.StripPrefix("/static", http.FileServer(filesDir)))

	// Public routes (no API key required)
	r.Get("/", s.HelloWorldHandler)

	// Swagger documentation (no API key required)
	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/", http.StatusMovedPermanently)
	})
	r.Handle("/docs/*", httpSwagger.WrapHandler)

	// Routes that require API key
	r.Group(func(r chi.Router) {
		r.Use(middleware.ApiKeyMiddleware(s.apiKey))

		// Health check
		r.Get("/health", s.healthHandler)

		// API routes
		r.Route("/api", func(r chi.Router) {
			s.registerAuthRoutes(r)
			s.registerTodoRoutes(r)
		})
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
	workDir, _ := os.Getwd()
	htmlPath := filepath.Join(workDir, "web", "index.html")

	htmlContent, err := os.ReadFile(htmlPath)
	if err != nil {
		log.Printf("Failed to read HTML file: %v", err)
		// Fallback to JSON response if HTML file not found
		resp := map[string]string{"message": "Welcome to TODO List API", "status": "HTML file not found"}
		jsonResp, _ := json.Marshal(resp)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResp)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if _, err := w.Write(htmlContent); err != nil {
		log.Printf("Failed to write HTML response: %v", err)
	}
}

// @Summary Health check
// @Description Check the health status of the API and database
// @Tags health
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /health [get]
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
