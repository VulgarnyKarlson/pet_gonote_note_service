package handlers

import (
	"net/http"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/customerrors"
)

type readNoteResponse struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (h *NoteHandlers) ReadNoteByID(r *http.Request) (*domain.HTTPResponse, error) {
	noteID := r.URL.Query().Get("note_id")
	if noteID == "" {
		return nil, customerrors.Create(http.StatusBadRequest, "Note ID is required")
	}

	user := r.Context().Value(domain.UserCtxKey).(*domain.User)
	note, err := h.noteServicePort.ReadByID(r.Context(), user, noteID)
	if err != nil {
		return nil, customerrors.Create(http.StatusInternalServerError, err.Error())
	}

	if note == nil {
		return nil, customerrors.ErrNotFound
	}

	resp := readNoteResponse{
		Title:   note.Title,
		Content: note.Content,
	}

	return &domain.HTTPResponse{Data: resp, Status: http.StatusOK}, nil
}
