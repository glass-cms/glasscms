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
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateItem creates a new item.
func (s *Service) CreateItem(ctx context.Context, item Item) (*Item, error) {
	createdItem := &Item{}

	err := s.repo.Transactionally(ctx, func(tx *sql.Tx) error {
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

	err := s.repo.Transactionally(ctx, func(tx *sql.Tx) error {
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

	err := s.repo.Transactionally(ctx, func(tx *sql.Tx) error {
		var err error

		items, err = s.repo.ListItems(ctx, tx, fieldmask)
		if err != nil {
			return err
		}

		return nil
	})

	return items, err
}
