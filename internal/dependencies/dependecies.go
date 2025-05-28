package dependencies

import (
	"github.com/go-playground/validator/v10"
	"github.com/spinmozgJr/note-service/internal/service"
	"log/slog"
)

type Dependencies struct {
	Validator   *validator.Validate
	Log         *slog.Logger
	UserService service.UserService
}

func New(
	validator *validator.Validate,
	log *slog.Logger,
	userService service.UserService) *Dependencies {
	return &Dependencies{
		Validator:   validator,
		Log:         log,
		UserService: userService,
	}
}
