#!/usr/bin/env sh
set -eu

# Generates the MongoDB internal-auth keyfile used by --keyFile when running
# as a replica set. Idempotent: if the keyfile already exists, only the
# strict ownership and permissions mongod requires are re-applied.
#
# The keyfile is stored on a Docker-managed named volume that this service
# and the mongodb service share. It is never written to the host filesystem.

KEYFILE=/keyfile/mongo-keyfile
MONGO_UID=999
MONGO_GID=999

install -d -m 0700 "$(dirname "$KEYFILE")"

if [ -s "$KEYFILE" ]; then
  echo "Keyfile already present at $KEYFILE; reusing."
else
  echo "Generating MongoDB keyfile at $KEYFILE..."
  openssl rand -base64 756 > "$KEYFILE"
fi

chown "$MONGO_UID:$MONGO_GID" "$KEYFILE" "$(dirname "$KEYFILE")"
chmod 0400 "$KEYFILE"

echo "Keyfile ready:"
ls -l "$KEYFILE"
