package handlers

import "gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"

type NoteHandlers struct {
	noteServicePort domain.NoteServicePort
}

func New(n domain.NoteServicePort) *NoteHandlers {
	return &NoteHandlers{
		noteServicePort: n,
	}
}
