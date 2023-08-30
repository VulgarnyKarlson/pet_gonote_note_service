package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/http/handlers"

	"github.com/gorilla/mux"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/auth"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/customerrors"

	"github.com/rs/zerolog/log"
)

type Server struct {
	cfg         *Config
	auth        auth.Client
	router      *mux.Router
	handlers    *handlers.NoteHandlers
	httpAdapter http.Server
}

func NewServer(cfg *Config, authClient auth.Client, noteHandlers *handlers.NoteHandlers) *Server {
	router := mux.NewRouter()
	s := &Server{
		cfg:      cfg,
		auth:     authClient,
		router:   router,
		handlers: noteHandlers,
		httpAdapter: http.Server{
			ReadTimeout: time.Duration(cfg.ReadTimeout) * time.Second,
			Addr:        cfg.Addr,
			Handler:     router,
		},
	}
	s.initRouter()
	return s
}

func (s *Server) Run() {
	log.Info().Msgf("Starting HTTP server on %s", s.httpAdapter.Addr)
	if err := s.httpAdapter.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Msgf("Failed to listen and serve: %v", err)
	}
}

func (s *Server) Stop() {
	log.Info().Msg("Stopping HTTP server")
	if err := s.httpAdapter.Shutdown(nil); err != nil {
		log.Fatal().Msgf("Failed to shutdown HTTP server: %v", err)
	}
}

func (s *Server) handlerErrors(h func(*http.Request) (*domain.HTTPResponse, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response, err := h(r)
		if err != nil {
			var customErr *customerrors.HTTPError
			ok := errors.As(err, &customErr)
			if !ok {
				customErr = customerrors.ErrInternalServer
			}
			w.WriteHeader(customErr.Code)
			response = &domain.HTTPResponse{Error: customErr.Message}
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(response.Status)
		}
		_ = json.NewEncoder(w).Encode(response)
	}
}
