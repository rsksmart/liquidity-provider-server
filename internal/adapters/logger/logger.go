package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"time"

	otellog "go.opentelemetry.io/otel/log"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/logger/internal/core"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/logger/internal/handler"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/logger/internal/redact"
)

// callerSkip is the number of stack frames between runtime.Callers (invoked
// inside log) and the user call site: Callers -> log -> level method -> caller.
const callerSkip = 3

// errMissingField builds the error returned by New when a mandatory identity
// field is empty.
func errMissingField(name string) error {
	return fmt.Errorf("logger: Config.%s is required", name)
}

// Config declares how a Logger behaves. Service, Environment and Version are
// the identity of the emitting service and are stamped on every line. The zero
// value is not usable directly; construct a Config (optionally via
// ConfigFromEnv) and pass it to New.
type Config struct {
	// Service is the service name (base field "service"). Required.
	Service string
	// Environment is the deploy environment (base field "environment"), e.g.
	// "production", "testnet". Required.
	Environment string
	// Version is the service version or git SHA (base field "version").
	// Required.
	Version string

	// Level is the minimum level that will be emitted. Defaults to LevelInfo.
	Level Level
	// Format is the output format. Defaults to FormatJSON.
	Format Format
	// Output is where records are written. Defaults to os.Stderr.
	Output io.Writer

	// AddSource includes the source file and line of the call site.
	AddSource bool
	// DisableStackTrace turns off the automatic "errorStack" field that is
	// otherwise attached to records at error and fatal levels. The zero value
	// keeps stack traces enabled.
	DisableStackTrace bool

	// Redaction controls censorship of sensitive data. See RedactionConfig.
	Redaction RedactionConfig

	// TimeZone is the location used to render timestamps. Defaults to UTC.
	TimeZone *time.Location

	// OTelLoggerProvider receives records when Format is FormatOTel. If nil,
	// the OTel global LoggerProvider is used (a no-op until the consumer sets
	// one). Ignored for other formats. Output is not used for FormatOTel.
	OTelLoggerProvider otellog.LoggerProvider

	// clock returns the current time. Overridable for deterministic tests.
	clock func() time.Time
	// exit terminates the process after a fatal log. Overridable for tests.
	exit func(int)
}

// Logger is a structured, leveled logger. It is safe for concurrent use by
// multiple goroutines. Create one with New and derive scoped children with
// With.
type Logger struct {
	handler slog.Handler
	clock   func() time.Time
	exit    func(int)
	stack   bool
}

// New builds a Logger from cfg, applying defaults for any unset fields. It
// returns an error when any of the mandatory identity fields (Service,
// Environment, Version) is empty, so misconfiguration fails fast at startup
// rather than emitting logs that violate the standard.
func New(cfg Config) (*Logger, error) {
	if cfg.Service == "" {
		return nil, errMissingField("Service")
	}
	if cfg.Environment == "" {
		return nil, errMissingField("Environment")
	}
	if cfg.Version == "" {
		return nil, errMissingField("Version")
	}

	if cfg.Output == nil {
		cfg.Output = os.Stderr
	}
	if cfg.Format == "" {
		cfg.Format = FormatJSON
	}
	if cfg.TimeZone == nil {
		cfg.TimeZone = time.UTC
	}
	if cfg.clock == nil {
		cfg.clock = time.Now
	}
	if cfg.exit == nil {
		cfg.exit = os.Exit
	}

	h := handler.New(handler.Options{
		Output:       cfg.Output,
		Format:       cfg.Format,
		Level:        cfg.Level,
		AddSource:    cfg.AddSource,
		TimeZone:     cfg.TimeZone,
		Redactor:     redact.New(cfg.Redaction),
		Service:      cfg.Service,
		Environment:  cfg.Environment,
		Version:      cfg.Version,
		OTelProvider: cfg.OTelLoggerProvider,
	})

	return &Logger{
		handler: h,
		clock:   cfg.clock,
		exit:    cfg.exit,
		stack:   !cfg.DisableStackTrace,
	}, nil
}

// With returns a child Logger that includes the given fields on every record.
// Use it to bind business context (e.g. an operationId) once instead of
// repeating it on each call.
func (l *Logger) With(fields ...Field) *Logger {
	child := *l
	child.handler = l.handler.WithAttrs(fields)
	return &child
}

// Enabled reports whether records at the given level would be emitted.
func (l *Logger) Enabled(ctx context.Context, level Level) bool {
	return l.handler.Enabled(ctx, level.SlogLevel())
}

// Trace logs at LevelTrace.
func (l *Logger) Trace(ctx context.Context, msg string, fields ...Field) {
	l.log(ctx, LevelTrace, msg, fields...)
}

// Debug logs at LevelDebug.
func (l *Logger) Debug(ctx context.Context, msg string, fields ...Field) {
	l.log(ctx, LevelDebug, msg, fields...)
}

// Info logs at LevelInfo.
func (l *Logger) Info(ctx context.Context, msg string, fields ...Field) {
	l.log(ctx, LevelInfo, msg, fields...)
}

// Warn logs at LevelWarn.
func (l *Logger) Warn(ctx context.Context, msg string, fields ...Field) {
	l.log(ctx, LevelWarn, msg, fields...)
}

// Error logs at LevelError. A stack trace is attached unless DisableStackTrace
// is set.
func (l *Logger) Error(ctx context.Context, msg string, fields ...Field) {
	l.log(ctx, LevelError, msg, fields...)
}

// Fatal logs at LevelFatal and then terminates the process via the configured
// exit function (os.Exit(1) by default). A fatal log means the service cannot
// continue running.
func (l *Logger) Fatal(ctx context.Context, msg string, fields ...Field) {
	l.log(ctx, LevelFatal, msg, fields...)
	l.exit(1)
}

// log is the single choke point that builds and dispatches a record. It uses
// the configured clock and captures the call site.
func (l *Logger) log(ctx context.Context, level Level, msg string, fields ...Field) { //nolint:contextcheck // Background is only a fallback for a nil ctx, where there is nothing to inherit from.
	if ctx == nil {
		ctx = context.Background()
	}
	if !l.handler.Enabled(ctx, level.SlogLevel()) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(callerSkip, pcs[:])

	record := slog.NewRecord(l.clock(), level.SlogLevel(), msg, pcs[0])
	record.AddAttrs(fields...)
	if l.stack && level >= LevelError {
		// +1 accounts for the captureStack frame itself, so the trace begins
		// at the user's call site.
		record.AddAttrs(slog.String(string(core.FieldErrorStack), captureStack(callerSkip+1)))
	}

	// Errors returned by handlers are not actionable at the call site; the
	// logger must never panic or interrupt program flow.
	_ = l.handler.Handle(ctx, record) //nolint:errcheck // a broken sink cannot be handled here.
}
