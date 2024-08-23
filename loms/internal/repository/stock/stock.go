package stock

import "github.com/jackc/pgx/v5"

type RepositoryStock struct {
	psqlConnection *pgx.Conn
}

func New(psqlConnection *pgx.Conn) *RepositoryStock {
	return &RepositoryStock{
		psqlConnection: psqlConnection,
	}
}
