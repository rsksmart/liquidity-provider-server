package otel_test

import (
	"context"
	"log/slog"
	"sync"
	"testing"
	"time"

	otellog "go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"

	otelhandler "github.com/rsksmart/liquidity-provider-server/pkg/logger/internal/handler/otel"
	"github.com/rsksmart/liquidity-provider-server/pkg/logger/internal/redact"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func attrMap(r sdklog.Record) map[string]string {
	out := make(map[string]string)
	r.WalkAttributes(func(kv otellog.KeyValue) bool {
		flatten(out, kv.Key, kv.Value)
		return true
	})
	return out
}

func flatten(out map[string]string, key string, v otellog.Value) {
	switch v.Kind() {
	case otellog.KindMap:
		for _, child := range v.AsMap() {
			childKey := child.Key
			if key != "" {
				childKey = key + "." + child.Key
			}
			flatten(out, childKey, child.Value)
		}
	default:
		out[key] = v.AsString()
	}
}

func newTestHandler(t *testing.T, exp *memExporter, level slog.Level) slog.Handler {
	t.Helper()
	provider := sdklog.NewLoggerProvider(sdklog.WithProcessor(sdklog.NewSimpleProcessor(exp)))
	t.Cleanup(func() {
		require.NoError(t, provider.Shutdown(context.Background()))
	})
	return otelhandler.New(otelhandler.Options{
		Service:  "svc",
		Level:    level,
		Redactor: redact.New(redact.Config{}),
		Provider: provider,
	})
}

func TestNewRedactsAttributes(t *testing.T) {
	exp := &memExporter{}
	h := newTestHandler(t, exp, slog.LevelDebug)

	record := slog.NewRecord(time.Now().UTC(), slog.LevelInfo, "leak", 0)
	record.AddAttrs(slog.String("privateKey", "0xabc"), slog.String("txHash", "0x1"))
	require.NoError(t, h.Handle(context.Background(), record))

	records := exp.get()
	require.Len(t, records, 1)
	assert.Equal(t, "leak", records[0].Body().AsString())

	attrs := attrMap(records[0])
	assert.Equal(t, "0x1", attrs["txHash"])
	assert.Equal(t, redact.DefaultPlaceholder, attrs["privateKey"])
}

func TestNewRedactsNestedGroups(t *testing.T) {
	exp := &memExporter{}
	h := newTestHandler(t, exp, slog.LevelDebug)

	record := slog.NewRecord(time.Now().UTC(), slog.LevelInfo, "nested", 0)
	record.AddAttrs(
		slog.Group("credentials",
			slog.String("privateKey", "0xabc"),
			slog.String("user", "alice"),
		),
		slog.Group("",
			slog.String("apiKey", "supersecretab3f"),
		),
	)
	require.NoError(t, h.Handle(context.Background(), record))

	attrs := attrMap(exp.get()[0])
	assert.Equal(t, "alice", attrs["credentials.user"])
	assert.Equal(t, redact.DefaultPlaceholder, attrs["credentials.privateKey"])
	assert.Equal(t, "****ab3f", attrs["apiKey"])
}

func TestWithAttrsRedacts(t *testing.T) {
	exp := &memExporter{}
	h := newTestHandler(t, exp, slog.LevelDebug).WithAttrs([]slog.Attr{
		slog.String("privateKey", "0xabc"),
	})

	record := slog.NewRecord(time.Now().UTC(), slog.LevelInfo, "bound", 0)
	require.NoError(t, h.Handle(context.Background(), record))

	attrs := attrMap(exp.get()[0])
	assert.Equal(t, redact.DefaultPlaceholder, attrs["privateKey"])
}

func TestWithAttrsEmptyIsNoop(t *testing.T) {
	exp := &memExporter{}
	base := newTestHandler(t, exp, slog.LevelDebug)
	assert.Equal(t, base, base.WithAttrs(nil))
}

func TestWithGroupEmptyIsNoop(t *testing.T) {
	exp := &memExporter{}
	base := newTestHandler(t, exp, slog.LevelDebug)
	assert.Equal(t, base, base.WithGroup(""))
}

func TestWithGroupNestsAttributes(t *testing.T) {
	exp := &memExporter{}
	h := newTestHandler(t, exp, slog.LevelDebug).WithGroup("req")

	record := slog.NewRecord(time.Now().UTC(), slog.LevelInfo, "grouped", 0)
	record.AddAttrs(slog.String("privateKey", "0xabc"), slog.String("method", "POST"))
	require.NoError(t, h.Handle(context.Background(), record))

	attrs := attrMap(exp.get()[0])
	assert.Equal(t, "POST", attrs["req.method"])
	assert.Equal(t, redact.DefaultPlaceholder, attrs["req.privateKey"])
}

func TestEnabledRespectsMinimumLevel(t *testing.T) {
	exp := &memExporter{}
	h := newTestHandler(t, exp, slog.LevelWarn)

	assert.False(t, h.Enabled(context.Background(), slog.LevelInfo))
	assert.True(t, h.Enabled(context.Background(), slog.LevelWarn))
	assert.True(t, h.Enabled(context.Background(), slog.LevelError))
}

func TestNewWithoutProviderUsesGlobal(t *testing.T) {
	// Nil provider is allowed; otelslog falls back to the global provider (a no-op).
	h := otelhandler.New(otelhandler.Options{
		Service:  "svc",
		Level:    slog.LevelInfo,
		Redactor: redact.New(redact.Config{}),
	})
	require.NotNil(t, h)
	assert.NotPanics(t, func() {
		record := slog.NewRecord(time.Now().UTC(), slog.LevelInfo, "noop", 0)
		require.NoError(t, h.Handle(context.Background(), record))
	})
}
