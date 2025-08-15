package config

import (
	"fmt"
	"strconv"
)

type AppConfig struct {
	Server   ServerConfig
	Database DatabaseConfig
	App      ApplicationConfig
}

type ServerConfig struct {
	Host string
	Port int
}

type ApplicationConfig struct {
	Environment string
	LogLevel    string
	Debug       bool
}

func NewAppConfig() (*AppConfig, error) {
	port, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid PORT value: %w", err)
	}

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

func (c *AppConfig) ServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

func (c *AppConfig) IsProduction() bool {
	return c.App.Environment == "production"
}

func (c *AppConfig) IsDevelopment() bool {
	return c.App.Environment == "development"
}
