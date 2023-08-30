package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/rs/zerolog"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/customerrors"

	"go.uber.org/mock/gomock"

	"github.com/stretchr/testify/assert"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
)

func TestCreateNote(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	user := &domain.User{
		ID:       "user_id",
		UserName: "user_name",
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var ctx context.Context
	var inputNoteChan chan *domain.Note
	var outputNoteIDsChan chan string
	var errChan chan error
	tests := []struct {
		name            string
		reqBody         []byte
		preReqFunc      func()
		expectedStatus  int
		expectedErr     string
		expectedNoteRes *domain.HTTPResponse
	}{
		{
			name:    "Successful case",
			reqBody: []byte(`[{"title": "Note 1", "content": "Content 1"}]`),
			preReqFunc: func() {
				inputNoteChan <- &domain.Note{
					ID:      "123",
					UserID:  user.ID,
					Title:   "Note 1",
					Content: "Content 1",
				}
				outputNoteIDsChan <- "123"
				close(outputNoteIDsChan)
			},
			expectedStatus: http.StatusOK,
			expectedErr:    "",
			expectedNoteRes: &domain.HTTPResponse{
				Data: &сreateNoteResponse{
					TotalNotes: 1,
					NoteIDs:    []string{"123"},
				},
				Status: http.StatusOK,
			},
		},
		{
			name:    "Request canceled",
			reqBody: []byte(`[{"title": "Note 1", "content": "Content 1"}]`),
			preReqFunc: func() {
				ctxx, cancel := context.WithCancel(ctx)
				ctx = ctxx
				cancel()
			},
			expectedStatus: http.StatusBadRequest,
			expectedErr:    "request-canceled",
		},
		{
			name:           "Nodata JSON",
			reqBody:        []byte(`[{"no": "data"}]`),
			expectedStatus: http.StatusBadRequest,
			expectedErr:    "invalid-json",
		},
		{
			name:           "Invalid StartJSON",
			reqBody:        []byte(`["title": "Note 1", "content": "Content 1"}, "title": "Note 1", "content": "Content 1"}]`),
			expectedStatus: http.StatusBadRequest,
			expectedErr:    "invalid-json",
		},
		{
			name:           "Invalid MiddleJSON",
			reqBody:        []byte(`[{"title": "Note 1", "content": "Content 1"}, "title": "Note 1", "content": "Content 1"}]`),
			expectedStatus: http.StatusBadRequest,
			expectedErr:    "invalid-json",
		},
		{
			name:           "Invalid EndJSON",
			reqBody:        []byte(`[{"title": "Note 1", "content": "Content 1"}, {"title": "Note 1", "content": "Content 1"}`),
			expectedStatus: http.StatusBadRequest,
			expectedErr:    "invalid-json",
		},
		// TODO - temporary disabled, need way to handle repository-errors
		//{
		//	name:    "Repository Error",
		//	reqBody: []byte(`[{"title": "Note 1", "content": "Content 1"}]`),
		//	preReqFunc: func() {
		//		errChan <- fmt.Errorf("repository error")
		//	},
		//	expectedStatus: http.StatusInternalServerError,
		//	expectedErr:    "repository-error",
		// },
	}

	for _, tt := range tests {
		inputNoteChan = make(chan *domain.Note, 5)
		outputNoteIDsChan = make(chan string, 5)
		errChan = make(chan error, 1)
		mockNoteService := domain.NewMockNoteService(ctrl)
		h := &NoteHandlers{
			noteServicePort: mockNoteService,
		}
		t.Run(tt.name, func(t *testing.T) {
			ctx = context.WithValue(context.Background(), domain.UserCtxKey, user)
			if tt.preReqFunc != nil {
				tt.preReqFunc()
			}
			r, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/create_note", bytes.NewBuffer(tt.reqBody))
			mockNoteService.EXPECT().
				Create(r.Context(), user).
				Return(inputNoteChan, outputNoteIDsChan, errChan)

			resp, err := h.CreateNote(r)
			if err != nil {
				var httpErr *customerrors.HTTPError
				errors.As(err, &httpErr)
				assert.Equal(t, tt.expectedErr, httpErr.Message)
				assert.Equal(t, tt.expectedStatus, httpErr.Code)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.expectedNoteRes, resp)
			}
		})
	}
}
