package item

import (
	"context"
	"database/sql"
	"errors"

	"github.com/glass-cms/glasscms/internal/database"
	"github.com/glass-cms/glasscms/pkg/resource"
)

// Service is a service for managing items.
type Service struct {
	db   *sql.DB
	repo Repository
}

func NewService(db *sql.DB, repo Repository) *Service {
	return &Service{
		db:   db,
		repo: repo,
	}
}

// CreateItem creates a new item.
func (s *Service) CreateItem(ctx context.Context, item Item) (*Item, error) {
	createdItem := &Item{}

	err := database.Transactionally(ctx, s.db, func(tx *sql.Tx) error {
		var err error

		createdItem, err = s.repo.CreateItem(ctx, tx, item)
		if errors.Is(err, database.ErrDuplicatePrimaryKey) {
			return resource.NewAlreadyExistsError(item.Name, ItemResource, err)
		}

		return err
	})
	if err != nil {
		return &Item{}, err
	}

	return createdItem, nil
}

// GetItem retrieves an item by name.
func (s *Service) GetItem(ctx context.Context, name string) (*Item, error) {
	var item *Item

	err := database.Transactionally(ctx, s.db, func(tx *sql.Tx) error {
		var err error

		item, err = s.repo.GetItem(ctx, tx, name)
		if errors.Is(err, database.ErrNotFound) {
			return resource.NewNotFoundError(name, ItemResource, err)
		}

		return err
	})

	return item, err
}

// ListItems retrieves a list of items.
func (s *Service) ListItems(ctx context.Context, fieldmask []string) ([]*Item, error) {
	var items []*Item

	err := database.Transactionally(ctx, s.db, func(tx *sql.Tx) error {
		var err error

		items, err = s.repo.ListItems(ctx, tx, fieldmask)
		if err != nil {
			return err
		}

		return nil
	})

	return items, err
}

// UpsertItems upserts a list of items.
func (s *Service) UpsertItems(ctx context.Context, items []Item) ([]*Item, error) {
	upsertedItems := make([]*Item, len(items))

	err := database.Transactionally(ctx, s.db, func(tx *sql.Tx) error {
		var err error

		for i, item := range items {
			upsertedItems[i], err = s.repo.UpsertItem(ctx, tx, item)
		}

		return err
	})
	if err != nil {
		return nil, err
	}

	return upsertedItems, nil
}

// DeleteItems deletes a list of items by the unique names.
func (s *Service) DeleteItems(ctx context.Context, names []string) error {
	return database.Transactionally(ctx, s.db, func(tx *sql.Tx) error {
		return s.repo.DeleteItems(ctx, tx, names)
	})
}
