package middleware

import (
	"net/http"
	"slices"

	"github.com/glass-cms/glasscms/pkg/mediatype"
)

// ContentType generates a handler that writes a 415 Unsupported Media Type header
// if the request's Content-Type header does not match the provided media type.
func ContentType(accepted ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				header := r.Header.Get("Content-Type")
				if header == "" {
					http.Error(w, "Invalid Media Type", http.StatusBadRequest)
					return
				}

				mdt, err := mediatype.Parse(header)
				if err != nil {
					http.Error(w, "Invalid Media Type", http.StatusBadRequest)
					return
				}

				if !slices.Contains(accepted, mdt.MediaType) {
					http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
