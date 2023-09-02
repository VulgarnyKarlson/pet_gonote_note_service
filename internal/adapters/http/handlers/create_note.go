package handlers

import (
	"net/http"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/customerrors"

	"github.com/rs/zerolog/log"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
)

type сreateNoteResponse struct {
	NoteIDs    []string `json:"note_id"`
	TotalNotes int      `json:"total_notes"`
}

func (h *NoteHandlers) CreateNote(r *http.Request) (*domain.HTTPResponse, error) {
	user := r.Context().Value(domain.UserCtxKey).(*domain.User)
	st, ctx := domain.NewStream(r.Context())
	defer st.Destroy()
	h.noteServicePort.Create(ctx, user, st)
	go func() {
		err := readNotes(r.Body, st)
		if err != nil {
			st.Fail(err)
		}
	}()

	noteCounter, noteIDs := 0, make([]string, 0)

	for {
		select {
		case <-st.Done():
			if err := st.Err(); err != nil {
				if err.Error() == "context canceled" {
					return nil, customerrors.ErrRequestCanceled
				}
				log.Err(err).Msg("error while creating note")
				return nil, err
			}
		case err := <-st.ErrChan():
			log.Err(err).Msg("error while creating note")
			return nil, err
		case noteID, ok := <-st.OutRead():
			if ok {
				noteCounter++
				noteIDs = append(noteIDs, noteID)
			} else {
				return &domain.HTTPResponse{
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
