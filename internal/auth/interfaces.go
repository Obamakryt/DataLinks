package auth

import (
	"DataLinks/internal/storages/postgreSQL/TakeConnect"
	"log/slog"
)

type AuthService interface {
	Register(pool *TakeConnect.PgxConnect, logger *slog.Logger) error
	Login()
}
