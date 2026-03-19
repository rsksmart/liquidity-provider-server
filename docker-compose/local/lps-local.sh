#!/bin/bash

set_defaults() {
  # local env defaults
  : "${LPS_UID=$(id -u)}"  ; export LPS_UID
  export ENABLE_MANAGEMENT_API=true
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
docker compose --env-file "$ENV_FILE" up -d --wait

### Funding wallets ###
if [[ "$FUND_WALLETS" == "true" ]]; then
  echo "Funding wallets..."
  docker compose -f docker-compose.yml -f wallet-funder/docker-compose.funder.yml --env-file "$ENV_FILE" up --wait
fi

### Contract deployment ###
if [[ "$DEPLOY_CONTRACTS" == "true" ]]; then
  echo "Deploying contracts..."
  docker compose -f docker-compose.yml -f wallet-funder/docker-compose.funder.yml -f lbc-deployer/docker-compose.lbc-deployer.yml --env-file "$ENV_FILE" up -d --wait
  docker wait lbc-deployer
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
fi

### LPS (always runs) ###
docker compose -f docker-compose.yml -f lps/docker-compose.lps-local.yml --env-file "$ENV_FILE" up -d --wait

if [[ "$CREATE_MONITORING" == "true" ]]; then
  docker compose -f docker-compose.yml -f metrics/docker-compose.metrics.yml --env-file "$ENV_FILE" up --wait
fi

echo ""
echo "============================================"
echo "✓ LPS environment is ready!"
echo "  LPS API:    http://localhost:8080"
echo "  Health:     http://localhost:8080/health"
echo "  Management: http://localhost:8080/management"
[[ "$CREATE_MONITORING" == "true" ]] && echo "  Grafana:    http://localhost:3000"
echo "============================================"
