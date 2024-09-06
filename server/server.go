package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/glass-cms/glasscms/lib/mediatype"
	"github.com/glass-cms/glasscms/lib/middleware"
	"github.com/glass-cms/glasscms/server/handler"
)

const (
	ShutdownGracePeriod = 10 * time.Second
	DefaultPort         = 8080
	DefaultReadTimeout  = 5 * time.Second
	DefaultWriteTimeout = 10 * time.Second
)

type Server struct {
	logger *slog.Logger
	server *http.Server
}

func New(
	logger *slog.Logger,
	handlers []handler.VersionedHandler,
	opts ...Option,
) (*Server, error) {
	server := &Server{
		logger: logger,
	}

	serveMux := http.NewServeMux()

	for _, h := range handlers {
		_ = h.Handler(serveMux, []func(http.Handler) http.Handler{
			middleware.ContentType(mediatype.ApplicationJSON, mediatype.ApplicationXML),
		})
	}

	server.server = &http.Server{
		Handler:      serveMux,
		Addr:         fmt.Sprintf(":%v", DefaultPort),
		ReadTimeout:  DefaultReadTimeout,
		WriteTimeout: DefaultWriteTimeout,
	}

	for _, opt := range opts {
		if err := opt(server); err != nil {
			return nil, err
		}
	}

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
