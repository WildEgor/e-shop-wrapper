package configs

import (
	"github.com/caarlos0/env/v7"
	"log/slog"
)

type ApiKeyConfig struct {
	Key string `env:"API_KEY" envDefault:""`
}

func NewApiKeyConfig(c *Configurator) *ApiKeyConfig {
	cfg := ApiKeyConfig{}

	if err := env.Parse(&cfg); err != nil {
		slog.Error("error parse apikey config")
	}

	return &cfg
}
