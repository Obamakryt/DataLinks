package auth

import (
	"DataLinks/internal/api/rest"
	auth "DataLinks/internal/storages/postgreSQL/reg_auth"
	"context"
	"fmt"
	"log/slog"
)

type RegLogic struct {
	DataWithoutHash rest.RequestRegister
	Logger          *slog.Logger
}
type AuthLogic struct {
	Data   rest.RequestLogIn
	Logger *slog.Logger
	JWTSigh
}

type Storage struct {
	Storage *auth.PostgresPool
}

func NewStorage(storage *auth.PostgresPool) *Storage {
	return &Storage{Storage: storage}
}

func (r *RegLogic) NewUser(ctx context.Context, storage Storage) error {
	hashpass := HashingPass(r.DataWithoutHash.Password)
	reg := auth.NewStorageRegister(hashpass, r.DataWithoutHash)
	err := storage.Storage.Registration(reg, ctx)
	if err != nil {
		r.Logger.Warn("Failed create new user", slog.String("Error", err.Error()))
		return fmt.Errorf("failed create new user")
	}
	return nil
}

func (a *AuthLogic) NewAuth(ctx context.Context, storage Storage) (string, error) {
	DataUser, err := storage.Storage.Authorization(a.Data, ctx)
	if err != nil {
		a.Logger.Warn("Failed find user", slog.String("Error", err.Error()))
		return "", fmt.Errorf("failed find user")
	}
	PassRight := CheckHashPass(a.Data.Password, DataUser.Password)
	if !PassRight {
		a.Logger.Warn("Failed incorrect password")
		return "", fmt.Errorf("incorrect password")
	}
	if a.JWTSigh.secret != "" {
		token, err := a.JWTSigh.generateJWT(DataUser.Id, 15)
		if err != nil {
			a.Logger.Warn("something wrong on part of sigh token")
			return "", fmt.Errorf("failed create token")
		}
		return token, nil
	} else {
		a.Logger.Info("where secret sigh?")
		return "", fmt.Errorf("failed create token")
	}
}
