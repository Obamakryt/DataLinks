package auth_reg

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type JWTSigh struct {
	secret string `env:"JWT_SIGH"`
}

func (j *JWTSigh) CreateSigh() error {
	err := cleanenv.ReadEnv(&j)
	if err != nil {
		return fmt.Errorf("couldnt find postgres password in env %w", err)
	}
	return nil
}

func (j *JWTSigh) generateJWT(id string, tl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   id,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tl)),
	})

	return token.SignedString([]byte(j.secret))
}

func (j *JWTSigh) parseJWT(reqToken string) (string, error) {
	token, err := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected sigh method")
		}
		return []byte(j.secret), nil
	})
	if err != nil || !token.Valid {
		return "", fmt.Errorf("no valid token")
	}
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("failed claim payload")
	}
	subId, ok := claim["sub"].(string)
	if !ok {
		return "", fmt.Errorf("token dont have expected payload")
	}
	return subId, nil

}
