package service

import (
	"DataLinks/internal/dto/request"
	"DataLinks/internal/slogger"
	"DataLinks/internal/storages/postgreSQL/storage_crud"
	"DataLinks/pkg/context_helper"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"golang.org/x/sync/errgroup"
	"log/slog"

	"time"
)

type LogicLinkAdd struct {
	Data   request.Add
	Logger *slog.Logger
}

func (l *LogicLinkAdd) NewLink(ctx context.Context, storage storage_crud.NewLinks) error {
	userID, ok := context_helper.ContextHaveUserID(ctx)
	if !ok {
		return errors.New("user id is required")
	}
	tx, err := storage.Transaction.Begin(ctx)
	if err != nil {
		l.Logger.Warn("failed start transaction", slog.String("error", err.Error()))
		return err
	}

	timeNow := time.Now()
	defer func() {
		timeoutime := time.Second * 3
		timeTx := time.Since(timeNow) / 5
		if timeTx < timeoutime {
			timeoutime = timeTx
		}
		txRollback, cancel := context.WithTimeout(context.Background(), timeoutime)
		defer cancel()
		ErrorTx := tx.Rollback(txRollback)
		if ErrorTx != nil && !errors.Is(ErrorTx, pgx.ErrTxClosed) {
			l.Logger.Warn("failed rollback transaction",
				slog.String("error", ErrorTx.Error()))
			//slog.Any("id_request", ctx.Value("id_request").(string))
		}
	}()

	idLink, err := storage.Storage.InsertOrFindUrl(ctx, l.Data.Url)
	if err != nil {
		err = slogger.LoggerQueryRow(err, l.Logger, "")
		return err
	}

	insert, err := storage.Storage.InsertNewLink(ctx, userID, idLink)
	if err != nil {
		err = slogger.LoggerExecInsert(insert, err, l.Logger, "CreateNewLinksUser")
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		l.Logger.Warn("failed commit transaction", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func TakeChart(ctx context.Context, storage storage_crud.UserLinks, logger *slog.Logger) ([]string, error) {
	userID, ok := context_helper.ContextHaveUserID(ctx)
	if !ok {
		return nil, errors.New("user id is required")
	}
	links, err := storage.Storage.TableUserLinks(ctx, userID)
	if err != nil {
		err = slogger.LoggerQuery(err, logger)
		return nil, err
	}
	if len(links) == 0 {
		logger.Info("no links found for user", slog.Int("id", userID))
		return links, nil
	}
	logger.Info("user make request - result finding",
		slog.Int("count", len(links)),
		slog.Int("id", userID))
	return links, nil
}

type LogicUpdateLink struct {
	Data   request.Change
	Logger *slog.Logger
}

func (l *LogicUpdateLink) ChangeCurrentLink(ctx context.Context, storage storage_crud.DBUpdateLink) error {
	userID, ok := context_helper.ContextHaveUserID(ctx)
	if !ok {
		return errors.New("user id is required")
	}
	g, ctx := errgroup.WithContext(ctx)
	var newID, oldID int
	g.Go(func() error {
		oldid, err := storage.FindOldLink(ctx, l.Data.Url)
		if err != nil {
			err = slogger.LoggerQueryRow(err, l.Logger, "FindOldLink")
			return err
		}
		oldID = oldid
		return nil
	})
	g.Go(func() error {
		newid, err := storage.ChangeLink(ctx, l.Data.NewUrl)
		if err != nil {
			err = slogger.LoggerQueryRow(err, l.Logger, "ReplacementOldLink")
			return err
		}
		newID = newid
		return nil

	})
	if err := g.Wait(); err != nil {
		return err
	}
	data := storage_crud.DataUpdateUserLink{IdUser: userID, IdOldLink: oldID, IdLink: newID}

	err := storage.UpdateUserLink(ctx, data)
	if err != nil {
		err = slogger.LoggerExecInsert(false, err, l.Logger, "ChangeLinksUser")
	}
	return nil
}
