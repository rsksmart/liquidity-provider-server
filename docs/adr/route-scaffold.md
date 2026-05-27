# Design proposal: a route-registration scaffold for the LPS REST layer

**Status:** Draft — circulated for team review. Not a decision.
**Type:** Spike / design proposal (no production code in this ticket)
**Follow-up:** an implementation ticket will be filed only after the team converges on §5–§6.

---

## 1. Current-state map

### How a request flows today

Routing is owned by `gorilla/mux`. A single entrypoint configures everything:

- `internal/adapters/entrypoints/rest/routes/routes.go:42-57` — `ConfigureRoutes` wires CORS, builds a cookie store, calls `registerPublicRoutes`, conditionally calls `registerManagementRoutes`, and tacks on a global `OPTIONS` handler.
- `routes.go:59-67` — public route loop. For each `PublicEndpoint`, if `RequiresCaptcha` is true the handler is wrapped with the captcha middleware; otherwise it is mounted bare.
- `routes.go:70-87` — management route loop. Every route gets the CSRF middleware; routes not in `AllowedPaths` additionally get the session validator. `AllowedPaths` is the hand-maintained allowlist `{LoginPath, UiPath, StaticPath, IconPath}` declared at `management.go:13-20`.
- `routes.go:89-94` — `useMiddlewares`, a small left-to-right composer used by both registration helpers.

There is exactly one declarative metadata field on a route today — `PublicEndpoint.RequiresCaptcha` (`public.go:11-14`). Everything else is implied by which file the entry lives in (`public.go` vs `management.go`) or by the `AllowedPaths` allowlist. That `RequiresCaptcha` flag is the seed pattern this doc extends.

### Counts (verified from source on this branch)

| Bucket                   | File                                                                  | Count  |
|--------------------------|-----------------------------------------------------------------------|-------:|
| Public endpoints         | `internal/adapters/entrypoints/rest/routes/public.go:18-141`          |     17 |
| Management endpoints     | `internal/adapters/entrypoints/rest/routes/management.go:25-177`      |     27 |
| **Total mounted routes** |                                                                       | **44** |

*Note:* the original ticket survey reports "45 (18 + 27)". I count 17 public on `main`. The discrepancy may be `/metrics` (`public.go:134-140` — wired as `promhttp.Handler()` rather than a custom handler, so it is plausible the survey counted or excluded it differently). Worth a one-line reconcile before §6 lands.

OpenAPI coverage ("~38 of 45 documented") is taken from the ticket survey — not independently re-verified against `OpenApi.yml` for this draft. Whichever direction the true number lands, the qualitative claim — *the spec drifts because nothing enforces sync* — is what matters for §2.

### Existing helpers worth knowing about

The codebase already centralises five things that any scaffold would otherwise reinvent. Cited from `internal/adapters/entrypoints/rest/common.go`:

- `RequestValidator` (line 29) + custom validators registered in `registerValidations` (line 122) — the canonical validator chain.
- `DecodeRequest[T]` (line 197) — JSON decode with `DisallowUnknownFields`; writes the 400 response on failure.
- `ValidateRequest[T]` (line 234) — runs the validator and emits per-field error details via `getValidationMessage`.
- `JsonResponseWithBody[T]` / `JsonResponse` / `JsonErrorResponse` (lines 163-183) — the encoded response funnel; already generic over body type.
- `ErrorResponse` shape + `NewErrorResponseWithDetails` (lines 142-155) — the wire contract for errors.

There is also a partial error-mapping helper at `internal/adapters/entrypoints/rest/handlers/common.go:15-39` (`HandleAcceptQuoteError`) that collapses seven `errors.Is` branches into a single switch. **This is significant**: half of the "boilerplate" we want to eliminate has already been recognised and partially extracted in handler-local code. We are not building from scratch; we are finishing a refactor the team has already started.

### Test seam

The `EndpointFactory` interface (`routes.go:23-26`) exists explicitly so tests can mock the slice of endpoints. The interface is small (`GetPublic`, `GetPrivate`) and trivially preservable under any of the candidate designs in §3 — but it is worth naming up front because losing that mock surface would be a regression.

