package postgres

import (
	"context"
	"time"
)

type NoteStorage struct {
	Repository
}

func NewNoteStorage(connector *DBConnector) NoteStorage {
	return NoteStorage{Repository{Conn: connector.Conn}}
}

func (s NoteStorage) CreateNote(ctx context.Context, title, content string, userID int) (int, error) {
	query := `
			INSERT INTO notes (user_id, title, content, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id;
	`
	var taskId int
	err := s.Conn.QueryRow(ctx, query, userID, title, content, time.Now(), time.Now()).Scan(&taskId)
	if err != nil {
		return -1, ErrCreateNote
	}
	return taskId, nil
}
