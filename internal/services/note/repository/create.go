package repository

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/customerrors"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/stream"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
)

func (r *repositoryImpl) CreateNote(
	ctx context.Context,
	st stream.Stream,
) {
	err := r.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		batch := make([]*domain.Note, 0, r.cfg.CreateNotesBatchSize)
	loop:
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-st.Done():
				return nil
			case note, ok := <-st.InProxyRead():
				if !ok {
					break loop
				}
				batch = append(batch, note)
				if len(batch) == r.cfg.CreateNotesBatchSize {
					if err := r.insertBatch(ctx, tx, batch, st); err != nil {
						return err
					}

					batch = batch[:0]
				}
			}
		}

		if len(batch) > 0 {
			if err := r.insertBatch(ctx, tx, batch, st); err != nil {
				return nil
			}
		}

		if err := st.Err(); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		r.logger.Err(err).Msg("error creating transaction")
		st.Fail(customerrors.ErrRepositoryError)
	}
	st.OutClose()
	st.Close()
}

func (r *repositoryImpl) insertBatch(
	ctx context.Context,
	tx pgx.Tx,
	notes []*domain.Note,
	st stream.Stream,
) error {
	noteBatch := &pgx.Batch{}
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	mBatch := make([]*DBModel, 0, len(notes))
	for _, note := range notes {
		if note.ID() == 0 {
			noteID, err := r.idGenerator.NextID()
			if err != nil {
				return fmt.Errorf("error generating uuid: %w", err)
			}
			note.SetID(noteID)
		}
		domainNote := noteDomainToDBModel(note)
		mBatch = append(mBatch, domainNote)
	}

	for _, note := range mBatch {
		query, args, _ := psql.Insert("notes").
			Columns("note_id", "user_id", "title", "content").
			Values(note.ID, note.UserID, note.Title, note.Content).
			Suffix("ON CONFLICT DO NOTHING").
			ToSql()

		noteBatch.Queue(query, args...)
	}

	br := tx.SendBatch(ctx, noteBatch)

	_, err := br.Exec()
	if err != nil {
		return fmt.Errorf("trx err: %w", err)
	}
	err = br.Close()
	if err != nil {
		return fmt.Errorf("error closing batch: %w", err)
	}

	for _, note := range notes {
		err = r.outboxRepo.Create(ctx, tx, note)
		if err != nil {
			return fmt.Errorf("error creating note outbox: %w", err)
		}
		st.OutWrite(note.ID())
	}

	return nil
}
