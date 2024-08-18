package log

import (
	"context"
	"log/slog"
)

// Noop is a [Handler] which is always disabled and therefore logs nothing.
var Noop slog.Handler = noopHandler{}

type noopHandler struct{}

func (noopHandler) Enabled(context.Context, slog.Level) bool  { return false }
func (noopHandler) Handle(context.Context, slog.Record) error { return nil }
func (d noopHandler) WithAttrs([]slog.Attr) slog.Handler      { return d }
func (d noopHandler) WithGroup(string) slog.Handler           { return d }

// NoopLogger returns a new slog.Logger that logs nothing.
func NoopLogger() *slog.Logger {
	return slog.New(Noop)
}
