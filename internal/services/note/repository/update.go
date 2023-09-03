package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/customerrors"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
)

func (r *repositoryImpl) UpdateNote(ctx context.Context, user *domain.User, note *domain.Note) error {
	tx, err := r.db.NewTransaction(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("error creating transaction: %w", err)
	}
	defer tx.ConnRelease()
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				err = errors.Join(err, rollbackErr)
			}
		}
	}()

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	query, args, _ := psql.Update("notes").
		Set("title", note.Title()).
		Set("content", note.Content()).
		Set("updated_at", note.UpdatedAt().Format(time.RFC3339)).
		Where(squirrel.Eq{"id": note.ID(), "user_id": user.ID()}).
		Suffix("RETURNING id").
		ToSql()

	rows, err := tx.Query(query, args...)
	if err != nil {
		return fmt.Errorf("trx err: %w", err)
	}
	defer rows.Close()
	noteID := ""
	for rows.Next() {
		if scanErr := rows.Scan(&noteID); scanErr != nil {
			log.Err(scanErr).Msg("can't scan noteID")
			return customerrors.ErrNotFoundNoteID
		}
	}
	if noteID == "" {
		return customerrors.ErrNotFoundNoteID
	}

	err = r.outboxRepo.Update(tx, note)
	if err != nil {
		return fmt.Errorf("error creating note outbox: %w", err)
	}

	return tx.Commit()
}
