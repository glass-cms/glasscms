package server

import (
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

func New(
	logger *slog.Logger,
	opts ...Option,
) (*Server, error) {
	s := &Server{
		logger: logger,
		server: &http.Server{
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// TODO.
			}),
			Addr:         fmt.Sprintf(":%v", DefaultPort),
			ReadTimeout:  DefaultReadTimeout,
			WriteTimeout: DefaultWriteTimeout,
		},
	}

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *Server) ListenAndServer() error {
	s.logger.Info("server is listening on :8080")
	return s.server.ListenAndServe()
}
