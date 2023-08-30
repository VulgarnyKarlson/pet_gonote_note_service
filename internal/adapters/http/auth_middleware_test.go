package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/auth"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
)

func TestAuthMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthClient := auth.NewMockClient(ctrl)

	tests := []struct {
		name                 string
		givenToken           string
		mockValidationResult *auth.ValidateTokenResponse
		mockValidationError  error
		expectedStatusCode   int
	}{
		{
			name:       "Valid token",
			givenToken: "Bearer validToken",
			mockValidationResult: &auth.ValidateTokenResponse{
				User:  &domain.User{ID: "1", UserName: "JohnDoe"},
				Valid: true,
			},
			mockValidationError: nil,
			expectedStatusCode:  http.StatusOK,
		},
		{
			name:                 "Invalid token",
			givenToken:           "Bearer invalidToken",
			mockValidationResult: &auth.ValidateTokenResponse{Valid: false},
			mockValidationError:  nil,
			expectedStatusCode:   http.StatusUnauthorized,
		},
		{
			name:                 "No token",
			givenToken:           "",
			mockValidationResult: nil,
			mockValidationError:  nil,
			expectedStatusCode:   http.StatusUnauthorized,
		},
		{
			name:                 "Token with wrong format",
			givenToken:           "invalidFormat",
			mockValidationResult: nil,
			mockValidationError:  nil,
			expectedStatusCode:   http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &Server{auth: mockAuthClient}

			req, err := http.NewRequestWithContext(context.TODO(), "GET", "/somepath", http.NoBody)
			assert.NoError(t, err)

			if tt.givenToken != "" {
				req.Header.Set("Authorization", tt.givenToken)
			}

			rr := httptest.NewRecorder()

			mockAuthClient.EXPECT().
				ValidateToken(gomock.Any(), gomock.Any()).
				Return(tt.mockValidationResult, tt.mockValidationError).MaxTimes(1)

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
			server.AuthMiddleware(handler).ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}
