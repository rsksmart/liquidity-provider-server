# logrus demo checklist

This package is the current LPS baseline for the FLY-2308 comparison in
[docs/spikes/logs.md](../../../docs/spikes/logs.md).

Run it with:

```sh
LOG_FORMAT=json LOG_LEVEL=debug go run ./spike/logging/logrus
go test ./spike/logging/logrus
```

## Requirement checklist

| Requirement | Evidence | Notes |
|-------------|----------|-------|
| `LOG_FORMAT` / `LOG_LEVEL` | `config.go`, `main.go` | JSON and text are built in. `logfmt` uses text output. `otel` falls back to JSON because there is no first-party maintained logrus bridge. |
| Base fields | `baseEntry(ctx)` in `utils.go` | Works, but every call site must use the derived entry. Weaker than slog handler-level enrichment. |
| Real LPS examples | `real_project_examples.go` | Imports `internal/entities` and `internal/usecases` to show pegout-rebalance logs with real domain values. Also covers the sentinel + `WrapUseCaseError` migration: useCase id → `errorCode`, wrapping chain → `errorMessage`, `errors.Is` against the sentinel still matches. |
| Inbound trace extraction | `tracing.go` (`demoInboundTraceExtraction`, `traceMiddleware`) | Parses W3C `traceparent`, mints a trace when absent, stores trace state in `context.Context`. |
| Structured access log | `tracing.go` (`traceMiddleware`) | Replaces `gorilla/handlers.LoggingHandler` with one structured record containing the six required HTTP fields. |
| Outbound trace injection | `tracing.go` (`demoOutboundTraceInjection`, `traceTransport`) | Injects the same trace-id and a fresh span-id into `traceparent`. |
| Error shape | `errors.go` | `demoErrorFields` / `demoFatal` at top; `logError` emits all four error fields; stack capture is custom. |
| PII redaction | `pii.go` | `demoPIIRedaction` at top; `piiRedactHook` mutates `Entry.Data`. Usable and better than writer-level redaction. |
| Fatal-equivalent | `errors.go`, `main.go` | `logFatal` logs with the standard error shape and exits; demo path is guarded by `DEMO_FATAL_EXIT`. |
| Test capture | `capture_test.go` | Matches current LPS-style `SetOutput` capture, then parses JSON lines for assertions. |

## Candidate caveats

- logrus is in maintenance mode, so keeping it would need most of the same LPS
  wrapper work while staying on the rejected direction from the Go guideline.
- Context propagation is not first-class. `WithContext` stores context, but it
  does not add trace fields; the demo must call `baseEntry(ctx)` everywhere.
- The current LPS codebase barely uses logrus structured fields, so the
  migration is still broad even if the library is kept.
