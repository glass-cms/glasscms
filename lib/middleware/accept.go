package middleware

import (
	"net/http"
	"slices"

	"github.com/glass-cms/glasscms/lib/mediatype"
)

// Accept generates a handler that writes a 415 Unsupported Media Type header
// if the request's Accept header does not match the provided media type.
func Accept(accepted ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Accept") == "" {
				// Set default media type of Accept header to application/json
				r.Header.Set("Accept", "application/json")
			}

			header := r.Header.Get("Accept")
			mdt, err := mediatype.Parse(header)
			if err != nil {
				http.Error(w, "Invalid media type for accept header", http.StatusBadRequest)
				return
			}

			if !slices.Contains(accepted, mdt.MediaType) {
				http.Error(w, "Unsupported Accept Media Type", http.StatusNotAcceptable)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
