package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/glass-cms/glasscms/internal/database"
	"github.com/glass-cms/glasscms/internal/item"
)

var _ item.Repository = &ItemRepository{}

type executor interface {
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

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
func (r *ItemRepository) CreateItem(ctx context.Context, tx *sql.Tx, item *item.Item) error {
	var exec executor
	if tx != nil {
		exec = tx
	} else {
		exec = r.db
	}

	query := `
        INSERT INTO items (
			name, 
			display_name, 
			create_time, 
			update_time, 
			delete_time, 
			hash, 
			content, 
			properties, 
			metadata)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `

	properties, err := json.Marshal(item.Properties)
	if err != nil {
		return r.errorHandler.HandleError(ctx, err)
	}

	metadata, err := json.Marshal(item.Metadata)
	if err != nil {
		return r.errorHandler.HandleError(ctx, err)
	}

	stmt, err := exec.PrepareContext(ctx, query)
	if err != nil {
		return r.errorHandler.HandleError(ctx, err)
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		item.Name,
		item.DisplayName,
		item.CreateTime,
		item.UpdateTime,
		item.DeleteTime,
		item.Hash,
		item.Content,
		properties,
		metadata,
	)

	if err != nil {
		return r.errorHandler.HandleError(ctx, err)
	}

	return nil
}

// GetItem retrieves an item from the database by its resource name.
func (r *ItemRepository) GetItem(ctx context.Context, name string) (*item.Item, error) {
	query := `
        SELECT name, display_name, create_time, update_time, delete_time, hash, content, properties, metadata
        FROM items
        WHERE name = $1 AND delete_time IS NULL
    `
	var item item.Item
	var propertiesJSON []byte
	var metadataJSON []byte

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, name).Scan(
		&item.Name,
		&item.DisplayName,
		&item.CreateTime,
		&item.UpdateTime,
		&item.DeleteTime,
		&item.Hash,
		&item.Content,
		&propertiesJSON,
		&metadataJSON,
	)

	if err != nil {
		return nil, r.errorHandler.HandleError(ctx, err)
	}

	err = json.Unmarshal(propertiesJSON, &item.Properties)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal properties: %w", err)
	}

	err = json.Unmarshal(metadataJSON, &item.Metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return &item, nil
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
