package core_test

import (
	"log/slog"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/pkg/logger/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseLevel(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    core.Level
		wantErr bool
	}{
		{name: "trace", input: "trace", want: core.LevelTrace},
		{name: "debug uppercase", input: "DEBUG", want: core.LevelDebug},
		{name: "info", input: "info", want: core.LevelInfo},
		{name: "empty defaults to info", input: "", want: core.LevelInfo},
		{name: "warn", input: "warn", want: core.LevelWarn},
		{name: "warning alias", input: "warning", want: core.LevelWarn},
		{name: "error", input: "error", want: core.LevelError},
		{name: "err alias", input: "err", want: core.LevelError},
		{name: "fatal", input: "fatal", want: core.LevelFatal},
		{name: "unknown", input: "verbose", want: core.LevelInfo, wantErr: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := core.ParseLevel(tc.input)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestLevelString(t *testing.T) {
	cases := map[core.Level]string{
		core.LevelTrace: "trace",
		core.LevelDebug: "debug",
		core.LevelInfo:  "info",
		core.LevelWarn:  "warn",
		core.LevelError: "error",
		core.LevelFatal: "fatal",
	}
	for level, want := range cases {
		assert.Equal(t, want, level.String(), want)
	}
	assert.Equal(t, "ERROR+2", core.Level(10).String())
}

func TestLevelSlogLevel(t *testing.T) {
	assert.Equal(t, slog.Level(0), core.LevelInfo.SlogLevel())
	assert.Equal(t, slog.Level(8), core.LevelError.SlogLevel())
	assert.Equal(t, slog.Level(-8), core.LevelTrace.SlogLevel())
}
