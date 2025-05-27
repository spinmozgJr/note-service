package handlers

type RegisterUserInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password"`
}
