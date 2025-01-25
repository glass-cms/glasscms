package middleware

import (
	"context"
	"net/http"
)

//go:generate moq -out mock_auth.go . Authentication

type Authentication interface {
	ValidateToken(ctx context.Context, token string) (bool, error)
}

// AuthMiddleware creates an http middleware that validates auth tokens in requests.
//
// It takes an Authentication interface and returns a middleware function that checks
// for valid Authorization header tokens, responding with 401 Unauthorized if
// validation fails.
func AuthMiddleware(auth Authentication) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			valid, err := auth.ValidateToken(r.Context(), token)
			if err != nil || !valid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
