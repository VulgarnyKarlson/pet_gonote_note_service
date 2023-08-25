package handlers

import (
	"net/http"

	"github.com/rs/zerolog/log"
	adapterHttp "gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/http"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/customerrors"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
)

type updateNoteResponse struct {
	TotalNotes int `json:"total_notes"`
}

func (h *NoteHandlers) UpdateNote(r *http.Request) (*adapterHttp.Response, error) {
	noteChan := make(chan *noteRequest)
	user := r.Context().Value(adapterHttp.UserCtxKey).(*domain.User)
	doneChan, errChan := readNotes(r.Context(), r.Body, noteChan)
	updatesCounter := 0
	for {
		select {
		case <-r.Context().Done():
			log.Info().Msg("request canceled")
			return nil, customerrors.ErrBadRequest
		case err := <-errChan:
			log.Printf("error while parsing note json: %+v", err)
			return nil, customerrors.Create(customerrors.ErrBadRequest.Code, "invalid-json")
		case note := <-noteChan:
			updatesCounter++
			log.Printf("[%d] received node: %+v", updatesCounter, note)
			if note.ID == "" {
				log.Error().Msg("note id is empty")
				return nil, customerrors.Create(customerrors.ErrBadRequest.Code, "invalid-json: note-id-is-empty")
			}
			n := noteHTTPToDomain(note)
			err := h.noteServicePort.Update(r.Context(), user, n)
			if err != nil {
				return nil, customerrors.Create(customerrors.ErrBadRequest.Code, err.Error())
			}
		case <-doneChan:
			log.Printf("finished reading notes")
			return &adapterHttp.Response{Data: &updateNoteResponse{TotalNotes: updatesCounter}, Status: http.StatusOK}, nil
		}
	}
}
