package middleware

import (
	"net/http"
	"slices"

	"github.com/glass-cms/glasscms/lib/mediatype"
)

// MediaType generates a handler that writes a 415 Unsupported Media Type header
// if the request's Content-Type header does not match the provided media type.
func MediaType(accepted ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mdt, err := mediatype.Parse(r.Header.Get("Content-Type"))
			if err != nil {
				http.Error(w, "Invalid Media Type", http.StatusBadRequest)
				return
			}

			if !slices.Contains(accepted, mdt.MediaType) {
				http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
