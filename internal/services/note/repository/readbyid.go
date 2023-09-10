package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
)

func (r *repositoryImpl) ReadNoteByID(ctx context.Context, user *domain.User, id uint64) (*domain.Note, error) {
	var note *domain.Note
	err := r.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

		query, args, err := psql.Select("note_id", "title", "content", "created_at", "updated_at").
			From("notes").
			Where(squirrel.Eq{"note_id": id, "user_id": user.ID()}).
			ToSql()

		if err != nil {
			return err
		}

		var model DBModel
		err = tx.QueryRow(ctx, query, args...).Scan(&model.ID, &model.Title, &model.Content, &model.CreatedAt, &model.UpdatedAt)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil
			}
			return err
		}

		note, err = noteDBModelToDomain(&model)
		note.SetUserID(user.ID())
		if err != nil {
			return err
		}

		err = r.outboxRepo.FindByID(ctx, tx, note)
		if err != nil {
			return fmt.Errorf("error creating note outbox: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error creating transaction: %w", err)
	}

	return note, nil
}
