package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

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
	application, err := app.NewAppNote(ctx)
	if err != nil {
		return err
	}

	go application.Adapters.HTTP.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	application.Adapters.HTTP.Stop()
	return nil
}
