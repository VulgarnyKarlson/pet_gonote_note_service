package repository

import (
	"context"
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
			Where(squirrel.Eq{"id": note.ID(), "user_id": note.UserID()}).
			Suffix("RETURNING id").
			ToSql()

		rows, err := tx.Query(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("trx err: %w", err)
		}
		defer rows.Close()
		noteID := ""
		for rows.Next() {
			if scanErr := rows.Scan(&noteID); scanErr != nil {
				r.logger.Err(scanErr).Msg("can't scan noteID")
				return customerrors.ErrNotFoundNoteID
			}
		}
		if noteID == "" {
			return customerrors.ErrNotFoundNoteID
		}

		err = r.outboxRepo.Update(ctx, tx, note)
		if err != nil {
			return fmt.Errorf("error creating note outbox: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("error creating transaction: %w", err)
	}
	return nil
}
