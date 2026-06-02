# slog demo checklist

This package is the recommended target for the FLY-2308 comparison in
[docs/spikes/logs.md](../../../docs/spikes/logs.md).

Run it with:

```sh
LOG_FORMAT=json LOG_LEVEL=debug go run ./spike/logging/slog
go test ./spike/logging/slog
```

## Requirement checklist

| Requirement | Evidence | Notes |
|-------------|----------|-------|
| `LOG_FORMAT` / `LOG_LEVEL` | `config.go`, `main.go` | JSON and text are stdlib. `logfmt` uses slog text as a close local-dev format. `otel` falls back to JSON in the runnable demo; the document points to `otelslog` for real wiring. |
| Base fields | `initLogger` in `utils.go` | `service`, `environment`, and `version` are attached once with `logger.With(...)`. |
| Real LPS examples | `real_project_examples.go` | Imports `internal/entities` and `internal/usecases` to show pegout-rebalance logs with real domain values. Also covers the sentinel + `WrapUseCaseError` migration: useCase id → `errorCode`, wrapping chain → `errorMessage`, and `errors.Is` against the sentinel still matches. |
| Request trace fields | `contextHandler` in `utils.go` | Handler reads `context.Context` and adds `traceId` / `spanId` without per-call-site fields. |
| Inbound trace extraction | `tracing.go` (`demoInboundTraceExtraction`, `traceMiddleware`) | Parses W3C `traceparent`, mints a trace when absent, stores trace state in `context.Context`. |
| Structured access log | `tracing.go` (`traceMiddleware`) | Replaces `gorilla/handlers.LoggingHandler` with one structured record containing the six required HTTP fields. |
| Outbound trace injection | `tracing.go` (`demoOutboundTraceInjection`, `traceTransport`) | Injects the same trace-id and a fresh span-id into `traceparent`. |
| Error shape | `errors.go` | `demoErrorFields` / `demoFatal` at top; `logError` emits `errorCode`, `errorMessage`, `errorStack`, and `errorContext`; stack capture is centralized. |
| Standard key names and PII redaction | `pii.go` | `demoPIIRedaction` at top; `HandlerOptions.ReplaceAttr` renames `time`/`msg` to `timestamp`/`message` and redacts denied PII keys. Cleanest interception point in the comparison. |
| Fatal-equivalent | `errors.go`, `main.go` | `logFatal` logs once and exits; demo path is guarded by `DEMO_FATAL_EXIT`. |
| Test capture | `capture_test.go` | Buffer-backed JSON handler with the same redaction and context handler stack. |

## Candidate caveats

- slog has no built-in fatal level. LPS should own `logging.Fatal`.
- slog has no built-in stack field. LPS should own `logging.Error` and capture
  `errorStack` centrally.
- Lazy evaluation is explicit through `slog.LogValuer` or an LPS
  `logging.Lazy` helper.
- Exact `logfmt` requires a small handler or third-party handler if the
  standards review requires more than slog text output.
