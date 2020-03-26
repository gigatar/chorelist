package gigatarerrors

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrDuplicate  = errors.New("duplicate")
	ErrBadRequest = errors.New("bad request")
)

func WebErrorCodes(err error) int {
	switch err {

	case ErrBadRequest:
		return http.StatusBadRequest
	case ErrDuplicate:
		return http.StatusConflict
	case ErrNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}

}
