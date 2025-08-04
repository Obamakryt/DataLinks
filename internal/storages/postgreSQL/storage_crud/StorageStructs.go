package storage_crud

import (
	"DataLinks/internal/dto/request"
	"DataLinks/internal/slogger"
	"DataLinks/internal/storages/postgreSQL/connect_pool"
	"context"
	"github.com/jackc/pgx/v5"
)

type PostgresPool struct {
	client connect_pool.PgxConnect
	logger slogger.Setup
}

func NewPostgresPool(pool connect_pool.PgxConnect, log slogger.Setup) *PostgresPool {
	return &PostgresPool{logger: log, client: pool}
}

type AuthReg struct {
	Storage DBAuthReg
}
type HandlerStorage struct {
	ServiceNewLinks
	ServiceUserLinks
	ServiceUpdateLink
	ServiceDeleteLink
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
	InsertOrFindUrlTx(ctx context.Context, tx pgx.Tx, url string) (int, error)
	InsertNewUserLinkTx(ctx context.Context, tx pgx.Tx, idUser int, idLink int) error
}

type NewLinkTransaction interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

type DBUserLinks interface {
	TableUserLinks(ctx context.Context, idUser int) ([]string, error)
}

type DBUpdateLink interface {
	InsertOrFindUrl(ctx context.Context, url string) (int, error)
	FindLink(ctx context.Context, url string) (int, error)
	ChangeUserLink(ctx context.Context, data DataUpdateUserLink) error
}
type DBDeleteLink interface {
	FindLink(ctx context.Context, url string) (int, error)
	DeleteUserLinkAssociation(ctx context.Context, urlID int, userid int) error
}

type ServiceNewLinks interface {
	NewLink(ctx context.Context) error
}
type ServiceUserLinks interface {
	TakeChart(ctx context.Context) ([]string, error)
}
type ServiceUpdateLink interface {
	ChangeCurrentLink(ctx context.Context) error
}
type ServiceDeleteLink interface {
	DeleteLink(ctx context.Context) error
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
