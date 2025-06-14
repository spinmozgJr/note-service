package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/spinmozgJr/note-service/internal/models"
	"github.com/spinmozgJr/note-service/internal/storage"
	"time"
)

type NoteStorage struct {
	Repository
}

func NewNoteStorage(connector *DBConnector) storage.NoteStorage {
	return &NoteStorage{Repository{Conn: connector.Conn}}
}

func (s *NoteStorage) CreateNote(ctx context.Context, title, content string, userID int) (int, error) {
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

func (s *NoteStorage) DeleteNoteByID(ctx context.Context, noteID int) error {
	query := `
		DELETE FROM notes
		WHERE id = $1;
	`
	tag, err := s.Conn.Exec(ctx, query, noteID)
	if err != nil {
		return ErrDeleteNote
	}
	if tag.RowsAffected() == 0 {
		return ErrNoteNotFound
	}

	return nil
}

func (s *NoteStorage) GetAllNotes(ctx context.Context, input *storage.InputGetAllNotes) ([]models.NoteDB, error) {
	query := fmt.Sprintf(`
        SELECT id, user_id, title, content, created_at, updated_at 
        FROM notes
        WHERE user_id = $1
        ORDER BY created_at %s
        LIMIT $2
        OFFSET $3`, input.Sort)

	rows, err := s.Conn.Query(ctx, query, input.UserID, input.Limit, input.Offset)
	if err != nil {
		return nil, ErrSelect
	}
	tasks, err := scanNoteListRow(rows)
	if err != nil {
		return nil, ErrSelect
	}
	return tasks, nil
}

func (s *NoteStorage) GetNoteByID(ctx context.Context, noteId int) (*models.NoteDB, error) {
	query := `
			SELECT id, user_id, title, content, created_at, updated_at 
			FROM notes
			WHERE id = $1;
	`

	row := s.Conn.QueryRow(ctx, query, noteId)
	note, err := scanNoteRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoteNotFound
		}
		return nil, ErrSelect
	}

	return note, nil
}

func (s *NoteStorage) UpdateNote(ctx context.Context, input *storage.UpdateNote) error {
	query := `
		UPDATE notes
		SET title = $1,
			content = $2,
			updated_at = $3
		WHERE id = $4;
	`
	tag, err := s.Conn.Exec(ctx, query, input.Title, input.Content, time.Now(), input.ID)
	if err != nil {
		return ErrUpdateNote
	}
	rows := tag.RowsAffected()
	if rows == 0 {
		return ErrNoteNotFound
	}

	return nil
}

func scanNoteRow(row pgx.Row) (*models.NoteDB, error) {
	var note models.NoteDB

	if err := row.Scan(
		&note.ID,
		&note.UserID,
		&note.Title,
		&note.Content,
		&note.CreatedAt,
		&note.UpdateAt,
	); err != nil {
		return nil, err
	}

	return &note, nil
}

func scanNoteListRow(rows pgx.Rows) ([]models.NoteDB, error) {
	var items []models.NoteDB
	for rows.Next() {
		var task models.NoteDB

		err := rows.Scan(
			&task.ID,
			&task.UserID,
			&task.Title,
			&task.Content,
			&task.CreatedAt,
			&task.UpdateAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, task)
	}
	return items, nil
}
