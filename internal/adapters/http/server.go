package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"

	"github.com/gorilla/mux"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/auth"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/customerrors"
)

type Server struct {
	logger      *zerolog.Logger
	cfg         *Config
	auth        auth.Client
	router      *mux.Router
	httpAdapter http.Server
}

func NewServer(logger *zerolog.Logger, cfg *Config, authClient auth.Client) *Server {
	router := mux.NewRouter()
	s := &Server{
		logger: logger,
		cfg:    cfg,
		auth:   authClient,
		router: router,
		httpAdapter: http.Server{
			ReadTimeout: time.Duration(cfg.ReadTimeout) * time.Second,
			Addr:        cfg.Addr,
			Handler:     router,
		},
	}
	s.initRouter()
	return s
}

func (s *Server) Run() error {
	s.logger.Info().Msgf("Starting HTTP server on %s", s.httpAdapter.Addr)
	if err := s.httpAdapter.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Fatal().Msgf("Failed to listen and serve: %v", err)
	}
	return nil
}

func (s *Server) Stop() error {
	s.logger.Info().Msg("Stopping HTTP server")
	if err := s.httpAdapter.Shutdown(nil); err != nil {
		return fmt.Errorf("failed to shutdown HTTP server: %v", err)
	}
	return nil
}

func (s *Server) handlerErrors(h func(*http.Request) (*Response, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response, err := h(r)
		if err != nil {
			var customErr *customerrors.HTTPError
			ok := errors.As(err, &customErr)
			if !ok {
				customErr = customerrors.ErrInternalServer
			}
			w.WriteHeader(customErr.Code)
			response = &Response{Error: customErr.Message}
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(response.Status)
		}
		_ = json.NewEncoder(w).Encode(response)
	}
}
