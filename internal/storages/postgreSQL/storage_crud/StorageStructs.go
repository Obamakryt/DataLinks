package storage_crud

import (
	"DataLinks/internal/dto/request"
	"DataLinks/internal/storages/postgreSQL/connect_pool"
	"context"
	"github.com/jackc/pgx/v5"
	"log/slog"
)

type PostgresPool struct {
	client connect_pool.PgxConnect
	logger *slog.Logger
}

func NewPostgresPool(pool connect_pool.PgxConnect, log *slog.Logger) *PostgresPool {
	return &PostgresPool{logger: log, client: pool}
}

type AuthReg struct {
	Storage DBAuthReg
}

type NewLinks struct {
	Transaction NewLinkTransaction
	Storage     DBNewLinks
}

type UserLinks struct {
	Storage DBUserLinks
}
type UpdateLink struct {
	Storage DBUpdateLink
}
type DeleteLink struct {
	Storage DBDeleteLink
}

type DBAuthReg interface {
	Authorization(req request.LogIn, ctx context.Context) (StorageAuth, error)
	Registration(r *StorageRegister, ctx context.Context) error
}
type DBNewLinks interface {
	InsertOrFindUrl(ctx context.Context, tx pgx.Tx, url string) (int, error)
	InsertNewUserLink(ctx context.Context, tx pgx.Tx, idUser int, idLink int) error
}

type NewLinkTransaction interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

type DBUserLinks interface {
	TableUserLinks(ctx context.Context, idUser int) ([]string, error)
}

type DBUpdateLink interface {
	InsertNewLink(ctx context.Context, url string) (int, error)
	FindLink(ctx context.Context, url string) (int, error)
	UpdateUserLink(ctx context.Context, data DataUpdateUserLink) error
}
type DBDeleteLink interface {
	FindLink(ctx context.Context, url string) (int, error)
	DeleteLink(ctx context.Context, urlID int, userid int) error
}
type DataUpdateUserLink struct {
	IdUser    int
	IdLink    int
	IdOldLink int
}

type StorageAuth struct {
	Password string
	Id       string
}
type StorageRegister struct {
	Name     string
	Email    string
	HashPass string
}
