package note

import (
	"context"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/services/note/repository"
)

type Service interface {
	Create(
		ctx context.Context,
		user *domain.User,
		note chan *domain.Note,
	) (noteIDsChan chan string, doneChan chan struct{}, errChan chan error)
	ReadByID(ctx context.Context, user *domain.User, id string) (*domain.Note, error)
	Update(ctx context.Context, user *domain.User, note *domain.Note) error
	Delete(ctx context.Context, user *domain.User, id string) (bool, error)
	Search(
		ctx context.Context,
		user *domain.User,
		criteria *domain.SearchCriteria,
	) (noteIDS []*domain.Note, err error)
}

type serviceImpl struct {
	repo repository.Repository
}

func NewService(r repository.Repository) Service {
	return &serviceImpl{repo: r}
}

func (s *serviceImpl) Create(
	ctx context.Context,
	user *domain.User,
	note chan *domain.Note,
) (notes chan string, doneChan chan struct{}, err chan error) {
	return s.repo.Create(ctx, user, note)
}

func (s *serviceImpl) ReadByID(ctx context.Context, user *domain.User, id string) (*domain.Note, error) {
	return s.repo.ReadByID(ctx, user, id)
}

func (s *serviceImpl) Update(ctx context.Context, user *domain.User, note *domain.Note) error {
	return s.repo.Update(ctx, user, note)
}

func (s *serviceImpl) Delete(ctx context.Context, user *domain.User, id string) (bool, error) {
	return s.repo.Delete(ctx, user, id)
}

func (s *serviceImpl) Search(
	ctx context.Context,
	user *domain.User,
	criteria *domain.SearchCriteria,
) ([]*domain.Note, error) {
	return s.repo.Search(ctx, user, criteria)
}
