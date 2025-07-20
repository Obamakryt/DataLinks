package auth

import (
	"DataLinks/internal/api/rest"
	"DataLinks/internal/storages/postgreSQL/TakeConnect"
	"context"
	"log/slog"
)

type StorageRegister struct {
	Name     string
	Email    string
	HashPass string
}

type RegistrationDB struct {
	logger *slog.Logger
	pool   TakeConnect.PgxConnect
}

func NewStorageRegister(pass string, lastData rest.RequestRegister) *StorageRegister {
	return &StorageRegister{HashPass: pass, Name: lastData.Name, Email: lastData.Email}
}
func NewRegistratDB(pool TakeConnect.PgxConnect, log *slog.Logger) *RegistrationDB {
	return &RegistrationDB{logger: log, pool: pool}
}

func (s *StorageRegister) Registration(r *RegistrationDB, ctx context.Context) error {
	q := `INSERT INTO users(name, email, password)  VALUES($1, $2, $3) ON CONFLICT (email) DO NOTHING`

	rec, err := r.pool.Pool.Exec(ctx, q, s.Name, s.Email, s.HashPass)

	if err != nil {
		returnErr := LoggerRegistration(rec, err, r.logger)
		return returnErr
	}
	r.logger.Info("User registered",
		slog.String("name", s.Name), slog.String("email", s.Email))
	return nil
}
