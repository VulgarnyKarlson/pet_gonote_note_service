package outboxproducer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/rabbitmq"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/services/noteoutbox"
)

type OutBoxProducer struct {
	Producer   *rabbitmq.Publisher
	db         *pgxpool.Pool
	outboxRepo noteoutbox.Repository
}

type outboxMessage struct {
	UserID    uint64 `json:"user_id"`
	Action    string `json:"action"`
	NoteID    uint64 `json:"note_id"`
	Timestamp string `json:"timestamp"`
}

func NewOutBoxProducer(db *pgxpool.Pool, outboxRepo noteoutbox.Repository, publisher *rabbitmq.Publisher) *OutBoxProducer {
	return &OutBoxProducer{Producer: publisher, db: db, outboxRepo: outboxRepo}
}

func (o *OutBoxProducer) Produce(ctx context.Context) (int, error) {
	var count int
	txPB, errPB := o.Producer.Tx()
	if errPB != nil {
		return 0, fmt.Errorf("error creating producer transaction: %w", errPB)
	}
	defer func(txPB *amqp.Channel) {
		closeErr := txPB.Close()
		if closeErr != nil {
			errPB = errors.Join(errPB, closeErr)
		}
	}(txPB)

	errTx := o.db.BeginTxFunc(ctx, pgx.TxOptions{}, func(txDB pgx.Tx) error {
		notesOutbox, err := o.outboxRepo.GetAllOutbox(ctx, txDB)
		if err != nil {
			return err
		}
		count = len(notesOutbox)

		for _, noteOutbox := range notesOutbox {
			b, errMarshal := json.Marshal(outboxMessage{
				UserID:    noteOutbox.UserID,
				Action:    noteOutbox.Action,
				NoteID:    noteOutbox.NoteID,
				Timestamp: time.Now().Format(time.RFC3339),
			})
			if errMarshal != nil {
				return errMarshal
			}
			err = o.Producer.Publish(ctx, txPB, b)
			if err != nil {
				return err
			}
			err = o.outboxRepo.MarkAsSent(ctx, txDB, noteOutbox)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if errTx != nil {
		errPB = txPB.TxRollback()
		if errPB != nil {
			errTx = errors.Join(errTx, errPB)
		}
		return 0, fmt.Errorf("error creating db transaction: %w", errTx)
	}

	return count, errTx
}
