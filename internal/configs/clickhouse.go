package configs

import (
	"github.com/caarlos0/env/v7"
	"log/slog"
)

type ClickhouseConfig struct {
	DSN string `env:"CLICKHOUSE_DSN,required"`
	DB  string `env:"CLICKHOUSE_DB,required"`
}

func NewClickhouseConfig(c *Configurator) *ClickhouseConfig {
	cfg := ClickhouseConfig{}

	if err := env.Parse(&cfg); err != nil {
		slog.Error("error parse clickhouse config")
	}

	return &cfg
}
