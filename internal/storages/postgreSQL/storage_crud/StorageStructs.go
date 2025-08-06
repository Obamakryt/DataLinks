package storage_crud

import (
	"DataLinks/internal/dto/request"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresPool struct {
	Pool *pgxpool.Pool
}

func NewPostgresPool(pool *pgxpool.Pool) *PostgresPool {
	return &PostgresPool{Pool: pool}
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
	NewLink(ctx context.Context, Data request.Add) error
}
type ServiceUserLinks interface {
	TakeChart(ctx context.Context, Data request.TakeChart) ([]string, error)
}
type ServiceUpdateLink interface {
	ChangeCurrentLink(ctx context.Context, Data request.Swap) error
}
type ServiceDeleteLink interface {
	DeleteLink(ctx context.Context, Data request.Delete) error
}

type HandlerStorage struct {
	ServiceNewLinks
	ServiceUserLinks
	ServiceUpdateLink
	ServiceDeleteLink
}
type NewLinks struct {
	Transaction NewLinkTransaction
	DBNewLinks
}

type UserLinks struct {
	DBUserLinks
}
type UpdateLink struct {
	DBUpdateLink
}
type DeleteLink struct {
	DBDeleteLink
}

type AuthReg struct {
	DBAuthReg
}

type DBAuthReg interface {
	Authorization(ctx context.Context, email string) (StorageAuth, error)
	Registration(ctx context.Context, r *StorageRegister) error
}

type Registration interface {
	NewUser(ctx context.Context, DataWithoutHash request.Register) error
}
type Login interface {
	NewAuth(ctx context.Context, Data request.LogIn) (string, error)
}

type HandlerRegAuthStorage struct {
	Registration
	Login
}

type DataUpdateUserLink struct {
	IdUser    int
	IdLink    int
	IdOldLink int
}

type StorageAuth struct {
	Password string
	Id       int
}
type StorageRegister struct {
	Name     string
	Email    string
	HashPass string
}

func NewStorageRegister(pass string, lastData request.Register) *StorageRegister {
	return &StorageRegister{HashPass: pass, Name: lastData.Name, Email: lastData.Email}
}
