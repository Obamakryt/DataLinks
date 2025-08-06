package service

import (
	"DataLinks/internal/dto/request"
	"DataLinks/internal/slogger"
	"DataLinks/internal/storages/postgreSQL/storage_crud"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/gommon/log"
	"log/slog"
	"sync"
	"time"
)

type LogicLinkAdd struct {
	Logger  slogger.Setup
	Storage storage_crud.NewLinks
}

func (l *LogicLinkAdd) NewLink(ctx context.Context, Data request.Add) error {
	tx, err := l.Storage.Transaction.Begin(ctx)
	if err != nil {
		l.Logger.Info("failed start transaction", slog.String("error", err.Error()))
		return err
	}

	defer func() {
		timeoutTime := time.Second * 3
		txRollback, cancel := context.WithTimeout(context.Background(), timeoutTime)
		defer cancel()
		ErrorTx := tx.Rollback(txRollback)
		if ErrorTx != nil && !errors.Is(ErrorTx, pgx.ErrTxClosed) {
			l.Logger.Warn("failed rollback transaction",
				slog.String("error", ErrorTx.Error()))
		}
	}()

	idLink, err := l.Storage.InsertOrFindUrlTx(ctx, tx, Data.Url)
	if err != nil {
		_ = slogger.LoggerQueryRow(err, l.Logger, slogger.E1011)
	}

	err = l.Storage.InsertNewUserLinkTx(ctx, tx, Data.UserId, idLink)
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
	Storage storage_crud.UserLinks
	Logger  slogger.Setup
}

func (l *LogicTakeChart) TakeChart(ctx context.Context, Data request.TakeChart) ([]string, error) {
	links, err := l.Storage.TableUserLinks(ctx, Data.UserId)
	if err != nil {
		err = slogger.LoggerQuery(err, l.Logger, slogger.E102)
		return nil, err
	}
	if len(links) == 0 {
		l.Logger.Info("no one links found", slog.Any("op", slogger.E102))
		return links, nil
	}
	l.Logger.Info("user make request - result finding",
		slog.Int("count", len(links)))

	return links, nil
}

type LogicUpdateLink struct {
	Logger  slogger.Setup
	Storage storage_crud.UpdateLink
}

func (l *LogicUpdateLink) ChangeCurrentLink(ctx context.Context, Data request.Swap) error {
	var wg sync.WaitGroup
	var newID, oldID int
	errchan := make(chan error, 1)
	wg.Add(2)
	go func() {
		defer wg.Done()
		log.Info(Data.Url)

		oldId, err := l.Storage.FindLink(ctx, Data.Url)
		if err != nil {
			err = slogger.LoggerQueryRow(err, l.Logger, slogger.E1041)
			errchan <- err
		}
		oldID = oldId

		return
	}()
	go func() {
		defer wg.Done()

		newId, err := l.Storage.InsertOrFindUrl(ctx, Data.NewUrl)
		if err != nil {
			_ = slogger.LoggerQueryRow(err, l.Logger, slogger.E1042)
		}
		newID = newId
		return

	}()
	wg.Wait()
	close(errchan)

	if err := ctx.Err(); err != nil {
		log.Info(err.Error())

		return slogger.DBError
	}

	err, ok := <-errchan
	if ok {
		return err
	}

	err = l.Storage.ChangeUserLink(ctx,
		storage_crud.DataUpdateUserLink{IdUser: Data.UserId, IdOldLink: oldID, IdLink: newID})

	if err != nil {
		return slogger.LoggerExecInsert(err, l.Logger, slogger.E104)
	}

	l.Logger.Info("Update link successfully", slog.Int("old", oldID), slog.Int("new", newID))
	return nil

}

type LogicDeleteLink struct {
	Logger  slogger.Setup
	Storage storage_crud.DeleteLink
}

func (l *LogicDeleteLink) DeleteLink(ctx context.Context, Data request.Delete) error {
	urlID, err := l.Storage.FindLink(ctx, Data.Url)
	if err != nil {
		return slogger.LoggerQueryRow(err, l.Logger, slogger.E1031)
	}
	err = l.Storage.DeleteUserLinkAssociation(ctx, urlID, Data.UserId)
	if err != nil {
		return slogger.LoggerExecInsert(err, l.Logger, slogger.E103)
	}
	l.Logger.Info("Successful delete Link")
	return nil
}
