// Test-output capture: buffer-backed JSON handler mirroring initLogger's stack,
// with field assertions. Replaces logrus SetOutput/SetLevel patterns in tests.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

// newCaptureLogger returns a logger that writes JSON records to buf,
// wired with the same handler stack (PII redactor + context handler)
// that initLogger installs. This is the slog "test capture story" —
// see logs.md "Test capture ergonomics".
func newCaptureLogger(buf *bytes.Buffer) *slog.Logger {
	base := slog.NewJSONHandler(buf, &slog.HandlerOptions{
		Level:       slog.LevelDebug,
		ReplaceAttr: replaceAttr,
	})
	return slog.New(contextHandler{Handler: base}).With(
		slog.String("service", "liquidity-provider-server"),
		slog.String("environment", "test"),
		slog.String("version", "test"),
	)
}

// records parses each non-empty line of the buffer as one slog record.
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
	logger := newCaptureLogger(buf)

	tc := TraceContext{TraceID: "4bf92f3577b34da6a3ce929d0e0e4736", SpanID: "00f067aa0ba902b7"}
	ctx := contextWithTrace(context.Background(), tc)
	logger.InfoContext(ctx, "ping", slog.String("requestId", "req-001"))

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
	logger := newCaptureLogger(buf)

	logger.Warn("attempt",
		slog.String("privateKey", "0xdeadbeef"),
		slog.String("apiKey", "sk-live-XYZ"),
		slog.String("safe", "ok"),
	)

	recs := records(t, buf)
	if len(recs) != 1 {
		t.Fatalf("expected 1 record, got %d", len(recs))
	}
	r := recs[0]

	for _, key := range []string{"privateKey", "apiKey"} {
		if got, _ := r[key].(string); got != "[REDACTED]" {
			t.Errorf("expected %q to be redacted, got %q", key, got)
		}
	}
	if got, _ := r["safe"].(string); got != "ok" {
		t.Errorf("non-PII field clobbered: got %q", got)
	}
}

func TestCapture_ErrorShape(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := newCaptureLogger(buf)
	slog.SetDefault(logger) // logError uses slog.LogAttrs on the default logger.

	logError(context.Background(), "demo", "DEMO_CODE", errPlaceholder("boom"),
		slog.String("integrationName", "rsk-lbc"))

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
		t.Errorf("errorContext: not a group, got %T", r["errorContext"])
	}
}
