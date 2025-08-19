// Package api contains HTTP API handlers and routing for the Packulator application.
package api

import (
	"context"

	"github.com/kliuchnikovv/engi"
	"github.com/kliuchnikovv/engi/definition/middlewares/auth"
	"github.com/kliuchnikovv/engi/definition/middlewares/cors"
	"github.com/kliuchnikovv/packulator/internal/store"
)

// HealthAPI provides health check endpoints for monitoring service status.
type HealthAPI struct {
	store store.Store // Database store for connectivity checks
}

// NewHealthAPI creates a new health API instance with the given store.
func NewHealthAPI(store store.Store) *HealthAPI {
	return &HealthAPI{
		store: store,
	}
}

// Prefix returns the URL prefix for all health check endpoints.
func (h *HealthAPI) Prefix() string {
	return "health"
}

// Middlewares returns the middleware stack for health check endpoints.
// Allows all origins, headers, methods and requires no authentication.
func (h *HealthAPI) Middlewares() []engi.Middleware {
	return []engi.Middleware{
		cors.AllowedOrigins("*"),
		cors.AllowedHeaders("*"),
		cors.AllowedMethods("*"),
		auth.NoAuth(),
	}
}

// Routers defines the available health check routes.
// GET /health/check - Returns service health status
func (h *HealthAPI) Routers() engi.Routes {
	return engi.Routes{
		engi.GET("check"): engi.Handle(h.HealthCheck),
	}
}

// HealthStatus represents the response structure for health checks.
type HealthStatus struct {
	Status   string `json:"status"`   // Overall service status (ok, degraded)
	Database string `json:"database"` // Database connectivity status (ok, error)
	Version  string `json:"version"`  // Service version number
}

// HealthCheck handles GET /health/check requests.
// It checks database connectivity and returns the overall service health status.
func (h *HealthAPI) HealthCheck(
	ctx context.Context,
	_ engi.Request,
	response engi.Response,
) error {
	// Check database connectivity
	dbStatus := "ok"
	if err := h.store.HealthCheck(ctx); err != nil {
		dbStatus = "error"
	}

	status := HealthStatus{
		Status:   "ok",
		Database: dbStatus,
		Version:  "1.0.0",
	}

	if dbStatus != "ok" {
		status.Status = "degraded"
	}

	return response.OK(status)
}
