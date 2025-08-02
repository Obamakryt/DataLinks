package storage_crud

import (
	"DataLinks/internal/dto/request"
	"DataLinks/internal/slogger"
	"context"
	"log/slog"
)

func NewStorageRegister(pass string, lastData request.Register) *StorageRegister {
	return &StorageRegister{HashPass: pass, Name: lastData.Name, Email: lastData.Email}
}

// TODO FK THIS SHIT REMAKE HATE THIS REALIZATION

func (p *PostgresPool) Registration(r *StorageRegister, ctx context.Context) error {
	q := `INSERT INTO users(name, email, password)  VALUES($1, $2, $3) ON CONFLICT (email) DO NOTHING`

	commandTag, err := p.Client.Pool.Exec(ctx, q, r.Name, r.Email, r.HashPass)
	if err != nil {
		returnErr := slogger.LoggerExecInsert(true, err, p.logger, slogger.InsertNewUser)
		return returnErr
	}
	p.logger.Info("User registered",
		slog.String("name", r.Name), slog.String("email", r.Email))
	return nil
}
