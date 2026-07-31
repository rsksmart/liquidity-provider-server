package trace

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Identifier sizes in bytes; hex-encoding doubles them (16 bytes -> 32 hex
// characters for a trace id, 8 bytes -> 16 hex characters for a span id).
const (
	traceIDBytes = 16
	spanIDBytes  = 8
)

const (
	// defaultVersion is the W3C Trace Context version emitted for new traces.
	defaultVersion = "00"
	// defaultTraceFlags marks a trace as sampled per W3C Trace Context.
	defaultTraceFlags = "01"
	// invalidVersion is the W3C Trace Context version that is forbidden.
	invalidVersion = "ff"
)

// traceparentPattern validates the W3C traceparent header:
// {version}-{trace-id}-{parent-id}-{trace-flags}.
var traceparentPattern = regexp.MustCompile(`^[0-9a-f]{2}-[0-9a-f]{32}-[0-9a-f]{16}-[0-9a-f]{2}$`)

// ErrInvalidTraceparent indicates that a W3C traceparent header is malformed or
// uses the forbidden ff version.
var ErrInvalidTraceparent = errors.New("invalid traceparent")

// TraceContext holds the W3C Trace Context identifiers for the current
// operation. TraceID correlates every log line of a request across services;
// SpanID identifies this service's span; ParentID links to the caller's span.
type TraceContext struct {
	// Version is the two-character hex W3C Trace Context version.
	Version string
	// TraceID is a 32-character hex string, stable across the whole trace.
	TraceID string
	// SpanID is a 16-character hex string identifying the current span.
	SpanID string
	// ParentID is the 16-character hex span id of the caller, empty at the
	// root of a trace.
	ParentID string
	// Flags is the two-character hex sampling flag. Defaults to "01".
	Flags string
}

// NewTraceContext generates a brand new trace with a fresh trace id and span
// id. Use it at the entry point of a new trace (e.g. a request without an
// inbound traceparent, or an externally triggered job).
func NewTraceContext() (TraceContext, error) {
	traceID, err := randomHex(traceIDBytes)
	if err != nil {
		return TraceContext{}, fmt.Errorf("generate trace ID: %w", err)
	}
	spanID, err := randomHex(spanIDBytes)
	if err != nil {
		return TraceContext{}, fmt.Errorf("generate span ID: %w", err)
	}
	return TraceContext{
		Version: defaultVersion,
		TraceID: traceID,
		SpanID:  spanID,
		Flags:   defaultTraceFlags,
	}, nil
}

// WithNewSpan returns a copy of the trace context for a new span within the
// same trace: the current span becomes the parent and a new span id is
// generated. Use it when this service starts a new unit of work under an
// inbound trace.
func (tc TraceContext) WithNewSpan() (TraceContext, error) {
	spanID, err := randomHex(spanIDBytes)
	if err != nil {
		return TraceContext{}, fmt.Errorf("generate span ID: %w", err)
	}
	tc.ParentID = tc.SpanID
	tc.SpanID = spanID
	if tc.Version == "" {
		tc.Version = defaultVersion
	}
	if tc.Flags == "" {
		tc.Flags = defaultTraceFlags
	}
	return tc, nil
}

// Traceparent renders the trace context as a W3C traceparent header value.
func (tc TraceContext) Traceparent() string {
	version := tc.Version
	if version == "" {
		version = defaultVersion
	}
	flags := tc.Flags
	if flags == "" {
		flags = defaultTraceFlags
	}
	return version + "-" + tc.TraceID + "-" + tc.SpanID + "-" + flags
}

// ParseTraceparent parses a W3C traceparent header. It returns
// ErrInvalidTraceparent when the header is malformed. On success the incoming
// span id is stored as ParentID (it is the caller's span).
func ParseTraceparent(header string) (TraceContext, error) {
	if !traceparentPattern.MatchString(header) {
		return TraceContext{}, ErrInvalidTraceparent
	}
	parts := strings.Split(header, "-")
	if parts[0] == invalidVersion {
		return TraceContext{}, ErrInvalidTraceparent
	}
	if parts[1] == "00000000000000000000000000000000" || parts[2] == "0000000000000000" {
		return TraceContext{}, ErrInvalidTraceparent
	}
	return TraceContext{
		Version:  parts[0],
		TraceID:  parts[1],
		ParentID: parts[2],
		Flags:    parts[3],
	}, nil
}

// NewSpanID returns a fresh 16-character hex span id. It is used by callers
// that establish this service's span under an inbound parent.
func NewSpanID() (string, error) {
	return randomHex(spanIDBytes)
}

// randomHex returns n cryptographically-random bytes as a hex string of length 2n.
func randomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("read crypto randomness: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// traceContextKey is the private context key type for the trace context.
type traceContextKey struct{}

// ContextWithTrace returns a copy of ctx carrying the trace context. The logger
// reads it automatically so every log line within the operation includes the
// traceId.
func ContextWithTrace(ctx context.Context, tc TraceContext) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, traceContextKey{}, tc)
}

// FromContext extracts the trace context stored by ContextWithTrace. The bool
// reports whether one was present.
func FromContext(ctx context.Context) (TraceContext, bool) {
	if ctx == nil {
		return TraceContext{}, false
	}
	tc, ok := ctx.Value(traceContextKey{}).(TraceContext)
	return tc, ok
}