---

## 2. Pain-point inventory

### What a developer touches to add one endpoint today

Walking through `POST /pegin/acceptQuote` as the canonical example (handler at `internal/adapters/entrypoints/rest/handlers/accept_pegin_quote.go`):

1. **New handler file** under `internal/adapters/entrypoints/rest/handlers/` — a constructor returning `http.HandlerFunc` (90 files in that folder today, roughly half implementation and half tests).
2. **Decode / validate / respond glue.** Lines 26-33 of `accept_pegin_quote.go` show the inline pattern: `DecodeRequest` → early return, `ValidateRequest` → early return.
3. **Quote-hash-shape validation** (line 35) — domain-specific, does not fit the generic validator tags, has to be invoked by hand.
4. **Use-case invocation** (line 41) — single line, the part we *want* developers to be writing.
5. **Error-mapping ladder** (lines 42-62) — five `errors.Is` branches mapping domain errors to HTTP status codes + user-visible messages. Mostly duplicated across accept-quote handlers, hence the partial extraction in `handlers/common.go`.
6. **Response construction** (lines 64-68) — DTO assembly + `JsonResponseWithBody`.
7. **Route entry in `public.go`** (line 40-47) — six lines: path, method, handler constructor wire-up, `RequiresCaptcha: true`.
8. **OpenAPI block in `OpenApi.yml`** — hand-edited path entry plus any new schemas in `components/schemas`.
9. **OpenAPI doc-comment** above the handler constructor (lines 19-24) — `@Title`, `@Description`, `@Param`, `@Success`, `@Route`. *Nothing reads these.* They are informational; the codebase has no `swag` build step, no `oapi-codegen`, no validation that they match either the handler or `OpenApi.yml`. Today they document intent for humans, nothing more.
10. **DTO definitions** in `pkg/` with `json:` and `validate:` tags.
11. **Handler-level test** in `..._test.go` next to the handler.

Steps 2, 5, 6, 8, 9 are essentially the same shape on every endpoint. Step 8 + 9 in particular drift independently: a developer can change the handler's behaviour without touching `OpenApi.yml`, and CI will pass.

### Concretely, by the numbers

A representative non-trivial handler — `internal/adapters/entrypoints/rest/handlers/get_report_summaries.go` — is 82 lines. Of those:

- ~6 lines are use-case invocation, the irreducible code.
- ~12 lines are query-param wiring + domain validation, partially scaffoldable.
- ~25 lines are `singleflight` plumbing (genuinely endpoint-specific — *not* something to fold into a scaffold).
- The remaining ~40 lines are decode/validate/respond/error glue.

For accept-quote shaped handlers, the ratio is worse: ~6 lines of intent buried in ~45 lines of boilerplate (`accept_pegin_quote.go:25-70`).

---

## 3. Candidate approaches

Each sketch below shows roughly what "adding a new route" looks like under that pattern. They are **pseudocode** — not buildable Go — and chosen to be syntactically plausible so the team can argue about ergonomics rather than syntax.

### A. Descriptor struct + generic wrapper

Extends today's `Endpoint` struct rather than replacing it. Generics are already used in this codebase (`DecodeRequest[T]`, `JsonResponseWithBody[T]`), so the leap is small.

