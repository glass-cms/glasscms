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
func (s *Service) CreateItem(ctx context.Context, item *Item) error {
	return s.repo.Transactionally(ctx, func(tx *sql.Tx) error {
		_, err := s.repo.CreateItem(ctx, tx, *item)
		if errors.Is(err, database.ErrDuplicatePrimaryKey) {
			return resource.NewAlreadyExistsError(item.Name, ItemResource, err)
		}

		return err
	})
}

// GetItem retrieves an item by name.
func (s *Service) GetItem(ctx context.Context, name string) (*Item, error) {
	item, err := s.repo.GetItem(ctx, name)
	if errors.Is(err, database.ErrNotFound) {
		return nil, resource.NewNotFoundError(name, ItemResource, err)
	}

	return item, err
}
