# `spike/logging/` — FLY-2308 logging library demos

Throwaway code that supports the comparison in
[`docs/spikes/logs.md`](../../docs/spikes/logs.md). Each subdirectory is a
self-contained `main` package that exercises the same scenario against a
different logger, so the comparison axes in the SPIKE doc are backed by
runnable evidence.

> Status: SPIKE-only. Not wired into the application. Not for production
> use. Will be deleted after the migration tickets land per
> [`docs/spikes/logs.md`](../../docs/spikes/logs.md#suggested-migration-pr-split).

## Layout

```text
spike/logging/
  logrus/    baseline — what LPS uses today
  slog/      stdlib — the recommended target (see logs.md)
  zap/       Uber's library — strongest stack + test story
  zerolog/   rs/zerolog — fluent, low-allocation
```

Each demo is split into per-deliverable files so the comparison axes in
[`docs/spikes/logs.md`](../../docs/spikes/logs.md#prototype-deliverables)
can be diffed across libraries one axis at a time. Within each file,
**demo functions are at the top** (what `main.go` calls), then the helpers
they exercise:

| File              | What it shows                                                                                                      |
|-------------------|--------------------------------------------------------------------------------------------------------------------|
| `main.go`         | End-to-end orchestration: numbered checklist of the ten Base Logging Standards deliverables                        |
| `config.go`       | `demoLogFormatSupport`, `demoLogLevelConfiguredAtStartup`, `demoStructuredOutputOnly`, `demoRequiredFields`          |
| `utils.go`        | `Config`, `readConfig`, `initLogger`, env defaults, package state, library-specific base-field / trace helper      |
| `real_project_examples.go` | `demoRealLPSUseCaseLogging` (info path) and `demoRealLPSUseCaseError` (sentinel + `WrapUseCaseError` migrated to the structured error shape, with `errors.Is` still matching) |
| `tracing.go`      | Trace demos, `hit`, W3C plumbing, `traceMiddleware`, `traceTransport`, `sampleHandler`                             |
| `errors.go`       | `demoErrorFields`, `demoFatal`, then error shape (`errorCode`, …), `logError` / `logFatal`, stack capture          |
| `pii.go`          | `demoPIIRedaction`, then deny-list redaction (`privateKey`, `seed`, `mnemonic`, `apiKey`, `password`, `authorization`) |
| `capture_test.go` | Test-output capture helper + assertions over required fields                                                       |

The library-specific call style is what varies between demos, so each file
is a one-screen side-by-side comparison point. Every `.go` file starts with
a package comment (immediately before `package main`) describing what that
file demonstrates.

Each package also has a `README.md` with the same checklist, so reviewers can
check candidate-specific caveats without reading the whole implementation.

## Scenario contract

All four demos run the same scenario. This keeps the code useful for
side-by-side comparison rather than just proving that each library can print
one log line.

| Evidence file     | Requirement proved                                                                 |
|-------------------|-------------------------------------------------------------------------------------|
| `main.go`         | orchestration of all deliverables                                                   |
| `config.go`       | `LOG_FORMAT`, `LOG_LEVEL`, structured-only output, required base fields             |
| `utils.go`        | logger init and base-field attachment                                               |
| `real_project_examples.go` | project imports, real LPS-domain structured fields, and sentinel + `WrapUseCaseError` mapped to the structured error shape |
| `tracing.go`      | inbound `traceparent` extraction, structured access log, outbound injection         |
| `errors.go`       | `errorCode`, `errorMessage`, `errorStack`, `errorContext`, fatal-equivalent         |
| `pii.go`          | deny-list redaction for `privateKey`, `seed`, `mnemonic`, `apiKey`, `password`, `authorization` |
| `capture_test.go` | test-output capture strategy and assertions over required fields                    |

## Run a demo

```sh
LOG_FORMAT=json LOG_LEVEL=debug go run ./spike/logging/slog
LOG_FORMAT=text LOG_LEVEL=info  go run ./spike/logging/logrus
LOG_FORMAT=json LOG_LEVEL=info  go run ./spike/logging/zap
LOG_FORMAT=json LOG_LEVEL=info  go run ./spike/logging/zerolog
```

Each demo:

1. Initialises a logger with the requested format/level and the six required
   base fields (`service`, `environment`, `version`; `timestamp`, `level`,
   `message` come from the handler; `traceId` is attached per-request).
2. Logs representative LPS pegout-rebalance fields using real project
   packages (`internal/entities` and `internal/usecases`), then re-logs
   the same scenario as an error to show sentinel + `WrapUseCaseError`
   flowing into the structured error shape — `errors.Is` against the
   sentinel still matches after wrapping.
3. Spins up an inbound `httptest` server wrapped by the trace middleware,
   then sends one request with an existing `traceparent` and one without.
4. Calls an outbound `httptest` server through a `RoundTripper` that
   injects `traceparent` with a fresh span-id and logs the four
   integration fields.
5. Emits one error log with all four error fields including a stack.
6. Attempts to log a `privateKey` value and shows what the deny-list does.
7. Logs the "fatal" line (guarded — does not call `os.Exit` so the demo
   stays inspectable; the wrapper that *does* exit is shown in the file).

Run `go test ./spike/logging/...` to exercise the test-capture helpers.

## Demo limitations

- `LOG_FORMAT=otel` is paper-evaluated in this spike and falls back to JSON in
  runnable demos. Full OpenTelemetry SDK wiring is out of scope for FLY-2308.
- `LOG_FORMAT=logfmt` is approximate for some candidates. The comparison
  calls out where an exact handler or encoder would be needed.
- PII handling intentionally uses each library's natural interception point,
  even when that makes the candidate look worse. For example, zerolog uses a
  writer-level redactor to demonstrate why arbitrary key redaction is awkward.
- Fatal-equivalent logs are guarded by default so demos and tests remain
  inspectable. Set `DEMO_FATAL_EXIT=1` to exercise the actual exit path where
  a library or wrapper supports it.

## How the demos map to the comparison axes

The Comparison-axes section of
[`docs/spikes/logs.md`](../../docs/spikes/logs.md#comparison-across-required-axes)
links into the file in each demo that owns that axis, so a reviewer can
diff the four implementations of the same thing side-by-side.
