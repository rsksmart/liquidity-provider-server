#!/bin/sh

set -e

if [ -z "${LBC_ADDR}" ]; then
  echo "LBC_ADDR is not set. Waiting for contract to be deployed"

  sleep 5
  while [ ! -f /lps/contracts/LiquidityBridgeContract.json ]; do sleep 1; done
  sleep 5

  LBC_ADDR=$(cat /lps/contracts/LiquidityBridgeContract.json | jq -r '.networks."33".address')
fi

echo "LBC_ADDR: $LBC_ADDR"

PROVIDER_BALANCE=$(curl -X POST "http://rskj:4444" -H "Content-Type: application/json" -d "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getBalance\",\"params\": [\"$LIQUIDITY_PROVIDER_RSK_ADDR\",\"latest\"],\"id\":1}" | jq -r ".result")
PROVIDER_TX_COUNT=$(curl -X POST "http://rskj:4444" -H "Content-Type: application/json" -d "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getTransactionCount\",\"params\": [\"$LIQUIDITY_PROVIDER_RSK_ADDR\",\"latest\"],\"id\":1}" | jq -r ".result")
if [[ "$PROVIDER_BALANCE" = "0x0" && "$PROVIDER_TX_COUNT" = "0x0" ]]; then
  echo "Transferring funds to $LIQUIDITY_PROVIDER_RSK_ADDR..."

  TX_HASH=$(curl -X POST "http://rskj:4444" -H "Content-Type: application/json" -d "{\"jsonrpc\":\"2.0\",\"method\":\"eth_sendTransaction\",\"params\": [{\"from\": \"0xcd2a3d9f938e13cd947ec05abc7fe734df8dd826\", \"to\": \"$LIQUIDITY_PROVIDER_RSK_ADDR\", \"value\": \"0x1000000\"}],\"id\":1}" | jq -r ".result")
  echo "Result: $TX_HASH"
  sleep 10
else
  echo "No need to fund the '$LIQUIDITY_PROVIDER_RSK_ADDR' provider. Balance: $PROVIDER_BALANCE, nonce: $PROVIDER_TX_COUNT"
fi

if [[ ! -d ./geth_keystore ]]; then
  mkdir ./geth_keystore && echo $LIQUIDITY_PROVIDER_RSK_KEY > ./geth_keystore/test
fi
if [[ ! -f ./pwd.txt ]]; then
  echo $LIQUIDITY_PROVIDER_RSK_KEY_PASS > ./pwd.txt
fi

echo "Starting LP Server..."
/liquidity-provider-server
