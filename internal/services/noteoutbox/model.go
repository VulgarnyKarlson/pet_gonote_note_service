package noteoutbox

import (
	"fmt"

	uuid "github.com/hashicorp/go-uuid"
)

type NoteOutBoxAction int

const (
	NoteActionNull NoteOutBoxAction = iota
	NoteActionCreated
	NoteActionRead
	NoteActionUpdated
	NoteActionDeleted
	NoteActionSearch
)

func (n NoteOutBoxAction) String() string {
	return [...]string{"Null", "CreateNote", "Read", "UpdateNote", "DeleteNote", "SearchNote"}[n]
}

type NoteOutbox struct {
	ID      int    `json:"id"`
	EventID string `json:"event_id"`
	Action  string `json:"action"`
	UserID  string `json:"user_id"`
	NoteID  string `json:"note_id"`
	Sent    bool   `json:"sent"`
}

func NewNoteOutbox(noteID string, eventType NoteOutBoxAction, userID string) (*NoteOutbox, error) {
	eventID, err := uuid.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf("error generating uuid: %w", err)
	}

	ent := &NoteOutbox{
		EventID: eventID,
		Action:  eventType.String(),
		UserID:  userID,
		NoteID:  noteID,
		Sent:    false,
	}

	return ent, nil
}
