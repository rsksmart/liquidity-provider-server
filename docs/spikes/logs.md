# Logging library comparison and recommendation

## Decision

Adopt `log/slog` behind a thin LPS-owned `internal/logging` package.

The wrapper should expose the native slog key/value API plus a small set of helpers for LPS-specific gaps: initialization, error fields, fatal-equivalent logging, lazy values, trace propagation, PII redaction, and test capture.

Replace `gorilla/handlers.LoggingHandler` with structured request-logging middleware during the migration.

## Demo packages

The four runnable prototypes live under [spike/logging/](../../spike/logging/). Each is its own `main` package with the same per-deliverable file layout (see [Prototype deliverables](#prototype-deliverables) below), so the comparison axes below can be checked against working code:

| Library  | Demo                                         | Notable file                                                                              |
|----------|----------------------------------------------|-------------------------------------------------------------------------------------------|
| logrus   | [spike/logging/logrus/](../../spike/logging/logrus/)   | [pii.go](../../spike/logging/logrus/pii.go) — hook-based redaction                       |
| **slog** | [spike/logging/slog/](../../spike/logging/slog/)       | [pii.go](../../spike/logging/slog/pii.go) — `ReplaceAttr` redaction (the cleanest)        |
| zap      | [spike/logging/zap/](../../spike/logging/zap/)         | [pii.go](../../spike/logging/zap/pii.go) — `zapcore.Core` wrapper                         |
| zerolog  | [spike/logging/zerolog/](../../spike/logging/zerolog/) | [pii.go](../../spike/logging/zerolog/pii.go) — io.Writer regexp (awkward, deliberately)   |

Run any demo with:

```sh
LOG_FORMAT=json LOG_LEVEL=info go run ./spike/logging/slog
go test ./spike/logging/...
```

See [spike/logging/README.md](../../spike/logging/README.md) for what each demo exercises.

## Why this is the recommendation

The SPIKE did not find a hard blocker against the Rootstock Go guideline that already recommends `log/slog`. All four candidates can be made compliant with the Base Logging Standards, but they do not have the same migration cost or long-term shape.

`slog` is the best fit for LPS because it is standard library, aligns with the org guideline, has a clean handler-level redaction point, supports context-aware logging natively, has an OpenTelemetry bridge in the OpenTelemetry Go contrib project, and avoids introducing a new application-wide logging dependency during a migration that already touches about 325 call sites.

The main slog gaps are real but small and centralized:

- No built-in fatal level.
- No built-in stack trace field.
- Lazy evaluation is explicit, not automatic.
- `logfmt` needs either a small handler or a third-party handler if it is required beyond local text output.

Those gaps fit naturally in `internal/logging`. They do not justify choosing zap or zerolog for this codebase.

## Current LPS baseline

Measured on the current branch:

- Go files importing `github.com/sirupsen/logrus`: 80.
- Unstructured `log.Info`, `log.Warn`, `log.Error`, `log.Debug`, `log.Fatal`, `log.Trace`, `log.Panic`, `log.Print` calls: about 260.
- Formatted `log.<Level>f(...)` calls: about 65.
- Structured `log.WithField` or `log.WithFields` calls: 0.
- Test calls to `log.SetOutput(...)`: 5 calls across 4 test files.
- Production call to `log.SetOutput(file)`: 1, in `cmd/application/main.go`.
- Test calls to `log.SetLevel(log.DebugLevel)`: 32 calls across 7 test files.
- Logger initialization points: `cmd/application/main.go` and `internal/adapters/entrypoints/rest/server/server.go`.
- Go version in `go.mod`: `go 1.25.9`.

The conclusion: LPS does not meaningfully use logrus structured logging today. Almost every call site is either string concatenation or `Printf`-style formatting.

This matters for the decision. We cannot keep most call sites unchanged and become compliant. Every option requires rewriting the logging surface; the real comparison is which rewrite is most mechanical and leaves the smallest long-term API.

## Hard requirements from the task

The Base Logging Standards require every candidate demo to show:

- `LOG_FORMAT` supporting `json`, `logfmt`, and `otel` where supported.
- `LOG_LEVEL` configurable at startup.
- No plain unstructured output outside local development.
- Required fields on every log line: `timestamp`, `level`, `service`, `environment`, `version`, `traceId`, `message`.
- Inbound W3C `traceparent` extraction, with a generated trace if the header is absent.
- `traceId` carried through `context.Context` and emitted on every request-scoped log line.
- Outbound HTTP/RPC `traceparent` injection with a fresh span-id.
- Error logs containing `errorCode`, `errorMessage`, `errorStack`, and `errorContext`.
- PII drop or redaction by key for `privateKey`, `seed`, `mnemonic`, `apiKey`, `password`, and `authorization`.
- Either compatibility with `gorilla/handlers.LoggingHandler` through `io.Writer`, or a structured replacement.
- A test-output capture story for asserting on log content.

None of the four libraries fails these requirements outright. The difference is how much LPS code must be built around each one.

## Candidate assessment

### logrus

Best argument for it:

- It is already installed and wired into the app.
- `log.SetOutput` and `log.SetLevel` match the current test patterns.
- Hooks can mutate `Entry.Data`, which is a usable PII redaction point.
- `WriterLevel` already provides the `io.Writer` adapter used by `gorilla/handlers.LoggingHandler`.

Problems:

- It is in maintenance mode. The README explicitly says there will be no new features.
- The current LPS codebase uses none of its structured APIs.
- Context propagation is not first-class. `WithContext` stores context, but it does not extract trace data or attach fields by itself.
- Stack traces, error shape, W3C Trace Context, and format switching still require custom code.
- Keeping it contradicts the Go guideline unless the SPIKE finds a blocker in slog, zap, or zerolog. It did not.

Decision: do not choose logrus. It remains useful only as the baseline.

### log/slog

Best argument for it:

- It is standard library since Go 1.21 and requires no new dependency.
- It matches the Rootstock Go guideline direction.
- `InfoContext`, `WarnContext`, `ErrorContext`, and `LogAttrs` make context-aware logging part of the native API.
- `HandlerOptions.ReplaceAttr` is the cleanest PII interception point in the comparison.
- Handler wrapping is a natural place to add request `traceId` from `context.Context`.
- `otelslog` exists in `go.opentelemetry.io/contrib/bridges/otelslog`; paper check shows v0.18.0 in the OpenTelemetry Go contrib project.
- Test capture can be centralized with a buffer-backed JSON handler and a shared `slog.LevelVar`.

Problems:

- No native fatal method.
- No native stack trace field.
- Lazy evaluation requires explicit use of `slog.LogValuer`, `LogAttrs`, or an LPS helper.
- `logfmt` is not in the standard library.

Decision: choose slog with a thin wrapper. Its gaps are small, centralized, and mostly the same cross-cutting code LPS needs regardless of library.

### zap

Best argument for it:

- Strong production reputation and active maintenance.
- Excellent performance.
- Best native stack-trace story: `zap.Stack(...)`.
- Best test capture story: `zaptest/observer`.
- OTel bridge exists in `go.opentelemetry.io/contrib/bridges/otelzap`.
- Lazy evaluation is strong when using typed fields and `Check()`.

Problems:

- The typed `zap.Field` API makes the migration less mechanical. Each call site needs decisions such as `zap.String`, `zap.Stringer`, `zap.Int`, `zap.Any`, or a custom marshaler.
- `SugaredLogger` reduces migration friction but gives up part of the reason to choose zap.
- It introduces a new application-wide logging dependency.
- The org guideline rejects zap as the direct API.
- PII redaction is possible with `zapcore.WrapCore` and a custom encoder, but it is more code than slog's `ReplaceAttr`.

Decision: do not choose zap. It is technically strong, but its advantages do not offset the migration cost for LPS.

### zerolog

Best argument for it:

- Very fast and low-allocation.
- Fluent API gives good lazy-evaluation behavior.
- Small dependency footprint.
- Maintained; paper check shows recent v1.35.x versions.

Problems:

- Fluent chaining is the most different shape from logrus and slog. It makes the migration the least mechanical.
- PII redaction by arbitrary key is awkward because fields are built through a fluent event API, not a record map.
- OTel bridge maturity is weaker than slog and zap; it is community-driven rather than an OTel contrib bridge.
- The org guideline rejects zerolog as the direct API.

Decision: do not choose zerolog. Its runtime performance is not the limiting factor for LPS, and its migration/API shape is the weakest fit.

## Executive comparison matrix

This is a table in plain-text form to keep Cursor/VSCode markdown preview stable.

```text
Axis                                logrus             slog                zap                 zerolog
Format switch                       Medium             Good                Good                Medium
Migration from logrus               Poor if compliant  Best                Expensive           Expensive
Test capture                        Good               Good                Best                Medium
io.Writer adapter                   Best               Small wrapper       zapio.Writer        Small wrapper
Fatal semantics                     Built in           Wrapper             Built in            Built in
context.Context propagation         Weak               Best                Medium              Medium
W3C Trace Context                   DIY                DIY + clean handler DIY                 DIY
Base-field attachment               Medium             Best                Good                Good
Stack for errorStack                Weak               Wrapper             Best                Hook needed
Lazy field evaluation               Weak               Explicit helper     Good                Best
PII deny-list hook                  Good               Best                Medium              Weak
OTel bridge maturity                Weak               Good                Good                Weak
Dependency footprint                Existing external  Best                Medium              Good
Maintenance signal                  Weak               Best                Good                Medium
Org guideline alignment             Rejected           Recommended         Rejected            Rejected
Overall fit for LPS                 Do not choose      Choose              Do not choose       Do not choose
```

## Comparison across required axes

### Format switch ergonomics

Winner: `slog`, with a small caveat for `logfmt`.

- logrus: JSON and text are built in; `logfmt` needs a custom Formatter; OTel requires a third-party path.
- slog: JSON and text are built in; OTel has `otelslog`; `logfmt` needs a small handler or a third-party handler.
- zap: JSON and console are built in; `logfmt` needs a custom Encoder; OTel has `otelzap`.
- zerolog: JSON and console are built in; `logfmt` needs custom work; OTel bridge is community-only.

Recommendation for LPS: implement `json` now, support `otel` through `otelslog` in the demo, and keep `logfmt` as a small handler only if the standards review insists on it. For local development, slog text output is acceptable unless the standards owners require exact logfmt.

### Migration cost from logrus

Winner: `slog`.

- logrus has zero migration cost only if LPS accepts the current standards gaps, which the task explicitly does not.
- slog is the most mechanical rewrite: string logs become `slog.InfoContext(ctx, "message", "key", value)` or `slog.Info("message", ...)`.
- zap requires per-field type choices across hundreds of call sites.
- zerolog requires rewriting each call into fluent chains.

### Test capture ergonomics

Winner in isolation: zap. Winner for LPS migration: slog.

- zap's `zaptest/observer` is the strongest test API.
- slog can use a buffer-backed JSON handler and assertion helpers in `internal/logging/logtest`.
- logrus already works with `SetOutput`, but that keeps the current library.
- zerolog can write JSON to a buffer, but tests then parse JSON manually or through helpers.

LPS has 5 `SetOutput` calls and 32 level toggles. A slog test helper is enough; zap's nicer test API does not justify the heavier migration.

See: [zap `observer.New`](../../spike/logging/zap/capture_test.go) (typed records, no JSON parsing) vs. buffer-based capture in [slog](../../spike/logging/slog/capture_test.go), [logrus](../../spike/logging/logrus/capture_test.go), and [zerolog](../../spike/logging/zerolog/capture_test.go). The buffer approach is repetitive across libraries, but the assertion site stays similar — JSON-parse, then check fields.

### `io.Writer` adapter complexity

Winner if keeping `gorilla/handlers.LoggingHandler`: logrus.

- logrus already uses `WriterLevel`.
- slog needs a small `io.Writer` wrapper.
- zap has `zapio.Writer`.
- zerolog needs a small custom wrapper.

However, this axis should not decide the library. The existing `LoggingHandler` emits plain text and violates the standard. The correct LPS decision is to replace it with structured request middleware.

### Fatal semantics

Winner: logrus, zap, zerolog.

- logrus, zap, and zerolog have fatal APIs.
- slog intentionally does not.

For LPS this is not a blocker. Fatal usage is startup-oriented and can be centralized as `logging.Fatal(ctx, msg, attrs...)`, which logs and calls `os.Exit(1)`.

### `context.Context` propagation

Winner: `slog`.

- slog has native context-aware methods.
- logrus can store context but does not interpret it.
- zap and zerolog need helper packages or manual context/logger plumbing.

LPS needs request-scoped trace fields. Slog's handler model is the cleanest match.

See: [`slog`'s `contextHandler`](../../spike/logging/slog/utils.go) vs. the per-call-site [`baseEntry(ctx)`](../../spike/logging/logrus/utils.go) (logrus), [`loggerWithTrace(ctx)`](../../spike/logging/zap/utils.go) (zap), and [`loggerWithTrace(ctx)`](../../spike/logging/zerolog/utils.go) (zerolog) — only slog adds trace fields without the call site having to ask.

### W3C Trace Context support

No library wins outright.

Every candidate needs LPS middleware to parse inbound `traceparent`, generate a trace if absent, store trace state in `context.Context`, and inject outbound `traceparent` with a fresh span-id.

Slog still has the best fit after that middleware exists, because the handler can read `ctx` and add `traceId` to every record.

### Base-field auto-attachment

Winner: `slog`.

- slog supports base fields through `logger.With(...)` and can set the configured logger as default.
- logrus supports `WithFields`, but LPS would need to ensure the derived logger is used everywhere.
- zap supports base fields through `zap.Fields(...)`.
- zerolog supports base fields through `log.With().Str(...).Logger()`.

Slog's default logger plus handler attrs keeps call sites small and consistent with the Go guideline.

### Stack trace capture for `errorStack`

Winner: zap.

- zap has the best native stack support.
- slog needs an LPS helper around `runtime.Callers`.
- logrus needs custom code.
- zerolog needs a stack-marshaling hook.

For LPS, slog's gap is acceptable because error shape should be centralized anyway. `logging.Error` should be the only place that decides how `errorCode`, `errorMessage`, `errorContext`, and `errorStack` are emitted.

See: [`zap.Stack("errorStack")`](../../spike/logging/zap/errors.go) (one-liner native) vs. the hand-rolled `captureStack` in [slog](../../spike/logging/slog/errors.go), [logrus](../../spike/logging/logrus/errors.go), and [zerolog](../../spike/logging/zerolog/errors.go) — three near-identical 20-line helpers that fit naturally inside `internal/logging` for slog.

### Lazy evaluation of field values

Winner: zerolog and zap.

- zerolog's event chain avoids most work when the level is disabled.
- zap supports lazy behavior through typed fields, object marshalers, and `Check()`.
- slog supports lazy values through `slog.LogValuer`, but callers must use it.
- logrus does not help; arguments are evaluated before the call.

For LPS, the recommendation is to add `logging.Lazy(func() any) slog.LogValuer` and document it for expensive values such as transaction serialization. Most current log calls are cheap string/error values, so this does not justify choosing another library.

### PII deny-list interception

Winner: `slog`.

- slog has `HandlerOptions.ReplaceAttr`, which sees attributes before serialization and can recurse into groups.
- logrus hooks can mutate `Entry.Data`, which is also good.
- zap can wrap the core or encoder, but the implementation is more involved.
- zerolog hooks are less natural for arbitrary key-based redaction.

This is one of the strongest slog arguments for the Base Logging Standards.

See: [slog `piiRedactor`](../../spike/logging/slog/pii.go) (one function, recurses into groups) vs. [logrus `piiRedactHook`](../../spike/logging/logrus/pii.go) (hook mutating `Entry.Data`), [zap `piiCore`](../../spike/logging/zap/pii.go) (full `zapcore.Core` wrapper), and [zerolog `piiWriter`](../../spike/logging/zerolog/pii.go) (regex over the serialized JSON — workable but the weakest of the four).

### OTel bridge maturity

Winner: slog and zap.

- slog: `go.opentelemetry.io/contrib/bridges/otelslog`, maintained in the OTel Go contrib project.
- zap: `go.opentelemetry.io/contrib/bridges/otelzap`, also maintained in OTel Go contrib.
- logrus: community paths only.
- zerolog: community paths only.

The task only asks for Level-1 paper evaluation. Full OTel SDK adoption remains out of scope.

### Dependency footprint

Winner: `slog`.

- slog adds no logging dependency.
- logrus is already present but remains an external dependency in maintenance mode.
- zap adds `go.uber.org/zap` and transitive support packages such as `multierr`.
- zerolog adds `github.com/rs/zerolog`.

Binary-size deltas should be measured in the prototype branch if reviewers want exact numbers. For the recommendation, the important point is dependency shape: slog is already in the toolchain.

### Active maintenance signal

Winner: slog for risk profile; zap for third-party library maturity.

- slog is maintained as part of Go.
- zap is actively maintained by Uber and has a large production user base.
- zerolog is maintained and has recent v1.35.x versions, but the maintainer base is smaller.
- logrus is maintained for compatibility and fixes, but explicitly not for new features.

## Recommended LPS design

Create `internal/logging` with this limited surface:

- `Init(config Config)`: reads `LOG_FORMAT`, `LOG_LEVEL`, service name, environment, version, and output target; builds the handler stack; calls `slog.SetDefault`.
- `Error(ctx context.Context, msg string, err error, attrs ...slog.Attr)`: emits `errorCode`, `errorMessage`, `errorStack`, and `errorContext`.
- `Fatal(ctx context.Context, msg string, attrs ...slog.Attr)`: logs once, flushes if needed, then calls `os.Exit(1)`.
- `Lazy(func() any) slog.LogValuer`: canonical helper for expensive fields.
- `WithTrace(ctx context.Context, trace TraceContext) context.Context`: stores parsed trace state in context.
- `TraceMiddleware(next http.Handler) http.Handler`: extracts inbound `traceparent` or creates a trace.
- `TraceTransport(next http.RoundTripper) http.RoundTripper`: injects outbound `traceparent`.
- `logtest.NewBuffer(t)` and `logtest.SetLevel(t, level)`: replace direct `SetOutput` and `SetLevel` usage in tests.

The rest of the code should use slog's native API:

```go
slog.InfoContext(ctx, "quote accepted",
    "quoteHash", quoteHash,
    "operation", "acceptQuote",
)
```

Avoid building another fluent layer on top. The wrapper should enforce shared behavior at initialization and handler boundaries, not by creating a new logging language at every call site.

## Base Logging Standards compliance for the recommended option

### Format and output

`LOG_FORMAT`:

- `json`: production default, backed by `slog.NewJSONHandler`.
- `otel`: prototype through `otelslog.NewHandler`.
- `logfmt`: supported through a small handler only if the standards review requires exact logfmt output. Otherwise, use slog text for local development.

`LOG_LEVEL`:

- Parsed at startup into a `slog.LevelVar`.
- Test helper can temporarily set it for assertions.

No plain unstructured output:

- Replace `gorilla/handlers.LoggingHandler`.
- Remove string-only access logs.
- Migrate call sites to key/value fields.

### Required base fields

Attached at initialization:

- `service`
- `environment`
- `version`

Emitted by the handler:

- `timestamp`
- `level`
- `message`

Attached per request through context:

- `traceId`

### Distributed tracing

Inbound:

- Middleware reads `traceparent`.
- If missing or invalid, middleware generates a new trace.
- Trace state is stored in `context.Context`.
- Slog handler wrapper adds `traceId` to each record.

Outbound:

- HTTP/RPC transport wrapper reads trace state from `context.Context`.
- It injects `traceparent` with the same trace-id and a fresh span-id.
- Integration logs include `integrationName`, `integrationMethod`, `integrationDurationMs`, and `integrationStatusCode`.

### Error logs

`logging.Error` is responsible for the standard shape:

- `errorCode`: use case id or sentinel/domain code.
- `errorMessage`: root error message.
- `errorStack`: captured at the logging boundary.
- `errorContext`: wrapping chain and relevant attrs.

Fatal-equivalent logging reuses the same error field rules where an error exists.

### Sensitive data

Install PII redaction in `HandlerOptions.ReplaceAttr`.

The deny list is:

- `privateKey`
- `seed`
- `mnemonic`
- `apiKey`
- `password`
- `authorization`

Policy: redact by default rather than silently dropping, unless the standards owner requires dropping. Redaction keeps the operational signal that a field was attempted without leaking the value.

### LPS-specific access logs

Replace `gorilla/handlers.LoggingHandler` with middleware that emits one structured record per request:

- `httpMethod`
- `httpPath`
- `httpStatusCode`
- `durationMs`
- `userAgent`
- `requestId`

Because the trace middleware runs earlier, `traceId` is already available to the handler wrapper and does not need to be passed manually.

### Test capture

Replace direct logrus mutation in tests with `internal/logging/logtest`:

- Buffered JSON handler for assertions.
- Level override helper backed by `slog.LevelVar`.
- Helpers to assert message, level, and fields.

This is enough for the current test usage and avoids exposing logger internals across tests.

## Open tensions - resolved

### Tension 1: FluentLogger vs direct library API

Position: use direct slog-style key/value calls through a thin wrapper. Do not build a fluent builder.

Reasons:

- The Go guideline already points toward `slog.Default` from a shared logging package.
- Required base fields are better attached at initialization than repeated through a builder.
- PII rejection belongs in `ReplaceAttr`, not in call-site builders.
- Trace enrichment belongs in a handler that reads `context.Context`.
- Lazy evaluation is handled by `slog.LogValuer` plus `logging.Lazy`.
- A fluent builder would create a second logging API and make future migration to a shared org package harder.

Scope for other Go services:

- The principle should propagate: prefer native slog API plus a small shared wrapper.
- LPS should not define an org-wide builder pattern from this SPIKE.
- If Rootstock later ships `org/log`, LPS should be able to swap imports with minimal code churn.

### Tension 2: keep or replace `gorilla/handlers.LoggingHandler`

Position: replace it.

Reasons:

- The handler emits plain-text Apache-style access logs.
- The Base Logging Standards explicitly reject unstructured output outside local development.
- An `io.Writer` adapter would make the library compatible but would not make the output compliant.
- A structured middleware is small and directly produces the six required HTTP fields.
- `gorilla/handlers` is only used for this logging handler in the repo; removing that dependency does not affect `gorilla/mux`, `csrf`, `securecookie`, or `sessions`.

## Prototype deliverables

The four throwaway packages on this SPIKE branch live at:

```text
spike/logging/
  logrus/    -> spike/logging/logrus/
  slog/      -> spike/logging/slog/
  zap/       -> spike/logging/zap/
  zerolog/   -> spike/logging/zerolog/
```

Each package splits the standard into one file per deliverable so the
implementations diff cleanly. Within each file, **demo functions come
first** (what `main.go` calls), then the helpers they exercise — so a
reviewer sees the requirement before the machinery.

| File              | Purpose                                                                                                                                 | slog (recommended)                                                              |
|-------------------|-----------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------|
| `main.go`         | End-to-end orchestration — numbered checklist of the ten Base Logging Standards                                                       | [spike/logging/slog/main.go](../../spike/logging/slog/main.go)                  |
| `config.go`       | `demoLogFormatSupport`, `demoLogLevelConfiguredAtStartup`, `demoStructuredOutputOnly`, `demoRequiredFields`                             | [spike/logging/slog/config.go](../../spike/logging/slog/config.go)              |
| `utils.go`        | `Config`, env parsing, `initLogger`, package state, library-specific base-field / trace helper                                          | [spike/logging/slog/utils.go](../../spike/logging/slog/utils.go)                |
| `real_project_examples.go` | `demoRealLPSUseCaseLogging` (info path) and `demoRealLPSUseCaseError` (sentinel + `WrapUseCaseError` migrated to the structured error shape, with `errors.Is` still matching) | [spike/logging/slog/real_project_examples.go](../../spike/logging/slog/real_project_examples.go) |
| `tracing.go`      | Trace demos, `hit`, W3C plumbing, `traceMiddleware`, `traceTransport`, `sampleHandler`                                                  | [spike/logging/slog/tracing.go](../../spike/logging/slog/tracing.go)            |
| `errors.go`       | `demoErrorFields`, `demoFatal`, then `logError` / `logFatal` and stack capture                                                          | [spike/logging/slog/errors.go](../../spike/logging/slog/errors.go)              |
| `pii.go`          | `demoPIIRedaction`, then deny-list redaction (`privateKey`, `seed`, `mnemonic`, …)                                                      | [spike/logging/slog/pii.go](../../spike/logging/slog/pii.go)                    |
| `capture_test.go` | Test-output capture helper + assertions                                                                                                 | [spike/logging/slog/capture_test.go](../../spike/logging/slog/capture_test.go)  |

Each file inside a demo uses its corresponding library at every log site. In particular **`traceMiddleware` in `tracing.go`** is the structured request-logging middleware the SPIKE recommends as the gorilla `LoggingHandler` replacement — implemented four times so the trace-decision log and the six standard HTTP fields read in the library's native idiom. The only shared code is the W3C trace-context plumbing (parse / mint / context get-set / `statusRecorder`), which is data, not logging — its near-identical duplication across all four demos is deliberate, so side-by-side diffs of *logging* code are clean.

Each demo must produce comparable output for:

- Startup setup with `LOG_FORMAT` and `LOG_LEVEL`.
- Base fields auto-attached at initialization.
- Real LPS pegout-rebalance domain fields imported from `internal/entities` and `internal/usecases`, including the sentinel + `WrapUseCaseError` pattern as a structured error log with `errors.Is` still matching the sentinel after wrapping.
- Request info log with `httpMethod`, `httpPath`, `httpStatusCode`, `durationMs`, `userAgent`, and `requestId`.
- Outbound integration log with `integrationName`, `integrationMethod`, `integrationDurationMs`, and `integrationStatusCode`.
- Error log with `errorCode`, `errorMessage`, `errorStack`, and `errorContext`.
- PII attempt using `privateKey`.
- Fatal-equivalent logging.
- Test capture.

If a candidate cannot produce one item cleanly, the demo should state the reason in its package README.

## Migration effort and risk

### Recommended path: slog + wrapper

Estimated effort: one sprint, 1 PR.

Risk: low to medium.

Why:

- The migration is broad but mostly mechanical.
- The wrapper centralizes the non-stdlib behavior.
- No new logging dependency is introduced.
- The biggest behavior change is replacing access logging and converting unstructured strings to key/value fields.

### zap

Estimated effort: one sprint, 1 PR.

Risk: medium.

Why:

- Better stack and test tools.
- More per-call-site thought because of typed fields.
- New dependency and direct API not aligned with the guideline.

### zerolog

Estimated effort: one sprint, 1 PR.

Risk: medium-high.

Why:

- Fast runtime behavior.
- Largest API-shape change from current logrus calls.
- Awkward PII interception story.
- Weaker OTel bridge maturity.

### logrus

Estimated effort: 0 days if doing nothing; more if made compliant.

Risk: medium-high and increasing.

Why:

- Does not satisfy the standards in the current codebase.
- Library is in maintenance mode.
- Keeping it would require most of the same wrapper/middleware work while staying on a deprecated direction.

## Acceptance status

Covered by this document:

- Candidate comparison against all requested axes.
- Section-by-section compliance check for the recommended library.
- Recommendation with reasoning.
- Migration effort and risk per option.
- Explicit positions on FluentLogger vs direct slog API.
- Explicit position on replacing `gorilla/handlers.LoggingHandler`.
- Suggested follow-up PR split.

Delivered on this SPIKE branch:

- Four demo packages under [`spike/logging/`](../../spike/logging/) (logrus, slog, zap, zerolog).
- Each demo includes a project-aware example that imports LPS packages and logs representative pegout-rebalance fields, plus the sentinel + `WrapUseCaseError` pattern flowing into the structured error shape (`real_project_examples.go`).
- Each demo includes inbound trace extraction and outbound injection (`tracing.go`).
- Each demo includes error logging with all four required error fields and a stack (`errors.go`).
- Each demo includes the PII deny-list attempt (`pii.go`).
- Each demo includes a test-capture helper with three asserting tests (`capture_test.go`); all pass under `go test ./spike/logging/...`.

Still required before closing FLY-2308:

- Notion publication of this document.
- Review by at least one team member.

## References

- FLY-2308 task description.
- Rootstock Observability Standards draft: Base Logging Standards, Golang sub-page, FluentLogger Appendix.
- logrus README and releases: maintenance-mode notice, latest observed v1.9.4.
- Go `log/slog` package documentation.
- `go.opentelemetry.io/contrib/bridges/otelslog` package documentation.
- `go.opentelemetry.io/contrib/bridges/otelzap` package documentation.
- zap package documentation, including `zaptest/observer` and `zapio`.
- zerolog package documentation and releases.
- LPS code anchors: `cmd/application/main.go`, `internal/adapters/entrypoints/rest/server/server.go`, and `test/utils.go`.
