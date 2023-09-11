package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/customerrors"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/stream"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
)

func noteHTTPToDomain(n *noteRequest, user *domain.User) (*domain.Note, error) {
	return domain.NewNote(n.ID, user.ID(), n.Title, n.Content)
}

func searchCriteriaHTTPToDomain(s *searchNoteRequest) (*domain.SearchCriteria, error) {
	fromDate := time.Time{}
	if s.FromDate != "" {
		var err error
		fromDate, err = time.Parse(time.RFC3339, s.FromDate)
		if err != nil {
			return nil, fmt.Errorf("failed to parse from date: %w", err)
		}
	}

	toDate := time.Time{}
	if s.ToDate != "" {
		var err error
		toDate, err = time.Parse(time.RFC3339, s.ToDate)
		if err != nil {
			return nil, fmt.Errorf("failed to parse to date: %w", err)
		}
	}

	return &domain.SearchCriteria{
		Title:    s.Title,
		Content:  s.Content,
		FromDate: fromDate,
		ToDate:   toDate,
	}, nil
}

func readNotes(r io.Reader, st stream.Stream, user *domain.User) error {
	decoder := json.NewDecoder(r)
	if delim, err := decoder.Token(); delim != json.Delim('[') || err != nil {
		return customerrors.ErrInvalidJSONOpenDelimiter
	}

	for decoder.More() {
		select {
		case <-st.Done():
			return nil
		default:
			var note noteRequest
			if err := decoder.Decode(&note); err != nil {
				return customerrors.ErrInvalidJSON
			}
			domainNote, err := noteHTTPToDomain(&note, user)
			if err != nil {
				return customerrors.ErrInvalidNote
			}
			st.InWrite(domainNote)
		}
	}

	if delim, err := decoder.Token(); delim != json.Delim(']') || err != nil {
		return customerrors.ErrInvalidJSONCloseDelimiter
	}

	st.InClose()

	return nil
}
