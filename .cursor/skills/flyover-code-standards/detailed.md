# Flyover Code Standards — Detailed Reference

Patterns with rationale and real before/after examples derived from code review
comments in the liquidity-provider-server repo.

---

## 1. Architecture & Layer Separation

This is the most frequently flagged category (~80 comments).

### Entity layer boundaries

The entity layer must not import infrastructure packages. If a function in
entities needs cryptographic hashing, receive the hash function as a parameter
instead of importing the Ethereum crypto package directly.

```go
// BAD — entity layer importing Ethereum crypto package
package liquidity_provider

import "github.com/ethereum/go-ethereum/crypto"

func VerifyHash(data []byte) []byte {
    return crypto.Keccak256(data)
}
```

```go
// GOOD — receive the hash function as a parameter
package liquidity_provider

type HashFunc func([]byte) []byte

func VerifyHash(data []byte, hash HashFunc) []byte {
    return hash(data)
}
```

> "you should avoid this, you're importing a function from ethereum package in
> the entity layer which shouldn't use libraries at all... to avoid this you
> should receive the hash function as a parameter" — PR #697

### Use case design

- Use cases expose only `Run()`. If you need additional public methods, create
  separate use cases.
- Use cases return domain objects. The handler is responsible for converting to
  DTOs using a converter function.
- If both the use case and repository need the same interface, define it at the
  entity layer so both can import it without circular dependencies.

```go
// BAD — use case returns a DTO
func (uc *GetAssetsUseCase) Run(ctx context.Context) (pkg.AssetsReportDTO, error) {
    // ...
}
```

```go
// GOOD — use case returns domain object; handler converts
func (uc *GetAssetsUseCase) Run(ctx context.Context) (reports.AssetsResult, error) {
    // ...
}

// In the handler:
func handleGetAssets(w http.ResponseWriter, r *http.Request) {
    result, err := useCase.Run(r.Context())
    dto := toAssetsDTO(result) // conversion happens here
    json.NewEncoder(w).Encode(dto)
}
```

> "don't return a DTO in the use case, return a domain object and then convert
> it to DTO in the handler" — PR #709

### Repository organization

- Each MongoDB collection gets its own file and repository struct.
- Functions that operate on different collections must not live in the same
  repository file.
- Create unique indexes where business rules require uniqueness (e.g., pegout
  deposits, trusted accounts).

### Configuration reads

Never read configurations directly from the database repository. Always go
through the `LiquidityProvider` interface, which internally calls the repo and
validates the signature. Without this, storing signatures adds no security value
because a compromised DB would serve unvalidated data.

### State initialization

If you have a use case dedicated to state initialization
(`InitializeStateConfigurationUseCase`), all initialization logic belongs there.
Don't scatter initialization across other use cases — it becomes hidden and hard
to track.

---

## 2. Function Design (Go)

### Multiple return values

Only use multiple returns when it's intuitive to the consumer what each value
represents (e.g., `(result, error)`). For anything beyond that, return a struct.

```go
// BAD — caller can't tell what each *big.Int means
func ListQuotesByDateRange(...) ([]Quote, []RetainedQuote, error) {
    // ...
}
```

```go
// GOOD — struct fields are self-documenting
type ListQuotesResult struct {
    Quotes         []Quote
    RetainedQuotes []RetainedQuote
}

func ListQuotesByDateRange(...) (ListQuotesResult, error) {
    // ...
}
```

> "Prefer structs to address both things, the multiple return is usually
> recommended only if it is intuitive to the consumer what is being returned"
> — PR #672

### Parameter reduction

- If a function takes 2+ fields from the same struct, pass the struct.
- If a function takes 4+ parameters, create a parameter struct.
- Functions that are methods of a use case struct should use the struct's fields
  directly instead of receiving them as parameters.

```go
// BAD — extracting two fields from the same struct
func transfer(btcTarget *big.Int, rbtcTarget *big.Int) { ... }

coldConfig := getColdWalletConfig()
transfer(coldConfig.BtcTarget, coldConfig.RbtcTarget)
```

```go
// GOOD — pass the struct
func transfer(config ColdWalletConfig) { ... }

transfer(getColdWalletConfig())
```

> "if you need to pass two fields to a function, and they are part of the same
> struct, you can just pass the struct, is cleaner and easier to maintain" — PR #901

### Flow control

