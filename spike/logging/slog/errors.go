// Error log shape (errorCode, errorMessage, errorStack, errorContext) and
// fatal-equivalent logging. Demonstrates centralized wrappers because slog has
// no native fatal level or stack field.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func demoErrorFields(ctx context.Context) {
	root := errors.New("ECONNREFUSED 127.0.0.1:4444")
	wrapped := fmt.Errorf("lbc.getQuote: %w", root)
	logError(ctx, "outbound rpc failed", "RPC_DIAL_FAILED", wrapped,
		slog.String("integrationName", "rsk-lbc"),
		slog.String("integrationMethod", "getQuote"),
	)
}

// demoFatal shows the fatal-equivalent log line. The real helper calls
// os.Exit(1); we set DEMO_FATAL_EXIT=1 to opt into that — by default we
// just log so the SPIKE output remains inspectable.
func demoFatal(ctx context.Context) {
	exitNow := os.Getenv("DEMO_FATAL_EXIT") == "1"
	if exitNow {
		logFatal(ctx, "would exit now", "DEMO_FATAL", errors.New("forced"))
		return
	}
	slog.LogAttrs(ctx, slog.LevelError, "fatal-equivalent (DEMO_FATAL_EXIT=1 to actually exit)",
		slog.String("errorCode", "DEMO_FATAL"),
		slog.String("errorMessage", "forced"),
	)
}

// captureStack returns a textual stack trace suitable for the
// errorStack field. slog has no native stack support; centralising
// capture here means call sites stay short.
//
// skip = 2 drops captureStack and its immediate caller (logError /
// logFatal), so the first frame is the application call site.
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

// logError emits the four-field error log shape required by the Base
// Logging Standards: errorCode, errorMessage, errorStack, errorContext.
//
// extra carries call-site context (e.g. integrationName) and ends up in
// errorContext as a slog.Group.
func logError(ctx context.Context, msg, code string, err error, extra ...slog.Attr) {
	if err == nil {
		err = errPlaceholder("nil error")
	}
	all := append([]slog.Attr{
		slog.String("errorCode", code),
		slog.String("errorMessage", err.Error()),
		slog.String("errorStack", captureStack(3)),
		slog.Any("errorContext", attrsAsGroup(extra)),
	}, extra...)
	slog.LogAttrs(ctx, slog.LevelError, msg, all...)
}

// logFatal logs once and exits 1. Centralising fatal here is the
// "Fatal semantics" entry in logs.md — slog has no native fatal level.
func logFatal(ctx context.Context, msg, code string, err error, extra ...slog.Attr) {
	logError(ctx, msg, code, err, extra...)
	os.Exit(1)
}

func attrsAsGroup(attrs []slog.Attr) slog.Value {
	out := make([]slog.Attr, len(attrs))
	copy(out, attrs)
	return slog.GroupValue(out...)
}

type errPlaceholder string

func (e errPlaceholder) Error() string { return string(e) }
