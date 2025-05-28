package httpx

import (
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/spinmozgJr/note-service/internal/dependencies"
	"log/slog"
	"net/http"
)

func DecodeAndValidateBody(w http.ResponseWriter,
	r *http.Request,
	deps *dependencies.Dependencies,
	dst interface{}) error {
	err := render.DecodeJSON(r.Body, &dst)
	if err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err)

		return err
	}

	deps.Log.Info("request body decoded", slog.Any("request", dst))

	if err := deps.Validator.Struct(dst); err != nil {
		validateErr := err.(validator.ValidationErrors)

		SendErrorJSON(w, r, http.StatusBadRequest, validateErr)

		return err
	}

	return nil
}
