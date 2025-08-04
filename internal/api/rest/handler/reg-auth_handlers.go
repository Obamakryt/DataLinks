package handler

import (
	"DataLinks/internal/dto/request"
	"DataLinks/internal/dto/responce"
	"DataLinks/internal/service"
	"DataLinks/internal/service/jwt_hash"
	servise "DataLinks/internal/storages/postgreSQL/storage_crud"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"log/slog"
)

type AuthReg struct {
	Validator *validator.Validate
	logger    *slog.Logger
	Storage   servise.AuthReg
}

func NewHandler(validator *validator.Validate, logger *slog.Logger, storage servise.AuthReg) *AuthReg {
	return &AuthReg{Validator: validator, logger: logger, Storage: storage}
}

func (h *AuthReg) RegHandler(e echo.Context) error {
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

func (h *AuthReg) LogInHandler(e echo.Context) error {
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

type LinkHandler struct {
	Validator *validator.Validate
	logger    *slog.Logger
	Storage   servise.HandlerStorage
}

func AddNewLinkHandler(e echo.Context) error {
	da
}
