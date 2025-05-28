package register_user

import (
	"github.com/go-chi/render"
	"github.com/spinmozgJr/note-service/internal/dependencies"
	"github.com/spinmozgJr/note-service/internal/handlers"
	"github.com/spinmozgJr/note-service/internal/httpx"
	"log/slog"
	"net/http"
)

func New(deps *dependencies.Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.register_user.New"

		deps.Log = deps.Log.With(
			slog.String("op", op),
			//slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var signIn handlers.SignInRequest
		if err := httpx.DecodeAndValidateBody(w, r, deps, &signIn); err != nil {
			return
		}

		input := handlers.SignInRequest{
			Username: signIn.Username,
			Password: signIn.Password,
		}

		response, err := deps.UserService.SignIn(r.Context(), input)
		if err != nil {
			deps.Log.Error("failed to add user", "error", err)

			render.JSON(w, r, http.StatusInternalServerError)

			return
		}

		render.JSON(w, r, response)
	}
}
