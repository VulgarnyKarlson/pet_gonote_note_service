package stream

import (
	"context"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
)

type Stream interface {
	Close()
	Fail(err error)
	Done() <-chan struct{}
	Destroy()
	Drain()
	Err() error
	InWrite(note *domain.Note)
	InProxyWrite(note *domain.Note)
	OutWrite(note string)
	InRead() <-chan *domain.Note
	InProxyRead() <-chan *domain.Note
	OutRead() <-chan string
	ErrChan() <-chan error
	InClose()
	InProxyClose()
	OutClose()
	ErrClose()
}

type StreamImpl struct {
	inChan        chan *domain.Note
	inProxy       chan *domain.Note
	outChan       chan string
	errChan       chan error
	err           error
	ctx           context.Context
	ctxCancelFunc context.CancelFunc
	isClosed      bool
}

func NewStream(originalCtx context.Context) (*StreamImpl, context.Context) {
	ctx, cancel := context.WithCancel(originalCtx)
	s := &StreamImpl{
		inChan:        make(chan *domain.Note),
		inProxy:       make(chan *domain.Note),
		outChan:       make(chan string),
		errChan:       make(chan error),
		ctx:           ctx,
		ctxCancelFunc: cancel,
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-s.Done():
				return
			}
		}
	}()
	return s, ctx
}

func (s *StreamImpl) Close() {
	if s.isClosed {
		return
	}
	s.isClosed = true
	s.ctxCancel()
}

func (s *StreamImpl) ctxCancel() {
	if s.ctxCancelFunc != nil {
		s.ctxCancelFunc()
	}
}

func (s *StreamImpl) Fail(err error) {
	if s.isClosed {
		return
	}
	s.isClosed = true
	s.err = err
	s.errChan <- err
}

func (s *StreamImpl) Done() <-chan struct{} {
	return s.ctx.Done()
}

func (s *StreamImpl) Destroy() {
	s.isClosed = true
	s.Drain()
	s.ctxCancel()
}

func (s *StreamImpl) Drain() {
	for {
		select {
		case <-s.InRead():
			s.InClose()
			s.inChan = nil
		case <-s.InProxyRead():
			s.InProxyClose()
			s.inProxy = nil
		case <-s.OutRead():
			s.OutClose()
			s.outChan = nil
		case <-s.ErrChan():
			s.ErrClose()
			s.errChan = nil
		default:
			return
		}
	}
}

func (s *StreamImpl) Err() error {
	err := s.err
	if err == nil {
		err = s.ctx.Err()
	}
	return err
}
