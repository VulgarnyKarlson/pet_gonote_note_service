package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/customerrors"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
)

func noteHTTPToDomain(n *noteRequest) *domain.Note {
	return &domain.Note{
		ID:      n.ID,
		Title:   n.Title,
		Content: n.Content,
	}
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

func readNotes(r io.Reader, st domain.Stream) error {
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

			st.InWrite(noteHTTPToDomain(&note))
		}
	}

	if delim, err := decoder.Token(); delim != json.Delim(']') || err != nil {
		return customerrors.ErrInvalidJSONCloseDelimiter
	}

	st.InClose()

	return nil
}
