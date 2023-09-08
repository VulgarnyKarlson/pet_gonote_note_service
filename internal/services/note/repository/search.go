package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
)

func (r *repositoryImpl) SearchNote(
	ctx context.Context,
	user *domain.User,
	criteria *domain.SearchCriteria,
) ([]*domain.Note, error) {
	tx, err := r.db.NewTransaction(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("error search creating transaction: %w", err)
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
	queryBuilder := psql.Select("id", "user_id", "title", "content", "created_at", "updated_at").From("notes")
	queryBuilder = queryBuilder.Where(squirrel.Eq{"user_id": user.ID()})
	if criteria.Title != "" {
		queryBuilder = queryBuilder.Where("title LIKE ?", fmt.Sprintf("%%%s%%", criteria.Title))
	}

	if criteria.Content != "" {
		queryBuilder = queryBuilder.Where("content LIKE ?", fmt.Sprintf("%%%s%%", criteria.Content))
	}

	if !criteria.FromDate.IsZero() {
		queryBuilder = queryBuilder.Where("created_at >= ?", criteria.FromDate)
	}

	if !criteria.ToDate.IsZero() {
		queryBuilder = queryBuilder.Where("created_at <= ?", criteria.ToDate)
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error search creating query: %w", err)
	}

	rows, err := tx.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error search notes: %w", err)
	}
	defer rows.Close()

	out := make([]*DBModel, 0)
	for rows.Next() {
		var note DBModel
		if err = rows.Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error search scan note: %w", err)
		}
		out = append(out, &note)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	err = r.outboxRepo.Search(tx, user)
	if err != nil {
		return nil, fmt.Errorf("error creating note outbox: %w", err)
	}

	domainOut := make([]*domain.Note, len(out))
	for i, note := range out {
		domainOut[i], err = noteDBModelToDomain(note)
		if err != nil {
			return nil, fmt.Errorf("error db -> domain note: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("error creating note outbox: %w", err)
	}

	return domainOut, nil
}
