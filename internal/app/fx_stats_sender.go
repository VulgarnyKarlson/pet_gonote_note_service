package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/postgres"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/rabbitmq"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/logger"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/config"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/services/noteoutbox"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/services/outboxproducer"
	"go.uber.org/fx"
)

func NewStatsSenderApp() *fx.App {
	return fx.New(
		fx.Options(
			postgres.NewModule(),
			rabbitmq.NewModule(),
			noteoutbox.NewModule(),
			outboxproducer.NewModule(),
		),
		fx.Provide(
			config.NewConfig,
			logger.SetupLogger,
			logger.NewConfig,
		),
		fx.WithLogger(logger.WithZerolog(&log.Logger)),
		fx.Invoke(initTimer),
	)
}

func initTimer(_ fx.Lifecycle, n *outboxproducer.OutBoxProducer) {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		for {
			select {
			case <-c:
				log.Info().Msgf("Terminating service")
				return
			case <-time.After(1 * time.Second):
				count, err := n.Produce(context.TODO())
				if err != nil {
					log.Error().Err(err).Msgf("Error while producing messages")
					return
				}
				log.Info().Msgf("Produced %d messages", count)
			}
		}
	}()
}
