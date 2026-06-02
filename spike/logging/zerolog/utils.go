// Shared logger setup: env-driven Config, zerolog logger with piiWriter, and
// loggerWithTrace for init-time base fields plus per-request trace fields.
package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type Config struct {
	Format      string
	Level       zerolog.Level
	Service     string
	Environment string
	Version     string
}

func readConfig() Config {
	cfg := Config{
		Format:      strings.ToLower(envDefault("LOG_FORMAT", "json")),
		Service:     envDefault("LOG_SERVICE", "liquidity-provider-server"),
		Environment: envDefault("LOG_ENVIRONMENT", "spike"),
		Version:     envDefault("LOG_VERSION", "spike-fly-2308"),
	}
	switch strings.ToLower(envDefault("LOG_LEVEL", "info")) {
	case "debug":
		cfg.Level = zerolog.DebugLevel
	case "warn", "warning":
		cfg.Level = zerolog.WarnLevel
	case "error":
		cfg.Level = zerolog.ErrorLevel
	default:
		cfg.Level = zerolog.InfoLevel
	}
	return cfg
}

// initLogger builds the zerolog logger. zerolog's fluent API does not
// expose a record-level redaction hook the way slog's ReplaceAttr does;
// instead we wrap the io.Writer (see piiWriter in pii.go) so denied
// keys are stripped on the way out. This is why zerolog ranks weakest
// on the "PII deny-list interception" axis.
//
// Note: zerolog.MessageFieldName / TimestampFieldName / LevelFieldName
// are process-global; setting them here flows to any zerolog usage in
// the same binary, which is fine for the spike but worth knowing.
func initLogger(out io.Writer, cfg Config) zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.MessageFieldName = "message"
	zerolog.TimestampFieldName = "timestamp"
	zerolog.LevelFieldName = "level"

	var w io.Writer = out
	switch cfg.Format {
	case "text", "logfmt":
		w = zerolog.ConsoleWriter{Out: out, TimeFormat: time.RFC3339Nano}
	case "otel":
		fmt.Fprintln(os.Stderr, "LOG_FORMAT=otel has no first-party zerolog bridge; falling back to JSON")
	}
	w = piiWriter{Out: w}

	return zerolog.New(w).Level(cfg.Level).With().
		Timestamp().
		Str("service", cfg.Service).
		Str("environment", cfg.Environment).
		Str("version", cfg.Version).
		Logger()
}

// loggerWithTrace returns a child logger carrying traceId/spanId. Note
// that the result is returned by pointer because zerolog's event
// builders (Info/Warn/Error/...) are defined on *zerolog.Logger — a
// small ergonomics tax that the other libraries don't have.
//
// zerolog has no native ctx-handler — call sites opt in.
func loggerWithTrace(ctx context.Context) *zerolog.Logger {
	if tc, ok := traceFromContext(ctx); ok {
		l := logger.With().Str("traceId", tc.TraceID).Str("spanId", tc.SpanID).Logger()
		return &l
	}
	return &logger
}

func envDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

var (
	cfg    Config
	logger zerolog.Logger
)
