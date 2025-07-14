package postgreSQL

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxConnect struct {
	Pool *pgxpool.Pool
}

func CreatePool() *pgxpool.Pool {
	pgxpool.New(context.Background())
}
