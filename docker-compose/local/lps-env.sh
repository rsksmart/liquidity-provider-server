#!/bin/bash

set -e

COMMIT_HASH=$(git rev-parse HEAD)
COMMIT_TAG=$(git describe --exact-match --tags || echo "")
export COMMIT_HASH
export COMMIT_TAG

# Detect OS
OS_TYPE="$(uname)"

if [[ "$OS_TYPE" == "Darwin" ]]; then
    # macOS
    echo "Running on macOS"
    SED_INPLACE=("sed" "-i" "")
elif [[ "$OS_TYPE" == "Linux" ]]; then
    # Assume Ubuntu or other Linux
    echo "Running on Linux"
    SED_INPLACE=("sed" "-i")
else
    echo "Unsupported OS: $OS_TYPE"
    exit 1
fi

if [ -z "${LPS_STAGE}" ]; then
  echo "LPS_STAGE is not set. Exit 1"
  exit 1
elif [ "$LPS_STAGE" = "regtest" ]; then
  cp ../../sample-config.env .env.regtest
  ENV_FILE=".env.regtest"
elif [ "$LPS_STAGE" = "testnet" ]; then
  ENV_FILE=".env.testnet"
else
  echo "Invalid LPS_STAGE: $LPS_STAGE"
  exit 1
fi

if [ -z "${LPS_UID}" ]; then
  LPS_UID=$(id -u)
  export LPS_UID
  if [ "$LPS_UID" = "0" ]; then
    echo "Please set LPS_UID env var or run as a non-root user"
    exit 1
  fi
fi

echo "LPS_STAGE: $LPS_STAGE; ENV_FILE: $ENV_FILE; LPS_UID: $LPS_UID"

# Force Management API to be enabled
if [ -f "$ENV_FILE" ]; then
  "${SED_INPLACE[@]}" 's/^ENABLE_MANAGEMENT_API=.*/ENABLE_MANAGEMENT_API=true/' "$ENV_FILE"
fi

SCRIPT_CMD=$1
if [ -z "${SCRIPT_CMD}" ]; then
  echo "Command is not provided"
  exit 1
elif [ "$SCRIPT_CMD" = "up" ]; then
  echo "Starting LPS env up..."
elif [ "$SCRIPT_CMD" = "down" ]; then
  echo "Shutting LPS env down..."
  docker compose --env-file "$ENV_FILE" -f docker-compose.yml -f docker-compose.lps.yml down
  exit 0
elif [ "$SCRIPT_CMD" = "build" ]; then
  echo "Building LPS env..."
  docker compose --env-file "$ENV_FILE" -f docker-compose.yml -f docker-compose.lbc-deployer.yml -f docker-compose.lps.yml build
  exit 0
elif [ "$SCRIPT_CMD" = "stop" ]; then
  echo "Stopping LPS env..."
  docker compose --env-file "$ENV_FILE" -f docker-compose.yml -f docker-compose.lps.yml stop
  exit 0
elif [ "$SCRIPT_CMD" = "ps" ]; then
  echo "List of running services:"
  docker compose --env-file "$ENV_FILE" -f docker-compose.yml -f docker-compose.lps.yml ps
  exit 0
elif [ "$SCRIPT_CMD" = "deploy" ]; then
  echo "Stopping LPS..."
  docker compose --env-file "$ENV_FILE" -f docker-compose.yml -f docker-compose.lps.yml stop lps
  echo "Building LPS..."
  docker compose --env-file "$ENV_FILE" -f docker-compose.yml -f docker-compose.lps.yml build lps
  echo "Starting LPS..."
  docker compose --env-file "$ENV_FILE" -f docker-compose.yml -f docker-compose.lps.yml start lps
  exit 0
elif [ "$SCRIPT_CMD" = "import-rsk-db" ]; then
  echo "Importing rsk db..."
  docker compose --env-file "$ENV_FILE" run -d rskj java -Xmx6g -Drpc.providers.web.http.bind_address=0.0.0.0 -Drpc.providers.web.http.hosts.0=localhost -Drpc.providers.web.http.hosts.1=rskj -cp rskj-core.jar co.rsk.Start --"${LPS_STAGE}" --import
  exit 0
