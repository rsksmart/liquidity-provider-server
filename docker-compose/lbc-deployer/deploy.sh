#!/bin/bash

set -e

if [[ ! -z "${LBC_ADDR}" ]]; then
  echo "LBC_ADDR is set to $LBC_ADDR. No need do deploy the contract"

  exit 0
fi

CONTRACTS_DIR="/contracts"
if [[ -d $CONTRACTS_DIR && "$(ls -A $CONTRACTS_DIR)" ]]; then
    echo "Cleaning up $CONTRACTS_DIR directory..."
    rm $CONTRACTS_DIR/*.*
fi

echo "Deploying contracts to RskJ..."

cd /code/lbc

npx truffle deploy --network rskRegtest

cp -r ./build/contracts /

echo "Deployment succeeded"
