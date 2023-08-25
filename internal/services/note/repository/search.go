package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"

	"github.com/Masterminds/squirrel"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
)

var searchNoteOutbox = &domain.Note{ID: "b40fae8f-7689-a545-d431-14f6374a79cc"}

func (r *repositoryImpl) Search(
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
	queryBuilder := psql.Select("id", "title", "content", "created_at", "updated_at").From("notes")

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

	out := make([]*domain.Note, 0)
	for rows.Next() {
		var note domain.Note
		if err = rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error search scan note: %w", err)
		}
		out = append(out, &note)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	err = r.outboxRepo.Search(tx, user, searchNoteOutbox)
	if err != nil {
		return nil, fmt.Errorf("error creating note outbox: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("error creating note outbox: %w", err)
	}

	return out, nil
}
