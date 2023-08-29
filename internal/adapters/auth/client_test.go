package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/proto"
)

func TestValidateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var tests = []struct {
		name           string
		token          string
		mockResp       *proto.ValidateTokenResponse
		mockErr        error
		expectedResult *ValidateTokenResponse
		expectedErr    error
	}{
		{
			name:  "Valid token",
			token: "valid_token",
			mockResp: &proto.ValidateTokenResponse{
				Valid: true,
				User: &proto.User{
					Id:       "1",
					Username: "JohnDoe",
				},
			},
			mockErr: nil,
			expectedResult: &ValidateTokenResponse{
				User: &domain.User{
					ID:       "1",
					UserName: "JohnDoe",
				},
				Valid: true,
			},
			expectedErr: nil,
		},
		{
			name:           "Invalid token",
			token:          "invalid_token",
			mockResp:       &proto.ValidateTokenResponse{Valid: false},
			mockErr:        nil,
			expectedResult: &ValidateTokenResponse{Valid: false},
			expectedErr:    nil,
		},
	}
	mockAuthService := proto.NewMockAuthServiceClient(ctrl)

	config := &Config{Address: "localhost:5000"}
	wrapper := NewWrapper(config)
	wrapper.service = mockAuthService

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthService.EXPECT().ValidateToken(
				gomock.Any(),
				&proto.ValidateTokenRequest{Token: tt.token},
			).Return(tt.mockResp, tt.mockErr)

			result, err := wrapper.ValidateToken(context.Background(), tt.token)

			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
