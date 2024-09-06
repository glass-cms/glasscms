package item

import (
	"context"
	"errors"

	"github.com/glass-cms/glasscms/database"
	"github.com/glass-cms/glasscms/lib/resource"
)

type Service struct {
	repo Repository
}

// NewService returns a new instance of Service.
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateItem(ctx context.Context, item *Item) error {
	err := s.repo.CreateItem(ctx, nil, item)
	if errors.Is(err, database.ErrDuplicatePrimaryKey) {
		return resource.NewAlreadyExistsError(item.Name, ItemResource, err)
	}

	return err
}
