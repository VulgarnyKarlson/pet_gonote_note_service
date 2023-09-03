package handlers

import (
	"net/http"

	adapterHTTP "gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/http"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/customerrors"
)

type readNoteResponse struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (h *NoteHandlers) ReadNoteByID(r *http.Request) (*adapterHTTP.Response, error) {
	noteID := r.URL.Query().Get("note_id")
	user := r.Context().Value(adapterHTTP.UserCtxKey).(*domain.User)
	note, err := h.noteServicePort.ReadByID(r.Context(), user, noteID)
	if err != nil {
		return nil, err
	}

	if note == nil {
		return nil, customerrors.ErrNotFound
	}

	resp := readNoteResponse{
		Title:   note.Title(),
		Content: note.Content(),
	}

	return &adapterHTTP.Response{Data: resp, Status: http.StatusOK}, nil
}