elif [ "$SCRIPT_CMD" = "start-bitcoind" ]; then
  echo "Starting bitcoind..."
  docker compose --env-file "$ENV_FILE" -f docker-compose.yml up -d bitcoind
  exit 0
else
  echo "Invalid command: $SCRIPT_CMD"
  exit 1
fi

BTCD_HOME="${BTCD_HOME:-./volumes/bitcoind}"
POWPEG_PEGIN_HOME="${POWPEG_PEGIN_HOME:-./volumes/powpeg/pegin}"
POWPEG_PEGOUT_HOME="${POWPEG_PEGOUT_HOME:-./volumes/powpeg/pegout}"
RSKJ_HOME="${RSKJ_HOME:-./volumes/rskj}"
LPS_HOME="${LPS_HOME:-./volumes/lps}"
MONGO_HOME="${MONGO_HOME:-./volumes/mongo}"
LOCALSTACK_HOME="${LOCALSTACK_HOME:-./volumes/localstack}"

[ -d "$BTCD_HOME" ] || mkdir -p "$BTCD_HOME" && chown "$LPS_UID" "$BTCD_HOME"
[ -d "$RSKJ_HOME" ] || mkdir -p "$RSKJ_HOME/db" && mkdir -p "$RSKJ_HOME/logs" && chown -R "$LPS_UID" "$RSKJ_HOME"
[ -d "$POWPEG_PEGIN_HOME" ] || mkdir -p "$POWPEG_PEGIN_HOME/db" && mkdir -p "$POWPEG_PEGIN_HOME/logs" && chown -R "$LPS_UID" "$POWPEG_PEGIN_HOME" && chmod -R 777 "$POWPEG_PEGIN_HOME"
[ -d "$POWPEG_PEGOUT_HOME" ] || mkdir -p "$POWPEG_PEGOUT_HOME/db" && mkdir -p "$POWPEG_PEGOUT_HOME/logs" && chown -R "$LPS_UID" "$POWPEG_PEGOUT_HOME" && chmod -R 777 "$POWPEG_PEGOUT_HOME"
[ -d "$LPS_HOME" ] || mkdir -p "$LPS_HOME/logs" && chmod -R 777 "$LPS_HOME"
[ -d "$MONGO_HOME" ] || mkdir -p "$MONGO_HOME/db" && chown -R "$LPS_UID" "$MONGO_HOME"
[ -d "$LOCALSTACK_HOME" ] || mkdir -p "$LOCALSTACK_HOME/db" && mkdir -p "$LOCALSTACK_HOME/logs" && chown -R "$LPS_UID" "$LOCALSTACK_HOME"

echo "LPS_UID: $LPS_UID; BTCD_HOME: '$BTCD_HOME'; RSKJ_HOME: '$RSKJ_HOME'; LPS_HOME: '$LPS_HOME'; MONGO_HOME: '$MONGO_HOME'; POWPEG_PEGIN_HOME: '$POWPEG_PEGIN_HOME'; POWPEG_PEGOUT_HOME: '$POWPEG_PEGOUT_HOME'; LOCALSTACK_HOME: '$LOCALSTACK_HOME'"

# start bitcoind and RSKJ dependant services
docker compose --env-file "$ENV_FILE" up -d bitcoind rskj mongodb localstack

# shellcheck disable=SC1090
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
  curl -s "http://127.0.0.1:5555" -X POST --user "$BTC_USERNAME:$BTC_PASSWORD" -H "Content-Type: application/json" -d '{"jsonrpc": "1.0", "method": "getnetworkinfo", "params": [], "id":"1"}' | grep "\"result\":{" \
    && echo "Bitcoind is up and running" \
    && break
done

