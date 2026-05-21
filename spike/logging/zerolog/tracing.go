// W3C Trace Context and structured HTTP logging via zerolog's fluent event API:
// traceMiddleware, traceTransport, and integration-call logging.
package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
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

func demoOutboundTraceInjection(url string) {
	ctx := contextWithTrace(context.Background(), newTraceContext())
	client := &http.Client{Transport: traceTransport{Base: http.DefaultTransport}}

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	start := time.Now()
	resp, err := client.Do(req)
	dur := time.Since(start).Milliseconds()
	if err != nil {
		logError(ctx, "integration call failed", "OUTBOUND_FAILED", err, map[string]string{
			"integrationName":   "rsk-lbc",
			"integrationMethod": "getQuote",
		})
		return
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	loggerWithTrace(ctx).Info().
		Str("integrationName", "rsk-lbc").
		Str("integrationMethod", "getQuote").
		Int64("integrationDurationMs", dur).
		Int("integrationStatusCode", resp.StatusCode).
		Msg("integration call")
}

func hit(url, traceparent string) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("X-Request-Id", "req-001")
	if traceparent != "" {
		req.Header.Set("traceparent", traceparent)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error().Err(err).Msg("inbound demo request failed")
		return
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
}

type TraceContext struct {
	TraceID string
	SpanID  string
}

type traceCtxKey struct{}

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

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

// traceMiddleware is the zerolog-flavoured structured request
// middleware. Every log call is a fluent chain — Event builders for
// Debug() / Info() / Warn() followed by typed .Str / .Int / .Bool
// methods and a terminal .Msg(). That chain is the heart of the zerolog
// API surface and what makes the migration from logrus the most
// invasive of the four candidates: every existing call site has to be
// rewritten into this shape.
func traceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := r.Header.Get("traceparent")
		tc, adopted := parseTraceparent(raw)
		if !adopted {
			tc = newTraceContext()
		}
		ctx := contextWithTrace(r.Context(), tc)

		loggerWithTrace(ctx).Debug().
			Bool("traceAdopted", adopted).
			Str("traceSource", traceSource(adopted)).
			Str("incomingTraceparent", raw).
			Msg("trace context resolved")

		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r.WithContext(ctx))

		loggerWithTrace(ctx).Info().
			Str("httpMethod", r.Method).
			Str("httpPath", r.URL.Path).
			Int("httpStatusCode", rec.status).
			Int64("durationMs", time.Since(start).Milliseconds()).
			Str("userAgent", r.UserAgent()).
			Str("requestId", r.Header.Get("X-Request-Id")).
			Msg("http request")
	})
}

// traceTransport is the zerolog-flavoured outbound RoundTripper. Each
// branch uses zerolog's fluent event API directly.
type traceTransport struct{ Base http.RoundTripper }

func (t traceTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}
	ctx := req.Context()

	tc, ok := traceFromContext(ctx)
	if !ok {
		loggerWithTrace(ctx).Warn().
			Str("integrationUrl", req.URL.String()).
			Msg("outbound call with no trace in context")
		return base.RoundTrip(req)
	}

	fresh := newTraceContext()
	header := "00-" + tc.TraceID + "-" + fresh.SpanID + "-01"
	req = req.Clone(ctx)
	req.Header.Set("traceparent", header)

	loggerWithTrace(ctx).Debug().
		Str("parentSpanId", tc.SpanID).
		Str("childSpanId", fresh.SpanID).
		Str("integrationUrl", req.URL.String()).
		Msg("outbound traceparent injected")

	return base.RoundTrip(req)
}

// sampleHandler is the business handler. traceMiddleware emits the six
// standard HTTP fields; here we only emit the business log via the
// zerolog fluent chain with trace fields added by loggerWithTrace.
func sampleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"quoteHash":"0xabc"}`))

	loggerWithTrace(r.Context()).Info().
		Str("quoteHash", "0xabc").
		Msg("quote returned")
}
