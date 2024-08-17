// Package v1 implements the API handlers for the v1 version of the Glass CMS API.
package v1

import (
	"log/slog"
	"net/http"

	v1 "github.com/glass-cms/glasscms/api/v1"
	"github.com/glass-cms/glasscms/item"
	"github.com/glass-cms/glasscms/server/handler"
)

type APIHandler struct {
	logger     *slog.Logger
	repository *item.Repository
}

// NewAPIHandler returns a new instance of ApiHandler.
func NewAPIHandler(
	logger *slog.Logger,
	repo *item.Repository,
) *APIHandler {
	return &APIHandler{
		logger:     logger,
		repository: repo,
	}
}

// Handler returns an http.Handler that implements the API.
func (s *APIHandler) Handler(
	baseRouter *http.ServeMux,
	middlewares []func(http.Handler) http.Handler,
) http.Handler {
	convertedMiddlewares := make([]v1.MiddlewareFunc, len(middlewares))
	for i, mw := range middlewares {
		convertedMiddlewares[i] = v1.MiddlewareFunc(mw)
	}

	return v1.HandlerWithOptions(s, v1.StdHTTPServerOptions{
		BaseURL:     "/v1",
		BaseRouter:  baseRouter,
		Middlewares: convertedMiddlewares,
	})
}

var _ v1.ServerInterface = (*APIHandler)(nil)
var _ handler.VersionedHandler = (*APIHandler)(nil)