curl -s "http://127.0.0.1:5555" --user "$BTC_USERNAME:$BTC_PASSWORD" -H "Content-Type: application/json" -d '{"jsonrpc": "1.0", "method": "getwalletinfo", "params": [], "id":"getwallet"}' | grep "{\"result\":null,\"error\":{\"code\":-18" \
  && echo "No default wallet" \
  && echo "Creating wallet" \
  && curl -s "http://127.0.0.1:5555" --user "$BTC_USERNAME:$BTC_PASSWORD" -H "Content-Type: application/json" -d '{"jsonrpc": "1.0", "method": "createwallet", "params": ["main", false, false, "test-password", true, false, true], "id":"createwallet"}' \
  && curl -s "http://127.0.0.1:5555" --user "$BTC_USERNAME:$BTC_PASSWORD" -H "Content-Type: application/json" -d '{"jsonrpc": "1.0", "method": "getnewaddress", "params": ["main"], "id":"getnewaddress"}' \
  | jq .result | xargs -I ADDRESS curl -s "http://127.0.0.1:5555" --user "$BTC_USERNAME:$BTC_PASSWORD" -H "Content-Type: application/json" -d '{"jsonrpc": "1.0", "method": "generatetoaddress", "params": [500, "ADDRESS"], "id":"generatetoaddress"}' \
  && echo "Wallet created and generated 500 blocks" \
  && curl -s "http://127.0.0.1:5555" --user "$BTC_USERNAME:$BTC_PASSWORD" -H "Content-Type: application/json" -d '{"jsonrpc": "1.0", "method": "createwallet", "params": ["rsk-wallet", true, true, "", false, false, true], "id":"createwallet"}' \
  && curl -s "http://127.0.0.1:5555/wallet/rsk-wallet" --user "$BTC_USERNAME:$BTC_PASSWORD" -H "Content-Type: application/json" -d '{"jsonrpc": "1.0", "method": "importpubkey", "params": ["0232858a5faa413101831afe7a880da9a8ac4de6bd5e25b4358d762ba450b03c22", "", false], "id":"importpubkey"}' \
  && curl -s "http://127.0.0.1:5555/wallet/rsk-wallet" --user "$BTC_USERNAME:$BTC_PASSWORD" -H "Content-Type: application/json" -d '{"jsonrpc": "1.0", "method": "importaddress", "params": ["n1jGDaxCW6jemLZyd9wmDHddseZwEMV9C6", "", false], "id":"importaddress"}' \
  && curl -s "http://127.0.0.1:5555/wallet/main" --user "$BTC_USERNAME:$BTC_PASSWORD" -H "Content-Type: application/json" -d '{"jsonrpc": "1.0", "method": "walletpassphrase", "params": ["test-password", 30000], "id":"walletpassphrase"}' \
  && curl -s "http://127.0.0.1:5555/wallet/main" --user "$BTC_USERNAME:$BTC_PASSWORD" -H "Content-Type: application/json" -d '{"jsonrpc": "1.0", "method": "sendtoaddress", "params": { "amount": 5, "fee_rate": 25, "address": "n1jGDaxCW6jemLZyd9wmDHddseZwEMV9C6" }, "id":"sendtoaddress"}' \
  && curl -s "http://127.0.0.1:5555/wallet/main" --user "$BTC_USERNAME:$BTC_PASSWORD" -H "Content-Type: application/json" -d '{"jsonrpc": "1.0", "method": "getnewaddress", "params": ["main"], "id":"getnewaddress"}' \
    | jq .result | xargs -I ADDRESS curl -s "http://127.0.0.1:5555" --user "$BTC_USERNAME:$BTC_PASSWORD" -H "Content-Type: application/json" -d '{"jsonrpc": "1.0", "method": "generatetoaddress", "params": [1, "ADDRESS"], "id":"generatetoaddress"}'

