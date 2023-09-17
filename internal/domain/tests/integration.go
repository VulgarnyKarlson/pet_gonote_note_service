package tests

import (
	"net"
	"net/url"
	"os"
	"strconv"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type TestIntegration struct {
	t         *testing.M
	l         zerolog.Logger
	pool      *dockertest.Pool
	resources []*dockertest.Resource
	Configs   *configStorage
}

func NewTestIntegration(t *testing.M) *TestIntegration {
	return &TestIntegration{t: t, l: log.Level(zerolog.DebugLevel), Configs: newConfigStorage()}
}

func (t *TestIntegration) RunServices(enabledServices ...TestIntegrationService) {
	for _, service := range enabledServices {
		switch service {
		case TestIntegrationServicePostgres:
			t.PostgresUP()
		case TestIntegrationServiceRedis:
			t.RedisUP()
		case TestIntegrationServiceRabbitMQ:
			t.RabbitMQUP()
		}
	}
}

func (t *TestIntegration) Setup() {
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.l.Panic().Err(err).Msg("Could not construct pool")
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		t.l.Panic().Err(err).Msg("Could not connect to Docker")
	}
	t.l.Info().Msg("Connected to Docker")

	t.pool = pool
}

func (t *TestIntegration) Teardown() {
	for _, resource := range t.resources {
		if err := t.pool.Purge(resource); err != nil {
			t.l.Panic().Err(err).Msg("Could not purge resource")
		}
	}
}

func (t *TestIntegration) RedisUP() {
	t.l.Info().Msg("Redis started")
}

func (t *TestIntegration) RabbitMQUP() {
	t.l.Info().Msg("RabbitMQ started")
}

func getHostPort(resource *dockertest.Resource, id string) (host string, port int) {
	l := log.Level(zerolog.DebugLevel)
	var err error
	dockerURL := os.Getenv("DOCKER_HOST")
	if dockerURL == "" {
		var portStr string
		hostPort := resource.GetHostPort(id)
		host, portStr, err = net.SplitHostPort(hostPort)
		if err != nil {
			l.Panic().Err(err).Msg("Could not split host port")
		}
		port, _ = strconv.Atoi(portStr)

		return host, port
	}
	u, err := url.Parse(dockerURL)
	if err != nil {
		l.Panic().Err(err).Msg("Could not parse docker url")
	}

	port, err = strconv.Atoi(resource.GetPort(id))
	if err != nil {
		l.Panic().Err(err).Msg("Could not convert port to int")
	}

	return u.Hostname(), port
}
