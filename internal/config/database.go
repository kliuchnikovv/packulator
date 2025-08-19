package config

import (
	"fmt"
	"os"
)

// DatabaseConfig holds PostgreSQL database connection settings
type DatabaseConfig struct {
	Host     string // Database host address
	Port     string // Database port number
	User     string // Database username
	Password string // Database password
	Database string // Database name
	SSLMode  string // SSL connection mode
}

// NewDatabaseConfig creates a new database configuration with values
// loaded from environment variables and sensible defaults.
func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", ""),
		Database: getEnv("DB_NAME", "packulator"),
		SSLMode:  getEnv("DB_SSL_MODE", "disable"),
	}
}

// DSN returns the PostgreSQL data source name (connection string)
// formatted for use with GORM and the postgres driver.
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode)
}

// getEnv retrieves an environment variable value or returns a default value
// if the environment variable is not set or is empty.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
