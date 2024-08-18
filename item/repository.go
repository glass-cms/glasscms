package item

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// CreateItem creates a new item in the database.
func (r *Repository) CreateItem(ctx context.Context, item *Item) error {
	query := `
        INSERT INTO items (
			uid, 
			name, 
			display_name, 
			create_time, 
			update_time, 
			delete_time, 
			hash, 
			content, 
			properties, 
			metadata)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `

	properties, err := json.Marshal(item.Properties)
	if err != nil {
		return fmt.Errorf("failed to marshal properties: %w", err)
	}

	metadata, err := json.Marshal(item.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		item.UID,
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
		return fmt.Errorf("failed to insert item: %w", err)
	}

	return nil
}

// GetItem retrieves an item from the database by its UID.
func (r *Repository) GetItem(ctx context.Context, uid string) (*Item, error) {
	query := `
        SELECT uid, name, display_name, create_time, update_time, delete_time, hash, content, properties, metadata
        FROM items
        WHERE uid = $1
    `
	var item Item
	var propertiesJSON []byte
	var metadataJSON []byte

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, uid).Scan(
		&item.UID,
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("item not found: %w", err)
		}
		return nil, fmt.Errorf("failed to retrieve item: %w", err)
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
func (r *Repository) UpdateItem(ctx context.Context, item *Item) error {
	query := `
		UPDATE items
		SET update_time = $1, hash = $2, name = $3, display_name = $4, content = $5, properties = $6, metadata = $7
		WHERE uid = $8
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
		item.UID,
	)

	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
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

// DeleteItem deletes an item from the database by its UID.
func (r *Repository) DeleteItem(_ context.Context, _ string) error {
	// TODO: Reimplement with soft delete.
	/*query := `DELETE FROM items WHERE uid = $1`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, uid)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("item not found")
	}*/

	return nil
}
