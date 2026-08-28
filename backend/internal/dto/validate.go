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

	// valuetype checks that the field's raw JSON text decodes to the scalar shape named by
	// the sibling ValueType field (one of "string", "number", "boolean").
	validate.RegisterValidation("valuetype", func(fl validator.FieldLevel) bool {
		valueType := fl.Parent().FieldByName("ValueType").String()
		fieldBytes := []byte(fl.Field().String())

		switch valueType {
		case "string":
			var v string
			return json.Unmarshal(fieldBytes, &v) == nil
		case "number":
			var v float64
			return json.Unmarshal(fieldBytes, &v) == nil
		case "boolean":
			var v bool
			return json.Unmarshal(fieldBytes, &v) == nil
		case "price":
			var v PTPrice
			if json.Unmarshal(fieldBytes, &v) != nil {
				return false
			}
			return Validate(v) == nil
		}
		return false
	})
}

func Validate(s any) error {
	return validate.Struct(s)
}
