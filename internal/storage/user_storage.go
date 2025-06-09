package storage

import (
	"context"
	"github.com/spinmozgJr/note-service/internal/models"
)

type UserStorage interface {
	AddUser(ctx context.Context, user, hashPass string) (int, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
}
