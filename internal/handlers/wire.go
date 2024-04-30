package handlers

import (
	eh "github.com/WildEgor/e-shop-fiber-wrapper/internal/handlers/errors"
	hch "github.com/WildEgor/e-shop-fiber-wrapper/internal/handlers/health_check"
	rch "github.com/WildEgor/e-shop-fiber-wrapper/internal/handlers/ready_check"
	"github.com/WildEgor/e-shop-fiber-wrapper/internal/handlers/sql"
	"github.com/WildEgor/e-shop-fiber-wrapper/internal/repositories"
	"github.com/WildEgor/e-shop-fiber-wrapper/internal/services"
	"github.com/google/wire"
)

var HandlersSet = wire.NewSet(
	repositories.RepositoriesSet,
	services.ServicesSet,
	eh.NewErrorsHandler,
	hch.NewHealthCheckHandler,
	rch.NewReadyCheckHandler,
	sql.NewSQLHandler,
)
