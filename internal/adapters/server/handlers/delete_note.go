package handlers

import (
	"encoding/json"
	"net/http"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/server"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/server/middlewares"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/customerrors"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/services/note"
)

type deleteNoteHandler struct {
	noteServicePort note.Service
}

func RegisterDeleteNote(s server.Server, n note.Service) {
	h := deleteNoteHandler{noteServicePort: n}
	s.AddEndpoint(server.Endpoint{Method: "POST", Path: "/delete", Middlewares: activeMiddlewares, Handler: h.handle})
}

type deleteNoteRequest struct {
	NoteID string `json:"id"`
}

type deleteNoteResponse struct {
	Status string `json:"status"`
}

func (h *deleteNoteHandler) handle(r *http.Request) (*server.Response, error) {
	var req deleteNoteRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, customerrors.ErrBadRequest
	}

	user := r.Context().Value(middlewares.UserCtxKey).(*domain.User)
	if isUpdated, err := h.noteServicePort.Delete(r.Context(), user, req.NoteID); err != nil {
		return nil, customerrors.ErrInternalServer
	} else if !isUpdated {
		return &server.Response{Data: &deleteNoteResponse{Status: "NotFound"}, Status: http.StatusNotFound}, nil
	}

	return &server.Response{Data: &deleteNoteResponse{Status: "Success"}, Status: http.StatusOK}, nil
}
