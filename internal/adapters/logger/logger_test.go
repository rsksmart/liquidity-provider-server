package logger_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	sdklog "go.opentelemetry.io/otel/sdk/log"

	"github.com/rsksmart/liquidity-provider-server/pkg/logcheck"
	"github.com/rsksmart/liquidity-provider-server/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type codedError struct {
	msg  string
	code string
}

func (e codedError) Error() string { return e.msg }
func (e codedError) Code() string  { return e.code }

// failingWriter always returns an error so handler.Handle failures can be
// exercised without panicking the logger.
type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) {
	return 0, errors.New("write failed")
}

// fixedTime is the deterministic timestamp used across logger tests. Its
// nanosecond component is 123ms so the millisecond formatting is exercised.
var fixedTime = time.Date(2026, time.April, 14, 15, 30, 0, 123_000_000, time.UTC)

func fixedClock() func() time.Time {
	return func() time.Time { return fixedTime }
}

func newTestLogger(t *testing.T, buf *bytes.Buffer, format logger.Format) *logger.Logger {
	t.Helper()
	cfg := logger.Config{
		Service:     "lps",
		Environment: "production",
		Version:     "v2.5.2",
		Level:       logger.LevelTrace,
		Format:      format,
		Output:      buf,
	}.WithClock(fixedClock())
	log, err := logger.New(cfg)
	require.NoError(t, err)
	return log
}

// decodeLine unmarshals the single JSON record written to buf.
func decodeLine(t *testing.T, buf *bytes.Buffer) map[string]any {
	t.Helper()
	line := strings.TrimSpace(buf.String())
	var m map[string]any
	require.NoError(t, json.Unmarshal([]byte(line), &m))
	return m
}

func TestJSONHandlerWritesBaseFields(t *testing.T) {
	var buf bytes.Buffer
	log := newTestLogger(t, &buf, logger.FormatJSON)

	log.Info(context.Background(), "Bridge transaction initiated", logger.String("txHash", "0x1234abcd"))

	m := decodeLine(t, &buf)
	assert.Equal(t, "2026-04-14T15:30:00.123Z", m["timestamp"])
	assert.Equal(t, "info", m["level"])
	assert.Equal(t, "rsk-bridge-api", m["service"])
	assert.Equal(t, "production", m["environment"])
	assert.Equal(t, "v1.4.2", m["version"])
	assert.Equal(t, "Bridge transaction initiated", m["message"])
	assert.Equal(t, "0x1234abcd", m["txHash"])
	assert.Contains(t, m, "traceId")
}

func TestJSONHandlerLowercasesCustomLevels(t *testing.T) {
	var buf bytes.Buffer
	log := newTestLogger(t, &buf, logger.FormatJSON)

	log.Trace(context.Background(), "trace line")

	assert.Equal(t, "trace", decodeLine(t, &buf)["level"])
}

func TestLevelBelowThresholdIsDropped(t *testing.T) {
	var buf bytes.Buffer
	cfg := logger.Config{
		Service:     "svc",
		Environment: "production",
		Version:     "v1",
		Level:       logger.LevelInfo,
		Output:      &buf,
	}.WithClock(fixedClock())
	log, err := logger.New(cfg)
	require.NoError(t, err)

	log.Debug(context.Background(), "should be filtered")

	assert.Empty(t, buf.String())
}

func TestWithBindsBusinessFields(t *testing.T) {
	var buf bytes.Buffer
	log := newTestLogger(t, &buf, logger.FormatJSON).With(logger.String("operationId", "op-42"))

	log.Info(context.Background(), "processing")

	assert.Equal(t, "op-42", decodeLine(t, &buf)["operationId"])
}

func TestTraceContextIsInjected(t *testing.T) {
	var buf bytes.Buffer
	log := newTestLogger(t, &buf, logger.FormatJSON)

	tc := logger.TraceContext{TraceID: "4bf92f3577b34da6a3ce929d0e0e4736", SpanID: "00f067aa0ba902b7", Flags: "01"}
	ctx := logger.ContextWithTrace(context.Background(), tc)

	log.Info(ctx, "with trace")

	m := decodeLine(t, &buf)
	assert.Equal(t, "4bf92f3577b34da6a3ce929d0e0e4736", m["traceId"])
	assert.Equal(t, "00f067aa0ba902b7", m["spanId"])
}

