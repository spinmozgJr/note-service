package models

import "time"

type NoteDTO struct {
	ID        int
	Title     string
	Content   string
	CreatedAt time.Time
	UpdateAt  time.Time
}
