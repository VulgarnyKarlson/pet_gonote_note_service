package stream

import "gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"

func (s *Impl) SetUser(user *domain.User) {
	s.user = user
}

func (s *Impl) User() *domain.User {
	return s.user
}