```go
type Route[Req, Resp any] struct {
    Method      string
    Path        string
    Summary     string                                  // -> OpenAPI summary
    Description string                                  // -> OpenAPI description
    UseCase     func(ctx context.Context, in Req) (Resp, error)
    Middlewares []Middleware                            // captcha, csrf, singleflight, ...
    ErrorMap    []ErrorMapping                          // {Is: err, Code: int, Msg: string}
}

func (r Route[Req, Resp]) Mount(router *mux.Router) {
    h := wrap(r)
    for _, mw := range r.Middlewares { h = mw(h) }
    router.Path(r.Path).Methods(r.Method).Handler(h)
}

func wrap[Req, Resp any](r Route[Req, Resp]) http.Handler { /* decode -> validate -> r.UseCase -> r.ErrorMap -> respond */ }

// Adding /pegin/acceptQuote becomes:
var AcceptPeginQuote = Route[pkg.AcceptQuoteRequest, pkg.AcceptPeginRespose]{
    Method: "POST", Path: "/pegin/acceptQuote",
    Summary: "Accept a peg-in quote and obtain a deposit address",
    UseCase: func(ctx context.Context, in pkg.AcceptQuoteRequest) (pkg.AcceptPeginRespose, error) {
        return useCases.AcceptPeginQuote().Run(ctx, in.QuoteHash, "")
    },
    Middlewares: []Middleware{Captcha},
    ErrorMap: []ErrorMapping{
        {Is: usecases.QuoteNotFoundError,    Code: 404, Msg: "quote not found"},
        {Is: usecases.ExpiredQuoteError,     Code: 410, Msg: "expired quote"},
        {Is: usecases.NoLiquidityError,      Code: 409, Msg: "not enough liquidity"},
        {Is: blockchain.ContractPausedError, Code: 503, Msg: "protocol is paused"},
    },
}
```

### B. Fluent builder

Builder API that finalises to a `Route` value internally (so the test seam is preserved).

```go
routes.POST("/pegin/acceptQuote").
    Summary("Accept a peg-in quote and obtain a deposit address").
    Body[pkg.AcceptQuoteRequest]().
    Returns[pkg.AcceptPeginRespose]().
    Captcha().
    MapError(usecases.QuoteNotFoundError,    404, "quote not found").
    MapError(usecases.ExpiredQuoteError,     410, "expired quote").
    MapError(usecases.NoLiquidityError,      409, "not enough liquidity").
    MapError(blockchain.ContractPausedError, 503, "protocol is paused").
    Handle(func(ctx context.Context, in pkg.AcceptQuoteRequest) (pkg.AcceptPeginRespose, error) {
        return useCases.AcceptPeginQuote().Run(ctx, in.QuoteHash, "")
    })

// And a trivial GET with no body:
routes.GET("/version").
    Summary("Server version info").
    Returns[pkg.ServerInfoDTO]().
    Handle(func(ctx context.Context, _ struct{}) (pkg.ServerInfoDTO, error) {
        info, err := useCases.ServerInfo().Run()
        return pkg.ToServerInfoDTO(info), err
    })
```

Tradeoff: reads nicer in line, but Go generics on chained method calls are awkward — `Body[T]()` and `Returns[T]()` need careful API design to avoid the builder collapsing to `any` mid-chain. Some libraries (e.g. `huma`) solve this with a top-level generic constructor instead of a fluent chain.

### C. Lightweight codegen

A `go:generate` step over a YAML (or annotated Go) descriptor that emits both the wiring code **and** the matching `OpenApi.yml` fragment. Single source of truth, but a toolchain commitment.

```yaml
# internal/adapters/entrypoints/rest/routes/routes.yaml
routes:
  - id: AcceptPeginQuote
    method: POST
    path: /pegin/acceptQuote
    summary: Accept a peg-in quote and obtain a deposit address
    requires_captcha: true
    request:  pkg.AcceptQuoteRequest
    response: pkg.AcceptPeginRespose
    use_case: AcceptPeginQuoteUseCase
    errors:
      - { is: usecases.QuoteNotFoundError,    code: 404, message: "quote not found" }
      - { is: usecases.ExpiredQuoteError,     code: 410, message: "expired quote" }
      - { is: usecases.NoLiquidityError,      code: 409, message: "not enough liquidity" }
      - { is: blockchain.ContractPausedError, code: 503, message: "protocol is paused" }
```

```go
//go:generate routegen -in ./routes.yaml -out ./routes_gen.go -openapi ../../../OpenApi.yml
// generated file mounts the route and stamps the matching path entry into OpenApi.yml
```

