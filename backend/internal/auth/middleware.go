package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/vectorsight/serverless-tool/backend/internal/types"
)

type contextKey string

const (
	// UserIDKey is the context key for user ID.
	UserIDKey contextKey = "userID"
	// EmailKey is the context key for user email.
	EmailKey contextKey = "email"
)

// AuthMiddleware intercepts requests, validates the JWT, and sets claims in context.
func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondWithError(w, http.StatusUnauthorized, "Missing Authorization Header")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				respondWithError(w, http.StatusUnauthorized, "Invalid Authorization Header Format")
				return
			}

			tokenStr := parts[1]
			claims, err := ValidateToken(tokenStr, jwtSecret)
			if err != nil {
				respondWithError(w, http.StatusUnauthorized, "Invalid or Expired Token: "+err.Error())
				return
			}

			// Add claims to request context
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, EmailKey, claims.Email)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserIDFromContext retrieves the user ID from the request context.
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(UserIDKey).(string)
	return val, ok
}

// GetEmailFromContext retrieves the user email from the request context.
func GetEmailFromContext(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(EmailKey).(string)
	return val, ok
}

// respondWithError formats and returns a standard JSON error response.
func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	errBody := types.ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	}
	_ = json.NewEncoder(w).Encode(errBody)
}
