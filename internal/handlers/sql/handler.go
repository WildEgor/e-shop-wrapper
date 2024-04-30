package sql

import (
	dtos "github.com/WildEgor/e-shop-fiber-wrapper/internal/dtos/sql"
	"github.com/WildEgor/e-shop-fiber-wrapper/internal/repositories"
	"github.com/WildEgor/e-shop-fiber-wrapper/internal/validators"
	core_dtos "github.com/WildEgor/e-shop-gopack/pkg/core/dtos"
	"github.com/WildEgor/e-shop-gopack/pkg/libs/logger/models"
	"github.com/gofiber/fiber/v3"
	"log/slog"
)

type SQLHandler struct {
	rr *repositories.RecordsRepository
}

func NewSQLHandler(
	rr *repositories.RecordsRepository,
) *SQLHandler {
	return &SQLHandler{
		rr,
	}
}

// Sql godoc
//
//	@Param			body body	dtos.SQLRequestDto	true	"Body"
//	@Param			X-API-KEY header	string	true	"123"
//	@Summary		SQL Query
//	@Description	SQL Query
//	@Tags			SQL Controller
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dtos.SQLResponseDto
//	@Router			/api/v1/sql [post]
func (h *SQLHandler) Handle(ctx fiber.Ctx) error {
	dto := &dtos.SQLRequestDto{}
	if resp := validators.ParseAndValidate(ctx, dto); resp != nil {
		return resp.JSON()
	}

	resp := core_dtos.NewResponse(ctx)

	slog.Debug("handle sql", models.LogEntryAttr(&models.LogEntry{
		Props: map[string]interface{}{
			"sql": dto.Sql,
		},
	}))

	rows, err := h.rr.GetRecords(ctx.Context(), dto.Sql)
	if err != nil {
		resp.SetStatus(fiber.StatusInternalServerError)

		return nil
	}

	resp.SetStatus(fiber.StatusOK)
	resp.SetData(rows)

	return resp.JSON()
}
