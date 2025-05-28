package storage

import (
	"context"
)

type Storage interface {
	AddUser(ctx context.Context, user, hashPass string) (int, error)
	Close() error
}
