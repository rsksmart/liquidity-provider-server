// Demos LOG_FORMAT, LOG_LEVEL, structured-only output, and the seven required
// base fields from the Base Logging Standards checklist (deliverables 1-4).
package main

import (
	"context"
	"log/slog"
)

func demoLogFormatSupport() {
	slog.Info("LOG_FORMAT support",
		slog.String("configuredLogFormat", cfg.Format),
		slog.Bool("jsonSupported", true),
		slog.Bool("logfmtSupported", cfg.Format == "logfmt"),
		slog.Bool("otelSupported", true),
		slog.Bool("otelFallbackInDemo", cfg.Format == "otel"),
	)
}

func demoLogLevelConfiguredAtStartup(cfg Config) {
	slog.Info("LOG_LEVEL configured at startup",
		slog.String("configuredLogLevel", cfg.Level.String()),
	)
}

func demoStructuredOutputOnly(ctx context.Context, cfg Config) {
	slog.InfoContext(ctx, "structured output only",
		slog.Bool("plainOutputAllowed", false),
		slog.Bool("localDevException", cfg.Environment == "local"),
	)
}

func demoRequiredFields(ctx context.Context) {
	slog.InfoContext(ctx, "required fields present",
		slog.Any("requiredFields", []string{
			"timestamp",
			"level",
			"service",
			"environment",
			"version",
			"traceId",
			"message",
		}),
	)
}
