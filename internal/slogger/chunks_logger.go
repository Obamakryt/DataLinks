package slogger

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"log/slog"
)

var RowsAffectedError = errors.New("RowsAffectedError")

// Exec only Storage logger! for pgx Exec operation its help find RowsAffected immediately and in service layer we test suffix on
// / "no changes were made" after that we have custom logger for RowsAffected type error
func Exec(err error, tag pgconn.CommandTag) error {
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return RowsAffectedError
	}
	return nil
}
func LoggerExecCommandTag(code CodeOperation, logger *slog.Logger, err error) error {
	switch code {
	case E101:
		logger.Info("Link already exists",
			slog.String("code", string(code)),
			slog.String("error", err.Error()))
		return fmt.Errorf("user already has this link")

	case E103:
		logger.Info("non existing link in user chart",
			slog.String("code", string(code)),
			slog.String("error", err.Error()))
		return fmt.Errorf("user does not have this link")

	default:
		logger.Error("Unhandled code",
			slog.String("code", string(code)),
			slog.String("error", err.Error()))
		return fmt.Errorf("unexpected operation failure")
	}

}

func LoggerPgError(logger *slog.Logger, err error, code CodeOperation, pgError *pgconn.PgError) error {
	switch code {
	case E1041:
		if pgError.Code == "23505" {
			logger.Error("Link change failed",
				slog.String("code", string(code)),
				slog.String("error", err.Error()))
			return fmt.Errorf("cannot change to already existing link")
		}
	}
	return nil
}

func LoggerErrNoRows(logger *slog.Logger, code CodeOperation, err error) error {
	switch code {
	//case IsAuthorization:
	//	logger.Info("Find Error", slog.String("User", "Not Find"))
	//	return fmt.Errorf("user with such email not exist")
	case E1011:
		logger.Info("Create Error",
			slog.String("url", "already exist"),
			slog.String("code", string(code)))
		slog.String("error", err.Error())
		return nil
	case E1042, E1031:
		logger.Info("Find Error",
			slog.String("error", err.Error()),
			slog.String("code", string(code)))
		return fmt.Errorf("non-existent link")
	case E1043:
		logger.Info("Create Error",
			slog.String("url", "already exist"),
			slog.String("code", string(code)))
		slog.String("error", err.Error())
		return nil
	}
	return nil
}

func Context(err error, logger *slog.Logger, code CodeOperation) (error, bool) {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		logger.Info("DB Timeout",
			slog.String("error", err.Error()),
			slog.Any("Operation", code))
		return fmt.Errorf("db timeout"), true

	case errors.Is(err, context.Canceled):
		logger.Info("Connect canceled",
			slog.String("error", err.Error()),
			slog.Any("Operation", code))
		return fmt.Errorf("rupture connection"), true
	}
	return err, false
}
