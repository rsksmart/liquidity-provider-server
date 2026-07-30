package trace

import "net/http"

// traceparentHeader is the W3C Trace Context HTTP header name.
const traceparentHeader = "traceparent"

// Middleware is a net/http middleware that establishes a trace context for each
// inbound request. It extracts the traceparent header when present (starting a
// new child span for this service under the inbound parent) or generates a
// fresh trace when absent, stores the context so the logger picks up the
// traceId automatically, and echoes the resulting traceparent on the response.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceCtx, err := contextForRequest(r.Header.Get(traceparentHeader))
		if err != nil {
			http.Error(w, "failed to establish trace context", http.StatusInternalServerError)
			return
		}

		w.Header().Set(traceparentHeader, traceCtx.Traceparent())
		next.ServeHTTP(w, r.WithContext(ContextWithTrace(r.Context(), traceCtx)))
	})
}

// contextForRequest continues a valid inbound trace or starts a new one. An
// invalid or absent traceparent is treated as the start of a new trace.
func contextForRequest(header string) (TraceContext, error) {
	traceCtx, err := ParseTraceparent(header)
	if err != nil {
		return NewTraceContext()
	}

	spanID, err := NewSpanID()
	if err != nil {
		return TraceContext{}, err
	}
	// The inbound parent-id is the caller's span; this service generates its
	// own span under it.
	traceCtx.SpanID = spanID
	return traceCtx, nil
}
