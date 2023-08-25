package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/rabbitmq"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/services/outboxproducer"

	"github.com/rs/zerolog/log"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/postgres"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/logger"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/config"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/services/noteoutbox"
)

func main() {
	if err := mainWithErr(); err != nil {
		log.Fatal().Err(err).Msg("error while starting service")
	}
}

func mainWithErr() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	logger.SetupLogger(cfg.Common.Logger)
	log.Info().Msgf("Starting service")
	ctx := context.Background()
	pgPool, err := postgres.New(ctx, cfg.Adapters.Postgres)
	if err != nil {
		return err
	}
	noteOutBoxRepo := noteoutbox.NewRepository(pgPool)
	msgProducer, err := rabbitmq.NewPublisher(cfg.Adapters.RabbitMQ)
	if err != nil {
		return err
	}

	outBoxProducer := outboxproducer.NewOutBoxProducer(pgPool, noteOutBoxRepo, msgProducer)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	for {
		select {
		case <-ctx.Done():
			log.Info().Msgf("Context done")
			return nil
		case <-c:
			log.Info().Msgf("Terminating service")
			return nil
		case <-time.After(1 * time.Second):
			count, err := outBoxProducer.Produce(ctx)
			if err != nil {
				return err
			}
			log.Info().Msgf("Produced %d messages", count)
		}
	}
}
