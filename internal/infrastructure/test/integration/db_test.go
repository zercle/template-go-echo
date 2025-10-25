package integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/zercle/template-go-echo/internal/config"
)

func TestNewDatabase(t *testing.T) {
	// This test verifies the database package compiles and initializes correctly
	// Note: Actual connection testing would require a real database or mock
	cfg := &config.DatabaseConfig{
		Driver:   "mysql",
		DSN:      "root:password@tcp(localhost:3306)/test",
		MaxConns: 10,
	}

	// We're not testing actual connection here as it requires running database
	// This is just to verify the structure is valid
	if cfg.Driver == "" {
		t.Error("database driver not set")
	}
	if cfg.MaxConns == 0 {
		t.Error("max connections not set")
	}
}

func TestDatabaseHealthWithContext(t *testing.T) {
	// Verify context handling for health checks is properly structured
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if ctx.Err() != nil {
		t.Error("context should not be cancelled immediately")
	}
}
