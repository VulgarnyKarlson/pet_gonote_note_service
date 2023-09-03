package repository

import (
	"time"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
)

type DBModel struct {
	ID        string    `json:"id,omitempty"`
	UserID    string    `json:"user_id,omitempty"`
	Title     string    `json:"title,omitempty"`
	Content   string    `json:"content,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func noteDomainToDBModel(note *domain.Note) *DBModel {
	return &DBModel{
		ID:        note.ID(),
		UserID:    note.UserID(),
		Title:     note.Title(),
		Content:   note.Content(),
		CreatedAt: note.CreatedAt(),
		UpdatedAt: note.UpdatedAt(),
	}
}

func noteDBModelToDomain(dbNote *DBModel) (*domain.Note, error) {
	var err error
	var note *domain.Note
	note, err = domain.NewNote(dbNote.ID, dbNote.UserID, dbNote.Title, dbNote.Content)
	if err != nil {
		return nil, err
	}
	note.SetCreatedAt(dbNote.CreatedAt)
	note.SetUpdatedAt(dbNote.UpdatedAt)
	return note, nil
}
