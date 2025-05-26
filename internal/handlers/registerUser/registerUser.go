package registerUser

import (
	"context"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/spinmozgJr/note-service/internal/models"
	"log/slog"
	"net/http"
	"strings"
)

type RegisterUser interface {
	AddUser(ctx context.Context, user models.User) error
}

func New(ctx context.Context, log *slog.Logger, registerUser RegisterUser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.registerUser.New"

		log = log.With(
			slog.String("op", op),
			//slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var user models.User

		// TODO: возращать более подробные ошибки?
		err := render.DecodeJSON(r.Body, &user)
		if err != nil {
			log.Error("failed to decode request body", "error", err.Error())

			render.JSON(w, r, http.StatusBadRequest)

			return
		}

		log.Info("request body decoded", slog.Any("request", user))

		if err := validator.New().Struct(user); err != nil {
			//validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", "error", err.Error())

			// было бы неплохо как-то вернуть validateErr
			render.JSON(w, r, http.StatusBadRequest)

			return
		}

		err = registerUser.AddUser(ctx, user)
		if err != nil {
			log.Error("failed to add user", "error", err)

			render.JSON(w, r, http.StatusInternalServerError)

			return
		}

		if strings.TrimSpace(user.Username) == "" {
			render.JSON(w, r, http.StatusBadRequest)

			return
		}

		render.JSON(w, r, http.StatusCreated)
	}
}
