package connect_pool

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"time"
)

type PgxConnect struct {
	Pool *pgxpool.Pool
}

func CreatePgxUrl(url PostgresUrl) string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		url.Username,
		url.Pass,
		url.Hostname,
		url.Port,
		url.DBName)
}

func NewPool(connurl string, try int, logger *slog.Logger) (*pgxpool.Pool, error) {
	var err error
	for i := 0; i < try; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		logger.Info("Connect pgx try num - ", try)
		pool, e := pgxpool.New(ctx, connurl)
		cancel()
		if e != nil {
			err = e
			logger.Warn("Connect failed", slog.String("error", err.Error()))
			continue
		} else {
			ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
			e = pool.Ping(ctx)
			cancel()
			if e != nil {
				pool.Close()
				logger.Warn("failed ping pool connect ", err)
				err = e
				continue
			}

			return pool, nil
		}
	}

	var pgErr *pgconn.PgError
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return nil, fmt.Errorf("no one of %d attempts couldn connect in the specified time", try)
	case errors.As(err, &pgErr):
		return nil, fmt.Errorf("pgx error: %s, %s", pgErr.Code, pgErr.Message)
	default:
		return nil, fmt.Errorf("0 info about connection error")

	}

}
