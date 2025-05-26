package postgres

import (
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/spinmozgJr/note-service/internal/config"
	"log"
)

// принимает config.Postgres, так как pgx открывает connection, а для миграций нужна db
func MigrateDB(pg config.Postgres) error {
	migrationsDir := "migrations/"

	pgxConfig, err := pgx.ParseConfig(getConnectionString(pg))
	if err != nil {
		return fmt.Errorf("failed to parse DSN: %v", err)
	}

	db := stdlib.OpenDB(*pgxConfig)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	currentVersion, _ := goose.GetDBVersion(db)
	migrations, err := goose.CollectMigrations(migrationsDir, 0, goose.MaxVersion)

	if err != nil {
		return err
	}

	latestVersion := int64(0)
	for _, m := range migrations {
		if m.Version > latestVersion {
			latestVersion = m.Version
		}
	}

	if currentVersion < latestVersion {
		if err := goose.Up(db, migrationsDir); err != nil {
			return err
		}
	} else {
		log.Println("no migrations required")
	}

	return nil
}
