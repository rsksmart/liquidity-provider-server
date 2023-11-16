#!/bin/sh

set -e

echo "Detected LPS_STAGE: $LPS_STAGE, LBC_ADDR: $LBC_ADDR, BTCD_RPC_USER: $BTCD_RPC_USER, RSK_CHAIN_ID: $RSK_CHAIN_ID"

echo "Testing if we have a default wallet"

curl -s "http://bitcoind01:5555" --user "$BTCD_RPC_USER:$BTCD_RPC_PASS" -H "Content-Type: application/json" -d '{"jsonrpc": "1.0", "method": "getwalletinfo", "params": [], "id":"getwallet"}' | grep "{\"result\":null,\"error\":{\"code\":-18" \
  && echo "No default wallet" \
  && echo "Creating wallet" \
  && curl -s "http://bitcoind01:5555" --user "$BTCD_RPC_USER:$BTCD_RPC_PASS" -H "Content-Type: application/json" -d '{"jsonrpc": "1.0", "method": "createwallet", "params": ["main", false, false, "test-password", true, false, true], "id":"createwallet"}'

echo "Starting LP Server..."
liquidity-provider-server
