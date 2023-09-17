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
	OutWrite(note uint64)
	InRead() <-chan *domain.Note
	InProxyRead() <-chan *domain.Note
	OutRead() <-chan uint64
	ErrChan() <-chan error
	InClose()
	InProxyClose()
	OutClose()
	ErrClose()
}

type Impl struct {
	inChan   chan *domain.Note
	inProxy  chan *domain.Note
	outChan  chan uint64
	errChan  chan error
	err      error
	ctx      context.Context
	isClosed bool
}

func NewStream(ctx context.Context) (*Impl, context.Context) {
	s := &Impl{
		inChan:  make(chan *domain.Note),
		inProxy: make(chan *domain.Note),
		outChan: make(chan uint64),
		errChan: make(chan error),
		ctx:     ctx,
	}

	return s, ctx
}

func (s *Impl) Close() {
	if s.isClosed {
		return
	}
	s.isClosed = true
}

func (s *Impl) Fail(err error) {
	if s.isClosed {
		return
	}
	s.isClosed = true
	s.err = err
	s.errChan <- err
}

func (s *Impl) Done() <-chan struct{} {
	return s.ctx.Done()
}

func (s *Impl) Destroy() {
	s.isClosed = true
	s.Drain()
}

func (s *Impl) Drain() {
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

func (s *Impl) Err() error {
	err := s.err
	if err == nil {
		err = context.Cause(s.ctx)
	}
	return err
}
