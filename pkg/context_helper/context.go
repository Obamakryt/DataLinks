package context_helper

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"time"
)

func WithTimeout(timeout time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) error {
			ctx, cancel := context.WithTimeout(e.Request().Context(), timeout)
			defer cancel()
			NewReq := e.Request().WithContext(ctx)
			e.SetRequest(NewReq)
			return next(e)
		}
	}
}
func RequestId(next echo.HandlerFunc) echo.HandlerFunc {
	return func(e echo.Context) error {
		requestID := e.Request().Header.Get("X-Request-Id")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		ctx := context.WithValue(e.Request().Context(), "request_id", requestID)
		NewReq := e.Request().WithContext(ctx)
		e.SetRequest(NewReq)
		return next(e)
	}
}

func ContextGetUserID(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		return -1, false
	}
	return userID, true
}
func ContextGetRequestID(ctx context.Context) (int, bool) {
	requestID, ok := ctx.Value("request_id").(int)
	if !ok {
		return -1, false
	}
	return requestID, true
}
