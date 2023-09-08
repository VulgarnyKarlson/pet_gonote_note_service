package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
)

func (r *repositoryImpl) DeleteNote(ctx context.Context, user *domain.User, id string) (bool, error) {
	tx, err := r.db.NewTransaction(ctx, pgx.TxOptions{})
	if err != nil {
		return false, fmt.Errorf("error creating transaction: %w", err)
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

	query, args, err := psql.Delete("notes").
		Where(squirrel.Eq{"id": id, "user_id": user.ID()}).
		ToSql()

	if err != nil {
		return false, err
	}

	ct, err := tx.Exec(query, args...)
	if err != nil {
		return false, fmt.Errorf("error delete note: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return false, nil
	}

	note := new(domain.Note)
	note.SetID(id)
	note.SetUserID(user.ID())
	err = r.outboxRepo.Delete(tx, note)
	if err != nil {
		return false, fmt.Errorf("error creating note outbox: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return false, fmt.Errorf("error creating note outbox: %w", err)
	}

	return true, nil
}
