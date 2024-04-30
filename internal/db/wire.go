package db

import (
	"github.com/WildEgor/e-shop-fiber-wrapper/internal/db/clickhouse"
	"github.com/google/wire"
)

var DbSet = wire.NewSet(
	clickhouse.NewClickhouseConnection,
)
