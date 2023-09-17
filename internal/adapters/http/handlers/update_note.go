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

type updateNoteResponse struct {
	TotalNotes int `json:"total_notes"`
}

func (h *NoteHandlers) UpdateNote(r *http.Request) (*adapterHTTP.Response, error) {
	user := r.Context().Value(adapterHTTP.UserCtxKey).(*domain.User)
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
		case <-r.Context().Done():
			h.logger.Info().Msg("request canceled")
			return nil, customerrors.ErrBadRequest
		case err := <-st.ErrChan():
			h.logger.Err(err).Msgf("error while updating note")
			return nil, err
		case <-st.Done():
			err := st.Err()
			if errors.Is(err, context.Canceled) {
				return nil, customerrors.ErrRequestCanceled
			}
			return nil, err
		case note, ok := <-st.InRead():
			if !ok {
				return &adapterHTTP.Response{Data: &updateNoteResponse{TotalNotes: updatesCounter}, Status: http.StatusOK}, nil
			}
			updatesCounter++
			err := h.noteServicePort.Update(ctx, user, note)
			if err != nil {
				h.logger.Err(err).Msg("error while updating note")
				return nil, err
			}
		}
	}
}
