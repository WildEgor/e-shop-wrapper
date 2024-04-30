package services

import (
	"github.com/WildEgor/e-shop-fiber-wrapper/internal/configs"
	api_jemiddlewares "github.com/WildEgor/e-shop-fiber-wrapper/internal/middlewares/api_key"
	"github.com/WildEgor/e-shop-gopack/pkg/libs/logger/models"
	"log/slog"
	"strings"
)

type ApiKeyValidator struct {
	cfg *configs.ApiKeyConfig
}

func NewApiKeyValidator(
	cfg *configs.ApiKeyConfig,
) *ApiKeyValidator {
	return &ApiKeyValidator{
		cfg,
	}
}

func (v *ApiKeyValidator) Validate(key string) error {

	slog.Debug("handle api key", models.LogEntryAttr(&models.LogEntry{
		Props: map[string]interface{}{
			"key": key,
		},
	}))

	if !strings.EqualFold(v.cfg.Key, key) {
		return api_jemiddlewares.ErrWrongAPIKey
	}

	return nil
}
