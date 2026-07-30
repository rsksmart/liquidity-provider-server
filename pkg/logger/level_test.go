package logger_test

import (
	"testing"

	"github.com/rsksmart/liquidity-provider-server/pkg/logger"
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
	assert.Equal(t, "info", logger.LevelInfo.String())
	assert.Equal(t, "fatal", logger.LevelFatal.String())
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
