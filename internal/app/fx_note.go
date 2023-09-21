package app

import (
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/auth"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/postgres"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/redis"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/server"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/server/handlers"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/server/middlewares"
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
			server.NewModule(),
		),
		fx.Provide(
			config.NewConfig,
			logger.SetupLogger,
			logger.NewConfig,
		),
		fx.WithLogger(logger.WithZerolog(&log.Logger)),
		fx.Invoke(middlewares.RegisterAuthMiddleware, middlewares.RegisterLoggerMiddleware),
		fx.Invoke(
			handlers.RegisterCreateNote,
			handlers.RegisterReadNote,
			handlers.RegisterUpdateNote,
			handlers.RegisterUpdateNote,
			handlers.RegisterDeleteNote,
			handlers.RegisterSearchNote,
		),
	)
}
