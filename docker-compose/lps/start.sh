#!/bin/sh

set -e

echo "Detected LPS_STAGE: $LPS_STAGE, LBC_ADDR: $LBC_ADDR, BTCD_RPC_USER: $BTCD_RPC_USER, RSK_CHAIN_ID: $RSK_CHAIN_ID"

jq \
  --arg lbcAddr "$LBC_ADDR" \
  --arg btcUser "$BTCD_RPC_USER" \
  --arg btcPass "$BTCD_RPC_PASS" \
  --arg btcNetwork "$LPS_STAGE" \
  --arg chainId "$RSK_CHAIN_ID" \
  '.rsk.lbcAddr=$lbcAddr | .btc.username=$btcUser | .btc.password=$btcPass | .btc.network=$btcNetwork | .provider.chainId=($chainId | tonumber)' \
  config.json > updated_config.json && mv updated_config.json config.json

if [[ ! -d ./geth_keystore ]]; then
  mkdir ./geth_keystore && echo $LIQUIDITY_PROVIDER_RSK_KEY > ./geth_keystore/key
fi
if [[ ! -f ./pwd.txt ]]; then
  echo $LIQUIDITY_PROVIDER_RSK_KEY_PASS > ./pwd.txt
fi

echo "Starting LP Server..."
liquidity-provider-server
