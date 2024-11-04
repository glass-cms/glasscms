package item

import (
	"context"
	"database/sql"
)

type Repository interface {
	Transactionally(ctx context.Context, f func(tx *sql.Tx) error) (err error)
	CreateItem(ctx context.Context, tx *sql.Tx, item Item) (*Item, error)
	GetItem(ctx context.Context, tx *sql.Tx, name string) (*Item, error)
	UpdateItem(ctx context.Context, tx *sql.Tx, item Item) (*Item, error)
	DeleteItem(ddctx context.Context, tx *sql.Tx, name string) error
	ListItems(ctx context.Context, tx *sql.Tx, fieldmasks []string) ([]*Item, error)
}

// TODO: Delete Many Items.
