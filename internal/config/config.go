package config

import (
	"fmt"
)

// Config holds all application configuration.
type Config struct {
	Server ServerConfig
	DB     DatabaseConfig
	JWT    JWTConfig
	CORS   CORSConfig
	Log    LogConfig
}

// ServerConfig holds server configuration.
type ServerConfig struct {
	Port string
	Env  string
}

// DatabaseConfig holds database configuration.
type DatabaseConfig struct {
	Driver string
	DSN    string
}

// JWTConfig holds JWT configuration.
type JWTConfig struct {
	Secret string
}

// CORSConfig holds CORS configuration.
type CORSConfig struct {
	AllowedOrigins []string
}

// LogConfig holds logging configuration.
type LogConfig struct {
	Level string
}

// Load loads configuration from environment variables.
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		DB: DatabaseConfig{
			Driver: getEnv("DB_DRIVER", "mysql"),
			DSN:    getEnv("DB_DSN", "root@tcp(localhost:3306)/template_go_echo"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "dev-secret-key-change-in-production"),
		},
		CORS: CORSConfig{
			AllowedOrigins: parseStringSlice(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")),
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("PORT is required")
	}

	if c.DB.Driver == "" {
		return fmt.Errorf("DB_DRIVER is required")
	}

	if c.DB.DSN == "" {
		return fmt.Errorf("DB_DSN is required")
	}

	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}

	return nil
}

// IsDevelopment returns true if the environment is development.
func (c *Config) IsDevelopment() bool {
	return c.Server.Env == "development"
}

// IsProduction returns true if the environment is production.
func (c *Config) IsProduction() bool {
	return c.Server.Env == "production"
}
