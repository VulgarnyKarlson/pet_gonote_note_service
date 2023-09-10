package stream

import "gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"

func (s *Impl) InRead() <-chan *domain.Note {
	return s.inChan
}

func (s *Impl) InProxyRead() <-chan *domain.Note {
	return s.inProxy
}

func (s *Impl) OutRead() <-chan string {
	return s.outChan
}

func (s *Impl) ErrChan() <-chan error {
	return s.errChan
}
