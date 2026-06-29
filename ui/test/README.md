# LPS Management UI — Playwright E2E

Browser E2E tests for the **next** management UI shell served by LPS at `/management/next/`. Specs assert server-templated security properties (CSRF meta, CSP nonces) and placeholder page headings.

## Layout

```
ui/
├── playwright.config.ts       # testDir: ./test/e2e
├── package.json               # test:e2e, test:e2e:ui scripts
└── test/
    ├── e2e/                   # runnable *.spec.ts only
    ├── fixtures/              # session helper, shared test exports
    └── .auth/                 # gitignored storageState (session.setup.ts)
```

Vitest unit tests stay under `tests/unit/` and are **not** run by Playwright.

## Prerequisites

1. **Embedded dist** — LPS must serve the built shell, not Vite dev stubs:

   ```bash
   make ui-build
   ```

2. **Running LPS** on `:8080` with next UI enabled (regtest stack or local dev):

   ```bash
   cd docker-compose/local && LPS_STAGE=regtest ./lps-local.sh
   ```

3. **Browser binaries** (one-time per machine):

   ```bash
   cd ui
   pnpm exec playwright install --with-deps chromium
   ```

## Commands

```bash
cd ui
pnpm install

# List specs (no LPS required)
pnpm test:e2e --list

# Full run against running LPS
LPS_E2E_BASE_URL=http://localhost:8080/management/next/ \
LPS_E2E_USER=<username> \
LPS_E2E_PASSWORD=<password> \
  pnpm test:e2e

# Headed UI mode (debug)
pnpm test:e2e:ui
```

From repo root:

```bash
make ui-e2e
```

## Environment variables

| Variable | Default | Purpose |
|----------|---------|---------|
| `LPS_E2E_BASE_URL` | `http://localhost:8080/management/next/` | Playwright `baseURL` — must be LPS-served shell (trailing slash keeps relative navigations under `/management/next/`) |
| `LPS_E2E_USER` | — | Management login username (authenticated specs only) |
| `LPS_E2E_PASSWORD` | — | Management login password (never commit) |

Credentials come from `docker-compose/local` test config or CI secrets. **Never** hardcode passwords in specs.

## Scope guard

- All specs use `baseURL` under `/management/next/` only.
- No tests against legacy `/management` HTML dashboard routes.
- No POST to collateral, configuration, or trusted-accounts APIs.
- Security assertions use **served HTML** via `request.get()` (not Vite `:5173` dev server).

## shadcn/ui

The next management UI uses [shadcn/ui](https://ui.shadcn.com) (Tailwind v4 + copied components under `src/components/ui/`).

| Component | Use |
| --------- | --- |
| `sonner` | Session-expired and logout-failure toasts |
| `alert` | Inline login error banner |
| `button`, `input`, `label`, `card` | Login form and logout control |

Add components with the shadcn CLI from `ui/`:

```bash
pnpm dlx shadcn@latest add <component>
```

Config: `components.json` (`@/` → `src/`).

## Not yet covered

Legacy `assets/management.html` areas without Playwright specs in this directory:

| Area | Legacy surface |
|------|----------------|
| Collateral overview | `/management` collateral tab |
| Trusted accounts | trusted-accounts APIs + UI |
| Provider cards | provider configuration surfaces |
| Configuration editors | `/configuration` POST flows |

Playwright is not wired into `.github/workflows/e2e.yml` yet; use `pnpm test:e2e` or `make ui-e2e` locally.
