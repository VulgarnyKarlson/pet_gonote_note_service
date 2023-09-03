package main

import (
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/app"

	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msgf("Starting service")
	fxApp := app.NewStatsSenderApp()
	fxApp.Run()
}
