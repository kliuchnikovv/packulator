package api

import (
	"context"

	"github.com/kliuchnikovv/engi"
	"github.com/kliuchnikovv/engi/definition/middlewares/auth"
	"github.com/kliuchnikovv/engi/definition/middlewares/cors"
	"github.com/kliuchnikovv/packulator/internal/store"
)

type HealthAPI struct {
	store store.Store
}

func NewHealthAPI(store store.Store) *HealthAPI {
	return &HealthAPI{
		store: store,
	}
}

func (h *HealthAPI) Prefix() string {
	return "health"
}

func (h *HealthAPI) Middlewares() []engi.Middleware {
	return []engi.Middleware{
		cors.AllowedOrigins("*"),
		cors.AllowedHeaders("*"),
		cors.AllowedMethods("*"),
		auth.NoAuth(),
	}
}

func (h *HealthAPI) Routers() engi.Routes {
	return engi.Routes{
		engi.GET("check"): engi.Handle(h.HealthCheck),
	}
}

type HealthStatus struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	Version  string `json:"version"`
}

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