package middleware

import (
	"net/http"
	"strings"
)

// ApiKeyMiddleware validates the X-API-Key header
func ApiKeyMiddleware(validApiKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Always allow OPTIONS requests (CORS preflight)
			if r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			// Skip API key validation for public endpoints
			publicPaths := []string{
				"/health",
				"/",
				"/api/auth/login",
				"/api/auth/register",
				"/api/auth/refresh",
			}

			for _, path := range publicPaths {
				if r.URL.Path == path || strings.HasPrefix(r.URL.Path, path) {
					next.ServeHTTP(w, r)
					return
				}
			}

			apiKey := r.Header.Get("X-API-Key")
			if apiKey == "" {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			if apiKey != validApiKey {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
