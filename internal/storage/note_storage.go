package storage

import (
	"context"
	"github.com/spinmozgJr/note-service/internal/models"
)

type NoteStorage interface {
	CreateNote(ctx context.Context, title, content string, userID int) (int, error)
	GetNoteByID(ctx context.Context, noteId int) (*models.NoteDB, error)
	//GetAllNotes(ctx context.Context) error
	UpdateNote(ctx context.Context, input *UpdateNote) error
	DeleteNoteByID(ctx context.Context, noteID int) error
}
