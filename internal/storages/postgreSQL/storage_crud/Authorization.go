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
	err := p.Client.Pool.QueryRow(ctx, q, req.Email).Scan(&scanqr)

	if err != nil {
		ReturnErr := slogger.LoggerQueryRow(err, p.logger, "Authorization")
		return scanqr, ReturnErr
	}
	p.logger.Info("Use is find", slog.String("email", req.Email))
	return scanqr, nil
}
