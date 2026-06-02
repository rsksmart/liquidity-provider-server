// Error log shape and fatal-equivalent logging; stack capture is hand-rolled.
// zerolog.Fatal is available; demoFatal is guarded by DEMO_FATAL_EXIT.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func demoErrorFields(ctx context.Context) {
	root := errors.New("ECONNREFUSED 127.0.0.1:4444")
	wrapped := fmt.Errorf("lbc.getQuote: %w", root)
	logError(ctx, "outbound rpc failed", "RPC_DIAL_FAILED", wrapped, map[string]string{
		"integrationName":   "rsk-lbc",
		"integrationMethod": "getQuote",
	})
}

// demoFatal shows the fatal-equivalent log line. zerolog has Fatal() — it
// logs at Fatal level and calls os.Exit(1) — DEMO_FATAL_EXIT=1 opts in.
// By default we log at Error so the SPIKE output stays inspectable.
func demoFatal(ctx context.Context) {
	exitNow := os.Getenv("DEMO_FATAL_EXIT") == "1"
	if exitNow {
		logger.Fatal().
			Str("errorCode", "DEMO_FATAL").
			Str("errorMessage", "forced").
			Msg("fatal-equivalent")
		return
	}
	loggerWithTrace(ctx).Error().
		Str("errorCode", "DEMO_FATAL").
		Str("errorMessage", "forced").
		Msg("fatal-equivalent (DEMO_FATAL_EXIT=1 to actually exit)")
}

// captureStack — zerolog has zerolog.ErrorStackMarshaler but it depends
// on a third-party errors package (pkg/errors style). We capture the
// stack ourselves for symmetry with the other demos and to avoid pulling
// another dep into this spike. This is why zerolog ranks "hook needed"
// on the stack-trace axis.
func captureStack(skip int) string {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(skip, pcs[:])
	if n == 0 {
		return ""
	}
	frames := runtime.CallersFrames(pcs[:n])
	var b strings.Builder
	for {
		f, more := frames.Next()
		b.WriteString(f.Function)
		b.WriteByte('\n')
		b.WriteByte('\t')
		b.WriteString(f.File)
		b.WriteByte(':')
		b.WriteString(strconv.Itoa(f.Line))
		b.WriteByte('\n')
		if !more {
			break
		}
	}
	return b.String()
}

// logError emits the four standard error fields. extra is folded into
// errorContext as a sub-object via .Interface so the JSON shape matches
// the slog and zap demos.
func logError(ctx context.Context, msg, code string, err error, extra map[string]string) {
	loggerWithTrace(ctx).Error().
		Str("errorCode", code).
		Str("errorMessage", err.Error()).
		Str("errorStack", captureStack(3)).
		Interface("errorContext", extra).
		Msg(msg)
}
