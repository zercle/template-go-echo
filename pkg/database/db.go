package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// Config holds database configuration.
type Config struct {
	Driver      string
	DSN         string
	MaxOpenConn int
	MaxIdleConn int
	MaxLifetime time.Duration
}

// Open opens a database connection with the given configuration.
func Open(cfg Config) (*sql.DB, error) {
	if cfg.MaxOpenConn == 0 {
		cfg.MaxOpenConn = 25
	}
	if cfg.MaxIdleConn == 0 {
		cfg.MaxIdleConn = 5
	}
	if cfg.MaxLifetime == 0 {
		cfg.MaxLifetime = 5 * time.Minute
	}

	db, err := sql.Open(cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(cfg.MaxOpenConn)
	db.SetMaxIdleConns(cfg.MaxIdleConn)
	db.SetConnMaxLifetime(cfg.MaxLifetime)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// Close closes the database connection.
func Close(db *sql.DB) error {
	if db == nil {
		return nil
	}
	return db.Close()
}