Tradeoff: solves the OpenAPI drift problem natively, but introduces a tool the team has to own, a `go generate` step in CI, and a less-obvious debugging path when something goes wrong (you are now reading generated code).

---

## 4. Decision matrix

| Axis                                  | A. Descriptor + generic wrapper                                                                            | B. Fluent builder                                                | C. Codegen                                                              |
|---------------------------------------|------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------|-------------------------------------------------------------------------|
| **Boilerplate reduction**             | High — collapses decode/validate/respond/error ladder                                                      | High — same, plus reads more like a DSL                          | Highest — descriptor is the only thing a dev writes                     |
| **Type safety / generics ergonomics** | Strong; `Req`/`Resp` inferred at the descriptor type itself                                                | Strong but fragile; chained `Body[T]/Returns[T]` is tricky in Go | Strong; codegen emits explicit typed wiring per route                   |
| **OpenAPI integration story**         | Walk descriptors at startup → emit / validate spec (both viable)                                           | Same — builder collects the same metadata                        | Native — gen owns the spec; drift is impossible by construction         |
| **Test seam (`EndpointFactory` mock)**| Preserved — `[]Route` is still a slice; factory unchanged                                                  | Preserved — builder finalises to a slice                         | Preserved — generated code outputs the same slice                       |
| **Migration cost**                    | Low — adapter from current `[]Endpoint` is trivial; incremental                                            | Medium — builder API needs to stabilise before mass-migrate      | High — toolchain authoring + CI wiring + team training                  |
| **Reversibility**                     | High — descriptors are plain data; any route can drop back to hand-rolled `http.Handler` as an escape hatch | Medium — DSL leaks into many call sites; backing out is grep-and-rewrite | Low — `go generate` in tree is sticky; removing it is disruptive       |
| **New runtime deps**                  | None                                                                                                       | None (or one small lib if we adopt e.g. `huma`)                  | One codegen tool we author or adopt                                     |
| **Effort to ship a viable v1**        | ~1 dev-week including 2–3 migrated routes as proof                                                         | ~1.5 dev-weeks (API design eats time)                            | ~2–3 dev-weeks (tool + tests + CI)                                      |
| **What goes wrong loudly?**           | Compile errors in descriptor field types                                                                   | Compile errors, but stack traces include builder internals       | "Why did `go generate` not run / not commit?" — silent until CI catches |

---

## 5. Open decisions for the team

These are genuinely open. Each one shifts the recommendation and the AC for the follow-up implementation ticket.

1. **OpenAPI: generate or validate?** Do we make the descriptor the single source of truth and *generate* `OpenApi.yml` (highest drift protection, biggest tooling commitment), or do we keep `OpenApi.yml` hand-maintained and add a CI step that *validates* it against the descriptors at PR time (lower commitment, surfaces drift instead of preventing it)? Validate-first is cheaper and reversible; generate-first is more correct but harder to undo.

2. **Per-route middleware: typed tags or explicit slice?** Today `RequiresCaptcha bool` is the only declarative middleware tag (`public.go:13`). Do we extend this with more booleans / an enum (`{Captcha, CSRF, SessionRequired, SingleFlight}`), or move to an explicit `Middlewares []Middleware` slice on the descriptor? The boolean approach is cleaner *if* the set is small and closed; the slice approach handles the long tail (e.g. `singleflight` on the three `/reports/*` endpoints, which is currently a constructor-time concern in `get_report_summaries.go:27-31`).

3. **Migration: big-bang or opportunistic?** Migrate all 44 routes in one PR (clean cutover, high review burden, one risky merge), or migrate opportunistically as endpoints get touched for unrelated work (mixed code-styles during the transition, indefinite tail, but zero risky single PR)? The descriptor-vs-hand-rolled coexistence is mechanically easy in approach A; it is harder in C.

4. **Router choice.** Do we stay on `gorilla/mux` (well-known, but the project has been in maintenance mode since 2022 and was briefly archived), or use this scaffold work as cover to evaluate alternatives — `chi`, `httprouter`, or the stdlib `http.ServeMux` improvements in Go 1.22+? The scaffold itself is router-agnostic; the question is whether we burn the migration budget on this at the same time, or split it.

