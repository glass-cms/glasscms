package server

import (
	"errors"
	"fmt"
)

type Option func(*Server) error

// WithPort is an option that sets the port the server listens on.
func WithPort(port string) func(*Server) error {
	return func(s *Server) error {
		if port == "" {
			return errors.New("port cannot be empty")
		}

		if s.server != nil {
			s.server.Addr = fmt.Sprintf(":%s", port)
		}

		return nil
	}
}
