package app

import (
	"context"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/postgres"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/rabbitmq"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/logger"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/config"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/services/noteoutbox"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/services/outboxproducer"
)

type StatsSenderApplication struct {
	Config   *config.Config
	Adapters struct {
		Postgres *postgres.Pool
		RabbitMQ *rabbitmq.Publisher
	}
	Services struct {
		OutBoxProducer *outboxproducer.OutBoxProducer
	}
}

func NewAppStatsSender(ctx context.Context) (app *StatsSenderApplication, err error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	logger.SetupLogger(cfg.Common.Logger)
	pgPool, err := postgres.New(ctx, cfg.Adapters.Postgres)
	if err != nil {
		return nil, err
	}
	noteOutBoxRepo := noteoutbox.NewRepository(pgPool)
	msgProducer, err := rabbitmq.NewPublisher(cfg.Adapters.RabbitMQ)
	if err != nil {
		return nil, err
	}

	outBoxProducer := outboxproducer.NewOutBoxProducer(pgPool, noteOutBoxRepo, msgProducer)
	app = new(StatsSenderApplication)
	app.Config = cfg
	app.Adapters.Postgres = pgPool
	app.Adapters.RabbitMQ = msgProducer
	app.Services.OutBoxProducer = outBoxProducer
	return app, nil
}
