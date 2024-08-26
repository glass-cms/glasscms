package item

import "context"

type Service interface {
	CreateItem(ctx context.Context, item *Item) error
}

var _ Service = &service{}

type service struct {
	repo Repository
}

// NewService returns a new instance of Service.
func NewService(repo Repository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateItem(ctx context.Context, item *Item) error {
	return nil
}
