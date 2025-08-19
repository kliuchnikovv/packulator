// Package config provides configuration management for the Packulator application.
// It loads configuration from environment variables with sensible defaults.
package config

import (
	"fmt"
	"strconv"
)

// AppConfig holds the complete application configuration
type AppConfig struct {
	Server   ServerConfig      // HTTP server configuration
	Database DatabaseConfig    // Database connection configuration
	App      ApplicationConfig // Application-specific settings
}

// ServerConfig contains HTTP server settings
type ServerConfig struct {
	Host string // Server host address
	Port int    // Server port number
}

// ApplicationConfig contains general application settings
type ApplicationConfig struct {
	Environment string // Application environment (development, production)
	LogLevel    string // Logging level (debug, info, warn, error)
	Debug       bool   // Debug mode flag
}

// NewAppConfig creates a new application configuration by loading values
// from environment variables with fallback to default values.
func NewAppConfig() (*AppConfig, error) {
	// Parse server port from environment variable
	port, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid PORT value: %w", err)
	}

	// Parse debug flag from environment variable
	debug, err := strconv.ParseBool(getEnv("DEBUG", "false"))
	if err != nil {
		return nil, fmt.Errorf("invalid DEBUG value: %w", err)
	}

	return &AppConfig{
		Server: ServerConfig{
			Host: getEnv("HOST", "0.0.0.0"),
			Port: port,
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Database: getEnv("DB_NAME", "packulator"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		App: ApplicationConfig{
			Environment: getEnv("ENVIRONMENT", "development"),
			LogLevel:    getEnv("LOG_LEVEL", "info"),
			Debug:       debug,
		},
	}, nil
}

// ServerAddress returns the formatted server address as host:port
func (c *AppConfig) ServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// IsProduction returns true if the application is running in production environment
func (c *AppConfig) IsProduction() bool {
	return c.App.Environment == "production"
}

// IsDevelopment returns true if the application is running in development environment
func (c *AppConfig) IsDevelopment() bool {
	return c.App.Environment == "development"
}
