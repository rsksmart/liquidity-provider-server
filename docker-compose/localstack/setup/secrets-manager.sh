#!/bin/bash

BTC_WALLET_PASSWORD=test-password
RSK_ENCRYPTED_JSON_PASSWORD=test

awslocal secretsmanager create-secret --name FlyoverTestEnv/LPS-LOCAL-BTC-WALLET-PASSWORD --secret-string $BTC_WALLET_PASSWORD
awslocal secretsmanager create-secret --name FlyoverTestEnv/LPS-LOCAL-PASSWORD --secret-string $RSK_ENCRYPTED_JSON_PASSWORD
awslocal secretsmanager create-secret --name FlyoverTestEnv/LPS-LOCAL-KEY --secret-string file:///tmp/local-key.json
rm /tmp/local-key.json