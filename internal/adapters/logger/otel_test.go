package logger_test

import (
	"context"
	"sync"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/logger"

	otellog "go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// memExporter is a test-only OTel log exporter that retains records in memory.
type memExporter struct {
	mu      sync.Mutex
	records []sdklog.Record
}

func (e *memExporter) Export(_ context.Context, records []sdklog.Record) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, r := range records {
		e.records = append(e.records, r.Clone())
	}
	return nil
}

func (e *memExporter) Shutdown(context.Context) error   { return nil }
func (e *memExporter) ForceFlush(context.Context) error { return nil }

func (e *memExporter) get() []sdklog.Record {
	e.mu.Lock()
	defer e.mu.Unlock()
	out := make([]sdklog.Record, len(e.records))
	copy(out, e.records)
	return out
}

func newOTelTestLogger(t *testing.T, exp *memExporter) *logger.Logger {
	t.Helper()
	provider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewSimpleProcessor(exp)),
	)
	t.Cleanup(func() {
		if err := provider.Shutdown(context.Background()); err != nil {
			t.Errorf("OTel LoggerProvider shutdown: %v", err)
		}
	})

	log, err := logger.New(logger.Config{
		Service:            "lps",
		Environment:        "production",
		Version:            "v1.4.2",
		Level:              logger.LevelTrace,
		Format:             logger.FormatOTel,
		OTelLoggerProvider: provider,
	}.WithClock(fixedClock()))
	require.NoError(t, err)
	return log
}

func attrMap(r sdklog.Record) map[string]string {
	out := make(map[string]string)
	r.WalkAttributes(func(kv otellog.KeyValue) bool {
		flattenOTelValue(out, kv.Key, kv.Value)
		return true
	})
	return out
}

func flattenOTelValue(out map[string]string, key string, v otellog.Value) {
	switch v.Kind() {
	case otellog.KindMap:
		for _, child := range v.AsMap() {
			childKey := child.Key
			if key != "" {
				childKey = key + "." + child.Key
			}
			flattenOTelValue(out, childKey, child.Value)
		}
	default:
		out[key] = v.AsString()
	}
}

func TestOTelHandlerShape(t *testing.T) {
	exp := &memExporter{}
	log := newOTelTestLogger(t, exp)

	tc := logger.TraceContext{
		TraceID: "4bf92f3577b34da6a3ce929d0e0e4736",
		SpanID:  "00f067aa0ba902b7",
		Flags:   "01",
	}
	ctx := logger.ContextWithTrace(context.Background(), tc)

	log.Info(ctx, "Bridge transaction initiated", logger.String("txHash", "0x1234abcd"))

	records := exp.get()
	require.Len(t, records, 1)
	r := records[0]

	assert.Equal(t, "Bridge transaction initiated", r.Body().AsString())
	assert.Equal(t, otellog.SeverityInfo, r.Severity())

	attrs := attrMap(r)
	assert.Equal(t, "0x1234abcd", attrs["txHash"])
	assert.Equal(t, "lps", attrs["service"])
	assert.Equal(t, "production", attrs["environment"])
	assert.Equal(t, "v1.4.2", attrs["version"])
	assert.Equal(t, "4bf92f3577b34da6a3ce929d0e0e4736", attrs["traceId"])
	assert.Equal(t, "00f067aa0ba902b7", attrs["spanId"])
}

func TestOTelHandlerRedactsAttributes(t *testing.T) {
	exp := &memExporter{}
	log := newOTelTestLogger(t, exp)

	log.Info(context.Background(), "leak", logger.String("privateKey", "0xabc123"))

	records := exp.get()
	require.Len(t, records, 1)

	attrs := attrMap(records[0])
	assert.Equal(t, "[REDACTED]", attrs["privateKey"])
}

func TestOTelHandlerRedactsNestedGroups(t *testing.T) {
	exp := &memExporter{}
	log := newOTelTestLogger(t, exp)

	log.Info(context.Background(), "nested",
		logger.Group("credentials",
			logger.String("privateKey", "0xabc123"),
			logger.String("user", "alice"),
		),
		// Empty-key group is inlined by slog; still must be redacted.
		logger.Group("",
			logger.String("apiKey", "supersecretab3f"),
		),
	)

	records := exp.get()
	require.Len(t, records, 1)
	attrs := attrMap(records[0])

	assert.Equal(t, "alice", attrs["credentials.user"])
	assert.Equal(t, "[REDACTED]", attrs["credentials.privateKey"])
	assert.Equal(t, "****ab3f", attrs["apiKey"])
}

func TestOTelWithBindsAndRedacts(t *testing.T) {
	exp := &memExporter{}
	log := newOTelTestLogger(t, exp).With(logger.String("privateKey", "0xabc123"))

	log.Info(context.Background(), "bound")

	records := exp.get()
	require.Len(t, records, 1)
	attrs := attrMap(records[0])
	assert.Equal(t, "[REDACTED]", attrs["privateKey"])
}
