package auth

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/proto"
	"google.golang.org/grpc/credentials/insecure"
)

type Client interface {
	ValidateToken(ctx context.Context, token string) (*ValidateTokenResponse, error)
}

type ValidateTokenResponse struct {
	User  *domain.User
	Valid bool
}

type ClientImpl struct {
	conn    *grpc.ClientConn
	service proto.AuthServiceClient
}

func NewWrapper(cnf *Config) *ClientImpl {
	conn, err := grpc.Dial(cnf.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().Msgf("Failed to connect to AuthService: %v", err)
	}

	return &ClientImpl{
		conn:    conn,
		service: proto.NewAuthServiceClient(conn),
	}
}

func (c *ClientImpl) ValidateToken(ctx context.Context, token string) (*ValidateTokenResponse, error) {
	resp, err := c.service.ValidateToken(ctx, &proto.ValidateTokenRequest{Token: token})
	if err != nil {
		return nil, err
	}
	validateTokenResponse := &ValidateTokenResponse{Valid: resp.Valid}
	if resp.Valid {
		validateTokenResponse.User = domain.NewUser(resp.User.GetId(), resp.User.GetUsername())
	}
	return validateTokenResponse, nil
}

func (c *ClientImpl) Close() {
	if c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			log.Info().Msgf("Failed to close connection: %v", err)
			return
		}
	}
}
