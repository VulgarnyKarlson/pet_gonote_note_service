package note

import (
	"context"
	"strconv"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/customerrors"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/stream"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/services/note/repository"
)

type Service interface {
	Create(ctx context.Context, st stream.Stream)
	ReadByID(ctx context.Context, user *domain.User, id string) (*domain.Note, error)
	Update(ctx context.Context, user *domain.User, note *domain.Note) error
	Delete(ctx context.Context, user *domain.User, id string) (bool, error)
	Search(ctx context.Context, user *domain.User, criteria *domain.SearchCriteria) ([]*domain.Note, error)
}

type serviceImpl struct {
	cfg  *Config
	repo repository.Repository
}

func NewService(cfg *Config, r repository.Repository) Service {
	return &serviceImpl{cfg: cfg, repo: r}
}

func (s *serviceImpl) Create(
	ctx context.Context,
	st stream.Stream,
) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-st.Done():
				return
			case note, ok := <-st.InRead():
				if !ok {
					st.InProxyClose()
					return
				}
				if len(note.Title()) > s.cfg.MaxTitleLength {
					st.Fail(customerrors.ErrTitleTooLong)
					return
				}

				if len(note.Content()) > s.cfg.MaxContentLength {
					st.Fail(customerrors.ErrContentTooLong)
					return
				}

				st.InProxyWrite(note)
			}
		}
	}()
	go s.repo.CreateNote(ctx, st)
}

func (s *serviceImpl) ReadByID(ctx context.Context, user *domain.User, id string) (*domain.Note, error) {
	if id == "" {
		return nil, customerrors.ErrInvalidNoteID
	}
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil || idUint == 0 {
		return nil, customerrors.ErrInvalidNoteID
	}
	return s.repo.ReadNoteByID(ctx, user, idUint)
}

func (s *serviceImpl) Update(ctx context.Context, user *domain.User, note *domain.Note) error {
	if note.ID() == 0 {
		return customerrors.ErrInvalidNoteID
	}
	return s.repo.UpdateNote(ctx, user, note)
}

func (s *serviceImpl) Delete(ctx context.Context, user *domain.User, id string) (bool, error) {
	if id == "" {
		return false, customerrors.ErrInvalidNoteID
	}
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil || idUint == 0 {
		return false, customerrors.ErrInvalidNoteID
	}
	return s.repo.DeleteNote(ctx, user, idUint)
}

func (s *serviceImpl) Search(
	ctx context.Context,
	user *domain.User,
	criteria *domain.SearchCriteria,
) ([]*domain.Note, error) {
	return s.repo.SearchNote(ctx, user, criteria)
}
