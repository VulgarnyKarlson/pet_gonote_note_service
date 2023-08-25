package repository

import (
	"context"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/services/noteoutbox"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/postgres"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
)

type Repository interface {
	Create(
		ctx context.Context,
		user *domain.User,
		noteChan chan *domain.Note,
	) (noteIDs chan string, doneChan chan struct{}, errChan chan error)
	ReadByID(ctx context.Context, user *domain.User, id string) (*domain.Note, error)
	Update(ctx context.Context, user *domain.User, note *domain.Note) error
	Delete(ctx context.Context, user *domain.User, id string) (bool, error)
	Search(
		ctx context.Context,
		user *domain.User,
		criteria *domain.SearchCriteria,
	) ([]*domain.Note, error)
}

type repositoryImpl struct {
	cfg        *Config
	db         *postgres.Pool
	outboxRepo noteoutbox.Repository
}

func NewRepository(cfg *Config, db *postgres.Pool, outboxRepo noteoutbox.Repository) Repository {
	return &repositoryImpl{cfg: cfg, db: db, outboxRepo: outboxRepo}
}
