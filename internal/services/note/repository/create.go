package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/stream"

	"github.com/Masterminds/squirrel"
	"github.com/hashicorp/go-uuid"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/postgres"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/customerrors"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
)

func (r *repositoryImpl) CreateNote(
	ctx context.Context,
	user *domain.User,
	st stream.Stream,
) {
	tx, err := r.db.NewTransaction(ctx, pgx.TxOptions{})
	if err != nil {
		log.Err(err).Msg("error creating transaction")
		st.Fail(customerrors.ErrRepositoryError)
		return
	}
	var txErr error
	var isCommitted bool
	defer tx.ConnRelease()
	defer func() {
		if !isCommitted {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				txErr = errors.Join(txErr, rollbackErr)
			}
			st.Fail(txErr)
		}
	}()

	batch := make([]*domain.Note, 0, r.cfg.CreateNotesBatchSize)
loop:
	for {
		select {
		case <-ctx.Done():
			return
		case <-st.Done():
			return
		case note, ok := <-st.InProxyRead():
			if !ok {
				break loop
			}
			batch = append(batch, note)
			if len(batch) == r.cfg.CreateNotesBatchSize {
				if err = r.insertBatch(tx, batch, user, st); err != nil {
					txErr = err
					return
				}

				batch = batch[:0]
			}
		}
	}

	if len(batch) > 0 {
		if err = r.insertBatch(tx, batch, user, st); err != nil {
			txErr = err
			return
		}
	}

	if err = st.Err(); err != nil {
		txErr = err
		return
	}

	if err = tx.Commit(); err != nil {
		txErr = err
		return
	}
	isCommitted = true

	st.OutClose()
	st.Close()
}

func (r *repositoryImpl) insertBatch(tx *postgres.Transaction, notes []*domain.Note, user *domain.User, st stream.Stream) error {
	noteBatch := &pgx.Batch{}
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	mBatch := make([]*DBModel, 0, len(notes))
	for _, note := range notes {
		noteID, err := uuid.GenerateUUID()
		if err != nil {
			return fmt.Errorf("error generating uuid: %w", err)
		}
		note.SetID(noteID)
		domainNote := noteDomainToDBModel(note)
		mBatch = append(mBatch, domainNote)
	}

	for _, note := range mBatch {
		query, args, _ := psql.Insert("notes").
			Columns("id", "user_id", "title", "content", "created_at", "updated_at").
			Values(note.ID, user.ID(), note.Title, note.Content, note.CreatedAt.Format(time.RFC3339), note.UpdatedAt.Format(time.RFC3339)).
			ToSql()

		noteBatch.Queue(query, args...)
	}

	br := tx.SendBatch(noteBatch)

	_, err := br.Exec()
	if err != nil {
		return fmt.Errorf("trx err: %w", err)
	}
	err = br.Close()
	if err != nil {
		return fmt.Errorf("error closing batch: %w", err)
	}

	for _, note := range notes {
		err = r.outboxRepo.Create(tx, note)
		if err != nil {
			return fmt.Errorf("error creating note outbox: %w", err)
		}
		st.OutWrite(note.ID())
	}

	return nil
}
