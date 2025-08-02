package responce

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	GoodCode = http.StatusOK
	BadCode  = http.StatusBadRequest
)

type ResponseAuth struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func Success(e echo.Context, code int, mes string, data interface{}) error {
	if data != nil {
		return e.JSON(code, ResponseAuth{
			Message: mes,
			Data:    data,
		})
	} else {
		return e.JSON(code, ResponseAuth{
			Message: mes,
		})
	}
}

func Failed(e echo.Context, code int, err string) error {
	return e.JSON(code, ResponseAuth{
		Error: err,
	})
}
