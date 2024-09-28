// Package v1 implements the API handlers for the v1 version of the Glass CMS API.
package v1

import (
	"log/slog"
	"net/http"
	"reflect"

	v1 "github.com/glass-cms/glasscms/api/v1"
	"github.com/glass-cms/glasscms/internal/resource"
	"github.com/glass-cms/glasscms/item"
	"github.com/glass-cms/glasscms/server/handler"
)

type APIHandler struct {
	logger      *slog.Logger
	itemService *item.Service

	errorHandler *ErrorHandler
}

// NewAPIHandler returns a new instance of ApiHandler.
func NewAPIHandler(
	logger *slog.Logger,
	service *item.Service,
) *APIHandler {
	return &APIHandler{
		logger:       logger,
		itemService:  service,
		errorHandler: NewErrorHandler(),
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

	s.registerErrorMappers()

	return v1.HandlerWithOptions(s, v1.StdHTTPServerOptions{
		BaseURL:          "/v1",
		BaseRouter:       baseRouter,
		Middlewares:      convertedMiddlewares,
		ErrorHandlerFunc: s.errorHandler.HandleError,
	})
}

func (s *APIHandler) registerErrorMappers() {
	s.errorHandler.RegisterErrorMapper(
		reflect.TypeOf(&resource.AlreadyExistsError{}),
		ErrorMapperAlreadyExistsError,
	)

	s.errorHandler.RegisterErrorMapper(
		reflect.TypeOf(&resource.NotFoundError{}),
		ErrorMapperNotFoundError,
	)
}

var _ v1.ServerInterface = (*APIHandler)(nil)
var _ handler.VersionedHandler = (*APIHandler)(nil)
