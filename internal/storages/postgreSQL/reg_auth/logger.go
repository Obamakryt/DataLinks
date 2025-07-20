package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"log/slog"
)

const someerr = "some unknown error"

func LoggerRegistration(rec pgconn.CommandTag, err error, logger *slog.Logger) error {
	if errors.Is(err, context.DeadlineExceeded) {
		logger.Info("DB Timeout", slog.String("error", err.Error()))
		return fmt.Errorf("failed to connect please try again later")
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
	return fmt.Errorf(someerr)
}

func LoggerAuthorization(err error, logger *slog.Logger) error {
	if errors.Is(err, context.DeadlineExceeded) {
		logger.Info("DB Timeout", slog.String("error", err.Error()))
		return fmt.Errorf("failed to connect please try again later")
	}
	if errors.Is(err, pgx.ErrNoRows) {
		logger.Info("Find Error", slog.String("User", "Not Find"))
		return fmt.Errorf("use with such email not exist")
	}
	return fmt.Errorf(someerr)
}