if [ "$LPS_STAGE" = "regtest" ]; then
  PROVIDER_TX_COUNT=$(curl -s -X POST "http://127.0.0.1:4444" -H "Content-Type: application/json" -d "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getTransactionCount\",\"params\": [\"$LIQUIDITY_PROVIDER_RSK_ADDR\",\"latest\"],\"id\":1}" | jq -r ".result")
  if [ "$PROVIDER_TX_COUNT" = "0x0" ]; then
    echo "Transferring funds to $LIQUIDITY_PROVIDER_RSK_ADDR..."

    TX_HASH=$(curl -s -X POST "http://127.0.0.1:4444" -H "Content-Type: application/json" -d "{\"jsonrpc\":\"2.0\",\"method\":\"eth_sendTransaction\",\"params\": [{\"from\": \"0xcd2a3d9f938e13cd947ec05abc7fe734df8dd826\", \"to\": \"$LIQUIDITY_PROVIDER_RSK_ADDR\", \"value\": \"0x8AC7230489E80000\"}],\"id\":1}" | jq -r ".result")
    echo "Result: $TX_HASH"
    sleep 10
  else
    echo "No need to fund the '$LIQUIDITY_PROVIDER_RSK_ADDR' provider. Nonce: $PROVIDER_TX_COUNT"
  fi

  if [ -z "${LBC_ADDR}" ]; then
    echo "LBC_ADDR is not set. Deploying LBC contract..."

    (grep GITHUB_TOKEN | head -n 1 | tr -d '\r' | awk '{gsub("GITHUB_TOKEN=",""); print}' > gh_token.txt) < $ENV_FILE
    # deploy LBC contracts to RSKJ
    LBC_ADDR_LINE=$(docker compose --env-file "$ENV_FILE" -f docker-compose.yml -f docker-compose.lbc-deployer.yml run --rm lbc-deployer bash deploy-lbc.sh | grep LBC_ADDR | head -n 1 | tr -d '\r')
    export LBC_ADDR="${LBC_ADDR_LINE#"LBC_ADDR="}"
  fi
fi

if [ -z "${LBC_ADDR}" ]; then
  docker compose down
  echo "LBC_ADDR is not set up. Exit"
  exit 1
fi
echo "LBC deployed at $LBC_ADDR"

docker compose --env-file "$ENV_FILE" up -d powpeg-pegin powpeg-pegout
# start LPS

docker compose --env-file "$ENV_FILE" -f docker-compose.yml -f docker-compose.lps.yml build lps
docker compose --env-file "$ENV_FILE" -f docker-compose.yml -f docker-compose.lps.yml up -d lps

FAIL=true
for _ in $(seq 1 10);
do
  sleep 5
  curl -s "http://localhost:8080/health" \
    && echo "LPS is up and running" \
    && FAIL=false
  if [ "$FAIL" = false ]; then
    break
  fi
done

if [ "$FAIL" = true ]; then
  echo "LPS failed to start"
  exit 1
fi

rm -f cookie_jar.txt

PASSWORD_FILE_PATH="/tmp/management_password.txt"

echo "Checking for management_password.txt..."
if ! docker exec lps01 test -f "$PASSWORD_FILE_PATH"; then
  echo "management_password.txt not found. Skipping configuration steps"
  exit 0
fi

echo "management_password.txt found. Proceeding with configuration."

MANAGEMENT_PWD=$(docker exec lps01 cat "$PASSWORD_FILE_PATH")

CSRF_TOKEN=$(curl -s -c cookie_jar.txt \
                      -H 'Accept: */*' \
                      -H 'Connection: keep-alive' \
                      -H 'Content-Type: application/json' \
                      -H 'Origin: http://localhost:8080' \
                      -H 'Sec-Fetch-Dest: empty' \
                      -H 'Sec-Fetch-Mode: cors' \
                      -H 'Sec-Fetch-Site: same-origin' \
  "http://localhost:8080/management" | sed -n 's/.*name="csrf"[^>]*value="\([^"]*\)".*/\1/p')

