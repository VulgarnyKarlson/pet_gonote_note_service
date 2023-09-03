package stream

import "gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"

func (s *StreamImpl) InRead() <-chan *domain.Note {
	return s.inChan
}

func (s *StreamImpl) InProxyRead() <-chan *domain.Note {
	return s.inProxy
}

func (s *StreamImpl) OutRead() <-chan string {
	return s.outChan
}

func (s *StreamImpl) ErrChan() <-chan error {
	return s.errChan
}