Avoid `continue` in loops. It fragments the logic. Prefer putting the body
inside the condition instead.

```go
// BAD — continue fragments the logic
for _, quote := range quotes {
    if !quote.IsAccepted() {
        continue
    }
    if quote.IsExpired() {
        continue
    }
    results = append(results, quote)
}
```

```go
// GOOD — logic is contained in one block
for _, quote := range quotes {
    if quote.IsAccepted() && !quote.IsExpired() {
        results = append(results, quote)
    }
}
```

> "I normally recommend to avoid the continue because it makes harder to follow
> the logic. I rather have the append inside the if" — PR #776

Reduce nesting: use early returns, avoid nested `if/else` chains.
`errors.Is()` returns false on nil, so `if err != nil` before `errors.Is` is
redundant.

```go
// BAD — redundant nil check
if err != nil {
    if errors.Is(err, ErrNotFound) {
        return nil, err
    }
}
```

```go
// GOOD — errors.Is returns false on nil
if errors.Is(err, ErrNotFound) {
    return nil, err
} else if errors.Is(err, ErrConfigNotFound) {
    return nil, ErrNotFound
} else if err != nil {
    return nil, ErrTampered
}
```

> "errors.Is returns false if the err is nil, there is no need of checking
> first if error != nil, you can skip that to avoid the nesting" — PR #704

### Pointer discipline

- Prefer values over pointers. Pointers increase nil-pointer risk and may cause
  unnecessary heap allocation.
- When a function returns a pointer, always nil-check it at the call site — even
  if you "know" the function won't return nil. Defensive coding prevents panics
  in watchers/goroutines from crashing the server.
- For optional fields, prefer zero-value structs over nil pointers. Modify
  deserializers to replace nil with zero values.

### Constructors

Use constructor functions for struct creation. When a new field is added, you
update one place. With inline struct literals scattered across the codebase, you
must remember every instantiation site.

### Go idioms

Use time constants, not magic strings:

```go
// BAD
t, _ := time.Parse("2006-01-02", dateStr)
```

```go
// GOOD
t, _ := time.Parse(time.DateOnly, dateStr)
```

> "please use the constants to set the format... you're repeating the same
> string which doesn't makes sense unless you're familiar with how is
> implemented the standard lib internally" — PR #672

Use `chan struct{}` for signal-only channels:

```go
// BAD — implies true/false semantics
close chan bool

// GOOD — zero-size, clear intent
close chan struct{}
```

> "is preferable to use chan struct{} if you're using the channel only for the
> signal... the size of struct{} is zero... A chan bool might imply that true
> and false should be handled differently" — PR #745

---

## 3. Naming & Conventions

### JSON and BSON tags

- JSON tags: always camelCase (`json:"callFee"`)
- BSON tags: snake_case is acceptable (`bson:"call_fee"`)
- Never mix casing conventions within the same struct.

```go
// BAD
type TrustedAccount struct {
    Address string `json:"trusted_address" bson:"trusted_address"`
}
```

```go
// GOOD
type TrustedAccount struct {
    Address string `json:"trustedAddress" bson:"trusted_address"`
}
```

> "We use camelCase in this repo, please try to avoid snake_case... same for
> the json (in the case of bson snake_case is fine)" — PR #690

### Context-aware naming

Don't repeat the package/struct context in field or method names. In a `pegin`
package, use `CallForUserGasCost`, not `PeginCallForUserGasCost`.

```go
// BAD — "Pegin" is redundant inside a pegin package/struct
type RetainedPeginQuote struct {
    PeginCallForUserGasCost    *entities.Wei
    PeginRegisterPeginGasCost  *entities.Wei
}
```

```go
// GOOD — context is already provided by the type
type RetainedPeginQuote struct {
    CallForUserGasCost    *entities.Wei
    RegisterPeginGasCost  *entities.Wei
}
```

> "The fact this field is in this struct already tells you that is related to
> the pegin, so it could be called just CallForUserGasCost" — PR #730

### Use "Result" not "Response" for use case returns

"Response" implies HTTP/handler context. Use cases are domain layer.

```go
// BAD — implies HTTP handler context
type SummaryResponse struct { ... }
```

```go
// GOOD — neutral domain name
type SummaryResult struct { ... }
```

> "probably SummaryResult is a better name because by using Response you might
> imply that this is going to be returned in a handler, which shouldn't happen
> because its not a DTO" — PR #676

