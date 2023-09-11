package handlers

import (
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
			if err := st.Err(); err != nil {
				if err.Error() == "context canceled" {
					return nil, customerrors.ErrRequestCanceled
				}
				h.logger.Err(err).Msg("error while creating note")
				return nil, err
			}
		case err := <-st.ErrChan():
			h.logger.Err(err).Msg("error while creating note")
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
