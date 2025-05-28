package models

type BaseResponse[T any] struct {
	Data  *T         `json:"data"`
	Error *BaseError `json:"error"`
}

type BaseError struct {
	Message string `json:"message"`
}

type AuthData struct {
	Username    string `json:"username"`
	AccessToken string `json:"access_token"`
}
