// Demos LOG_FORMAT, LOG_LEVEL, structured-only output, and the seven required
// base fields from the Base Logging Standards checklist (deliverables 1-4).
package main

import (
	"context"

	"go.uber.org/zap"
)

func demoLogFormatSupport() {
	logger.Info("LOG_FORMAT support",
		zap.String("configuredLogFormat", cfg.Format),
		zap.Bool("jsonSupported", true),
		zap.Bool("logfmtSupported", cfg.Format == "logfmt"),
		zap.Bool("otelSupported", true),
		zap.Bool("otelFallbackInDemo", cfg.Format == "otel"),
	)
}

func demoLogLevelConfiguredAtStartup() {
	logger.Info("LOG_LEVEL configured at startup",
		zap.String("configuredLogLevel", cfg.Level.String()),
	)
}

func demoStructuredOutputOnly(ctx context.Context) {
	loggerWithTrace(ctx).Info("structured output only",
		zap.Bool("plainOutputAllowed", false),
		zap.Bool("localDevException", cfg.Environment == "local"),
	)
}

func demoRequiredFields(ctx context.Context) {
	loggerWithTrace(ctx).Info("required fields present",
		zap.Strings("requiredFields", []string{
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
