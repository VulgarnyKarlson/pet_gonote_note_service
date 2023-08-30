package customerrors

import "net/http"

var (
	ErrInvalidNoteID = &HTTPError{
		Code:    http.StatusBadRequest,
		Message: "invalid-note-id",
	}
	ErrTitleTooLong = &HTTPError{
		Code:    http.StatusBadRequest,
		Message: "title-too-long",
	}
	ErrContentTooLong = &HTTPError{
		Code:    http.StatusBadRequest,
		Message: "content-too-long",
	}
)
