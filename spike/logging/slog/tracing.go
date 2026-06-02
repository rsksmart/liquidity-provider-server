// W3C Trace Context and structured HTTP logging: inbound traceparent extraction,
// middleware replacing gorilla/handlers.LoggingHandler, outbound traceparent
// injection, and integration-call logging. All log sites use slog directly.
package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"time"
)

func demoInboundTraceExtraction(url string) {
	hit(url, "")
}

func demoRequestScopedTrace(url string) {
	hit(url, "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01")
}

func demoStructuredRequestLoggingReplacement(url string) {
	hit(url, "00-11111111111111111111111111111111-2222222222222222-01")
}

// demoOutbound exercises traceTransport end-to-end and logs the four
// integration fields required by the standard via slog.
func demoOutboundTraceInjection(url string) {
	ctx := contextWithTrace(context.Background(), newTraceContext())
	client := &http.Client{Transport: traceTransport{Base: http.DefaultTransport}}

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	start := time.Now()
	resp, err := client.Do(req)
	dur := time.Since(start).Milliseconds()
	if err != nil {
		logError(ctx, "integration call failed", "OUTBOUND_FAILED", err,
			slog.String("integrationName", "rsk-lbc"),
			slog.String("integrationMethod", "getQuote"),
			slog.Int64("integrationDurationMs", dur),
		)
		return
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	slog.LogAttrs(ctx, slog.LevelInfo, "integration call",
		slog.String("integrationName", "rsk-lbc"),
		slog.String("integrationMethod", "getQuote"),
		slog.Int64("integrationDurationMs", dur),
		slog.Int("integrationStatusCode", resp.StatusCode),
	)
}

func hit(url, traceparent string) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("X-Request-Id", "req-001")
	if traceparent != "" {
		req.Header.Set("traceparent", traceparent)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("inbound demo request failed", slog.String("error", err.Error()))
		return
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
}

// TraceContext is the minimal W3C Trace Context state we carry through
// context. A real implementation would also keep traceFlags + tracestate.
type TraceContext struct {
	TraceID string // 32 hex chars
	SpanID  string // 16 hex chars
}

type traceCtxKey struct{}

// W3C traceparent: version "-" trace-id "-" parent-id "-" flags
//
//	00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01
var traceparentPattern = regexp.MustCompile(`^([0-9a-f]{2})-([0-9a-f]{32})-([0-9a-f]{16})-([0-9a-f]{2})$`)

func parseTraceparent(v string) (TraceContext, bool) {
	m := traceparentPattern.FindStringSubmatch(v)
	if m == nil {
		return TraceContext{}, false
	}
	return TraceContext{TraceID: m[2], SpanID: m[3]}, true
}

func newTraceContext() TraceContext {
	var t [16]byte
	var s [8]byte
	_, _ = rand.Read(t[:])
	_, _ = rand.Read(s[:])
	return TraceContext{TraceID: hex.EncodeToString(t[:]), SpanID: hex.EncodeToString(s[:])}
}

func contextWithTrace(ctx context.Context, tc TraceContext) context.Context {
	return context.WithValue(ctx, traceCtxKey{}, tc)
}

func traceFromContext(ctx context.Context) (TraceContext, bool) {
	if ctx == nil {
		return TraceContext{}, false
	}
	tc, ok := ctx.Value(traceCtxKey{}).(TraceContext)
	return tc, ok
}

func traceSource(adopted bool) string {
	if adopted {
		return "inbound-header"
	}
	return "minted"
}

// statusRecorder lets the middleware see the response status the
// downstream handler wrote.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

// traceMiddleware is the slog-flavoured implementation of the
// structured request-logging middleware that the SPIKE recommends as
// the replacement for gorilla/handlers.LoggingHandler. It:
//
//  1. Extracts or mints a W3C trace context.
//  2. Logs the trace decision at DEBUG via slog (contextHandler picks
//     up the ctx).
//  3. After serving, logs the six standard HTTP fields at INFO via
//     slog.InfoContext.
//
// Every log call uses slog directly so the middleware reads as a slog
// demo end-to-end.
func traceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := r.Header.Get("traceparent")
		tc, adopted := parseTraceparent(raw)
		if !adopted {
			tc = newTraceContext()
		}
		ctx := contextWithTrace(r.Context(), tc)

		slog.LogAttrs(ctx, slog.LevelDebug, "trace context resolved",
			slog.Bool("traceAdopted", adopted),
			slog.String("traceSource", traceSource(adopted)),
			slog.String("incomingTraceparent", raw),
		)

		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r.WithContext(ctx))

		slog.LogAttrs(ctx, slog.LevelInfo, "http request",
			slog.String("httpMethod", r.Method),
			slog.String("httpPath", r.URL.Path),
			slog.Int("httpStatusCode", rec.status),
			slog.Int64("durationMs", time.Since(start).Milliseconds()),
			slog.String("userAgent", r.UserAgent()),
			slog.String("requestId", r.Header.Get("X-Request-Id")),
		)
	})
}

// traceTransport is the slog-flavoured outbound RoundTripper. Every
// branch logs through slog so this file reads as a slog demo end-to-end:
//
//   - injection success: DEBUG with the parent + child span ids;
//   - no trace in ctx: WARN through slog so the operator notices the
//     missing context propagation instead of silently shipping a header.
type traceTransport struct{ Base http.RoundTripper }

func (t traceTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}
	ctx := req.Context()

	tc, ok := traceFromContext(ctx)
	if !ok {
		slog.WarnContext(ctx, "outbound call with no trace in context",
			slog.String("integrationUrl", req.URL.String()),
		)
		return base.RoundTrip(req)
	}

	fresh := newTraceContext()
	header := "00-" + tc.TraceID + "-" + fresh.SpanID + "-01"
	req = req.Clone(ctx)
	req.Header.Set("traceparent", header)

	slog.LogAttrs(ctx, slog.LevelDebug, "outbound traceparent injected",
		slog.String("parentSpanId", tc.SpanID),
		slog.String("childSpanId", fresh.SpanID),
		slog.String("integrationUrl", req.URL.String()),
	)
	return base.RoundTrip(req)
}

// sampleHandler is the inbound business handler. The six standard HTTP
// fields are emitted by traceMiddleware, so this handler only logs the
// business event. contextHandler adds traceId without explicit attrs here.
func sampleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"quoteHash":"0xabc"}`))

	slog.InfoContext(r.Context(), "quote returned",
		slog.String("quoteHash", "0xabc"),
	)
}
