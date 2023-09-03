package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
)

func (r *repositoryImpl) ReadNoteByID(ctx context.Context, user *domain.User, id string) (*domain.Note, error) {
	tx, err := r.db.NewTransaction(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("error creating transaction: %w", err)
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

	query, args, err := psql.Select("id", "title", "content", "created_at", "updated_at").
		From("notes").
		Where(squirrel.Eq{"id": id, "user_id": user.ID()}).
		ToSql()

	if err != nil {
		return nil, err
	}

	var model DBModel
	err = tx.QueryRow(query, args...).Scan(&model.ID, &model.Title, &model.Content, &model.CreatedAt, &model.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	var note *domain.Note
	note, err = noteDBModelToDomain(&model)
	if err != nil {
		return nil, err
	}

	err = r.outboxRepo.FindByID(tx, note)
	if err != nil {
		return nil, fmt.Errorf("error creating note outbox: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("error creating note outbox: %w", err)
	}

	return note, nil
}
