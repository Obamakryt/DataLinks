package storage_crud

import (
	"DataLinks/internal/dto/request"
	"DataLinks/internal/slogger"
	"context"
	"log/slog"
)

// TODO FK THIS SHIT REMAKE
func (p *PostgresPool) Authorization(req request.LogIn, ctx context.Context) (StorageAuth, error) {
	q := `SELECT password, id FROM users WHERE email=$1;`
	scanqr := StorageAuth{}
	err := p.client.Pool.QueryRow(ctx, q, req.Email).Scan(&scanqr)

	if err != nil {
		ReturnErr := slogger.LoggerQueryRow(err, p.logger, "Authorization")
		return scanqr, ReturnErr
	}
	p.logger.Info("Use is find", slog.String("email", req.Email))
	return scanqr, nil
}
func NewStorageRegister(pass string, lastData request.Register) *StorageRegister {
	return &StorageRegister{HashPass: pass, Name: lastData.Name, Email: lastData.Email}
}

// TODO FK THIS SHIT REMAKE HATE THIS REALIZATION

func (p *PostgresPool) Registration(r *StorageRegister, ctx context.Context) error {
	q := `INSERT INTO users(name, email, password)  VALUES($1, $2, $3) ON CONFLICT (email) DO NOTHING`

	commandTag, err := p.client.Pool.Exec(ctx, q, r.Name, r.Email, r.HashPass)
	if err != nil {
		returnErr := slogger.LoggerExecInsert(err, p.logger, slogger.InsertNewUser)
		return returnErr
	}
	p.logger.Info("User registered",
		slog.String("name", r.Name), slog.String("email", r.Email))
	return nil
}
