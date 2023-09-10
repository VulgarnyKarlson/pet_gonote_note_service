package noteoutbox

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
)

var nullNote = &domain.Note{}

type Repository interface {
	Create(ctx context.Context, tx pgx.Tx, note *domain.Note) (err error)
	Update(ctx context.Context, tx pgx.Tx, note *domain.Note) (err error)
	Delete(ctx context.Context, tx pgx.Tx, note *domain.Note) (err error)
	FindByID(ctx context.Context, tx pgx.Tx, note *domain.Note) (err error)
	Search(ctx context.Context, tx pgx.Tx, user *domain.User) (err error)
	GetAllOutbox(ctx context.Context, tx pgx.Tx) (notesOutbox []*NoteOutbox, err error)
	MarkAsSent(ctx context.Context, tx pgx.Tx, notesOutbox *NoteOutbox) error
}

type repositoryImpl struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	nullNote.SetID("b40fae8f-7689-a545-d431-14f6374a79cc")
	return &repositoryImpl{db: db}
}

func (r *repositoryImpl) Create(ctx context.Context, tx pgx.Tx, note *domain.Note) (err error) {
	return r.insert(ctx, tx, note, NoteActionCreated)
}

func (r *repositoryImpl) Update(ctx context.Context, tx pgx.Tx, note *domain.Note) (err error) {
	return r.insert(ctx, tx, note, NoteActionUpdated)
}

func (r *repositoryImpl) Delete(ctx context.Context, tx pgx.Tx, note *domain.Note) (err error) {
	return r.insert(ctx, tx, note, NoteActionDeleted)
}

func (r *repositoryImpl) FindByID(ctx context.Context, tx pgx.Tx, note *domain.Note) (err error) {
	return r.insert(ctx, tx, note, NoteActionRead)
}

func (r *repositoryImpl) Search(ctx context.Context, tx pgx.Tx, user *domain.User) (err error) {
	tmp := nullNote.Copy()
	tmp.SetUserID(user.ID())
	return r.insert(ctx, tx, tmp, NoteActionSearch)
}

func (r *repositoryImpl) GetAllOutbox(ctx context.Context, tx pgx.Tx) (notesOutbox []*NoteOutbox, err error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	query, args, err := psql.Select("id", "event_id", "action", "user_id", "note_id", "sent").
		From("notes_outbox").
		Where(squirrel.Eq{"sent": false}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("SQL build error: %w", err)
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("trx err: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var noteOutbox NoteOutbox
		err = rows.Scan(&noteOutbox.ID, &noteOutbox.EventID, &noteOutbox.Action, &noteOutbox.UserID, &noteOutbox.NoteID, &noteOutbox.Sent)
		if err != nil {
			return nil, fmt.Errorf("trx err: %w", err)
		}

		notesOutbox = append(notesOutbox, &noteOutbox)
	}

	return notesOutbox, nil
}

func (r *repositoryImpl) MarkAsSent(ctx context.Context, tx pgx.Tx, notesOutbox *NoteOutbox) error {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	query, args, err := psql.Update("notes_outbox").
		Set("sent", true).
		Where(squirrel.Eq{"id": notesOutbox.ID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("SQL build error: %w", err)
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("trx err: %w", err)
	}

	return nil
}

func (r *repositoryImpl) insert(ctx context.Context, tx pgx.Tx, note *domain.Note, actionType NoteOutBoxAction) (err error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	noteOutbox, err := NewNoteOutbox(note.ID(), actionType, note.UserID())
	if err != nil {
		return fmt.Errorf("error creating note outbox: %w", err)
	}
	noteOutbox.UserID = note.UserID()

	query, args, err := psql.Insert("notes_outbox").
		Columns("event_id", "action", "user_id", "note_id").
		Values(noteOutbox.EventID, noteOutbox.Action, noteOutbox.UserID, noteOutbox.NoteID).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return fmt.Errorf("SQL build error: %w", err)
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("trx err: %w", err)
	}
	return nil
}
