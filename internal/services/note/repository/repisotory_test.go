package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
	"go.uber.org/mock/gomock"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/postgres"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/customerrors"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/stream"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain/tests"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/services/noteoutbox"
)

var testInteggration *tests.TestIntegration

func TestMain(m *testing.M) {
	testInteggration = tests.NewTestIntegration(m)
	testInteggration.Setup()
	testInteggration.RunServices(tests.TestIntegrationServicePostgres)
	goleak.VerifyTestMain(m, goleak.Cleanup(func(exitCode int) {
		testInteggration.Teardown()
		os.Exit(exitCode)
	}))
}

func TestMockRepository_All(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	tests.TestIsIntegration(t)
	pgxPool, err := postgres.New(testInteggration.Configs.GetPGConfig())
	if err != nil {
		log.Panic().Err(err).Msg("error creating postgres pool")
	}
	defer pgxPool.Close()
	var st *stream.MockStream
	var testRepository Repository
	var outboxRepo *noteoutbox.MockRepository
	content := "TEST"
	note1, _ := domain.NewNote(1, 1, "Note 1", content)
	note2, _ := domain.NewNote(2, 1, "Note 2", content)
	note3, _ := domain.NewNote(3, 1, "Note 3", content)
	note4, _ := domain.NewNote(4, 1, "Note 4", content)
	user := domain.NewUser(1, "user1")
	testCases := []struct {
		name string
		run  func()
	}{
		{
			name: "context cancel",
			run: func() {
				ctx, cancel := context.WithCancel(context.Background())
				proxyNotes := make(chan *domain.Note)
				notes := []*domain.Note{note1, note2, note3}
				syncCh := make(chan struct{})
				st.EXPECT().Done().Return(ctx.Done()).AnyTimes()
				st.EXPECT().InProxyRead().Return(proxyNotes).AnyTimes()
				outboxRepo.EXPECT().Create(ctx, gomock.Any(), gomock.Any()).Times(2)
				st.EXPECT().OutWrite(gomock.Any()).Times(2)
				st.EXPECT().Fail(gomock.Any()).Times(1)
				st.EXPECT().OutClose().Times(1)
				st.EXPECT().Close().Times(1)
				go func() {
					// let db transaction start
					for _, note := range notes {
						proxyNotes <- note
					}
					cancel()
					syncCh <- struct{}{}
				}()
				testRepository.CreateNote(ctx, st)
				<-syncCh
			},
		},
		{
			name: "success with 1 batch",
			run: func() {
				ctx := context.Background()
				proxyNotes := make(chan *domain.Note)
				notes := []*domain.Note{note1, note2, note3}
				syncCh := make(chan struct{})
				st.EXPECT().Done().Return(ctx.Done()).AnyTimes()
				st.EXPECT().InProxyRead().Return(proxyNotes).AnyTimes()
				outboxRepo.EXPECT().Create(ctx, gomock.Any(), gomock.Any()).Times(len(notes))
				st.EXPECT().OutWrite(gomock.Any()).Times(len(notes))
				st.EXPECT().OutClose().Times(1)
				st.EXPECT().Close().Times(1)
				st.EXPECT().Err().Return(nil).Times(1)
				go func() {
					// let db transaction start
					for _, note := range notes {
						proxyNotes <- note
					}
					close(proxyNotes)
					syncCh <- struct{}{}
				}()
				testRepository.CreateNote(ctx, st)
				<-syncCh

				outboxRepo.EXPECT().Search(ctx, gomock.Any(), gomock.Any()).Times(1)
				searchNotes, _ := testRepository.SearchNote(
					ctx,
					user,
					&domain.SearchCriteria{Content: content},
				)
				assert.Len(t, searchNotes, len(notes))

				searchNotes[0].SetTitle("Changed title")
				searchNotes[0].SetUpdatedAt(time.Now())
				outboxRepo.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).Times(1)
				err = testRepository.UpdateNote(ctx, user, searchNotes[0])
				assert.Nil(t, err)
				err = testRepository.UpdateNote(ctx, user, note4)
				assert.ErrorIs(t, err, customerrors.ErrNotFoundNoteID)

				outboxRepo.EXPECT().FindByID(ctx, gomock.Any(), gomock.Any()).Times(1)
				note, _ := testRepository.ReadNoteByID(ctx, user, searchNotes[0].ID())
				noteNil, _ := testRepository.ReadNoteByID(ctx, user, 0)
				assert.Nil(t, noteNil)
				assert.Equal(t, searchNotes[0].ID(), note.ID())
				assert.Equal(t, searchNotes[0].Content(), note.Content())
				assert.Equal(t, searchNotes[0].Title(), note.Title())
				assert.Equal(t, searchNotes[0].CreatedAt(), note.CreatedAt())
				assert.Equal(t, searchNotes[0].UserID(), note.UserID())

				outboxRepo.EXPECT().Delete(ctx, gomock.Any(), gomock.Any()).Times(len(searchNotes))
				for _, note = range searchNotes {
					_, _ = testRepository.DeleteNote(
						ctx,
						user,
						note.ID(),
					)
				}

				outboxRepo.EXPECT().Search(ctx, gomock.Any(), gomock.Any()).Times(1)
				searchNotes, _ = testRepository.SearchNote(
					ctx,
					user,
					&domain.SearchCriteria{},
				)
				assert.Len(t, searchNotes, 0)
			},
		},
		{
			name: "success without batch",
			run: func() {
				ctx := context.Background()
				proxyNotes := make(chan *domain.Note)
				notes := []*domain.Note{note1}
				syncCh := make(chan struct{})
				st.EXPECT().Done().Return(ctx.Done()).AnyTimes()
				st.EXPECT().InProxyRead().Return(proxyNotes).AnyTimes()
				outboxRepo.EXPECT().Create(ctx, gomock.Any(), gomock.Any()).Times(len(notes))
				st.EXPECT().OutWrite(gomock.Any()).Times(len(notes))
				st.EXPECT().OutClose().Times(1)
				st.EXPECT().Close().Times(1)
				st.EXPECT().Err().Return(nil).Times(1)
				go func() {
					// let db transaction start
					for _, note := range notes {
						proxyNotes <- note
					}
					close(proxyNotes)
					syncCh <- struct{}{}
				}()
				testRepository.CreateNote(ctx, st)
				<-syncCh

				outboxRepo.EXPECT().Search(ctx, gomock.Any(), gomock.Any()).Times(1)
				searchNotes, _ := testRepository.SearchNote(
					ctx,
					user,
					&domain.SearchCriteria{Content: "TEST"},
				)
				assert.Len(t, searchNotes, len(notes))

				// clean up
				outboxRepo.EXPECT().Delete(ctx, gomock.Any(), gomock.Any()).Times(len(searchNotes))
				for _, note := range searchNotes {
					_, _ = testRepository.DeleteNote(
						ctx,
						user,
						note.ID(),
					)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			st = stream.NewMockStream(ctrl)
			outboxRepo = noteoutbox.NewMockRepository(ctrl)
			testRepository = NewRepository(&log.Logger, &Config{CreateNotesBatchSize: 2}, pgxPool, outboxRepo)
			tc.run()
		})
	}
}
