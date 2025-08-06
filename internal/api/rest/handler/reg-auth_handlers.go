package handler

import (
	"DataLinks/internal/dto/request"
	"DataLinks/internal/dto/responce"
	"DataLinks/internal/slogger"
	servise "DataLinks/internal/storages/postgreSQL/storage_crud"
	"DataLinks/pkg/context_helper"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type AuthRegHandler struct {
	Validator *validator.Validate
	Logger    slogger.Setup
	Storage   servise.HandlerRegAuthStorage
}

func (h *AuthRegHandler) RegHandler(e echo.Context) error {
	var data request.Register
	//h.Logger.LoggerUserRequestID(e.Request().Context(), context_helper.ContextRequestIDKey)
	err := e.Bind(&data)
	if err != nil {
		return responce.Failed(e, responce.BadCode, responce.InvalidData)
	}
	log.Info("password  ", data.Password)

	err = h.Validator.Struct(data)
	log.Info("password ", data.Password)
	if err != nil {
		DataErr := ErrorValid(err)
		LoggerValidatorError(h.Logger, DataErr)
		return responce.Failed(e, responce.BadCode, responce.InvalidData)
	}

	err = h.Storage.Registration.NewUser(e.Request().Context(), data)
	if err != nil {
		return responce.Failed(e, responce.BadCode, err.Error())
	}

	return responce.Success(e, responce.GoodCode, "user create successful", "/login")
}

func (h *AuthRegHandler) LoginHandler(e echo.Context) error {
	var data request.LogIn
	err := e.Bind(&data)
	if err != nil {
		return responce.Failed(e, responce.BadCode, responce.InvalidData)
	}
	err = h.Validator.Struct(data)
	if err != nil {
		DataErr := ErrorValid(err)
		LoggerValidatorError(h.Logger, DataErr)
		return responce.Failed(e, responce.BadCode, responce.InvalidData)
	}
	token, err := h.Storage.Login.NewAuth(e.Request().Context(), data)

	if err != nil {
		return responce.Failed(e, responce.BadCode, err.Error())
	}
	return responce.Success(e, responce.GoodCode, "successful login", token)
}

type LinkHandler struct {
	Validator *validator.Validate
	Logger    slogger.Setup
	Storage   servise.HandlerStorage
}

func (l *LinkHandler) AddNewLinkHandler(e echo.Context) error {
	var data request.Add

	ok := BindAndValidate(&data, e, l.Validator, l.Logger)
	if !ok {
		return responce.Failed(e, responce.BadCode, responce.InvalidData)
	}
	ctx := e.Request().Context()
	userID, ok := context_helper.ContextGetUserID(ctx)
	if !ok {
		return responce.Failed(e, responce.BadCode, responce.InvalidData)
	}

	data.UserId = userID
	err := l.Storage.ServiceNewLinks.NewLink(ctx, data)
	if err != nil {
		return responce.Failed(e, responce.BadCode, err.Error())
	}
	return responce.Success(e, responce.GoodCode, "link add", nil)
}

func (l *LinkHandler) GetUserLinksHandler(e echo.Context) error {
	ctx := e.Request().Context()
	userID, ok := context_helper.ContextGetUserID(ctx)
	if !ok {
		return responce.Failed(e, responce.BadCode, responce.InvalidData)
	}
	data := request.TakeChart{UserId: userID}
	chart, err := l.Storage.ServiceUserLinks.TakeChart(ctx, data)
	if err != nil {
		return responce.Failed(e, responce.BadCode, err.Error())
	}
	return responce.Success(e, responce.GoodCode, "", chart)
}

func (l *LinkHandler) DeleteLinkHandler(e echo.Context) error {
	var data request.Delete
	ok := BindAndValidate(&data, e, l.Validator, l.Logger)
	if !ok {
		return responce.Failed(e, responce.BadCode, responce.InvalidData)
	}
	ctx := e.Request().Context()
	userID, ok := context_helper.ContextGetUserID(ctx)
	if !ok {
		return responce.Failed(e, responce.BadCode, responce.InvalidData)
	}
	data.UserId = userID
	err := l.Storage.ServiceDeleteLink.DeleteLink(ctx, data)
	if err != nil {
		return responce.Failed(e, responce.BadCode, err.Error())
	}
	return responce.Success(e, responce.GoodCode, "link delete", nil)
}

func (l *LinkHandler) SwapLinks(e echo.Context) error {
	var data request.Swap
	ok := BindAndValidate(&data, e, l.Validator, l.Logger)
	if !ok {
		return responce.Failed(e, responce.BadCode, responce.InvalidData)
	}
	ctx := e.Request().Context()
	userID, ok := context_helper.ContextGetUserID(ctx)
	if !ok {
		return responce.Failed(e, responce.BadCode, responce.InvalidData)
	}
	data.UserId = userID

	err := l.Storage.ServiceUpdateLink.ChangeCurrentLink(ctx, data)
	if err != nil {
		return responce.Failed(e, responce.BadCode, err.Error())
	}
	return responce.Success(e, responce.GoodCode, "link change", nil)
}
