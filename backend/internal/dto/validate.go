package dto

import (
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
}

func Validate(s any) error {
	return validate.Struct(s)
}
