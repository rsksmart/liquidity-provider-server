// Error log shape with zap.Stack for errorStack and native Fatal support.
// demoFatal is guarded by DEMO_FATAL_EXIT so spike output stays inspectable.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"go.uber.org/zap"
)

func demoErrorFields(ctx context.Context) {
	root := errors.New("ECONNREFUSED 127.0.0.1:4444")
	wrapped := fmt.Errorf("lbc.getQuote: %w", root)
	logError(ctx, "outbound rpc failed", "RPC_DIAL_FAILED", wrapped,
		zap.String("integrationName", "rsk-lbc"),
		zap.String("integrationMethod", "getQuote"),
	)
}

// demoFatal shows the fatal-equivalent log line. zap has a real Fatal
// that calls os.Exit(1) after sync — DEMO_FATAL_EXIT=1 opts in. By
// default we log at Error so the SPIKE output stays inspectable.
func demoFatal(ctx context.Context) {
	exitNow := os.Getenv("DEMO_FATAL_EXIT") == "1"
	if exitNow {
		logger.Fatal("fatal-equivalent",
			zap.String("errorCode", "DEMO_FATAL"),
			zap.String("errorMessage", "forced"),
		)
		return
	}
	loggerWithTrace(ctx).Error("fatal-equivalent (DEMO_FATAL_EXIT=1 to actually exit)",
		zap.String("errorCode", "DEMO_FATAL"),
		zap.String("errorMessage", "forced"),
	)
}

// logError emits the four-field error log shape. zap's killer feature
// here is zap.Stack — native, well-tested stack capture — which is the
// reason zap wins the "Stack trace capture" axis in logs.md.
//
// We attach extra fields as a "errorContext" namespace so the structure
// is comparable to the slog demo.
func logError(ctx context.Context, msg, code string, err error, extra ...zap.Field) {
	fields := []zap.Field{
		zap.String("errorCode", code),
		zap.String("errorMessage", err.Error()),
		zap.Stack("errorStack"),
		zap.Namespace("errorContext"),
	}
	fields = append(fields, extra...)
	loggerWithTrace(ctx).Error(msg, fields...)
}
