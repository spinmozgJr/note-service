package service

import "errors"

var (
	ErrTaskNotFound  = errors.New("задача не найдена")
	ErrTaskForbidden = errors.New("задача недоступна для этого пользователя")
)
