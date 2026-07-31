package logger

import (
	"errors"
	"log/slog"
	"runtime"
	"strconv"
	"strings"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/logger/internal/core"
)

// maxStackDepth caps the number of frames captured for errorStack.
const maxStackDepth = 32

// Coder is satisfied by errors that expose an internal error code. When an
// error passed to Err implements it (directly or via wrapping), the code is
// emitted as "errorCode".
//
// Domain packages can satisfy this interface without importing logger by
// defining Code() string on their error types.
type Coder interface {
	Code() string
}

// Err converts an error into fields describing it: "errorMessage" always, and
// "errorCode" when the error (or one it wraps) implements [Coder]. Optional
// context fields are nested under "errorContext" for state needed to reproduce
// the failure (input params, partial results, and similar). All of these are
// inlined at the top level of the record.
//
// The stack trace is added automatically by the logger at error and fatal
// levels (unless Config.DisableStackTrace is set), so Err intentionally does
// not capture it.
//
// If err is nil, Err returns an empty field, which slog discards.
func Err(err error, context ...Field) Field {
	if err == nil {
		return Field{}
	}

	fields := []Field{slog.String(string(core.FieldErrorMessage), err.Error())}
	var c Coder
	if errors.As(err, &c) {
		if code := c.Code(); code != "" {
			fields = append(fields, slog.String(string(core.FieldErrorCode), code))
		}
	}
	if len(context) > 0 {
		fields = append(fields, Group(string(core.FieldErrorContext), context...))
	}
	return inlineGroup(fields)
}

// inlineGroup returns a group field with an empty key. slog inlines the
// attributes of an empty-keyed group, so they appear at the top level.
func inlineGroup(fields []Field) Field {
	args := make([]any, len(fields))
	for i, f := range fields {
		args[i] = f
	}
	return slog.Group("", args...)
}

// captureStack returns a compact, newline-separated stack trace, formatted as
// "function\n\tfile:line" per frame. skip is the number of leading frames to
// omit, counted from runtime.Callers (0 = runtime.Callers, 1 = captureStack).
func captureStack(skip int) string {
	var pcs [maxStackDepth]uintptr
	n := runtime.Callers(skip, pcs[:])
	if n == 0 {
		return ""
	}

	frames := runtime.CallersFrames(pcs[:n])
	var b strings.Builder
	for more := true; more; {
		var frame runtime.Frame
		frame, more = frames.Next()
		if frame.Function != "" {
			writeFrame(&b, frame, more)
		}
	}
	return b.String()
}

func writeFrame(b *strings.Builder, frame runtime.Frame, more bool) {
	b.WriteString(frame.Function)
	b.WriteString("\n\t")
	b.WriteString(frame.File)
	b.WriteByte(':')
	b.WriteString(strconv.Itoa(frame.Line))
	if more {
		b.WriteByte('\n')
	}
}
