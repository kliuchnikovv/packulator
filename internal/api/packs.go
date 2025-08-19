package api

import (
	"context"

	"github.com/kliuchnikovv/engi"
	"github.com/kliuchnikovv/engi/definition/middlewares/auth"
	"github.com/kliuchnikovv/engi/definition/middlewares/cors"
	"github.com/kliuchnikovv/engi/definition/parameter"
	"github.com/kliuchnikovv/engi/definition/parameter/placing"
	"github.com/kliuchnikovv/engi/definition/parameter/query"
	"github.com/kliuchnikovv/engi/definition/validate"
	"github.com/kliuchnikovv/packulator/internal/model"
	"github.com/kliuchnikovv/packulator/internal/service"
	"github.com/kliuchnikovv/packulator/internal/store"
)

// PacksAPI provides endpoints for managing pack configurations.
type PacksAPI struct {
	packService service.PackService // Service layer for pack operations
}

// NewPacksAPI creates a new packs API instance with the given store.
func NewPacksAPI(store store.Store) *PacksAPI {
	return &PacksAPI{
		packService: service.NewPackService(store),
	}
}

// Prefix returns the URL prefix for all pack management endpoints.
func (c *PacksAPI) Prefix() string {
	return "packs"
}

// Middlewares returns the middleware stack for pack management endpoints.
// Allows all origins, headers, methods and requires no authentication.
func (c *PacksAPI) Middlewares() []engi.Middleware {
	return []engi.Middleware{
		cors.AllowedOrigins("*"),
		cors.AllowedHeaders("*"),
		cors.AllowedMethods("*"),
		auth.NoAuth(),
	}
}

// Routers defines the available pack management routes:
// POST /packs/create - Create new pack configuration
// GET /packs/list - List all available packs  
// GET /packs/id - Get specific pack by ID
// GET /packs/hash - Get packs by version hash
// DELETE /packs/delete - Delete pack configuration
func (c *PacksAPI) Routers() engi.Routes {
	return engi.Routes{
		engi.PST("create"): engi.Handle(
			c.CreatePacks,
			parameter.Body(new(model.CreatePacksRequest)),
		),
		engi.GET("list"): engi.Handle(c.ListPacks),
		engi.GET("id"): engi.Handle(
			c.GetPackByID,
			query.String("id", validate.NotEmpty),
		),
		engi.GET("hash"): engi.Handle(
			c.GetPackByHash,
			query.String("hash", validate.NotEmpty),
		),
		engi.DEL("delete"): engi.Handle(
			c.DeletePack,
			query.String("id", validate.NotEmpty),
		),
	}
}

// CreatePacks handles POST /packs/create requests.
// It creates a new pack configuration with the provided pack sizes.
func (c *PacksAPI) CreatePacks(
	ctx context.Context,
	request engi.Request,
	response engi.Response,
) error {
	var body = request.Body().(*model.CreatePacksRequest)

	if len(body.Packs) == 0 {
		return response.InternalServerError("packs can't be empty")
	}

	versionHash, err := c.packService.CreatePacks(ctx, body.Packs...)
	if err != nil {
		return response.InternalServerError("can't create packs: %s", err)
	}

	return response.OK(model.CreatePacksResponse{
		VersionHash: versionHash,
	})
}

// ListPacks handles GET /packs/list requests.
// It returns a list of all available pack configurations.
func (c *PacksAPI) ListPacks(
	ctx context.Context,
	_ engi.Request,
	response engi.Response,
) error {
	packs, err := c.packService.ListPacks(ctx)
	if err != nil {
		return response.InternalServerError("can't list packs: %s", err)
	}
	return response.OK(packs)
}

// GetPackByID handles GET /packs/id requests.
// It retrieves a specific pack configuration by its unique ID.
func (c *PacksAPI) GetPackByID(
	ctx context.Context,
	request engi.Request,
	response engi.Response,
) error {
	var id = request.String("id", placing.InQuery)

	pack, err := c.packService.GetPackByID(ctx, id)
	if err != nil {
		return response.InternalServerError("can't get pack by id - %s: %s", id, err)
	}

	return response.OK(pack)
}

// GetPackByHash handles GET /packs/hash requests.
// It retrieves pack configurations by their version hash.
func (c *PacksAPI) GetPackByHash(
	ctx context.Context,
	request engi.Request,
	response engi.Response,
) error {
	var hash = request.String("hash", placing.InQuery)

	pack, err := c.packService.GetPackByHash(ctx, hash)
	if err != nil {
		return response.InternalServerError("can't get pack by hash - %s: %s", hash, err)
	}

	return response.OK(pack)
}

// DeletePack handles DELETE /packs/delete requests.
// It removes a pack configuration by its unique ID.
func (c *PacksAPI) DeletePack(
	ctx context.Context,
	request engi.Request,
	response engi.Response,
) error {
	var id = request.String("id", placing.InQuery)

	if err := c.packService.DeletePack(ctx, id); err != nil {
		return response.InternalServerError("can't delete pack: %s", err)
	}

	return response.NoContent()
}
