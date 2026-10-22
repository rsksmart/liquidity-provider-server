#!/bin/bash

BTC_WALLET_PASSWORD=test-password
RSK_ENCRYPTED_JSON_PASSWORD="test"

awslocal secretsmanager create-secret --name FlyoverTestEnv/LPS-LOCAL-BTC-WALLET-PASSWORD --secret-string $BTC_WALLET_PASSWORD
awslocal secretsmanager create-secret --name FlyoverTestEnv/LPS-LOCAL-PASSWORD --secret-string $RSK_ENCRYPTED_JSON_PASSWORD
awslocal secretsmanager create-secret --name FlyoverTestEnv/LPS-LOCAL-KEY --secret-string file:///tmp/local-key.json
rm /tmp/local-key.json
# Fireblocks secrets are only needed when running with WALLET=fireblocks
if [ -n "$FIREBLOCKS_API_KEY" ]; then
  awslocal secretsmanager create-secret --name FlyoverTestEnv/LPS-LOCAL-FIREBLOCKS-API-KEY --secret-string "$FIREBLOCKS_API_KEY"
else
  echo "FIREBLOCKS_API_KEY not set, skipping FlyoverTestEnv/LPS-LOCAL-FIREBLOCKS-API-KEY"
fi
if [ -f /tmp/fireblocks_secret.key ]; then
  awslocal secretsmanager create-secret --name FlyoverTestEnv/LPS-LOCAL-FIREBLOCKS-SECRET-KEY --secret-string file:///tmp/fireblocks_secret.key
  rm /tmp/fireblocks_secret.key
else
  echo "fireblocks_secret.key not found, skipping FlyoverTestEnv/LPS-LOCAL-FIREBLOCKS-SECRET-KEY"
fi
