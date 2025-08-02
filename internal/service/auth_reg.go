package service

import (
	"DataLinks/internal/dto/request"
	"DataLinks/internal/service/jwt_hash"
	auth "DataLinks/internal/storages/postgreSQL/storage_crud"
	"context"
	"fmt"
	"log/slog"
)

type LogicReg struct {
	DataWithoutHash request.Register
	Logger          *slog.Logger
}
type AuthLogic struct {
	Data   request.LogIn
	Logger *slog.Logger
	jwt_hash.JWTSigh
}

func (r *LogicReg) NewUser(ctx context.Context, storage auth.AuthReg) error {
	hashpass := jwt_hash.HashingPass(r.DataWithoutHash.Password)
	reg := auth.NewStorageRegister(hashpass, r.DataWithoutHash)
	err := storage.Storage.Registration(reg, ctx)
	if err != nil {
		r.Logger.Warn("Failed create new user", slog.String("Error", err.Error()))
		return fmt.Errorf("failed create new user")
	}
	return nil
}

func (a *AuthLogic) NewAuth(ctx context.Context, storage auth.AuthReg) (string, error) {
	DataUser, err := storage.Storage.Authorization(a.Data, ctx)
	if err != nil {
		a.Logger.Warn("Failed find user", slog.String("Error", err.Error()))
		return "", fmt.Errorf("failed find user")
	}
	PassRight := jwt_hash.CheckHashPass(a.Data.Password, DataUser.Password)
	if !PassRight {
		a.Logger.Warn("Failed incorrect password")
		return "", fmt.Errorf("incorrect password")
	}
	if a.JWTSigh.Secret != "" {
		token, err := a.JWTSigh.GenerateJWT(DataUser.Id, 15)
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
