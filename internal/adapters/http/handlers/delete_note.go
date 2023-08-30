package handlers

import (
	"encoding/json"
	"net/http"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/customerrors"
)

type deleteNoteRequest struct {
	NoteID string `json:"id"`
}

type deleteNoteResponse struct {
	Status string `json:"status"`
}

func (h *NoteHandlers) DeleteNoteByID(r *http.Request) (*domain.HTTPResponse, error) {
	var req deleteNoteRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, customerrors.ErrBadRequest
	}

	user := r.Context().Value(domain.UserCtxKey).(*domain.User)
	if isUpdated, err := h.noteServicePort.Delete(r.Context(), user, req.NoteID); err != nil {
		return nil, customerrors.ErrInternalServer
	} else if !isUpdated {
		return &domain.HTTPResponse{Data: &deleteNoteResponse{Status: "NotFound"}, Status: http.StatusNotFound}, nil
	}

	return &domain.HTTPResponse{Data: &deleteNoteResponse{Status: "Success"}, Status: http.StatusOK}, nil
}
