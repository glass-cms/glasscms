package item

import (
	"context"
	"database/sql"
)

type Repository interface {
	Transactionally(ctx context.Context, f func(tx *sql.Tx) error) (err error)
	CreateItem(ctx context.Context, tx *sql.Tx, item Item) (*Item, error)
	GetItem(ctx context.Context, name string) (*Item, error)
	UpdateItem(ctx context.Context, item *Item) error
}
