package api

import (
	"context"
	"errors"

	"github.com/kliuchnikovv/engi"
	"github.com/kliuchnikovv/engi/definition/middlewares/auth"
	"github.com/kliuchnikovv/engi/definition/middlewares/cors"
	"github.com/kliuchnikovv/engi/definition/parameter/placing"
	"github.com/kliuchnikovv/engi/definition/parameter/query"
	"github.com/kliuchnikovv/engi/definition/validate"
	"github.com/kliuchnikovv/packulator/internal/service"
	"github.com/kliuchnikovv/packulator/internal/store"
)

// PackagingService provides endpoints for pack calculation operations.
type PackagingService struct {
	store store.Store // Database store for pack retrieval
}

// NewPackagingService creates a new packaging service instance with the given store.
func NewPackagingService(store store.Store) *PackagingService {
	return &PackagingService{
		store: store,
	}
}

// Prefix returns the URL prefix for all packaging calculation endpoints.
func (c *PackagingService) Prefix() string {
	return "packaging"
}

// Middlewares returns the middleware stack for packaging calculation endpoints.
// Allows all origins, headers, methods and requires no authentication.
func (c *PackagingService) Middlewares() []engi.Middleware {
	return []engi.Middleware{
		cors.AllowedOrigins("*"),
		cors.AllowedHeaders("*"),
		cors.AllowedMethods("*"),
		auth.NoAuth(),
	}
}

// Routers defines the available packaging calculation routes:
// GET /packaging/number_of_packages - Calculate optimal pack combination for given amount
func (c *PackagingService) Routers() engi.Routes {
	return engi.Routes{
		engi.GET("number_of_packages"): engi.Handle(
			c.NumberOfPackages,
			query.Integer("amount", validate.Greater(0)),  // Required: amount > 0
			query.String("packs_hash", validate.NotEmpty), // Required: pack configuration hash
		),
	}
}

// NumberOfPackages handles GET /packaging/number_of_packages requests.
// It calculates the optimal combination of packs needed for a given amount
// using the pack configuration identified by the provided hash.
func (c *PackagingService) NumberOfPackages(
	ctx context.Context,
	request engi.Request,
	response engi.Response,
) error {
	// Extract query parameters
	var (
		amount      = request.Integer("amount", placing.InQuery)    // Amount to be packed
		versionHash = request.String("packs_hash", placing.InQuery) // Pack configuration hash
	)

	// Retrieve pack configuration by hash
	pack, err := c.store.GetPackByHash(ctx, versionHash)
	switch {
	case err == nil:
		// Pack found successfully
	case errors.Is(err, store.ErrNotFound):
		return response.NotFound("packs not found by hash: %s", versionHash)
	default:
		return response.InternalServerError("failed to get packs: %s", err)
	}

	// Calculate optimal pack combination
	result, err := service.NumberOfPacks(ctx, amount, pack.GetPacks())
	if err != nil {
		return response.InternalServerError("can't calculate number of packages: %s", err)
	}

	return response.OK(result)
}
