package auth

import (
	"context"
	"database/sql"
)

// Repository provides an interface for token persistence operations.
type Repository interface {
	// CreateToken stores a new token in the database.
	CreateToken(ctx context.Context, tx *sql.Tx, token Token) error

	// GetToken retrieves a token from the database by its hash.
	GetToken(ctx context.Context, tx *sql.Tx, hash string) (*Token, error)

	// DeleteToken removes a token from the database by its ID.
	DeleteToken(ctx context.Context, tx *sql.Tx, id string) error
}
