#!/bin/sh

# Generates the RSA key pair used by the LPS to authenticate as a Fireblocks API user.
#   - fireblocks_secret.key: private key, loaded by the LPS from the secrets source. Never share it.
#   - fireblocks.csr: certificate signing request, sent to the Fireblocks workspace admin to create the API user.
#
# Optional env vars:
#   FORCE  set to "true" to overwrite an existing key pair. The API user created with the previous CSR
#          stops working, so a new CSR must be sent to the admin.

set -e

OUTPUT_DIR=${1:-docker-compose/local/localstack}
KEY_FILE="$OUTPUT_DIR/fireblocks_secret.key"
CSR_FILE="$OUTPUT_DIR/fireblocks.csr"

if [ -e "$KEY_FILE" ] && [ "$FORCE" != "true" ]; then
  echo "$KEY_FILE already exists. A new key invalidates the Fireblocks API user created with the previous CSR."
  echo "Run again with FORCE=true to overwrite it."
  exit 1
fi

mkdir -p "$OUTPUT_DIR"
# Fireblocks only uses the public key from the CSR, so the subject is left empty
openssl req -new -newkey rsa:4096 -nodes -keyout "$KEY_FILE" -out "$CSR_FILE" -subj "/"
chmod 600 "$KEY_FILE"

echo "Private key: $KEY_FILE (keep it, never share it)"
echo "CSR: $CSR_FILE (send it to the Fireblocks workspace admin to create your API user)"
