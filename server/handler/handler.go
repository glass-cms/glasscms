package handler

import "net/http"

// VersionedHandler is an interface that defines an HTTP handler.
type VersionedHandler interface {
	// HttpHandler returns an http.Handler that implements the API for a specific version.
	Handler(baseRouter *http.ServeMux, middlewares []func(http.Handler) http.Handler) http.Handler
}
