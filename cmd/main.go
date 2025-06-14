package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/spinmozgJr/note-service/internal/config"
	"github.com/spinmozgJr/note-service/internal/dependencies"
	"github.com/spinmozgJr/note-service/internal/handlers"
	mw "github.com/spinmozgJr/note-service/internal/middlewares"
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
	connector, err := postgres.NewDBConnector(ctx, cfg.Postgres)
	if err != nil {
		slog.Error("failed to init connector", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := connector.Close(); err != nil {
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

	userStorage := postgres.NewUserStorage(connector)
	userService := service.NewUserService(userStorage, tokenManager, cfg)

	noteStorage := postgres.NewNoteStorage(connector)
	noteService := service.NewNoteService(noteStorage)

	deps := dependencies.New(v, log, userService, noteService, tokenManager)

	router := chi.NewRouter()

	router.Use(mw.NewLoggerMiddleware(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/", func(r chi.Router) {
		//r.Post("/sign-in", h.PostSignIn)
		//r.Post("/login", h.PostLogin)
		router.Post("/users", handlers.Registration(deps))
		router.Post("/login", handlers.Login(deps))
	})

	router.With(mw.JwtAuthMiddleware(deps)).Group(func(r chi.Router) {
		r.Post("/users/{id}/notes", handlers.CreateNote(deps))
		r.Get("/users/{id}/notes/{note_id}", handlers.GetNoteByID(deps))
		r.Get("/users/{id}/notes", handlers.GetAllNotes(deps))
		r.Put("/users/{id}/notes/{note_id}", handlers.UpdateNote(deps))
		r.Delete("/users/{id}/notes/{note_id}", handlers.DeleteNote(deps))
	})

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
