#!/usr/bin/env bash
#
# Checkout LPS + LBC refs, deploy the local regtest stack, then run Flyover smokes
# from flyover-lps-api-e2e against that pair.
#
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LPS_LOCAL="$(cd "${SCRIPT_DIR}/.." && pwd)"
LPS_REPO="${LPS_REPO:-$(cd "${LPS_LOCAL}/../.." && pwd)}"
LBC_REPO="${LBC_REPO:-${LPS_REPO}/../liquidity-bridge-contract}"
E2E_REPO="${E2E_REPO:-${LPS_REPO}/../flyover-lps-api-e2e}"

LPS_REF="${LPS_REF:-QA-Test}"
LBC_REF="${LBC_REF:-QA-Test}"
LPS_ENV_FILE="${LPS_ENV_FILE:-}"
SKIP_LPS_UP="${SKIP_LPS_UP:-false}"
START_EXISTING="${START_EXISTING:-false}"
CLEAN_BEFORE_UP="${CLEAN_BEFORE_UP:-false}"
SMOKE_TESTS=""
SMOKE_TESTS_EXPLICIT=false
SMOKE_PASSED=0
SMOKE_FAILED=0
SMOKE_TOTAL=0
LPS_WORKTREE=""
LPS_DEPLOY_ROOT=""

cleanup() {
  if [[ -n "$LPS_WORKTREE" && -d "$LPS_WORKTREE" ]]; then
    git -C "$LPS_REPO" worktree remove --force "$LPS_WORKTREE" 2>/dev/null || rm -rf "$LPS_WORKTREE"
  fi
}
trap cleanup EXIT

lps_compose_lps_files() {
  local lps_local="$1"
  if [[ -f "${lps_local}/lps/docker-compose.lps-local.yml" ]]; then
    echo "-f docker-compose.yml -f lps/docker-compose.lps-local.yml"
  elif [[ -f "${lps_local}/docker-compose.lps.yml" ]]; then
    echo "-f docker-compose.yml -f docker-compose.lps.yml"
  else
    echo "-f docker-compose.yml"
  fi
}

wait_for_lps_health() {
  local url="$1"
  local attempts="${2:-60}"
  local i=1
  while [[ $i -le $attempts ]]; do
    if curl -sf "${url}/health" >/dev/null 2>&1; then
      echo "LPS health OK: ${url}/health"
      return 0
    fi
    sleep 5
    i=$((i + 1))
  done
  echo "ERROR: LPS not healthy after ${attempts} attempts (${url}/health)" >&2
  return 1
}

# Wipe bind-mount volumes so a fresh deploy can chown cleanly
# (avoids WSL/Docker "Operation not permitted" when container-owned files linger).
clean_lps_stack_volumes() {
  local lps_local="$1"
  local env_file="$2"
  local lps_uid="${LPS_UID:-$(id -u)}"
  local -a volume_dirs=(
    ./volumes/bitcoind
    ./volumes/rskj
    ./volumes/powpeg
    ./volumes/lps
    ./volumes/mongo
    ./volumes/localstack
  )

  echo "Cleaning LPS stack volumes before deploy..."
  (
    cd "$lps_local"
    export LPS_UID="$lps_uid"

    if [[ -f "$env_file" ]]; then
      docker compose --env-file "$env_file" \
        -f docker-compose.yml \
        -f docker-compose.lps.yml \
        -f docker-compose.lbc-deployer.yml \
        down --remove-orphans 2>/dev/null || true
      # Newer layouts use lps/docker-compose.lps-local.yml instead of docker-compose.lps.yml
      docker compose --env-file "$env_file" \
        -f docker-compose.yml \
        -f lps/docker-compose.lps-local.yml \
        down --remove-orphans 2>/dev/null || true
      docker compose --env-file "$env_file" down --remove-orphans 2>/dev/null || true
    else
      docker stop lps01 mongo01 rskj01 bitcoind01 localstack powpeg-pegin powpeg-pegout 2>/dev/null || true
      docker rm -f lps01 mongo01 rskj01 bitcoind01 localstack powpeg-pegin powpeg-pegout 2>/dev/null || true
    fi

    if [[ -d ./volumes ]]; then
      docker run --rm -v "$(pwd)/volumes:/volumes" alpine:3.20 \
        sh -c 'rm -rf /volumes/bitcoind /volumes/rskj /volumes/powpeg /volumes/lps /volumes/mongo /volumes/localstack' \
        2>/dev/null || rm -rf "${volume_dirs[@]}" 2>/dev/null || true
    fi

    mkdir -p ./volumes/bitcoind
    mkdir -p ./volumes/rskj/db ./volumes/rskj/logs
    mkdir -p ./volumes/powpeg/pegin/db ./volumes/powpeg/pegin/logs
    mkdir -p ./volumes/powpeg/pegout/db ./volumes/powpeg/pegout/logs
    mkdir -p ./volumes/lps/logs
    mkdir -p ./volumes/mongo/db ./volumes/mongo/logs
    mkdir -p ./volumes/localstack/db ./volumes/localstack/logs

    chown -R "$lps_uid" ./volumes/bitcoind ./volumes/rskj ./volumes/mongo ./volumes/localstack
    chown -R "$lps_uid" ./volumes/powpeg ./volumes/lps
    chmod -R 777 ./volumes/powpeg ./volumes/lps
  )
  echo "LPS volume cleanup complete."
}

