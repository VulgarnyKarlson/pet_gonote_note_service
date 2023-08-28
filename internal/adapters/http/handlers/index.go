package handlers

import "gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"

type NoteHandlers struct {
	noteServicePort domain.NoteService
}

func New(n domain.NoteService) *NoteHandlers {
	return &NoteHandlers{
		noteServicePort: n,
	}
}
