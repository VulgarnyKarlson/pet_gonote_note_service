package http

import (
	"net/http"
	"time"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/logger"
)

func (s *Server) requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		l := logger.Get()

		next.ServeHTTP(w, r)

		l.
			Info().
			Str("method", r.Method).
			Str("url", r.URL.RequestURI()).
			Str("user_agent", r.UserAgent()).
			Dur("elapsed_ms", time.Since(start)).
			Msg("incoming request")
	})
}
