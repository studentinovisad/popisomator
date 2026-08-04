package dto

import (
	"encoding/json"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

var (
	hasUpper = regexp.MustCompile(`[A-Z]`)
	hasLower = regexp.MustCompile(`[a-z]`)
	hasDigit = regexp.MustCompile(`[0-9]`)
)

func init() {
	validate.RegisterValidation("password_complexity", func(fl validator.FieldLevel) bool {
		pw := fl.Field().String()
		return hasUpper.MatchString(pw) && hasLower.MatchString(pw) && hasDigit.MatchString(pw)
	})

	// valuetype checks that the field's raw JSON text decodes to the shape named by the
	// sibling ValueType field (one of "string", "number", "boolean", "object", "array").
	validate.RegisterValidation("valuetype", func(fl validator.FieldLevel) bool {
		valueType := fl.Parent().FieldByName("ValueType").String()

		var v any
		if err := json.Unmarshal([]byte(fl.Field().String()), &v); err != nil {
			return false
		}

		switch valueType {
		case "string":
			_, ok := v.(string)
			return ok
		case "number":
			_, ok := v.(float64)
			return ok
		case "boolean":
			_, ok := v.(bool)
			return ok
		case "object":
			_, ok := v.(map[string]any)
			return ok
		case "array":
			_, ok := v.([]any)
			return ok
		}
		return false
	})
}

func Validate(s any) error {
	return validate.Struct(s)
}
