package logging_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/logging"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/stretchr/testify/require"
)

func TestInit_ParsesFormats(t *testing.T) {
	swapDefaultLogger(t)

	for _, format := range []string{"json", "text", "logfmt"} {
		t.Run(format, func(t *testing.T) {
			require.NoError(t, logging.Init(logging.Config{
				Service:     "liquidity-provider-server",
				Environment: "regtest",
				Version:     "v1",
				Format:      format,
				Level:       "info",
			}))
		})
	}
}

func TestInit_RejectsOtelFormat(t *testing.T) {
	swapDefaultLogger(t)
	err := logging.Init(logging.Config{
		Service: "liquidity-provider-server",
		Format:  "otel",
		Level:   "info",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "not supported yet")
}

func TestInit_RejectsInvalidLevel(t *testing.T) {
	swapDefaultLogger(t)
	err := logging.Init(logging.Config{
		Service: "liquidity-provider-server",
		Format:  "json",
		Level:   "panic",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported log level")
}

func TestInit_AttachesBaseFields(t *testing.T) {
	var buf bytes.Buffer
	swapDefaultLogger(t)
	level := logging.LevelVar()
	level.Set(slog.LevelInfo)
	logger := slog.New(logging.NewTestHandler(&buf, level)).With(
		slog.String("service", "liquidity-provider-server"),
		slog.String("environment", "testnet"),
		slog.String("version", "v2.6.0"),
	)

	logger.Info("startup")

	entry := decodeEntry(t, buf.Bytes())
	require.Equal(t, "liquidity-provider-server", entry["service"])
	require.Equal(t, "testnet", entry["environment"])
	require.Equal(t, "v2.6.0", entry["version"])
}

func TestLevelVar_FiltersByLevel(t *testing.T) {
	var buf bytes.Buffer
	level := logging.LevelVar()
	level.Set(slog.LevelWarn)
	logger := slog.New(logging.NewTestHandler(&buf, level))

	logger.Debug("hidden")
	logger.Warn("visible")

	require.NotContains(t, buf.String(), "hidden")
	require.Contains(t, buf.String(), "visible")
}

func TestError_EmitsStandardShape(t *testing.T) {
	var buf bytes.Buffer
	swapDefaultLogger(t)
	level := logging.LevelVar()
	level.Set(slog.LevelInfo)
	slog.SetDefault(slog.New(logging.NewTestHandler(&buf, level)))

	wrapped := usecases.WrapUseCaseError(usecases.AcceptPeginQuoteId, usecases.ExpiredQuoteError)
	logging.Error(context.Background(), "accept failed", wrapped)

	entry := decodeEntry(t, buf.Bytes())
	require.Equal(t, "ERROR", entry["level"])
	require.Equal(t, "accept failed", entry["msg"])
	require.Equal(t, usecases.ExpiredQuoteError.Error(), entry["errorCode"])
	require.Equal(t, wrapped.Error(), entry["errorMessage"])
	require.Equal(t, string(usecases.AcceptPeginQuoteId), entry["errorContext"])
	require.NotEmpty(t, entry["errorStack"])
}

func TestError_PreservesErrorsIs(t *testing.T) {
	swapDefaultLogger(t)
	level := logging.LevelVar()
	level.Set(slog.LevelInfo)
	slog.SetDefault(slog.New(logging.NewTestHandler(&bytes.Buffer{}, level)))

	wrapped := usecases.WrapUseCaseError(usecases.AcceptPeginQuoteId, usecases.ExpiredQuoteError)
	logging.Error(context.Background(), "accept failed", wrapped)

	require.ErrorIs(t, wrapped, usecases.ExpiredQuoteError)
}

func TestError_HandlesJoinedErrors(t *testing.T) {
	var buf bytes.Buffer
	swapDefaultLogger(t)
	level := logging.LevelVar()
	level.Set(slog.LevelInfo)
	slog.SetDefault(slog.New(logging.NewTestHandler(&buf, level)))

	joined := errors.Join(
		usecases.WrapUseCaseError(usecases.SendPegoutId, usecases.NonRecoverableError),
		usecases.QuoteNotFoundError,
	)
	logging.Error(context.Background(), "send failed", joined)

	entry := decodeEntry(t, buf.Bytes())
	require.Equal(t, usecases.NonRecoverableError.Error(), entry["errorCode"])
}

func TestLazy_DefersEvaluation(t *testing.T) {
	var buf bytes.Buffer
	level := logging.LevelVar()
	level.Set(slog.LevelInfo)
	logger := slog.New(logging.NewTestHandler(&buf, level))

	called := false
	logger.Info("lazy", "value", logging.Lazy(func() any {
		called = true
		return "expensive"
	}))
	require.True(t, called)
	require.Contains(t, buf.String(), "expensive")
}

func TestInit_WritesToLogFile(t *testing.T) {
	swapDefaultLogger(t)
	logPath := filepath.Join(t.TempDir(), "logs", "lps.log")
	require.NoError(t, logging.Init(logging.Config{
		Service:     "liquidity-provider-server",
		Environment: "regtest",
		Version:     "test",
		Format:      "json",
		Level:       "info",
		File:        logPath,
	}))

	slog.Info("file output test")

	content, err := os.ReadFile(logPath)
	require.NoError(t, err)
	require.Contains(t, string(content), "file output test")
}

func swapDefaultLogger(t *testing.T) {
	t.Helper()
	previous := slog.Default()
	logging.ResetLevelVar()
	t.Cleanup(func() {
		slog.SetDefault(previous)
		logging.ResetLevelVar()
	})
}

func decodeEntry(t *testing.T, raw []byte) map[string]any {
	t.Helper()
	var entry map[string]any
	require.NoError(t, json.Unmarshal(raw, &entry))
	return entry
}
