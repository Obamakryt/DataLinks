package context_helper

import (
	"context"
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

func ContextHaveUserID(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		return -1, false
	}
	return userID, true
}
