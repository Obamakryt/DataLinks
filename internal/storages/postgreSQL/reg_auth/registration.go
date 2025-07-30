package reg_auth

import (
	"DataLinks/internal/dto/request"
	"context"
	"log/slog"
)

type StorageRegister struct {
	Name     string
	Email    string
	HashPass string
}

func NewStorageRegister(pass string, lastData request.Register) *StorageRegister {
	return &StorageRegister{HashPass: pass, Name: lastData.Name, Email: lastData.Email}
}

func (p *PostgresPool) Registration(r *StorageRegister, ctx context.Context) error {
	q := `INSERT INTO users(name, email, password)  VALUES($1, $2, $3) ON CONFLICT (email) DO NOTHING`

	rec, err := p.pool.Pool.Exec(ctx, q, r.Name, r.Email, r.HashPass)

	if err != nil {
		returnErr := LoggerRegistration(rec, err, p.logger)
		return returnErr
	}
	p.logger.Info("User registered",
		slog.String("name", r.Name), slog.String("email", r.Email))
	return nil
}
