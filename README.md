# DataLinks
type PgxConnect struct {
Pool *pgxpool.Pool
}

func CreatePgxUrl(url PostgresUrl) string{
return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
url.Username,
url.Pass,
url.Hostname,
url.Port,
url.DBName)
}

func CreatePool(connurl string, try int, logger *slog.Logger) (*pgxpool.Pool, error) {
ctx, cancel := WithTimeout(3)
defer cancel()
var pool *pgxpool.Pool
var err error
poolcreate := make(chan struct{})
poolerr := make(chan error)
for i := 0; i < try; i++ {

		logger.Info("Connect pgx try num - ", try)
		go func() {
			p, err := pgxpool.New(ctx, connurl)
			if p == nil {
				pool = p
				poolcreate <- struct{}{}
			}
			if err != nil {
				poolerr <- err
				return
			}
		}()
		if i == try-1 {
			close(poolcreate)
			close(poolerr)
		}
		select {
		case <-ctx.Done():
			logger.Warn("connect time exceeded")
			continue
		case err := <-poolerr:
			logger.Warn("Connect failed", slog.String("error", err.Error()))
			continue
		case <-poolcreate:
			close(poolcreate)
			close(poolerr)
			err = pool.Ping(ctx)
			if err != nil {
				pool.Close()
				logger.Warn("failed ping pool ", err)
				continue
			}
			return pool, nil
		}
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return nil, fmt.Errorf("no one of %d attempts couldn connect in the specified time", try)
	}
	if PgErr, ok := err.(*pgconn.PgError); ok {
		return nil, fmt.Errorf("pgx error: %s, %s", PgErr.Code, PgErr.Message)
	}
		return nil, fmt.Errorf("0 info about connection error")
}