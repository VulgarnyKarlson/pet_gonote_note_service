package http

import (
	"net/http"

	"github.com/rs/zerolog/log"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"
)

type Endpoint struct {
	method  string
	path    string
	auth    bool
	handler func(r *http.Request) (*domain.HTTPResponse, error)
}

func (s *Server) AddEndpoint(e Endpoint) {
	handler := s.handlerErrors(e.handler)
	if e.auth {
		handler = s.AuthMiddleware(handler)
	}
	s.router.HandleFunc(e.path, handler).Methods(e.method)
	log.Info().Msgf("Added endpoint: %s %s", e.method, e.path)
}

func (s *Server) initRouter() {
	endpoints := []Endpoint{
		{method: "POST", path: "/create", auth: true, handler: s.handlers.CreateNote},
		{method: "GET", path: "/read", auth: true, handler: s.handlers.ReadNoteByID},
		{method: "POST", path: "/update", auth: true, handler: s.handlers.UpdateNote},
		{method: "POST", path: "/delete", auth: true, handler: s.handlers.DeleteNoteByID},
		{method: "GET", path: "/search", auth: true, handler: s.handlers.SearchNote},
	}

	for _, e := range endpoints {
		s.AddEndpoint(e)
	}

	s.router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}).Methods("GET")

	s.router.Use(s.requestLogger)
}
