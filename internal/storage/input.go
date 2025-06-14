package storage

type AddUserInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateNote struct {
	ID      int
	Title   string
	Content string
}

type InputGetAllNotes struct {
	UserID int
	Limit  int
	Offset int
	Sort   string
}
