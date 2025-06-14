package service

type UserInput struct {
	Username string
	Password string
}

type NoteInput struct {
	UserID  int
	Title   string
	Content string
}

type UpdateNote struct {
	ID      int
	UserID  int
	Title   string
	Content string
}

type ServiceQueryParams struct {
	ID     int
	Limit  int
	Offset int
	Sort   string
}
