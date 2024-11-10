// Package sync provides a way to synchronize items from a source to the server.
package sync

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/glass-cms/glasscms/internal/parser"
	"github.com/glass-cms/glasscms/internal/sourcer"
	"github.com/glass-cms/glasscms/internal/sourcer/fs"
	"github.com/glass-cms/glasscms/pkg/api"
	"github.com/glass-cms/glasscms/pkg/log"
)

// Syncer synchronizes items from a source to the server.
type Syncer struct {
	sourcer sourcer.Sourcer
	client  *api.ClientWithResponses
	logger  *slog.Logger
}

// NewSyncer returns a new syncer.
func NewSyncer(s sourcer.Sourcer, c *api.ClientWithResponses) (*Syncer, error) {
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

	sourceMap := s.transformItemMap(sourceItems)
	s.logger.DebugContext(ctx, "collected source items", "item_count", len(sourceMap))

	serverItems, err := s.getServerItems(ctx)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to get items from server", "error", err)
		return err
	}

	serverMap := s.transformItemMap(serverItems)
	s.logger.DebugContext(ctx, "collected server items", "item_count", len(serverMap))

	// TODO: Implement sync logic.
	// Iterate over the source items and compare them to the server items
	// when an item is found that is not on the server, create it.
	// when an item is found that is on the server, update it if the hash is different and the update time of the source is later than the server.
	// when an item is found that is on the server, delete it if the item is not on the source (by updating the item with a deleted timestamp).

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
		if i != nil {
			items = append(items, i)
		}
	}
	return items, nil
}

// getServerItems retrieves a list of items from the server.
func (s *Syncer) getServerItems(ctx context.Context) ([]*api.Item, error) {
	params := api.ItemsListParams{
		Fields: func() *[]string {
			fields := []string{"name", "hash", "update_time"}
			return &fields
		}(),
	}

	response, err := s.client.ItemsListWithResponse(ctx, &params)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status from api: %s", response.Status())
	}

	items := make([]*api.Item, len(*response.JSON200))
	for i, item := range *response.JSON200 {
		items[i] = &item
	}
	return items, nil
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