### CRUD naming

Use distinct verbs: "Add" for creation, "Update" for modification. Avoid
"Set" — it's ambiguous (could mean create or update).

### Interface naming

Give the "default name" to the interface: `AcceptQuoteUseCase` (interface),
`AcceptQuoteUseCaseImpl` (struct). Only implementations carry suffixes.

### File naming

Be consistent: don't mix `snakeCase` and `hyphen-case` in the same directory.

---

## 4. HTTP / API Design

### HTTP method semantics

- `GET`: no request body. Use query parameters.
- `DELETE`: no request body. Use path variables.
- Paths express **resources**, not actions. The HTTP verb conveys the action.

```go
// BAD — using request body in a GET endpoint
func GetReports(w http.ResponseWriter, r *http.Request) {
    var req ReportRequest
    json.NewDecoder(r.Body).Decode(&req)
}
```

```go
// GOOD — use query parameters
func GetReports(w http.ResponseWriter, r *http.Request) {
    startDate := r.URL.Query().Get("startDate")
    endDate := r.URL.Query().Get("endDate")
}
```

> "a GET request is supposed to not have body... You should use the query
> params for this" — PR #673

### Error responses

Return accurate HTTP status codes. Known errors (duplicate, not found) get
their proper code (409, 404). Don't return 500 for everything.

```go
// BAD — 500 for all errors
if err != nil {
    rest.WriteError(w, http.StatusInternalServerError, err.Error())
    return
}
```

```go
// GOOD — match the error to the code
if errors.Is(err, liquidity_provider.ErrNotFound) {
    rest.WriteError(w, http.StatusNotFound, err.Error())
} else if errors.Is(err, liquidity_provider.ErrDuplicate) {
    rest.WriteError(w, http.StatusConflict, err.Error())
} else {
    rest.WriteError(w, http.StatusInternalServerError, err.Error())
}
```

> "what about the case where the account doesn't exist? Shouldn't it return
> 404 instead of 500?" — PR #690

### DTOs

- Handlers return DTOs, never domain entities.
- Use `*big.Int` in DTOs, not `*entities.Wei` (which is a domain type).
- Document DTO fields when their names aren't self-explanatory.

### OpenApi.yml

The `deprecated` field is maintained manually because the codegen tool doesn't
support it. Never remove it during auto-regeneration.

---

## 5. Testing

### Assertion discipline

Assert the **error type**, not just that an error occurred.

```go
// BAD
require.Error(t, err)
```

```go
// GOOD
require.ErrorIs(t, err, ErrAmountTooLow)
```

> "can we assert the AmountTooLow error?" — PR #369

Assert **all fields** of result structs, not a subset. Use the zero-value
counter utility to detect new unasserted fields. In smart contract tests: assert
**events** (penalization, balance changes, call success/failure). Expected value
goes on the **left** in assertions.

### Mock discipline

Use **mockery** for generating mocks. Don't create manual mock structs.

```go
// BAD — hand-written mock struct
type mockSummaryUseCase struct {
    result *SummaryResult
    err    error
}

func (m *mockSummaryUseCase) Run(...) (*SummaryResult, error) {
    return m.result, m.err
}
```

```go
// GOOD — generated with mockery
mock := mocks.NewMockSummaryUseCase(t)
mock.On("Run", mock.Anything).Return(expectedResult, nil)
```

> "we have the mockery library to create the repos in a specific folder so
> this is not needed" — PR #676

#### Mockery setup in this repo

Configuration lives in `.mockery.yaml` at the repo root. Key conventions from that file:

- **Output directory**: `test/mocks/`
- **File name**: `{{ .InterfaceName | snakecase }}_mock.go` → e.g. `summary_use_case_mock.go`
- **Struct name**: `{{ .InterfaceName | firstUpper }}Mock` → e.g. `SummaryUseCaseMock`
- **Constructor**: `mocks.NewSummaryUseCaseMock(t)`

To regenerate all mocks after changing an interface:

```bash
mockery
```

Run from the repo root. Requires the pinned version to be installed:

```bash
go install github.com/vektra/mockery/v2@v2.53.1
```

When you add a **new interface** that needs a mock, register it in `.mockery.yaml` under its package before running mockery. Never add the interface to an arbitrary package entry — it must go under the package where the interface is defined.

Use `AssertNotCalled` for functions that shouldn't execute in error paths,
not just `AssertExpectations`.