5. **Error mapping ownership.** Where does "domain error → HTTP status + user-visible message" live? Today it is per-handler ladders (`accept_pegin_quote.go:42-62`) with a partial extraction in `handlers/common.go:15-39`. Options: (a) on the route descriptor (per §3 sketches), (b) a router-wide error translator with a global registry, (c) a typed `httpcontract` interface on use-case errors so the mapping is owned by the use-case package. Option (c) is the most architecturally clean but the biggest behavioural change.

6. **Handler doc-comments (`@Title`, `@Route`).** Once descriptors carry the same metadata, do the doc-comments become noise (delete them) or do we keep them as fallback documentation? If we adopt `swag`/`oapi-codegen` style, that decides this for us.

---

## 6. Recommendation (proposal — not a decision)

**My pick: Approach A (descriptor + generic wrapper), shipped incrementally, with OpenAPI handled by *validation in CI* first and *generation* deferred to a follow-on ticket once the descriptor shape has settled.**

The reasoning, briefly: A is the smallest step that solves the actual pain points enumerated in §2 — decode/validate/respond/error-map collapse into wrapper code, the existing `Endpoint`/`PublicEndpoint`/`RequiresCaptcha` shape extends naturally so the migration is mechanically incremental, the `EndpointFactory` test seam is preserved verbatim, and any single route can revert to a hand-rolled `http.HandlerFunc` if we hit an edge case (e.g. the `singleflight`-wrapped report endpoints). Starting with *validation* rather than *generation* of `OpenApi.yml` gives us the same drift-detection win for ~10% of the implementation cost, and keeps the door open to generation later if we decide validation is too coarse. The fluent builder (B) is appealing on a screenshot but harder to get right in Go generics today, and codegen (C) is the wrong commitment to make before we have lived with the descriptor shape for a release or two. If the team disagrees, the matrix in §4 should make clear which axis dominates that disagreement.

---

## Appendix A — Worked example: `/pegin/acceptQuote` under Approach A

**Before** (today). The relevant pieces are spread across three locations totalling ~55 lines of glue:

```go
// internal/adapters/entrypoints/rest/handlers/accept_pegin_quote.go (lines 25-70, abridged)
func NewAcceptPeginQuoteHandler(useCase AcceptQuoteUseCase) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        var err error
        acceptRequest := pkg.AcceptQuoteRequest{}
        if err = rest.DecodeRequest(w, req, &acceptRequest); err != nil { return }
        if err = rest.ValidateRequest(w, &acceptRequest); err != nil { return }

        if err = quote.ValidateQuoteHash(acceptRequest.QuoteHash); err != nil {
            jsonErr := rest.NewErrorResponseWithDetails("invalid quote hash", rest.DetailsFromError(err), true)
            rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
            return
        }

        acceptedQuote, err := useCase.Run(req.Context(), acceptRequest.QuoteHash, "")
        if errors.Is(err, usecases.QuoteNotFoundError)    { /* 404 */ ; return }
        if errors.Is(err, usecases.ExpiredQuoteError)     { /* 410 */ ; return }
        if errors.Is(err, usecases.NoLiquidityError)      { /* 409 */ ; return }
        if errors.Is(err, blockchain.ContractPausedError) { /* 503 */ ; return }
        if err != nil                                     { /* 500 */ ; return }

        response := pkg.AcceptPeginRespose{
            Signature:                 acceptedQuote.Signature,
            BitcoinDepositAddressHash: acceptedQuote.DepositAddress,
        }
        rest.JsonResponseWithBody(w, http.StatusOK, &response)
    }
}

// internal/adapters/entrypoints/rest/routes/public.go (lines 40-47)
{
    Endpoint: Endpoint{
        Path:    "/pegin/acceptQuote",
        Method:  http.MethodPost,
        Handler: handlers.NewAcceptPeginQuoteHandler(useCaseRegistry.GetAcceptPeginQuoteUseCase()),
    },
    RequiresCaptcha: true,
},

// OpenApi.yml — hand-edited path entry (omitted)
```

