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

type PacksAPI struct {
	packService service.PackService
}

func NewPacksAPI(store store.Store) *PacksAPI {
	return &PacksAPI{
		packService: service.NewPackService(store),
	}
}

func (c *PacksAPI) Prefix() string {
	return "packs"
}

func (c *PacksAPI) Middlewares() []engi.Middleware {
	return []engi.Middleware{

		cors.AllowedOrigins("*"),
		cors.AllowedHeaders("*"),
		cors.AllowedMethods("*"),
		auth.NoAuth(),
	}
}

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
