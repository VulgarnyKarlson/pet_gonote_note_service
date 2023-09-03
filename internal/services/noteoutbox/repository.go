package noteoutbox

import (
	"fmt"

	"github.com/Masterminds/squirrel"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/postgres"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
)

var nullNote = &domain.Note{}

type Repository interface {
	Create(tx *postgres.Transaction, note *domain.Note) (err error)
	Update(tx *postgres.Transaction, note *domain.Note) (err error)
	Delete(tx *postgres.Transaction, note *domain.Note) (err error)
	FindByID(tx *postgres.Transaction, note *domain.Note) (err error)
	Search(tx *postgres.Transaction, user *domain.User) (err error)
	GetAllOutbox(tx *postgres.Transaction) (notesOutbox []*NoteOutbox, err error)
	MarkAsSent(tx *postgres.Transaction, notesOutbox *NoteOutbox) error
}

type repositoryImpl struct {
	db *postgres.Pool
}

func NewRepository(db *postgres.Pool) Repository {
	nullNote.SetID("b40fae8f-7689-a545-d431-14f6374a79cc")
	return &repositoryImpl{db: db}
}

func (r *repositoryImpl) Create(tx *postgres.Transaction, note *domain.Note) (err error) {
	return r.insert(tx, note, NoteActionCreated)
}

func (r *repositoryImpl) Update(tx *postgres.Transaction, note *domain.Note) (err error) {
	return r.insert(tx, note, NoteActionUpdated)
}

func (r *repositoryImpl) Delete(tx *postgres.Transaction, note *domain.Note) (err error) {
	return r.insert(tx, note, NoteActionDeleted)
}

func (r *repositoryImpl) FindByID(tx *postgres.Transaction, note *domain.Note) (err error) {
	return r.insert(tx, note, NoteActionRead)
}

func (r *repositoryImpl) Search(tx *postgres.Transaction, user *domain.User) (err error) {
	tmp := nullNote.Copy()
	tmp.SetUserID(user.ID())
	return r.insert(tx, tmp, NoteActionSearch)
}

func (r *repositoryImpl) GetAllOutbox(tx *postgres.Transaction) (notesOutbox []*NoteOutbox, err error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	query, args, err := psql.Select("id", "event_id", "action", "user_id", "note_id", "sent").
		From("notes_outbox").
		Where(squirrel.Eq{"sent": false}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("SQL build error: %w", err)
	}

	rows, err := tx.Query(query, args...)
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

func (r *repositoryImpl) MarkAsSent(tx *postgres.Transaction, notesOutbox *NoteOutbox) error {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	query, args, err := psql.Update("notes_outbox").
		Set("sent", true).
		Where(squirrel.Eq{"id": notesOutbox.ID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("SQL build error: %w", err)
	}

	_, err = tx.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("trx err: %w", err)
	}

	return nil
}

func (r *repositoryImpl) insert(tx *postgres.Transaction, note *domain.Note, actionType NoteOutBoxAction) (err error) {
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

	_, err = tx.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("trx err: %w", err)
	}
	return nil
}
