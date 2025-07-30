package reg_auth

import (
	"DataLinks/internal/dto/request"
	"DataLinks/internal/storages/postgreSQL/connect_pool"
	"context"
	"log/slog"
)

type StorageAuth struct {
	Password string
	Id       string
}
type PostgresPool struct {
	pool   connect_pool.PgxConnect
	logger *slog.Logger
}

func NewPostgresPool(pool connect_pool.PgxConnect, log *slog.Logger) *PostgresPool {
	return &PostgresPool{logger: log, pool: pool}
}

type GlobalStorage struct {
	Storage DBStorage
}

func NewStorage(storage DBStorage) GlobalStorage {
	return GlobalStorage{Storage: storage}
}

type DBStorage interface {
	Authorization(req request.LogIn, ctx context.Context) (StorageAuth, error)
	Registration(r *StorageRegister, ctx context.Context) error
}

func (p *PostgresPool) Authorization(req request.LogIn, ctx context.Context) (StorageAuth, error) {
	q := `SELECT password, id FROM users WHERE email=$1`
	scanqr := StorageAuth{}
	err := p.pool.Pool.QueryRow(ctx, q, req.Email).Scan(&scanqr)

	if err != nil {
		ReturnErr := LoggerAuthorization(err, p.logger)
		return scanqr, ReturnErr
	}
	p.logger.Info("Use is find", slog.String("email", req.Email))
	return scanqr, nil
}
