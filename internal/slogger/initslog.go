package slogger

import (
	"log/slog"
	"os"
)

const (
	LocalLog = iota
	ProdLog
)

type Setup struct {
	*slog.Logger
}

func (l *Setup) SetupLogger(level int) {
	switch level {
	case LocalLog:
		l.Logger = slog.New(slog.NewTextHandler(os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelDebug}))
	case ProdLog:
		l.Logger = slog.New(slog.NewJSONHandler(os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelError}))
	}
}

//func (l *Setup) LoggerUserRequestID(ctx context.Context, keys ...any) error {
//	var NewLogger *Setup
//	for _, key := range keys {
//		switch key {
//		case ctxhelp.ContextUserIdKey:
//			userID, ok := ctx.Value(ctxhelp.ContextUserIdKey).(int)
//			if !ok {
//				return
//			}
//			NewLogger = &Setup{l.With(slog.Int("user_id", userID))}
//		case ctxhelp.ContextRequestIDKey:
//			requestID, ok := ctx.Value(ctxhelp.ContextRequestIDKey).(string)
//			if !ok {
//				return l
//			}
//			NewLogger = &Setup{l.With(slog.String("request_id", requestID))}
//		}
//	}
//	return NewLogger
//}
