// Test-output capture via a buffer-backed logger and JSON-line parsing, mirroring
// the slog demo's buffer approach.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

// newCaptureLogger wires a logger that writes JSON through piiWriter
// into a buffer. zerolog's "test capture story" is essentially this —
// SetOutput a buffer, then JSON-parse lines. Workable, but less
// expressive than zap's observer.
func newCaptureLogger(buf *bytes.Buffer) zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.MessageFieldName = "message"
	zerolog.TimestampFieldName = "timestamp"
	zerolog.LevelFieldName = "level"

	l := zerolog.New(piiWriter{Out: buf}).Level(zerolog.DebugLevel).With().
		Timestamp().
		Str("service", "liquidity-provider-server").
		Str("environment", "test").
		Str("version", "test").
		Logger()
	logger = l
	return l
}

func records(t *testing.T, buf *bytes.Buffer) []map[string]any {
	t.Helper()
	out := []map[string]any{}
	for _, line := range strings.Split(buf.String(), "\n") {
		if line == "" {
			continue
		}
		var rec map[string]any
		if err := json.Unmarshal([]byte(line), &rec); err != nil {
			t.Fatalf("unmarshal %q: %v", line, err)
		}
		out = append(out, rec)
	}
	return out
}

func TestCapture_BaseFieldsAndTraceId(t *testing.T) {
	buf := &bytes.Buffer{}
	newCaptureLogger(buf)

	tc := TraceContext{TraceID: "4bf92f3577b34da6a3ce929d0e0e4736", SpanID: "00f067aa0ba902b7"}
	ctx := contextWithTrace(context.Background(), tc)
	loggerWithTrace(ctx).Info().Str("requestId", "req-001").Msg("ping")

	recs := records(t, buf)
	if len(recs) != 1 {
		t.Fatalf("expected 1 record, got %d", len(recs))
	}
	r := recs[0]
	for k, want := range map[string]string{
		"service":     "liquidity-provider-server",
		"environment": "test",
		"version":     "test",
		"traceId":     tc.TraceID,
		"spanId":      tc.SpanID,
		"message":     "ping",
		"requestId":   "req-001",
	} {
		if got, _ := r[k].(string); got != want {
			t.Errorf("field %q: got %q, want %q", k, got, want)
		}
	}
}

func TestCapture_PIIRedaction(t *testing.T) {
	buf := &bytes.Buffer{}
	newCaptureLogger(buf)
	loggerWithTrace(context.Background()).Warn().
		Str("privateKey", "0xdeadbeef").
		Str("apiKey", "sk-live-XYZ").
		Str("safe", "ok").
		Msg("attempt")

	recs := records(t, buf)
	if len(recs) != 1 {
		t.Fatalf("expected 1 record, got %d", len(recs))
	}
	r := recs[0]
	for _, key := range []string{"privateKey", "apiKey"} {
		if got, _ := r[key].(string); got != "[REDACTED]" {
			t.Errorf("expected %q redacted, got %q", key, got)
		}
	}
	if got, _ := r["safe"].(string); got != "ok" {
		t.Errorf("non-PII field clobbered: got %q", got)
	}
}

func TestCapture_ErrorShape(t *testing.T) {
	buf := &bytes.Buffer{}
	newCaptureLogger(buf)
	logError(context.Background(), "demo", "DEMO_CODE", errors.New("boom"), map[string]string{
		"integrationName": "rsk-lbc",
	})

	recs := records(t, buf)
	if len(recs) != 1 {
		t.Fatalf("expected 1 record, got %d", len(recs))
	}
	r := recs[0]
	if got, _ := r["errorCode"].(string); got != "DEMO_CODE" {
		t.Errorf("errorCode: got %q", got)
	}
	if got, _ := r["errorMessage"].(string); got != "boom" {
		t.Errorf("errorMessage: got %q", got)
	}
	if got, _ := r["errorStack"].(string); got == "" {
		t.Errorf("errorStack: empty")
	}
	if _, ok := r["errorContext"].(map[string]any); !ok {
		t.Errorf("errorContext: not a map, got %T", r["errorContext"])
	}
}
