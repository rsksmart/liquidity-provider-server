#!/bin/bash

set_defaults() {
  # local env defaults
  : "${LPS_UID=$(id -u)}"  ; export LPS_UID
  export ENABLE_MANAGEMENT_API=true
  export ENABLE_MANAGEMENT_UI_NEXT=true
  export LPS_STAGE=regtest
}

if [[ "$1" == "--help" || "$1" == "-h" ]]; then
  echo "Usage: $0 [OPTIONS]"
  echo ""
  echo "Options:"
  echo "  -r, --reset    Reset the environment by stopping containers and removing volumes"
  echo "  -h, --help     Show this help message and exit"
  exit 0
fi

if [[ "$1" == "--reset" || "$1" == "-r" ]]; then
  echo "Resetting environment..."
  docker compose -p local down
  rm -rf volumes
  rm -f .env.regtest # delete default
fi

: "${ENV_FILE=".env.regtest"}"  ; export ENV_FILE
if [ ! -f "$ENV_FILE" ]; then
  echo "Creating $ENV_FILE from sample-config.env..."; cp ../../sample-config.env "$ENV_FILE"
else
  echo "Using existing $ENV_FILE"
fi
set -a
# shellcheck disable=SC1090
source "$ENV_FILE"
set +a
set_defaults

### Create base (always runs) ###
docker compose --progress plain --env-file "$ENV_FILE" up -d --wait

### Funding wallets ###
if [[ "$FUND_WALLETS" == "true" ]]; then
  echo "Funding wallets..."
  docker compose -f docker-compose.yml -f wallet-funder/docker-compose.funder.yml --env-file "$ENV_FILE" up --wait
fi

### Contract deployment ###
if [[ "$DEPLOY_CONTRACTS" == "true" ]]; then
  echo "Deploying contracts..."
  # LBC_IMAGE is a pinned digest that never changes between runs, so compose would
  # treat a previous run's "service_completed_successfully" as still satisfied and
  # skip re-running the deployer. After a chain reset (rm -rf volumes) that leaves the
  # contracts undeployed. Remove the one-shot containers so they always run fresh.
  docker rm -f lbc-deployer lps-approver >/dev/null 2>&1 || true
  docker compose -f docker-compose.yml -f wallet-funder/docker-compose.funder.yml -f lbc-deployer/docker-compose.lbc-deployer.yml --env-file "$ENV_FILE" up -d --wait
  EXIT_CODE=$(docker wait lbc-deployer)
  # docker wait returns 0 for a container that was created but never started, so also
  # require that the deployer actually ran before trusting the exit code.
  STARTED_AT=$(docker inspect -f '{{.State.StartedAt}}' lbc-deployer 2>/dev/null)
  if [ "$EXIT_CODE" != "0" ] || [ "$STARTED_AT" = "0001-01-01T00:00:00Z" ]; then
    echo "ERROR: Contract deployment failed (exit code $EXIT_CODE, started $STARTED_AT)"
    exit 1
  fi
  echo "Contracts deployed"
  set -a
  # shellcheck disable=SC1090
  source "$ENV_FILE"
  set +a
  set_defaults
fi

### Powpeg ###
if [[ "$CREATE_POWPEG" == "true" ]]; then
  docker compose -f docker-compose.yml -f powpeg/docker-compose.powpeg.yml --env-file "$ENV_FILE" up -d

  ### Federation segwit migration ###
  if [[ "$MIGRATE_FEDERATION" == "true" ]]; then
    echo "Migrating federation to segwit..."
    docker compose -f docker-compose.yml -f fed-migrator/docker-compose.fed-migrator.yml --env-file "$ENV_FILE" up --wait
    EXIT_CODE=$(docker wait fed-migrator)
    if [ "$EXIT_CODE" != "0" ]; then
      echo "ERROR: Federation migration failed (exit code $EXIT_CODE)"
      exit 1
    fi
    echo "Federation migrated to segwit"
    set -a
    
    source "$ENV_FILE"
    set +a
    set_defaults
  fi
fi

### LPS (always runs) ###
# lps01 may briefly crash on first boot during a BTC
# wallet rescan. lps-configurer polls /health internally and signals completion via
# its exit code
# TODO: Change the fatal.log when lps01 crashes
docker compose --progress plain -f docker-compose.yml -f lps/docker-compose.lps-local.yml --env-file "$ENV_FILE" up -d
echo "Configuring LPS..."
EXIT_CODE=$(docker wait lps-configurer)
if [ "$EXIT_CODE" != "0" ]; then
  echo "ERROR: LPS configuration failed (exit code $EXIT_CODE)"
  exit 1
fi
echo "LPS configured"

if [[ "$CREATE_MONITORING" == "true" ]]; then
  docker compose -f docker-compose.yml -f monitoring/docker-compose.monitoring.yml --env-file "$ENV_FILE" up --wait
fi

echo ""
echo "============================================"
echo "✓ LPS environment is ready!"
echo "  LPS API:    http://localhost:8080"
echo "  Health:     http://localhost:8080/health"
echo "  Management: http://localhost:8080/management"
[[ "$CREATE_MONITORING" == "true" ]] && echo "  Grafana:    http://localhost:3000"
echo "============================================"
