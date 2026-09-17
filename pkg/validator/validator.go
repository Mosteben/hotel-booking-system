package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func init() {
	// Report validation errors using the request's JSON field name (e.g.
	// "first_name") instead of the Go struct field name (e.g. "FirstName"),
	// so FieldErrors below can hand the client field names it recognizes
	// from what it actually sent.
	Validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" || name == "" {
			return field.Name
		}
		return name
	})
}

func ValidateStruct(data any) error {
	return Validate.Struct(data)
}

// FieldErrors converts a validation failure from ValidateStruct into a
// client-safe map of {field: message}, e.g. {"first_name": "is required"}.
// It returns nil if err is not a validator.ValidationErrors, so callers can
// fall back to treating err as an ordinary (already client-safe) error.
func FieldErrors(err error) map[string]string {
	var verrs validator.ValidationErrors
	if !errors.As(err, &verrs) {
		return nil
	}

	fields := make(map[string]string, len(verrs))
	for _, fe := range verrs {
		fields[fe.Field()] = fieldMessage(fe)
	}
	return fields
}

func fieldMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return "must be at least " + fe.Param() + " characters"
	case "max":
		return "must be at most " + fe.Param() + " characters"
	case "oneof":
		return "must be one of: " + fe.Param()
	default:
		return fmt.Sprintf("is invalid (%s)", fe.Tag())
	}
}
