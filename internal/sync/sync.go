// Package sync provides a way to synchronize items from a source to the server.
package sync

import (
	"context"
	"log/slog"

	"github.com/glass-cms/glasscms/internal/sourcer"
	"github.com/glass-cms/glasscms/pkg/client"
	"github.com/glass-cms/glasscms/pkg/log"
)

// Syncer synchronizes items from a source to the server.
type Syncer struct {
	sourcer *sourcer.Sourcer
	client  *client.Client
	logger  *slog.Logger
}

// NewSyncer returns a new syncer.
func NewSyncer(s *sourcer.Sourcer, c *client.Client) (*Syncer, error) {
	logger, err := log.NewLogger()
	if err != nil {
		return nil, err
	}

	return &Syncer{
		sourcer: s,
		client:  c,
		logger:  logger,
	}, nil
}

// Sync synchronizes items from a source to the server.
func (s *Syncer) Sync(_ context.Context, _ bool) error {
	return nil
}
