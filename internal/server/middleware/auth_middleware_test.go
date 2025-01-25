package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glass-cms/glasscms/internal/server/middleware"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	mockAuth := &middleware.AuthenticationMock{}

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := middleware.AuthMiddleware(mockAuth)
	wrappedHandler := middleware(handler)

	t.Run("Valid Token", func(t *testing.T) {
		mockAuth.ValidateTokenFunc = func(_ context.Context, _ string) (bool, error) {
			return true, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "valid-token")
		w := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		mockAuth.ValidateTokenFunc = func(_ context.Context, _ string) (bool, error) {
			return false, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "invalid-token")
		w := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Missing Token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