**After** (descriptor approach). One descriptor, one place to look:

```go
var AcceptPeginQuote = Route[pkg.AcceptQuoteRequest, pkg.AcceptPeginRespose]{
    Method: "POST", Path: "/pegin/acceptQuote",
    Summary: "Accept a peg-in quote and obtain a deposit address",
    Middlewares: []Middleware{Captcha},
    PreValidate: func(in pkg.AcceptQuoteRequest) error {
        return quote.ValidateQuoteHash(in.QuoteHash) // domain-specific check
    },
    UseCase: func(ctx context.Context, in pkg.AcceptQuoteRequest) (pkg.AcceptPeginRespose, error) {
        accepted, err := useCases.AcceptPeginQuote().Run(ctx, in.QuoteHash, "")
        if err != nil { return pkg.AcceptPeginRespose{}, err }
        return pkg.AcceptPeginRespose{
            Signature:                 accepted.Signature,
            BitcoinDepositAddressHash: accepted.DepositAddress,
        }, nil
    },
    ErrorMap: []ErrorMapping{
        {Is: usecases.QuoteNotFoundError,    Code: 404, Msg: "quote not found"},
        {Is: usecases.ExpiredQuoteError,     Code: 410, Msg: "expired quote"},
        {Is: usecases.NoLiquidityError,      Code: 409, Msg: "not enough liquidity"},
        {Is: blockchain.ContractPausedError, Code: 503, Msg: "protocol is paused"},
    },
}
```

OpenAPI: a CI step walks every `Route[Req, Resp]` value, derives `(method, path, request schema, response schema, error codes)`, and asserts `OpenApi.yml` matches. PRs that change the descriptor without updating the spec fail at CI time instead of merging silently.

---

## Appendix B — Out of scope for this proposal

- Picking the final approach. The doc proposes; the team decides.
- Migrating any routes.
- Auth/CSRF/session refactors beyond what middleware-as-data implies.
- Replacing `gorilla/mux` (that is question 4 in §5, deliberately surfaced as a *decision*, not a *bundled change*).
- DTO / `validator` tag conventions in `pkg/` — out of scope; the scaffold consumes whatever DTOs we already have.
- Query-parameter routes (e.g. `/reports/summaries?startDate=&endDate=`). Pseudocode here covers JSON-body and no-body routes; query-param decoding is a v2 concern.

---

## References (all on this branch)

- `internal/adapters/entrypoints/rest/routes/routes.go:16-94` — `Endpoint`, `EndpointFactory`, `ConfigureRoutes`, public/management register helpers.
- `internal/adapters/entrypoints/rest/routes/public.go:11-141` — `PublicEndpoint` (with `RequiresCaptcha`), 17 public route entries.
- `internal/adapters/entrypoints/rest/routes/management.go:13-178` — `AllowedPaths`, 27 management route entries.
- `internal/adapters/entrypoints/rest/common.go:29, 122, 163-206, 234-250` — validator registration, `DecodeRequest[T]`, `ValidateRequest[T]`, `JsonResponseWithBody[T]`, error-response funnel.
- `internal/adapters/entrypoints/rest/handlers/accept_pegin_quote.go:19-70` — canonical handler shape, the `@Title/@Route` doc-comment style, and the five-branch error ladder.
- `internal/adapters/entrypoints/rest/handlers/common.go:15-39` — `HandleAcceptQuoteError`, evidence of an already-started consolidation effort.
- `internal/adapters/entrypoints/rest/handlers/get_report_summaries.go:27-82` — a heavier handler with `singleflight` plumbing; the scaffold must not bulldoze this case.
- `internal/adapters/entrypoints/rest/registry/registry.go:11-54` — `UseCaseRegistry` interface; the surface any scaffold would call into.
- `OpenApi.yml` — repo root; hand-maintained spec referenced in §2 step 8.
