# CI and local unit tests: surfacing failures (FLY-2282)

## Problem

The default `go test -v` stream is long. Failures appear inline; finding them often means searching the log (for example for `FAIL:`). The goal is a **short end-of-run recap** of failing tests (similar in spirit to Jest’s summary) **without** changing whether the run passes or fails.

## Chosen approach: gotestsum

We standardize on **[gotestsum](https://github.com/gotestyourself/gotestsum)** (`gotest.tools/gotestsum`), invoked via pinned `go run`:

- **End-of-run summary:** failed packages and tests are grouped at the bottom of the output.
- **Exit code:** `gotestsum` forwards `go test`’s exit code — CI still goes red when tests fail (Stage 03 regression proof).
- **No extra install in CI:** `go run …@v1.13.0` uses the module proxy like any other Go tool run; local devs get the same pin.
- **Coverage unchanged:** all `-race`, `-covermode`, `-coverpkg`, `-coverprofile` flags are passed **after** `--` so they apply to `go test` exactly as before; the existing Makefile filter for generated bindings still runs.

### Formatter

We use `--format standard-verbose` so output stays close to familiar `go test -v` style while still getting gotestsum’s summary.

## Tradeoffs

| Topic | Notes |
|-------|--------|
| **Cold `go run` cache** | First CI run after cache miss may download the gotestsum module; subsequent runs use build cache. |
| **Pin bumps** | Updating `@v1.13.0` is a deliberate change; record in PR when bumping. |
| **Integration tests** | This doc targets **`make test`** / PR **unit-tests** job (`./pkg/...`, `./internal/...`, `./cmd/...`). Integration tests under `test/integration/` stay separate; see `test/integration/Readme.md` and optional `go tool test2json` there. |

## Alternatives (not selected)

- **`go test -json` + custom script:** maximum control, higher maintenance.
- **GitHub Actions only:** wrapping only CI would diverge from `make test` / pre-commit (`make test`); we keep **one** entry point.
- **JUnit + `GITHUB_STEP_SUMMARY`:** useful later; not required for the first iteration.

## Where it is wired

- **`Makefile`** targets `test`, `coverage`, and `coverage-report` use `go run gotest.tools/gotestsum@v1.13.0 -- …`.
- **`.github/workflows/ci.yml`** still runs `make test` — no workflow duplication.

## Optional: install gotestsum on PATH

`make tools` also installs `gotestsum` so you can run `gotestsum` directly if you prefer; `make test` does not require it.
