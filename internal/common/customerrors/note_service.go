package customerrors

import "net/http"

var (
	ErrInvalidNoteID = &HTTPError{
		Code:    http.StatusBadRequest,
		Message: "invalid-note-id",
	}
)
