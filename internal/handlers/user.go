package handlers

import (
	"github.com/go-chi/render"
	"github.com/spinmozgJr/note-service/internal/dependencies"
	"github.com/spinmozgJr/note-service/internal/httpx"
	"github.com/spinmozgJr/note-service/internal/service"
	"log/slog"
	"net/http"
)

func Registration(deps *dependencies.Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.Registration"

		deps.Log = deps.Log.With(
			slog.String("op", op),
			//slog.String("request_id", middlewares.GetReqID(r.Context())),
		)

		var signIn SignInRequest
		if err := httpx.DecodeAndValidateBody(w, r, deps, &signIn); err != nil {
			return
		}

		input := service.UserInput{
			Username: signIn.Username,
			Password: signIn.Password,
		}

		response, err := deps.UserService.SignIn(r.Context(), input)
		if err != nil {
			deps.Log.Error("failed to add user", "error", err)

			//render.JSON(w, r, http.StatusInternalServerError)
			httpx.SendErrorJSON(w, r, http.StatusInternalServerError, err)

			return
		}

		render.JSON(w, r, response)
	}
}

func Login(deps *dependencies.Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.Login"

		deps.Log = deps.Log.With(
			slog.String("op", op),
			//slog.String("request_id", middlewares.GetReqID(r.Context())),
		)

		var signIn SignInRequest
		if err := httpx.DecodeAndValidateBody(w, r, deps, &signIn); err != nil {
			return
		}

		input := service.UserInput{
			Username: signIn.Username,
			Password: signIn.Password,
		}

		response, err := deps.UserService.Login(r.Context(), input)
		if err != nil {
			deps.Log.Error("failed to login user", "error", err)

			//render.JSON(w, r, http.StatusInternalServerError)
			httpx.SendErrorJSON(w, r, http.StatusInternalServerError, err)

			return
		}

		render.JSON(w, r, response)
	}
}
