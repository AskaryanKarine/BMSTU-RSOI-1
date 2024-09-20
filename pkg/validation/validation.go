package validation

import (
	"github.com/go-playground/validator/v10"
)

type customValidator struct {
	validator *validator.Validate
}

func MustRegisterCustomValidator(v *validator.Validate) *customValidator {
	return &customValidator{validator: v}
}

func (cv *customValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
