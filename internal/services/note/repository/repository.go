package repository

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
	"github.com/sony/sonyflake"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/stream"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/services/noteoutbox"
)

type Repository interface {
	CreateNote(
		ctx context.Context,
		st stream.Stream,
	)
	ReadNoteByID(ctx context.Context, user *domain.User, id uint64) (*domain.Note, error)
	UpdateNote(ctx context.Context, user *domain.User, note *domain.Note) error
	DeleteNote(ctx context.Context, user *domain.User, id uint64) (bool, error)
	SearchNote(
		ctx context.Context,
		user *domain.User,
		criteria *domain.SearchCriteria,
	) ([]*domain.Note, error)
}

type repositoryImpl struct {
	cfg         *Config
	db          *pgxpool.Pool
	outboxRepo  noteoutbox.Repository
	logger      *zerolog.Logger
	idGenerator *sonyflake.Sonyflake
}

func NewRepository(
	logger *zerolog.Logger,
	cfg *Config,
	db *pgxpool.Pool,
	outboxRepo noteoutbox.Repository,
) Repository {
	idGenerator := sonyflake.NewSonyflake(sonyflake.Settings{})
	return &repositoryImpl{logger: logger, cfg: cfg, db: db, outboxRepo: outboxRepo, idGenerator: idGenerator}
}
