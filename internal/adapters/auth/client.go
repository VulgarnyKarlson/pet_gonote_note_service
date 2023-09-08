package auth

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/redis"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/circuitbreaker"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/customerrors"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/proto"
)

type Client interface {
	SetProtoService(service any)
	ValidateToken(ctx context.Context, token string) (*ValidateTokenResponse, error)
	Connect() error
	Close() error
}

type ValidateTokenResponse struct {
	User  *domain.User
	Valid bool
}

type ClientImpl struct {
	conn           *grpc.ClientConn
	service        proto.AuthServiceClient
	logger         *zerolog.Logger
	config         *Config
	circuitbreaker circuitbreaker.CircuitBreaker
	storage        redis.Client
}

func NewWrapper(logger *zerolog.Logger, cnf *Config, cb circuitbreaker.CircuitBreaker, storage redis.Client) (Client, error) {
	return &ClientImpl{
		config:         cnf,
		logger:         logger,
		circuitbreaker: cb,
		storage:        storage,
	}, nil
}

func (c *ClientImpl) SetProtoService(service any) {
	c.service = service.(proto.AuthServiceClient)
}

func (c *ClientImpl) Connect() error {
	conn, err := grpc.Dial(c.config.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to auth service: %w", err)
	}
	c.conn = conn
	c.service = proto.NewAuthServiceClient(conn)
	c.logger.Info().Msg("connected to auth service")
	return nil
}

func (c *ClientImpl) ValidateToken(ctx context.Context, token string) (*ValidateTokenResponse, error) {
	if res, err := c.get(ctx, token); err == nil {
		return res, nil
	}
	err := c.circuitbreaker.Attempt()
	if err != nil {
		return nil, err
	}
	resp, err := c.service.ValidateToken(ctx, &proto.ValidateTokenRequest{Token: token})
	if err != nil {
		c.circuitbreaker.Fail(err)
		return nil, customerrors.ErrAuthServiceError
	}
	c.circuitbreaker.Success()
	validateTokenResponse := &ValidateTokenResponse{Valid: resp.Valid}
	if resp.Valid {
		validateTokenResponse.User = domain.NewUser(resp.User.GetId(), resp.User.GetUsername())
	}
	err = c.store(ctx, token, validateTokenResponse)
	if err != nil {
		c.logger.Err(err).Msg("failed to store token")
		return nil, customerrors.ErrAuthServiceError
	}
	return validateTokenResponse, nil
}

func (c *ClientImpl) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
