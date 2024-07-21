package database

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

// MigrateDatabase migrates the database to the latest version.
func MigrateDatabase(db *sql.DB, cfg Config) error {
	if err := goose.SetDialect(cfg.Driver); err != nil {
		return err
	}
	goose.SetBaseFS(embedMigrations)

	return goose.Up(db, "migrations")
}
