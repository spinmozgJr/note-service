package handlers

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/spinmozgJr/note-service/internal/dependencies"
	"github.com/spinmozgJr/note-service/internal/httpx"
	"github.com/spinmozgJr/note-service/internal/service"
	"log/slog"
	"net/http"
	"strconv"
)

func CreateNote(deps *dependencies.Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.note.CreateNote"

		deps.Log = deps.Log.With(
			slog.String("op", op),
		)

		userIDStr := chi.URLParam(r, "id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			deps.Log.Error("id не является числом", "error", err)
			httpx.SendErrorJSON(w, r, http.StatusBadRequest, errors.New("id не является числом"))
			return
		}

		//userID, err := auth.GetUserIDFromRequest(r)
		//if err != nil {
		//	deps.Log.Error("unknown user id", "error", err)
		//	httpx.SendErrorJSON(w, r, http.StatusUnauthorized, errors.New("пользователь не авторизован"))
		//	return
		//}

		var inputNote InputNote
		if err := httpx.DecodeAndValidateBody(w, r, deps, &inputNote); err != nil {
			return
		}

		input := &service.NoteInput{
			Title:   inputNote.Title,
			Content: inputNote.Content,
			UserID:  userID,
		}

		response, err := deps.NoteService.CreateNote(r.Context(), input)
		if err != nil {
			deps.Log.Error("failed to add note", "error", err)

			httpx.SendErrorJSON(w, r, http.StatusInternalServerError, err)
			return
		}

		render.JSON(w, r, response)
	}
}
