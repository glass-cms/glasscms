package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
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

var _ ServerInterface = (*Server)(nil)

func New(
	logger *slog.Logger,
	opts ...Option,
) (*Server, error) {
	server := &Server{
		logger: logger,
	}

	server.server = &http.Server{
		Handler:      HandlerFromMux(server, http.NewServeMux()),
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

func (s *Server) ItemsDelete(w http.ResponseWriter, _ *http.Request) {
	// TODO.
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) ItemsList(w http.ResponseWriter, _ *http.Request) {
	// TODO.
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) ItemsCreate(w http.ResponseWriter, _ *http.Request) {
	// TODO.
	w.WriteHeader(http.StatusNotImplemented)
}