func TestTraceIDIsAlwaysPresentEvenWithoutContext(t *testing.T) {
	var buf bytes.Buffer
	log := newTestLogger(t, &buf, logger.FormatJSON)

	log.Info(context.Background(), "no trace")

	m := decodeLine(t, &buf)
	require.Contains(t, m, "traceId")
	assert.Empty(t, m["traceId"])
}

func TestErrorLevelAttachesStack(t *testing.T) {
	var buf bytes.Buffer
	log := newTestLogger(t, &buf, logger.FormatJSON)

	log.Error(context.Background(), "boom", logger.Err(context.DeadlineExceeded))

	m := decodeLine(t, &buf)
	assert.Equal(t, context.DeadlineExceeded.Error(), m["errorMessage"])
	require.Contains(t, m, "errorStack")
	assert.Contains(t, m["errorStack"], "logger_test.TestErrorLevelAttachesStack")
}

func TestErrNestsErrorContext(t *testing.T) {
	var buf bytes.Buffer
	log := newTestLogger(t, &buf, logger.FormatJSON)

	log.Error(context.Background(), "db write failed", logger.Err(context.DeadlineExceeded,
		logger.String("quoteId", "q-1"),
		logger.Int("retries", 3),
	))

	m := decodeLine(t, &buf)
	assert.Equal(t, context.DeadlineExceeded.Error(), m["errorMessage"])
	assert.NotContains(t, m, "quoteId")
	assert.NotContains(t, m, "retries")
	ctx, ok := m["errorContext"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "q-1", ctx["quoteId"])
	assert.EqualValues(t, 3, ctx["retries"])
}

func TestErrNilIsDiscarded(t *testing.T) {
	var buf bytes.Buffer
	log := newTestLogger(t, &buf, logger.FormatJSON)

	log.Error(context.Background(), "boom", logger.Err(nil))

	m := decodeLine(t, &buf)
	assert.NotContains(t, m, "errorMessage")
	assert.NotContains(t, m, "errorCode")
	assert.NotContains(t, m, "errorContext")
}

func TestErrEmitsErrorCodeFromCoder(t *testing.T) {
	var buf bytes.Buffer
	log := newTestLogger(t, &buf, logger.FormatJSON)

	err := codedError{msg: "connection refused", code: "DB_CONN"}
	log.Error(context.Background(), "db write failed", logger.Err(err))

	m := decodeLine(t, &buf)
	assert.Equal(t, "connection refused", m["errorMessage"])
	assert.Equal(t, "DB_CONN", m["errorCode"])
}

func TestErrOmitsEmptyErrorCode(t *testing.T) {
	var buf bytes.Buffer
	log := newTestLogger(t, &buf, logger.FormatJSON)

	err := codedError{msg: "no code", code: ""}
	log.Error(context.Background(), "boom", logger.Err(err))

	m := decodeLine(t, &buf)
	assert.Equal(t, "no code", m["errorMessage"])
	assert.NotContains(t, m, "errorCode")
}

func TestInfoLevelHasNoStack(t *testing.T) {
	var buf bytes.Buffer
	log := newTestLogger(t, &buf, logger.FormatJSON)

	log.Info(context.Background(), "no stack here")

	assert.NotContains(t, decodeLine(t, &buf), "errorStack")
}

func TestFatalLogsAndExits(t *testing.T) {
	var buf bytes.Buffer
	var gotCode int
	called := false
	cfg := logger.Config{
		Service:     "svc",
		Environment: "production",
		Version:     "v1",
		Output:      &buf,
	}.WithClock(fixedClock()).WithExit(func(code int) {
		called = true
		gotCode = code
	})
	log, err := logger.New(cfg)
	require.NoError(t, err)

	log.Fatal(context.Background(), "cannot continue")

	require.True(t, called)
	assert.Equal(t, 1, gotCode)
	assert.Equal(t, "fatal", decodeLine(t, &buf)["level"])
}

