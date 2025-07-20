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
	logger          *slog.Logger
}
type Storage struct {
	storage *auth.RegistrationDB
}

func NewStorage(storage *auth.RegistrationDB) *Storage {
	return &Storage{storage: storage}
}

func (r *RegLogic) NewUser(ctx context.Context, storage Storage) error {
	hashpass := HashingPass(r.DataWithoutHash.Password)
	reg := auth.NewStorageRegister(hashpass, r.DataWithoutHash)
	err := reg.Registration(storage.storage, ctx)
	if err != nil {
		r.logger.Warn("Failed create new user", slog.String("Error", err.Error()))
		return fmt.Errorf("failed create new user")
	}
	return nil
}
