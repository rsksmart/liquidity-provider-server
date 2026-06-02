// Demos LOG_FORMAT, LOG_LEVEL, structured-only output, and the seven required
// base fields from the Base Logging Standards checklist (deliverables 1-4).
package main

import (
	"context"

	"github.com/sirupsen/logrus"
)

func demoLogFormatSupport() {
	baseEntry(contextless()).WithFields(logrus.Fields{
		"configuredLogFormat": cfg.Format,
		"jsonSupported":       true,
		"logfmtSupported":     cfg.Format == "logfmt",
		"otelSupported":       false,
		"otelFallback":        cfg.Format == "otel",
	}).Info("LOG_FORMAT support")
}

func demoLogLevelConfiguredAtStartup() {
	baseEntry(contextless()).WithFields(logrus.Fields{
		"configuredLogLevel": cfg.Level.String(),
	}).Info("LOG_LEVEL configured at startup")
}

func demoStructuredOutputOnly(ctx context.Context) {
	baseEntry(ctx).WithFields(logrus.Fields{
		"plainOutputAllowed": false,
		"localDevException":  cfg.Environment == "local",
	}).Info("structured output only")
}

func demoRequiredFields(ctx context.Context) {
	baseEntry(ctx).WithFields(logrus.Fields{
		"requiredFields": []string{
			"timestamp",
			"level",
			"service",
			"environment",
			"version",
			"traceId",
			"message",
		},
	}).Info("required fields present")
}
