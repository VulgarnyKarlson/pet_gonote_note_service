package config

import (
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/auth"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/http"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/postgres"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/rabbitmq"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/logger"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/services/note"
)

type Config struct {
	Common struct {
		Logger *logger.Config
	}
	Services struct {
		Note *note.Config
	}
	Adapters struct {
		HTTP     *http.Config
		Postgres *postgres.Config
		RabbitMQ *rabbitmq.Config
		Auth     *auth.Config
	}
}
