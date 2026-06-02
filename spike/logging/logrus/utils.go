// Shared logger setup: env-driven Config, global logrus logger, and baseEntry for
// init-time base fields plus per-request trace fields on every call site.
package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Format      string
	Level       logrus.Level
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
	lvl, err := logrus.ParseLevel(envDefault("LOG_LEVEL", "info"))
	if err != nil {
		lvl = logrus.InfoLevel
	}
	cfg.Level = lvl
	return cfg
}

// initLogger builds the global logger. Because logrus has no native
// context handler, base fields must live on a *logrus.Entry that every
// call site has to thread through (see baseEntry below). That is the
// "Base-field auto-attachment" weakness called out in logs.md.
func initLogger(out io.Writer, cfg Config) *logrus.Logger {
	l := logrus.New()
	l.SetOutput(out)
	l.SetLevel(cfg.Level)
	switch cfg.Format {
	case "json", "otel":
		// "otel" has no first-party logrus bridge — log a notice and
		// fall back to JSON for the spike.
		if cfg.Format == "otel" {
			fmt.Fprintln(os.Stderr, "LOG_FORMAT=otel has no maintained logrus bridge; falling back to JSON")
		}
		l.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime: "timestamp",
				logrus.FieldKeyMsg:  "message",
			},
		})
	case "text", "logfmt":
		l.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: time.RFC3339Nano,
			DisableColors:   true,
			FullTimestamp:   true,
		})
	default:
		l.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339Nano})
	}
	l.AddHook(&piiRedactHook{})
	return l
}

// baseEntry returns a *logrus.Entry pre-populated with the three init
// fields and (if present) the trace fields from ctx. Every call site
// in a "real" logrus migration would call this — there is no slog-style
// handler that adds them automatically.
func baseEntry(ctx context.Context) *logrus.Entry {
	e := logger.WithFields(logrus.Fields{
		"service":     cfg.Service,
		"environment": cfg.Environment,
		"version":     cfg.Version,
	})
	if tc, ok := traceFromContext(ctx); ok {
		e = e.WithFields(logrus.Fields{
			"traceId": tc.TraceID,
			"spanId":  tc.SpanID,
		})
	}
	return e
}

func envDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// contextless returns context.Background for the non-request demos. The
// logrus call sites still have to take a context to thread base fields
// through baseEntry, so this saves repeating context.Background() at
// every non-request demo site.
func contextless() context.Context {
	return context.Background()
}

// Package-level state. logrus's ergonomic API is package-global; storing
// cfg and logger here matches how a real migration would use it.
var (
	cfg    Config
	logger *logrus.Logger
)
