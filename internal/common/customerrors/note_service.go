package customerrors

import "net/http"

var (
	ErrInvalidNoteID = &HTTPError{
		Code:    http.StatusBadRequest,
		Message: "invalid-note-id",
	}
	ErrNotFoundNoteID = &HTTPError{
		Code:    http.StatusNotFound,
		Message: "not-found-note-id",
	}
	ErrTitleTooLong = &HTTPError{
		Code:    http.StatusBadRequest,
		Message: "title-too-long",
	}
	ErrContentTooLong = &HTTPError{
		Code:    http.StatusBadRequest,
		Message: "content-too-long",
	}
	ErrRepositoryError = &HTTPError{
		Code:    http.StatusInternalServerError,
		Message: "repository-error",
	}
)