func TestNewRequiresIdentityFields(t *testing.T) {
	cases := map[string]logger.Config{
		"missing service":     {Environment: "production", Version: "v1"},
		"missing environment": {Service: "svc", Version: "v1"},
		"missing version":     {Service: "svc", Environment: "production"},
	}
	for name, cfg := range cases {
		t.Run(name, func(t *testing.T) {
			log, err := logger.New(cfg)
			require.Error(t, err)
			assert.Nil(t, log)
		})
	}
}

func TestDisableStackTraceOmitsStack(t *testing.T) {
	var buf bytes.Buffer
	cfg := logger.Config{
		Service:           "svc",
		Environment:       "production",
		Version:           "v1",
		Level:             logger.LevelTrace,
		Output:            &buf,
		DisableStackTrace: true,
	}.WithClock(fixedClock())
	log, err := logger.New(cfg)
	require.NoError(t, err)

	log.Error(context.Background(), "boom", logger.Err(context.DeadlineExceeded))

	m := decodeLine(t, &buf)
	assert.Equal(t, context.DeadlineExceeded.Error(), m["errorMessage"])
	assert.NotContains(t, m, "errorStack")
}

func TestGroupNestingKeepsBaseFieldsTopLevel(t *testing.T) {
	var buf bytes.Buffer
	log := newTestLogger(t, &buf, logger.FormatJSON)

	tc := logger.TraceContext{TraceID: "4bf92f3577b34da6a3ce929d0e0e4736", SpanID: "00f067aa0ba902b7", Flags: "01"}
	ctx := logger.ContextWithTrace(context.Background(), tc)

	log.Info(ctx, "grouped",
		logger.Group("request",
			logger.String("method", "POST"),
			logger.Int("status", 201),
		),
	)

	m := decodeLine(t, &buf)
	// Base fields stay at the top level even when the caller nests via Group.
	assert.Equal(t, "4bf92f3577b34da6a3ce929d0e0e4736", m["traceId"])
	assert.Equal(t, "rsk-bridge-api", m["service"])

	request, ok := m["request"].(map[string]any)
	require.True(t, ok, "request should be a nested object")
	assert.Equal(t, "POST", request["method"])
	assert.EqualValues(t, 201, request["status"])
}

func TestLogfmtHandlerOutput(t *testing.T) {
	var buf bytes.Buffer
	log := newTestLogger(t, &buf, logger.FormatLogfmt)

	log.Info(context.Background(), "block processed")

	line := buf.String()
	assert.Contains(t, line, "level=info")
	assert.Contains(t, line, "service=rsk-bridge-api")
	assert.Contains(t, line, `message="block processed"`)
	assert.Contains(t, line, "timestamp=2026-04-14T15:30:00.123Z")
}

func TestWarnLogsAtWarnLevel(t *testing.T) {
	var buf bytes.Buffer
	log := newTestLogger(t, &buf, logger.FormatJSON)

	log.Warn(context.Background(), "slow dependency")

	assert.Equal(t, "warn", decodeLine(t, &buf)["level"])
}

func TestEnabledRespectsConfiguredLevel(t *testing.T) {
	var buf bytes.Buffer
	cfg := logger.Config{
		Service:     "svc",
		Environment: "production",
		Version:     "v1",
		Level:       logger.LevelWarn,
		Output:      &buf,
	}.WithClock(fixedClock())
	log, err := logger.New(cfg)
	require.NoError(t, err)

	assert.False(t, log.Enabled(context.Background(), logger.LevelInfo))
	assert.True(t, log.Enabled(context.Background(), logger.LevelWarn))
	assert.True(t, log.Enabled(context.Background(), logger.LevelError))
}

