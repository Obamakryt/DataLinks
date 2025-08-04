package storage_crud

import (
	"DataLinks/internal/slogger"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

func (p *PostgresPool) InsertOrFindUrlTx(ctx context.Context, tx pgx.Tx, url string) (int, error) {
	var UrlId int
	const query = `INSERT INTO urls (unique_url) 
              VALUES ($1) 
              ON CONFLICT (unique_url) 
              DO UPDATE SET 
              unique_url = urls.unique_url 
              RETURNING id;`

	err := tx.QueryRow(ctx, query, url).Scan(&UrlId)

	return UrlId, err
}
func (p *PostgresPool) InsertNewUserLinkTx(ctx context.Context, tx pgx.Tx, idUser int, idLink int) error {
	const query = `INSERT INTO user_links (user_id,urls_id) 
			  VALUES ($1,$2) 
			  ON CONFLICT (user_id,urls_id) 
			  DO NOTHING;`
	commandTag, err := tx.Exec(ctx, query, idUser, idLink)

	return slogger.Exec(err, commandTag)

}

func (p *PostgresPool) Begin(ctx context.Context) (pgx.Tx, error) {
	return p.client.Pool.BeginTx(ctx, pgx.TxOptions{})
}

func (p *PostgresPool) TableUserLinks(ctx context.Context, idUser int) ([]string, error) {
	var chartLinks []string
	const query = `SELECT urls.unique_url
				   FROM user_links
				   JOIN urls ON user_links.urls_id = urls.id
				   WHERE user_links.user_id = $1;`

	rows, err := p.client.Pool.Query(ctx, query, idUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var link string
		if err = rows.Scan(&link); err != nil {
			return nil, fmt.Errorf("%w, %v", slogger.ScanError, err)
		}
		chartLinks = append(chartLinks, link)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%w, %v", slogger.RowsError, err)
	}
	return chartLinks, nil
}

func (p *PostgresPool) InsertOrFindUrl(ctx context.Context, url string) (int, error) {
	var urlLink int
	const query = `INSERT INTO urls (unique_url) 
              VALUES ($1) 
              ON CONFLICT (unique_url) 
              DO UPDATE SET 
              unique_url = urls.unique_url 
              RETURNING id;`

	err := p.client.Pool.QueryRow(ctx, query, url).Scan(&urlLink)

	return urlLink, err
}

func (p *PostgresPool) FindLink(ctx context.Context, url string) (int, error) {
	var oldLink int
	const query = `SELECT id FROM urls WHERE unique_url = $1;`

	err := p.client.Pool.QueryRow(ctx, query, url).Scan(&oldLink)
	if err != nil {
		return -1, err
	}
	return oldLink, nil
}

func (p *PostgresPool) ChangeUserLink(ctx context.Context, data DataUpdateUserLink) error {
	const query = `UPDATE user_links
				   SET urls_id = $1 
				   WHERE user_id = $2 
				   AND urls_id = $3;`

	commandTag, err := p.client.Pool.Exec(ctx, query, data.IdLink, data.IdUser, data.IdOldLink)

	return slogger.Exec(err, commandTag)
}

func (p *PostgresPool) DeleteUserLinkAssociation(ctx context.Context, urlID int, userid int) error {
	const query = `DELETE FROM user_links
       			   WHERE user_id = $1
       			   AND urls_id = $2;`

	commandTag, err := p.client.Pool.Exec(ctx, query, userid, urlID)

	return slogger.Exec(err, commandTag)

}
