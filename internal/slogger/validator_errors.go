package slogger

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"log/slog"
)

const (
	InsertNewUser = "InsertNewUser"
)

func LoggerExecInsert(err error, logger *slog.Logger, code CodeOperation) error {
	if errors.Is(err, RowsAffectedError) {
		return LoggerExecCommandTag(code, logger, err)
	}
	if err, ok := Context(err, logger, code); ok {
		return err
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return LoggerPgError(logger, err, code, pgErr)
	}
	logger.Info("unknown error", slog.String("error", err.Error()))
	return fmt.Errorf(SomeError)
}

func LoggerQueryRow(err error, logger *slog.Logger, code CodeOperation) error {
	if err, ok := Context(err, logger, code); ok {
		return err
	}
	if errors.Is(err, pgx.ErrNoRows) {
		err = LoggerErrNoRows(logger, code, err)
	}
	logger.Info("unknown error", slog.String("error", err.Error()))
	return fmt.Errorf(SomeError)
}

func LoggerQuery(err error, logger *slog.Logger, code CodeOperation) error {
	if err, ok := Context(err, logger, code); ok {
		return err
	}
	if errors.Is(err, ScanError) {
		logger.Info("parse error",
			slog.String("error", err.Error()),
			slog.String("op", string(code)))
		return ServerError

	}
	if errors.Is(err, RowsError) {
		logger.Info("Context deadline or rupture",
			slog.String("error", err.Error()),
			slog.String("op", string(code)))
		return ServerError
	}
	logger.Info("unknown error", slog.String("error", err.Error()))
	return fmt.Errorf(SomeError)
}
