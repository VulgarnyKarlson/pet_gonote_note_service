package note

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
	"go.uber.org/mock/gomock"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/customerrors"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/stream"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain/tests"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/services/note/repository"
)

func TestNoteService_Create(t *testing.T) {
	tests.TestIsUnit(t)
	defer goleak.VerifyNone(t)
	var mockRepo *repository.MockRepository
	var st *stream.MockStream
	cfg := &Config{
		MaxTitleLength:   7,
		MaxContentLength: 11,
	}

	successNote1, err := domain.NewNote(1, 1, "Note 1", "Content 1")
	if err != nil {
		t.Fatal(err)
	}
	successNote2, err := domain.NewNote(2, 2, "Note 2", "Content 2")
	if err != nil {
		t.Fatal(err)
	}
	titleLongNote, err := domain.NewNote(3, 3, "Note 1111", "Content 3")
	if err != nil {
		t.Fatal(err)
	}
	contentLongNote, err := domain.NewNote(4, 4, "Note 4", "Content 2222222222 2")
	if err != nil {
		t.Fatal(err)
	}
	testCases := []struct {
		name            string
		notes           []*domain.Note
		expectedErr     error
		proxyWriteCount int
	}{
		{
			name: "valid",
			notes: []*domain.Note{
				successNote1,
				successNote2,
			},
			expectedErr:     nil,
			proxyWriteCount: 2,
		},
		{
			name: "title too long",
			notes: []*domain.Note{
				successNote2,
				titleLongNote,
			},
			expectedErr:     customerrors.ErrTitleTooLong,
			proxyWriteCount: 1,
		},
		{
			name: "content too long",
			notes: []*domain.Note{
				successNote1,
				contentLongNote,
			},
			expectedErr:     customerrors.ErrContentTooLong,
			proxyWriteCount: 1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			st = stream.NewMockStream(ctrl)
			mockRepo = repository.NewMockRepository(ctrl)
			mockRepo.EXPECT().CreateNote(
				gomock.Any(),
				st,
			).Times(1)
			st.EXPECT().Done().AnyTimes()
			if tc.expectedErr != nil {
				st.EXPECT().Fail(tc.expectedErr).Times(1)
			} else {
				st.EXPECT().InProxyClose().Times(1)
			}

			proxyWriteChan, syncCh := make(chan *domain.Note), make(chan struct{})
			st.EXPECT().InProxyRead().Return(proxyWriteChan).AnyTimes()
			st.EXPECT().InRead().Return(proxyWriteChan).AnyTimes()
			st.EXPECT().InProxyWrite(gomock.Any()).Times(tc.proxyWriteCount)
			go func() {
				for _, note := range tc.notes {
					proxyWriteChan <- note
				}
				close(proxyWriteChan)
				// let's just wait for all calls to be done
				time.Sleep(10 * time.Millisecond)
				close(syncCh)
			}()

			s := NewService(cfg, mockRepo)
			go s.Create(context.TODO(), st)

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			select {
			case <-ctx.Done():
				t.Fatal("timeout")
			case <-syncCh:
				return
			}
		})
	}
}

func TestNoteService_ReadByID(t *testing.T) {
	tests.TestIsUnit(t)
	defer goleak.VerifyNone(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var mockRepo *repository.MockRepository
	cfg := &Config{
		MaxTitleLength:   7,
		MaxContentLength: 11,
	}
	testCases := []struct {
		name        string
		id          uint64
		expectedErr error
	}{
		{
			name:        "valid",
			id:          1,
			expectedErr: nil,
		},
		{
			name:        "invalid id",
			id:          0,
			expectedErr: customerrors.ErrInvalidNoteID,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo = repository.NewMockRepository(ctrl)
			if tc.expectedErr == nil {
				mockRepo.EXPECT().ReadNoteByID(
					gomock.Any(),
					gomock.Any(),
					tc.id,
				).Times(1).Return(&domain.Note{}, nil)
			}

			s := NewService(cfg, mockRepo)
			_, err := s.ReadByID(context.Background(), domain.NewUser(1, ""), strconv.Itoa(int(tc.id)))
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestNoteService_Update(t *testing.T) {
	tests.TestIsUnit(t)
	defer goleak.VerifyNone(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var mockRepo *repository.MockRepository
	cfg := &Config{
		MaxTitleLength:   7,
		MaxContentLength: 11,
	}
	successNote, err := domain.NewNote(1, 1, "Note 1", "Content 1")
	if err != nil {
		t.Fatal(err)
	}
	invalidNote, err := domain.NewNote(0, 1, "Note 1", "Content 1")
	if err != nil {
		t.Fatal(err)
	}
	testCases := []struct {
		name        string
		note        *domain.Note
		expectedErr error
	}{
		{
			name:        "valid",
			note:        successNote,
			expectedErr: nil,
		},
		{
			name:        "invalid id",
			note:        invalidNote,
			expectedErr: customerrors.ErrInvalidNoteID,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo = repository.NewMockRepository(ctrl)
			if tc.expectedErr == nil {
				mockRepo.EXPECT().UpdateNote(
					gomock.Any(),
					gomock.Any(),
					tc.note,
				).Times(1).Return(nil)
			}

			s := NewService(cfg, mockRepo)
			err := s.Update(context.Background(), domain.NewUser(1, ""), tc.note)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestNoteService_Delete(t *testing.T) {
	tests.TestIsUnit(t)
	defer goleak.VerifyNone(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var mockRepo *repository.MockRepository
	cfg := &Config{
		MaxTitleLength:   7,
		MaxContentLength: 11,
	}
	testCases := []struct {
		name        string
		id          uint64
		expectedErr error
	}{
		{
			name:        "valid",
			id:          1,
			expectedErr: nil,
		},
		{
			name:        "invalid id",
			id:          0,
			expectedErr: customerrors.ErrInvalidNoteID,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo = repository.NewMockRepository(ctrl)
			if tc.expectedErr == nil {
				mockRepo.EXPECT().DeleteNote(
					gomock.Any(),
					gomock.Any(),
					tc.id,
				).Times(1).Return(true, nil)
			}

			s := NewService(cfg, mockRepo)
			_, err := s.Delete(context.Background(), domain.NewUser(1, ""), strconv.Itoa(int(tc.id)))
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestNoteService_Search(t *testing.T) {
	tests.TestIsUnit(t)
	defer goleak.VerifyNone(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var mockRepo *repository.MockRepository
	cfg := &Config{
		MaxTitleLength:   7,
		MaxContentLength: 11,
	}
	testCases := []struct {
		name        string
		criteria    *domain.SearchCriteria
		expectedErr error
	}{
		{
			name:        "valid",
			criteria:    &domain.SearchCriteria{},
			expectedErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo = repository.NewMockRepository(ctrl)
			if tc.expectedErr == nil {
				mockRepo.EXPECT().SearchNote(
					gomock.Any(),
					gomock.Any(),
					tc.criteria,
				).Times(1).Return([]*domain.Note{}, nil)
			}

			s := NewService(cfg, mockRepo)
			_, err := s.Search(context.Background(), domain.NewUser(1, ""), tc.criteria)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
