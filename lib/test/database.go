package test

import (
	"database/sql"

	"github.com/glass-cms/glasscms/database"
)

// NewDB creates a new in-memory SQLite database with the necessary schema
// for testing purposes.
func NewDB() (*sql.DB, error) {
	config := database.Config{
		Driver: database.DriverName[int32(database.DriverSqlite)],
		DSN:    "file::memory:?cache=shared",
	}

	db, err := database.NewConnection(config)
	if err != nil {
		return nil, err
	}

	if err := database.MigrateDatabase(db, config); err != nil {
		return nil, err
	}

	return db, nil
}
