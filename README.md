# DataLinks
func CreatePool(connurl string, try int, logger *slog.Logger) (*pgxpool.Pool, error) {


	var err error
	for i := 0; i < try; i++ {
		ctx, cancel := WithTimeout(3)
		defer cancel()
		logger.Info("Connect pgx try num - ", try)
		pool, e := pgxpool.New(ctx, connurl)
		if e != nil {
			err = e
			logger.Warn("Connect failed", slog.String("error", err.Error()))
			continue
		}else{
			e = pool.Ping(ctx)
			
			if e != nil{
				pool.Close()
				logger.Warn("failed ping pool connect ", err)
				err = e
				continue
			}
			
				return pool,  nil
			}
		}
		
	var pgErr *pgconn.PgError
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return nil, fmt.Errorf("no one of %d attempts couldn connect in the specified time", try)
	case errors.As(err, &pgErr):
		return nil, fmt.Errorf("pgx error: %s, %s", pgErr.Code, pgErr.Message)
	default:
		return nil, fmt.Errorf("0 info about connection error")

	}

}