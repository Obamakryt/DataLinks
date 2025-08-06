package handler

import (
	"DataLinks/internal/slogger"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"log/slog"
)

func BindAndValidate[T any](RecStruct *T, e echo.Context, validator *validator.Validate, logger slogger.Setup) bool {
	if err := e.Bind(&RecStruct); err != nil {
		logger.Info("Bind error", slog.String("error", err.Error()))
		return true
	}
	err := validator.Struct(RecStruct)
	if err != nil {
		DataErr := ErrorValid(err)
		LoggerValidatorError(logger, DataErr)
		return false
	}
	return true
}
