package handlers

import (
	"net/http"

	"github.com/rs/zerolog/log"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/customerrors"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"

	adapterHttp "gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/http"
)

type сreateNoteResponse struct {
	NoteIDs    []string `json:"note_id"`
	TotalNotes int      `json:"total_notes"`
}

func (h *NoteHandlers) CreateNote(r *http.Request) (*adapterHttp.Response, error) {
	user := r.Context().Value(adapterHttp.UserCtxKey).(*domain.User)
	noteChan, errReadBodyChan := readNotes(r.Context(), r.Body)
	noteIDsCreatedChan, errCreateChan := h.noteServicePort.Create(r.Context(), user, noteChan)
	noteCounter, noteIDs := 0, make([]string, 0)

	for {
		select {
		case <-r.Context().Done():
			log.Info().Msg("request canceled")
			return nil, customerrors.ErrBadRequest
		case err := <-errReadBodyChan:
			log.Err(err).Msg("error while parsing note json")
			return nil, customerrors.Create(customerrors.ErrBadRequest.Code, "invalid-json")
		case err := <-errCreateChan:
			log.Err(err).Msg("repository error")
			return nil, customerrors.Create(customerrors.ErrInternalServer.Code, "repository-error")
		case noteID, ok := <-noteIDsCreatedChan:
			if ok {
				noteCounter++
				noteIDs = append(noteIDs, noteID)
			} else {
				return &adapterHttp.Response{
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
