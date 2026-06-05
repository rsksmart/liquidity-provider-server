#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
COMPOSE_FILE="$SCRIPT_DIR/../docker-compose/docker-compose.yml"

cleanup() {
  echo "Stopping MongoDB..."
  docker compose -f "$COMPOSE_FILE" down -v || true
}
trap cleanup EXIT

if ! command -v gotestsum >/dev/null 2>&1; then
  echo "Error: gotestsum not found in PATH. Run 'make tools' to install it."
  exit 1
fi

# Start MongoDB
echo "Starting MongoDB..."
docker compose -f "$COMPOSE_FILE" up -d --wait

# Run tests
echo "Running database integration tests..."
cd "$PROJECT_ROOT"
MONGODB_HOST=localhost \
MONGODB_PORT=27018 \
MONGODB_USER=test \
MONGODB_PASSWORD=test \
gotestsum --format testname -- -tags integration -timeout 5m ./test/mongodb/...
