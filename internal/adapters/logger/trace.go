package logger

import (
	"context"
	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/logger/internal/trace"
)

// TraceContext holds the W3C Trace Context identifiers (trace id, span id,
// parent id, version, flags) for the current operation. The logger reads it
// from the context so every line carries the traceId.
type TraceContext = trace.TraceContext

// ErrInvalidTraceparent indicates that a W3C traceparent header is malformed.
var ErrInvalidTraceparent = trace.ErrInvalidTraceparent

// NewTraceContext generates a brand new trace with a fresh trace id and span
// id. Use it at the entry point of a new trace.
func NewTraceContext() (TraceContext, error) { return trace.NewTraceContext() }

// ParseTraceparent parses a W3C traceparent header. It returns
// ErrInvalidTraceparent when the header is malformed.
func ParseTraceparent(header string) (TraceContext, error) {
	return trace.ParseTraceparent(header)
}

// ContextWithTrace returns a copy of ctx carrying the trace context so the
// logger includes the traceId on every line within the operation.
func ContextWithTrace(ctx context.Context, tc TraceContext) context.Context {
	return trace.ContextWithTrace(ctx, tc)
}

// TraceFromContext extracts the trace context stored by ContextWithTrace. The
// bool reports whether one was present.
func TraceFromContext(ctx context.Context) (TraceContext, bool) {
	return trace.FromContext(ctx)
}

// TraceMiddleware is a net/http middleware that establishes a trace context for
// each inbound request (continuing an inbound traceparent or starting a fresh
// trace), stores it so the logger picks up the traceId, and echoes the
// resulting traceparent on the response.
func TraceMiddleware(next http.Handler) http.Handler {
	return trace.Middleware(next)
}
