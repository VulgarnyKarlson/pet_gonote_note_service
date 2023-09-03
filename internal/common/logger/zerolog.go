package logger

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx/fxevent"
)

const constructorPathRemove = "gitlab.karlson.dev/individual/pet_gonote/note_service/internal/"

func SetupLogger(conf *Config) *zerolog.Logger {
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	lvl, err := zerolog.ParseLevel(conf.Level)
	if err != nil {
		lvl = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(lvl)
	return &log.Logger
}

type logger struct {
	Logger *zerolog.Logger
}

//nolint:gocyclo
func (l *logger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		shortFunction := strings.ReplaceAll(e.FunctionName, constructorPathRemove, "")
		shortCaller := strings.ReplaceAll(e.CallerName, constructorPathRemove, "")
		l.Logger.Info().Str("callee", shortFunction).
			Str("caller", shortCaller).
			Msg("OnStart hook executing")
	case *fxevent.OnStartExecuted:
		shortFunction := strings.ReplaceAll(e.FunctionName, constructorPathRemove, "")
		shortCaller := strings.ReplaceAll(e.CallerName, constructorPathRemove, "")
		if e.Err != nil {
			l.Logger.Warn().Err(e.Err).
				Str("callee", shortFunction).
				Str("caller", shortCaller).
				Msg("OnStart hook failed")
		} else {
			l.Logger.Info().Str("callee", shortCaller).
				Str("caller", shortCaller).
				Str("runtime", e.Runtime.String()).
				Msg("OnStart hook executed")
		}
	case *fxevent.OnStopExecuting:
		shortFunction := strings.ReplaceAll(e.FunctionName, constructorPathRemove, "")
		shortCaller := strings.ReplaceAll(e.CallerName, constructorPathRemove, "")
		l.Logger.Info().Str("callee", shortFunction).
			Str("caller", shortCaller).
			Msg("OnStop hook executing")
	case *fxevent.OnStopExecuted:
		shortFunction := strings.ReplaceAll(e.FunctionName, constructorPathRemove, "")
		shortCaller := strings.ReplaceAll(e.CallerName, constructorPathRemove, "")
		if e.Err != nil {
			l.Logger.Warn().Err(e.Err).
				Str("callee", shortFunction).
				Str("callee", shortCaller).
				Msg("OnStop hook failed")
		} else {
			l.Logger.Info().Str("callee", shortFunction).
				Str("caller", shortCaller).
				Str("runtime", e.Runtime.String()).
				Msg("OnStop hook executed")
		}
	case *fxevent.Supplied:
		if e.Err != nil {
			l.Logger.Warn().Err(e.Err).Str("type", e.TypeName).Msg("supplied")
		} else {
			l.Logger.Info().Str("type", e.TypeName).Msg("supplied")
		}
	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			shortConstructor := strings.ReplaceAll(e.ConstructorName, constructorPathRemove, "")
			l.Logger.Info().Str("type", rtype).
				Str("constructor", shortConstructor).
				Msg("provided")
		}
		if e.Err != nil {
			l.Logger.Error().Err(e.Err).Msg("error encountered while applying options")
		}
	case *fxevent.Invoking:
		// Do nothing. Will log on Invoked.

	case *fxevent.Invoked:
		shortFunction := strings.ReplaceAll(e.FunctionName, constructorPathRemove, "")
		if e.Err != nil {
			l.Logger.Error().Err(e.Err).Str("stack", e.Trace).
				Str("function", shortFunction).Msg("invoke failed")
		} else {
			l.Logger.Info().Str("function", shortFunction).Msg("invoked")
		}
	case *fxevent.Stopping:
		l.Logger.Info().Str("signal", strings.ToUpper(e.Signal.String())).Msg("received signal")
	case *fxevent.Stopped:
		if e.Err != nil {
			l.Logger.Error().Err(e.Err).Msg("stop failed")
		}
	case *fxevent.RollingBack:
		l.Logger.Error().Err(e.StartErr).Msg("start failed, rolling back")
	case *fxevent.RolledBack:
		if e.Err != nil {
			l.Logger.Error().Err(e.Err).Msg("rollback failed")
		}
	case *fxevent.Started:
		if e.Err != nil {
			l.Logger.Error().Err(e.Err).Msg("start failed")
		} else {
			l.Logger.Info().Msg("started")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			l.Logger.Error().Err(e.Err).Msg("custom logger initialization failed")
		} else {
			l.Logger.Info().Str("function", e.ConstructorName).Msg("initialized custom fxevent.Logger")
		}
	}
}

// WithZerolog customize zerolog instance for fxevent.
func WithZerolog(l *zerolog.Logger) func() fxevent.Logger {
	SetupLogger(&Config{Level: "info"})
	return func() fxevent.Logger {
		return &logger{Logger: l}
	}
}
