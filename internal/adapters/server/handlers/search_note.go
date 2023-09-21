package handlers

import (
	"net/http"

	"github.com/pkg/errors"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/server"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/server/middlewares"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/customerrors"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
	noteService "gitlab.karlson.dev/individual/pet_gonote/note_service/internal/services/note"
)

type searchNoteHandler struct {
	noteServicePort noteService.Service
}

func RegisterSearchNote(s server.Server, n noteService.Service) {
	h := searchNoteHandler{noteServicePort: n}
	s.AddEndpoint(server.Endpoint{Method: "GET", Path: "/search", Middlewares: activeMiddlewares, Handler: h.handle})
}

type searchNoteRequest struct {
	Title    string
	Content  string
	FromDate string
	ToDate   string
}

type searcNoteResponse struct {
	Notes []*noteResponse `json:"notes"`
	Total int             `json:"total"`
}

func (h *searchNoteHandler) handle(r *http.Request) (*server.Response, error) {
	var req searchNoteRequest
	req.Title = r.URL.Query().Get("title")
	req.Content = r.URL.Query().Get("content")
	req.FromDate = r.URL.Query().Get("from_date")
	req.ToDate = r.URL.Query().Get("to_date")
	searchNoteDomain, err := searchCriteriaHTTPToDomain(&req)
	if err != nil {
		return nil, errors.Wrap(err, "error while converting search criteria")
	}

	user := r.Context().Value(middlewares.UserCtxKey).(*domain.User)
	notes, err := h.noteServicePort.Search(r.Context(), user, searchNoteDomain)
	if err != nil {
		return nil, customerrors.ErrInternalServer
	}

	if len(notes) == 0 {
		return nil, customerrors.ErrNotFound
	}

	var noteHTTPresp []*noteResponse
	for _, note := range notes {
		noteHTTPresp = append(noteHTTPresp, noteDomainToHTTP(note))
	}

	resp := &searcNoteResponse{
		Notes: noteHTTPresp,
		Total: len(noteHTTPresp),
	}

	return &server.Response{Data: resp, Status: http.StatusOK}, nil
}
