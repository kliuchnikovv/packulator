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

type PackagingService struct {
	store            store.Store
	packagingService service.PackagingService
}

func NewPackagingService(store store.Store) *PackagingService {
	return &PackagingService{
		store:            store,
		packagingService: service.NewPackagingService(),
	}
}

func (c *PackagingService) Prefix() string {
	return "packaging"
}

func (c *PackagingService) Middlewares() []engi.Middleware {
	// Defines middlewares for all requests to this service.
	// CORS, auth, etc...
	return []engi.Middleware{
		cors.AllowedOrigins("*"),
		cors.AllowedHeaders("*"),
		cors.AllowedMethods("*"),
		auth.NoAuth(),
	}
}

func (c *PackagingService) Routers() engi.Routes {
	return engi.Routes{
		engi.GET("number_of_packages"): engi.Handle(
			c.NumberOfPackages,
			query.Integer("amount", validate.Greater(0)),
			query.String("packs_hash", validate.NotEmpty),
		),
	}
}

// NumberOfPackages
func (c *PackagingService) NumberOfPackages(
	ctx context.Context,
	request engi.Request,
	response engi.Response,
) error {
	var (
		amount      = request.Integer("amount", placing.InQuery)
		versionHash = request.String("packs_hash", placing.InQuery)
	)

	invariants, err := c.store.GetPacksInvariantsByHash(ctx, versionHash)
	switch {
	case err == nil:
	case errors.Is(err, store.ErrNotFound):
		return response.NotFound("packs not found by hash: %s", versionHash)
	default:
		return response.InternalServerError("failed to get packs: %s", err)
	}

	result, err := c.packagingService.NumberOfPacks(ctx, amount, invariants)
	if err != nil {
		return response.InternalServerError("can't calculate number of packages: %s", err)
	}

	return response.OK(result)
}
