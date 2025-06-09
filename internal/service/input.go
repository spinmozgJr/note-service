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
