package repository

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
)

func (r *repositoryImpl) DeleteNote(ctx context.Context, user *domain.User, id uint64) (bool, error) {
	var isDeleted bool
	err := r.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

		query, args, err := psql.Delete("notes").
			Where(squirrel.Eq{"note_id": id, "user_id": user.ID()}).
			ToSql()

		if err != nil {
			return err
		}

		ct, err := tx.Exec(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("error delete note: %w", err)
		}

		if ct.RowsAffected() == 0 {
			return nil
		}
		isDeleted = true

		note := new(domain.Note)
		note.SetID(id)
		note.SetUserID(user.ID())
		err = r.outboxRepo.Delete(ctx, tx, note)
		if err != nil {
			return fmt.Errorf("error creating note outbox: %w", err)
		}

		return nil
	})
	if err != nil {
		return false, fmt.Errorf("error creating transaction: %w", err)
	}

	return isDeleted, nil
}
