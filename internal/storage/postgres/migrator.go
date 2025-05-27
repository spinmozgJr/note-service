package postgres

import (
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/spinmozgJr/note-service/internal/config"
)

// принимает config.Postgres, так как pgx открывает connection, а для миграций нужна db
func MigrateDB(pg config.Postgres) error {
	migrationsDir := "internal/storage/migrations/"

	pgxConfig, err := pgx.ParseConfig(getConnectionString(pg))
	if err != nil {
		return fmt.Errorf("failed to parse DSN: %w", err)
	}

	db := stdlib.OpenDB(*pgxConfig)
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		return err
	}

	return nil
}
