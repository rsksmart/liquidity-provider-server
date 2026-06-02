// Shared logger setup: env-driven Config, zap core stack with piiCore, init-time
// base fields, and loggerWithTrace for per-request trace fields.
package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Format      string
	Level       zapcore.Level
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
		cfg.Level = zapcore.DebugLevel
	case "warn", "warning":
		cfg.Level = zapcore.WarnLevel
	case "error":
		cfg.Level = zapcore.ErrorLevel
	default:
		cfg.Level = zapcore.InfoLevel
	}
	return cfg
}

// initLogger builds the zap core stack:
//
//	piiCore (drops PII fields)
//	  -> base core (JSON or console encoder)
//
// Base fields are attached with zap.Fields and then frozen onto the
// returned logger. Trace fields ride along by way of loggerWithTrace.
func initLogger(out io.Writer, cfg Config) *zap.Logger {
	enc := buildEncoder(cfg.Format)
	core := zapcore.NewCore(enc, zapcore.AddSync(out), cfg.Level)
	core = piiCore{Core: core}

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(0)).With(
		zap.String("service", cfg.Service),
		zap.String("environment", cfg.Environment),
		zap.String("version", cfg.Version),
	)
	return logger
}

func buildEncoder(format string) zapcore.Encoder {
	encCfg := zap.NewProductionEncoderConfig()
	encCfg.TimeKey = "timestamp"
	encCfg.MessageKey = "message"
	encCfg.LevelKey = "level"
	encCfg.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	encCfg.EncodeLevel = zapcore.LowercaseLevelEncoder
	switch format {
	case "text", "logfmt":
		// Console is human-readable, not strict logfmt; a real
		// logfmt would need a custom encoder.
		return zapcore.NewConsoleEncoder(encCfg)
	case "otel":
		fmt.Fprintln(os.Stderr, "LOG_FORMAT=otel: would use otelzap; falling back to JSON for the spike")
		return zapcore.NewJSONEncoder(encCfg)
	default:
		return zapcore.NewJSONEncoder(encCfg)
	}
}

// loggerWithTrace returns the package logger with traceId/spanId
// attached if ctx carries a TraceContext. zap has no native handler hook
// that reads ctx — call sites must opt in.
func loggerWithTrace(ctx context.Context) *zap.Logger {
	if tc, ok := traceFromContext(ctx); ok {
		return logger.With(zap.String("traceId", tc.TraceID), zap.String("spanId", tc.SpanID))
	}
	return logger
}

func envDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

var (
	cfg    Config
	logger *zap.Logger
)
