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
	noteRequestChan, noteServiceChan := make(chan *noteRequest), make(chan *domain.Note)
	user := r.Context().Value(adapterHttp.UserCtxKey).(*domain.User)

	doneReadChan, errReadChan := readNotes(r.Context(), r.Body, noteRequestChan)
	noteIDChan, doneCreateChan, errCreateChan := h.noteServicePort.Create(r.Context(), user, noteServiceChan)

	noteCounter := 0
	noteIDs := make([]string, 0)
	doneRequestChan := make(chan struct{})
	errRequestChan := make(chan error)

	go func() {
		for {
			select {
			case <-r.Context().Done():
				errRequestChan <- customerrors.ErrBadRequest
				return
			case err := <-errReadChan:
				log.Err(err).Msg("error while parsing note json")
				errRequestChan <- customerrors.Create(customerrors.ErrBadRequest.Code, "invalid-json")
				return
			case noteReq := <-noteRequestChan:
				noteServiceChan <- noteHTTPToDomain(noteReq)
			case <-doneReadChan:
				close(noteServiceChan)
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case noteID := <-noteIDChan:
				noteCounter++
				noteIDs = append(noteIDs, noteID)
			case err := <-errCreateChan:
				log.Err(err).Msg("repository error")
				errRequestChan <- customerrors.Create(customerrors.ErrInternalServer.Code, "repository-error")
				return
			case <-r.Context().Done():
				errRequestChan <- customerrors.ErrBadRequest
				return
			case <-doneCreateChan:
				doneRequestChan <- struct{}{}
			}
		}
	}()

	select {
	case <-r.Context().Done():
		log.Info().Msg("request canceled")
		return nil, customerrors.ErrBadRequest
	case err := <-errRequestChan:
		log.Err(err).Msg("eror while request")
		return nil, err
	case <-doneCreateChan:
		return &adapterHttp.Response{
			Data: &сreateNoteResponse{
				TotalNotes: noteCounter,
				NoteIDs:    noteIDs,
			},
			Status: http.StatusOK,
		}, nil
	}
}
