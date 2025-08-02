package jwt_hash

import (
	"DataLinks/internal/dto/responce"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type JWTSigh struct {
	Secret string `env:"JWT_SIGH"`
}

func (j *JWTSigh) CreateSigh() error {
	err := cleanenv.ReadEnv(&j)
	if err != nil {
		return fmt.Errorf("couldnt find postgres password in env %w", err)
	}
	return nil
}

func (j *JWTSigh) GenerateJWT(id string, tl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   id,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tl)),
	})

	return token.SignedString([]byte(j.Secret))
}

func (j *JWTSigh) parseJWT(reqToken string) (int, error) {
	token, err := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected sigh method")
		}
		return []byte(j.Secret), nil
	})
	if err != nil || !token.Valid {
		return -1, fmt.Errorf("no valid token")
	}
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return -1, fmt.Errorf("failed claim payload")
	}
	subId, ok := claim["sub"].(string)
	if !ok {
		return -1, fmt.Errorf("token dont have expected payload")
	}
	intSubId, err := strconv.Atoi(subId)
	if err != nil {
		return -1, fmt.Errorf("invalid payload")
	}
	return intSubId, nil

}
func Middleware(sigh JWTSigh) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(token, "Bearer ") {
				return c.Redirect(http.StatusUnauthorized, "http://localhost:8080/login")
			}
			token = strings.TrimPrefix(token, "Bearer ")
			subId, err := sigh.parseJWT(token)
			if err != nil {
				return responce.Failed(c, http.StatusUnauthorized, "invalid token")
			}
			c.Set("user_id", subId)
			return next(c)
		}
	}
}