# Recover LPS when bitcoind wallet rescan blocks startup ("rescan started" in lps01 logs).
recover_lps_btc_wallet_rescan() {
  local lps_local="$1"
  local env_file="$2"

  if ! docker logs lps01 2>&1 | grep -q 'rescan started'; then
    return 1
  fi

  local btc_user btc_pass compose_lps
  btc_user="$(grep -E '^BTC_USERNAME=' "$env_file" | cut -d= -f2- | tr -d '"')"
  btc_pass="$(grep -E '^BTC_PASSWORD=' "$env_file" | cut -d= -f2- | tr -d '"')"
  btc_user="${btc_user:-test}"
  btc_pass="${btc_pass:-test}"
  compose_lps="$(lps_compose_lps_files "$lps_local")"

  echo "BTC wallet rescan in progress — waiting, then restarting LPS..."
  local i=1
  while [[ $i -le 60 ]]; do
    if docker exec bitcoind01 bitcoin-cli -rpcuser="$btc_user" -rpcpassword="$btc_pass" -rpcport=5555 \
      -rpcwallet=rsk-wallet getwalletinfo 2>/dev/null | grep -q '"scanning": false'; then
      break
    fi
    sleep 5
    i=$((i + 1))
  done

  (
    cd "$lps_local"
    export LPS_UID="${LPS_UID:-$(id -u)}"
    # shellcheck disable=SC2086
    docker compose $compose_lps --env-file "$env_file" up -d lps
  )

  wait_for_lps_health "http://127.0.0.1:8080" 36
}

usage() {
  cat <<'EOF'
Usage: run-pair.sh [options] [<lbc-ref> <lps-ref>]

Checkout LPS + LBC git refs, deploy the local regtest stack, and run
flyover-lps-api-e2e smoke tests against that pair.

Options:
  --lps-repo PATH       LPS repository (default: this repo)
  --lbc-repo PATH       LBC repository (default: ../liquidity-bridge-contract)
  --e2e-repo PATH       flyover-lps-api-e2e repo (default: ../flyover-lps-api-e2e)
  --lps-ref REF         Git branch/tag in LPS (default: QA-Test)
  --lbc-ref REF         Git branch/tag in LBC (default: QA-Test)
  <lbc-ref> <lps-ref>   Positional pair (LBC first, then LPS)
  --lps-env-file FILE   LPS env file after deploy (default: .env.regtest)
  --skip-lps-up         Do not start LPS; expect stack already running
  --clean-before-up     Stop LPS stack and wipe bind-mount volumes before deploy
  --start-existing      Start stopped local containers only (no image pull / deploy)
  --smoke-tests LIST    Comma-separated smoke ids from compat/lps-lbc-matrix.yaml
  -h, --help            Show this help

Examples:
  ./run-pair.sh
  ./run-pair.sh master master
  ./run-pair.sh --skip-lps-up
  ./run-pair.sh --clean-before-up --lps-ref QA-Test --lbc-ref master
EOF
}

POSITIONAL=()

