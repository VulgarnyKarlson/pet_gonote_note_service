package handlers

import (
	"net/http"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/server"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/server/middlewares"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
	noteService "gitlab.karlson.dev/individual/pet_gonote/note_service/internal/services/note"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/customerrors"
)

type readNoteHandler struct {
	noteServicePort noteService.Service
}

func RegisterReadNote(s server.Server, n noteService.Service) {
	h := readNoteHandler{noteServicePort: n}
	s.AddEndpoint(server.Endpoint{Method: "GET", Path: "/read", Middlewares: activeMiddlewares, Handler: h.handle})
}

type readNoteResponse struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (h *readNoteHandler) handle(r *http.Request) (*server.Response, error) {
	noteID := r.URL.Query().Get("note_id")
	user := r.Context().Value(middlewares.UserCtxKey).(*domain.User)
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

	return &server.Response{Data: resp, Status: http.StatusOK}, nil
}
