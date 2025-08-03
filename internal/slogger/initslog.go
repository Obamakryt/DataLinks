package slogger

import (
	"log/slog"
	"os"
)

const (
	LocalLog = iota
	DevLog
	ProdLog
)

func Setup(level int) *slog.Logger {
	var logger *slog.Logger
	switch level {
	case 0:
		logger = slog.New(slog.NewTextHandler(os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelInfo}))
	case 1:
		logger = slog.New(slog.NewJSONHandler(os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelDebug}))
	case 2:
		logger = slog.New(slog.NewJSONHandler(os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		logger = slog.Default()
	}
	return logger

}
