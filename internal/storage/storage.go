package storage

import (
	"context"
	"github.com/spinmozgJr/note-service/internal/handlers"
)

type Storage interface {
	AddUser(ctx context.Context, user handlers.RegisterUserInput) error
	Close() error
}
