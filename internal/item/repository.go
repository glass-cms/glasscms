package item

import (
	"context"
	"database/sql"
)

// TODO Everything should accept a transaction.

type Repository interface {
	Transactionally(ctx context.Context, f func(tx *sql.Tx) error) (err error)
	CreateItem(ctx context.Context, tx *sql.Tx, item Item) (*Item, error)
	GetItem(ctx context.Context, name string) (*Item, error)
	UpdateItem(ctx context.Context, item *Item) error
	ListItems(ctx context.Context, fieldmasks []string) ([]Item, error)
}
