package handler

import (
	"DataLinks/internal/dto/request"
	"DataLinks/internal/dto/responce"
	"DataLinks/internal/service"
	"DataLinks/internal/service/jwt_hash"
	auth "DataLinks/internal/storages/postgreSQL/storage_crud"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"log/slog"
)

type Handler struct {
	Validator *validator.Validate
	logger    *slog.Logger
	Storage   auth.AuthReg
}

func NewHandler(validator *validator.Validate, logger *slog.Logger, storage auth.AuthReg) *Handler {
	return &Handler{Validator: validator, logger: logger, Storage: storage}
}

func (h *Handler) RegHandler(e echo.Context) error {
	JsonStruct := request.Register{}

	err := e.Bind(&JsonStruct)
	if err != nil {
		return responce.Failed(e, responce.BadCode, "invalid data")
	}
	PasswordValidator(h.Validator)
	err = h.Validator.Struct(JsonStruct)
	if err != nil {
		DataErr := ErrorValid(err)
		LoggerValidatorError(h.logger, DataErr)
		return responce.Failed(e, responce.BadCode, "invalid data")
	}
	RegService := service.LogicReg{DataWithoutHash: JsonStruct, Logger: h.logger}

	err = RegService.NewUser(e.Request().Context(), h.Storage)
	if err != nil {
		return responce.Failed(e, responce.BadCode, "during process happened error")
	}
	return responce.Success(e, responce.GoodCode, "New user created", "/login")
}

func (h *Handler) LogHandler(e echo.Context) error {
	JsonStruct := request.LogIn{}

	err := e.Bind(&JsonStruct)
	if err != nil {
		return responce.Failed(e, responce.BadCode, "invalid data")
	}
	PasswordValidator(h.Validator)
	err = h.Validator.Struct(JsonStruct)
	if err != nil {
		DataErr := ErrorValid(err)
		LoggerValidatorError(h.logger, DataErr)
		return responce.Failed(e, responce.BadCode, "invalid data")
	}
	jwtsight := jwt_hash.JWTSigh{}
	err = jwtsight.CreateSigh()
	if err != nil {
		return responce.Failed(e, responce.BadCode, "something go wrong")
	}
	LogService := service.AuthLogic{Data: JsonStruct, Logger: h.logger, JWTSigh: jwtsight}

	token, err := LogService.NewAuth(e.Request().Context(), h.Storage)
	if err != nil {
		return responce.Failed(e, responce.BadCode, "during process happened error")
	}
	return responce.Success(e, responce.BadCode, "you success log in", token)
}
