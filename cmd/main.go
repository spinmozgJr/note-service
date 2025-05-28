package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/spinmozgJr/note-service/internal/config"
	"github.com/spinmozgJr/note-service/internal/dependencies"
	"github.com/spinmozgJr/note-service/internal/handlers/register_user"
	mwLogger "github.com/spinmozgJr/note-service/internal/middleware"
	"github.com/spinmozgJr/note-service/internal/service"
	"github.com/spinmozgJr/note-service/internal/storage/postgres"
	"github.com/spinmozgJr/note-service/pkg/auth"
	"github.com/spinmozgJr/note-service/pkg/logger"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	cfg := config.MustLoad()

	log := logger.New(cfg.Env)

	log.Info("starting note-service", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	ctx := context.Background()
	storage, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		slog.Error("failed to init storage", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := storage.Close(); err != nil {
			log.Info("close error: %w", err)
		}
	}()

	err = postgres.MigrateDB(cfg.Postgres)
	if err != nil {
		slog.Error("failed to migrate db", "error", err)
		os.Exit(1)
	}

	v := validator.New()
	tokenManager, err := auth.NewManager(cfg.SigningKey)
	if err != nil {
		slog.Error("failed to init token manager", "error", err)
	}
	userService := service.NewUserService(storage, tokenManager, cfg)

	deps := dependencies.New(v, log, userService)

	router := chi.NewRouter()

	router.Use(mwLogger.NewLoggerMiddleware(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/users", register_user.New(deps))

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.TimeOut,
		WriteTimeout: cfg.HTTPServer.TimeOut,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}
