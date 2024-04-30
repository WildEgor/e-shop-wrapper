package domain

import (
	core_dtos "github.com/WildEgor/e-shop-gopack/pkg/core/dtos"
	"github.com/gofiber/fiber/v3"
)

var ErrCodesMessages = map[int]string{
	99:  "unknown error",
	100: "wrong api key",
	101: "missing api key",
}

func SetWrongApiKeyError(resp *core_dtos.ResponseDto) {
	resp.SetStatus(fiber.StatusUnauthorized)
	resp.SetError(100, ErrCodesMessages[100])
}

func SetMissingApiKeyError(resp *core_dtos.ResponseDto) {
	resp.SetStatus(fiber.StatusBadRequest)
	resp.SetError(101, ErrCodesMessages[101])
}
