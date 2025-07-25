package handler

import (
	"DataLinks/internal/api/rest"
	"DataLinks/internal/service/auth"
	"github.com/go-playground/validator/v10"
	"log/slog"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	Validator *validator.Validate
	logger    *slog.Logger
	Storage   auth.Storage
}

func NewHandler(validator *validator.Validate, logger *slog.Logger, storage auth.Storage) *Handler {
	return &Handler{Validator: validator, logger: logger, Storage: storage}
}

func (h *Handler) RegHandler(e echo.Context) error {
	JsonStruct := rest.RequestRegister{}

	err := e.Bind(&JsonStruct)
	if err != nil {
		return Failed(e, BadCode, "invalid data")
	}
	PasswordValidator(h.Validator)
	err = h.Validator.Struct(JsonStruct)
	if err != nil {
		DataErr := ErrorValid(err)
		LoggerValidatorError(h.logger, DataErr)
		return Failed(e, BadCode, "invalid data")
	}
	RegService := auth.RegLogic{DataWithoutHash: JsonStruct, Logger: h.logger}

	err = RegService.NewUser(e.Request().Context(), h.Storage)
	if err != nil {
		return Failed(e, BadCode, "during process happened error")
	}
	return Success(e, GoodCode, "New user created", "/login")
}

func (h *Handler) LogHandler(e echo.Context) error {
	JsonStruct := rest.RequestLogIn{}

	err := e.Bind(&JsonStruct)
	if err != nil {
		return Failed(e, BadCode, "invalid data")
	}
	PasswordValidator(h.Validator)
	err = h.Validator.Struct(JsonStruct)
	if err != nil {
		DataErr := ErrorValid(err)
		LoggerValidatorError(h.logger, DataErr)
		return Failed(e, BadCode, "invalid data")
	}
	jwtsight := auth.JWTSigh{}
	err = jwtsight.CreateSigh()
	if err != nil {
		return Failed(e, BadCode, "something go wrong")
	}
	LogService := auth.AuthLogic{Data: JsonStruct, Logger: h.logger, JWTSigh: jwtsight}

	token, err := LogService.NewAuth(e.Request().Context(), h.Storage)
	if err != nil {
		return Failed(e, BadCode, "during process happened error")
	}
	return Success(e, BadCode, "you success log in", token)
}
