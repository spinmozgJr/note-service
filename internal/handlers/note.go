package handlers

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/spinmozgJr/note-service/internal/auth"
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

		//userIDStr := chi.URLParam(r, "id")
		//userID, err := strconv.Atoi(userIDStr)
		//if err != nil {
		//	deps.Log.Error("id не является числом", "error", err)
		//	httpx.SendErrorJSON(w, r, http.StatusBadRequest, errors.New("id не является числом"))
		//	return
		//}

		userID, err := auth.GetUserIDFromRequest(r)
		if err != nil {
			deps.Log.Error("пользователь не авторизован", "error", err)
			httpx.SendErrorJSON(w, r, http.StatusUnauthorized, errors.New("пользователь не авторизован"))
			return
		}

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
			deps.Log.Error("ошибка при добавлении заметки", "error", err)

			httpx.SendErrorJSON(w, r, http.StatusInternalServerError, err)
			return
		}

		render.JSON(w, r, response)
	}
}

func DeleteNote(deps *dependencies.Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.note.DeleteNote"
		deps.Log = deps.Log.With(
			slog.String("op", op),
		)

		ctx := r.Context()
		userID, err := auth.GetUserIDFromRequest(r)
		if err != nil {
			deps.Log.Error("пользователь не авторизован", "error", err)
			httpx.SendErrorJSON(w, r, http.StatusUnauthorized, errors.New("пользователь не авторизован"))
			return
		}
		taskIDStr := chi.URLParam(r, "id")
		taskID, err := strconv.Atoi(taskIDStr)
		if err != nil {
			deps.Log.Error("неверный формат идентификатора задачи", "error", err)
			httpx.SendErrorJSON(w, r, http.StatusBadRequest, errors.New("неверный формат идентификатора задачи"))
			return
		}
		response, err := deps.NoteService.DeleteNote(ctx, userID, taskID)
		if err != nil {
			if errors.Is(err, service.ErrTaskNotFound) {
				httpx.SendErrorJSON(w, r, http.StatusNotFound, err)
				return
			}
			httpx.SendErrorJSON(w, r, http.StatusInternalServerError, err)
			return
		}
		render.JSON(w, r, response)
	}
}

//func GetAllNotes(deps *dependencies.Dependencies) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		const op = "handlers.note.GetAllNotes"
//
//		deps.Log = deps.Log.With(
//			slog.String("op", op),
//		)
//
//		ctx := r.Context()
//		userId, err := auth.GetUserIDFromRequest(r)
//		if err != nil {
//			deps.Log.Error("пользователь не авторизован", "error", err)
//			httpx.SendErrorJSON(w, r, http.StatusUnauthorized, errors.New("пользователь не авторизован"))
//			return
//		}
//		response, err := deps.NoteService.GetAllNotes(ctx, userId)
//		if err != nil {
//			httpx.SendErrorJSON(w, r, http.StatusInternalServerError, err)
//			return
//		}
//
//		render.JSON(w, r, response)
//	}
//}

func GetNoteByID(deps *dependencies.Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.note.GetNoteByID"

		deps.Log = deps.Log.With(
			slog.String("op", op),
		)

		ctx := r.Context()
		userId, err := auth.GetUserIDFromRequest(r)
		if err != nil {
			deps.Log.Error("пользователь не авторизован", "error", err)
			httpx.SendErrorJSON(w, r, http.StatusUnauthorized, errors.New("пользователь не авторизован"))
			return
		}
		noteIdStr := chi.URLParam(r, "note_id")
		noteId, err := strconv.Atoi(noteIdStr)
		if err != nil {
			deps.Log.Error("неверный формат идентификатора заметки", "error", err)
			httpx.SendErrorJSON(w, r, http.StatusBadRequest, errors.New("неверный формат идентификатора заметки"))
			return
		}

		response, err := deps.NoteService.GetNoteByID(ctx, userId, noteId)
		if err != nil {
			if errors.Is(err, service.ErrTaskNotFound) {
				deps.Log.Error("заметка не найдена", "error", err)
				httpx.SendErrorJSON(w, r, http.StatusNotFound, err)
				return
			}
			deps.Log.Error("заметка не найдена", "error", err)
			httpx.SendErrorJSON(w, r, http.StatusInternalServerError, err)
			return
		}

		render.JSON(w, r, response)
	}
}

func UpdateNote(deps *dependencies.Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.note.UpdateNote"
		deps.Log = deps.Log.With(
			slog.String("op", op),
		)

		ctx := r.Context()
		userId, err := auth.GetUserIDFromRequest(r)
		if err != nil {
			deps.Log.Error("пользователь не авторизован", "error", err)
			httpx.SendErrorJSON(w, r, http.StatusUnauthorized, errors.New("пользователь не авторизован"))
			return
		}
		noteIdStr := chi.URLParam(r, "note_id")
		noteId, err := strconv.Atoi(noteIdStr)
		if err != nil {
			deps.Log.Error("неверный формат идентификатора заметки", "error", err)
			httpx.SendErrorJSON(w, r, http.StatusBadRequest, errors.New("неверный формат идентификатора заметки"))
			return
		}
		var request UpdateNoteRequest
		if err := httpx.DecodeAndValidateBody(w, r, deps, &request); err != nil {
			deps.Log.Error("ошибка при декодировании тела запроса", "error", err)
			httpx.SendErrorJSON(w, r, http.StatusBadRequest, errors.New("ошибка при декодировании тела запроса"))
			return
		}

		input := &service.UpdateNote{
			ID:      noteId,
			UserID:  userId,
			Title:   request.Title,
			Content: request.Content,
		}
		response, err := deps.NoteService.UpdateNote(ctx, input)
		if err != nil {
			if errors.Is(err, service.ErrTaskNotFound) {
				deps.Log.Error("заметка не найдена", "error", err)
				httpx.SendErrorJSON(w, r, http.StatusNotFound, err)
				return
			}
			deps.Log.Error("ошибка при обновлении заметки", "error", err)
			httpx.SendErrorJSON(w, r, http.StatusInternalServerError, err)
			return
		}
		render.JSON(w, r, response)
	}
}
