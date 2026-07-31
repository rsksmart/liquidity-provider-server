package logger_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/logger"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseLevel(t *testing.T) {
	cases := map[string]logger.Level{
		"trace":   logger.LevelTrace,
		"DEBUG":   logger.LevelDebug,
		"info":    logger.LevelInfo,
		"":        logger.LevelInfo,
		"warning": logger.LevelWarn,
		"error":   logger.LevelError,
		"fatal":   logger.LevelFatal,
	}
	for input, expected := range cases {
		got, err := logger.ParseLevel(input)
		require.NoError(t, err, input)
		assert.Equal(t, expected, got, input)
	}
}

func TestParseLevelUnknown(t *testing.T) {
	got, err := logger.ParseLevel("verbose")

	require.Error(t, err)
	assert.Equal(t, logger.LevelInfo, got)
}

func TestLevelString(t *testing.T) {
	cases := map[logger.Level]string{
		logger.LevelTrace: "trace",
		logger.LevelDebug: "debug",
		logger.LevelInfo:  "info",
		logger.LevelWarn:  "warn",
		logger.LevelError: "error",
		logger.LevelFatal: "fatal",
	}
	for level, want := range cases {
		assert.Equal(t, want, level.String(), want)
	}
	// Non-canonical levels fall back to slog's naming.
	assert.Equal(t, "ERROR+2", logger.Level(10).String())
}

func TestConfigFromEnv(t *testing.T) {
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("LOG_FORMAT", "logfmt")

	cfg := logger.ConfigFromEnv("svc", "production", "v1")

	assert.Equal(t, "svc", cfg.Service)
	assert.Equal(t, logger.LevelDebug, cfg.Level)
	assert.Equal(t, logger.FormatLogfmt, cfg.Format)
	// Stack traces are on by default (DisableStackTrace defaults to false).
	assert.False(t, cfg.DisableStackTrace)
}

func TestConfigFromEnvDefaults(t *testing.T) {
	t.Setenv("LOG_LEVEL", "")
	t.Setenv("LOG_FORMAT", "")

	cfg := logger.ConfigFromEnv("svc", "production", "v1")

	assert.Equal(t, logger.LevelInfo, cfg.Level)
	assert.Equal(t, logger.FormatJSON, cfg.Format)
}

func TestConfigFromEnvOTelAndInvalidLevel(t *testing.T) {
	t.Setenv("LOG_LEVEL", "verbose")
	t.Setenv("LOG_FORMAT", "otel")

	cfg := logger.ConfigFromEnv("svc", "production", "v1")

	assert.Equal(t, logger.LevelInfo, cfg.Level) // invalid level keeps default
	assert.Equal(t, logger.FormatOTel, cfg.Format)
}

func TestConfigFromEnvJSONFormat(t *testing.T) {
	t.Setenv("LOG_LEVEL", "error")
	t.Setenv("LOG_FORMAT", "json")

	cfg := logger.ConfigFromEnv("svc", "production", "v1")

	assert.Equal(t, logger.LevelError, cfg.Level)
	assert.Equal(t, logger.FormatJSON, cfg.Format)
}
