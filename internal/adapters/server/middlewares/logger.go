package middlewares

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/server"
)

const middlewareLoggerID = "logger"

type Logger struct {
	logger *zerolog.Logger
}

func RegisterLoggerMiddleware(logger *zerolog.Logger, s server.Server) {
	w := &Logger{logger: logger}
	s.Use(w.Middleware, middlewareLoggerID)
}

func LoggerID() string {
	return middlewareLoggerID
}

func (w *Logger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, reader *http.Request) {
		start := time.Now()
		next.ServeHTTP(writer, reader)
		w.logger.
			Info().
			Str("method", reader.Method).
			Str("url", reader.URL.RequestURI()).
			Str("user_agent", reader.UserAgent()).
			Dur("elapsed_ms", time.Since(start)).
			Msg("incoming request")
	})
}
