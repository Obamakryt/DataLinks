package auth

import (
	"DataLinks/internal/api/rest"
	"DataLinks/internal/storages/postgreSQL/TakeConnect"
	"context"
	"log/slog"
)

type DataFromDB struct {
	Password string
	Id       string
}
type PostgresPool struct {
	pool   TakeConnect.PgxConnect
	logger *slog.Logger
}

func NewPostgresPool(pool TakeConnect.PgxConnect, log *slog.Logger) *PostgresPool {
	return &PostgresPool{logger: log, pool: pool}
}

func (p *PostgresPool) Authorization(req rest.RequestLogIn, ctx context.Context) (DataFromDB, error) {
	q := `SELECT password, id FROM users WHERE email=$1`
	scanqr := DataFromDB{}
	err := p.pool.Pool.QueryRow(ctx, q, req.Email).Scan(&scanqr)

	if err != nil {
		retrurnErr := LoggerAuthorization(err, p.logger)
		return scanqr, retrurnErr
	}
	p.logger.Info("Use is find", slog.String("email", req.Email))
	return scanqr, nil
}
