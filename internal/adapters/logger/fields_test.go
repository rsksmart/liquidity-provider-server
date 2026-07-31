package logger_test

import (
	"log/slog"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFieldConstructors(t *testing.T) {
	ts := time.Date(2026, time.April, 14, 15, 30, 0, 0, time.UTC)
	dur := 250 * time.Millisecond

	cases := []struct {
		name string
		got  logger.Field
		want slog.Attr
	}{
		{"String", logger.String("k", "v"), slog.String("k", "v")},
		{"Int", logger.Int("n", 7), slog.Int("n", 7)},
		{"Int64", logger.Int64("n64", 42), slog.Int64("n64", 42)},
		{"Uint64", logger.Uint64("u", 9), slog.Uint64("u", 9)},
		{"Float64", logger.Float64("f", 1.5), slog.Float64("f", 1.5)},
		{"Bool", logger.Bool("ok", true), slog.Bool("ok", true)},
		{"Time", logger.Time("ts", ts), slog.Time("ts", ts)},
		{"Duration", logger.Duration("d", dur), slog.Duration("d", dur)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.True(t, tc.got.Equal(tc.want))
		})
	}

	// slog.Any with a string stores a string-kind value; the constructor still
	// must be exercised and return the expected key/value.
	anyField := logger.Any("any", "payload")
	assert.Equal(t, "any", anyField.Key)
	assert.Equal(t, "payload", anyField.Value.String())
}

func TestGroupConstructor(t *testing.T) {
	got := logger.Group("req", logger.String("method", "GET"), logger.Int("status", 200))

	require.Equal(t, "req", got.Key)
	require.Equal(t, slog.KindGroup, got.Value.Kind())
	children := got.Value.Group()
	require.Len(t, children, 2)
	assert.True(t, children[0].Equal(slog.String("method", "GET")))
	assert.True(t, children[1].Equal(slog.Int("status", 200)))
}
