# Fluent logger

A structured logging library built on top of the standard library's
`[log/slog](https://pkg.go.dev/log/slog)`. It implements the organisation's
logging standard so compliant logging is the default.

- **Structured & machine-parseable**: JSON (default), logfmt, or OpenTelemetry
via the official [otelslog](https://pkg.go.dev/go.opentelemetry.io/contrib/bridges/otelslog) bridge.
- **Mandatory base fields** on every line: `timestamp`, `level`, `service`,
`environment`, `traceId`, `message`, `version`.
- **W3C Trace Context**: extract/inject `traceparent`, propagate `traceId` /
`spanId` through `context.Context`.
- **Levels**: `trace`, `debug`, `info`, `warn`, `error`, `fatal`.
- **Sensitive-data censorship**: typed `Secret` / `Masked` wrappers, a
field-name denylist, and conservative value scanning.
- **Error stack traces** at `error` and `fatal` levels.
- **Extensible business fields** via attributes and child loggers.
- **Stdlib-only for JSON/logfmt**; OTel format pulls in the official OpenTelemetry
slog bridge.



## Install

TODO, initial version of the package is integrated directly in the application

## Quick start

```go
log, err := logger.New(logger.ConfigFromEnv("lps", "production", "v2.5.2"))
if err != nil {
    panic(err)
}

ctx := context.Background()
log.Info(ctx, "Bridge transaction initiated",
    logger.String("txHash", "0x1234abcd"),
    logger.Uint64("blockNumber", 12345),
)
```

```json
{"timestamp":"2026-04-14T15:30:00.123Z","level":"info","service":"rsk-bridge-api","environment":"production","version":"v1.4.2","message":"Bridge transaction initiated","txHash":"0x1234abcd","blockNumber":12345,"traceId":""}
```

`New` returns an error when any of the mandatory identity fields (`Service`, `Environment`, `Version`) is empty, so misconfiguration surfaces at startup. `ConfigFromEnv` reads `LOG_LEVEL` and `LOG_FORMAT`; you can also build a `logger.Config` directly for full control (output writer, timezone, redaction, stack traces via `DisableStackTrace`, source, OTel provider).

## Package layout

The root package is a small consumer-facing facade; the implementation is
hidden under `internal/` and cannot be imported directly.

```
pkg/logger/            facade: Logger, Config, New, level methods, Field, and
                       re-exported aliases (Level, Format, TraceContext,
                       Secret, Masked, RedactionConfig)
  internal/core/       Level, Format, canonical field names, timestamp layout
  internal/redact/     redactor, Secret/Masked, denylists, value scanning
  internal/trace/      W3C Trace Context, propagation, net/http middleware
  internal/handler/    shared handler assembly (JSON, logfmt, context)
    otel/              otelslog bridge, OTel redaction, level filtering
```

Attributes are created with `logger.Field` constructors (`logger.String`,
`logger.Int`, …); `Field` is an alias of `slog.Attr`, so existing slog
attributes work directly.

## Format

Set `LOG_FORMAT=json|logfmt|otel` (or `Config.Format`). One format per service.

### OpenTelemetry (`FormatOTel`)

`FormatOTel` bridges slog records into an OpenTelemetry `LoggerProvider` via `otelslog`. Delivery (OTLP, stdout, etc.) and resource attributes are owned by the provider you inject; `Config.Output` is ignored for this format.

```go
import (
    "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
    sdklog "go.opentelemetry.io/otel/sdk/log"
)

exporter, err := otlploghttp.New(ctx)
if err != nil {
    return err
}
provider := sdklog.NewLoggerProvider(
    sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
)
defer provider.Shutdown(ctx)

log, err := logger.New(logger.Config{
    Service:            "rsk-bridge-api",
    Environment:        "production",
    Version:            "v1.4.2",
    Format:             logger.FormatOTel,
    OTelLoggerProvider: provider,
})
```

If `OTelLoggerProvider` is nil, the OTel global provider is used (a no-op until you set one), so otel logs are silently dropped when unconfigured.

Base identity fields (`service`, `environment`, `version`) and the custom `traceId` / `spanId` from `ContextWithTrace` are emitted as record attributes. Consumers who prefer semantic-convention Resource attributes can also set those on their provider. Severity numbers match the OTel mapping; `SeverityText` for custom levels follows slog (`DEBUG-4` for trace, `ERROR+4` for fatal) rather than the lowercase names used in JSON/logfmt.

## Trace context

```go
// In HTTP middleware (a ready-made net/http helper is provided):
handler = logger.TraceMiddleware(handler)

// Anywhere with the request context, the traceId/spanId are added automatically:
log.Info(ctx, "processing")
```

Manually: `logger.ParseTraceparent`, `logger.NewTraceContext`, `TraceContext.WithNewSpan`, `TraceContext.Traceparent`, `logger.ContextWithTrace`, `logger.TraceFromContext`. Trace/span generation returns errors from the system entropy source; malformed headers return the sentinel `logger.ErrInvalidTraceparent`. Parsed contexts preserve their W3C traceparent version.

## Child loggers and business fields

```go
opLog := log.With(logger.String("operationId", "pegout-42"))
opLog.Warn(ctx, "anomaly detected", logger.String("btcTxHash", hash))
```

Nest related fields under a key with `logger.Group`; the mandatory base fields
(including `traceId`) always stay at the top level:

```go
log.Info(ctx, "request handled",
    logger.Group("http",
        logger.String("method", "POST"),
        logger.Int("status", 201),
    ),
)
```



## Errors

```go
log.Error(ctx, "db write failed", logger.Err(err,
    logger.String("quoteId", quoteID),
    logger.Int("retries", retries),
))
// -> errorMessage (+errorCode), errorContext{quoteId,retries}, plus errorStack
```

Stack traces are attached automatically at `error` and `fatal` levels. Set `Config.DisableStackTrace` to turn them off.

If `err` (or one it wraps) implements `logger.Coder` (`Code() string`), `Err` also emits `errorCode`. Domain errors can satisfy that method without importing this package.

Optional trailing fields are nested under `errorContext` for state needed to reproduce the issue. Do not put secrets there — redaction still applies, but error context should stay non-sensitive by design.

## Sensitive data

Never appears in logs, three complementary mechanisms:

```go
// 1. Typed wrappers (recommended for anything that could hold a secret):
type Wallet struct {
    Address    string
    PrivateKey logger.Secret // always "[REDACTED]" when logged / JSON-marshaled
    APIKey     logger.Masked // "****ab3f" when logged / JSON-marshaled
}

log.Info(ctx, "loaded", logger.Any("key", logger.Secret(privKey)))     // -> "[REDACTED]"
log.Info(ctx, "auth",   logger.Any("apiKey", logger.Masked(apiKey)))   // -> "****ab3f"
log.Info(ctx, "wallet", logger.Any("wallet", wallet)) // Secret/Masked fields stay safe
```

```go
// 2. Field-name denylist: values under keys like privateKey, seed, mnemonic,
//    password (dropped) or apiKey, token, authorization (masked) are censored.

// 3. Value scanning: mnemonic phrases and bearer tokens are redacted.
```

Domain identifiers such as `txHash` / `blockHash` (64-hex) are intentionally **not** scanned so they remain visible for cross-referencing.

**Depth limitation**: censorship is applied per attribute. String values are scanned, and `map[string]any` / `[]any` values are censored recursively (so a decoded JSON body is still protected). Arbitrary structs are not field-walked by the redactor — type sensitive fields as `logger.Secret` / `logger.Masked` (they implement `slog.LogValuer` and safe JSON marshaling), or log flat attributes.

## Conformance validator (`logcheck`)
TODO language agnostic validator