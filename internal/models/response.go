package models

type BaseResponse struct {
	Error *BaseError `json:"error,omitempty"`
}

type BaseError struct {
	Message string `json:"message,omitempty"`
}

type AuthResponse struct {
	BaseResponse
	Data *AuthData `data:"error,omitempty"`
}

type OperationResultResponse struct {
	BaseResponse
	Data *OperationResultData `data:"error,omitempty"`
}

type OperationResultData struct {
	Success bool `json:"success"`
}

type AuthData struct {
	Username    string `json:"username"`
	AccessToken string `json:"access_token"`
}

type NoteListResponse struct {
	BaseResponse
	Data *[]NoteDTO `data:"error,omitempty"`
}

type NoteResponse struct {
	BaseResponse
	Data *NoteDTO `data:"error,omitempty"`
}

type CreateNoteResponse struct {
	BaseResponse
	Data *CreateNoteData `json:"data,omitempty"`
}

type CreateNoteData struct {
	NoteID int `json:"note_id"`
}
