// W3C Trace Context and structured HTTP logging via baseEntry(ctx) at each log
// site: traceMiddleware, traceTransport, and integration-call logging.
package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/sirupsen/logrus"
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
		logError(ctx, "integration call failed", "OUTBOUND_FAILED", err, logrus.Fields{
			"integrationName":       "rsk-lbc",
			"integrationMethod":     "getQuote",
			"integrationDurationMs": dur,
		})
		return
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	baseEntry(ctx).WithFields(logrus.Fields{
		"integrationName":       "rsk-lbc",
		"integrationMethod":     "getQuote",
		"integrationDurationMs": dur,
		"integrationStatusCode": resp.StatusCode,
	}).Info("integration call")
}

func hit(url, traceparent string) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("X-Request-Id", "req-001")
	if traceparent != "" {
		req.Header.Set("traceparent", traceparent)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		baseEntry(context.Background()).WithError(err).Error("inbound demo request failed")
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

// traceMiddleware is the logrus-flavoured structured request middleware.
// Every log call goes through baseEntry(ctx).WithFields(...).<Level>(msg)
// — there is no auto-attachment, so the function reads as a sequence of
// explicit logrus.Entry mutations. That verbosity *is* the logrus
// ergonomics story this demo is showing.
func traceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := r.Header.Get("traceparent")
		tc, adopted := parseTraceparent(raw)
		if !adopted {
			tc = newTraceContext()
		}
		ctx := contextWithTrace(r.Context(), tc)

		baseEntry(ctx).WithFields(logrus.Fields{
			"traceAdopted":        adopted,
			"traceSource":         traceSource(adopted),
			"incomingTraceparent": raw,
		}).Debug("trace context resolved")

		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r.WithContext(ctx))

		baseEntry(ctx).WithFields(logrus.Fields{
			"httpMethod":     r.Method,
			"httpPath":       r.URL.Path,
			"httpStatusCode": rec.status,
			"durationMs":     time.Since(start).Milliseconds(),
			"userAgent":      r.UserAgent(),
			"requestId":      r.Header.Get("X-Request-Id"),
		}).Info("http request")
	})
}

// traceTransport is the logrus-flavoured outbound RoundTripper. Note
// the per-call cost: every branch has to call baseEntry(ctx) to attach
// the trace fields, then .WithFields(...), then .<Level>(msg). That is
// the same verbosity slog avoids structurally via its contextHandler.
type traceTransport struct{ Base http.RoundTripper }

func (t traceTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}
	ctx := req.Context()

	tc, ok := traceFromContext(ctx)
	if !ok {
		baseEntry(ctx).WithField("integrationUrl", req.URL.String()).
			Warn("outbound call with no trace in context")
		return base.RoundTrip(req)
	}

	fresh := newTraceContext()
	header := "00-" + tc.TraceID + "-" + fresh.SpanID + "-01"
	req = req.Clone(ctx)
	req.Header.Set("traceparent", header)

	baseEntry(ctx).WithFields(logrus.Fields{
		"parentSpanId":   tc.SpanID,
		"childSpanId":    fresh.SpanID,
		"integrationUrl": req.URL.String(),
	}).Debug("outbound traceparent injected")

	return base.RoundTrip(req)
}

// sampleHandler is the business handler. traceMiddleware now owns the
// six standard HTTP fields, so this handler only emits the business log
// via baseEntry(ctx), because logrus has no native context-to-fields bridge.
func sampleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"quoteHash":"0xabc"}`))

	baseEntry(r.Context()).WithField("quoteHash", "0xabc").Info("quote returned")
}
