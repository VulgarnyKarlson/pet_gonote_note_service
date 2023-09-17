package handlers

import (
	"context"
	"errors"
	"net/http"

	adapterHTTP "gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/http"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/customerrors"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/stream"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
)

type сreateNoteResponse struct {
	NoteIDs    []uint64 `json:"note_id"`
	TotalNotes int      `json:"total_notes"`
}

func (h *NoteHandlers) CreateNote(r *http.Request) (*adapterHTTP.Response, error) {
	st, ctx := stream.NewStream(r.Context())
	defer st.Destroy()
	user := r.Context().Value(adapterHTTP.UserCtxKey).(*domain.User)
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
				return &adapterHTTP.Response{
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
