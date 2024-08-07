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
        INSERT INTO items (uid, create_time, update_time, hash, name, path, content, properties)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `

	propertiesJSON, err := json.Marshal(item.Properties)
	if err != nil {
		return fmt.Errorf("failed to marshal properties: %w", err)
	}

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		item.UID,
		item.CreateTime,
		item.UpdateTime,
		item.Hash,
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
func (r *Repository) GetItem(ctx context.Context, uid string) (*Item, error) {
	query := `
        SELECT uid, create_time, update_time, hash, name, path, content, properties
        FROM items
        WHERE uid = $1
    `
	var item Item
	var propertiesJSON []byte

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, uid).Scan(
		&item.UID,
		&item.CreateTime,
		&item.UpdateTime,
		&item.Hash,
		&item.Name,
		&item.Path,
		&item.Content,
		&propertiesJSON,
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

	return &item, nil
}

// UpdateItem updates an existing item in the database.
func (r *Repository) UpdateItem(ctx context.Context, item *Item) error {
	query := `
		UPDATE items
		SET update_time = $1, hash = $2, name = $3, path = $4, content = $5, properties = $6
		WHERE uid = $8
	`

	propertiesJSON, err := json.Marshal(item.Properties)
	if err != nil {
		return fmt.Errorf("failed to marshal properties: %w", err)
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
		item.Path,
		item.Content,
		propertiesJSON,
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
func (r *Repository) DeleteItem(ctx context.Context, uid string) error {
	query := `DELETE FROM items WHERE uid = $1`

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
	}

	return nil
}
