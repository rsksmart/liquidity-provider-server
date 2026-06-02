// Shared logger setup: env-driven Config, handler stack (JSON/text/otel fallback),
// init-time base fields, and contextHandler for per-request traceId.
package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

// Config holds the env-driven settings shared across all four demos.
type Config struct {
	Format      string // json | text | logfmt | otel
	Level       slog.Level
	Service     string
	Environment string
	Version     string
}

// readConfig parses the env vars listed in the Base Logging Standards.
func readConfig() Config {
	cfg := Config{
		Format:      strings.ToLower(envDefault("LOG_FORMAT", "json")),
		Service:     envDefault("LOG_SERVICE", "liquidity-provider-server"),
		Environment: envDefault("LOG_ENVIRONMENT", "spike"),
		Version:     envDefault("LOG_VERSION", "spike-fly-2308"),
	}
	switch strings.ToLower(envDefault("LOG_LEVEL", "info")) {
	case "debug":
		cfg.Level = slog.LevelDebug
	case "warn", "warning":
		cfg.Level = slog.LevelWarn
	case "error":
		cfg.Level = slog.LevelError
	default:
		cfg.Level = slog.LevelInfo
	}
	return cfg
}

// initLogger builds the handler stack and installs it as slog.Default.
//
// The stack is, outermost to innermost:
//
//	contextHandler (adds traceId from ctx)
//	  -> base handler (JSON or text) with:
//	       HandlerOptions.ReplaceAttr = piiRedactor (drops PII keys)
//
// Base fields (service, environment, version) are attached once via
// logger.With and ride along on every record without being repeated at
// call sites — see docs/spikes/logs.md "Base-field auto-attachment".
func initLogger(out io.Writer, cfg Config) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level:       cfg.Level,
		ReplaceAttr: replaceAttr,
		AddSource:   false,
	}

	var base slog.Handler
	switch cfg.Format {
	case "text", "logfmt":
		// slog.NewTextHandler is logfmt-ish (key=value, space-separated).
		// Exact RFC-5424-style logfmt would need a third-party handler;
		// see logs.md "Format switch ergonomics".
		base = slog.NewTextHandler(out, opts)
	case "otel":
		// Paper-only in the SPIKE: a real wiring uses
		// go.opentelemetry.io/contrib/bridges/otelslog. We log a notice
		// and fall back to JSON so the demo stays self-contained.
		fmt.Fprintln(os.Stderr, "LOG_FORMAT=otel: would use otelslog.NewHandler; falling back to JSON for the spike")
		base = slog.NewJSONHandler(out, opts)
	default:
		base = slog.NewJSONHandler(out, opts)
	}

	logger := slog.New(contextHandler{Handler: base}).With(
		slog.String("service", cfg.Service),
		slog.String("environment", cfg.Environment),
		slog.String("version", cfg.Version),
	)
	slog.SetDefault(logger)
	return logger
}

// contextHandler enriches every record with the request-scoped traceId
// (and spanId, when present) so call sites don't pass it explicitly.
// This is the cleanest demonstration of slog's "context.Context
// propagation" win in the comparison.
type contextHandler struct{ slog.Handler }

func (h contextHandler) Handle(ctx context.Context, r slog.Record) error {
	if tc, ok := traceFromContext(ctx); ok {
		r.AddAttrs(
			slog.String("traceId", tc.TraceID),
			slog.String("spanId", tc.SpanID),
		)
	}
	return h.Handler.Handle(ctx, r)
}

func (h contextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return contextHandler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h contextHandler) WithGroup(name string) slog.Handler {
	return contextHandler{Handler: h.Handler.WithGroup(name)}
}

func envDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

var cfg Config
