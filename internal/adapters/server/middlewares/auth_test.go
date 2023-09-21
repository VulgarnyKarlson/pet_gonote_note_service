package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
	"go.uber.org/mock/gomock"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/auth"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain/tests"
)

func TestAuthMiddleware(t *testing.T) {
	tests.TestIsUnit(t)
	defer goleak.VerifyNone(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthClient := auth.NewMockClient(ctrl)

	testCases := []struct {
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
				User:  domain.NewUser(1, "JohnDoe"),
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

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			middleWare := AuthMiddleware{auth: mockAuthClient}

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
			middleWare.Middleware(handler).ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}