while [[ $# -gt 0 ]]; do
  case "$1" in
    --lps-repo)
      LPS_REPO="$2"
      shift 2
      ;;
    --lbc-repo)
      LBC_REPO="$2"
      shift 2
      ;;
    --e2e-repo)
      E2E_REPO="$2"
      shift 2
      ;;
    --lps-ref)
      LPS_REF="$2"
      shift 2
      ;;
    --lbc-ref)
      LBC_REF="$2"
      shift 2
      ;;
    --lps-env-file)
      LPS_ENV_FILE="$2"
      shift 2
      ;;
    --skip-lps-up)
      SKIP_LPS_UP=true
      shift
      ;;
    --clean-before-up)
      CLEAN_BEFORE_UP=true
      shift
      ;;
    --start-existing)
      START_EXISTING=true
      shift
      ;;
    --smoke-tests)
      SMOKE_TESTS="$2"
      SMOKE_TESTS_EXPLICIT=true
      shift 2
      ;;
    -h | --help)
      usage
      exit 0
      ;;
    --*)
      echo "Unknown option: $1" >&2
      usage >&2
      exit 1
      ;;
    *)
      POSITIONAL+=("$1")
      shift
      ;;
  esac
done

if [[ ${#POSITIONAL[@]} -eq 2 ]]; then
  LBC_REF="${POSITIONAL[0]}"
  LPS_REF="${POSITIONAL[1]}"
elif [[ ${#POSITIONAL[@]} -eq 1 ]]; then
  echo "ERROR: provide both LBC and LPS refs: <lbc-ref> <lps-ref>" >&2
  usage >&2
  exit 1
elif [[ ${#POSITIONAL[@]} -gt 2 ]]; then
  echo "ERROR: too many positional arguments (expected <lbc-ref> <lps-ref>)" >&2
  usage >&2
  exit 1
fi

MATRIX_CLI=(npx --yes tsx "${E2E_REPO}/scripts/compat-matrix-cli.ts")

matrix_cli() {
  (
    cd "$E2E_REPO"
    "${MATRIX_CLI[@]}" "$@"
  )
  return 0
}

resolve_smoke_settings() {
  if [[ "$SMOKE_TESTS_EXPLICIT" != "true" ]]; then
    SMOKE_TESTS="$(matrix_cli smoke-tests)"
  fi
}

require_cmd() {
  local cmd="$1"
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "ERROR: required command not found: ${cmd}" >&2
    exit 1
  fi
}

resolve_lps_env_file() {
  if [[ -n "$LPS_ENV_FILE" ]]; then
    echo "$LPS_ENV_FILE"
    return
  fi
  echo "${LPS_LOCAL}/.env.regtest"
}

# Checkout only — never patch tracked files in LPS/LBC repos.
checkout_git_ref() {
  local repo="$1"
  local ref="$2"
  local label="$3"
  local resolved=""
  local sha subject
  local -a candidates=()

  if [[ ! -d "$repo/.git" ]]; then
    echo "ERROR: ${label} repo not found: ${repo}" >&2
    exit 1
  fi

  echo "Resolving ${label} ref: ${ref}"
  if ! git -C "$repo" fetch origin --tags --quiet; then
    echo "WARNING: git fetch origin failed for ${label}; will try local refs only" >&2
  fi

  if [[ "$ref" == origin/* ]]; then
    candidates+=("$ref")
    candidates+=("${ref#origin/}")
  else
    candidates+=("origin/${ref}")
    candidates+=("$ref")
    if git -C "$repo" rev-parse --verify "refs/tags/${ref}" >/dev/null 2>&1; then
      candidates=("refs/tags/${ref}" "${candidates[@]}")
    fi
  fi

  for candidate in "${candidates[@]}"; do
    if git -C "$repo" checkout --detach "$candidate" 2>/dev/null; then
      resolved="$candidate"
      break
    fi
  done

  if [[ -z "$resolved" && "$label" == "LBC" ]]; then
    for alt in QA-Test master v2.5.0-fixes version-2.5.0-fixes; do
      if [[ "$alt" == "$ref" ]]; then
        continue
      fi
      if git -C "$repo" checkout --detach "origin/${alt}" 2>/dev/null; then
        echo "NOTE: LBC ref '${ref}' not found on origin; using origin/${alt}" >&2
        resolved="origin/${alt}"
        break
      fi
    done
  fi

  if [[ -z "$resolved" ]]; then
    echo "ERROR: could not checkout ${label} ref '${ref}' (tried origin/${ref} and local)" >&2
    exit 1
  fi

  sha="$(git -C "$repo" rev-parse --short HEAD)"
  subject="$(git -C "$repo" log -1 --format='%s' HEAD)"
  echo "${label} checkout: ${resolved} @ ${sha} — ${subject}"
}

load_flyover_env_from_lps() {
  local env_file="$1"
  if [[ ! -f "$env_file" ]]; then
    echo "ERROR: LPS env file not found: $env_file" >&2
    exit 1
  fi

  local -a required_vars=(
    PEGIN_CONTRACT_ADDRESS
    PEGOUT_CONTRACT_ADDRESS
    DISCOVERY_ADDRESS
    COLLATERAL_MANAGEMENT_ADDRESS
  )

  for var in "${required_vars[@]}"; do
    local line
    line="$(grep -E "^${var}=" "$env_file" | tail -1 || true)"
    if [[ -z "$line" ]]; then
      echo "ERROR: missing ${var} in ${env_file}" >&2
      exit 1
    fi
    local value="${line#*=}"
    value="${value%\"}"
    value="${value#\"}"
    if [[ "$value" == "0x0000000000000000000000000000000000000000" || -z "$value" ]]; then
      echo "ERROR: ${var} is unset in ${env_file}; run with DEPLOY_CONTRACTS=true" >&2
      exit 1
    fi
    export "$var=$value"
  done

  local port_line server_port
  port_line="$(grep -E '^SERVER_PORT=' "$env_file" | tail -1 || true)"
  server_port="${port_line#*=}"
  server_port="${server_port:-8080}"
  export FLYOVER_LPS_URL="http://127.0.0.1:${server_port}"
  export LPS_BASE_URL="${FLYOVER_LPS_URL}"
  export FLYOVER_RSK_RPC_URL="${FLYOVER_RSK_RPC_URL:-http://127.0.0.1:4444}"
  export RPC_URL="${FLYOVER_RSK_RPC_URL}"
}

reset_localstack_for_lps() {
  local lps_local="$1"
  local env_file="$2"
  echo "Resetting localstack (ensure flat keystore secret for LPS regtest)..."
  (
    cd "$lps_local"
    export LPS_UID="${LPS_UID:-$(id -u)}"
    local localstack_home="${LOCALSTACK_HOME:-./volumes/localstack}"
    docker compose --env-file "$env_file" stop localstack 2>/dev/null || true
    docker compose --env-file "$env_file" rm -f localstack 2>/dev/null || true
    rm -rf "$localstack_home"
    mkdir -p "${localstack_home}/db" "${localstack_home}/logs"
    chown -R "$LPS_UID" "$localstack_home"
    docker compose --env-file "$env_file" build localstack
    docker compose --env-file "$env_file" up -d --force-recreate localstack
  )
  sleep 10
}

start_lps_stack_lps_local() {
  local deploy_local="${LPS_DEPLOY_ROOT}/docker-compose/local"
  echo "Starting LPS regtest stack via lps-local.sh (LBC_COMMIT=${LBC_REF})..."
  (
    cd "$deploy_local"
    export LPS_UID="${LPS_UID:-$(id -u)}"
    export LBC_COMMIT="$LBC_REF"
    export DEPLOY_CONTRACTS=true
    export FUND_WALLETS=true
    ./lps-local.sh
  )
}

start_lps_stack_lps_env() {
  local deploy_local="${LPS_DEPLOY_ROOT}/docker-compose/local"
  local env_file
  env_file="$(resolve_lps_env_file)"

  reset_localstack_for_lps "$deploy_local" "$env_file"

  echo "Starting LPS regtest environment (lps-env.sh up)..."
  local lps_up_rc=0
  (
    cd "$deploy_local"
    export LPS_UID="${LPS_UID:-$(id -u)}"
    export LOG_FILE="${LOG_FILE:-}"
    LPS_STAGE=regtest bash ./lps-env.sh up
  ) || lps_up_rc=$?

  if [[ $lps_up_rc -ne 0 ]]; then
    if recover_lps_btc_wallet_rescan "$deploy_local" "$env_file"; then
      echo "LPS recovered after BTC wallet rescan — running lps-env.sh again for management configuration..."
      (
        cd "$deploy_local"
        export LPS_UID="${LPS_UID:-$(id -u)}"
        export LOG_FILE="${LOG_FILE:-}"
        LPS_STAGE=regtest bash ./lps-env.sh up
      ) || lps_up_rc=$?
      if [[ $lps_up_rc -ne 0 ]]; then
        echo "ERROR: lps-env.sh configuration pass failed after BTC rescan recovery" >&2
        exit 1
      fi
    else
      echo "ERROR: lps-env.sh up failed" >&2
      exit 1
    fi
  fi
}

start_existing_lps_stack() {
  local env_file compose_lps
  env_file="$(resolve_lps_env_file)"

  if [[ ! -f "$env_file" ]]; then
    echo "ERROR: no ${env_file}; run a full deploy first or pass --lps-env-file" >&2
    exit 1
  fi

  echo "Starting existing LPS containers (no pull, no contract deploy)..."
  export LPS_UID="${LPS_UID:-$(id -u)}"
  compose_lps="$(lps_compose_lps_files "$LPS_LOCAL")"
  (
    cd "$LPS_LOCAL"
    docker start mongo01 rskj01 bitcoind01 localstack 2>/dev/null || true
    # shellcheck disable=SC2086
    LOG_FILE= \
      docker compose $compose_lps --env-file "$env_file" up -d --no-build lps
  )
}

# Prefer a detached worktree when scripts live in LPS_REPO so checkout does not
# remove docker-compose/local/lps-lbc-compat from the runner's working tree.
prepare_lps_deploy_root() {
  local script_toplevel repo_toplevel safe_ref
  script_toplevel="$(cd "$SCRIPT_DIR" && git rev-parse --show-toplevel 2>/dev/null || true)"
  repo_toplevel="$(cd "$LPS_REPO" && git rev-parse --show-toplevel 2>/dev/null || true)"

  if [[ -n "$script_toplevel" && -n "$repo_toplevel" && "$script_toplevel" == "$repo_toplevel" ]]; then
    safe_ref="${LPS_REF//\//-}"
    LPS_WORKTREE="$(mktemp -d "${TMPDIR:-/tmp}/lps-compat-${safe_ref}-XXXXXX")"
    rm -rf "$LPS_WORKTREE"
    echo "Creating LPS worktree for deploy: ${LPS_WORKTREE} (${LPS_REF})"
    if ! git -C "$LPS_REPO" fetch origin --tags --quiet; then
      echo "WARNING: git fetch origin failed for LPS; will try local refs only" >&2
    fi
    if ! git -C "$LPS_REPO" worktree add --detach "$LPS_WORKTREE" "origin/${LPS_REF}" 2>/dev/null \
      && ! git -C "$LPS_REPO" worktree add --detach "$LPS_WORKTREE" "${LPS_REF}"; then
      echo "ERROR: could not create LPS worktree for ref '${LPS_REF}'" >&2
      exit 1
    fi
    LPS_DEPLOY_ROOT="$LPS_WORKTREE"
  else
    checkout_git_ref "$LPS_REPO" "$LPS_REF" "LPS"
    LPS_DEPLOY_ROOT="$LPS_REPO"
  fi

  LPS_LOCAL="${LPS_DEPLOY_ROOT}/docker-compose/local"
  if [[ ! -d "$LPS_LOCAL" ]]; then
    echo "ERROR: LPS local compose dir missing after checkout: ${LPS_LOCAL}" >&2
    exit 1
  fi
}

start_lps_stack() {
  if [[ ! -d "$LPS_REPO" ]]; then
    echo "ERROR: LPS repo not found: $LPS_REPO" >&2
    exit 1
  fi

  if [[ "$START_EXISTING" == "true" ]]; then
    LPS_DEPLOY_ROOT="$LPS_REPO"
    LPS_LOCAL="${LPS_DEPLOY_ROOT}/docker-compose/local"
    start_existing_lps_stack
    return
  fi

  prepare_lps_deploy_root
  checkout_git_ref "$LBC_REPO" "$LBC_REF" "LBC"

  if [[ "$CLEAN_BEFORE_UP" == "true" ]]; then
    clean_lps_stack_volumes "$LPS_LOCAL" "$(resolve_lps_env_file)"
  fi

  if [[ -f "${LPS_LOCAL}/lps-local.sh" ]]; then
    start_lps_stack_lps_local
  elif [[ -f "${LPS_LOCAL}/lps-env.sh" ]]; then
    start_lps_stack_lps_env
  else
    echo "ERROR: no supported LPS startup script under ${LPS_LOCAL}" >&2
    echo "  Expected lps-local.sh or lps-env.sh" >&2
    exit 1
  fi
}

smoke_file_for_id() {
  local smoke_id="$1"
  local files_tsv
  files_tsv="$(matrix_cli smoke-files)"
  local ids
  ids="$(matrix_cli smoke-tests)"
  local -a id_arr=()
  local -a file_arr=()
  IFS=',' read -r -a id_arr <<<"$ids"
  mapfile -t file_arr <<<"$files_tsv"
  local i
  for i in "${!id_arr[@]}"; do
    if [[ "${id_arr[$i]}" == "$smoke_id" ]]; then
      echo "${file_arr[$i]}"
      return 0
    fi
  done
  echo "ERROR: unknown smoke id '${smoke_id}'" >&2
  return 1
}

run_e2e_smokes() {
  local case_id test_rc smoke_file
  local -a case_ids=()

  if [[ ! -d "$E2E_REPO" ]]; then
    echo "ERROR: flyover-lps-api-e2e repo not found: ${E2E_REPO}" >&2
    exit 1
  fi

  cd "$E2E_REPO"

  if [[ ! -d node_modules ]]; then
    echo "Installing e2e dependencies (npm ci)..."
    npm ci
  fi

  IFS=',' read -r -a case_ids <<<"$SMOKE_TESTS"
  SMOKE_TOTAL=${#case_ids[@]}
  SMOKE_PASSED=0
  SMOKE_FAILED=0

  echo ""
  echo "Running ${SMOKE_TOTAL} smoke test(s) for LPS ${LPS_REF} x LBC ${LBC_REF}..."
  echo "  FLYOVER_LPS_URL=${FLYOVER_LPS_URL}"
  echo "  E2E_REPO=${E2E_REPO}"
  echo ""

  set +e
  for case_id in "${case_ids[@]}"; do
    smoke_file="$(smoke_file_for_id "$case_id")"
    echo "--- Smoke: ${case_id} (${smoke_file}) ---"
    npx jest --runInBand --verbose "$smoke_file"
    test_rc=$?
    if [[ $test_rc -eq 0 ]]; then
      SMOKE_PASSED=$((SMOKE_PASSED + 1))
      echo "  ${case_id}: PASS"
    else
      SMOKE_FAILED=$((SMOKE_FAILED + 1))
      echo "  ${case_id}: FAIL"
    fi
  done
  set -e

  echo "COMPAT_SMOKE_SUMMARY passed=${SMOKE_PASSED} failed=${SMOKE_FAILED} total=${SMOKE_TOTAL}"

  local status
  status="$(matrix_cli smoke-status "$SMOKE_PASSED" "$SMOKE_TOTAL")"
  case "$status" in
    pass) return 0 ;;
    partial) return 2 ;;
    *) return 1 ;;
  esac
}

main() {
  require_cmd git
  require_cmd docker
  require_cmd curl
  require_cmd npm
  require_cmd npx
  resolve_smoke_settings

  echo "=== LPS x LBC local compatibility check ==="
  echo "LPS repo:  ${LPS_REPO}"
  echo "LBC repo:  ${LBC_REPO}"
  echo "E2E repo:  ${E2E_REPO}"
  echo "LPS ref:   ${LPS_REF}"
  echo "LBC ref:   ${LBC_REF}"
  echo ""

  if [[ "$SKIP_LPS_UP" != "true" ]]; then
    start_lps_stack
  else
    echo "Skipping LPS startup (--skip-lps-up)"
    if ! docker ps --format '{{.Names}}' | grep -q 'lps01'; then
      echo "WARNING: lps01 container not running; smoke tests may fail" >&2
    fi
  fi

  local env_file
  env_file="$(resolve_lps_env_file)"
  load_flyover_env_from_lps "$env_file"
  wait_for_lps_health "$FLYOVER_LPS_URL"

  local smoke_rc=0
  set +e
  run_e2e_smokes
  smoke_rc=$?
  set -e

  local smoke_status
  smoke_status="$(matrix_cli smoke-status "$SMOKE_PASSED" "$SMOKE_TOTAL")"
  echo ""
  case "$smoke_status" in
    pass)
      echo "=== All smokes passed (${SMOKE_PASSED}/${SMOKE_TOTAL}) ==="
      ;;
    partial)
      echo "=== Some smokes passed (${SMOKE_PASSED}/${SMOKE_TOTAL}) ==="
      ;;
    *)
      echo "=== All smokes failed (0/${SMOKE_TOTAL}) ==="
      ;;
  esac
  echo "Pair: LPS ${LPS_REF} x LBC ${LBC_REF}"

  exit "$smoke_rc"
}

main "$@"
