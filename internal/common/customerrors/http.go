package customerrors

import (
	"net/http"
)

var (
	ErrUnauthorized = &HTTPError{
		Code:    http.StatusUnauthorized,
		Message: http.StatusText(http.StatusUnauthorized),
	}
	ErrNotFound = &HTTPError{
		Code:    http.StatusNotFound,
		Message: http.StatusText(http.StatusNotFound),
	}
	ErrInternalServer = &HTTPError{
		Code:    http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
	}
	ErrBadRequest = &HTTPError{
		Code:    http.StatusBadRequest,
		Message: http.StatusText(http.StatusBadRequest),
	}
	ErrInvalidJSON = &HTTPError{
		Code:    http.StatusBadRequest,
		Message: "invalid-json",
	}
	ErrRequestCanceled = &HTTPError{
		Code:    http.StatusBadRequest,
		Message: "request-canceled",
	}
)
