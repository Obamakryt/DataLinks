package storage_crud

import (
	"DataLinks/internal/slogger"
	"context"
)

func (p *PostgresPool) Authorization(ctx context.Context, email string) (StorageAuth, error) {
	const query = `SELECT password, id
		 		   FROM users
  		  		   WHERE email=$1`
	data := StorageAuth{}
	err := p.Pool.QueryRow(ctx, query, email).Scan(&data.Password, &data.Id)

	if err != nil {
		return data, err
	}
	return data, nil
}

func (p *PostgresPool) Registration(ctx context.Context, r *StorageRegister) error {
	q := `INSERT INTO users(name, email, password)
		  VALUES($1, $2, $3)`

	commandTag, err := p.Pool.Exec(ctx, q, r.Name, r.Email, r.HashPass)

	return slogger.Exec(err, commandTag)
}
