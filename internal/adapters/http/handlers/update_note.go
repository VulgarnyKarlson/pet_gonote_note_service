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
	inputNoteChan := make(chan *domain.Note)
	errChan := readNotes(r.Context(), r.Body, inputNoteChan)
	updatesCounter := 0
	for {
		select {
		case <-r.Context().Done():
			log.Info().Msg("request canceled")
			return nil, customerrors.ErrBadRequest
		case err := <-errChan:
			log.Printf("error while parsing note json: %+v", err)
			return nil, customerrors.ErrInvalidJSON
		case note, ok := <-inputNoteChan:
			if !ok {
				log.Printf("finished reading notes")
				return &domain.HTTPResponse{Data: &updateNoteResponse{TotalNotes: updatesCounter}, Status: http.StatusOK}, nil
			}
			updatesCounter++
			log.Printf("[%d] received node: %+v", updatesCounter, note)
			if note.ID == "" {
				log.Error().Msg("note id is empty")
				return nil, customerrors.ErrInvalidNoteID
			}
			err := h.noteServicePort.Update(r.Context(), user, note)
			if err != nil {
				return nil, err
			}
		}
	}
}
