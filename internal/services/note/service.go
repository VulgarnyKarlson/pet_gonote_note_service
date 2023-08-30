package note

import (
	"context"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/customerrors"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/services/note/repository"
)

type serviceImpl struct {
	repo repository.Repository
}

func NewService(r repository.Repository) domain.NoteService {
	return &serviceImpl{repo: r}
}

func (s *serviceImpl) Create(
	ctx context.Context,
	user *domain.User,
) (inputNoteChan chan *domain.Note, outputNoteIDsChan chan string, errChan chan error) {
	inputNoteChan = make(chan *domain.Note)
	outputNoteChan, errChan := s.repo.CreateNote(ctx, user, inputNoteChan)
	return inputNoteChan, outputNoteChan, errChan
}

func (s *serviceImpl) ReadByID(ctx context.Context, user *domain.User, id string) (*domain.Note, error) {
	if id == "" {
		return nil, customerrors.ErrInvalidNoteID
	}
	return s.repo.ReadNoteByID(ctx, user, id)
}

func (s *serviceImpl) Update(ctx context.Context, user *domain.User, note *domain.Note) error {
	if note.ID == "" {
		return customerrors.ErrInvalidNoteID
	}
	return s.repo.UpdateNote(ctx, user, note)
}

func (s *serviceImpl) Delete(ctx context.Context, user *domain.User, id string) (bool, error) {
	if id == "" {
		return false, customerrors.ErrInvalidNoteID
	}
	return s.repo.DeleteNote(ctx, user, id)
}

func (s *serviceImpl) Search(
	ctx context.Context,
	user *domain.User,
	criteria *domain.SearchCriteria,
) ([]*domain.Note, error) {
	return s.repo.SearchNote(ctx, user, criteria)
}
