package handlers

import (
	"context"
	"errors"
	"net/http"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/server"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/server/middlewares"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/customerrors"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/stream"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
	noteService "gitlab.karlson.dev/individual/pet_gonote/note_service/internal/services/note"
)

type updateNoteHandler struct {
	noteServicePort noteService.Service
}

func RegisterUpdateNote(s server.Server, n noteService.Service) {
	h := updateNoteHandler{noteServicePort: n}
	s.AddEndpoint(server.Endpoint{Method: "POST", Path: "/update", Middlewares: activeMiddlewares, Handler: h.handle})
}

type updateNoteResponse struct {
	TotalNotes int `json:"total_notes"`
}

func (h *updateNoteHandler) handle(r *http.Request) (*server.Response, error) {
	user := r.Context().Value(middlewares.UserCtxKey).(*domain.User)
	st, ctx := stream.NewStream(r.Context())
	defer st.Destroy()
	go func() {
		err := readNotes(r.Body, st, user)
		if err != nil {
			st.Fail(err)
		}
	}()
	updatesCounter := 0
	for {
		select {
		case <-st.Done():
			err := st.Err()
			if errors.Is(err, context.Canceled) {
				return nil, customerrors.ErrRequestCanceled
			}
			return nil, err
		case err := <-st.ErrChan():
			return nil, err
		case note, ok := <-st.InRead():
			if !ok {
				return &server.Response{Data: &updateNoteResponse{TotalNotes: updatesCounter}, Status: http.StatusOK}, nil
			}
			updatesCounter++
			err := h.noteServicePort.Update(ctx, user, note)
			if err != nil {
				return nil, err
			}
		}
	}
}
