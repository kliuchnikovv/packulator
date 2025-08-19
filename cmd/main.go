// Package main is the entry point for the Packulator application.
// Packulator is a Go-based HTTP API backend that calculates the number
// of shipping packs needed for customer orders.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/kliuchnikovv/engi"
	"github.com/kliuchnikovv/engi/definition/response"
	"github.com/kliuchnikovv/packulator/internal/api"
	"github.com/kliuchnikovv/packulator/internal/config"
	"github.com/kliuchnikovv/packulator/internal/store"
	"go.opentelemetry.io/otel"
	"gorm.io/driver/postgres"
)

// main is the application entry point that initializes configuration,
// sets up logging, establishes database connection, starts the HTTP server,
// and handles graceful shutdown.
func main() {
	// Load application configuration from environment variables
	cfg, err := config.NewAppConfig()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Configure log level based on configuration
	var logLevel slog.Level
	switch cfg.App.LogLevel {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn", "warning":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	// Initialize structured logger
	var (
		logHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
		logger     = slog.New(logHandler)
	)

	logger.Info("starting packulator application",
		"environment", cfg.App.Environment,
		"address", cfg.ServerAddress(),
		"debug", cfg.App.Debug,
	)

	// Initialize database store with PostgreSQL connection
	store, err := store.NewStore(postgres.Open(cfg.Database.DSN()))
	if err != nil {
		logger.Error("failed to create store", "error", err)
		os.Exit(1)
	}

	// Create HTTP server engine with JSON responses, logging, and tracing
	var engine = engi.New(cfg.ServerAddress(),
		engi.ResponseAsJSON(response.AsIs),
		engi.WithLogger(slog.NewTextHandler(os.Stdout,
			&slog.HandlerOptions{
				Level: logLevel,
			},
		)),
		engi.WithTracerProvider(otel.GetTracerProvider()),
	)

	// Register API services: pack management, packaging calculations, and health checks
	if err := engine.RegisterServices(
		api.NewPacksAPI(store),
		api.NewPackagingService(store),
		api.NewHealthAPI(store),
	); err != nil {
		logger.Error("failed to register services", "error", err)
		os.Exit(1)
	}

	// Start HTTP server in a separate goroutine
	go func() {
		logger.Info("server starting", "address", cfg.ServerAddress())
		if err := engine.Start(); err != nil {
			logger.Error("failed to start engine", "error", err)
		}
	}()

	// Set up graceful shutdown handling for SIGINT and SIGTERM
	var intSignal = make(chan os.Signal, 1)
	signal.Notify(intSignal, syscall.SIGINT, syscall.SIGTERM)

	// Wait for interrupt signal
	<-intSignal

	// Gracefully shutdown the server
	logger.Info("received interruption signal: shutting down")
	engine.Shutdown(context.TODO())
}
