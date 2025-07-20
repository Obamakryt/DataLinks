package auth

import (
	"DataLinks/internal/api/rest"
	"DataLinks/internal/storages/postgreSQL/TakeConnect"
	"context"
	"log/slog"
)

type ScanQR struct {
	Name     string
	Password string
	Id       string
}
type Auth struct {
	pool   TakeConnect.PgxConnect
	req    rest.RequestLogIn
	logger *slog.Logger
}

func Authorization(a *Auth, ctx context.Context) (ScanQR, error) {
	q := `SELECT name, password, id FROM users WHERE email=$1`
	scanqr := ScanQR{}
	err := a.pool.Pool.QueryRow(ctx, q, a.req.Email).Scan(&scanqr)

	if err != nil {
		retrurnErr := LoggerAuthorization(err, a.logger)
		return scanqr, retrurnErr
	}
	a.logger.Info("Use is find", slog.String("email", a.req.Email))
	return scanqr, nil
}
