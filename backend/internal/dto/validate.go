package dto

import (
	"encoding/json"
	"log"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/studentinovisad/popisomator/backend/internal/config"
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

	// valuetype checks that the field's raw JSON text decodes to the shape named by the sibling
	// ValueType field, whether that's a scalar ("string", "number", "boolean", "expiry") or a
	// structured property type ("price", "mass", "volume").
	validate.RegisterValidation("valuetype", func(fl validator.FieldLevel) bool {
		valueType := fl.Parent().FieldByName("ValueType").String()
		fieldBytes := fl.Field().Bytes()

		if config.CurrentConfig.DebugMode {
			log.Printf("Property type validation; type %v; value %v", valueType, string(fieldBytes))
		}

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
		case "expiry":
			var v string
			if json.Unmarshal(fieldBytes, &v) != nil {
				return false
			}
			_, err := time.Parse(time.DateOnly, v)
			return err == nil
		case "mass":
			var mass PTMass
			if json.Unmarshal(fieldBytes, &mass) != nil {
				return false
			}
			return Validate(mass) == nil
		case "volume":
			var volume PTVolume
			if json.Unmarshal(fieldBytes, &volume) != nil {
				return false
			}
			return Validate(volume) == nil
		}
		return false
	})
}

func Validate(s any) error {
	return validate.Struct(s)
}
