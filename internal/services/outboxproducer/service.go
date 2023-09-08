package outboxproducer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	amqp "github.com/rabbitmq/amqp091-go"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/postgres"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/rabbitmq"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/services/noteoutbox"
)

type OutBoxProducer struct {
	Producer   *rabbitmq.Publisher
	db         *postgres.Pool
	outboxRepo noteoutbox.Repository
}

type outboxMessage struct {
	UserID    string `json:"user_id"`
	Action    string `json:"action"`
	NoteID    string `json:"note_id"`
	Timestamp string `json:"timestamp"`
}

func NewOutBoxProducer(db *postgres.Pool, outboxRepo noteoutbox.Repository, publisher *rabbitmq.Publisher) *OutBoxProducer {
	return &OutBoxProducer{Producer: publisher, db: db, outboxRepo: outboxRepo}
}

func (o *OutBoxProducer) Produce(ctx context.Context) (count int, err error) {
	txDB, err := o.db.NewTransaction(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("error creating db transaction: %w", err)
	}
	txPB, err := o.Producer.Tx()
	if err != nil {
		return 0, fmt.Errorf("error creating producer transaction: %w", err)
	}
	defer func(txPB *amqp.Channel) {
		closeErr := txPB.Close()
		if closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}(txPB)
	defer func() {
		if err != nil {
			rollbackErr := txPB.TxRollback()
			if rollbackErr != nil {
				err = errors.Join(err, rollbackErr)
			}
		}
	}()

	defer txDB.ConnRelease()
	defer func() {
		if err != nil {
			rollbackErr := txDB.Rollback()
			if rollbackErr != nil {
				err = errors.Join(err, rollbackErr)
			}
		}
	}()

	notesOutbox, err := o.outboxRepo.GetAllOutbox(txDB)
	if err != nil {
		return 0, err
	}

	for _, noteOutbox := range notesOutbox {
		b, errMarshal := json.Marshal(outboxMessage{
			UserID:    noteOutbox.UserID,
			Action:    noteOutbox.Action,
			NoteID:    noteOutbox.NoteID,
			Timestamp: time.Now().Format(time.RFC3339),
		})
		if errMarshal != nil {
			return 0, err
		}
		err = o.Producer.Publish(ctx, txPB, b)
		if err != nil {
			return 0, err
		}
		err = o.outboxRepo.MarkAsSent(txDB, noteOutbox)
		if err != nil {
			return 0, err
		}
	}

	err = txPB.TxCommit()
	if err != nil {
		return 0, err
	}
	err = txDB.Commit()
	if err != nil {
		return 0, err
	}

	return len(notesOutbox), nil
}