```go
// BAD — only checks expectations were met
mock.AssertExpectations(t)
```

```go
// GOOD — explicitly verify functions that should NOT execute
mock.AssertNotCalled(t, "GetRbtcBalance")
mock.AssertExpectations(t)
```

> "instead of AssertExpectations you should use AssertNotCalled for the ones
> you didn't set up the mock because the function is supposed to return
> earlier" — PR #709

Avoid `Maybe()` in mock expectations when the call should always happen in
successful paths.

### Test design

- Test edge cases **separately**. Don't expire two conditions simultaneously —
  you won't detect if only one is implemented.
- Test with **real data** when possible (mainnet/testnet signatures, hashes,
  addresses).
- Split large tests into smaller ones rather than disabling `maintidx`. Disabling
  `funlen` in tests is acceptable.

### Test package naming

Test files must declare `package foo_test`, not `package foo`. This enforces
black-box testing: if a behavior can only be verified by reaching into unexported
identifiers, that is a signal the API surface needs improvement, not that the
test should be given internal access.

```go
// BAD — test lives inside the package, can access unexported symbols
package watcher

import "testing"

func TestSomething(t *testing.T) {
    s := internalHelper() // compiles only because it's in the same package
    ...
}
```

```go
// GOOD — test is an external consumer of the package
package watcher_test

import (
    "testing"
    "github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
)

func TestSomething(t *testing.T) {
    w := watcher.New(...)
    ...
}
```

If test helpers or shared setup are needed across test files in the same
directory, place them in a dedicated `common_test.go` file that also declares
`package foo_test`. Never use the production package name in test files.

### Coverage gaps

Every handler, converter, and public function should be tested. "Not tested" is
a common review comment — check coverage before submitting.

---

## 6. Security & Safety

### Management UI (JavaScript)

Use `textContent` instead of `innerHTML` to prevent XSS. If dynamic HTML is
unavoidable, sanitize all interpolated values.

```javascript
// BAD — XSS risk
cell.innerHTML = txId;
```

```javascript
// GOOD — safe text insertion
cell.textContent = txId;
```

> "I think it would be good to reduce (or eliminate when possible) the usage
> of innerHTML... prefer using textContent over innerHTML" — PR #914

Use `BigInt` instead of `Number` for large values (wei amounts overflow
`Number.MAX_SAFE_INTEGER`).

```javascript
// BAD — overflow risk with Number
const value = Number(weiAmount);
```

```javascript
// GOOD — safe for wei-scale values
const value = BigInt(weiAmount);
```

> "use bigint here to prevent overflow" — PR #690

### Git

Squash-merge PRs where sensitive data (mnemonics, keys) appeared in any
commit, even if later removed.

### Configuration integrity

- Always store the signed version of configurations and trusted accounts.
- Verify signature integrity on every read. Without verification, the signature
  field is security theater.
- Nil-check pointers in watchers and goroutines — a panic crashes the entire
  server.

---

## 7. Database Patterns (MongoDB)

### Query efficiency

Avoid N+1 patterns: fetch related data (quote + retained quote) in a single
aggregation, not N individual reads.

```go
// BAD — N reads to the DB
quotes := repo.ListQuotes(dateRange)
for _, q := range quotes {
    retained := repo.GetRetainedQuote(q.Hash) // N calls!
    // ...
}
```

```go
// GOOD — single aggregation returns both
result := repo.ListQuotesWithRetained(dateRange) // 1 call
for _, pair := range result.Pairs {
    quote := pair.Quote
    retained := pair.Retained
}
```

> "why are you refetching the quotes here? This is too expensive, we're
> talking about N reads to the DB where N is the dataset size" — PR #672

Use `singleflight` for expensive endpoints to deduplicate concurrent requests.

`cursor.All()` already closes the cursor — no need for `defer cursor.Close()`.

### Data handling

Decode MongoDB documents directly into Go structs. Don't roundtrip through
`bson.D` → marshal → unmarshal.

```go
// BAD — roundtrip through bson.D
var doc bson.D
cursor.Decode(&doc)
bytes, _ := bson.Marshal(doc)
bson.Unmarshal(bytes, &result)
```

```go
// GOOD — direct decode
var result StoredQuote
cursor.Decode(&result)
```

> "why are you marshalling and unmarshalling in the same function? why not
> decoding the doc directly into the struct?" — PR #676

