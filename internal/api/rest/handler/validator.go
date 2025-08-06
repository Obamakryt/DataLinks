package handler

import (
	"DataLinks/internal/slogger"
	"errors"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/url"
	"strings"
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
			"http":     "dont support http website",
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

func LoggerValidatorError(logger slogger.Setup, logErrs map[string]string) {
	for k, v := range logErrs {
		logger.Info("Validator Error", slog.String(k, v))
	}
}

func PasswordValidator(v *validator.Validate) {
	_ = v.RegisterValidation("password", func(f validator.FieldLevel) bool {
		value := f.Field().String()
		if err := v.Var(value, "min=8,max=24"); err != nil {
			return false
		}
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

func UrlValidator(v *validator.Validate) {
	_ = v.RegisterValidation("url", func(f validator.FieldLevel) bool {
		value := f.Field().String()
		if len(value) == 0 {
			return false
		}
		parsed, err := url.ParseRequestURI(value)
		if err != nil {
			return false
		}
		return parsed.Scheme != "" && parsed.Host != ""
	})
}

func HTTPValidator(v *validator.Validate) {
	_ = v.RegisterValidation("http", func(f validator.FieldLevel) bool {
		value := f.Field().String()
		if strings.HasPrefix(value, "http://") {
			return false
		}
		return true
	})
}
