# Content Security Policy in the Management UI

## What is CSP?

Content Security Policy (CSP) is a browser security mechanism delivered via an HTTP response header. It tells the browser which resources it is allowed to execute or apply. Its primary purpose is to prevent **Cross-Site Scripting (XSS)** attacks: even if an attacker manages to inject `<script>` or `<style>` tags into the HTML, the browser will refuse to run them if they are not explicitly approved by the policy.

---

## Where is it configured in this repo?

The CSP header for all management UI pages is set in:

```
internal/adapters/entrypoints/rest/handlers/management.go
```

Inside the `htmlTemplateSecurityHeaders` function. It is only sent when `ENABLE_SECURITY_HEADERS=true` in the environment configuration (the `management.go` handler checks this flag before calling the function).

The header format (simplified) is:

```
Content-Security-Policy:
  default-src 'self';
  font-src 'self' data:;
  style-src 'self' 'sha256-<HASH>';
  object-src 'none';
  frame-src 'self';
  script-src 'self' 'nonce-<RANDOM>';
  img-src 'self' data:;
  connect-src 'self';
```

The two directives that affect page development the most are `style-src` and `script-src`.

---

## How `script-src` works — the nonce pattern

Scripts use a **nonce** (number used once): a random token generated fresh on every HTTP request. The same token is placed both in the CSP header and in the `<script>` tag:

```html
<script nonce="{{ .ScriptNonce }}">
  // your code
</script>
```

The browser checks: does the tag's `nonce` attribute match the nonce in the header? If yes, it runs the script.

**Key property:** because the nonce is regenerated on every request and the tag gets the same value via Go's template engine, you can edit the script content freely and it will always be trusted. The nonce proves intent at request-time, not content-fingerprint.

---

## How `style-src` works — the hash pattern

Inline `<style>` blocks use a different mechanism: a **SHA-256 hash of the exact block content**, pre-computed and pinned in the header.

The browser, on page load:
1. Finds each inline `<style>` block.
2. Computes its SHA-256 hash.
3. Compares it against the hashes listed in `style-src`.
4. If the hashes match → applies the styles. If not → **silently discards the entire block**.

This means: **editing even a single character of an inline `<style>` block invalidates the pinned hash**, and the browser will discard all the styles in that block without any visible error in the UI.

### Why this is fragile

The hash must be recomputed and updated in `management.go` every time the `<style>` block changes. If you forget, the styles will silently stop working — the UI will look broken in production (where `ENABLE_SECURITY_HEADERS=true`) but may appear fine in local development (where the flag is often `false`).

---

## The established pattern in this codebase: use external stylesheets

To avoid the hash maintenance problem entirely, the correct approach (already used for the main dashboard) is to **never put styles in inline `<style>` blocks**. Instead, place all CSS in external `.css` files under:

```
internal/adapters/entrypoints/rest/assets/static/
```

These files are served from the same origin and are covered by `style-src 'self'`, so the browser trusts them unconditionally — no hash, no nonce required.

The Go embed directive in `rest/assets/` picks up the entire `static/` directory:

```go
//go:embed static favicon.ico
```

So adding a new `.css` file to `static/` is all that's needed — no other Go changes.

### Linking the stylesheet in an HTML template

```html
<link href="/static/your-page.css" rel="stylesheet" />
```

One CSS file per page (rather than appending to the shared `management.css`) is preferred when the page has styles that would conflict with selectors already in the shared file (e.g. `.table-responsive` has different `max-height` values per page).

---

## Historical context: the stale hash

As of the time of writing, `management.go` still contains a pinned SHA-256 hash.

This hash was computed for the **dashboard** (`management.html`) `<style>` block at the time it was written. The dashboard's inline `<style>` was subsequently emptied (its rules were moved to `management.css`), so this hash now matches an empty block and is effectively vestigial — it guards nothing today.

**Do not rely on this hash for new pages.** If you ever need to add a new inline style block, compute a fresh hash for it:

```bash
# Compute the hash for the content between <style>...</style> (not including the tags themselves)
echo -n "<your style content>" | openssl dgst -sha256 -binary | openssl base64
```

Then add it to `style-src` in `management.go`. But again, the preferred approach is to avoid inline styles altogether and use an external stylesheet.

---

## Practical rules for developers

| Situation | What to do |
|---|---|
| Adding new CSS for a management page | Create `static/<page-name>.css`, add `<link>` to the template. **Do not use `<style>` blocks.** |
| Adding JavaScript to a management template | Use `<script nonce="{{ .ScriptNonce }}">`. The nonce is injected automatically — no additional Go changes needed. |
| Debugging styles that aren't applying | Open DevTools → Console. Look for CSP violation messages. Check `document.styleSheets` — if your inline `<style>` block is missing from the list, CSP blocked it. |
| Testing with security headers | Set `ENABLE_SECURITY_HEADERS=true` in your local `.env` file, otherwise CSP is never sent and you won't catch these issues locally. |
