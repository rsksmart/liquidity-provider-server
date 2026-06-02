// Test-output capture using zaptest/observer, the strongest API in the comparison.
// Asserts on typed observed logs without JSON parsing.
package main

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// newCaptureLogger uses zaptest/observer — the cleanest test API of the
// four candidates (see logs.md "Test capture ergonomics"). It also wraps
// the observer core in piiCore so redaction is part of the assertion.
func newCaptureLogger(t *testing.T) (*zap.Logger, *observer.ObservedLogs) {
	t.Helper()
	obsCore, recorded := observer.New(zapcore.DebugLevel)
	core := piiCore{Core: obsCore}
	l := zap.New(core).With(
		zap.String("service", "liquidity-provider-server"),
		zap.String("environment", "test"),
		zap.String("version", "test"),
	)
	logger = l
	return l, recorded
}

func TestCapture_BaseFieldsAndTraceId(t *testing.T) {
	_, recorded := newCaptureLogger(t)

	tc := TraceContext{TraceID: "4bf92f3577b34da6a3ce929d0e0e4736", SpanID: "00f067aa0ba902b7"}
	ctx := contextWithTrace(context.Background(), tc)
	loggerWithTrace(ctx).Info("ping", zap.String("requestId", "req-001"))

	all := recorded.All()
	if len(all) != 1 {
		t.Fatalf("expected 1 record, got %d", len(all))
	}
	m := all[0].ContextMap()
	for k, want := range map[string]string{
		"service":     "liquidity-provider-server",
		"environment": "test",
		"version":     "test",
		"traceId":     tc.TraceID,
		"spanId":      tc.SpanID,
		"requestId":   "req-001",
	} {
		if got, _ := m[k].(string); got != want {
			t.Errorf("field %q: got %q, want %q", k, got, want)
		}
	}
	if all[0].Message != "ping" {
		t.Errorf("message: got %q", all[0].Message)
	}
}

func TestCapture_PIIRedaction(t *testing.T) {
	_, recorded := newCaptureLogger(t)
	loggerWithTrace(context.Background()).Warn("attempt",
		zap.String("privateKey", "0xdeadbeef"),
		zap.String("apiKey", "sk-live-XYZ"),
		zap.String("safe", "ok"),
	)
	m := recorded.All()[0].ContextMap()
	for _, key := range []string{"privateKey", "apiKey"} {
		if got, _ := m[key].(string); got != "[REDACTED]" {
			t.Errorf("expected %q redacted, got %q", key, got)
		}
	}
	if got, _ := m["safe"].(string); got != "ok" {
		t.Errorf("non-PII field clobbered: got %q", got)
	}
}

func TestCapture_ErrorShape(t *testing.T) {
	_, recorded := newCaptureLogger(t)
	logError(context.Background(), "demo", "DEMO_CODE", errors.New("boom"),
		zap.String("integrationName", "rsk-lbc"))

	all := recorded.All()
	if len(all) != 1 {
		t.Fatalf("expected 1 record, got %d", len(all))
	}
	m := all[0].ContextMap()
	if got, _ := m["errorCode"].(string); got != "DEMO_CODE" {
		t.Errorf("errorCode: got %q", got)
	}
	if got, _ := m["errorMessage"].(string); got != "boom" {
		t.Errorf("errorMessage: got %q", got)
	}
	if got, _ := m["errorStack"].(string); got == "" {
		t.Errorf("errorStack: empty")
	}
	// zap.Namespace nests subsequent fields under "errorContext"; the
	// observer represents that as a map[string]any.
	ctxMap, ok := m["errorContext"].(map[string]any)
	if !ok {
		t.Fatalf("errorContext: not a map, got %T", m["errorContext"])
	}
	if got, _ := ctxMap["integrationName"].(string); got != "rsk-lbc" {
		t.Errorf("errorContext.integrationName: got %q", got)
	}
}
