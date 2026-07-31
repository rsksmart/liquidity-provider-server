package otel

import (
	"context"
	"log/slog"
	"slices"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/logger/internal/redact"
)

// redactHandler applies the redactor to every attribute before forwarding to
// the wrapped handler. It is used for the OTel path, which has no
// ReplaceAttr hook equivalent to the JSON/logfmt handlers.
type redactHandler struct {
	slog.Handler
	redactor *redact.Redactor
	min      slog.Level
	groups   []string
}

// Enabled reports whether the level is at or above the configured minimum and
// the wrapped handler would also emit it.
func (h redactHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.min && h.Handler.Enabled(ctx, level)
}

// Handle redacts every attribute on the record (including nested groups) and
// forwards the rebuilt record to the wrapped handler.
func (h redactHandler) Handle(ctx context.Context, record slog.Record) error {
	out := slog.NewRecord(record.Time, record.Level, record.Message, record.PC)
	record.Attrs(func(a slog.Attr) bool {
		out.AddAttrs(h.redactAttr(h.groups, a)...)
		return true
	})
	return h.Handler.Handle(ctx, out)
}

// WithAttrs redacts the given attributes and returns a derived handler.
func (h redactHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	redacted := make([]slog.Attr, 0, len(attrs))
	for _, a := range attrs {
		redacted = append(redacted, h.redactAttr(h.groups, a)...)
	}
	return redactHandler{
		Handler:  h.Handler.WithAttrs(redacted),
		redactor: h.redactor,
		min:      h.min,
		groups:   h.groups,
	}
}

// WithGroup returns a handler that nests subsequent attributes under name.
func (h redactHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return redactHandler{
		Handler:  h.Handler.WithGroup(name),
		redactor: h.redactor,
		min:      h.min,
		groups:   append(slices.Clone(h.groups), name),
	}
}

// redactAttr applies redaction to a, recursing into groups. Empty attributes
// (discarded by the redactor) are omitted from the result.
func (h redactHandler) redactAttr(groups []string, a slog.Attr) []slog.Attr {
	a.Value = a.Value.Resolve()
	if a.Value.Kind() == slog.KindGroup {
		return h.redactGroup(groups, a)
	}

	a = h.redactor.Apply(groups, a)
	if a.Equal(slog.Attr{}) {
		return nil
	}
	return []slog.Attr{a}
}

// redactGroup redacts every child of a group attribute. An empty-keyed group
// is inlined; a named group is rebuilt around the redacted children.
func (h redactHandler) redactGroup(groups []string, a slog.Attr) []slog.Attr {
	nested := groups
	if a.Key != "" {
		nested = append(slices.Clone(groups), a.Key)
	}

	children := a.Value.Group()
	out := make([]slog.Attr, 0, len(children))
	for _, child := range children {
		out = append(out, h.redactAttr(nested, child)...)
	}
	if a.Key == "" {
		return out
	}
	return []slog.Attr{slog.Group(a.Key, attrsToAny(out)...)}
}

// attrsToAny converts a slice of attributes into the []any expected by
// slog.Group.
func attrsToAny(attrs []slog.Attr) []any {
	out := make([]any, len(attrs))
	for i, a := range attrs {
		out[i] = a
	}
	return out
}
