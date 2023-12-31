package handlers

import (
	"time"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/server/middlewares"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
)

var activeMiddlewares = []string{middlewares.LoggerID(), middlewares.AuthID()}

type noteRequest struct {
	ID      uint64 `json:"id,omitempty"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
}

type noteResponse struct {
	ID        uint64    `json:"id,omitempty"`
	Title     string    `json:"title,omitempty"`
	Content   string    `json:"content,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func noteDomainToHTTP(note *domain.Note) *noteResponse {
	return &noteResponse{
		ID:        note.ID(),
		Title:     note.Title(),
		Content:   note.Content(),
		CreatedAt: note.CreatedAt(),
		UpdatedAt: note.UpdatedAt(),
	}
}
