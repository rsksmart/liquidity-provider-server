# LPS × LBC Compatibility Matrix

Declared support window for **split-contract** Flyover deployments (PegIn, PegOut, FlyoverDiscovery, CollateralManagement).

Machine-readable source of truth lives in the e2e repo:

`flyover-lps-api-e2e/compat/lps-lbc-matrix.yaml`

## Declared matrix (PoC)

| LPS \ LBC | QA-Test | master |
|-----------|---------|--------|
| **QA-Test** | SUPPORTED | NOT TESTED |
| **master** | NOT TESTED | NOT TESTED |

- **SUPPORTED** — declared safe for split-contract deployments (smoke-validated locally)
- **NOT TESTED** — in matrix scope; not yet validated
- **UNSUPPORTED** — declared incompatible (none in the PoC window)

Pre-release tags (`*-rc*`), cross-MAJOR pairs, and legacy monolithic `LBC_ADDR` deployments are out of scope. See [LP Migration Utilities](./LP-Migration-Utils.md).

## How to run

Prerequisites: Docker, Git, Node.js 24+, sibling clones of `liquidity-bridge-contract` and `flyover-lps-api-e2e`.

From this repo:

```bash
cd docker-compose/local/lps-lbc-compat

# Single pair (default QA-Test × QA-Test)
./run-pair.sh

# Trunk pair (positional args: <lbc-ref> <lps-ref>)
./run-pair.sh master master

# Smokes only against an already-running stack
./run-pair.sh --skip-lps-up

# Full 2×2 matrix (cleans volumes before each cell)
./run-matrix.sh
```

Orchestration lives next to `lps-local.sh`:

| Script | Role |
|--------|------|
| `run-pair.sh` | Checkout LPS+LBC refs (LPS via worktree), deploy (incl. volume clean / BTC rescan recovery), run e2e smokes |
| `run-matrix.sh` | Loop declared pairs, grade cells, print green/orange/red grid |

Smoke tests and matrix grading live in **flyover-lps-api-e2e**:

```bash
cd ../flyover-lps-api-e2e
npm run test:regtest:smoke   # against a running LPS stack
npm run compat:print         # print declared matrix
npm run compat:cli -- pairs-tsv
```

## Runtime grades

After a matrix run, cells are graded from smoke counts:

- **Green** (`N/N`) — all smokes passed
- **Orange** (`k/N`) — some smokes passed
- **Red** (`0/N` or FAIL) — none passed or deploy failed
