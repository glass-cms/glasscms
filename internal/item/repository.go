package item

import (
	"context"
	"database/sql"
)

type Repository interface {
	Transactionally(ctx context.Context, f func(tx *sql.Tx) error) (err error)
	CreateItem(ctx context.Context, tx *sql.Tx, item *Item) error

	// TODO: Add transaction to all methods.
	// TODO: Add method to get a transaction.
	GetItem(ctx context.Context, name string) (*Item, error)
	UpdateItem(ctx context.Context, item *Item) error
}
