package logger_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTraceparentValid(t *testing.T) {
	tc, err := logger.ParseTraceparent("00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01")

	require.NoError(t, err)
	assert.Equal(t, "00", tc.Version)
	assert.Equal(t, "4bf92f3577b34da6a3ce929d0e0e4736", tc.TraceID)
	assert.Equal(t, "00f067aa0ba902b7", tc.ParentID)
	assert.Equal(t, "01", tc.Flags)
}

func TestParseTraceparentPreservesVersion(t *testing.T) {
	tc, err := logger.ParseTraceparent("01-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01")

	require.NoError(t, err)
	assert.Equal(t, "01", tc.Version)
	tc.SpanID = tc.ParentID
	assert.Equal(t, "01-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01", tc.Traceparent())
}

func TestParseTraceparentInvalid(t *testing.T) {
	cases := []string{
		"",
		"garbage",
		"ff-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
		"00-tooshort-00f067aa0ba902b7-01",
	}
	for _, header := range cases {
		_, err := logger.ParseTraceparent(header)
		assert.ErrorIs(t, err, logger.ErrInvalidTraceparent, header)
	}
}

func TestTraceparentRoundTrip(t *testing.T) {
	tc, err := logger.NewTraceContext()
	require.NoError(t, err)
	assert.Equal(t, "00", tc.Version)
	assert.Len(t, tc.TraceID, 32)
	assert.Len(t, tc.SpanID, 16)

	parsed, err := logger.ParseTraceparent(tc.Traceparent())
	require.NoError(t, err)
	assert.Equal(t, tc.Version, parsed.Version)
	assert.Equal(t, tc.TraceID, parsed.TraceID)
	assert.Equal(t, tc.SpanID, parsed.ParentID)
}

func TestWithNewSpanKeepsTraceIDAndPromotesParent(t *testing.T) {
	root, err := logger.NewTraceContext()
	require.NoError(t, err)

	child, err := root.WithNewSpan()
	require.NoError(t, err)

	assert.Equal(t, root.Version, child.Version)
	assert.Equal(t, root.TraceID, child.TraceID)
	assert.Equal(t, root.SpanID, child.ParentID)
	assert.NotEqual(t, root.SpanID, child.SpanID)
}

func TestWithNewSpanDefaultsEmptyVersionAndFlags(t *testing.T) {
	parent := logger.TraceContext{
		TraceID: "4bf92f3577b34da6a3ce929d0e0e4736",
		SpanID:  "00f067aa0ba902b7",
	}

	child, err := parent.WithNewSpan()
	require.NoError(t, err)

	assert.Equal(t, "00", child.Version)
	assert.Equal(t, "01", child.Flags)
	assert.Equal(t, parent.SpanID, child.ParentID)
}

func TestTraceparentDefaultsEmptyVersionAndFlags(t *testing.T) {
	tc := logger.TraceContext{
		TraceID: "4bf92f3577b34da6a3ce929d0e0e4736",
		SpanID:  "00f067aa0ba902b7",
	}

	assert.Equal(t, "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01", tc.Traceparent())
}

func TestTraceFromContextNil(t *testing.T) {
	_, ok := logger.TraceFromContext(nil) //nolint:staticcheck // intentional nil context
	assert.False(t, ok)
}

func TestTraceMiddlewareContinuesInboundTrace(t *testing.T) {
	var captured logger.TraceContext
	var capturedOK bool
	handler := logger.TraceMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		captured, capturedOK = logger.TraceFromContext(r.Context())
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("traceparent", "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.True(t, capturedOK)
	assert.Equal(t, "00", captured.Version)
	assert.Equal(t, "4bf92f3577b34da6a3ce929d0e0e4736", captured.TraceID)
	assert.Equal(t, "00f067aa0ba902b7", captured.ParentID)
	assert.NotEqual(t, "00f067aa0ba902b7", captured.SpanID)
	assert.Equal(t, captured.Traceparent(), rec.Header().Get("traceparent"))
}

func TestTraceMiddlewareStartsNewTraceWhenAbsent(t *testing.T) {
	var captured logger.TraceContext
	handler := logger.TraceMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		captured, _ = logger.TraceFromContext(r.Context())
	}))

	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))

	assert.Len(t, captured.TraceID, 32)
	assert.Len(t, captured.SpanID, 16)
	assert.Empty(t, captured.ParentID)
}

func TestTraceFromContextEmpty(t *testing.T) {
	_, ok := logger.TraceFromContext(context.Background())
	assert.False(t, ok)
}

func TestInvalidTraceparentSentinel(t *testing.T) {
	_, err := logger.ParseTraceparent("invalid")
	require.Error(t, err)
	assert.ErrorIs(t, err, logger.ErrInvalidTraceparent)
}
