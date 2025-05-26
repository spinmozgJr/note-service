package main

import (
	"github.com/spinmozgJr/note-service/internal/config"
	"github.com/spinmozgJr/note-service/internal/storage/postgres"
	"github.com/spinmozgJr/note-service/pkg/logger"
	"log/slog"
	"os"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	log.Info("starting note-service", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := postgres.New(cfg.Postgres)
	if err != nil {
		slog.Error("failed to init storage", "error", err)
		os.Exit(1)
	}
	defer storage.Close()

	err = postgres.MigrateDB(cfg.Postgres)
	if err != nil {
		slog.Error("failed to migrate db", "error", err)
		os.Exit(1)
	}
}
