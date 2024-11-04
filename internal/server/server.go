package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
	"time"

	"github.com/glass-cms/glasscms/internal/item"
	"github.com/glass-cms/glasscms/pkg/api"
	"github.com/glass-cms/glasscms/pkg/fieldmask"
	"github.com/glass-cms/glasscms/pkg/mediatype"
	"github.com/glass-cms/glasscms/pkg/middleware"
	"github.com/glass-cms/glasscms/pkg/resource"
)

const (
	ShutdownGracePeriod = 10 * time.Second
	DefaultPort         = 8080
	DefaultReadTimeout  = 5 * time.Second
	DefaultWriteTimeout = 10 * time.Second
)

var _ api.ServerInterface = (*Server)(nil)

type Server struct {
	logger *slog.Logger
	server *http.Server

	itemService  *item.Service
	errorHandler *ErrorHandler

	handler http.Handler
}

func New(
	logger *slog.Logger,
	itemService *item.Service,
	opts ...Option,
) (*Server, error) {
	serveMux := http.NewServeMux()

	server := &Server{
		logger:       logger,
		itemService:  itemService,
		errorHandler: NewErrorHandler(),
	}

	middlewares := []func(http.Handler) http.Handler{
		middleware.RequestID,
		middleware.ContentType(mediatype.ApplicationJSON),
		middleware.Accept(mediatype.ApplicationJSON),
	}
	convertedMiddlewares := make([]api.MiddlewareFunc, len(middlewares))
	for i, mw := range middlewares {
		convertedMiddlewares[i] = api.MiddlewareFunc(mw)
	}

	server.handler = api.HandlerWithOptions(server, api.StdHTTPServerOptions{
		BaseURL:     "",
		BaseRouter:  serveMux,
		Middlewares: convertedMiddlewares,
	})

	server.server = &http.Server{
		Handler:      server.handler,
		Addr:         fmt.Sprintf(":%v", DefaultPort),
		ReadTimeout:  DefaultReadTimeout,
		WriteTimeout: DefaultWriteTimeout,
	}

	for _, opt := range opts {
		if err := opt(server); err != nil {
			return nil, err
		}
	}

	server.registerErrorMappers()
	return server, nil
}

// ListenAndServe starts the server.
func (s *Server) ListenAndServer() error {
	s.logger.Info("server is listening on :8080")
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the underlying server without interrupting any active connections.
func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownGracePeriod)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error("could not gracefully shutdown the server:", "err", err)
		return
	}

	s.logger.Info("server stopped")
}

func (s *Server) Handler() http.Handler {
	return s.handler
}

func (s *Server) registerErrorMappers() {
	s.errorHandler.RegisterErrorMapper(
		reflect.TypeOf(&resource.AlreadyExistsError{}),
		ErrorMapperAlreadyExistsError,
	)
	s.errorHandler.RegisterErrorMapper(
		reflect.TypeOf(&resource.NotFoundError{}),
		ErrorMapperNotFoundError,
	)
	s.errorHandler.RegisterErrorMapper(
		reflect.TypeOf(&fieldmask.InvalidFieldMaskError{}),
		ErrorMapperInvalidFieldMaskError,
	)
	s.errorHandler.RegisterErrorMapper(
		reflect.TypeOf(&DeserializeError{}),
		ErrorMapperDeserializeError,
	)
}
