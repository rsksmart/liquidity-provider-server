// Error log shape and logFatal wrapper; stack capture is hand-rolled because
// logrus has no native errorStack field.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

func demoErrorFields(ctx context.Context) {
	root := errors.New("ECONNREFUSED 127.0.0.1:4444")
	wrapped := fmt.Errorf("lbc.getQuote: %w", root)
	logError(ctx, "outbound rpc failed", "RPC_DIAL_FAILED", wrapped, logrus.Fields{
		"integrationName":   "rsk-lbc",
		"integrationMethod": "getQuote",
	})
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
	baseEntry(ctx).WithFields(logrus.Fields{
		"errorCode":    "DEMO_FATAL",
		"errorMessage": "forced",
	}).Error("fatal-equivalent (DEMO_FATAL_EXIT=1 to actually exit)")
}

// captureStack mirrors the slog demo so the error log shape is
// comparable across libraries. logrus has no native stack capture.
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

// logError emits the four standard error fields. logrus's WithError
// stashes the error under "error"; the standard wants errorMessage so
// we duplicate it explicitly.
func logError(ctx context.Context, msg, code string, err error, extra logrus.Fields) {
	if extra == nil {
		extra = logrus.Fields{}
	}
	fields := logrus.Fields{
		"errorCode":    code,
		"errorMessage": err.Error(),
		"errorStack":   captureStack(3),
		"errorContext": extra,
	}
	baseEntry(ctx).WithFields(fields).Error(msg)
}

// logFatal exits 1 after logging. logrus has logrus.Fatal but it would
// not run hooks the same way for our wrapper, so we own the exit.
func logFatal(ctx context.Context, msg, code string, err error) {
	logError(ctx, msg, code, err, nil)
	os.Exit(1)
}
