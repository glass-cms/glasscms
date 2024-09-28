package test

import (
	"database/sql"
	"fmt"

	"github.com/glass-cms/glasscms/database"
	"github.com/google/uuid"
)

// NewDB creates a new in-memory SQLite database with the necessary schema
// for testing purposes.
func NewDB() (*sql.DB, error) {
	config := database.Config{
		Driver: database.DriverName[int32(database.DriverSqlite)],
		DSN:    fmt.Sprintf("file:%s?mode=memory&cache=shared", uuid.New().String()),
	}

	db, err := database.NewConnection(config)
	if err != nil {
		return nil, err
	}

	if err = database.MigrateDatabase(db, config); err != nil {
		return nil, err
	}

	return db, nil
}
