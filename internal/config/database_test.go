package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDatabaseConfig(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		// Clear database environment variables
		envVars := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSL_MODE"}
		for _, env := range envVars {
			os.Unsetenv(env)
		}

		cfg := NewDatabaseConfig()

		assert.Equal(t, "localhost", cfg.Host)
		assert.Equal(t, "5432", cfg.Port)
		assert.Equal(t, "postgres", cfg.User)
		assert.Equal(t, "", cfg.Password)
		assert.Equal(t, "packulator", cfg.Database)
		assert.Equal(t, "disable", cfg.SSLMode)
	})

	t.Run("custom environment variables", func(t *testing.T) {
		// Set custom environment variables
		os.Setenv("DB_HOST", "db.example.com")
		os.Setenv("DB_PORT", "5433")
		os.Setenv("DB_USER", "testuser")
		os.Setenv("DB_PASSWORD", "testpass")
		os.Setenv("DB_NAME", "testdb")
		os.Setenv("DB_SSL_MODE", "require")

		defer func() {
			envVars := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSL_MODE"}
			for _, env := range envVars {
				os.Unsetenv(env)
			}
		}()

		cfg := NewDatabaseConfig()

		assert.Equal(t, "db.example.com", cfg.Host)
		assert.Equal(t, "5433", cfg.Port)
		assert.Equal(t, "testuser", cfg.User)
		assert.Equal(t, "testpass", cfg.Password)
		assert.Equal(t, "testdb", cfg.Database)
		assert.Equal(t, "require", cfg.SSLMode)
	})

	t.Run("partial custom environment variables", func(t *testing.T) {
		// Set only some environment variables
		os.Setenv("DB_HOST", "custom.host.com")
		os.Setenv("DB_PASSWORD", "secretpass")

		defer func() {
			os.Unsetenv("DB_HOST")
			os.Unsetenv("DB_PASSWORD")
		}()

		cfg := NewDatabaseConfig()

		// Custom values
		assert.Equal(t, "custom.host.com", cfg.Host)
		assert.Equal(t, "secretpass", cfg.Password)

		// Default values for unset variables
		assert.Equal(t, "5432", cfg.Port)
		assert.Equal(t, "postgres", cfg.User)
		assert.Equal(t, "packulator", cfg.Database)
		assert.Equal(t, "disable", cfg.SSLMode)
	})
}

func TestDatabaseConfig_DSN(t *testing.T) {
	tests := []struct {
		name     string
		config   DatabaseConfig
		expected string
	}{
		{
			name: "full configuration",
			config: DatabaseConfig{
				Host:     "localhost",
				Port:     "5432",
				User:     "postgres",
				Password: "secret",
				Database: "testdb",
				SSLMode:  "disable",
			},
			expected: "host=localhost port=5432 user=postgres password=secret dbname=testdb sslmode=disable",
		},
		{
			name: "empty password",
			config: DatabaseConfig{
				Host:     "localhost",
				Port:     "5432",
				User:     "postgres",
				Password: "",
				Database: "testdb",
				SSLMode:  "disable",
			},
			expected: "host=localhost port=5432 user=postgres password= dbname=testdb sslmode=disable",
		},
		{
			name: "production-like configuration",
			config: DatabaseConfig{
				Host:     "prod-db.example.com",
				Port:     "5432",
				User:     "app_user",
				Password: "complex_password_123",
				Database: "production_db",
				SSLMode:  "require",
			},
			expected: "host=prod-db.example.com port=5432 user=app_user password=complex_password_123 dbname=production_db sslmode=require",
		},
		{
			name: "custom port",
			config: DatabaseConfig{
				Host:     "localhost",
				Port:     "5433",
				User:     "postgres",
				Password: "password",
				Database: "packulator",
				SSLMode:  "disable",
			},
			expected: "host=localhost port=5433 user=postgres password=password dbname=packulator sslmode=disable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.DSN()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetEnv(t *testing.T) {
	t.Run("existing environment variable", func(t *testing.T) {
		os.Setenv("TEST_VAR", "test_value")
		defer os.Unsetenv("TEST_VAR")

		result := getEnv("TEST_VAR", "default_value")
		assert.Equal(t, "test_value", result)
	})

	t.Run("non-existing environment variable", func(t *testing.T) {
		os.Unsetenv("NON_EXISTING_VAR")

		result := getEnv("NON_EXISTING_VAR", "default_value")
		assert.Equal(t, "default_value", result)
	})

	t.Run("empty environment variable", func(t *testing.T) {
		os.Setenv("EMPTY_VAR", "")
		defer os.Unsetenv("EMPTY_VAR")

		result := getEnv("EMPTY_VAR", "default_value")
		assert.Equal(t, "default_value", result)
	})

	t.Run("whitespace environment variable", func(t *testing.T) {
		os.Setenv("WHITESPACE_VAR", "   ")
		defer os.Unsetenv("WHITESPACE_VAR")

		result := getEnv("WHITESPACE_VAR", "default_value")
		assert.Equal(t, "   ", result) // getEnv doesn't trim whitespace
	})
}
