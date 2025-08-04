package context_helper

import (
	"DataLinks/internal/dto/responce"
	"DataLinks/internal/service/jwt_hash"
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
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
func MWRequestId(next echo.HandlerFunc) echo.HandlerFunc {
	return func(e echo.Context) error {
		requestID := e.Request().Header.Get("X-Request-Id")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		ctx := context.WithValue(e.Request().Context(), ContextRequestIDKey, requestID)
		newReq := e.Request().WithContext(ctx)
		e.SetRequest(newReq)
		return next(e)
	}
}
func MWUserId(sigh jwt_hash.JWTSigh) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(token, "Bearer ") {
				return c.Redirect(http.StatusUnauthorized, "http://localhost:8080/login")
			}
			token = strings.TrimPrefix(token, "Bearer ")
			subId, err := sigh.ParseJWT(token)
			if err != nil {
				return responce.Failed(c, http.StatusUnauthorized, "invalid token")
			}
			ctx := context.WithValue(c.Request().Context(), ContextUserIdKey, subId)
			ctxWithId := c.Request().WithContext(ctx)
			c.SetRequest(ctxWithId)
			return next(c)
		}
	}
}

type ContextId string

const (
	ContextUserIdKey    ContextId = "user_id"
	ContextRequestIDKey ContextId = "request_id"
)

func ContextGetUserID(ctx context.Context) (int, bool) {
	UserID := ctx.Value(ContextUserIdKey).(int)
	if UserID != 0 {
		return UserID, true
	}
	return -1, false
}
