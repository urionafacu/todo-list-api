package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"todo-list-api/internal/handlers"
	"todo-list-api/internal/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/", s.HelloWorldHandler)

	mux.HandleFunc("/health", s.healthHandler)

	s.registerTodoRoutes(mux)

	handler := middleware.AuthMiddleware(mux)
	handler = middleware.ApiKeyMiddleware(s.apiKey)(handler)
	handler = middleware.CorsMiddleware(handler)
	// Wrap the mux with CORS middleware
	return handler
}

func (s *Server) registerTodoRoutes(mux *http.ServeMux) {
	todoHandlers := handlers.NewTodoHandlers(s.db.GetDB())

	// Handle collection endpoints: /api/todos
	mux.HandleFunc("/api/todos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			todoHandlers.GetTodos(w, r)
		case http.MethodPost:
			todoHandlers.CreateTodo(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Handle specific todo endpoints: /api/todos/{id}
	mux.HandleFunc("/api/todos/", func(w http.ResponseWriter, r *http.Request) {
		// Check if this is a path with an ID (more than 3 segments)
		if strings.Count(r.URL.Path, "/") < 3 {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		switch r.Method {
		case http.MethodGet:
			todoHandlers.GetTodoByID(w, r)
		case http.MethodPut:
			todoHandlers.UpdateTodo(w, r)
		case http.MethodDelete:
			todoHandlers.DeleteTodo(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
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
