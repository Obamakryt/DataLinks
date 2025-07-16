package auth

import (
	"DataLinks/internal/servise/auth"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"log/slog"
)

func LoggerRegistration(rec pgconn.CommandTag, err error, logger *slog.Logger, req *auth.StorageRegister) error {
	if errors.Is(err, context.DeadlineExceeded) {
		logger.Info("DB Timeout", slog.String("error", err.Error()))
		fmt.Errorf("failed to connect please try again later")
	}
	if rec.RowsAffected() == 0 {
		logger.Info("Already exist", slog.String("error", err.Error()))
		return fmt.Errorf("user with this data already exist")
	}
	var PgErr *pgconn.PgError
	if errors.As(err, &PgErr) {
		logger.Info("Failed PGX ", PgErr.Message)
		return fmt.Errorf("failed please try again later")
	}
	return fmt.Errorf("some unknown error")

}
