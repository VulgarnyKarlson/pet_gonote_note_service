package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/customerrors"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
)

func (r *repositoryImpl) UpdateNote(ctx context.Context, user *domain.User, note *domain.Note) error {
	err := r.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		note.SetUserID(user.ID())
		psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
		query, args, _ := psql.Update("notes").
			Set("title", note.Title()).
			Set("content", note.Content()).
			Set("updated_at", note.UpdatedAt().Format(time.RFC3339)).
			Where(squirrel.Eq{"note_id": note.ID(), "user_id": note.UserID()}).
			Suffix("RETURNING note_id").
			ToSql()

		rows, err := tx.Query(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("trx err: %w", err)
		}
		defer rows.Close()
		var noteID uint64
		for rows.Next() {
			if scanErr := rows.Scan(&noteID); scanErr != nil {
				return customerrors.ErrNotFoundNoteID
			}
		}
		if noteID == 0 {
			return customerrors.ErrNotFoundNoteID
		}

		err = r.outboxRepo.Update(ctx, tx, note)
		if err != nil {
			return fmt.Errorf("error creating note outbox: %w", err)
		}

		return nil
	})
	if err != nil {
		return errors.Join(customerrors.ErrRepositoryError, err)
	}
	return nil
}
