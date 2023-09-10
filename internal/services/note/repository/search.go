package repository

import (
	"context"
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
	var domainOut []*domain.Note
	err := r.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
		queryBuilder := psql.Select("note_id", "user_id", "title", "content", "created_at", "updated_at").From("notes")
		queryBuilder = queryBuilder.Where(squirrel.Eq{"user_id": user.ID()})
		if criteria.Title != "" {
			queryBuilder = queryBuilder.Where("fts @@ plainto_tsquery('english', ?)", criteria.Title)
		}

		if criteria.Content != "" {
			queryBuilder = queryBuilder.Where("fts @@ plainto_tsquery('english', ?)", criteria.Content)
		}

		if !criteria.FromDate.IsZero() {
			queryBuilder = queryBuilder.Where("created_at >= ?", criteria.FromDate)
		}

		if !criteria.ToDate.IsZero() {
			queryBuilder = queryBuilder.Where("created_at <= ?", criteria.ToDate)
		}

		query, args, err := queryBuilder.ToSql()
		if err != nil {
			return fmt.Errorf("error search creating query: %w", err)
		}

		rows, err := tx.Query(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("error search notes: %w", err)
		}
		defer rows.Close()

		out := make([]*DBModel, 0)
		for rows.Next() {
			var note DBModel
			if err = rows.Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt); err != nil {
				return fmt.Errorf("error search scan note: %w", err)
			}
			out = append(out, &note)
		}

		if err = rows.Err(); err != nil {
			return err
		}

		err = r.outboxRepo.Search(ctx, tx, user)
		if err != nil {
			return fmt.Errorf("error creating note outbox: %w", err)
		}

		domainOut = make([]*domain.Note, len(out))
		for i, note := range out {
			domainOut[i], err = noteDBModelToDomain(note)
			if err != nil {
				return fmt.Errorf("error db -> domain note: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error search creating transaction: %w", err)
	}

	return domainOut, nil
}
