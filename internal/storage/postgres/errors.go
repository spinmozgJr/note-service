package postgres

import "errors"

var (
	ErrSelect       = errors.New("ошибка выборки из базы данных")
	ErrCreateUser   = errors.New("ошибка создания пользователя")
	ErrCreateNote   = errors.New("ошибка создания заметки")
	ErrUpdateNote   = errors.New("ошибка обновления заметки")
	ErrDeleteNote   = errors.New("ошибка удаления заметки")
	ErrNoteNotFound = errors.New("заметка не найдена")
)
