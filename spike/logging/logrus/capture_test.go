// Test-output capture via logrus.SetOutput to a buffer, then JSON-line parsing
// and field assertions. Matches current LPS test patterns.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

// newCaptureLogger swaps the package-level logger for one that writes to
// a buffer, matching the SetOutput pattern that the current LPS tests
// already use (logs.md "Test capture ergonomics").
func newCaptureLogger(buf *bytes.Buffer) *logrus.Logger {
	cfg = Config{
		Format:      "json",
		Level:       logrus.DebugLevel,
		Service:     "liquidity-provider-server",
		Environment: "test",
		Version:     "test",
	}
	l := logrus.New()
	l.SetOutput(buf)
	l.SetLevel(cfg.Level)
	l.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "timestamp",
			logrus.FieldKeyMsg:  "message",
		},
	})
	l.AddHook(&piiRedactHook{})
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
	baseEntry(ctx).WithField("requestId", "req-001").Info("ping")

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

	baseEntry(context.Background()).WithFields(logrus.Fields{
		"privateKey": "0xdeadbeef",
		"apiKey":     "sk-live-XYZ",
		"safe":       "ok",
	}).Warn("attempt")

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

	logError(context.Background(), "demo", "DEMO_CODE", errors.New("boom"), logrus.Fields{
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
