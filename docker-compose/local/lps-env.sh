#!/bin/sh

set -e

if [ -z "${LPS_STAGE}" ]; then
  echo "LPS_STAGE is not set. Exit 1"
  exit 1
elif [ "$LPS_STAGE" = "regtest" ]; then
  ENV_FILE=".env.regtest"
elif [ "$LPS_STAGE" = "testnet" ]; then
  ENV_FILE=".env.testnet"
else
  echo "Invalid LPS_STAGE: $LPS_STAGE"
  exit 1
fi

if [ -z "${LPS_UID}" ]; then
  export LPS_UID=$(id -u)
  if [ "$LPS_UID" = "0" ]; then
    echo "Please set LPS_UID env var or run as a non-root user"
    exit 1
  fi
fi

echo "LPS_STAGE: $LPS_STAGE; ENV_FILE: $ENV_FILE; LPS_UID: $LPS_UID"

SCRIPT_CMD=$1
if [ -z "${SCRIPT_CMD}" ]; then
  echo "Command is not provided"
  exit 1
elif [ "$SCRIPT_CMD" = "up" ]; then
  echo "Starting LPS env up..."
elif [ "$SCRIPT_CMD" = "down" ]; then
  echo "Shutting LPS env down..."
  docker-compose --env-file "$ENV_FILE" -f docker-compose.yml -f docker-compose.lps.yml down
  exit 0
elif [ "$SCRIPT_CMD" = "build" ]; then
  echo "Building LPS env..."
  docker-compose --env-file "$ENV_FILE" -f docker-compose.yml -f docker-compose.lbc-deployer.yml -f docker-compose.lps.yml build
  exit 0
elif [ "$SCRIPT_CMD" = "stop" ]; then
  echo "Stopping LPS env..."
  docker-compose --env-file "$ENV_FILE" -f docker-compose.yml -f docker-compose.lps.yml stop
  exit 0
elif [ "$SCRIPT_CMD" = "ps" ]; then
  echo "List of running services:"
  docker-compose --env-file "$ENV_FILE" -f docker-compose.yml -f docker-compose.lps.yml ps
  exit 0
elif [ "$SCRIPT_CMD" = "deploy" ]; then
  echo "Stopping LPS..."
  docker-compose --env-file "$ENV_FILE" -f docker-compose.yml -f docker-compose.lps.yml stop lps
  echo "Building LPS..."
  docker-compose --env-file "$ENV_FILE" -f docker-compose.yml -f docker-compose.lps.yml build lps
  echo "Starting LPS..."
  docker-compose --env-file "$ENV_FILE" -f docker-compose.yml -f docker-compose.lps.yml start lps
  exit 0
elif [ "$SCRIPT_CMD" = "import-rsk-db" ]; then
  echo "Importing rsk db..."
  docker-compose --env-file "$ENV_FILE" run -d rskj java -Xmx6g -Drpc.providers.web.http.bind_address=0.0.0.0 -Drpc.providers.web.http.hosts.0=localhost -Drpc.providers.web.http.hosts.1=rskj -cp rskj-core.jar co.rsk.Start --${LPS_STAGE} --import
  exit 0
elif [ "$SCRIPT_CMD" = "start-bitcoind" ]; then
  echo "Starting bitcoind..."
  docker-compose --env-file "$ENV_FILE" -f docker-compose.yml up -d bitcoind
  exit 0
else
  echo "Invalid command: $SCRIPT_CMD"
  exit 1
fi

BTCD_HOME="${BTCD_HOME:-./volumes/bitcoind}"
RSKJ_HOME="${RSKJ_HOME:-./volumes/rskj}"
LPS_HOME="${LPS_HOME:-./volumes/lps}"

[ -d "$BTCD_HOME" ] || mkdir -p "$BTCD_HOME" && chown "$LPS_UID" "$BTCD_HOME"
[ -d "$RSKJ_HOME" ] || mkdir -p "$RSKJ_HOME/db" && mkdir -p "$RSKJ_HOME/logs" && chown -R "$LPS_UID" "$RSKJ_HOME"
[ -d "$LPS_HOME" ] || mkdir -p "$LPS_HOME/db" && mkdir -p "$LPS_HOME/logs" && chown -R "$LPS_UID" "$LPS_HOME"

