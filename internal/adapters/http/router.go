package http

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

type Endpoint struct {
	Method  string
	Path    string
	Auth    bool
	Handler func(r *http.Request) (*Response, error)
}

func (s *Server) AddEndpoint(e Endpoint) {
	handler := s.handlerErrors(e.Handler)
	if e.Auth {
		handler = s.AuthMiddleware(handler)
	}
	s.router.HandleFunc(e.Path, handler).Methods(e.Method)
	log.Info().Msgf("Added endpoint: %s %s", e.Method, e.Path)
}

func (s *Server) initRouter() {
	s.router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}).Methods("GET")

	s.router.Use(s.requestLogger)
}
