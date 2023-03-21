#!/bin/bash

set -e

echo "Deploying contracts to RskJ..."

RSK_NETWORK="rsk${LPS_STAGE^}"
echo "RSK network: $RSK_NETWORK"

npx truffle deploy --network $RSK_NETWORK

echo "Deployment succeeded"

LBC_ADDR=$(cat ./config.json | jq -r '.rskRegtest.LiquidityBridgeContract.address')
echo "LBC_ADDR=$LBC_ADDR"
