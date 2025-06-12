package mapper

import "github.com/spinmozgJr/note-service/internal/models"

func MapNoteDTOFromTaskDb(db models.NoteDB) models.NoteDTO {
	dto := models.NoteDTO{
		ID:        db.ID,
		Title:     db.Title,
		Content:   db.Content,
		CreatedAt: db.CreatedAt,
		UpdateAt:  db.UpdateAt,
	}
	return dto
}
