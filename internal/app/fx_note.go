package app

import (
	"context"

	"github.com/rs/zerolog/log"
	"go.uber.org/fx"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/auth"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/http"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/http/handlers"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/postgres"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/redis"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/circuitbreaker"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/logger"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/config"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/services/note"
	noteRepo "gitlab.karlson.dev/individual/pet_gonote/note_service/internal/services/note/repository"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/services/noteoutbox"
)

func NewNoteApp() *fx.App {
	return fx.New(
		fx.Options(
			circuitbreaker.NewModule(),
			redis.NewModule(),
			auth.NewModule(),
			postgres.NewModule(),
			noteoutbox.NewModule(),
			noteRepo.NewModule(),
			note.NewModule(),
			http.NewModule(),
			handlers.NewModule(),
		),
		fx.Provide(
			config.NewConfig,
			logger.SetupLogger,
			logger.NewConfig,
		),
		fx.WithLogger(logger.WithZerolog(&log.Logger)),
		fx.Invoke(initHTTPEndpoints),
	)
}

func initHTTPEndpoints(lx fx.Lifecycle, h *handlers.NoteHandlers, n *http.Server) {
	endpoints := []http.Endpoint{
		{Method: "POST", Path: "/create", Auth: true, Handler: h.CreateNote},
		{Method: "GET", Path: "/read", Auth: true, Handler: h.ReadNoteByID},
		{Method: "POST", Path: "/update", Auth: true, Handler: h.UpdateNote},
		{Method: "POST", Path: "/delete", Auth: true, Handler: h.DeleteNoteByID},
		{Method: "GET", Path: "/search", Auth: true, Handler: h.SearchNote},
	}

	for _, e := range endpoints {
		n.AddEndpoint(e)
	}

	lx.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				err := n.Run()
				if err != nil {
					log.Fatal().Err(err).Msgf("Error while starting http server")
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return n.Stop()
		},
	})
}
