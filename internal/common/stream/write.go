package stream

import "gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"

func (s *Impl) InWrite(note *domain.Note) {
	if s.isClosed {
		return
	}
	s.inChan <- note
}

func (s *Impl) InProxyWrite(note *domain.Note) {
	if s.isClosed {
		return
	}
	s.inProxy <- note
}

func (s *Impl) OutWrite(note string) {
	if s.isClosed {
		return
	}
	s.outChan <- note
}
