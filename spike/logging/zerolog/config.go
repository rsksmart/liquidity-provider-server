// Demos LOG_FORMAT, LOG_LEVEL, structured-only output, and the seven required
// base fields from the Base Logging Standards checklist (deliverables 1-4).
package main

import "context"

func demoLogFormatSupport() {
	logger.Info().
		Str("configuredLogFormat", cfg.Format).
		Bool("jsonSupported", true).
		Bool("logfmtSupported", cfg.Format == "logfmt").
		Bool("otelSupported", false).
		Bool("otelFallback", cfg.Format == "otel").
		Msg("LOG_FORMAT support")
}

func demoLogLevelConfiguredAtStartup() {
	logger.Info().
		Str("configuredLogLevel", cfg.Level.String()).
		Msg("LOG_LEVEL configured at startup")
}

func demoStructuredOutputOnly(ctx context.Context) {
	loggerWithTrace(ctx).Info().
		Bool("plainOutputAllowed", false).
		Bool("localDevException", cfg.Environment == "local").
		Msg("structured output only")
}

func demoRequiredFields(ctx context.Context) {
	loggerWithTrace(ctx).Info().
		Strs("requiredFields", []string{
			"timestamp",
			"level",
			"service",
			"environment",
			"version",
			"traceId",
			"message",
		}).
		Msg("required fields present")
}
