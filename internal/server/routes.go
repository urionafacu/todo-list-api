package server

import (
	"encoding/json"
	"log"
	"net/http"
	"todo-list-api/internal/handlers"
	"todo-list-api/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.CorsMiddleware)
	r.Use(middleware.ApiKeyMiddleware(s.apiKey))
	r.Use(middleware.AuthMiddleware)

	// Basic routes
	r.Get("/", s.HelloWorldHandler)
	r.Get("/health", s.healthHandler)

	// API routes
	r.Route("/api", func(r chi.Router) {
		s.registerTodoRoutes(r)
		// Future entities can be easily added here
		// s.registerUserRoutes(r)
		// s.registerProjectRoutes(r)
	})

	return r
}

func (s *Server) registerTodoRoutes(r chi.Router) {
	todoHandlers := handlers.NewTodoHandlers(s.db.GetDB())

	r.Route("/todos", func(r chi.Router) {
		// Collection routes: /api/todos
		r.Get("/", todoHandlers.GetTodos)
		r.Post("/", todoHandlers.CreateTodo)

		// Individual item routes: /api/todos/{id}
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", todoHandlers.GetTodoByID)
			r.Put("/", todoHandlers.UpdateTodo)
			r.Delete("/", todoHandlers.DeleteTodo)
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
