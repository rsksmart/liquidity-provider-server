#!/usr/bin/env bash
set -euo pipefail

# Initiates the single-node replica set rs0 against the mongodb container.
# Idempotent: if the replica set is already configured, the script only
# waits for PRIMARY and exits. Safe to run on every `docker compose up`.

HOST=${MONGO_HOST:-mongo01}
PORT=${MONGO_PORT:-27017}
RS=${REPLICA_SET:-rs0}
USER=${MONGO_INITDB_ROOT_USERNAME:?MONGO_INITDB_ROOT_USERNAME is required}
PASS=${MONGO_INITDB_ROOT_PASSWORD:?MONGO_INITDB_ROOT_PASSWORD is required}

mongo_eval() {
  mongo --host "$HOST" --port "$PORT" --quiet \
    -u "$USER" -p "$PASS" --authenticationDatabase admin \
    --eval "$1"
}

echo "Waiting for MongoDB at $HOST:$PORT to accept authenticated commands..."
# On a fresh data volume the mongo entrypoint creates the root user during
# its bootstrap phase, then re-execs mongod with our --replSet command. This
# loop covers that whole window.
until mongo_eval 'db.adminCommand({ ping: 1 }).ok' >/dev/null 2>&1; do
  sleep 2
done

echo "MongoDB is reachable. Checking replica set state..."
RS_OK=$(mongo_eval 'try { rs.status().ok } catch (e) { print(0) }' 2>/dev/null | tr -d '[:space:]' || echo 0)

if [ "$RS_OK" = "1" ]; then
  echo "Replica set $RS already initialized."
else
  echo "Initializing replica set $RS with member $HOST:$PORT..."
  mongo_eval "rs.initiate({_id: '$RS', members: [{_id: 0, host: '$HOST:$PORT'}]})"
fi

echo "Waiting for $HOST:$PORT to become PRIMARY..."
until [ "$(mongo_eval 'rs.isMaster().ismaster' 2>/dev/null | tr -d '[:space:]')" = "true" ]; do
  sleep 2
done

echo "Replica set $RS is ready and PRIMARY is $HOST:$PORT."
