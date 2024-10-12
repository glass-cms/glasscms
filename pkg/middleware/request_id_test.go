package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glass-cms/glasscms/pkg/log"
	"github.com/glass-cms/glasscms/pkg/middleware"
	"github.com/google/uuid"
)

func TestRequestID(t *testing.T) {
	tests := []struct {
		name           string
		requestID      string
		expectedHeader string
	}{
		{
			name:           "existing request ID",
			requestID:      "existing-request-id",
			expectedHeader: "existing-request-id",
		},
		{
			name:           "generate new request ID",
			requestID:      "",
			expectedHeader: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := middleware.RequestID(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
				requestID := r.Context().Value(log.RequestIDContextKey).(string) //nolint:errcheck // Ignore.
				if tt.expectedHeader == "" {
					if _, err := uuid.Parse(requestID); err != nil {
						t.Errorf("expected a valid UUID, got %v", requestID)
					}
				} else if requestID != tt.expectedHeader {
					t.Errorf("expected request ID %v, got %v", tt.expectedHeader, requestID)
				}
			}))

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.requestID != "" {
				req.Header.Set(middleware.RequestIDHeader, tt.requestID)
			}
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)
		})
	}
}
