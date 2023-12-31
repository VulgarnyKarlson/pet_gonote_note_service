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
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/services/note"
)

type createNoteHandler struct {
	noteServicePort note.Service
}

func RegisterCreateNote(s server.Server, n note.Service) {
	h := createNoteHandler{noteServicePort: n}
	s.AddEndpoint(server.Endpoint{Method: "POST", Path: "/create", Middlewares: activeMiddlewares, Handler: h.handle})
}

type сreateNoteResponse struct {
	NoteIDs    []uint64 `json:"note_id"`
	TotalNotes int      `json:"total_notes"`
}

func (h *createNoteHandler) handle(r *http.Request) (*server.Response, error) {
	st, ctx := stream.NewStream(r.Context())
	defer st.Destroy()
	user := r.Context().Value(middlewares.UserCtxKey).(*domain.User)
	h.noteServicePort.Create(ctx, st)
	go func() {
		err := readNotes(r.Body, st, user)
		if err != nil {
			st.Fail(err)
		}
	}()

	noteCounter, noteIDs := 0, make([]uint64, 0)

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
		case noteID, ok := <-st.OutRead():
			if ok {
				noteCounter++
				noteIDs = append(noteIDs, noteID)
			} else {
				return &server.Response{
					Data: &сreateNoteResponse{
						TotalNotes: noteCounter,
						NoteIDs:    noteIDs,
					},
					Status: http.StatusOK,
				}, nil
			}
		}
	}
}
