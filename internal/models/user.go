package models

type User struct {
	ID           int    `json:"id,omitempty"`
	Username     string `json:"username" validate:"required"`
	PasswordHash string `json:"password"`
}
