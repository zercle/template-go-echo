package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Save original env
	originalPort := os.Getenv("PORT")
	originalEnv := os.Getenv("ENV")

	// Set test values
	//nolint:errcheck
	os.Setenv("PORT", "9000")
	//nolint:errcheck
	os.Setenv("ENV", "testing")

	defer func() {
		// Restore original env
		if originalPort != "" {
			//nolint:errcheck
			os.Setenv("PORT", originalPort)
		} else {
			//nolint:errcheck
			os.Unsetenv("PORT")
		}
		if originalEnv != "" {
			//nolint:errcheck
			os.Setenv("ENV", originalEnv)
		} else {
			//nolint:errcheck
			os.Unsetenv("ENV")
		}
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Server.Port != "9000" {
		t.Errorf("expected port 9000, got %s", cfg.Server.Port)
	}
	if cfg.Server.Env != "testing" {
		t.Errorf("expected env testing, got %s", cfg.Server.Env)
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				Server: ServerConfig{Port: "8080", Env: "development"},
				DB:     DatabaseConfig{Driver: "mysql", DSN: "dsn"},
				JWT:    JWTConfig{Secret: "secret"},
			},
			wantErr: false,
		},
		{
			name: "missing port",
			config: &Config{
				Server: ServerConfig{Port: "", Env: "development"},
				DB:     DatabaseConfig{Driver: "mysql", DSN: "dsn"},
				JWT:    JWTConfig{Secret: "secret"},
			},
			wantErr: true,
		},
		{
			name: "missing driver",
			config: &Config{
				Server: ServerConfig{Port: "8080", Env: "development"},
				DB:     DatabaseConfig{Driver: "", DSN: "dsn"},
				JWT:    JWTConfig{Secret: "secret"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigIsDevelopment(t *testing.T) {
	cfg := &Config{Server: ServerConfig{Env: "development"}}
	if !cfg.IsDevelopment() {
		t.Error("expected IsDevelopment to be true")
	}

	cfg.Server.Env = "production"
	if cfg.IsDevelopment() {
		t.Error("expected IsDevelopment to be false")
	}
}

func TestConfigIsProduction(t *testing.T) {
	cfg := &Config{Server: ServerConfig{Env: "production"}}
	if !cfg.IsProduction() {
		t.Error("expected IsProduction to be true")
	}

	cfg.Server.Env = "development"
	if cfg.IsProduction() {
		t.Error("expected IsProduction to be false")
	}
}