Return stored structs from the repository. Don't reconstruct them from raw
documents when a direct decode works.

---

## 8. Code Organization & PR Hygiene

### Scope discipline

Keep PRs focused. Changes unrelated to the PR topic should go in separate PRs.

### Completeness

When populating a struct from external data, populate **all fields**, not just
the ones your current task needs. The function doesn't know how it'll be
consumed in the future.

### File placement

Scripts that aren't part of the build don't belong in `pkg/`, `cmd/`, or
`internal/`. Place them in `docker-compose/` or a top-level `scripts/` dir.

### Documentation

- Keep README badges (security team requirement). Don't remove them.
- When introducing new env vars, add them to documentation AND all
  docker-compose files.
- Specify units for configuration values (e.g., "5 means 5x the fee, not 5%").

### Linter rules

- Never disable `cyclop` or `maintidx` in production code. Compose smaller
  functions instead.
- `funlen` may be disabled in tests but `maintidx` should not — split into
  separate test functions.
- If a linter rule is disabled, it requires explicit justification in the review.

---

## 9. Log Message Centralization

Inline log strings scattered across the codebase cause test drift: the message
and its test assertion live in two places and silently diverge when one is
changed. The pattern below was adopted starting from PR #1024 (FLY-2388) and
refined in PR #1025.

### The rule

> Every log message lives in the **same package** that emits it, in a
> `log_messages.go` file. Log messages never go in `entities`.

Log strings are observability details, not domain facts. `entities` is reserved
for true domain contracts — specifically, alert subjects already live in
`entities/alerts` because the external alerting system depends on their exact
text. That exception does not extend to log messages.

