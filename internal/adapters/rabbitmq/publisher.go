package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	conn   *amqp.Connection
	queue  *amqp.Queue
	config *Config
}

func NewPublisher(cfg *Config) (*Publisher, error) {
	return &Publisher{config: cfg}, nil
}

func (p *Publisher) Open() error {
	connString := fmt.Sprintf("amqp://%s:%s@%s:%d/", p.config.UserName, p.config.Password, p.config.Host, p.config.Port)
	conn, err := amqp.Dial(connString)
	if err != nil {
		return err
	}
	p.conn = conn
	return nil
}

func (p *Publisher) Close() error {
	return p.conn.Close()
}

func (p *Publisher) Tx() (ch *amqp.Channel, err error) {
	ch, err = p.conn.Channel()
	if err != nil {
		return nil, err
	}
	q, err := ch.QueueDeclare(
		p.config.QueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	p.queue = &q
	err = ch.Tx()
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func (p *Publisher) Publish(ctx context.Context, ch *amqp.Channel, message []byte) error {
	return ch.PublishWithContext(ctx, "", p.queue.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        message,
	})
}
