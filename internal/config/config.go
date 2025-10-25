package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

// ServerConfig holds the server configuration
type ServerConfig struct {
	Address string
	Port    string
	Debug   bool
}

// DatabaseConfig holds the database configuration
type DatabaseConfig struct {
	Driver   string
	DSN      string
	MaxConns int
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret string
	TTL    int
}

// Load loads configuration from environment variables
func Load() *Config {
	// Set default values
	viper.SetDefault("SERVER_ADDRESS", ":8080")
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("SERVER_DEBUG", false)
	viper.SetDefault("DB_DRIVER", "mysql")
	viper.SetDefault("DB_DSN", "root:password@tcp(localhost:3306)/template_go_echo")
	viper.SetDefault("DB_MAX_CONNS", 10)
	viper.SetDefault("JWT_SECRET", "your-secret-key")
	viper.SetDefault("JWT_TTL", 3600)

	// Read environment variables
	viper.AutomaticEnv()

	// Load from .env file if it exists
	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		viper.SetConfigType("env")
		if err := viper.ReadInConfig(); err != nil {
			log.Printf("Error reading .env file: %v", err)
		}
	}

	cfg := &Config{
		Server: ServerConfig{
			Address: viper.GetString("SERVER_ADDRESS"),
			Port:    viper.GetString("SERVER_PORT"),
			Debug:   viper.GetBool("SERVER_DEBUG"),
		},
		Database: DatabaseConfig{
			Driver:   viper.GetString("DB_DRIVER"),
			DSN:      viper.GetString("DB_DSN"),
			MaxConns: viper.GetInt("DB_MAX_CONNS"),
		},
		JWT: JWTConfig{
			Secret: viper.GetString("JWT_SECRET"),
			TTL:    viper.GetInt("JWT_TTL"),
		},
	}

	cfg.Validate()
	return cfg
}

// Validate validates the configuration
func (c *Config) Validate() {
	if c.Server.Address == "" {
		log.Fatal("SERVER_ADDRESS is required")
	}
	if c.Database.Driver == "" {
		log.Fatal("DB_DRIVER is required")
	}
	if c.Database.DSN == "" {
		log.Fatal("DB_DSN is required")
	}
	if c.Database.MaxConns <= 0 {
		log.Fatal("DB_MAX_CONNS must be greater than 0")
	}
	if c.JWT.Secret == "" {
		log.Fatal("JWT_SECRET is required")
	}
	if c.JWT.Secret == "your-secret-key" {
		log.Printf("WARNING: Using default JWT_SECRET. Set JWT_SECRET environment variable for production")
	}
	if c.JWT.TTL <= 0 {
		log.Fatal("JWT_TTL must be greater than 0")
	}
}
