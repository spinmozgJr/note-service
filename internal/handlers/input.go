package handlers

type SignInRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type InputNote struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type QueryParams struct {
	Limit  int
	Offset int
	Sort   string
}

type UpdateNoteRequest struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}
