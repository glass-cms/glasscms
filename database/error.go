package database

import (
	"context"
	"fmt"
)

// UniqueConstraintError is an error that occurs when a unique constraint is violated.
type UniqueConstraintError struct {
	Field string
	Err   error
}

func (e *UniqueConstraintError) Error() string {
	return e.Err.Error()
}

// ErrorHandler is an interface for handling database errors.
// Implementations of this interface should handle database-specific errors.
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
		return &SqliteErrorHandler{}, nil
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}
}

type SqliteErrorHandler struct{}

func (e *SqliteErrorHandler) HandleError(_ context.Context, err error) error {
	// TODO: Implement error handling.
	return err
}

type PostgresErrorHandler struct{}

func (e *PostgresErrorHandler) HandleError(_ context.Context, err error) error {
	// TODO: Implement error handling.
	return err
}
