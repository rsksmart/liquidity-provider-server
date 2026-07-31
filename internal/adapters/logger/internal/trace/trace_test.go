package trace_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/pkg/logger/internal/trace"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTraceparentValid(t *testing.T) {
	tc, err := trace.ParseTraceparent("00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01")

	require.NoError(t, err)
	assert.Equal(t, "00", tc.Version)
	assert.Equal(t, "4bf92f3577b34da6a3ce929d0e0e4736", tc.TraceID)
	assert.Equal(t, "00f067aa0ba902b7", tc.ParentID)
	assert.Equal(t, "01", tc.Flags)
}

func TestParseTraceparentInvalid(t *testing.T) {
	cases := []string{
		"",
		"garbage",
		"ff-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
		"00-tooshort-00f067aa0ba902b7-01",
	}
	for _, header := range cases {
		_, err := trace.ParseTraceparent(header)
		assert.ErrorIs(t, err, trace.ErrInvalidTraceparent, header)
	}
}

func TestNewTraceContextAndTraceparent(t *testing.T) {
	tc, err := trace.NewTraceContext()
	require.NoError(t, err)

	assert.Equal(t, "00", tc.Version)
	assert.Len(t, tc.TraceID, 32)
	assert.Len(t, tc.SpanID, 16)
	assert.Equal(t, "01", tc.Flags)

	parsed, err := trace.ParseTraceparent(tc.Traceparent())
	require.NoError(t, err)
	assert.Equal(t, tc.TraceID, parsed.TraceID)
	assert.Equal(t, tc.SpanID, parsed.ParentID)
}

func TestTraceparentDefaultsEmptyVersionAndFlags(t *testing.T) {
	tc := trace.TraceContext{
		TraceID: "4bf92f3577b34da6a3ce929d0e0e4736",
		SpanID:  "00f067aa0ba902b7",
	}
	assert.Equal(t, "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01", tc.Traceparent())
}

func TestWithNewSpan(t *testing.T) {
	root, err := trace.NewTraceContext()
	require.NoError(t, err)

	child, err := root.WithNewSpan()
	require.NoError(t, err)
	assert.Equal(t, root.TraceID, child.TraceID)
	assert.Equal(t, root.SpanID, child.ParentID)
	assert.NotEqual(t, root.SpanID, child.SpanID)
}

func TestWithNewSpanDefaultsEmptyVersionAndFlags(t *testing.T) {
	parent := trace.TraceContext{
		TraceID: "4bf92f3577b34da6a3ce929d0e0e4736",
		SpanID:  "00f067aa0ba902b7",
	}
	child, err := parent.WithNewSpan()
	require.NoError(t, err)
	assert.Equal(t, "00", child.Version)
	assert.Equal(t, "01", child.Flags)
}

func TestNewSpanID(t *testing.T) {
	id, err := trace.NewSpanID()
	require.NoError(t, err)
	assert.Len(t, id, 16)
}

func TestContextRoundTrip(t *testing.T) {
	tc := trace.TraceContext{TraceID: "abc", SpanID: "def"}
	ctx := trace.ContextWithTrace(context.Background(), tc)

	got, ok := trace.FromContext(ctx)
	require.True(t, ok)
	assert.Equal(t, tc, got)

	_, ok = trace.FromContext(context.Background())
	assert.False(t, ok)

	_, ok = trace.FromContext(nil) //nolint:staticcheck // intentional nil context
	assert.False(t, ok)
}

func TestMiddlewareContinuesInboundTrace(t *testing.T) {
	var captured trace.TraceContext
	var capturedOK bool
	handler := trace.Middleware(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		captured, capturedOK = trace.FromContext(r.Context())
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("traceparent", "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.True(t, capturedOK)
	assert.Equal(t, "4bf92f3577b34da6a3ce929d0e0e4736", captured.TraceID)
	assert.Equal(t, "00f067aa0ba902b7", captured.ParentID)
	assert.NotEqual(t, "00f067aa0ba902b7", captured.SpanID)
	assert.Equal(t, captured.Traceparent(), rec.Header().Get("traceparent"))
}

func TestMiddlewareStartsNewTraceWhenAbsent(t *testing.T) {
	var captured trace.TraceContext
	handler := trace.Middleware(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		captured, _ = trace.FromContext(r.Context())
	}))

	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))

	assert.Len(t, captured.TraceID, 32)
	assert.Len(t, captured.SpanID, 16)
	assert.Empty(t, captured.ParentID)
}
