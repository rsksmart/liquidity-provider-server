package logging_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

	for _, format := range []string{"json", "text", "logfmt", "otel"} {
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

func TestInit_RejectsInvalidFormat(t *testing.T) {
	swapDefaultLogger(t)
	err := logging.Init(logging.Config{
		Service: "liquidity-provider-server",
		Format:  "yaml",
		Level:   "info",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid log format")
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

func TestInit_RejectsUnknownLevel(t *testing.T) {
	swapDefaultLogger(t)
	err := logging.Init(logging.Config{
		Service: "liquidity-provider-server",
		Format:  "json",
		Level:   "verbose",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid log level")
}

func TestInit_ParsesLevels(t *testing.T) {
	swapDefaultLogger(t)

	for _, level := range []string{"trace", "debug", "info", "warn", "warning", "error", "fatal"} {
		t.Run(level, func(t *testing.T) {
			require.NoError(t, logging.Init(logging.Config{
				Service: "liquidity-provider-server",
				Format:  "json",
				Level:   level,
			}))
		})
	}
}

func TestInit_SetsErrorSentinelsFromConfig(t *testing.T) {
	swapDefaultLogger(t)
	setupErrorSentinels(t)

	require.NoError(t, logging.Init(logging.Config{
		Service:        "liquidity-provider-server",
		Format:         "json",
		Level:          "info",
		ErrorSentinels: usecaseSentinels(),
	}))

	var buf bytes.Buffer
	level := logging.LevelVar()
	level.Set(slog.LevelInfo)
	slog.SetDefault(slog.New(logging.NewTestHandler(&buf, level)))

	wrapped := usecases.WrapUseCaseError(usecases.AcceptPeginQuoteId, usecases.ExpiredQuoteError)
	logging.Error(context.Background(), "accept failed", wrapped)

	entry := decodeEntry(t, buf.Bytes())
	require.Equal(t, usecases.ExpiredQuoteError.Error(), entry["errorCode"])
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
	setupErrorSentinels(t)
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

func TestError_EmitsArgsInErrorContext(t *testing.T) {
	var buf bytes.Buffer
	swapDefaultLogger(t)
	setupErrorSentinels(t)
	level := logging.LevelVar()
	level.Set(slog.LevelInfo)
	slog.SetDefault(slog.New(logging.NewTestHandler(&buf, level)))

	args := usecases.ErrorArg("quoteHash", "0xabc")
	wrapped := usecases.WrapUseCaseErrorArgs(usecases.AcceptPeginQuoteId, usecases.QuoteNotFoundError, args)
	logging.Error(context.Background(), "accept failed", wrapped)

	entry := decodeEntry(t, buf.Bytes())
	errorContext, ok := entry["errorContext"].(string)
	require.True(t, ok)
	require.Equal(t, string(usecases.AcceptPeginQuoteId), errorContext[:len(usecases.AcceptPeginQuoteId)])
	require.Contains(t, errorContext, ". Args:")
	require.Contains(t, errorContext, "quoteHash")
}

func TestError_HandlesJoinedErrors(t *testing.T) {
	var buf bytes.Buffer
	swapDefaultLogger(t)
	setupErrorSentinels(t)
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

func TestError_UsesInnermostWhenNoSentinels(t *testing.T) {
	var buf bytes.Buffer
	swapDefaultLogger(t)
	level := logging.LevelVar()
	level.Set(slog.LevelInfo)
	slog.SetDefault(slog.New(logging.NewTestHandler(&buf, level)))

	root := errors.New("root cause")
	wrapped := fmt.Errorf("outer: %w", root)
	logging.Error(context.Background(), "failed", wrapped)

	entry := decodeEntry(t, buf.Bytes())
	require.Equal(t, "root cause", entry["errorCode"])
}

func TestError_OmitsShapeForNilError(t *testing.T) {
	var buf bytes.Buffer
	swapDefaultLogger(t)
	level := logging.LevelVar()
	level.Set(slog.LevelInfo)
	slog.SetDefault(slog.New(logging.NewTestHandler(&buf, level)))

	logging.Error(context.Background(), "no error", nil)

	entry := decodeEntry(t, buf.Bytes())
	_, hasCode := entry["errorCode"]
	_, hasMessage := entry["errorMessage"]
	_, hasStack := entry["errorStack"]
	_, hasContext := entry["errorContext"]
	require.False(t, hasCode)
	require.False(t, hasMessage)
	require.False(t, hasStack)
	require.False(t, hasContext)
}

func TestTraceIDFromContext(t *testing.T) {
	require.Empty(t, logging.TraceIDFromContext(context.Background()))

	ctx := logging.WithTraceID(context.Background(), "trace-123")
	require.Equal(t, "trace-123", logging.TraceIDFromContext(ctx))
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
	logging.ResetErrorSentinelsForTest()
	t.Cleanup(func() {
		slog.SetDefault(previous)
		logging.ResetLevelVar()
		logging.ResetErrorSentinelsForTest()
	})
}

func setupErrorSentinels(t *testing.T) {
	t.Helper()
	logging.SetErrorSentinelsForTest(usecaseSentinels())
}

func usecaseSentinels() []error {
	return []error{
		usecases.NonRecoverableError,
		usecases.TxBelowMinimumError,
		usecases.RskAddressNotSupportedError,
		usecases.QuoteNotFoundError,
		usecases.QuoteNotAcceptedError,
		usecases.ExpiredQuoteError,
		usecases.NoLiquidityError,
		usecases.ProviderConfigurationError,
		usecases.WrongStateError,
		usecases.NoEnoughConfirmationsError,
		usecases.InsufficientAmountError,
		usecases.RegistrationRejectedError,
		usecases.RegistrationWithdrawnError,
		usecases.IllegalQuoteStateError,
		usecases.LockingCapExceededError,
		usecases.NonPositiveWeiError,
		usecases.EmptyConfirmationsMapError,
		usecases.NonPositiveConfirmationKeyError,
		usecases.NonPositiveReimbursementWindowError,
	}
}

func decodeEntry(t *testing.T, raw []byte) map[string]any {
	t.Helper()
	var entry map[string]any
	require.NoError(t, json.Unmarshal(raw, &entry))
	return entry
}
