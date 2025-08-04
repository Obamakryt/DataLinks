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
	Data request.Add
	slogger.Setup
	storage storage_crud.NewLinks
}

func (l *LogicLinkAdd) NewLink(ctx context.Context) error {
	userID, ok := context_helper.ContextGetUserID(ctx)
	if !ok {
		return errors.New("user id is required")
	}
	tx, err := l.storage.Transaction.Begin(ctx)
	if err != nil {
		l.Logger.Info("failed start transaction", slog.String("error", err.Error()))
		return err
	}

	timeNow := time.Now()
	defer func() {
		timeoutTime := time.Second * 3
		timeTx := time.Since(timeNow) / 5
		if timeTx < timeoutTime {
			timeoutTime = timeTx
		}
		txRollback, cancel := context.WithTimeout(context.Background(), timeoutTime)
		defer cancel()
		ErrorTx := tx.Rollback(txRollback)
		if ErrorTx != nil && !errors.Is(ErrorTx, pgx.ErrTxClosed) {
			l.Logger.Warn("failed rollback transaction",
				slog.String("error", ErrorTx.Error()))
		}
	}()

	idLink, err := l.storage.Storage.InsertOrFindUrlTx(ctx, tx, l.Data.Url)
	if err != nil {
		_ = slogger.LoggerQueryRow(err, l.Logger, slogger.E1011)
	}

	err = l.storage.Storage.InsertNewUserLinkTx(ctx, tx, userID, idLink)
	if err != nil {
		err = slogger.LoggerExecInsert(err, l.Logger, slogger.E101)
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		l.Logger.Warn("failed commit transaction", slog.String("error", err.Error()))
		return err
	}
	return nil
}

type LogicTakeChart struct {
	storage storage_crud.UserLinks
	logger  slogger.Setup
}

func (l *LogicTakeChart) TakeChart(ctx context.Context) ([]string, error) {
	userID, ok := context_helper.ContextGetUserID(ctx)
	if !ok {
		return nil, errors.New("user id is required")
	}
	links, err := l.storage.Storage.TableUserLinks(ctx, userID)
	if err != nil {
		err = slogger.LoggerQuery(err, l.logger.Logger, slogger.E102)
		return nil, err
	}
	if len(links) == 0 {
		l.logger.Info("no one links found", slog.Any("op", slogger.E102))
		return links, nil
	}
	l.logger.Info("user make request - result finding",
		slog.Int("count", len(links)))

	return links, nil
}

type LogicUpdateLink struct {
	Data request.Change
	slogger.Setup
	storage storage_crud.DBUpdateLink
}

func (l *LogicUpdateLink) ChangeCurrentLink(ctx context.Context) error {
	userID, ok := context_helper.ContextGetUserID(ctx)
	if !ok {
		return errors.New("user id is required")
	}
	g, ctx := errgroup.WithContext(ctx)
	var newID, oldID int
	g.Go(func() error {
		oldid, err := l.storage.FindLink(ctx, l.Data.Url)
		if err != nil {
			err = slogger.LoggerQueryRow(err, l.Logger, slogger.E1042)
			return err
		}
		oldID = oldid
		return nil
	})
	g.Go(func() error {
		newid, err := l.storage.InsertOrFindUrl(ctx, l.Data.NewUrl)
		if err != nil {
			_ = slogger.LoggerQueryRow(err, l.Logger, slogger.E1043)
		}
		newID = newid
		return nil

	})
	if err := g.Wait(); err != nil {
		return err
	}
	data := storage_crud.DataUpdateUserLink{IdUser: userID, IdOldLink: oldID, IdLink: newID}

	err := l.storage.ChangeUserLink(ctx, data)
	if err != nil {
		return slogger.LoggerExecInsert(err, l.Logger, slogger.E1041)
	}
	l.Logger.Info("Update link successfully", slog.Int("old", oldID), slog.Int("new", newID))
	return nil
}

type LogicDeleteLink struct {
	Data request.Delete
	slogger.Setup
	storage storage_crud.DeleteLink
}

func (l *LogicDeleteLink) DeleteLink(ctx context.Context) error {
	userID, ok := context_helper.ContextGetUserID(ctx)
	if !ok {
		return errors.New("user id is required")
	}
	urlID, err := l.storage.Storage.FindLink(ctx, l.Data.Url)
	if err != nil {
		return slogger.LoggerQueryRow(err, l.Logger, slogger.E1031)
	}
	err = l.storage.Storage.DeleteUserLinkAssociation(ctx, urlID, userID)
	if err != nil {
		return slogger.LoggerExecInsert(err, l.Logger, slogger.E103)
	}
	l.Logger.Info("Successful delete Link")
	return nil
}
