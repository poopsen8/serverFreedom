package httperr

import (
	"errors"
	"net/http"
)

type ErrorHttp struct {
	Err           error
	StatusRequest int
}

var ErrInvalidID = ErrorHttp{
	Err:           errors.New("invalid id"),
	StatusRequest: http.StatusBadRequest, // 400
}

var ErrInvalidJSON = ErrorHttp{
	Err:           errors.New("invalid JSON"),
	StatusRequest: http.StatusBadRequest, // 400
}

var ErrIDNotFound = ErrorHttp{
	Err:           errors.New("ID not found"),
	StatusRequest: http.StatusNotFound, // 404
}

func ValidateUserID(id int64) ErrorHttp {
	if id <= 0 {
		return ErrInvalidID
	}
	return ErrorHttp{nil, 0}
}
