package storage

import (
	"context"
)

type NoteStorage interface {
	CreateNote(ctx context.Context, title, content string, userID int) (int, error)
	//GetAllNotes(ctx context.Context) error
	//GetNoteByID(ctx context.Context) error
	//UpdateNote(ctx context.Context) error
	//DeleteNote(ctx context.Context) error
}
