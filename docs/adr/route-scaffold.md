# Route registration scaffold — design proposal

Status: draft

## 1. Current state

Routes are registered with `gorilla/mux`. The entrypoint is `ConfigureRoutes` in `internal/adapters/entrypoints/rest/routes/routes.go:42-57`. It calls two helpers:

- `registerPublicRoutes` (lines 59-67) iterates `[]PublicEndpoint`. If `RequiresCaptcha` is true, the handler is wrapped in the captcha middleware.
- `registerManagementRoutes` (lines 70-87) iterates `[]Endpoint`. Every route gets CSRF. Routes outside the `AllowedPaths` allowlist (`management.go:13-20`) also get the session validator.

Counts on `main`:

- Public: 17 (`public.go:18-141`)
- Management: 27 (`management.go:25-177`)
- Total: 44

The ticket survey says 45. The difference is probably `/metrics`, wired as `promhttp.Handler()` at `public.go:134-140`.

OpenAPI: ticket reports ~38 of 45 paths documented in `OpenApi.yml`. I didn't recount. The point isn't the exact number, it's that the spec drifts because nothing enforces it.

The only declarative routing metadata today is `PublicEndpoint.RequiresCaptcha` (`public.go:11-14`). Everything else — middleware choice, OpenAPI presence, error mapping — is implicit in which file the route entry lives in.

Things already centralized in `internal/adapters/entrypoints/rest/common.go`:

- `RequestValidator` + custom validators (lines 29, 122)
- `DecodeRequest[T]` (197) — disallows unknown fields
- `ValidateRequest[T]` (234) — emits per-field error details
- `JsonResponseWithBody[T]`, `JsonResponse`, `JsonErrorResponse` (163-183)
- `ErrorResponse` + `NewErrorResponseWithDetails` (142-155)

And in `internal/adapters/entrypoints/rest/handlers/common.go:15-39`: `HandleAcceptQuoteError` already collapses seven `errors.Is` branches into one switch. Half of what a scaffold would extract has already been pulled out; we'd be finishing that work, not starting it.

Test seam: `EndpointFactory` (`routes.go:23-26`) is the interface tests mock. Any scaffold has to keep that working.

## 2. What it takes to add an endpoint today

Walking through `POST /pegin/acceptQuote` (`handlers/accept_pegin_quote.go`):

1. New handler file in `rest/handlers/` with a constructor returning `http.HandlerFunc`.
2. Inline `DecodeRequest` + `ValidateRequest` with early returns (lines 26-33).
3. Domain-specific check (`quote.ValidateQuoteHash`, line 35) — doesn't fit validator tags.
4. Use case call (line 41).
5. Five-branch `errors.Is` ladder mapping domain errors to HTTP codes (lines 42-62).
6. Response DTO assembly + `JsonResponseWithBody`.
7. Entry in `routes/public.go`.
8. Hand-edit `OpenApi.yml`.
9. `@Title`/`@Route` doc comment above the handler. Nothing reads these. No `swag`, no `oapi-codegen`. They go stale.
10. DTOs in `pkg/` with json + validate tags.
11. Handler test.

Steps 2, 5, 6, 8, 9 look the same on every endpoint. 8 and 9 drift silently: you can change a handler without touching `OpenApi.yml` and CI passes.

Sizes: `accept_pegin_quote.go` is 70 lines, of which one calls the use case. `get_report_summaries.go` is 82 lines, of which ~25 are singleflight plumbing that a scaffold should leave alone.

## 3. Three approaches

Pseudocode only. None of this compiles.

### A. Descriptor struct + generic wrapper

```go
type Route[Req, Resp any] struct {
    Method, Path, Summary string
    UseCase     func(ctx context.Context, in Req) (Resp, error)
    Middlewares []Middleware
    PreValidate func(Req) error      // optional, domain-specific
    ErrorMap    []ErrorMapping       // {Is, Code, Message}
}

var AcceptPeginQuote = Route[pkg.AcceptQuoteRequest, pkg.AcceptPeginRespose]{
    Method: "POST", Path: "/pegin/acceptQuote",
    Summary: "Accept a peg-in quote",
    Middlewares: []Middleware{Captcha},
    PreValidate: func(in pkg.AcceptQuoteRequest) error {
        return quote.ValidateQuoteHash(in.QuoteHash)
    },
    UseCase: func(ctx context.Context, in pkg.AcceptQuoteRequest) (pkg.AcceptPeginRespose, error) {
        return useCases.AcceptPeginQuote().Run(ctx, in.QuoteHash, "")
    },
    ErrorMap: []ErrorMapping{
        {Is: usecases.QuoteNotFoundError,    Code: 404, Msg: "quote not found"},
        {Is: usecases.ExpiredQuoteError,     Code: 410, Msg: "expired quote"},
        {Is: usecases.NoLiquidityError,      Code: 409, Msg: "not enough liquidity"},
        {Is: blockchain.ContractPausedError, Code: 503, Msg: "protocol is paused"},
    },
}
```

