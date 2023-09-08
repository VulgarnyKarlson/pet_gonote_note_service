package main

import (
	"github.com/rs/zerolog/log"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/app"
)

func main() {
	log.Info().Msgf("Starting service")
	fxApp := app.NewStatsSenderApp()
	fxApp.Run()
}
