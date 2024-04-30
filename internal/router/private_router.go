package router

import (
	"github.com/WildEgor/e-shop-fiber-wrapper/internal/handlers/sql"
	api_jemiddlewares "github.com/WildEgor/e-shop-fiber-wrapper/internal/middlewares/api_key"
	"github.com/WildEgor/e-shop-fiber-wrapper/internal/services"
	"github.com/gofiber/fiber/v3"
)

type PrivateRouter struct {
	vs *services.ApiKeyValidator
	sh *sql.SQLHandler
}

func NewPrivateRouter(
	vs *services.ApiKeyValidator,
	sh *sql.SQLHandler,
) *PrivateRouter {
	return &PrivateRouter{
		vs,
		sh,
	}
}

func (r *PrivateRouter) Setup(app *fiber.App) {
	v1 := app.Group("/api/v1")

	akm := api_jemiddlewares.NewApiKeyMiddleware(api_jemiddlewares.ApiKeyMiddlewareConfig{
		KeyLookup: "header:x-api-key",
		Validator: func(ctx fiber.Ctx, key string) (bool, error) {
			err := r.vs.Validate(key)
			if err != nil {
				return false, err
			}

			return true, nil
		},
	})

	v1.Post("sql", akm, r.sh.Handle)
}
