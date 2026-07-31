package otel

import (
	"log/slog"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	otellog "go.opentelemetry.io/otel/log"

	"github.com/rsksmart/liquidity-provider-server/pkg/logger/internal/redact"
)

// Options configures the OpenTelemetry handler backend.
type Options struct {
	Service  string
	Level    slog.Level
	Redactor *redact.Redactor
	Provider otellog.LoggerProvider
}

// New builds an otelslog handler with redaction and minimum-level filtering.
// Delivery is owned by the consumer's LoggerProvider exporter.
func New(opts Options) slog.Handler {
	var bridgeOpts []otelslog.Option
	if opts.Provider != nil {
		bridgeOpts = append(bridgeOpts, otelslog.WithLoggerProvider(opts.Provider))
	}
	return redactHandler{
		Handler:  otelslog.NewHandler(opts.Service, bridgeOpts...),
		redactor: opts.Redactor,
		min:      opts.Level,
	}
}
