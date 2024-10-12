package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/glass-cms/glasscms/internal/database"
	"github.com/glass-cms/glasscms/internal/item"
	"github.com/glass-cms/glasscms/internal/item/repository/query"
)

var _ item.Repository = &ItemRepository{}

type ItemRepository struct {
	db           *sql.DB
	errorHandler database.ErrorHandler
}

func NewRepository(db *sql.DB, errorHandler database.ErrorHandler) *ItemRepository {
	return &ItemRepository{
		db:           db,
		errorHandler: errorHandler,
	}
}

// Transactionally executes a function within a database transaction. It commits the transaction
// if the function succeeds, otherwise it rolls back. If rollback fails, both errors are returned.
func (r *ItemRepository) Transactionally(ctx context.Context, f func(tx *sql.Tx) error) (err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("panic occurred: %v, rollback error: %w", p, rbErr)
			} else {
				err = fmt.Errorf("panic occurred: %v", p)
			}
		}
	}()

	err = f(tx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction rollback error: %w, original error: %w", rbErr, err)
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// CreateItem creates a new item in the database.
// If a transaction is provided, the item will be created within the transaction.
func (r *ItemRepository) CreateItem(ctx context.Context, tx *sql.Tx, item item.Item) (*item.Item, error) {
	q := query.New(tx)

	propertiesJSON, err := json.Marshal(item.Properties)
	if err != nil {
		return nil, r.errorHandler.HandleError(ctx, err)
	}

	metadataJSON, err := json.Marshal(item.Metadata)
	if err != nil {
		return nil, r.errorHandler.HandleError(ctx, err)
	}

	params := query.CreateItemParams{
		Name:        item.Name,
		DisplayName: item.DisplayName,
		CreateTime:  item.CreateTime,
		UpdateTime:  item.UpdateTime,
		DeleteTime:  sql.NullTime{},
		Hash: sql.NullString{
			String: item.Hash,
			Valid:  true,
		},
		Content: sql.NullString{
			String: item.Content,
			Valid:  true,
		},
		Properties: propertiesJSON,
		Metadata:   metadataJSON,
	}

	i, err := q.CreateItem(ctx, params)
	if err != nil {
		return nil, r.errorHandler.HandleError(ctx, err)
	}

	newItem, err := ConvertQueryItem(i)
	if err != nil {
		return nil, r.errorHandler.HandleError(ctx, err)
	}

	return newItem, nil
}

// GetItem retrieves an item from the database by its resource name.
func (r *ItemRepository) GetItem(ctx context.Context, name string) (*item.Item, error) {
	q := query.New(r.db)

	i, err := q.GetItem(ctx, name)
	if err != nil {
		return nil, r.errorHandler.HandleError(ctx, err)
	}

	foundItem, err := ConvertQueryItem(i)
	if err != nil {
		return nil, r.errorHandler.HandleError(ctx, err)
	}

	return foundItem, nil
}

// UpdateItem updates an existing item in the database.
func (r *ItemRepository) UpdateItem(ctx context.Context, item *item.Item) error {
	query := `
		UPDATE items
		SET update_time = $1, hash = $2, name = $3, display_name = $4, content = $5, properties = $6, metadata = $7
		WHERE name = $8
	`

	propertiesJSON, err := json.Marshal(item.Properties)
	if err != nil {
		return fmt.Errorf("failed to marshal properties: %w", err)
	}

	metadataJSON, err := json.Marshal(item.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
		item.UpdateTime,
		item.Hash,
		item.Name,
		item.DisplayName,
		item.Content,
		propertiesJSON,
		metadataJSON,
		item.Name,
	)

	if err != nil {
		return r.errorHandler.HandleError(ctx, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("item not found")
	}

	return nil
}
