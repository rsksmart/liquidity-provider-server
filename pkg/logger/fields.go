package logger

import (
	"log/slog"
	"time"
)

// Field is a single structured key/value pair attached to a log record. It is
// an alias of slog.Attr, so slog attributes can be passed directly and the
// constructors below can be used without importing log/slog.
type Field = slog.Attr

// String returns a string field.
func String(key, value string) Field { return slog.String(key, value) }

// Int returns an int field.
func Int(key string, value int) Field { return slog.Int(key, value) }

// Int64 returns an int64 field.
func Int64(key string, value int64) Field { return slog.Int64(key, value) }

// Uint64 returns a uint64 field.
func Uint64(key string, value uint64) Field { return slog.Uint64(key, value) }

// Float64 returns a float64 field.
func Float64(key string, value float64) Field { return slog.Float64(key, value) }

// Bool returns a boolean field.
func Bool(key string, value bool) Field { return slog.Bool(key, value) }

// Time returns a time field.
func Time(key string, value time.Time) Field { return slog.Time(key, value) }

// Duration returns a duration field.
func Duration(key string, value time.Duration) Field { return slog.Duration(key, value) }

// Any returns a field for an arbitrary value.
func Any(key string, value any) Field { return slog.Any(key, value) }

// Group returns a field whose value is a group of the given fields, nested
// under key. Use it for structured nesting; the base fields stay at the top
// level regardless.
func Group(key string, fields ...Field) Field {
	args := make([]any, len(fields))
	for i, f := range fields {
		args[i] = f
	}
	return slog.Group(key, args...)
}
