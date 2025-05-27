package register_user

import (
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/spinmozgJr/note-service/internal/handlers"
	"github.com/spinmozgJr/note-service/internal/httpx"
	"github.com/spinmozgJr/note-service/internal/storage"
	"log/slog"
	"net/http"
)

func New(log *slog.Logger, storage storage.Storage, v *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.register_user.New"

		log = log.With(
			slog.String("op", op),
			//slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var user handlers.RegisterUserInput

		// TODO: возращать более подробные ошибки?
		err := render.DecodeJSON(r.Body, &user)
		if err != nil {
			httpx.SendErrorJSON(w, r, http.StatusBadRequest, err)

			return
		}

		log.Info("request body decoded", slog.Any("request", user))

		if err := v.Struct(user); err != nil {
			validateErr := err.(validator.ValidationErrors)

			httpx.SendErrorJSON(w, r, http.StatusBadRequest, validateErr)

			return
		}

		err = storage.AddUser(r.Context(), user)
		if err != nil {
			log.Error("failed to add user", "error", err)

			render.JSON(w, r, http.StatusInternalServerError)

			return
		}

		render.JSON(w, r, http.StatusCreated)
	}
}
