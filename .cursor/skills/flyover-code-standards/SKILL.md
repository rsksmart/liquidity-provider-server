---
name: flyover-code-standards
description: >-
  Enforce Flyover protocol coding standards derived from lead reviewer patterns.
  Use when writing, reviewing, or modifying Go or JavaScript code in the
  liquidity-provider-server. Covers architecture, naming, testing, API design,
  MongoDB patterns, security, and management UI CSP rules.
---

# Flyover Code Standards

Standards extracted from code review comments across the Flyover repos.
Apply these when writing or modifying code to reduce review friction.

## Quick Checklist

Before submitting code, verify:

- [ ] **Layer boundaries respected** — entities don't import infra; use cases return domain objects, not DTOs
- [ ] **Use case helpers are methods** — private helpers used only inside a use case are receivers of the use case struct, not package-level functions
- [ ] **Functions are small** — no disabled `cyclop`/`maintidx` linter rules; compose smaller functions instead
- [ ] **Naming is accurate** — camelCase JSON tags, no redundant prefixes, names match actual behavior
- [ ] **HTTP verbs are correct** — no body in GET/DELETE; paths are resources, verbs are actions
- [ ] **Errors are specific** — proper HTTP status codes; `errors.Is` for type checks; `fmt.Errorf` for wrapping
- [ ] **Pointers are justified** — prefer values; nil-check every pointer; use zero-value structs over nil
- [ ] **Tests are complete** — assert error types, all result fields, events; use mockery; expected on left
- [ ] **Test package uses `_test` suffix** — test files declare `package foo_test`, not `package foo`
- [ ] **Security basics** — `textContent` not `innerHTML` in JS; `BigInt` not `Number` for large values in JS; verify config signatures in Go
- [ ] **Management UI styles** — never use inline `<style>` blocks in HTML templates; put CSS in `static/<page>.css` (trusted via `style-src 'self'`); inline blocks are blocked by CSP when `ENABLE_SECURITY_HEADERS=true`

## Architecture (Go)

- **Entity layer**: no infra imports (no Ethereum crypto, no Mongo operators). Receive dependencies as parameters.
- **Use cases**: expose only `Run()`. Need more? Create separate use cases.
- **Use case helpers**: private helper functions used exclusively inside a use case must be methods of the use case struct, not package-level functions. E.g. `func (useCase *MyUseCase) calculateFoo(...)` instead of `func calculateFoo(...)`.
- **Handlers**: convert domain objects to DTOs here, not in use cases. Return DTOs, never domain entities.
- **Shared interfaces**: define at entity layer, not in use cases.
- **DB repos**: one collection = one file. Different collection = different repository struct.
- **Config reads**: always go through the `LiquidityProvider` interface (validates signature integrity), never read directly from the repo.

## Function Design (Go)

- Avoid `continue` — put logic inside the `if` block instead.
- Avoid multiple returns when semantics aren't obvious — use a result struct.
- Group related params into a struct. If two fields come from the same struct, pass the struct.
- Use constructors for struct creation. Callbacks go as the last parameter.
- Use `chan struct{}` for signal-only channels, not `chan bool`.
- Use `log.Debug`/`log.Error`, not `fmt.Printf`.
- Use `time.DateOnly`/`time.RFC3339` constants, never hardcode `"2006-01-02"`.
- `*big.Int` in DTOs, `*entities.Wei` in domain. Gas used = `big.Int` (not a currency). Gas price = `Wei`.
- Move strings to constants on the third occurrence.

## Management UI (CSS / CSP)

- All page-specific CSS goes in `static/<page-name>.css`, linked via `<link href="/static/<page-name>.css">`.
- Never add or edit inline `<style>` blocks in management HTML templates. The CSP `style-src` directive only allows external files (`'self'`) and a single pinned hash. Any edit to a `<style>` block changes its hash → the browser silently discards the **entire** block.
- Scripts already use a per-request nonce (`{{ .ScriptNonce }}`); styles do not — external files are the only safe pattern.
- Test with `ENABLE_SECURITY_HEADERS=true` locally. With the flag off, CSP is never sent and style blocks appear to work even when broken.

See [detailed.md](detailed.md) for full patterns with rationale and examples.
