package handlers

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
	"go.uber.org/mock/gomock"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/customerrors"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/stream"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain/tests"
)

func TestReadNotes(t *testing.T) {
	tests.TestIsUnit(t)
	defer goleak.VerifyNone(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name        string
		inputjson   string
		writeCount  int
		expectedErr error
	}{
		{
			name: "valid json",
			inputjson: `[{"id": 1, "user_id": "user1", "title": "Note 1", "content": "Content 1"},
				  {"id": 2, "user_id": 1, "title": "Note 2", "content": "Content 2"}]`,
			writeCount:  2,
			expectedErr: nil,
		},
		{
			name:        "invalid json - no open delimiter",
			inputjson:   `{"id": 1, "user_id": 1, "title": "Note 1", "content": "Content 1"}]`,
			expectedErr: customerrors.ErrInvalidJSONOpenDelimiter,
		},
		{
			name:        "invalid json - no close delimiter",
			inputjson:   `[{"id": 1, "user_id": 1, "title": "Note 1", "content": "Content 1"}`,
			writeCount:  1,
			expectedErr: customerrors.ErrInvalidJSONCloseDelimiter,
		},
		{
			name:        "invalid json",
			inputjson:   `[{"id": 1, "user_id": 1, title": "Note 1", "content": "Content 1"}]`,
			expectedErr: customerrors.ErrInvalidJSON,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockStream := stream.NewMockStream(ctrl)
			reader := bytes.NewBufferString(tc.inputjson)
			mockStream.EXPECT().Done().AnyTimes()
			mockStream.EXPECT().InWrite(gomock.Any()).Times(tc.writeCount)
			if tc.writeCount != 0 && tc.expectedErr == nil {
				mockStream.EXPECT().InClose().Times(1)
			}

			syncCh := make(chan struct{})
			go func() {
				err := readNotes(reader, mockStream, domain.NewUser(1, "user1"))
				assert.Equal(t, tc.expectedErr, err)
				close(syncCh)
			}()

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			select {
			case <-syncCh:
				return
			case <-ctx.Done():
				t.Fatal("timeout")
			}
		})
	}
}
