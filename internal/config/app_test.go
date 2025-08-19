package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAppConfig(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		// Clear environment variables
		envVars := []string{
			"HOST", "PORT", "DB_HOST", "DB_PORT", "DB_USER",
			"DB_PASSWORD", "DB_NAME", "DB_SSL_MODE", "ENVIRONMENT",
			"LOG_LEVEL", "DEBUG",
		}
		for _, env := range envVars {
			os.Unsetenv(env)
		}

		cfg, err := NewAppConfig()
		require.NoError(t, err)

		// Server defaults
		assert.Equal(t, "0.0.0.0", cfg.Server.Host)
		assert.Equal(t, 8080, cfg.Server.Port)

		// Database defaults
		assert.Equal(t, "localhost", cfg.Database.Host)
		assert.Equal(t, "5432", cfg.Database.Port)
		assert.Equal(t, "postgres", cfg.Database.User)
		assert.Equal(t, "", cfg.Database.Password)
		assert.Equal(t, "packulator", cfg.Database.Database)
		assert.Equal(t, "disable", cfg.Database.SSLMode)

		// App defaults
		assert.Equal(t, "development", cfg.App.Environment)
		assert.Equal(t, "info", cfg.App.LogLevel)
		assert.False(t, cfg.App.Debug)
	})

	t.Run("custom environment variables", func(t *testing.T) {
		// Set custom environment variables
		os.Setenv("HOST", "127.0.0.1")
		os.Setenv("PORT", "3000")
		os.Setenv("DB_HOST", "db.example.com")
		os.Setenv("DB_PORT", "5433")
		os.Setenv("DB_USER", "testuser")
		os.Setenv("DB_PASSWORD", "testpass")
		os.Setenv("DB_NAME", "testdb")
		os.Setenv("DB_SSL_MODE", "require")
		os.Setenv("ENVIRONMENT", "production")
		os.Setenv("LOG_LEVEL", "debug")
		os.Setenv("DEBUG", "true")

		defer func() {
			envVars := []string{
				"HOST", "PORT", "DB_HOST", "DB_PORT", "DB_USER",
				"DB_PASSWORD", "DB_NAME", "DB_SSL_MODE", "ENVIRONMENT",
				"LOG_LEVEL", "DEBUG",
			}
			for _, env := range envVars {
				os.Unsetenv(env)
			}
		}()

		cfg, err := NewAppConfig()
		require.NoError(t, err)

		// Server custom values
		assert.Equal(t, "127.0.0.1", cfg.Server.Host)
		assert.Equal(t, 3000, cfg.Server.Port)

		// Database custom values
		assert.Equal(t, "db.example.com", cfg.Database.Host)
		assert.Equal(t, "5433", cfg.Database.Port)
		assert.Equal(t, "testuser", cfg.Database.User)
		assert.Equal(t, "testpass", cfg.Database.Password)
		assert.Equal(t, "testdb", cfg.Database.Database)
		assert.Equal(t, "require", cfg.Database.SSLMode)

		// App custom values
		assert.Equal(t, "production", cfg.App.Environment)
		assert.Equal(t, "debug", cfg.App.LogLevel)
		assert.True(t, cfg.App.Debug)
	})

	t.Run("invalid PORT value", func(t *testing.T) {
		os.Setenv("PORT", "invalid")
		defer os.Unsetenv("PORT")

		cfg, err := NewAppConfig()
		assert.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "invalid PORT value")
	})

	t.Run("invalid DEBUG value", func(t *testing.T) {
		os.Setenv("DEBUG", "invalid")
		defer os.Unsetenv("DEBUG")

		cfg, err := NewAppConfig()
		assert.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "invalid DEBUG value")
	})
}

func TestAppConfig_ServerAddress(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		port     int
		expected string
	}{
		{
			name:     "default values",
			host:     "0.0.0.0",
			port:     8080,
			expected: "0.0.0.0:8080",
		},
		{
			name:     "localhost with custom port",
			host:     "localhost",
			port:     3000,
			expected: "localhost:3000",
		},
		{
			name:     "IPv6 localhost",
			host:     "::1",
			port:     8080,
			expected: "::1:8080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &AppConfig{
				Server: ServerConfig{
					Host: tt.host,
					Port: tt.port,
				},
			}

			result := cfg.ServerAddress()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAppConfig_IsProduction(t *testing.T) {
	tests := []struct {
		name        string
		environment string
		expected    bool
	}{
		{
			name:        "production environment",
			environment: "production",
			expected:    true,
		},
		{
			name:        "development environment",
			environment: "development",
			expected:    false,
		},
		{
			name:        "staging environment",
			environment: "staging",
			expected:    false,
		},
		{
			name:        "empty environment",
			environment: "",
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &AppConfig{
				App: ApplicationConfig{
					Environment: tt.environment,
				},
			}

			result := cfg.IsProduction()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAppConfig_IsDevelopment(t *testing.T) {
	tests := []struct {
		name        string
		environment string
		expected    bool
	}{
		{
			name:        "development environment",
			environment: "development",
			expected:    true,
		},
		{
			name:        "production environment",
			environment: "production",
			expected:    false,
		},
		{
			name:        "staging environment",
			environment: "staging",
			expected:    false,
		},
		{
			name:        "empty environment",
			environment: "",
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &AppConfig{
				App: ApplicationConfig{
					Environment: tt.environment,
				},
			}

			result := cfg.IsDevelopment()
			assert.Equal(t, tt.expected, result)
		})
	}
}
