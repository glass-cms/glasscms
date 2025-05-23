// Package sync provides a way to synchronize items from a source to the server.
package sync

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"net/http"
	"time"

	"github.com/glass-cms/glasscms/internal/parser"
	"github.com/glass-cms/glasscms/internal/sourcer"
	"github.com/glass-cms/glasscms/internal/sourcer/fs"
	"github.com/glass-cms/glasscms/pkg/api"
)

const (
	MetadataKeySyncID     = "sync_id"
	MetadataKeySyncSource = "sync_source"
)

var (
	ErrUnexpectedStatusCode = errors.New("unexpected status code")
)

// Syncer synchronizes items from a source to the server.
type Syncer struct {
	id *ID

	client  *api.ClientWithResponses
	config  parser.Config
	logger  *slog.Logger
	sourcer sourcer.Sourcer
}

// NewSyncer returns a new syncer.
func NewSyncer(
	id *ID,
	sourcer sourcer.Sourcer,
	c *api.ClientWithResponses,
	l *slog.Logger,
	config *parser.Config,
) (*Syncer, error) {
	// TODO: Add MetadataKeySyncSource to the config.

	// Merge the additional metadata from the config with the default metadata.
	config.AdditionalMetadata = maps.Clone(config.AdditionalMetadata)
	if config.AdditionalMetadata == nil {
		config.AdditionalMetadata = make(map[string]any)
	}
	config.AdditionalMetadata[MetadataKeySyncID] = id.String()

	return &Syncer{
		id:      id,
		sourcer: sourcer,
		client:  c,
		logger:  l,
		config:  *config,
	}, nil
}

// Sync synchronizes items from a source to the server.
func (s *Syncer) Sync(ctx context.Context, livemode bool) error {
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

	upsertItems := s.createUpsertSlice(ctx, sourceMap, serverMap)
	s.logger.DebugContext(ctx, "upserting items", "item_count", len(upsertItems))

	if !livemode {
		s.logger.InfoContext(ctx, "dry run complete, exiting")
		return nil
	}

	if len(upsertItems) == 0 {
		s.logger.InfoContext(ctx, "no items to upsert, exiting")
		return nil
	}

	return s.upsertItems(ctx, upsertItems)
}

// createUpsertSlice generates a slice of items that need to be upserted (created or updated) or deleted
// based on the differences between the source and server maps.
func (s *Syncer) createUpsertSlice(ctx context.Context, sourceMap, serverMap map[string]*api.Item) []*api.Item {
	var upsertItems []*api.Item

	// Iterate over the source items and compare them to the server items.
	for name, sourceItem := range sourceMap {
		serverItem, exists := serverMap[name]
		if !exists {
			s.logger.DebugContext(ctx, "creating item", "name", name)
			upsertItems = append(upsertItems, sourceItem)
			continue
		}

		// TODO: Add option to ignore update time.

		if sourceItem.Hash != serverItem.Hash && sourceItem.UpdateTime.After(serverItem.UpdateTime) {
			s.logger.DebugContext(ctx, "updating item", "name", name)
			upsertItems = append(upsertItems, sourceItem)
		}
	}

	// Check for items that are on the server but not on the source, these items should be deleted.
	for name, serverItem := range serverMap {
		_, ok := sourceMap[name]
		if !ok {
			s.logger.DebugContext(ctx, "deleting item", "name", name)

			now := time.Now()
			serverItem.DeleteTime = &now
			upsertItems = append(upsertItems, serverItem)
		}
	}

	return upsertItems
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
		i, err = parser.ParseWithConfig(src, s.config)
		if errors.Is(err, parser.ErrItemHidden) {
			// Item is hidden, skip it
			s.logger.DebugContext(ctx, "skipping hidden item", "name", src.Name())
			continue
		}
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
		s.logger.ErrorContext(ctx, "failed to list items", "error", err)
		return nil, err
	}

	if response.StatusCode() != http.StatusOK {
		s.logger.ErrorContext(
			ctx, "received unexpected status code while listing items", "err", err, "status_code", response.StatusCode())
		return nil, fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, response.StatusCode())
	}

	items := make([]*api.Item, len(*response.JSON200))
	for i, item := range *response.JSON200 {
		items[i] = &item
	}
	return items, nil
}

// upsertItems upserts a slice of items to the server.
func (s *Syncer) upsertItems(ctx context.Context, items []*api.Item) error {
	upsertItems := make([]api.ItemUpsert, len(items))
	for i, item := range items {
		upsertItems[i] = api.ItemUpsert{
			Content:     item.Content,
			CreateTime:  item.CreateTime,
			DeleteTime:  item.DeleteTime,
			DisplayName: item.DisplayName,
			Metadata:    item.Metadata,
			Name:        item.Name,
			Properties:  item.Properties,
			UpdateTime:  item.UpdateTime,
		}
	}

	response, err := s.client.ItemsUpsertWithResponse(ctx, upsertItems)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to upsert items", "error", err)
		return err
	}

	if response.StatusCode() != http.StatusOK {
		s.logger.ErrorContext(
			ctx, "received unexpected status code while upserting items", "error", err, "status_code", response.StatusCode())
		return fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, response.StatusCode())
	}

	return nil
}

// transformItemMap transforms a slice of items into a map where the key is the item name.
func (s *Syncer) transformItemMap(items []*api.Item) map[string]*api.Item {
	if items == nil {
		return make(map[string]*api.Item)
	}

	itemMap := make(map[string]*api.Item, len(items))
	for _, i := range items {
		if i != nil && i.Name != "" { // Add validation for empty names
			itemMap[i.Name] = i
		}
	}
	return itemMap
}