### B. Fluent builder

```go
routes.POST("/pegin/acceptQuote").
    Summary("Accept a peg-in quote").
    Body[pkg.AcceptQuoteRequest]().
    Returns[pkg.AcceptPeginRespose]().
    Captcha().
    MapError(usecases.QuoteNotFoundError, 404, "quote not found").
    MapError(usecases.ExpiredQuoteError,  410, "expired quote").
    Handle(func(ctx context.Context, in pkg.AcceptQuoteRequest) (pkg.AcceptPeginRespose, error) {
        return useCases.AcceptPeginQuote().Run(ctx, in.QuoteHash, "")
    })
```

Reads nicer per line. The issue is `Body[T]()` and `Returns[T]()` mid-chain. Go generics don't compose cleanly through fluent builders. You usually end up with `any` and a runtime type assertion somewhere, or with awkward top-level constructor functions like `huma` uses.

### C. Codegen

```yaml
# routes.yaml
- id: AcceptPeginQuote
  method: POST
  path: /pegin/acceptQuote
  requires_captcha: true
  request:  pkg.AcceptQuoteRequest
  response: pkg.AcceptPeginRespose
  errors:
    - { is: usecases.QuoteNotFoundError, code: 404, message: "quote not found" }
    - { is: usecases.ExpiredQuoteError,  code: 410, message: "expired quote" }
```

`go:generate` over the YAML emits the wiring code and the matching `OpenApi.yml` fragment. Single source of truth. But now we own a generator.

## 4. Comparison

|                       | A: descriptor                       | B: builder              | C: codegen                       |
|-----------------------|-------------------------------------|-------------------------|----------------------------------|
| Boilerplate cut       | high                                | high                    | highest                          |
| Generics ergonomics   | good                                | fragile                 | good (gen emits explicit types)  |
| OpenAPI sync          | validate or generate from descriptors | same                  | native; drift impossible         |
| `EndpointFactory` mock| preserved                           | preserved               | preserved                        |
| Migration             | incremental, low risk               | medium                  | one big lift or all-or-nothing   |
| Reversibility         | high (descriptors are data)         | medium (DSL spreads)    | low (gen step is sticky)         |
| New runtime deps      | none                                | none, or one small lib  | one toolchain we own             |
| Time to a viable v1   | ~1 week                             | ~1.5 weeks              | ~2-3 weeks                       |

## 5. Open decisions

1. **OpenAPI: generate or validate?** Generate is correct-by-construction but commits us to a code generator. Validate-in-CI is cheap, catches drift at PR time, and doesn't prevent it.
2. **Middleware on the descriptor: bools or slice?** Today there's one bool, `RequiresCaptcha`. Adding more bools (`CSRF`, `Session`, `SingleFlight`) is fine if the set stays small. A `Middlewares []Middleware` slice handles the long tail — e.g. the three `/reports/*` endpoints that wrap singleflight at construction time (`get_report_summaries.go:27-31`).
3. **Migration: big-bang or rolling?** All 44 in one PR is one risky merge. Rolling means mixed styles for a while. Approach A coexists with the current `[]Endpoint` cleanly; C doesn't.
4. **Keep `gorilla/mux`?** It's been in maintenance mode since 2022. Swap to `chi`, `httprouter`, or stdlib `http.ServeMux` (1.22+) while we're here, or out of scope? The scaffold itself is router-agnostic.
5. **Where do error-to-HTTP mappings live?** Per-route descriptor (as sketched), a global registry, or a typed interface on use-case errors so the mapping lives in the use-case package. The last is the cleanest, also the biggest change.
6. **Keep the `@Title`/`@Route` doc comments?** If descriptors carry the same metadata, the comments are duplicate. Delete or keep as a redundant human-readable copy?

## 6. Recommendation

Approach A, rolling migration, OpenAPI handled by CI validation first. Defer codegen until we've used the descriptor shape for a release or two.

A is the smallest change that fixes what hurts in §2. Decode, validate, response, and error mapping move into the wrapper. The existing `Endpoint`/`RequiresCaptcha` shape extends naturally, so routes can migrate one at a time without a flag day. Anything that doesn't fit (e.g. the singleflight report handlers) can stay hand-rolled; the scaffold and the original `[]Endpoint` coexist. `EndpointFactory` doesn't change.

Validation in CI gets most of the drift-protection win for a small fraction of the cost of generation. If we later want generation, the descriptor shape we settle on now becomes its input.

B looks nice in screenshots and is harder to get right in Go generics today. C is plausibly the right end state, but it's the wrong first step — we'd be locking in a descriptor shape before we've felt the rough edges.

