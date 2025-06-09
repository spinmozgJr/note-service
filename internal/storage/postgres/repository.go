package postgres

import "github.com/jackc/pgx/v5"

type Repository struct {
	Conn *pgx.Conn
}
