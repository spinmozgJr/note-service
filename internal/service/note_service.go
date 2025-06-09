package service

import (
	"context"
	"github.com/spinmozgJr/note-service/internal/models"
	"github.com/spinmozgJr/note-service/internal/storage"
)

type NoteService struct {
	NoteRepository storage.NoteStorage
}

func NewNoteService(noteRepository storage.NoteStorage) NoteService {
	return NoteService{
		NoteRepository: noteRepository,
	}
}

func (s *NoteService) CreateNote(ctx context.Context, note *NoteInput) (*models.BaseResponse[int], error) {
	noteID, err := s.NoteRepository.CreateNote(ctx, note.Title, note.Content, note.UserID)
	if err != nil {
		return nil, err
	}

	response := &models.BaseResponse[int]{
		Data: &noteID,
	}
	return response, nil
}
