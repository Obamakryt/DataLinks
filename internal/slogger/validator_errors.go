package slogger

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"log/slog"
	"strings"
)

// TODO куча говна переделать
const (
	InsertNewUser      = "InsertNewUser"
	CreateNewLinksUser = "CreateNewLinksUser"
	ChangeLinksUser    = "ChangeLinksUser"
	IsAuthorization    = "Authorization"
	CreateNewURL       = "CreateNewURL"
	FindOldLink        = "FindOldLink"
	ReplacementOldLink = "ReplacementOldLink"
	SomeError          = "internal server error"
)

func LoggerExecInsert(err error, logger *slog.Logger, typeLogger string) error {
	if strings.HasSuffix(err.Error(), "no changes were made") {
		switch typeLogger {
		case InsertNewUser:
			logger.Info("user already exist", slog.String("error", err.Error()))
			return fmt.Errorf("user with this data already exist")
		case CreateNewLinksUser:
			logger.Info("link already exist", slog.String("error", err.Error()))
			return fmt.Errorf("user already has this link")
		case ChangeLinksUser:
			logger.Info("during swap link something go wrong", slog.String("error", err.Error()))
			return fmt.Errorf(SomeError)
		}
	}
	if err, ok := Context(err, logger); ok {
		return err
	}
	logger.Info("unknown error", slog.String("error", err.Error()))
	return fmt.Errorf(SomeError)
}

func LoggerQueryRow(err error, logger *slog.Logger, typeLogger string) error {
	if err, ok := Context(err, logger); ok {
		return err
	}
	if errors.Is(err, pgx.ErrNoRows) {
		switch typeLogger {
		case IsAuthorization:
			logger.Info("Find Error", slog.String("User", "Not Find"))
			return fmt.Errorf("user with such email not exist")
		case CreateNewURL:
			logger.Info("Create Error", slog.String("Urls", "Not Add"))
			return fmt.Errorf("failed to create new URL")
		case FindOldLink:
			logger.Info("Find Error", slog.String("during swap process we dont find", "current url"))
			return fmt.Errorf("invalid current url")
		case ReplacementOldLink:
			logger.Info("Insert Error user already has this link")
			return fmt.Errorf("try repiatedly insert link")
		}
	}
	logger.Info("unknown error", slog.String("error", err.Error()))
	return fmt.Errorf(SomeError)
}

func LoggerQuery(err error, logger *slog.Logger) error {
	if err, ok := Context(err, logger); ok {
		return err
	}
	if strings.HasPrefix(err.Error(), "rows error") {
		logger.Info("Rows Error", slog.String("storage give error", err.Error()))
		return fmt.Errorf("something went wrong")
	}
	if strings.HasPrefix(err.Error(), "failed to scan link") {
		logger.Info("Link Error", slog.String("storage give error", err.Error()))
	}
	logger.Info("unknown error", slog.String("error", err.Error()))
	return fmt.Errorf(SomeError)
}

func Context(err error, logger *slog.Logger) (error, bool) {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		logger.Info("DB Timeout", slog.String("error", err.Error()))
		return fmt.Errorf("failed to connect please try again later"), true
	case errors.Is(err, context.Canceled):
		logger.Info("Connect canceled", slog.String("error", err.Error()))
		return fmt.Errorf("user close connection"), true
	}

	return err, false
}

func Exec(err error, tag pgconn.CommandTag, msg string) error {
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s - no changes were made", msg)
	}
	return nil
}
