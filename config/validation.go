package config

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// contract
type Validation interface {
	ValidateRequest(request interface{}) interface{}
}

func NewValidation() Validation {
	validate := validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &ValidationImpl{
		Validate: validate,
	}
}

func (validateImpl *ValidationImpl) ValidateRequest(request interface{}) interface{} {
	err := validateImpl.Validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)

		out := make([]ErrorMessage, len(validationErrors))
		for i, fieldError := range validationErrors {
			out[i] = ErrorMessage{
				Field:   fieldError.Field(),
				Message: GetErrorMsg(fieldError),
			}
		}
		return out
	}
	return nil
}

type ValidationImpl struct {
	Validate *validator.Validate
}
