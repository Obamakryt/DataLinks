package service

import (
	"DataLinks/internal/dto/request"
	"DataLinks/internal/service/jwt_hash"
	"DataLinks/internal/slogger"
	auth "DataLinks/internal/storages/postgreSQL/storage_crud"
	"context"
	"log/slog"
)

type LogicReg struct {
	Logger  slogger.Setup
	Storage auth.AuthReg
}

func (r *LogicReg) NewUser(ctx context.Context, DataWithoutHash request.Register) error {
	hashpass := jwt_hash.HashingPass(DataWithoutHash.Password)
	reg := auth.NewStorageRegister(hashpass, DataWithoutHash)
	err := r.Storage.Registration(ctx, reg)
	if err != nil {
		returnErr := slogger.LoggerExecInsert(err, r.Logger, slogger.E001)
		return returnErr
	}

	r.Logger.Info("User registered",
		slog.String("name", DataWithoutHash.Name),
		slog.String("email", DataWithoutHash.Email))
	return nil
}

type LogicAuth struct {
	Logger  slogger.Setup
	Storage auth.AuthReg
}

func (a *LogicAuth) NewAuth(ctx context.Context, Data request.LogIn) (string, error) {
	DataUser, err := a.Storage.Authorization(ctx, Data.Email)
	if err != nil {
		return "", slogger.LoggerQueryRow(err, a.Logger, slogger.E002)
	}
	PassRight := jwt_hash.VerifyPassword(Data.Password, DataUser.Password)
	if !PassRight {
		a.Logger.Info("Failed Login incorrect password", slog.String("email", Data.Email))
		return "", slogger.AuthError
	}
	sight := jwt_hash.JWTSigh{}
	err = sight.CreateSigh()
	if err != nil {
		return "", slogger.ServerError
	}
	if sight.Secret == "" {
		a.Logger.Info("lost jwt sigh")
		return "", slogger.ServerError
	}
	token, err := sight.GenerateJWT(DataUser.Id, jwt_hash.TimeLive)
	if err != nil {
		a.Logger.Info("something wrong on part of sigh token",
			slog.String("error", err.Error()))
		return "", slogger.ServerError
	}
	a.Logger.Info("User login", slog.String("email", Data.Email))
	return token, nil
}
