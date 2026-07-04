package logger

import (
	"context"
	"log/slog"

	"github.com/rs/zerolog"
)

var _ slog.Handler = (*zerologHandler)(nil)

// zerologHandler is an internal slog.Handler that delegates to zerolog.
type zerologHandler struct {
	logger zerolog.Logger
	attrs  []slog.Attr
	group  string
}

func newZerologHandler(l zerolog.Logger) slog.Handler {
	return &zerologHandler{logger: l}
}

func (h *zerologHandler) Enabled(_ context.Context, level slog.Level) bool {
	return h.logger.GetLevel() <= slogToZerologLevel(level)
}

func (h *zerologHandler) Handle(_ context.Context, r slog.Record) error {
	ev := h.logger.WithLevel(slogToZerologLevel(r.Level))
	if ev == nil {
		return nil
	}
	for _, a := range h.attrs {
		ev = addAttr(ev, h.group, a)
	}
	r.Attrs(func(a slog.Attr) bool {
		ev = addAttr(ev, h.group, a)
		return true
	})
	ev.Msg(r.Message)
	return nil
}

func (h *zerologHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newH := *h
	newH.attrs = make([]slog.Attr, len(h.attrs)+len(attrs))
	copy(newH.attrs, h.attrs)
	copy(newH.attrs[len(h.attrs):], attrs)
	return &newH
}

func (h *zerologHandler) WithGroup(name string) slog.Handler {
	newH := *h
	newH.group = name
	return &newH
}

func slogToZerologLevel(l slog.Level) zerolog.Level {
	switch {
	case l >= slog.LevelError:
		return zerolog.ErrorLevel
	case l >= slog.LevelWarn:
		return zerolog.WarnLevel
	case l >= slog.LevelInfo:
		return zerolog.InfoLevel
	default:
		return zerolog.DebugLevel
	}
}

func addAttr(ev *zerolog.Event, group string, a slog.Attr) *zerolog.Event {
	key := a.Key
	if group != "" {
		key = group + "." + key
	}
	v := a.Value.Resolve()
	switch v.Kind() {
	case slog.KindString:
		return ev.Str(key, v.String())
	case slog.KindInt64:
		return ev.Int64(key, v.Int64())
	case slog.KindFloat64:
		return ev.Float64(key, v.Float64())
	case slog.KindBool:
		return ev.Bool(key, v.Bool())
	case slog.KindTime:
		return ev.Time(key, v.Time())
	case slog.KindDuration:
		return ev.Dur(key, v.Duration())
	case slog.KindAny:
		if err, ok := v.Any().(error); ok {
			return ev.Err(err)
		}
		return ev.Interface(key, v.Any())
	default:
		return ev.Interface(key, v.Any())
	}
}
