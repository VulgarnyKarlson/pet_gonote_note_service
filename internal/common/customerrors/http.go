package customerrors

import (
	"fmt"
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
)

func Create(code int, message string, args ...any) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: fmt.Sprintf(message, args...),
	}
}
