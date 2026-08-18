package handler

import (
	"context"
	"io"
	"log/slog"
	"slices"
	"time"

	otellog "go.opentelemetry.io/otel/log"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/logger/internal/core"
	otelhandler "github.com/rsksmart/liquidity-provider-server/internal/adapters/logger/internal/handler/otel"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/logger/internal/redact"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/logger/internal/trace"
)

// contextHandler is the outermost handler. It injects the traceId (always)
// and spanId (when available) at the root.
type contextHandler struct {
	inner slog.Handler
	steps []handlerStep
}

// handlerStep records a WithGroup or WithAttrs call that must run after
// root-level trace injection so slog does not nest traceId under an open group.
type handlerStep struct {
	group string
	attrs []slog.Attr
}

// Options declares how the handler chain is built. It is populated by the
// logger facade from its Config.
type Options struct {
	Output      io.Writer
	Format      core.Format
	Level       core.Level
	AddSource   bool
	TimeZone    *time.Location
	Redactor    *redact.Redactor
	Service     string
	Environment string
	Version     string
	// OTelProvider receives records when Format is FormatOTel. If nil, the
	// OTel global LoggerProvider is used. Ignored for other formats.
	OTelProvider otellog.LoggerProvider
}

// New assembles the handler chain for the given options.
func New(opts Options) slog.Handler {
	ab := &attrBuilder{redactor: opts.Redactor, tz: opts.TimeZone}

	slogOpts := &slog.HandlerOptions{
		Level:       opts.Level.SlogLevel(),
		AddSource:   opts.AddSource,
		ReplaceAttr: ab.replaceAttr,
	}

	var inner slog.Handler
	switch opts.Format {
	case core.FormatLogfmt:
		inner = slog.NewTextHandler(opts.Output, slogOpts)
	case core.FormatOTel:
		inner = otelhandler.New(otelhandler.Options{
			Service:  opts.Service,
			Level:    opts.Level.SlogLevel(),
			Redactor: opts.Redactor,
			Provider: opts.OTelProvider,
		})
	case core.FormatJSON:
		inner = slog.NewJSONHandler(opts.Output, slogOpts)
	default:
		inner = slog.NewJSONHandler(opts.Output, slogOpts)
	}

	base := []slog.Attr{
		slog.String(string(core.FieldService), opts.Service),
		slog.String(string(core.FieldEnvironment), opts.Environment),
		slog.String(string(core.FieldVersion), opts.Version),
	}

	return contextHandler{inner: inner}.WithAttrs(base)
}

// Enabled reports whether the inner handler enables the level.
func (h contextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

// Handle adds root-level trace identifiers, then applies deferred group/attr
// steps and forwards the record.
func (h contextHandler) Handle(ctx context.Context, record slog.Record) error {
	handler := h.inner.WithAttrs(h.traceAttrs(ctx))
	for _, step := range h.steps {
		if step.group != "" {
			handler = handler.WithGroup(step.group)
		}
		if len(step.attrs) > 0 {
			handler = handler.WithAttrs(step.attrs)
		}
	}
	return handler.Handle(ctx, record)
}

// WithAttrs preserves the contextHandler wrapper. Attrs applied before any
// group land on the inner handler immediately (top-level). Attrs after a
// group are deferred so they nest correctly while traceId stays at the root.
func (h contextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	if len(h.steps) == 0 {
		return contextHandler{inner: h.inner.WithAttrs(attrs)}
	}
	return contextHandler{
		inner: h.inner,
		steps: append(slices.Clip(h.steps), handlerStep{attrs: attrs}),
	}
}

// WithGroup defers opening the group until Handle so trace injection can run
// on the ungrouped inner handler first.
func (h contextHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return contextHandler{
		inner: h.inner,
		steps: append(slices.Clip(h.steps), handlerStep{group: name}),
	}
}

// traceAttrs returns the root-level correlation fields for this record.
// Missing context yields a freshly generated W3C trace so every line carries
// a usable traceId.
func (h contextHandler) traceAttrs(ctx context.Context) []slog.Attr {
	traceID, spanID := "", ""
	if tc, ok := trace.FromContext(ctx); ok && tc.TraceID != "" {
		traceID, spanID = tc.TraceID, tc.SpanID
	} else if tc, err := trace.NewTraceContext(); err == nil {
		traceID, spanID = tc.TraceID, tc.SpanID
	}

	attrs := []slog.Attr{slog.String(string(core.FieldTraceID), traceID)}
	if spanID != "" {
		attrs = append(attrs, slog.String(string(core.FieldSpanID), spanID))
	}
	return attrs
}
