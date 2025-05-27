package models

import "time"

type User struct {
	ID        int       `json:"id,omitempty"`
	Username  string    `json:"username" validate:"required"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Password  string    `json:"password"`
}
