package handlers

import (
	"github.com/rs/zerolog"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/services/note"
)

type NoteHandlers struct {
	noteServicePort note.Service
	logger          *zerolog.Logger
}

func New(n note.Service, logger *zerolog.Logger) *NoteHandlers {
	return &NoteHandlers{
		noteServicePort: n,
		logger:          logger,
	}
}
