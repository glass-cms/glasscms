package item

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// CreateItem creates a new item in the database.
func (r *Repository) CreateItem(ctx context.Context, tx *sql.Tx, item *Item) error {
	query := `
        INSERT INTO items (uid, create_time, update_time, hash, display_name, name, path, content, properties)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `

	propertiesJSON, err := json.Marshal(item.Properties)
	if err != nil {
		return fmt.Errorf("failed to marshal properties: %w", err)
	}

	// Execute the query
	_, err = tx.ExecContext(ctx, query,
		item.UID,
		item.CreateTime,
		item.UpdateTime,
		item.Hash,
		item.DisplayName,
		item.Name,
		item.Path,
		item.Content,
		propertiesJSON,
	)

	if err != nil {
		return fmt.Errorf("failed to insert item: %w", err)
	}

	return nil
}

// GetItem retrieves an item from the database by its UID.
func (r *Repository) GetItem(_ context.Context, _ *sql.Tx, _ string) (*Item, error) {
	// TODO: Implement this method.
	return nil, nil
}

// UpdateItem updates an existing item in the database.
func (r *Repository) UpdateItem(_ context.Context, _ *sql.Tx, _ *Item) error {
	// TODO: Implement this method.
	return nil
}