# shellcheck disable=SC2001
CSRF_TOKEN=$(echo "$CSRF_TOKEN" | sed 's/&#43;/+/g')
curl -s -b cookie_jar.txt -c cookie_jar.txt "http://localhost:8080/management/login" \
  -H "X-CSRF-Token: $CSRF_TOKEN" \
  -H 'Content-Type: application/json' \
  -H 'Accept: */*' \
  -H 'Connection: keep-alive' \
  -H 'Origin: http://localhost:8080' \
  -H 'Referer: http://localhost:8080/management' \
  -H 'Sec-Fetch-Dest: empty' \
  -H 'Sec-Fetch-Mode: cors' \
  -H 'Sec-Fetch-Site: same-origin' \
  --data "{
     \"username\": \"admin\",
     \"password\": \"$MANAGEMENT_PWD\"
  }" || { echo "Error: login to Management UI failed"; exit 1; }

echo "Setting up general regtest configuration"
curl -sfS -b cookie_jar.txt 'http://localhost:8080/configuration' \
  -H "X-CSRF-Token: $CSRF_TOKEN" \
  -H 'Content-Type: application/json' \
  -H 'Accept: */*' \
  -H 'Connection: keep-alive' \
  -H 'Origin: http://localhost:8080' \
  -H 'Referer: http://localhost:8080/management' \
  -H 'Sec-Fetch-Dest: empty' \
  -H 'Sec-Fetch-Mode: cors' \
  -H 'Sec-Fetch-Site: same-origin' \
  --data '{
      "configuration": {
          "rskConfirmations": {
              "100000000000000000": 4,
              "2000000000000000000": 20,
              "400000000000000000": 12,
              "4000000000000000000": 40,
              "8000000000000000000": 80
          },
          "btcConfirmations": {
              "100000000000000000": 2,
              "2000000000000000000": 10,
              "400000000000000000": 6,
              "4000000000000000000": 20,
              "8000000000000000000": 40
          },
          "publicLiquidityCheck": true
      }
  }' || { echo "Error in configuring general regtest configuration"; exit 1; }

echo "Setting up pegin regtest configuration"
curl -sfS -b cookie_jar.txt 'http://localhost:8080/pegin/configuration' \
  -H "X-CSRF-Token: $CSRF_TOKEN" \
  -H 'Content-Type: application/json' \
  -H 'Accept: */*' \
  -H 'Connection: keep-alive' \
  -H 'Origin: http://localhost:8080' \
  -H 'Referer: http://localhost:8080/management' \
  -H 'Sec-Fetch-Dest: empty' \
  -H 'Sec-Fetch-Mode: cors' \
  -H 'Sec-Fetch-Site: same-origin' \
  --data '{
      "configuration": {
          "timeForDeposit": 3600,
          "callTime": 7200,
          "penaltyFee": "1000000000000000",
          "callFee": "10000000000000000",
          "maxValue": "10000000000000000000",
          "minValue": "600000000000000000"
      }
  }' || { echo "Error in configuring pegin regtest configuration"; exit 1; }

echo "Setting up pegout regtest configuration"
curl -sfS -b cookie_jar.txt 'http://localhost:8080/pegout/configuration' \
  -H "X-CSRF-Token: $CSRF_TOKEN" \
  -H 'Content-Type: application/json' \
  -H 'Accept: */*' \
  -H 'Connection: keep-alive' \
  -H 'Origin: http://localhost:8080' \
  -H 'Referer: http://localhost:8080/management' \
  -H 'Sec-Fetch-Dest: empty' \
  -H 'Sec-Fetch-Mode: cors' \
  -H 'Sec-Fetch-Site: same-origin' \
  --data '{
      "configuration": {
          "timeForDeposit": 3600,
          "expireTime": 10800,
          "penaltyFee": "1000000000000000",
          "callFee": "10000000000000000",
          "maxValue": "10000000000000000000",
          "minValue": "600000000000000000",
          "expireBlocks": 500,
          "bridgeTransactionMin": "1500000000000000000"
      }
  }' || { echo "Error in configuring pegout regtest configuration"; exit 1; }
