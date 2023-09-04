package config

import (
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/auth"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/http"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/postgres"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/rabbitmq"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/circuitbreaker"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/logger"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/services/note"
)

type Config struct {
	Services struct {
		Note *note.Config `yaml:"Note"`
	} `yaml:"Services"`
	Adapters struct {
		HTTP     *http.Config     `yaml:"HTTP"`
		Postgres *postgres.Config `yaml:"Postgres"`
		RabbitMQ *rabbitmq.Config `yaml:"RabbitMQ"`
		Auth     *auth.Config     `yaml:"Auth"`
	} `yaml:"Adapters"`
	Common struct {
		Logger         *logger.Config         `yaml:"Logger"`
		CircuitBreaker *circuitbreaker.Config `yaml:"CircuitBreaker"`
	} `yaml:"Common"`
}
