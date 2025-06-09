package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/spinmozgJr/note-service/internal/models"
	"time"
)

type UserStorage struct {
	Repository
}

func NewUserStorage(connector *DBConnector) *UserStorage {
	return &UserStorage{Repository{Conn: connector.Conn}}
}

func (s *UserStorage) AddUser(ctx context.Context, user, hashPass string) (int, error) {
	const op = "storage.postgres.AddUser"

	query := `INSERT INTO users (username, password, created_at) VALUES ($1, $2, $3) RETURNING id`

	userId := 0
	err := s.Repository.Conn.QueryRow(ctx, query, user, hashPass, time.Now()).Scan(&userId)
	return userId, err
}

func (s *UserStorage) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	const op = "storage.postgres.GetUserByUserName"

	var user *models.User

	query := `
				SELECT id, username, password FROM users
				WHERE username = $1;
	`

	row := s.Repository.Conn.QueryRow(ctx, query, username)
	user, err := scanUserRow(row)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func scanUserRow(row pgx.Row) (*models.User, error) {
	var user models.User

	if err := row.Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
	); err != nil {
		return nil, err
	}

	return &user, nil
}
