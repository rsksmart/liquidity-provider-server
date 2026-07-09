// Package logging is the LPS-owned slog wrapper for the logrus → slog migration.
//
// It is shipped additively (Phase 2): no production code imports this package until
// Phase 8 replaces setUpLogger in main with Init. That lets logrus keep running while
// later phases add trace plumbing, migrate call sites, and remove logrus in a final
// cutover—the only one-way door in the migration.
//
// Files are split per the migration plan, not because Go requires it:
//   - logging.go — public API (Init, Error, Fatal, Lazy) and handler assembly
//   - handler.go — slog.Handler wrappers (traceId from context, PII redaction)
//   - stack.go   — runtime.Callers capture for Error's errorStack field
//
// Follow-up phases (separate PRs): trace.go (W3C traceparent, middleware, StartJob),
// logtest helpers, LOG_FORMAT on Environment (Phase 1), then call-site migration.
//
// Config is standalone today; Phase 8 maps Environment and build metadata into it.
package logging

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"strings"

	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

// Config holds logger initialization inputs. Phase 8 wiring maps Environment and build vars here.
type Config struct {
	Service     string
	Environment string
	Version     string
	Format      string
	Level       string
	File        string
}

var levelVar = &slog.LevelVar{}

type traceIDKey struct{}

type lazyValuer struct {
	fn func() any
}

// LevelVar returns the shared level used by the default logger. Intended for logtest.SetLevel in Phase 5.
func LevelVar() *slog.LevelVar {
	return levelVar
}

// WithTraceID attaches a trace id to ctx for the context handler.
// Phase 3 trace.go will add W3C traceparent parsing, middleware, and StartJob; this
// minimal hook lets the handler stack be tested before that plumbing exists.
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, traceID)
}

// TraceIDFromContext returns the trace id stored in ctx, if any.
func TraceIDFromContext(ctx context.Context) string {
	traceID, ok := ctx.Value(traceIDKey{}).(string)
	if !ok {
		return ""
	}
	return traceID
}

// Init configures the default slog logger: handler stack, base fields, and optional file output.
func Init(cfg Config) error {
	if err := applyLevel(cfg.Level); err != nil {
		return err
	}

	output, err := openLogOutput(cfg.File)
	if err != nil {
		return err
	}

	baseHandler, err := newFormatHandler(cfg.Format, output, levelVar)
	if err != nil {
		return closeOutputOnError(output, err)
	}

	handler := newContextHandler(baseHandler)
	logger := slog.New(handler).With(
		slog.String("service", cfg.Service),
		slog.String("environment", cfg.Environment),
		slog.String("version", cfg.Version),
	)
	slog.SetDefault(logger)
	return nil
}

// Error logs msg with the standard LPS error shape. The err value is never wrapped or modified.
func Error(ctx context.Context, msg string, err error, attrs ...slog.Attr) {
	logAttrs := errorAttrs(err)
	logAttrs = append(logAttrs, attrs...)
	slog.Default().LogAttrs(ctx, slog.LevelError, msg, logAttrs...)
}

// Fatal logs once at error level and exits the process.
func Fatal(ctx context.Context, msg string, attrs ...slog.Attr) {
	slog.Default().LogAttrs(ctx, slog.LevelError, msg, attrs...)
	os.Exit(1)
}

// Lazy defers evaluation of expensive log field values.
func Lazy(fn func() any) slog.LogValuer {
	return lazyValuer{fn: fn}
}

func (v lazyValuer) LogValue() slog.Value {
	return slog.AnyValue(v.fn())
}

func applyLevel(level string) error {
	parsed, err := parseLevel(level)
	if err != nil {
		return err
	}
	levelVar.Set(parsed)
	return nil
}

func parseLevel(level string) (slog.Level, error) {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "trace":
		return slog.Level(-8), nil
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	case "fatal":
		return slog.LevelError, nil
	case "panic":
		return 0, fmt.Errorf("unsupported log level %q", level)
	default:
		return 0, fmt.Errorf("invalid log level %q", level)
	}
}

func newFormatHandler(format string, output io.Writer, level *slog.LevelVar) (slog.Handler, error) {
	opts := buildHandlerOptions(level)
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "json":
		return slog.NewJSONHandler(output, &opts), nil
	case "text", "logfmt":
		return slog.NewTextHandler(output, &opts), nil
	case "otel":
		// Phase 2: json and text/logfmt only. OTel export via otelslog lands in a follow-up.
		return nil, fmt.Errorf("log format %q is not supported yet; use json or text/logfmt", format)
	default:
		return nil, fmt.Errorf("invalid log format %q", format)
	}
}

func openLogOutput(filePath string) (io.Writer, error) {
	if filePath == "" {
		return os.Stderr, nil
	}
	if err := os.MkdirAll(path.Dir(filePath), 0o700); err != nil {
		return nil, fmt.Errorf("cannot create log file path (%s): %w", filePath, err)
	}
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, fmt.Errorf("cannot open log file %s: %w", filePath, err)
	}
	return file, nil
}

func closeOutputOnError(output io.Writer, err error) error {
	if file, ok := output.(*os.File); ok && file != os.Stderr {
		_ = file.Close()
	}
	return err
}

type errorShape struct {
	code    string
	message string
	stack   string
	context string
}

func errorAttrs(err error) []slog.Attr {
	if err == nil {
		return nil
	}
	shape := extractErrorShape(err)
	return []slog.Attr{
		slog.String("errorCode", shape.code),
		slog.String("errorMessage", shape.message),
		slog.String("errorStack", shape.stack),
		slog.String("errorContext", shape.context),
	}
}

func extractErrorShape(err error) errorShape {
	return errorShape{
		code:    sentinelCode(err),
		message: err.Error(),
		stack:   captureStack(stackSkipFrames),
		context: useCaseContext(err),
	}
}

func sentinelCode(err error) string {
	for _, sentinel := range lpsSentinels() {
		if errors.Is(err, sentinel) {
			return sentinel.Error()
		}
	}
	return innermostErrorMessage(err)
}

func innermostErrorMessage(err error) string {
	current := err
	var last string
	for current != nil {
		last = current.Error()
		current = errors.Unwrap(current)
	}
	return last
}

func useCaseContext(err error) string {
	msg := err.Error()
	useCase, remainder, found := strings.Cut(msg, ": ")
	if !found || useCase == "" || strings.Contains(useCase, " ") {
		return ""
	}
	if idx := strings.Index(remainder, ". Args: "); idx >= 0 {
		return useCase + remainder[idx:]
	}
	return useCase
}

func lpsSentinels() []error {
	return []error{
		usecases.NonRecoverableError,
		usecases.TxBelowMinimumError,
		usecases.RskAddressNotSupportedError,
		usecases.QuoteNotFoundError,
		usecases.QuoteNotAcceptedError,
		usecases.ExpiredQuoteError,
		usecases.NoLiquidityError,
		usecases.ProviderConfigurationError,
		usecases.WrongStateError,
		usecases.NoEnoughConfirmationsError,
		usecases.InsufficientAmountError,
		usecases.RegistrationRejectedError,
		usecases.RegistrationWithdrawnError,
		usecases.IllegalQuoteStateError,
		usecases.LockingCapExceededError,
		usecases.NonPositiveWeiError,
		usecases.EmptyConfirmationsMapError,
		usecases.NonPositiveConfirmationKeyError,
		usecases.NonPositiveReimbursementWindowError,
	}
}
