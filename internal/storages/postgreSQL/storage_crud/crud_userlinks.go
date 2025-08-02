package storage_crud

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

func (p *PostgresPool) InsertOrFindUrl(ctx context.Context, tx pgx.Tx, url string) (int, error) {
	var UrlId int
	const query = `INSERT INTO urls (unique_url) 
              VALUES ($1) 
              ON CONFLICT (unique_url) 
              DO UPDATE SET 
              unique_url = EXCLUDED.unique_url 
              RETURNING id;`

	err := tx.QueryRow(ctx, query, url).Scan(&UrlId)

	return UrlId, err
}
func (p *PostgresPool) InsertNewLink(ctx context.Context, tx pgx.Tx, idUser int, idLink int) (bool, error) {
	const query = `INSERT INTO user_links (user_id,urls_id) 
			  VALUES ($1,$2) 
			  ON CONFLICT (user_id,urls_id) 
			  DO NOTHING;`
	commandTag, err := tx.Exec(ctx, query, idUser, idLink)
	if err != nil {
		return false, err
	}
	if commandTag.RowsAffected() == 0 {
		return true, err
	}
	return true, nil
}

func (p *PostgresPool) Begin(ctx context.Context) (pgx.Tx, error) {
	return p.Client.Pool.BeginTx(ctx, pgx.TxOptions{})
}

func (p *PostgresPool) TableUserLinks(ctx context.Context, idUser int) ([]string, error) {
	chartLinks := []string{}
	const query = `SELECT urls.unique_url
				   FROM user_links
				   JOIN urls ON user_links.urls_id = urls.id
				   WHERE user_links.user_id = $1;`

	rows, err := p.Client.Pool.Query(ctx, query, idUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var link string
		if err = rows.Scan(&link); err != nil {
			return nil, fmt.Errorf("failed to scan link: %w", err)
		}
		chartLinks = append(chartLinks, link)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return chartLinks, nil
}

func (p *PostgresPool) ChangeLink(ctx context.Context, url string) (int, error) {
	var urlLink int
	const query = `INSERT INTO urls (unique_url)
                   VALUES ($1)
                   ON CONFLICT (unique_url)
                   DO NOTHING
                   RETURNING id;`
	err := p.Client.Pool.QueryRow(ctx, query, url).Scan(&urlLink)
	if err != nil {
		return -1, err
	}
	return urlLink, nil
}

func (p *PostgresPool) FindOldLink(ctx context.Context, url string) (int, error) {
	var oldLink int
	const query = `SELECT id FROM urls WHERE unique_url = $1;`

	err := p.Client.Pool.QueryRow(ctx, query, url).Scan(&oldLink)
	if err != nil {
		return -1, err
	}
	return oldLink, nil
}

func (p *PostgresPool) UpdateUserLink(ctx context.Context, data DataUpdateUserLink) error {
	const query = `UPDATE user_links
				   SET urls_id = $1 
				   WHERE user_id = $2 
				   AND urls_id = $3;`

	_, err := p.Client.Pool.Exec(ctx, query, data.IdLink, data.IdUser, data.IdOldLink)
	if err != nil {
		return err
	}
	return nil
}
