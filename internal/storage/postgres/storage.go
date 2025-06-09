package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/spinmozgJr/note-service/internal/config"
)

type DBConnector struct {
	Conn *pgx.Conn
}

func NewDBConnector(ctx context.Context, pg config.Postgres) (*DBConnector, error) {
	const op = "storage.postgres.New"

	connStr := getConnectionString(pg)

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s не удалось проверить подключение: %v", op, err)
	}

	return &DBConnector{Conn: conn}, nil
}

func getConnectionString(p config.Postgres) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		p.Host, p.Port, p.User, p.Password, p.DBName, p.SSLMode)
}

func (c *DBConnector) Close() error {
	err := c.Conn.Close(context.Background())
	if err != nil {
		return err
	}
	return nil
}
