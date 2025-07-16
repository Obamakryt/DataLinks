package auth

import (
	"DataLinks/internal/servise/auth"
	"DataLinks/internal/storages/postgreSQL/TakeConnect"
	"context"
	"log/slog"
)

type ScanQR struct {
	Name     string
	Password string
}
type Auth struct {
	ctx    *context.Context
	pool   TakeConnect.PgxConnect
	req    auth.RequestLogIn
	logger *slog.Logger
}

func Authorization(a *Auth) (ScanQR, error) {
	q := `SELECT name, password FROM users WHERE email=$1`
	scanqr := ScanQR{}
	err := a.pool.Pool.QueryRow(*a.ctx, q, a.req.Email).Scan(&scanqr)

	if err != nil {
		retrurnErr := LoggerAuthorization(err, a.logger)
		return scanqr, retrurnErr
	}
	a.logger.Info("Use is find", slog.String("email", a.req.Email))
	return scanqr, nil
}
