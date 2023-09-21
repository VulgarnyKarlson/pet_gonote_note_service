package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Wrapper) AddEndpoint(e Endpoint) {
	handler := s.handlerErrors(e.Handler)
	for _, id := range e.Middlewares {
		m, ok := s.middleWares[id]
		if !ok {
			s.logger.Fatal().Msgf("Middleware %s not found", id)
		}
		handler = m.Middleware(handler)
	}
	s.router.Handle(e.Path, handler).Methods(e.Method)
	s.logger.Info().Msgf("Added endpoint: %s %s", e.Method, e.Path)
}

func (s *Wrapper) initRouter() {
	s.router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}).Methods("GET")
}

func (s *Wrapper) Use(m mux.MiddlewareFunc, id string) {
	s.middleWares[id] = m
}
