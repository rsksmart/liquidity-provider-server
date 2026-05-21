# zap demo checklist

This package evaluates `go.uber.org/zap` for the FLY-2308 comparison in
[docs/spikes/logs.md](../../../docs/spikes/logs.md).

Run it with:

```sh
LOG_FORMAT=json LOG_LEVEL=debug go run ./spike/logging/zap
go test ./spike/logging/zap
```

## Requirement checklist

| Requirement | Evidence | Notes |
|-------------|----------|-------|
| `LOG_FORMAT` / `LOG_LEVEL` | `config.go`, `main.go` | JSON and console are built in. `logfmt` would need a custom encoder. `otel` falls back to JSON in the runnable demo; the document points to `otelzap` for real wiring. |
| Base fields | `initLogger` in `utils.go` | Uses `logger.With(...)` and typed `zap.Field` values. |
| Real LPS examples | `real_project_examples.go` | Imports `internal/entities` and `internal/usecases` to show pegout-rebalance logs with real domain values. Also covers the sentinel + `WrapUseCaseError` migration: useCase id → `errorCode`, wrapping chain → `errorMessage`, `errors.Is` against the sentinel still matches. |
| Request trace fields | `loggerWithTrace(ctx)` in `utils.go` | Works, but call sites must opt into the context-derived logger. |
| Inbound trace extraction | `tracing.go` (`demoInboundTraceExtraction`, `traceMiddleware`) | Parses W3C `traceparent`, mints a trace when absent, stores trace state in `context.Context`. |
| Structured access log | `tracing.go` (`traceMiddleware`) | Replaces `gorilla/handlers.LoggingHandler` with one structured record containing the six required HTTP fields. |
| Outbound trace injection | `tracing.go` (`demoOutboundTraceInjection`, `traceTransport`) | Injects the same trace-id and a fresh span-id into `traceparent`. |
| Error shape | `errors.go` | `demoErrorFields` / `demoFatal` at top; `zap.Stack` in `logError` is the strongest native stack story. |
| PII redaction | `pii.go` | `demoPIIRedaction` at top; `piiCore` wraps `zapcore.Core`. More involved than slog `ReplaceAttr`. |
| Fatal-equivalent | `errors.go`, `main.go` | Native `Fatal` exists and exits; demo path is guarded by `DEMO_FATAL_EXIT`. |
| Test capture | `capture_test.go` | Uses `zaptest/observer`, the strongest test capture API among the candidates. |

## Candidate caveats

- The typed field API makes migration from logrus less mechanical than slog.
- `SugaredLogger` would reduce migration friction, but also gives up part of
  the reason to choose zap.
- Direct zap usage does not align with the existing Go guideline preference for
  slog.
