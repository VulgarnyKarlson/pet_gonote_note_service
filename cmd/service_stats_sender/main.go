package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/app"

	"github.com/rs/zerolog/log"
)

func main() {
	if err := mainWithErr(); err != nil {
		log.Fatal().Err(err).Msg("error while starting service")
	}
}

func mainWithErr() error {
	log.Info().Msgf("Starting service")
	ctx := context.Background()
	application, err := app.NewAppStatsSender(ctx)
	if err != nil {
		return err
	}
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
			count, err := application.Services.OutBoxProducer.Produce(ctx)
			if err != nil {
				return err
			}
			log.Info().Msgf("Produced %d messages", count)
		}
	}
}
