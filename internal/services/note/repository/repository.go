package repository

import (
	"context"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/stream"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/services/noteoutbox"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/postgres"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
)

type Repository interface {
	CreateNote(
		ctx context.Context,
		user *domain.User,
		st stream.Stream,
	)
	ReadNoteByID(ctx context.Context, user *domain.User, id string) (*domain.Note, error)
	UpdateNote(ctx context.Context, user *domain.User, note *domain.Note) error
	DeleteNote(ctx context.Context, user *domain.User, id string) (bool, error)
	SearchNote(
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
