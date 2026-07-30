package logger

import "github.com/rsksmart/liquidity-provider-server/pkg/logger/internal/core"

// Level is the severity of a log record. See the level constants below.
type Level = core.Level

// Severity levels, ordered from most to least verbose.
const (
	// LevelTrace is the most verbose level, for fine-grained tracing.
	LevelTrace = core.LevelTrace
	// LevelDebug is for debugging information.
	LevelDebug = core.LevelDebug
	// LevelInfo is the default level for normal operational messages.
	LevelInfo = core.LevelInfo
	// LevelWarn is for conditions that deserve attention but are not errors.
	LevelWarn = core.LevelWarn
	// LevelError is for errors; a stack trace is attached by default.
	LevelError = core.LevelError
	// LevelFatal is for unrecoverable errors; logging at this level exits the
	// process.
	LevelFatal = core.LevelFatal
)

// Format selects the wire format of the emitted logs.
type Format = core.Format

// Supported output formats.
const (
	// FormatJSON emits one JSON object per line (default).
	FormatJSON = core.FormatJSON
	// FormatLogfmt emits key=value pairs per line.
	FormatLogfmt = core.FormatLogfmt
	// FormatOTel bridges records into an OpenTelemetry LoggerProvider via
	// otelslog. Config.Output is unused for this format.
	FormatOTel = core.FormatOTel
)

// ParseLevel converts a case-insensitive level name into a Level. Unknown names
// return an error and LevelInfo.
func ParseLevel(name string) (Level, error) {
	return core.ParseLevel(name)
}
