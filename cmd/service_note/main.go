package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/auth"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/http"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/http/handlers"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/postgres"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/logger"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/config"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/services/note"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/services/note/repository"
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
	noteRepo := repository.NewRepository(
		&repository.Config{CreateNotesBatchSize: cfg.Services.Note.CreateNotesBatchSize},
		pgPool, noteOutBoxRepo,
	)
	noteService := note.NewService(cfg.Services.Note, noteRepo)
	noteHandlers := handlers.New(noteService)
	authService := auth.NewWrapper(cfg.Adapters.Auth)
	httpServer := http.NewServer(cfg.Adapters.HTTP, authService, noteHandlers)
	go httpServer.Run()

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	httpServer.Stop()
	return nil
}
