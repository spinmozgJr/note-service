package dependencies

import (
	"github.com/go-playground/validator/v10"
	"github.com/spinmozgJr/note-service/internal/service"
	"github.com/spinmozgJr/note-service/pkg/auth"
	"log/slog"
)

type Dependencies struct {
	Validator    *validator.Validate
	Log          *slog.Logger
	UserService  service.UserService
	NoteService  service.NoteService
	TokenManager auth.TokenManager
}

func New(
	validator *validator.Validate,
	log *slog.Logger,
	userService service.UserService,
	noteService service.NoteService,
	tokenManager auth.TokenManager) *Dependencies {
	return &Dependencies{
		Validator:    validator,
		Log:          log,
		UserService:  userService,
		NoteService:  noteService,
		TokenManager: tokenManager,
	}
}
