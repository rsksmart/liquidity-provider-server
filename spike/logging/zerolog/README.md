# zerolog demo checklist

This package evaluates `github.com/rs/zerolog` for the FLY-2308 comparison in
[docs/spikes/logs.md](../../../docs/spikes/logs.md).

Run it with:

```sh
LOG_FORMAT=json LOG_LEVEL=debug go run ./spike/logging/zerolog
go test ./spike/logging/zerolog
```

## Requirement checklist

| Requirement | Evidence | Notes |
|-------------|----------|-------|
| `LOG_FORMAT` / `LOG_LEVEL` | `config.go`, `main.go` | JSON and console output are built in. `logfmt` is approximated by console output. `otel` falls back to JSON because there is no first-party zerolog bridge. |
| Base fields | `initLogger` in `utils.go` | Uses `logger.With().Str(...).Logger()`. |
| Real LPS examples | `real_project_examples.go` | Imports `internal/entities` and `internal/usecases` to show pegout-rebalance logs with real domain values. Also covers the sentinel + `WrapUseCaseError` migration: useCase id → `errorCode`, wrapping chain → `errorMessage`, `errors.Is` against the sentinel still matches. |
| Request trace fields | `loggerWithTrace(ctx)` in `utils.go` | Works, but call sites must opt into the context-derived logger. |
| Inbound trace extraction | `tracing.go` (`demoInboundTraceExtraction`, `traceMiddleware`) | Parses W3C `traceparent`, mints a trace when absent, stores trace state in `context.Context`. |
| Structured access log | `tracing.go` (`traceMiddleware`) | Replaces `gorilla/handlers.LoggingHandler` with one structured record containing the six required HTTP fields. |
| Outbound trace injection | `tracing.go` (`demoOutboundTraceInjection`, `traceTransport`) | Injects the same trace-id and a fresh span-id into `traceparent`. |
| Error shape | `errors.go` | `demoErrorFields` / `demoFatal` at top; custom stack capture in `logError`. |
| PII redaction | `pii.go` | `demoPIIRedaction` at top; `piiWriter` does writer-level JSON rewriting — intentionally awkward. |
| Fatal-equivalent | `errors.go`, `main.go` | Native `Fatal` exists and exits; demo path is guarded by `DEMO_FATAL_EXIT`. |
| Test capture | `capture_test.go` | Buffer capture plus JSON-line parsing. |

## Candidate caveats

- Fluent chaining is the largest API-shape change from current logrus calls.
- Arbitrary PII key redaction is weaker than slog or logrus because zerolog
  does not expose a natural record-level map/hook.
- OTel bridge maturity is weaker than slog and zap for this task's paper
  evaluation.
