package server

import (
	"context"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"adapters/server",
		fx.Provide(
			NewConfig,
			NewServer,
			mux.NewRouter,
		),
		fx.Invoke(func(lx fx.Lifecycle, s Server) {
			lx.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						err := s.Run()
						if err != nil {
							log.Fatal().Err(err).Msgf("Error while starting http server")
						}
					}()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					return s.Stop()
				},
			})
		}),
		fx.Decorate(func(log *zerolog.Logger) *zerolog.Logger {
			lg := log.With().Str("module", "server").Logger()
			return &lg
		}),
	)
}
