package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"
)

// ErrDuplicatePrimaryKey is returned when an insert operation fails because the primary key already exists.
var ErrDuplicatePrimaryKey = errors.New("primary key constraint violated")

// ErrUniqueConstraint is an error that occurs when a unique constraint is violated.
var ErrUniqueConstraint = errors.New("unique constraint violated")

// ErrNotFound is an error that occurs when a statement does not return any rows.
var ErrNotFound = errors.New("not found")

// ErrOperationFailed is a fallback error for when an operation fails.
var ErrOperationFailed = errors.New("operation failed")

// ErrorHandler is an interface for handling database errors.
// Implementations of this interface should handle database-driver specific errors.
type ErrorHandler interface {
	HandleError(ctx context.Context, err error) error
}

// NewErrorHandler creates a new ErrorHandler based on the provided configuration.
func NewErrorHandler(cfg Config) (ErrorHandler, error) {
	driver, ok := DriverValue[cfg.Driver]
	if !ok {
		return nil, fmt.Errorf("unrecognized database driver: %s", cfg.Driver)
	}

	switch driver {
	case int32(DriverPostgres):
		return &PostgresErrorHandler{}, nil
	case int32(DriverSqlite):
		return NewSqliteErrorHandler(), nil
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}
}

type SqliteErrorHandler struct{}

func NewSqliteErrorHandler() *SqliteErrorHandler {
	return &SqliteErrorHandler{}
}

func (e *SqliteErrorHandler) HandleError(_ context.Context, err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("not found: %w", ErrNotFound)
	}

	// Handle driver-specific errors.
	var sqliteErr sqlite3.Error

	if errors.As(err, &sqliteErr) {
		// Check for constraint violations.
		if sqliteErr.Code == sqlite3.ErrConstraint {
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintPrimaryKey:
				return fmt.Errorf("%w : %w", ErrDuplicatePrimaryKey, err)
			case sqlite3.ErrConstraintUnique:
				return fmt.Errorf("%w : %w", ErrUniqueConstraint, err)
			}
		}
	}

	return fmt.Errorf("%w : %w", ErrOperationFailed, err)
}

type PostgresErrorHandler struct{}

func (e *PostgresErrorHandler) HandleError(_ context.Context, err error) error {
	// TODO: Implement error handling.
	return err
}
