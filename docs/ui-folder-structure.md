# Management UI (`ui/`) — folder structure

Canonical layout for the Vite + React app under `liquidity-provider-server/ui/`. Intended for humans and AI agents adding code in later epic stories.

## Top level

```
ui/
├── src/
│   ├── shared/          # Cross-feature code (only area with components/types/utils at top level today)
│   ├── feature/         # One subdirectory per feature
│   ├── lib/             # One subdirectory per internal library
│   ├── api/             # HTTP/OpenAPI client modules
│   ├── App.tsx          # Root component (scaffold; may move later)
│   └── main.tsx
├── tests/               # All test code — never colocate *.test.* under src/
│   ├── unit/            # Vitest + RTL; mirrors src/ layout by feature
│   ├── fixtures/        # Shared test data (`index.ts`)
│   ├── utils/           # Test-only utilities (e.g. DOM seeding, render helpers)
│   └── setup/           # Vitest global setup
├── docs/                # (this file lives in repo-root docs/)
└── …config files
```

## `shared/` — use now

Shared code that is not tied to a single feature. **This is the only tree that currently contains `components/`, `types/`, and `utils/` at the first level below the area name.**

```
src/shared/
├── components/    # Reusable UI (buttons, layout, etc.)
├── types/         # Shared TypeScript types
└── utils/         # Pure helpers, formatters, guards
```

Import via path aliases, e.g. `@shared/components/...`, `@shared/types/...`, `@shared/utils/...`.

## `feature/` — one folder per feature (not yet used)

Do **not** put `components/`, `types/`, or `utils/` directly under `feature/`. Each feature gets its own named directory:

```
src/feature/
└── <feature-name>/           # e.g. auth, peg-in-quotes, provider-settings
    ├── components/
    ├── types/
    └── utils/
```

Example import: `@feature/auth/components/LoginForm`.

Until a story introduces the first feature, `feature/` stays empty (`.gitkeep` only).

## `lib/` — one folder per library (not yet used)

Same pattern as features, for reusable non-UI modules (e.g. `formatting`, `validation`):

```
src/lib/
└── <lib-name>/
    ├── components/   # only if the lib exposes UI primitives
    ├── types/
    └── utils/
```

Example import: `@lib/validation/utils/isAddress`.

Until a lib is introduced, `lib/` stays empty (`.gitkeep` only).

## `api/` — client modules (not yet used)

When the UI talks to the Management API, add named client areas (exact split can follow OpenAPI tags):

```
src/api/
└── <client-name>/     # e.g. management, quotes
    ├── types/
    └── utils/         # request helpers, mappers
```

Do not add a flat `api/components|types|utils` without a `<client-name>/` parent.

## `tests/` — unit tests (outside `src/`)

**Do not** place `*.test.ts(x)` under `src/`. Unit tests live under `tests/` at the `ui/` root.

```
tests/
├── unit/                    # Vitest + React Testing Library
│   ├── app/                 # Root App / router integration
│   ├── feature/<name>/      # Mirrors src/feature/<name>/
│   ├── api/<client>/        # Mirrors src/api/<client>/
│   └── shared/              # Mirrors src/shared/
├── fixtures/                # Typed payloads shared across unit tests
├── utils/                   # Test-only utilities (`index.tsx`: seeding, render helpers)
└── setup/                   # vitest-setup.ts (global beforeAll hooks)
```

**Unit test naming:** `tests/unit/feature/auth/AuthGuard.test.tsx` tests `@feature/auth/components/AuthGuard`.

**Commands:** `pnpm test` and `pnpm test:coverage` run Vitest against `tests/unit/**`.

## Path aliases (`tsconfig.app.json` / `tsconfig.test.json`)

| Alias | Maps to |
|-------|---------|
| `@/*` | `src/*` |
| `@shared/*` | `src/shared/*` |
| `@feature/*` | `src/feature/*` |
| `@lib/*` | `src/lib/*` |
| `@api/*` | `src/api/*` |
| `@tests/*` | `tests/*` (unit tests and utils only; `tsconfig.test.json`) |

Prefer aliases over deep relative imports (`../../`).

## What not to do

- Do not create `src/feature/components` (or `types` / `utils`) — that implies a single feature.
- Do not create `src/lib/components` at the lib root — use `src/lib/<lib-name>/...`.
- Do not colocate unit tests under `src/` — use `tests/unit/<area>/`.
- Do not embed `ui/dist/` in Go or change `management.html` / Management API routes in scaffold-only stories unless the ticket explicitly says so.

## Related docs

- LPS README: [Management UI workspace](../README.md#management-ui-workspace-ui) (commands, CI, hooks)
- Epic: FLY-2305 — Migrate Management UI to Vite + React
