package validators

import (
	core_dtos "github.com/WildEgor/e-shop-gopack/pkg/core/dtos"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"strings"
)

// NewValidator func for create a new validator for model fields.
func NewValidator() *validator.Validate {
	// Create a new validator for a Book model.
	validate := validator.New()

	// Custom validation for sql query
	_ = validate.RegisterValidation("sql", func(fl validator.FieldLevel) bool {

		str := fl.Field().String()
		if len(str) == 0 {
			return false
		}

		if !strings.Contains(strings.ToLower(str), "limit") {
			return false
		}

		if strings.Contains(strings.ToLower(str), "insert") {
			return false
		}

		if strings.Contains(strings.ToLower(str), "drop") {
			return false
		}

		if strings.Contains(strings.ToLower(str), "truncate") {
			return false
		}

		return true
	})

	return validate
}

// ParseAndValidate parser
func ParseAndValidate(ctx fiber.Ctx, out interface{}) *core_dtos.ResponseDto {
	resp := core_dtos.NewResponse(ctx)

	// Checking received data from JSON body. Return status 400 and error message.
	if err := ctx.Bind().Body(&out); err != nil {
		resp.SetStatus(fiber.StatusBadRequest)
		return resp
	}

	// TODO
	// Create a new validator for a RegistrationRequestDto.
	// validate := NewValidator()

	// Validate fields.
	//if err := validate.Struct(&out); err != nil {
	//	resp.SetStatus(fiber.StatusBadRequest)
	//
	//	log.Errorf("Validation error: %s", err)
	//	// TODO: add validation logic
	//
	//	return resp
	//}

	return nil
}

// ValidatorErrors func for show validation errors for each invalid fields.
func validatorErrors(err error) map[string]string {
	// Define fields map.
	fields := map[string]string{}

	// TODO: add validation logic

	return fields
}
