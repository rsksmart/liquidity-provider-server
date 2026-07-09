package logging_test

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/logging"
	"github.com/stretchr/testify/require"
)

func TestHandler_RedactsTopLevelPII(t *testing.T) {
	var buf bytes.Buffer
	level := logging.LevelVar()
	level.Set(slog.LevelInfo)
	logger := slog.New(logging.NewTestHandler(&buf, level))

	logger.Info("test", "password", "secret-value", "quoteHash", "0xabc")

	entry := decodeEntry(t, buf.Bytes())
	require.Equal(t, "[REDACTED]", entry["password"])
	require.Equal(t, "0xabc", entry["quoteHash"])
}

func TestHandler_RedactsPIIInsideGroups(t *testing.T) {
	var buf bytes.Buffer
	level := logging.LevelVar()
	level.Set(slog.LevelInfo)
	logger := slog.New(logging.NewTestHandler(&buf, level))

	logger.Info("test", slog.Group("request", slog.String("apiKey", "key-123"), slog.String("path", "/quotes")))

	entry := decodeEntry(t, buf.Bytes())
	request, ok := entry["request"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "[REDACTED]", request["apiKey"])
	require.Equal(t, "/quotes", request["path"])
}

func TestHandler_AddsTraceIDFromContext(t *testing.T) {
	var buf bytes.Buffer
	level := logging.LevelVar()
	level.Set(slog.LevelInfo)
	logger := slog.New(logging.NewTestHandler(&buf, level))

	ctx := logging.WithTraceID(context.Background(), "trace-abc")
	logger.InfoContext(ctx, "test message")

	entry := decodeEntry(t, buf.Bytes())
	require.Equal(t, "trace-abc", entry["traceId"])
}

func TestHandler_OmitsTraceIDWhenAbsent(t *testing.T) {
	var buf bytes.Buffer
	level := logging.LevelVar()
	level.Set(slog.LevelInfo)
	logger := slog.New(logging.NewTestHandler(&buf, level))

	logger.InfoContext(context.Background(), "test message")

	entry := decodeEntry(t, buf.Bytes())
	_, ok := entry["traceId"]
	require.False(t, ok)
}
