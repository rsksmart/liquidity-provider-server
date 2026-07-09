package logging

import (
	"io"
	"log/slog"
)

// NewTestHandler builds the production handler stack for tests in package logging_test.
func NewTestHandler(output io.Writer, level *slog.LevelVar) slog.Handler {
	base, err := newFormatHandler("json", output, level)
	if err != nil {
		panic(err)
	}
	return newContextHandler(base)
}

// ResetLevelVar restores the package level to info between tests.
func ResetLevelVar() {
	levelVar.Set(slog.LevelInfo)
}


// CaptureStackForTest exposes stack capture for tests.
func CaptureStackForTest(skip int) string {
	return captureStack(skip)
}
