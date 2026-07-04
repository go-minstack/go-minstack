package logger

import (
	"strings"

	"github.com/rs/zerolog"
	"go.uber.org/fx/fxevent"
)

// fxZeroLogger is an internal fxevent.Logger that routes FX lifecycle events
// to zerolog. It is not exported — use Module() to register it.
type fxZeroLogger struct {
	logger zerolog.Logger
}

func (l *fxZeroLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.logger.Debug().
			Str("callee", e.FunctionName).
			Str("caller", e.CallerName).
			Msg("OnStart hook executing")
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			l.logger.Error().
				Str("callee", e.FunctionName).
				Str("caller", e.CallerName).
				Err(e.Err).
				Msg("OnStart hook failed")
		} else {
			l.logger.Debug().
				Str("callee", e.FunctionName).
				Str("caller", e.CallerName).
				Str("runtime", e.Runtime.String()).
				Msg("OnStart hook executed")
		}
	case *fxevent.OnStopExecuting:
		l.logger.Debug().
			Str("callee", e.FunctionName).
			Str("caller", e.CallerName).
			Msg("OnStop hook executing")
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			l.logger.Error().
				Str("callee", e.FunctionName).
				Str("caller", e.CallerName).
				Err(e.Err).
				Msg("OnStop hook failed")
		} else {
			l.logger.Debug().
				Str("callee", e.FunctionName).
				Str("caller", e.CallerName).
				Str("runtime", e.Runtime.String()).
				Msg("OnStop hook executed")
		}
	case *fxevent.Supplied:
		ev := l.logger.Debug().Str("type", e.TypeName)
		if e.ModuleName != "" {
			ev = ev.Str("module", e.ModuleName)
		}
		if e.Err != nil {
			ev.Err(e.Err).Msg("error encountered while applying options")
		} else {
			ev.Msg("supplied")
		}
	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			ev := l.logger.Debug().
				Str("constructor", e.ConstructorName).
				Str("type", rtype)
			if e.ModuleName != "" {
				ev = ev.Str("module", e.ModuleName)
			}
			ev.Msg("provided")
		}
		if e.Err != nil {
			ev := l.logger.Error()
			if e.ModuleName != "" {
				ev = ev.Str("module", e.ModuleName)
			}
			ev.Err(e.Err).Msg("error encountered while applying options")
		}
	case *fxevent.Replaced:
		for _, rtype := range e.OutputTypeNames {
			ev := l.logger.Debug().Str("type", rtype)
			if e.ModuleName != "" {
				ev = ev.Str("module", e.ModuleName)
			}
			ev.Msg("replaced")
		}
		if e.Err != nil {
			ev := l.logger.Error()
			if e.ModuleName != "" {
				ev = ev.Str("module", e.ModuleName)
			}
			ev.Err(e.Err).Msg("error encountered while replacing")
		}
	case *fxevent.Decorated:
		for _, rtype := range e.OutputTypeNames {
			ev := l.logger.Debug().
				Str("decorator", e.DecoratorName).
				Str("type", rtype)
			if e.ModuleName != "" {
				ev = ev.Str("module", e.ModuleName)
			}
			ev.Msg("decorated")
		}
		if e.Err != nil {
			ev := l.logger.Error()
			if e.ModuleName != "" {
				ev = ev.Str("module", e.ModuleName)
			}
			ev.Err(e.Err).Msg("error encountered while applying options")
		}
	case *fxevent.Run:
		if e.Err != nil {
			ev := l.logger.Error().
				Str("name", e.Name).
				Str("kind", e.Kind)
			if e.ModuleName != "" {
				ev = ev.Str("module", e.ModuleName)
			}
			ev.Err(e.Err).Msg("error returned")
		} else {
			ev := l.logger.Debug().
				Str("name", e.Name).
				Str("kind", e.Kind).
				Str("runtime", e.Runtime.String())
			if e.ModuleName != "" {
				ev = ev.Str("module", e.ModuleName)
			}
			ev.Msg("run")
		}
	case *fxevent.Invoking:
		ev := l.logger.Debug().Str("function", e.FunctionName)
		if e.ModuleName != "" {
			ev = ev.Str("module", e.ModuleName)
		}
		ev.Msg("invoking")
	case *fxevent.Invoked:
		if e.Err != nil {
			ev := l.logger.Error().
				Err(e.Err).
				Str("stack", e.Trace).
				Str("function", e.FunctionName)
			if e.ModuleName != "" {
				ev = ev.Str("module", e.ModuleName)
			}
			ev.Msg("invoke failed")
		}
	case *fxevent.Stopping:
		l.logger.Info().
			Str("signal", strings.ToUpper(e.Signal.String())).
			Msg("received signal")
	case *fxevent.Stopped:
		if e.Err != nil {
			l.logger.Error().Err(e.Err).Msg("stop failed")
		}
	case *fxevent.RollingBack:
		l.logger.Error().Err(e.StartErr).Msg("start failed, rolling back")
	case *fxevent.RolledBack:
		if e.Err != nil {
			l.logger.Error().Err(e.Err).Msg("rollback failed")
		}
	case *fxevent.Started:
		if e.Err != nil {
			l.logger.Error().Err(e.Err).Msg("start failed")
		} else {
			l.logger.Debug().Msg("bootstrap complete")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			l.logger.Error().Err(e.Err).Msg("custom logger initialization failed")
		} else {
			l.logger.Debug().
				Str("function", e.ConstructorName).
				Msg("initialized custom fxevent.Logger")
		}
	}
}