echo "LPS_UID: $LPS_UID; BTCD_HOME: '$BTCD_HOME'; RSKJ_HOME: '$RSKJ_HOME'; LPS_HOME: '$LPS_HOME'"

# start bitcoind and RSKJ dependant services
docker-compose --env-file "$ENV_FILE" up -d

# read env vars
. ./"$ENV_FILE"

echo "Waiting for RskJ to be up and running..."
while true
do
  sleep 3
  curl -s "http://127.0.0.1:4444" -X POST -H "Content-Type: application/json" -d '{"jsonrpc":"2.0","method":"eth_chainId","params": [],"id":1}' \
    && echo "RskJ is up and running" \
    && break
done

echo "Waiting for Bitcoind to be up and running..."
while true
do
  sleep 3
  curl -s "http://127.0.0.1:5555" -X POST --user "$BTCD_RPC_USER:$BTCD_RPC_PASS" -H "Content-Type: application/json" -d '{"jsonrpc": "1.0", "method": "getnetworkinfo", "params": [], "id":"1"}' | grep "\"result\":{" \
    && echo "Bitcoind is up and running" \
    && break
done

if [ "$LPS_STAGE" = "regtest" ]; then
  # pre-fund provider in regtest, if needed
  LIQUIDITY_PROVIDER_RSK_ADDR_LINE=$(cat "$ENV_FILE" | grep LIQUIDITY_PROVIDER_RSK_ADDR | head -n 1 | tr -d '\r')
  LIQUIDITY_PROVIDER_RSK_ADDR="${LIQUIDITY_PROVIDER_RSK_ADDR_LINE#"LIQUIDITY_PROVIDER_RSK_ADDR="}"
  PROVIDER_TX_COUNT=$(curl -s -X POST "http://127.0.0.1:4444" -H "Content-Type: application/json" -d "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getTransactionCount\",\"params\": [\"$LIQUIDITY_PROVIDER_RSK_ADDR\",\"latest\"],\"id\":1}" | jq -r ".result")
  if [ "$PROVIDER_TX_COUNT" = "0x0" ]; then
    echo "Transferring funds to $LIQUIDITY_PROVIDER_RSK_ADDR..."

    TX_HASH=$(curl -s -X POST "http://127.0.0.1:4444" -H "Content-Type: application/json" -d "{\"jsonrpc\":\"2.0\",\"method\":\"eth_sendTransaction\",\"params\": [{\"from\": \"0xcd2a3d9f938e13cd947ec05abc7fe734df8dd826\", \"to\": \"$LIQUIDITY_PROVIDER_RSK_ADDR\", \"value\": \"0x8AC7230489E80000\"}],\"id\":1}" | jq -r ".result")
    echo "Result: $TX_HASH"
    sleep 10
  else
    echo "No need to fund the '$LIQUIDITY_PROVIDER_RSK_ADDR' provider. Nonce: $PROVIDER_TX_COUNT"
  fi

  # if [ -z "${LBC_ADDR}" ]; then
    # echo "LBC_ADDR is not set. Deploying LBC contract..."

    # deploy LBC contracts to RSKJ
    # LBC_ADDR_LINE=$(docker-compose --env-file "$ENV_FILE" -f docker-compose.yml -f docker-compose.lbc-deployer.yml run --rm lbc-deployer bash deploy-lbc.sh | grep LBC_ADDR | head -n 1 | tr -d '\r')
    # export LBC_ADDR="0x8901a2Bbf639bFD21A97004BA4D7aE2BD00B8DA8"
  # fi
  export LBC_ADDR="0x8901a2Bbf639bFD21A97004BA4D7aE2BD00B8DA8"
fi

if [ -z "${LBC_ADDR}" ]; then
  docker-compose down
  echo "LBC_ADDR is not set up. Exit"
  exit 1
fi
echo "LBC deployed at $LBC_ADDR"

# start LPS
docker-compose --env-file "$ENV_FILE" -f docker-compose.yml -f docker-compose.lps.yml up -d lps
