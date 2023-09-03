package handlers

import (
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/services/note"
)

type NoteHandlers struct {
	noteServicePort note.Service
}

func New(n note.Service) *NoteHandlers {
	return &NoteHandlers{
		noteServicePort: n,
	}
}
