package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

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

func readNotes(ctx context.Context, r io.Reader, noteChan chan *noteRequest) (doneChan chan struct{}, errChan chan error) {
	errChan = make(chan error)
	doneChan = make(chan struct{})
	go func() {
		decoder := json.NewDecoder(r)

		if _, err := decoder.Token(); err != nil {
			errChan <- fmt.Errorf("failed to read opening delimiter: %w", err)
			return
		}

		for decoder.More() {
			if ctx.Err() != nil {
				errChan <- ctx.Err()
				return
			}

			var note noteRequest
			if err := decoder.Decode(&note); err != nil {
				errChan <- fmt.Errorf("failed to decode opening token: %w", err)
				return
			}

			noteChan <- &note
		}

		if _, err := decoder.Token(); err != nil {
			errChan <- fmt.Errorf("failed to read closing delimiter: %w", err)
			return
		}

		close(doneChan)
	}()

	return doneChan, errChan
}