The canonical filename is `log_messages.go` (renamed from `messages.go` per
review feedback in PR #1025).

### Constant vs typed function — the split

A message is either a **constant** or a **typed function**, decided by what the
caller must provide:

| Message shape | Form | Emitter |
|---|---|---|
| Fixed string, no placeholders (`"FooWatcher shut down"`) | constant | `log.Debug(LogFooShutdown)` |
| Error-only template, single trailing `%v` (`"error running check: %v"`) | constant | `log.Errorf(LogFooError, err)` |
| Any caller-provided value(s), with or without an error | **typed function** | `log.Info(LogFooChecking(txHash, quoteHash))` |

> "Distinguish the cases where you just prepend something to the error from the
> cases where you actually provide values... or even better a function so the
> values the caller needs to provide are typed:
> `LogPegoutRskGetDepositsError(start, end, err)`" — PR #1025 (Luisfc68)

Typed functions make the compiler enforce argument order, count, and types —
untyped `%s %d %v` templates silently accept anything. This applies to info and
debug templates too, not only error templates: a single-value template becomes
a one-parameter function for consistency.

```go
// internal/adapters/entrypoints/watcher/log_messages.go
package watcher

// Log message templates for pegout_rsk_watcher.go
const (
    LogPegoutRskPrefix      = "PegoutRskDepositWatcher: "
    LogPegoutRskChainHeight = LogPegoutRskPrefix + "error getting Rootstock chain height: %v"
    LogPegoutRskShutdown    = LogPegoutRskPrefix + "shut down"
)

func LogPegoutRskGetDepositsError(fromBlock, toBlock uint64, err error) string {
    return fmt.Sprintf(LogPegoutRskPrefix+"error executing getting deposits in range [%d, %d]: %v", fromBlock, toBlock, err)
}

func LogPegoutRskExpired(quoteHash string, expirationTime int64) string {
    return fmt.Sprintf(LogPegoutRskPrefix+"Quote %s expired at %d", quoteHash, expirationTime)
}
```

Type the parameters to match the call site exactly (`uint64` for block numbers,
`int64` for `time.Time.Unix()`, `uint` if the counter is `uint`, a concrete
struct for `%+v` values). Do not widen to `any` or `string` for convenience.

Do NOT convert the remaining constants to functions: prefixes must stay
constants (compile-time concatenation), no-placeholder strings gain nothing,
and error-only `%v` templates are already validated by `go vet`'s printf check
against the logrus `*f` variants.

### Structure

One `log_messages.go` per package. When a package has more than one source file
that emits logs, group per source file — one `const` block plus its functions,
under a comment naming that file. Do not group by ad-hoc themes; one group per
source file is the rule.

Message fragments that are concatenated at runtime (e.g. a reject reason built
from optional parts) are also typed functions; the emitting code concatenates
their return values:

```go
rejectReason := LogPegoutRskRejectReason(quoteHash)
if expired {
    rejectReason += LogPegoutRskRejectExpired(expirationTime, depositTime-expirationTime, depositTime)
}
log.Info(LogPegoutRskPrefix + rejectReason)
```

### Emitter

- Constant with format verbs → always the `*f` variant: `log.Errorf(LogFooError, err)`.
- Constant without verbs → non-`f` variant: `log.Debug(LogFooShutdown)`.
- Typed function → non-`f` variant with the call: `log.Error(LogFooUpdateError(quoteHash, err))`.

```go
// BAD — plain constant used with log.Error via concatenation
const LogTransferError = "...Error executing transfer to cold wallet: "
log.Error(LogTransferError, err) // relies on trailing space

// GOOD — format string used with Errorf
const LogTransferError = "...Error executing transfer to cold wallet: %v"
log.Errorf(LogTransferError, err)
```

### Test assertions

Tests derive the expected string from the **same constant or function** — never
re-type the string. Three rules, all from real review comments:

1. **Call the typed function directly** — no `fmt.Sprintf` needed:

```go
// BAD — re-types the string; drifts silently if the message changes
defer test.AssertLogContains(t, "Error executing transfer to cold wallet")()

// GOOD
defer test.AssertLogContains(t, watcher.LogTransferFailed("BTC", "transfer failed", transferError))()
```

2. **Assert the real error value, never an empty placeholder.** Formatting the
expected message with `""` makes the test pass even if the code stops logging
the error at all. When the use case wraps errors, assert the wrapped form:

```go
// BAD — passes even if the error is dropped from the log
checkFunc := test.AssertLogContains(t, fmt.Sprintf(watcher.LogPeginBtcUpdateExpiredError, test.AnyHash, ""))

// GOOD — asserts the exact error the use case produces
wrappedErr := usecases.WrapUseCaseError(usecases.ExpiredPeginQuoteId, assert.AnError)
checkFunc := test.AssertLogContains(t, watcher.LogPeginBtcUpdateExpiredError(test.AnyHash, wrappedErr))
```

3. **Reuse declared error variables** instead of recreating literals inline.

### Error wrapping: wrap once, at the use case boundary

Related rule surfaced while fixing the assertions above: `usecases.WrapUseCaseError`
must be applied exactly once, in the public `Run` method. Private helpers return
the raw error; wrapping in both places produces logs like
`"GetWatchedPegoutQuote: GetWatchedPegoutQuote: <err>"` and forces tests to
assert double-wrapped strings.

```go
// BAD — helper wraps, then Run wraps again
func (u *UseCase) getWatchedQuotes(...) (..., error) {
    if err != nil {
        return nil, usecases.WrapUseCaseError(usecases.GetWatchedPegoutQuoteId, err)
    }
}

// GOOD — helper returns raw error; Run is the single wrapping point
func (u *UseCase) getWatchedQuotes(...) (..., error) {
    if err != nil {
        return nil, err
    }
}
```

### Where log messages never go

Do not add log messages to any `entities` package. They always stay in the
package that emits them:

| Message example | Wrong location | Correct location |
|---|---|---|
| `"TransferColdWalletWatcher shut down"` | `entities/cold_wallet` | `adapters/entrypoints/watcher` |
| `"Server started at localhost:..."` | `entities/...` | `adapters/entrypoints/rest/server` |
| `"%s interaction with db: %+v"` | `entities/...` | `adapters/dataproviders/database/mongo` |
| `"Initializing application..."` | `entities/...` | `cmd/application` |

Alert subjects in `entities/alerts` are the single existing exception — they are
not log messages; they are external contracts consumed by the alerting system.

### Naming convention

Use a `Log` prefix for all log-message constants and functions to distinguish
them from alert subjects and error templates in the same codebase:

- Log messages → `LogTransferSuccess`, `LogPegoutRskExpired(quoteHash, ts)`
- Alert subjects → `AlertSubjectNodeReorg` (in `entities/alerts`)
- Error templates → `RskChainHeightErrorTemplate` (in `entities/blockchain`)

### Reference implementation

`internal/adapters/entrypoints/watcher/log_messages.go` is the pilot for the
repo-wide migration (FLY-2388). Follow it when centralizing the remaining
packages.
