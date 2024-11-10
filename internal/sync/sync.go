// Package sync provides a way to synchronize items from a source to the server.
package sync

import (
	"context"
	"errors"
	"log/slog"

	"github.com/glass-cms/glasscms/internal/parser"
	"github.com/glass-cms/glasscms/internal/sourcer"
	"github.com/glass-cms/glasscms/internal/sourcer/fs"
	"github.com/glass-cms/glasscms/pkg/api"
	"github.com/glass-cms/glasscms/pkg/client"
	"github.com/glass-cms/glasscms/pkg/log"
)

// Syncer synchronizes items from a source to the server.
type Syncer struct {
	sourcer sourcer.Sourcer
	client  *client.Client
	logger  *slog.Logger
}

// NewSyncer returns a new syncer.
func NewSyncer(s sourcer.Sourcer, c *client.Client) (*Syncer, error) {
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
func (s *Syncer) Sync(ctx context.Context, _ bool) error {
	s.logger.InfoContext(ctx, "syncing items")

	sourceItems, err := s.collectSourceItems(ctx)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to collect items from sourcer", "error", err)
		return err
	}

	_ = s.transformItemMap(sourceItems)

	_, err = s.getServerItems(ctx)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to get items from server", "error", err)
		return err
	}

	return nil
}

// collectItems returns a slice of parsed items collected from the source
// or an error if the retrieval process fails.
func (s *Syncer) collectSourceItems(ctx context.Context) ([]*api.Item, error) {
	size := s.sourcer.Size()
	items := make([]*api.Item, 0, size)

	for {
		// Check if context is cancelled.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// Get the next source.
		var src sourcer.Source
		src, err := s.sourcer.Next()
		if errors.Is(err, fs.ErrDone) {
			break
		}

		if err != nil {
			return nil, err
		}

		var i *api.Item
		i, err = parser.Parse(src)
		if err != nil {
			s.logger.WarnContext(ctx, "failed to parse item from source", "name", src.Name(), "error", err)
			continue
		}

		items = append(items, i)
	}
	return items, nil
}

// getServerItems retrieves a list of items from the server.
func (s *Syncer) getServerItems(_ context.Context) ([]*api.Item, error) {
	return nil, nil
}

// transformItemMap transforms a slice of items into a map where the key is the item name.
func (s *Syncer) transformItemMap(items []*api.Item) map[string]*api.Item {
	itemMap := make(map[string]*api.Item, len(items))
	for _, i := range items {
		if i != nil {
			itemMap[i.Name] = i
		}
	}
	return itemMap
}
