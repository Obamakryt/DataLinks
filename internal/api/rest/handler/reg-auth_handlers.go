package handler

import (
	"DataLinks/internal/dto/request"
	"DataLinks/internal/service/auth-reg"
	auth "DataLinks/internal/storages/postgreSQL/reg_auth"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"log/slog"
)

type Handler struct {
	Validator *validator.Validate
	logger    *slog.Logger
	Storage   auth.GlobalStorage
}

func NewHandler(validator *validator.Validate, logger *slog.Logger, storage auth.GlobalStorage) *Handler {
	return &Handler{Validator: validator, logger: logger, Storage: storage}
}

func (h *Handler) RegHandler(e echo.Context) error {
	JsonStruct := request.Register{}

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
	RegService := auth_reg.LogicReg{DataWithoutHash: JsonStruct, Logger: h.logger}

	err = RegService.NewUser(e.Request().Context(), h.Storage)
	if err != nil {
		return Failed(e, BadCode, "during process happened error")
	}
	return Success(e, GoodCode, "New user created", "/login")
}

func (h *Handler) LogHandler(e echo.Context) error {
	JsonStruct := request.LogIn{}

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
	jwtsight := auth_reg.JWTSigh{}
	err = jwtsight.CreateSigh()
	if err != nil {
		return Failed(e, BadCode, "something go wrong")
	}
	LogService := auth_reg.AuthLogic{Data: JsonStruct, Logger: h.logger, JWTSigh: jwtsight}

	token, err := LogService.NewAuth(e.Request().Context(), h.Storage)
	if err != nil {
		return Failed(e, BadCode, "during process happened error")
	}
	return Success(e, BadCode, "you success log in", token)
}
