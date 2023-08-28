package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/hashicorp/go-uuid"
	"github.com/jackc/pgx/v4"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/postgres"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
)

func (r *repositoryImpl) CreateNote(
	ctx context.Context,
	user *domain.User,
	noteChan chan *domain.Note,
) (noteIDs chan string, errChan chan error) {
	noteIDs = make(chan string)
	errChan = make(chan error)
	tx, err := r.db.NewTransaction(ctx, pgx.TxOptions{})
	if err != nil {
		errChan <- fmt.Errorf("error creating transaction: %w", err)
		return noteIDs, errChan
	}

	go func() {
		defer tx.ConnRelease()
		defer func() {
			if err != nil {
				rollbackErr := tx.Rollback()
				if rollbackErr != nil {
					err = errors.Join(err, rollbackErr)
				}
			}
		}()

		batch := make([]*domain.Note, 0, r.cfg.CreateNotesBatchSize)
		for note := range noteChan {
			batch = append(batch, note)
			if len(batch) == r.cfg.CreateNotesBatchSize {
				if err = r.insertBatch(tx, batch, user, noteIDs); err != nil {
					errChan <- err
					return
				}

				batch = batch[:0]
			}
		}

		if len(batch) > 0 {
			if err = r.insertBatch(tx, batch, user, noteIDs); err != nil {
				errChan <- err
				return
			}
		}

		if err = tx.Commit(); err != nil {
			errChan <- err
			return
		}
		close(noteIDs)
	}()
	return noteIDs, errChan
}

func (r *repositoryImpl) insertBatch(tx *postgres.Transaction, batch []*domain.Note, user *domain.User, noteIDsChan chan string) error {
	noteBatch := &pgx.Batch{}
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	for _, note := range batch {
		noteID, err := uuid.GenerateUUID()
		if err != nil {
			return fmt.Errorf("error generating uuid: %w", err)
		}
		note.ID = noteID
		note.CreatedAt = time.Now()
		note.UpdatedAt = time.Now()
		query, args, _ := psql.Insert("notes").
			Columns("id", "user_id", "title", "content", "created_at", "updated_at").
			Values(note.ID, user.ID, note.Title, note.Content, note.CreatedAt.Format(time.RFC3339), note.UpdatedAt.Format(time.RFC3339)).
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

	for _, note := range batch {
		err = r.outboxRepo.Create(tx, user, note)
		if err != nil {
			return fmt.Errorf("error creating note outbox: %w", err)
		}
		noteIDsChan <- note.ID
	}

	return nil
}
