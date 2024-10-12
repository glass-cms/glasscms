package middleware

import (
	"context"
	"net/http"

	"github.com/glass-cms/glasscms/pkg/log"
	"github.com/google/uuid"
)

// RequestIDHeader is the name of the HTTP Header which contains the request id.
var RequestIDHeader = "X-Request-Id"

func RequestID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestID := r.Header.Get(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.NewString()
		}
		ctx = context.WithValue(ctx, log.RequestIDContextKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
