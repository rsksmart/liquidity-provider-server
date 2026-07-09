package logging

// Stack capture for Error's errorStack attribute. Separated from logging.go so the
// public API file does not also own runtime frame-walking details.

import (
	"fmt"
	"runtime"
	"strings"
)

const stackSkipFrames = 3

func captureStack(skip int) string {
	const maxFrames = 32
	pcs := make([]uintptr, maxFrames)
	n := runtime.Callers(skip+1, pcs)
	if n == 0 {
		return ""
	}

	frames := runtime.CallersFrames(pcs[:n])
	var lines []string
	for {
		frame, more := frames.Next()
		if frame.Function == "" && frame.File == "" {
			break
		}
		if !strings.Contains(frame.Function, "github.com/rsksmart/liquidity-provider-server/internal/logging.") {
			lines = append(lines, fmt.Sprintf("%s:%d %s", frame.File, frame.Line, frame.Function))
		}
		if !more {
			break
		}
	}
	return strings.Join(lines, "\n")
}
