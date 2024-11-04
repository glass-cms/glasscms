// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package query

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.createItemStmt, err = db.PrepareContext(ctx, createItem); err != nil {
		return nil, fmt.Errorf("error preparing query CreateItem: %w", err)
	}
	if q.deleteItemStmt, err = db.PrepareContext(ctx, deleteItem); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteItem: %w", err)
	}
	if q.getItemStmt, err = db.PrepareContext(ctx, getItem); err != nil {
		return nil, fmt.Errorf("error preparing query GetItem: %w", err)
	}
	if q.listItemsStmt, err = db.PrepareContext(ctx, listItems); err != nil {
		return nil, fmt.Errorf("error preparing query ListItems: %w", err)
	}
	if q.updateItemStmt, err = db.PrepareContext(ctx, updateItem); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateItem: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.createItemStmt != nil {
		if cerr := q.createItemStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createItemStmt: %w", cerr)
		}
	}
	if q.deleteItemStmt != nil {
		if cerr := q.deleteItemStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteItemStmt: %w", cerr)
		}
	}
	if q.getItemStmt != nil {
		if cerr := q.getItemStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getItemStmt: %w", cerr)
		}
	}
	if q.listItemsStmt != nil {
		if cerr := q.listItemsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listItemsStmt: %w", cerr)
		}
	}
	if q.updateItemStmt != nil {
		if cerr := q.updateItemStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateItemStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db             DBTX
	tx             *sql.Tx
	createItemStmt *sql.Stmt
	deleteItemStmt *sql.Stmt
	getItemStmt    *sql.Stmt
	listItemsStmt  *sql.Stmt
	updateItemStmt *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:             tx,
		tx:             tx,
		createItemStmt: q.createItemStmt,
		deleteItemStmt: q.deleteItemStmt,
		getItemStmt:    q.getItemStmt,
		listItemsStmt:  q.listItemsStmt,
		updateItemStmt: q.updateItemStmt,
	}
}
