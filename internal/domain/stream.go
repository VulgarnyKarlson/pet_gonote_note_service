package domain

import (
	"context"
)

type Stream interface {
	Close()
	Fail(err error)
	Done() <-chan struct{}
	Destroy()
	Drain()
	Err() error
	InWrite(note *Note)
	InProxyWrite(note *Note)
	OutWrite(note string)
	InRead() <-chan *Note
	InProxyRead() <-chan *Note
	OutRead() <-chan string
	ErrChan() <-chan error
	InClose()
	InProxyClose()
	OutClose()
	ErrClose()
}

type StreamImpl struct {
	inChan        chan *Note
	inProxy       chan *Note
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
		inChan:        make(chan *Note),
		inProxy:       make(chan *Note),
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

func (s *StreamImpl) InWrite(note *Note) {
	if s.isClosed {
		return
	}
	s.inChan <- note
}

func (s *StreamImpl) InProxyWrite(note *Note) {
	if s.isClosed {
		return
	}
	s.inProxy <- note
}

func (s *StreamImpl) OutWrite(note string) {
	if s.isClosed {
		return
	}
	s.outChan <- note
}

func (s *StreamImpl) InRead() <-chan *Note {
	return s.inChan
}

func (s *StreamImpl) InProxyRead() <-chan *Note {
	return s.inProxy
}

func (s *StreamImpl) OutRead() <-chan string {
	return s.outChan
}

func (s *StreamImpl) ErrChan() <-chan error {
	return s.errChan
}

func (s *StreamImpl) InClose() {
	if s.isClosed {
		return
	}
	close(s.inChan)
}

func (s *StreamImpl) InProxyClose() {
	if s.isClosed {
		return
	}
	close(s.inProxy)
}

func (s *StreamImpl) OutClose() {
	if s.isClosed {
		return
	}
	close(s.outChan)
}

func (s *StreamImpl) ErrClose() {
	if s.isClosed {
		return
	}
	close(s.errChan)
}
