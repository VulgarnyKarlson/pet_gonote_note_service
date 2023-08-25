package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/customerrors"

	"github.com/rs/zerolog/log"
)

type Server struct {
	cfg *Config
	http.Server
}

func NewServer(cfg *Config, handler http.Handler) *Server {
	return &Server{
		cfg: cfg,
		Server: http.Server{
			ReadTimeout: time.Duration(cfg.ReadTimeout) * time.Second,
			Addr:        cfg.Addr,
			Handler:     handler,
		},
	}
}

func (w *Server) Run() {
	log.Info().Msgf("Starting HTTP server on %s", w.Addr)
	if err := w.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Msgf("Failed to listen and serve: %v", err)
	}
}

func (w *Server) Stop() {
	log.Info().Msg("Stopping HTTP server")
	if err := w.Shutdown(nil); err != nil {
		log.Fatal().Msgf("Failed to shutdown HTTP server: %v", err)
	}
}

type Response struct {
	Data   any    `json:"data,omitempty"`
	Status int    `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

func (w *Server) HandlerErrors(h func(*http.Request) (*Response, error)) http.HandlerFunc {
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
