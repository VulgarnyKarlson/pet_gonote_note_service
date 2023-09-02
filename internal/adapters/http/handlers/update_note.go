package handlers

import (
	"net/http"

	"github.com/rs/zerolog/log"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/customerrors"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
)

type updateNoteResponse struct {
	TotalNotes int `json:"total_notes"`
}

func (h *NoteHandlers) UpdateNote(r *http.Request) (*domain.HTTPResponse, error) {
	user := r.Context().Value(domain.UserCtxKey).(*domain.User)
	st, ctx := domain.NewStream(r.Context())
	defer st.Destroy()
	go func() {
		err := readNotes(r.Body, st)
		if err != nil {
			st.Fail(err)
		}
	}()
	updatesCounter := 0
	for {
		select {
		case <-r.Context().Done():
			log.Info().Msg("request canceled")
			return nil, customerrors.ErrBadRequest
		case err := <-st.ErrChan():
			log.Err(err).Msgf("error while updating note: %+v", err)
			return nil, err
		case <-st.Done():
			if err := st.Err(); err != nil {
				if err.Error() == "context canceled" {
					return nil, customerrors.ErrRequestCanceled
				}
				log.Err(err).Msg("error while updating note")
				return nil, err
			}
		case note, ok := <-st.InRead():
			if !ok {
				return &domain.HTTPResponse{Data: &updateNoteResponse{TotalNotes: updatesCounter}, Status: http.StatusOK}, nil
			}
			updatesCounter++
			err := h.noteServicePort.Update(ctx, user, note)
			if err != nil {
				return nil, err
			}
		}
	}
}
