package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zercle/template-go-echo/internal/config"
)

// Database represents the database connection
type Database struct {
	conn *sql.DB
}

// New creates a new database connection
func New(cfg *config.DatabaseConfig) (*Database, error) {
	db, err := sql.Open(cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxConns)
	db.SetMaxIdleConns(cfg.MaxConns / 2)
	db.SetConnMaxLifetime(time.Hour)
	db.SetConnMaxIdleTime(10 * time.Minute)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	slog.Info("database connection established",
		slog.String("driver", cfg.Driver),
		slog.Int("max_conns", cfg.MaxConns),
	)

	return &Database{conn: db}, nil
}

// GetConn returns the underlying database connection
func (d *Database) GetConn() *sql.DB {
	return d.conn
}

// Close closes the database connection
func (d *Database) Close() error {
	if d.conn != nil {
		return d.conn.Close()
	}
	return nil
}

// Health checks the database health
func (d *Database) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return d.conn.PingContext(ctx)
}

// BeginTx starts a new transaction
func (d *Database) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return d.conn.BeginTx(ctx, opts)
}
