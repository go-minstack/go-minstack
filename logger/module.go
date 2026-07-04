package logger

import (
	"log/slog"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

// Module configures zerolog from env vars, provides *slog.Logger, and registers
// the zerolog-backed FX event logger.
//
// Env vars:
//
//	MINSTACK_LOG_LEVEL  — trace | debug | info | warn | error  (default: info)
//	MINSTACK_LOG_FORMAT — json | console                        (default: json)
func Module() fx.Option {
	initLogger()
	return fx.Options(
		fx.WithLogger(func() fxevent.Logger { return &fxZeroLogger{log.Logger} }),
		fx.Module("logger",
			fx.Provide(newSlogLogger),
		),
	)
}

func newSlogLogger() *slog.Logger {
	return slog.New(newZerologHandler(log.Logger))
}

func initLogger() {
	switch os.Getenv("MINSTACK_LOG_FORMAT") {
	case "console":
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	default:
		log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	}

	switch os.Getenv("MINSTACK_LOG_LEVEL") {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
