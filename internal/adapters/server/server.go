package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/customerrors"
)

type Wrapper struct {
	logger      *zerolog.Logger
	cfg         *Config
	router      *mux.Router
	middleWares map[string]mux.MiddlewareFunc
	httpAdapter http.Server
}

func NewServer(logger *zerolog.Logger, cfg *Config, router *mux.Router) Server {
	s := &Wrapper{
		logger:      logger,
		cfg:         cfg,
		router:      router,
		middleWares: make(map[string]mux.MiddlewareFunc),
		httpAdapter: http.Server{
			ReadTimeout: time.Duration(cfg.ReadTimeout) * time.Second,
			Addr:        cfg.Addr,
			Handler:     router,
		},
	}
	s.initRouter()
	return s
}

func (s *Wrapper) Run() error {
	s.logger.Info().Msgf("Starting HTTP server on %s", s.httpAdapter.Addr)
	if err := s.httpAdapter.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Fatal().Msgf("Failed to listen and serve: %v", err)
	}
	return nil
}

func (s *Wrapper) Stop() error {
	s.logger.Info().Msg("Stopping HTTP server")
	if err := s.httpAdapter.Shutdown(nil); err != nil {
		return fmt.Errorf("failed to shutdown HTTP server: %v", err)
	}
	return nil
}

func (s *Wrapper) handlerErrors(h func(*http.Request) (*Response, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	})
}
