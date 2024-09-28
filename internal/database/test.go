package database

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// NewDB creates a new in-memory SQLite database with the necessary schema
// for testing purposes.
func NewTestDB() (*sql.DB, error) {
	config := Config{
		Driver: DriverName[int32(DriverSqlite)],
		DSN:    fmt.Sprintf("file:%s?mode=memory&cache=shared", uuid.New().String()),
	}

	db, err := NewConnection(config)
	if err != nil {
		return nil, err
	}

	if err = MigrateDatabase(db, config); err != nil {
		return nil, err
	}

	return db, nil
}
