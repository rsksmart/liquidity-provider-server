#!/bin/bash

WALLET_SECRET=FlyoverTestEnv/LPS-LOCAL-KEY
PASSWORD_SECRET=FlyoverTestEnv/LPS-LOCAL-PASSWORD
RSK_ENCRYPTED_JSON_PASSWORD="test"
HOT_WALLET_FILE=file:///tmp/local-key.json

awslocal secretsmanager create-secret --name $PASSWORD_SECRET --secret-string $RSK_ENCRYPTED_JSON_PASSWORD
awslocal secretsmanager create-secret --name $WALLET_SECRET --secret-string $HOT_WALLET_FILE
rm /tmp/local-key.json