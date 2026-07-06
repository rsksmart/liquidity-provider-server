# Why `LockingCapExceededError` and `TamperedTrustedAccountError` return `500` on public endpoints

Quote handlers map use-case errors to HTTP responses via shared helpers in
[`handlers/common.go`](../internal/adapters/entrypoints/rest/handlers/common.go).
Most of that is straightforward. One thing is worth calling out: on the public accept
endpoints, two errors come back as `500`. **That is intentional — please don't "fix" it.**

It can look like a mistake at first. The public accept endpoints return `500` with
`"unknown error"` for `LockingCapExceededError` and `TamperedTrustedAccountError`, while
the authenticated endpoints return `409` with a clear message for the same errors. The
asymmetry is deliberate.

Both errors belong to the trusted-account flow, which only runs when the use case is
invoked **with a signature**. Public handlers always pass an empty signature
(`useCase.Run(ctx, quoteHash, "")`), so that flow never runs and these errors should
not appear on a public route. Only authenticated handlers supply a signature.

If one of them does show up on a public endpoint, something is wrong on our side — for
example, auth-only logic wired onto a public route. That is a server bug, not something
the caller can work around. A `500` with `recoverable: false` says so plainly. Mapping
it to `409` would paper over the bug and suggest the client should retry.