package domain

import (
	"context"
	"time"
)

type Note struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SearchCriteria struct {
	Title    string
	Content  string
	FromDate time.Time
	ToDate   time.Time
}

type NoteService interface {
	Create(ctx context.Context, user *User) (inputNoteChan chan *Note, outputNoteIDsChan chan string, errChan chan error)
	ReadByID(ctx context.Context, user *User, id string) (*Note, error)
	Update(ctx context.Context, user *User, note *Note) error
	Delete(ctx context.Context, user *User, id string) (bool, error)
	Search(ctx context.Context, user *User, criteria *SearchCriteria) ([]*Note, error)
}
