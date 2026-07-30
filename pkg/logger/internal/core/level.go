package core

import (
	"fmt"
	"log/slog"
	"strings"
)

// Level is the severity of a log record. The values are aligned with
// [slog.Level] so a Level converts to an slog.Level with a plain numeric
// conversion. Two extra levels sit on top of the slog set: Trace (below Debug)
// and Fatal (above Error).
type Level int

// Severity levels ordered from most to least verbose. The numeric gaps match
// slog conventions (4 apart) so intermediate custom levels remain expressible.
const (
	LevelTrace Level = -8
	LevelDebug Level = -4
	LevelInfo  Level = 0
	LevelWarn  Level = 4
	LevelError Level = 8
	LevelFatal Level = 12
)

// String returns the lowercase canonical name of the level, which is what gets
// written to every log line (e.g. "info", "fatal").
func (l Level) String() string {
	switch l {
	case LevelTrace:
		return "trace"
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	default:
		return slog.Level(l).String()
	}
}

// SlogLevel converts the Level to its [slog.Level] equivalent.
func (l Level) SlogLevel() slog.Level {
	return slog.Level(l)
}

// ParseLevel converts a case-insensitive level name into a Level. It accepts
// the six canonical names (empty defaults to info). Unknown names return an
// error and LevelInfo so callers can safely fall back.
func ParseLevel(name string) (Level, error) {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "trace":
		return LevelTrace, nil
	case "debug":
		return LevelDebug, nil
	case "info", "":
		return LevelInfo, nil
	case "warn", "warning":
		return LevelWarn, nil
	case "error", "err":
		return LevelError, nil
	case "fatal":
		return LevelFatal, nil
	default:
		return LevelInfo, fmt.Errorf("logger: unknown level %q", name)
	}
}

// Format selects the wire format of the emitted logs.
type Format string

const (
	// FormatJSON emits one JSON object per line. Default and recommended for
	// production.
	FormatJSON Format = "json"
	// FormatLogfmt emits key=value pairs per line. Useful for local, human
	// readable development output.
	FormatLogfmt Format = "logfmt"
	// FormatOTel bridges records into an OpenTelemetry LoggerProvider via the
	// official otelslog handler. Config.Output is unused for this format.
	FormatOTel Format = "otel"
)