If you disagree, the matrix is the place to push back. Tell me which row is wrong.

## Appendix — `/pegin/acceptQuote` before and after

Today, in `handlers/accept_pegin_quote.go:25-70` (abridged) plus the `public.go` entry plus the `OpenApi.yml` block:

```go
func NewAcceptPeginQuoteHandler(useCase AcceptQuoteUseCase) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        acceptRequest := pkg.AcceptQuoteRequest{}
        if err := rest.DecodeRequest(w, req, &acceptRequest); err != nil { return }
        if err := rest.ValidateRequest(w, &acceptRequest); err != nil { return }

        if err := quote.ValidateQuoteHash(acceptRequest.QuoteHash); err != nil {
            // 400
            return
        }

        acceptedQuote, err := useCase.Run(req.Context(), acceptRequest.QuoteHash, "")
        if errors.Is(err, usecases.QuoteNotFoundError)    { /* 404 */ ; return }
        if errors.Is(err, usecases.ExpiredQuoteError)     { /* 410 */ ; return }
        if errors.Is(err, usecases.NoLiquidityError)      { /* 409 */ ; return }
        if errors.Is(err, blockchain.ContractPausedError) { /* 503 */ ; return }
        if err != nil                                     { /* 500 */ ; return }

        rest.JsonResponseWithBody(w, http.StatusOK, &pkg.AcceptPeginRespose{
            Signature: acceptedQuote.Signature,
            BitcoinDepositAddressHash: acceptedQuote.DepositAddress,
        })
    }
}
```

Under A — one descriptor, one place to look:

```go
var AcceptPeginQuote = Route[pkg.AcceptQuoteRequest, pkg.AcceptPeginRespose]{
    Method: "POST", Path: "/pegin/acceptQuote",
    Summary: "Accept a peg-in quote",
    Middlewares: []Middleware{Captcha},
    PreValidate: func(in pkg.AcceptQuoteRequest) error {
        return quote.ValidateQuoteHash(in.QuoteHash)
    },
    UseCase: func(ctx context.Context, in pkg.AcceptQuoteRequest) (pkg.AcceptPeginRespose, error) {
        accepted, err := useCases.AcceptPeginQuote().Run(ctx, in.QuoteHash, "")
        if err != nil { return pkg.AcceptPeginRespose{}, err }
        return pkg.AcceptPeginRespose{
            Signature: accepted.Signature,
            BitcoinDepositAddressHash: accepted.DepositAddress,
        }, nil
    },
    ErrorMap: []ErrorMapping{
        {Is: usecases.QuoteNotFoundError,    Code: 404, Message: "quote not found"},
        {Is: usecases.ExpiredQuoteError,     Code: 410, Message: "expired quote"},
        {Is: usecases.NoLiquidityError,      Code: 409, Message: "not enough liquidity"},
        {Is: blockchain.ContractPausedError, Code: 503, Message: "protocol is paused"},
    },
}
```

OpenAPI: the CI step walks every `Route[Req, Resp]` value, derives `(method, path, request schema, response schema, error codes)`, and asserts `OpenApi.yml` matches. PRs that change a descriptor without updating the spec fail in CI.

## Out of scope

- Picking the final approach. The doc proposes.
- Migrating any routes.
- Auth/CSRF/session changes beyond middleware-as-data.
- Replacing `gorilla/mux` (decision 4 above, deliberately surfaced rather than bundled).
- DTO and validator-tag conventions in `pkg/`. The scaffold consumes whatever DTOs we have.
- Query-parameter routes like `/reports/summaries?startDate=&endDate=`. Pseudocode here covers JSON-body and no-body cases; query params are a v2 concern.

## References

- `internal/adapters/entrypoints/rest/routes/routes.go:16-94` — `Endpoint`, `EndpointFactory`, `ConfigureRoutes`, register helpers.
- `internal/adapters/entrypoints/rest/routes/public.go:11-141` — `PublicEndpoint`, 17 entries.
- `internal/adapters/entrypoints/rest/routes/management.go:13-178` — `AllowedPaths`, 27 entries.
- `internal/adapters/entrypoints/rest/common.go:29, 122, 163-206, 234-250` — validator chain, generic helpers.
- `internal/adapters/entrypoints/rest/handlers/accept_pegin_quote.go:19-70` — canonical handler.
- `internal/adapters/entrypoints/rest/handlers/common.go:15-39` — `HandleAcceptQuoteError`.
- `internal/adapters/entrypoints/rest/handlers/get_report_summaries.go:27-82` — singleflight handler.
- `internal/adapters/entrypoints/rest/registry/registry.go:11-54` — `UseCaseRegistry`.
- `OpenApi.yml` — repo root.
