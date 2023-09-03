package handlers

import (
	"net/http"

	adapterHTTP "gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/http"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/customerrors"
)

type searchNoteRequest struct {
	Title    string
	Content  string
	FromDate string
	ToDate   string
}

type searcNoteResponse struct {
	Notes []*domain.Note `json:"notes"`
	Total int            `json:"total"`
}

func (h *NoteHandlers) SearchNote(r *http.Request) (*adapterHTTP.Response, error) {
	var req searchNoteRequest
	req.Title = r.URL.Query().Get("title")
	req.Content = r.URL.Query().Get("content")
	req.FromDate = r.URL.Query().Get("from_date")
	req.ToDate = r.URL.Query().Get("to_date")
	searchNoteDomain, err := searchCriteriaHTTPToDomain(&req)
	if err != nil {
		return nil, err
	}

	user := r.Context().Value(adapterHTTP.UserCtxKey).(*domain.User)
	notes, err := h.noteServicePort.Search(r.Context(), user, searchNoteDomain)
	if err != nil {
		return nil, customerrors.ErrInternalServer
	}

	if len(notes) == 0 {
		return nil, customerrors.ErrNotFound
	}

	resp := searcNoteResponse{
		Notes: notes,
		Total: len(notes),
	}

	return &adapterHTTP.Response{Data: resp, Status: http.StatusOK}, nil
}
