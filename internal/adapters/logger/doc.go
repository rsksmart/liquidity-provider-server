// Package logger is a structured logging library built on top of the standard
// library's log/slog.
//
// It enforces the organisation's logging standard: machine-parseable output
// (JSON by default, logfmt, and OpenTelemetry via the official otelslog
// bridge), mandatory base fields on every line (timestamp, level, service,
// environment, traceId, message, version), W3C Trace Context propagation,
// level support including a fatal level, automatic censorship of sensitive
// data, and error stack traces.
//
// JSON and logfmt paths depend only on the Go standard library. The OTel format
// additionally uses go.opentelemetry.io/contrib/bridges/otelslog and
// go.opentelemetry.io/otel/log so records can be delivered through a
// consumer-owned LoggerProvider.
//
// This root package is the consumer-facing facade: it re-exports the types
// consumers need (Level, Format, TraceContext, Secret, Masked, RedactionConfig,
// Field, Coder) with aliases and thin wrappers.
package logger
