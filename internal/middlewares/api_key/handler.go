package api_jemiddlewares

import (
	"errors"
	"github.com/WildEgor/e-shop-fiber-wrapper/internal/domain"
	core_dtos "github.com/WildEgor/e-shop-gopack/pkg/core/dtos"
	"github.com/WildEgor/e-shop-gopack/pkg/libs/logger/models"
	"github.com/gofiber/fiber/v3"
	"log/slog"
	"net/url"
	"strings"
)

type contextKey int

const (
	tokenKey contextKey = 0
)

const (
	query  = "query"
	form   = "form"
	param  = "param"
	cookie = "cookie"
)

var (
	ErrMissingOrMalformedAPIKey = errors.New("missing or malformed API Key")
	ErrWrongAPIKey              = errors.New("wrong api key")
)

type ApiKeyMiddlewareConfig struct {
	Pass           bool
	Next           func(fiber.Ctx) bool
	SuccessHandler fiber.Handler
	ErrorHandler   fiber.ErrorHandler

	// KeyLookup is a string in the form of "<source>:<name>" that is used to extract key from the request.
	// Optional. Default value "header:X-API-KEY".
	// Possible values:
	// - "header:<name>"
	// - "query:<name>"
	// - "form:<name>"
	// - "param:<name>"
	// - "cookie:<name>"
	KeyLookup string

	// AuthScheme to be used in the Authorization header.
	// Optional. Default value "Bearer".
	AuthScheme string

	// Validator is a function to validate key.
	Validator func(fiber.Ctx, string) (bool, error)
}

var AuthMiddlewareConfigDefault = ApiKeyMiddlewareConfig{
	SuccessHandler: func(ctx fiber.Ctx) error {
		return ctx.Next()
	},

	ErrorHandler: func(ctx fiber.Ctx, err error) error {
		slog.Error("error", models.LogEntryAttr(&models.LogEntry{
			Err: err,
		}))

		resp := core_dtos.NewResponse(ctx)

		if errors.Is(err, ErrMissingOrMalformedAPIKey) {
			domain.SetMissingApiKeyError(resp)
		} else if errors.Is(err, ErrWrongAPIKey) {
			domain.SetWrongApiKeyError(resp)
		} else {
			resp.SetStatus(fiber.StatusInternalServerError)
		}

		return resp.JSON()
	},
}

func configDefault(config ...ApiKeyMiddlewareConfig) ApiKeyMiddlewareConfig {
	if len(config) < 1 {
		return AuthMiddlewareConfigDefault
	}

	cfg := config[0]

	if cfg.SuccessHandler == nil {
		cfg.SuccessHandler = AuthMiddlewareConfigDefault.SuccessHandler
	}

	if cfg.ErrorHandler == nil {
		cfg.ErrorHandler = AuthMiddlewareConfigDefault.ErrorHandler
	}

	if cfg.KeyLookup == "" {
		cfg.KeyLookup = AuthMiddlewareConfigDefault.KeyLookup
		// set AuthScheme as "Bearer" only if KeyLookup is set to default.
		if cfg.AuthScheme == "" {
			cfg.AuthScheme = AuthMiddlewareConfigDefault.AuthScheme
		}
	}

	if cfg.Validator == nil {
		panic("fiber: keyauth middleware requires a validator function")
	}

	return cfg
}

func NewApiKeyMiddleware(config ApiKeyMiddlewareConfig) fiber.Handler {

	cfg := configDefault(config)

	parts := strings.Split(cfg.KeyLookup, ":")

	slog.Debug("extract header", models.LogEntryAttr(&models.LogEntry{
		Props: map[string]interface{}{
			"headers": parts,
		},
	}))

	extractor := keyFromHeader(parts[1], cfg.AuthScheme)
	switch parts[0] {
	case query:
		extractor = keyFromQuery(parts[1])
	case form:
		extractor = keyFromForm(parts[1])
	case param:
		extractor = keyFromParam(parts[1])
	case cookie:
		extractor = keyFromCookie(parts[1])
	}

	// Return middleware handler
	return func(c fiber.Ctx) error {

		slog.Debug("handle api key middleware")

		// Filter request to skip middleware
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// Extract and verify key
		key, err := extractor(c)
		if err != nil {
			return cfg.ErrorHandler(c, err)
		}

		valid, err := cfg.Validator(c, key)

		if err == nil && valid {
			c.Locals(tokenKey, key)
			return cfg.SuccessHandler(c)
		}

		return cfg.ErrorHandler(c, err)
	}
}

// TokenFromContext returns the bearer token from the request context.
// returns an empty string if the token does not exist
func TokenFromContext(c fiber.Ctx) string {
	token, ok := c.Locals(tokenKey).(string)
	if !ok {
		return ""
	}

	return token
}

// keyFromHeader returns a function that extracts api key from the request header.
func keyFromHeader(header, authScheme string) func(c fiber.Ctx) (string, error) {
	return func(c fiber.Ctx) (string, error) {
		auth := c.Get(header)

		slog.Debug("api-key", models.LogEntryAttr(&models.LogEntry{
			Props: map[string]interface{}{
				"header": header,
				"auth":   auth,
			},
		}))

		l := len(authScheme)
		if len(auth) > 0 && l == 0 {
			return auth, nil
		}

		if len(auth) > l+1 && auth[:l] == authScheme {
			return auth[l+1:], nil
		}

		return "", ErrMissingOrMalformedAPIKey
	}
}

// keyFromQuery returns a function that extracts api key from the query string.
func keyFromQuery(param string) func(c fiber.Ctx) (string, error) {
	return func(c fiber.Ctx) (string, error) {
		key := c.Query(param)
		if key == "" {
			return "", ErrMissingOrMalformedAPIKey
		}
		return key, nil
	}
}

// keyFromForm returns a function that extracts api key from the form.
func keyFromForm(param string) func(c fiber.Ctx) (string, error) {
	return func(c fiber.Ctx) (string, error) {
		key := c.FormValue(param)
		if key == "" {
			return "", ErrMissingOrMalformedAPIKey
		}
		return key, nil
	}
}

// keyFromParam returns a function that extracts api key from the url param string.
func keyFromParam(param string) func(c fiber.Ctx) (string, error) {
	return func(c fiber.Ctx) (string, error) {
		key, err := url.PathUnescape(c.Params(param))
		if err != nil {
			return "", ErrMissingOrMalformedAPIKey
		}
		return key, nil
	}
}

// keyFromCookie returns a function that extracts api key from the named cookie.
func keyFromCookie(name string) func(c fiber.Ctx) (string, error) {
	return func(c fiber.Ctx) (string, error) {
		key := c.Cookies(name)
		if key == "" {
			return "", ErrMissingOrMalformedAPIKey
		}
		return key, nil
	}
}
