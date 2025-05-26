package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/spinmozgJr/note-service/internal/config"
	"github.com/spinmozgJr/note-service/internal/models"
	"time"
)

type Storage struct {
	conn *pgx.Conn
}

func New(ctx context.Context, pg config.Postgres) (*Storage, error) {
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

	return &Storage{conn: conn}, nil
}

func (s *Storage) Close() {
	s.conn.Close(context.Background())
}

func (s *Storage) AddUser(ctx context.Context, user models.User) error {
	const op = "storage.postgres.AddUser"

	query := `INSERT INTO users (username, created_at) VALUES ($1, $2)`

	_, err := s.conn.Exec(ctx, query, user.Username, time.Now())
	return err
}

func getConnectionString(p config.Postgres) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		p.Host, p.Port, p.User, p.Password, p.DBName, p.SSLMode)
}
