package handler

import (
	"context"
	"io"
	"log/slog"
	"time"

	otellog "go.opentelemetry.io/otel/log"

	"github.com/rsksmart/liquidity-provider-server/pkg/logger/internal/core"
	otelhandler "github.com/rsksmart/liquidity-provider-server/pkg/logger/internal/handler/otel"
	"github.com/rsksmart/liquidity-provider-server/pkg/logger/internal/redact"
	"github.com/rsksmart/liquidity-provider-server/pkg/logger/internal/trace"
)

// contextHandler is the outermost handler. It injects the traceId (always)
// and spanId (when available) read from the request context, then delegates
// to the wrapped handler.
type contextHandler struct {
	slog.Handler
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

	return contextHandler{Handler: inner}.WithAttrs(base)
}

// Handle adds the trace identifiers to the record and forwards it.
func (h contextHandler) Handle(ctx context.Context, record slog.Record) error {
	traceID, spanID := "", ""
	if tc, ok := trace.FromContext(ctx); ok {
		traceID, spanID = tc.TraceID, tc.SpanID
	}
	record.AddAttrs(slog.String(string(core.FieldTraceID), traceID))
	if spanID != "" {
		record.AddAttrs(slog.String(string(core.FieldSpanID), spanID))
	}
	return h.Handler.Handle(ctx, record)
}

// WithAttrs preserves the contextHandler wrapper around the derived handler.
func (h contextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return contextHandler{Handler: h.Handler.WithAttrs(attrs)}
}

// WithGroup preserves the contextHandler wrapper around the derived handler.
func (h contextHandler) WithGroup(name string) slog.Handler {
	return contextHandler{Handler: h.Handler.WithGroup(name)}
}
