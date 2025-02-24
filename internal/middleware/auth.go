package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"movie-ticket-booking/internal/services"
	"net/http"
	"strings"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
)

func AuthMiddleware(authService *services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip authentication for login and register mutations
			if r.URL.Path == "/query" {
				// Read the request body
				body, err := io.ReadAll(r.Body)
				if err != nil {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`{"errors":[{"message":"Error reading request body"}]}`))                
					return
				}
				
				// Parse the request body
				var reqBody struct {
					Query     string                 `json:"query"`
					Variables map[string]interface{} `json:"variables"`
				}
				if err := json.Unmarshal(body, &reqBody); err != nil {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`{"errors":[{"message":"Invalid JSON request body"}]}`))                
					return
				}

				// Convert query to lowercase and remove all whitespace for consistent matching
				query := strings.ToLower(reqBody.Query)
				query = strings.ReplaceAll(query, " ", "")
				query = strings.ReplaceAll(query, "\n", "")
				query = strings.ReplaceAll(query, "\t", "")

				// Check if it's a login or register mutation
				if strings.Contains(query, "mutation{login(") || 
				   strings.Contains(query, "mutation{register(") || 
				   strings.Contains(query, "mutationlogin(") || 
				   strings.Contains(query, "mutationregister(") {
					// Create a new reader with the same body content
					r.Body = io.NopCloser(bytes.NewBuffer(body))
					next.ServeHTTP(w, r)
					return
				}
				// Reset the body for further processing
				r.Body = io.NopCloser(bytes.NewBuffer(body))
			}

			// Get the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"errors":[{"message":"Authorization header is required"}]}`))                
				return
			}

			// Extract the token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"errors":[{"message":"Invalid authorization header format. Use 'Bearer <token>'"}]}`))                
				return
			}

			tokenString := parts[1]

			// Validate the token
			claims, err := authService.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Add user ID to context
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID retrieves the user ID from the context
func GetUserID(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value(UserIDKey).(uint)
	return userID, ok
}