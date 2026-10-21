#!/bin/bash

BTC_WALLET_PASSWORD=test-password
RSK_ENCRYPTED_JSON_PASSWORD="test"
FIREBLOCKS_API_KEY=fireblocks-api-key-placeholder

awslocal secretsmanager create-secret --name FlyoverTestEnv/LPS-LOCAL-BTC-WALLET-PASSWORD --secret-string $BTC_WALLET_PASSWORD
awslocal secretsmanager create-secret --name FlyoverTestEnv/LPS-LOCAL-PASSWORD --secret-string $RSK_ENCRYPTED_JSON_PASSWORD
awslocal secretsmanager create-secret --name FlyoverTestEnv/LPS-LOCAL-KEY --secret-string file:///tmp/local-key.json
rm /tmp/local-key.json
awslocal secretsmanager create-secret --name FlyoverTestEnv/LPS-LOCAL-FIREBLOCKS-API-KEY --secret-string $FIREBLOCKS_API_KEY
awslocal secretsmanager create-secret --name FlyoverTestEnv/LPS-LOCAL-FIREBLOCKS-SECRET-KEY --secret-string file:///tmp/fireblocks_secret.key
rm /tmp/fireblocks_secret.key
