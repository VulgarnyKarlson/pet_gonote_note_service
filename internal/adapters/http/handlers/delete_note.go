package handlers

import (
	"encoding/json"
	"net/http"

	adapterHttp "gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/http"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/customerrors"
)

type deleteNoteRequest struct {
	NoteID string `json:"id"`
}

type deleteNoteResponse struct {
	Status string `json:"status"`
}

func (h *NoteHandlers) DeleteNoteByID(r *http.Request) (*adapterHttp.Response, error) {
	var req deleteNoteRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, customerrors.ErrBadRequest
	}

	user := r.Context().Value(adapterHttp.UserCtxKey).(*domain.User)
	if isUpdated, err := h.noteServicePort.Delete(r.Context(), user, req.NoteID); err != nil {
		return nil, customerrors.ErrInternalServer
	} else if !isUpdated {
		return &adapterHttp.Response{Data: &deleteNoteResponse{Status: "NotFound"}, Status: http.StatusNotFound}, nil
	}

	return &adapterHttp.Response{Data: &deleteNoteResponse{Status: "Success"}, Status: http.StatusOK}, nil
}
