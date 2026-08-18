package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/logger/internal/core"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/logger/internal/handler"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/logger/internal/redact"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/logger/internal/trace"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newHandler(t *testing.T, buf *bytes.Buffer, format core.Format) slog.Handler {
	t.Helper()
	return handler.New(handler.Options{
		Output:      buf,
		Format:      format,
		Level:       core.LevelTrace,
		TimeZone:    time.UTC,
		Redactor:    redact.New(redact.Config{}),
		Service:     "svc",
		Environment: "production",
		Version:     "v1",
	})
}

func decodeJSON(t *testing.T, buf *bytes.Buffer) map[string]any {
	t.Helper()
	var m map[string]any
	require.NoError(t, json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &m))
	return m
}

func TestNewJSONEmitsBaseFieldsAndRedacts(t *testing.T) {
	var buf bytes.Buffer
	h := newHandler(t, &buf, core.FormatJSON)

	record := slog.NewRecord(time.Date(2026, 4, 14, 15, 30, 0, 123_000_000, time.UTC), slog.LevelInfo, "hello", 0)
	record.AddAttrs(slog.String("privateKey", "0xabc"), slog.String("txHash", "0x1"))
	require.NoError(t, h.Handle(context.Background(), record))

	m := decodeJSON(t, &buf)
	assert.Equal(t, "hello", m["message"])
	assert.Equal(t, "info", m["level"])
	assert.Equal(t, "svc", m["service"])
	assert.Equal(t, "production", m["environment"])
	assert.Equal(t, "v1", m["version"])
	assert.Equal(t, "0x1", m["txHash"])
	assert.Equal(t, redact.DefaultPlaceholder, m["privateKey"])
	require.Contains(t, m, "traceId")
	assert.Len(t, m["traceId"], 32)
	assert.Len(t, m["spanId"], 16)
}

func TestNewInjectsTraceFromContext(t *testing.T) {
	var buf bytes.Buffer
	h := newHandler(t, &buf, core.FormatJSON)

	tc := trace.TraceContext{
		TraceID: "4bf92f3577b34da6a3ce929d0e0e4736",
		SpanID:  "00f067aa0ba902b7",
	}
	ctx := trace.ContextWithTrace(context.Background(), tc)

	record := slog.NewRecord(time.Now().UTC(), slog.LevelInfo, "traced", 0)
	require.NoError(t, h.Handle(ctx, record))

	m := decodeJSON(t, &buf)
	assert.Equal(t, tc.TraceID, m["traceId"])
	assert.Equal(t, tc.SpanID, m["spanId"])
}

func TestNewLogfmtFormat(t *testing.T) {
	var buf bytes.Buffer
	h := newHandler(t, &buf, core.FormatLogfmt)

	record := slog.NewRecord(time.Date(2026, 4, 14, 15, 30, 0, 123_000_000, time.UTC), slog.LevelInfo, "block", 0)
	require.NoError(t, h.Handle(context.Background(), record))

	line := buf.String()
	assert.Contains(t, line, "level=info")
	assert.Contains(t, line, "service=svc")
	assert.Contains(t, line, `message=block`)
}

func TestNewUnknownFormatFallsBackToJSON(t *testing.T) {
	var buf bytes.Buffer
	h := newHandler(t, &buf, core.Format("yaml"))

	record := slog.NewRecord(time.Now().UTC(), slog.LevelInfo, "fallback", 0)
	require.NoError(t, h.Handle(context.Background(), record))

	m := decodeJSON(t, &buf)
	assert.Equal(t, "fallback", m["message"])
}

func TestWithAttrsPreservesTraceInjection(t *testing.T) {
	var buf bytes.Buffer
	h := newHandler(t, &buf, core.FormatJSON).WithAttrs([]slog.Attr{
		slog.String("operationId", "op-1"),
	})

	record := slog.NewRecord(time.Now().UTC(), slog.LevelInfo, "bound", 0)
	require.NoError(t, h.Handle(context.Background(), record))

	m := decodeJSON(t, &buf)
	assert.Equal(t, "op-1", m["operationId"])
	require.Contains(t, m, "traceId")
	assert.Len(t, m["traceId"], 32)
}

func TestWithGroupPreservesHandlerWrapper(t *testing.T) {
	var buf bytes.Buffer
	h := newHandler(t, &buf, core.FormatJSON).WithGroup("req")

	record := slog.NewRecord(time.Now().UTC(), slog.LevelInfo, "grouped", 0)
	record.AddAttrs(slog.String("method", "POST"))
	require.NoError(t, h.Handle(context.Background(), record))

	m := decodeJSON(t, &buf)
	req, ok := m["req"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "POST", req["method"])
	// Base correlation fields stay at the root even with an open group.
	require.Contains(t, m, "traceId")
	assert.Len(t, m["traceId"], 32)
	assert.NotContains(t, req, "traceId")
}

func TestWithGroupThenAttrsNestsBoundFieldsOnly(t *testing.T) {
	var buf bytes.Buffer
	h := newHandler(t, &buf, core.FormatJSON).
		WithGroup("req").
		WithAttrs([]slog.Attr{slog.String("method", "POST")})

	record := slog.NewRecord(time.Now().UTC(), slog.LevelInfo, "grouped-attrs", 0)
	require.NoError(t, h.Handle(context.Background(), record))

	m := decodeJSON(t, &buf)
	req, ok := m["req"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "POST", req["method"])
	require.Contains(t, m, "traceId")
	assert.Len(t, m["traceId"], 32)
	assert.NotContains(t, req, "traceId")
	assert.Equal(t, "svc", m["service"])
}

func TestNestedGroupRedaction(t *testing.T) {
	var buf bytes.Buffer
	h := newHandler(t, &buf, core.FormatJSON)

	record := slog.NewRecord(time.Now().UTC(), slog.LevelInfo, "nested", 0)
	record.AddAttrs(slog.Group("credentials",
		slog.String("privateKey", "0xabc"),
		slog.String("user", "alice"),
	))
	require.NoError(t, h.Handle(context.Background(), record))

	m := decodeJSON(t, &buf)
	creds, ok := m["credentials"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "alice", creds["user"])
	assert.Equal(t, redact.DefaultPlaceholder, creds["privateKey"])
}
