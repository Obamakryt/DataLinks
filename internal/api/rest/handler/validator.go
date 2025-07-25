package handler

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"unicode"
)

func ErrorValid(err error) map[string]string {
	var v validator.ValidationErrors
	PackString := make(map[string]string)

	if errors.As(err, &v) {
		IncorrectPlaces := map[string]string{
			"required": "Empty field",
			"email":    "Incorrect email",
			"min=8":    "short pass",
			"max=24":   "long pass",
			"password": "password doesnt comply standard",
		}
		for _, e := range v {
			fail, ok := IncorrectPlaces[e.Tag()]
			if !ok {
				fail = "validator error"
			}
			PackString[e.Field()] = fail
		}
		return PackString
	}
	return PackString
}

func LoggerValidatorError(logger *slog.Logger, slogerr map[string]string) {
	for k, v := range slogerr {
		logger.Info("Validator Error", k, v)
	}
}

func PasswordValidator(v *validator.Validate) {
	_ = v.RegisterValidation("password", func(f validator.FieldLevel) bool {
		value := f.Field().String()
		var hasUpper, hasDigit bool
		for _, char := range value {
			if unicode.IsUpper(char) {
				hasUpper = true
			}
			if unicode.IsDigit(char) {
				hasDigit = true
			}
			if hasDigit && hasUpper {
				return true
			}
		}
		return false
	})
}