func TestHandlerWriteErrorDoesNotPanic(t *testing.T) {
	cfg := logger.Config{
		Service:     "svc",
		Environment: "production",
		Version:     "v1",
		Output:      failingWriter{},
	}.WithClock(fixedClock())
	log, err := logger.New(cfg)
	require.NoError(t, err)

	assert.NotPanics(t, func() {
		log.Info(context.Background(), "still safe")
	})
}

func TestNilContextIsAccepted(t *testing.T) {
	var buf bytes.Buffer
	log := newTestLogger(t, &buf, logger.FormatJSON)

	// Callers may pass a nil context; the logger substitutes Background.
	log.Info(nil, "nil ctx") //nolint:staticcheck // intentional nil context for the nil-guard path

	m := decodeLine(t, &buf)
	assert.Equal(t, "nil ctx", m["message"])
	require.Contains(t, m, "traceId")
	assert.Empty(t, m["traceId"])
}

func TestNewDefaultsNilOutputForOTel(t *testing.T) {
	provider := sdklog.NewLoggerProvider()
	t.Cleanup(func() {
		require.NoError(t, provider.Shutdown(context.Background()))
	})

	log, err := logger.New(logger.Config{
		Service:            "svc",
		Environment:        "production",
		Version:            "v1",
		Format:             logger.FormatOTel,
		OTelLoggerProvider: provider,
		// Output left nil; New defaults it to os.Stderr even though OTel ignores it.
	}.WithClock(fixedClock()))
	require.NoError(t, err)
	require.NotNil(t, log)

	assert.NotPanics(t, func() {
		log.Info(context.Background(), "otel without output")
	})
}

func TestNewDefaultsEmptyFormatToJSON(t *testing.T) {
	var buf bytes.Buffer
	log, err := logger.New(logger.Config{
		Service:     "svc",
		Environment: "production",
		Version:     "v1",
		Output:      &buf,
		// Format left empty; New defaults to JSON.
	}.WithClock(fixedClock()))
	require.NoError(t, err)

	log.Info(context.Background(), "default format")

	m := decodeLine(t, &buf)
	assert.Equal(t, "default format", m["message"])
}

func TestUnknownFormatFallsBackToJSON(t *testing.T) {
	var buf bytes.Buffer
	log, err := logger.New(logger.Config{
		Service:     "svc",
		Environment: "production",
		Version:     "v1",
		Format:      logger.Format("yaml"),
		Output:      &buf,
	}.WithClock(fixedClock()))
	require.NoError(t, err)

	log.Info(context.Background(), "fallback")

	m := decodeLine(t, &buf)
	assert.Equal(t, "fallback", m["message"])
}

// Ensure failingWriter satisfies io.Writer at compile time.
var _ io.Writer = failingWriter{}

// TestOutputPassesConformance is the key round-trip: records produced by the
// logger must validate against the standard for every supported format.
func TestOutputPassesConformance(t *testing.T) {
	formats := []struct {
		format     logger.Format
		confFormat logcheck.Format
	}{
		{logger.FormatJSON, logcheck.FormatJSON},
		{logger.FormatLogfmt, logcheck.FormatLogfmt},
	}

	for _, tc := range formats {
		t.Run(string(tc.format), func(t *testing.T) {
			var buf bytes.Buffer
			log := newTestLogger(t, &buf, tc.format)

			trace := logger.TraceContext{TraceID: "4bf92f3577b34da6a3ce929d0e0e4736", SpanID: "00f067aa0ba902b7", Flags: "01"}
			ctx := logger.ContextWithTrace(context.Background(), trace)
			log.Info(ctx, "block processed", logger.Uint64("blockNumber", 12345), logger.String("txHash", "0xdeadbeef"))
			log.Error(ctx, "db write failed", logger.Err(context.DeadlineExceeded))

			result, err := logcheck.Validate(bytes.NewReader(buf.Bytes()), logcheck.Options{Format: tc.confFormat, Strict: true})
			require.NoError(t, err)
			assert.True(t, result.OK(), "expected conformance to pass, result: %+v", result.Lines)
			assert.Equal(t, 2, result.Passed)
		})
	}
}
