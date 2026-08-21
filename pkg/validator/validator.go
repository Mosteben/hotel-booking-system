package validator

import "github.com/go-playground/validator/v10"

var Validate = validator.New()

func ValidateStruct(data any) error {
	return Validate.Struct(data)
}