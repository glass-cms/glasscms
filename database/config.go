package database

import (
	"database/sql"
	"fmt"
)

type Driver int32

const (
	DriverUnrecognized Driver = -1
	DriverUnspecified  Driver = iota
	DriverPostgres
	DriverSqlite

	MaxConnectionsDefault     = 5
	MaxIdleConnectionsDefault = 1
)

var (
	DriverName = map[int32]string{
		int32(DriverUnspecified): "unspecified",
		int32(DriverPostgres):    "postgres",
		int32(DriverSqlite):      "sqlite3",
	}
	DriverValue = map[string]int32{
		"unspecified": int32(DriverUnspecified),
		"postgres":    int32(DriverPostgres),
		"sqlite3":     int32(DriverSqlite),
	}
)

// Config represents the configuration for a database connection.
type Config struct {
	// Driver is the name of the database driver.
	Driver string

	// DSN is the Data Source Name. It specifies the username, password, and database name
	// that are used to connect to the database.
	DSN string

	// MaxConnection is the maximum number of connections that can be opened to the database.
	MaxConnections int

	// MaxIdleConnections is the maximum number of idle connections that can be maintained.
	// Idle connections are connections that are open but not in use.
	MaxIdleConnections int
}

// NewConnection creates a new database connection using the provided configuration.
// It returns a pointer to the sql.DB object and an error, if any occurred during the connection process.
// The sql.DB object represents a pool of zero or more underlying connections.
// It's safe for concurrent use by multiple goroutines.
func NewConnection(cfg Config) (*sql.DB, error) {
	db, err := sql.Open(cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if cfg.MaxConnections <= 0 {
		db.SetMaxOpenConns(MaxConnectionsDefault)
	} else {
		db.SetMaxOpenConns(cfg.MaxConnections)
	}

	if cfg.MaxIdleConnections <= 0 {
		db.SetMaxIdleConns(MaxIdleConnectionsDefault)
	} else {
		db.SetMaxIdleConns(cfg.MaxIdleConnections)
	}

	return db, err
}
