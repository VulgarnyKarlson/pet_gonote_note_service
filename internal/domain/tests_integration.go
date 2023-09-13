package domain

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/postgres"
)

type TestIntegration struct {
	t         *testing.M
	l         zerolog.Logger
	pool      *dockertest.Pool
	resources []*dockertest.Resource
}

func NewTestIntegration(t *testing.M) *TestIntegration {
	return &TestIntegration{t: t, l: log.Level(zerolog.DebugLevel)}
}

func (t *TestIntegration) Setup() {
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.l.Fatal().Err(err).Msg("Could not construct pool")
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		t.l.Fatal().Err(err).Msg("Could not connect to Docker")
	}
	t.l.Info().Msg("Connected to Docker")

	t.pool = pool
}

func (t *TestIntegration) Teardown() {
	for _, resource := range t.resources {
		if err := t.pool.Purge(resource); err != nil {
			t.l.Fatal().Err(err).Msg("Could not purge resource")
		}
	}
}

var pgConfig *postgres.Config

func (t *TestIntegration) PostgresUP() {
	resource, err := t.pool.Run("postgres", "14.3", []string{
		"POSTGRES_DB=note_service", "POSTGRES_USER=postgres", "POSTGRES_PASSWORD=1234",
	})
	if err != nil {
		t.l.Fatal().Err(err).Msg("Could not start postgres")
	}
	t.l.Info().Msg("Postgres started")

	if err = t.pool.Retry(func() error {
		host, port := t.getHostPort(resource, "5432/tcp")
		pgConfig = &postgres.Config{
			Host:     host,
			Port:     port,
			UserName: "postgres",
			Password: "1234",
			DBName:   "note_service",
		}
		postgresURL := fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=disable",
			pgConfig.UserName, pgConfig.Password, pgConfig.Host, pgConfig.Port, pgConfig.DBName,
		)
		cfg, errPG := pgx.ParseConfig(postgresURL)
		if errPG != nil {
			return fmt.Errorf("could not parse postgres url: %w", errPG)
		}
		con := stdlib.OpenDB(*cfg)
		err = con.Ping()
		if err != nil {
			return err
		}
		t.l.Info().Msg("Postgres connected")

		_ = os.Setenv("POSTGRES_URL", postgresURL)
		abs, _ := filepath.Abs("../")
		migrationsPath := strings.Split(abs, "internal")[0] + "db/migrations"

		err = goose.Up(con, migrationsPath)
		if err != nil {
			t.l.Err(err).Msg("Could not migrate postgres")
		}
		return err
	}); err != nil {
		t.l.Fatal().Err(err).Msg("Could not connect to postgres")
	} else {
		t.l.Info().Msg("Postgres migrated")
	}

	t.resources = append(t.resources, resource)
}

func TestIntegrationGetPostgresConfig() *postgres.Config {
	return pgConfig
}

func (t *TestIntegration) getHostPort(resource *dockertest.Resource, id string) (host string, port int) {
	var err error
	dockerURL := os.Getenv("DOCKER_HOST")
	if dockerURL == "" {
		var portStr string
		hostPort := resource.GetHostPort(id)
		host, portStr, err = net.SplitHostPort(hostPort)
		if err != nil {
			t.l.Fatal().Err(err).Msg("Could not split host port")
		}
		port, _ = strconv.Atoi(portStr)

		return host, port
	}
	u, err := url.Parse(dockerURL)
	if err != nil {
		t.l.Fatal().Err(err).Msg("Could not parse docker url")
	}

	port, err = strconv.Atoi(resource.GetPort(id))
	if err != nil {
		t.l.Fatal().Err(err).Msg("Could not convert port to int")
	}

	return u.Hostname(), port
}
