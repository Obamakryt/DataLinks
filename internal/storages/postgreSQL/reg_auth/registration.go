package auth

import (
	"DataLinks/internal/servise/auth"
	"DataLinks/internal/storages/postgreSQL/TakeConnect"
	"context"
	"log/slog"
)

type Reg struct {
	logger *slog.Logger
	ctx    *context.Context
	pool   TakeConnect.PgxConnect
	req    *auth.StorageRegister
}

func Registration(r *Reg) error {
	q := `INSERT INTO users(name, email, password)  VALUES($1, $2, $3) ON CONFLICT (email) DO NOTHING`

	rec, err := r.pool.Pool.Exec(*r.ctx, q, r.req.Name, r.req.Email, r.req.HashPass)

	if err != nil {
		retrurnErr := LoggerRegistration(rec, err, r.logger)
		return retrurnErr
	}
	r.logger.Info("User registered",
		slog.String("name", r.req.Name), slog.String("email", r.req.Email))
	return nil
}
