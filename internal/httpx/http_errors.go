package httpx

import (
	"github.com/go-chi/render"
	"github.com/spinmozgJr/note-service/internal/models"
	"net/http"
)

func SendErrorJSON(w http.ResponseWriter, r *http.Request, httpStatus int, err error) {
	render.Status(r, httpStatus)

	response := models.BaseResponse{
		Error: &models.BaseError{
			Message: err.Error(),
		},
	}

	render.JSON(w, r, response)
}
